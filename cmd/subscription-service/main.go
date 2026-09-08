package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	subscriptionv1 "github.com/andreysavvinov/petproject/gen/subscription/v1"
	"github.com/andreysavvinov/petproject/internal/db"
	"github.com/andreysavvinov/petproject/internal/events"
	"github.com/andreysavvinov/petproject/internal/grpcserver"
)

type server struct {
	subscriptionv1.UnimplementedSubscriptionServiceServer
	pool *pgxpool.Pool
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	databaseURL := env("DATABASE_URL", "postgres://music:music@localhost:5432/music?sslmode=disable")
	rabbitURL := env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	grpcAddr := env("GRPC_ADDR", ":50053")

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	consumer := events.NewConsumer(rabbitURL, func(ctx context.Context, event events.UserDeleted) error {
		log.Printf("subscription: removing subscriptions for user %s", event.UserID)
		_, err := pool.Exec(ctx, `DELETE FROM subscription_validaty WHERE users_id = $1`, event.UserID)
		return err
	})
	go func() {
		if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	s := grpcserver.New()
	subscriptionv1.RegisterSubscriptionServiceServer(s, &server{pool: pool})

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("subscription-service listening on %s", grpcAddr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (s *server) CreateSubscription(ctx context.Context, req *subscriptionv1.CreateSubscriptionRequest) (*subscriptionv1.Subscription, error) {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO subscription_validaty (users_id, track_list_id, start_date, end_date)
		VALUES ($1,$2,$3::date,$4::date)`,
		req.UsersId, req.TrackListId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &subscriptionv1.Subscription{
		UsersId: req.UsersId, TrackListId: req.TrackListId,
		StartDate: req.StartDate, EndDate: req.EndDate,
	}, nil
}

func (s *server) GetSubscription(ctx context.Context, req *subscriptionv1.GetSubscriptionRequest) (*subscriptionv1.Subscription, error) {
	var sub subscriptionv1.Subscription
	err := s.pool.QueryRow(ctx, `
		SELECT users_id, track_list_id, start_date::text, end_date::text
		FROM subscription_validaty WHERE users_id = $1 AND track_list_id = $2`,
		req.UsersId, req.TrackListId,
	).Scan(&sub.UsersId, &sub.TrackListId, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, status.Error(codes.NotFound, "subscription not found")
	}
	return &sub, nil
}

func (s *server) ListSubscriptions(ctx context.Context, _ *subscriptionv1.ListSubscriptionsRequest) (*subscriptionv1.ListSubscriptionsResponse, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT users_id, track_list_id, start_date::text, end_date::text
		FROM subscription_validaty ORDER BY users_id, track_list_id`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query: %v", err)
	}
	defer rows.Close()

	var subs []*subscriptionv1.Subscription
	for rows.Next() {
		var sub subscriptionv1.Subscription
		if err := rows.Scan(&sub.UsersId, &sub.TrackListId, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, status.Errorf(codes.Internal, "scan: %v", err)
		}
		subs = append(subs, &sub)
	}
	return &subscriptionv1.ListSubscriptionsResponse{Subscriptions: subs}, nil
}

func (s *server) UpdateSubscription(ctx context.Context, req *subscriptionv1.UpdateSubscriptionRequest) (*subscriptionv1.Subscription, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE subscription_validaty SET start_date=$3::date, end_date=$4::date
		WHERE users_id=$1 AND track_list_id=$2`,
		req.UsersId, req.TrackListId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "subscription not found")
	}
	return &subscriptionv1.Subscription{
		UsersId: req.UsersId, TrackListId: req.TrackListId,
		StartDate: req.StartDate, EndDate: req.EndDate,
	}, nil
}

func (s *server) DeleteSubscription(ctx context.Context, req *subscriptionv1.DeleteSubscriptionRequest) (*subscriptionv1.DeleteResponse, error) {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM subscription_validaty WHERE users_id=$1 AND track_list_id=$2`,
		req.UsersId, req.TrackListId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "subscription not found")
	}
	return &subscriptionv1.DeleteResponse{Success: true}, nil
}
