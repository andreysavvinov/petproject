package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/andreysavvinov/petproject/gen/user/v1"
	"github.com/andreysavvinov/petproject/internal/db"
	"github.com/andreysavvinov/petproject/internal/events"
	"github.com/andreysavvinov/petproject/internal/grpcserver"
	"github.com/andreysavvinov/petproject/internal/jwtutil"
)

type server struct {
	userv1.UnimplementedUserServiceServer
	pool      *pgxpool.Pool
	jwtSecret string
	publisher *events.Publisher
}

func main() {
	ctx := context.Background()
	databaseURL := env("DATABASE_URL", "postgres://music:music@localhost:5432/music?sslmode=disable")
	rabbitURL := env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	grpcAddr := env("GRPC_ADDR", ":50051")
	jwtSecret := env("JWT_SECRET", "dev-secret-change-me")

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	publisher, err := events.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}

	s := grpcserver.New()
	userv1.RegisterUserServiceServer(s, &server{
		pool:      pool,
		jwtSecret: jwtSecret,
		publisher: publisher,
	})

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("user-service listening on %s", grpcAddr)
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

func (s *server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name, email and password are required")
	}

	id, err := db.NextID(ctx, s.pool, "users")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate id: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash password: %v", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`INSERT INTO users (id, name, email, password, phone_number, country) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, req.Name, req.Email, string(hash), req.PhoneNumber, req.Country,
	); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user exists: %v", err)
	}

	var myMusicID int64
	if err := tx.QueryRow(ctx,
		`INSERT INTO my_music (users_id) VALUES ($1) RETURNING my_music_id`, id,
	).Scan(&myMusicID); err != nil {
		return nil, status.Errorf(codes.Internal, "create library: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "commit: %v", err)
	}

	user := &userv1.User{
		Id:          id,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Country:     req.Country,
		MyMusicId:   myMusicID,
	}

	token, err := jwtutil.Generate(s.jwtSecret, id, req.Email, 24*time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token: %v", err)
	}

	return &userv1.AuthResponse{Token: token, User: user}, nil
}

func (s *server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.AuthResponse, error) {
	var user userv1.User
	var passwordHash string
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.name, u.email, u.password, u.phone_number, u.country, COALESCE(m.my_music_id, 0)
		FROM users u
		LEFT JOIN my_music m ON m.users_id = u.id
		WHERE u.email = $1`, req.Email,
	).Scan(&user.Id, &user.Name, &user.Email, &passwordHash, &user.PhoneNumber, &user.Country, &user.MyMusicId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := jwtutil.Generate(s.jwtSecret, user.Id, user.Email, 24*time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token: %v", err)
	}

	return &userv1.AuthResponse{Token: token, User: &user}, nil
}

func (s *server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error) {
	var user userv1.User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.name, u.email, u.phone_number, u.country, COALESCE(m.my_music_id, 0)
		FROM users u
		LEFT JOIN my_music m ON m.users_id = u.id
		WHERE u.id = $1`, req.Id,
	).Scan(&user.Id, &user.Name, &user.Email, &user.PhoneNumber, &user.Country, &user.MyMusicId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &user, nil
}

func (s *server) ListUsers(ctx context.Context, _ *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT u.id, u.name, u.email, u.phone_number, u.country, COALESCE(m.my_music_id, 0)
		FROM users u
		LEFT JOIN my_music m ON m.users_id = u.id
		ORDER BY u.id`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query: %v", err)
	}
	defer rows.Close()

	var users []*userv1.User
	for rows.Next() {
		var user userv1.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PhoneNumber, &user.Country, &user.MyMusicId); err != nil {
			return nil, status.Errorf(codes.Internal, "scan: %v", err)
		}
		users = append(users, &user)
	}
	return &userv1.ListUsersResponse{Users: users}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	var myMusicID int64
	err := s.pool.QueryRow(ctx, `SELECT my_music_id FROM my_music WHERE users_id = $1`, req.Id).Scan(&myMusicID)
	if err == pgx.ErrNoRows {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lookup library: %v", err)
	}

	if err := s.publisher.PublishUserDeleted(ctx, req.Id, myMusicID); err != nil {
		return nil, status.Errorf(codes.Internal, "publish event: %v", err)
	}

	tag, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	log.Printf("user deleted: %s, event published", req.Id)
	return &userv1.DeleteUserResponse{Success: true}, nil
}
