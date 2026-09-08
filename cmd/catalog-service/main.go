package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogv1 "github.com/andreysavvinov/petproject/gen/catalog/v1"
	"github.com/andreysavvinov/petproject/internal/db"
	"github.com/andreysavvinov/petproject/internal/events"
	"github.com/andreysavvinov/petproject/internal/grpcserver"
)

type server struct {
	catalogv1.UnimplementedCatalogServiceServer
	pool *pgxpool.Pool
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	databaseURL := env("DATABASE_URL", "postgres://music:music@localhost:5432/music?sslmode=disable")
	rabbitURL := env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	grpcAddr := env("GRPC_ADDR", ":50052")

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	consumer := events.NewConsumer(rabbitURL, func(ctx context.Context, event events.UserDeleted) error {
		log.Printf("catalog: hiding tracks for library %d (user %s)", event.MyMusicID, event.UserID)
		if event.MyMusicID == 0 {
			return nil
		}
		_, err := pool.Exec(ctx, `UPDATE track_list SET hidden = TRUE WHERE my_music_my_music_id = $1`, event.MyMusicID)
		return err
	})
	go func() {
		if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	s := grpcserver.New()
	catalogv1.RegisterCatalogServiceServer(s, &server{pool: pool})

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("catalog-service listening on %s", grpcAddr)
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

func (s *server) CreateAlbum(ctx context.Context, req *catalogv1.CreateAlbumRequest) (*catalogv1.Album, error) {
	id, err := db.NextID(ctx, s.pool, "albums")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "id: %v", err)
	}
	_, err = s.pool.Exec(ctx,
		`INSERT INTO albums (id, album_name, release_year, genre) VALUES ($1,$2,$3,$4)`,
		id, req.AlbumName, req.ReleaseYear, req.Genre)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &catalogv1.Album{Id: id, AlbumName: req.AlbumName, ReleaseYear: req.ReleaseYear, Genre: req.Genre}, nil
}

func (s *server) GetAlbum(ctx context.Context, req *catalogv1.GetAlbumRequest) (*catalogv1.Album, error) {
	var a catalogv1.Album
	err := s.pool.QueryRow(ctx, `SELECT id, album_name, release_year, genre FROM albums WHERE id = $1`, req.Id).
		Scan(&a.Id, &a.AlbumName, &a.ReleaseYear, &a.Genre)
	if err != nil {
		return nil, status.Error(codes.NotFound, "album not found")
	}
	return &a, nil
}

func (s *server) ListAlbums(ctx context.Context, _ *catalogv1.ListAlbumsRequest) (*catalogv1.ListAlbumsResponse, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, album_name, release_year, genre FROM albums ORDER BY id`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query: %v", err)
	}
	defer rows.Close()

	var albums []*catalogv1.Album
	for rows.Next() {
		var a catalogv1.Album
		if err := rows.Scan(&a.Id, &a.AlbumName, &a.ReleaseYear, &a.Genre); err != nil {
			return nil, status.Errorf(codes.Internal, "scan: %v", err)
		}
		albums = append(albums, &a)
	}
	return &catalogv1.ListAlbumsResponse{Albums: albums}, nil
}

func (s *server) UpdateAlbum(ctx context.Context, req *catalogv1.UpdateAlbumRequest) (*catalogv1.Album, error) {
	tag, err := s.pool.Exec(ctx,
		`UPDATE albums SET album_name=$2, release_year=$3, genre=$4 WHERE id=$1`,
		req.Id, req.AlbumName, req.ReleaseYear, req.Genre)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "album not found")
	}
	return &catalogv1.Album{Id: req.Id, AlbumName: req.AlbumName, ReleaseYear: req.ReleaseYear, Genre: req.Genre}, nil
}

func (s *server) DeleteAlbum(ctx context.Context, req *catalogv1.DeleteAlbumRequest) (*catalogv1.DeleteResponse, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM albums WHERE id = $1`, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "album not found")
	}
	return &catalogv1.DeleteResponse{Success: true}, nil
}

func (s *server) CreatePerformer(ctx context.Context, req *catalogv1.CreatePerformerRequest) (*catalogv1.Performer, error) {
	id, err := db.NextID(ctx, s.pool, "performers")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "id: %v", err)
	}
	_, err = s.pool.Exec(ctx, `INSERT INTO performers (id, performer_name) VALUES ($1,$2)`, id, req.PerformerName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &catalogv1.Performer{Id: id, PerformerName: req.PerformerName}, nil
}

func (s *server) GetPerformer(ctx context.Context, req *catalogv1.GetPerformerRequest) (*catalogv1.Performer, error) {
	var p catalogv1.Performer
	err := s.pool.QueryRow(ctx, `SELECT id, performer_name FROM performers WHERE id = $1`, req.Id).
		Scan(&p.Id, &p.PerformerName)
	if err != nil {
		return nil, status.Error(codes.NotFound, "performer not found")
	}
	return &p, nil
}

func (s *server) ListPerformers(ctx context.Context, _ *catalogv1.ListPerformersRequest) (*catalogv1.ListPerformersResponse, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, performer_name FROM performers ORDER BY id`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query: %v", err)
	}
	defer rows.Close()

	var performers []*catalogv1.Performer
	for rows.Next() {
		var p catalogv1.Performer
		if err := rows.Scan(&p.Id, &p.PerformerName); err != nil {
			return nil, status.Errorf(codes.Internal, "scan: %v", err)
		}
		performers = append(performers, &p)
	}
	return &catalogv1.ListPerformersResponse{Performers: performers}, nil
}

func (s *server) UpdatePerformer(ctx context.Context, req *catalogv1.UpdatePerformerRequest) (*catalogv1.Performer, error) {
	tag, err := s.pool.Exec(ctx, `UPDATE performers SET performer_name=$2 WHERE id=$1`, req.Id, req.PerformerName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "performer not found")
	}
	return &catalogv1.Performer{Id: req.Id, PerformerName: req.PerformerName}, nil
}

func (s *server) DeletePerformer(ctx context.Context, req *catalogv1.DeletePerformerRequest) (*catalogv1.DeleteResponse, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM performers WHERE id = $1`, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "performer not found")
	}
	return &catalogv1.DeleteResponse{Success: true}, nil
}

func (s *server) CreateTrack(ctx context.Context, req *catalogv1.CreateTrackRequest) (*catalogv1.Track, error) {
	id, err := db.NextID(ctx, s.pool, "track_list")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "id: %v", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO track_list (id, track_name, genre, duration, albums_id, performers_id, my_music_my_music_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		id, req.TrackName, req.Genre, req.Duration, req.AlbumsId, req.PerformersId, req.MyMusicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert: %v", err)
	}
	return &catalogv1.Track{
		Id: id, TrackName: req.TrackName, Genre: req.Genre, Duration: req.Duration,
		AlbumsId: req.AlbumsId, PerformersId: req.PerformersId, MyMusicId: req.MyMusicId,
	}, nil
}

func (s *server) GetTrack(ctx context.Context, req *catalogv1.GetTrackRequest) (*catalogv1.Track, error) {
	var t catalogv1.Track
	err := s.pool.QueryRow(ctx, `
		SELECT id, track_name, genre, duration, COALESCE(albums_id,''), COALESCE(performers_id,''), COALESCE(my_music_my_music_id,0)
		FROM track_list WHERE id = $1 AND hidden = FALSE`, req.Id).
		Scan(&t.Id, &t.TrackName, &t.Genre, &t.Duration, &t.AlbumsId, &t.PerformersId, &t.MyMusicId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "track not found")
	}
	return &t, nil
}

func (s *server) ListTracks(ctx context.Context, _ *catalogv1.ListTracksRequest) (*catalogv1.ListTracksResponse, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, track_name, genre, duration, COALESCE(albums_id,''), COALESCE(performers_id,''), COALESCE(my_music_my_music_id,0)
		FROM track_list WHERE hidden = FALSE ORDER BY id`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query: %v", err)
	}
	defer rows.Close()

	var tracks []*catalogv1.Track
	for rows.Next() {
		var t catalogv1.Track
		if err := rows.Scan(&t.Id, &t.TrackName, &t.Genre, &t.Duration, &t.AlbumsId, &t.PerformersId, &t.MyMusicId); err != nil {
			return nil, status.Errorf(codes.Internal, "scan: %v", err)
		}
		tracks = append(tracks, &t)
	}
	return &catalogv1.ListTracksResponse{Tracks: tracks}, nil
}

func (s *server) UpdateTrack(ctx context.Context, req *catalogv1.UpdateTrackRequest) (*catalogv1.Track, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE track_list SET track_name=$2, genre=$3, duration=$4, albums_id=$5, performers_id=$6, my_music_my_music_id=$7
		WHERE id=$1 AND hidden = FALSE`,
		req.Id, req.TrackName, req.Genre, req.Duration, req.AlbumsId, req.PerformersId, req.MyMusicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "track not found")
	}
	return &catalogv1.Track{
		Id: req.Id, TrackName: req.TrackName, Genre: req.Genre, Duration: req.Duration,
		AlbumsId: req.AlbumsId, PerformersId: req.PerformersId, MyMusicId: req.MyMusicId,
	}, nil
}

func (s *server) DeleteTrack(ctx context.Context, req *catalogv1.DeleteTrackRequest) (*catalogv1.DeleteResponse, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM track_list WHERE id = $1`, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete: %v", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "track not found")
	}
	return &catalogv1.DeleteResponse{Success: true}, nil
}
