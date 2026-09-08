package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	catalogv1 "github.com/andreysavvinov/petproject/gen/catalog/v1"
	subscriptionv1 "github.com/andreysavvinov/petproject/gen/subscription/v1"
	userv1 "github.com/andreysavvinov/petproject/gen/user/v1"
	"github.com/andreysavvinov/petproject/internal/grpcjson"
	"github.com/andreysavvinov/petproject/internal/jwtutil"
)

type gateway struct {
	users         userv1.UserServiceClient
	catalog       catalogv1.CatalogServiceClient
	subscriptions subscriptionv1.SubscriptionServiceClient
	jwtSecret     string
}

func main() {
	userAddr := env("USER_SERVICE_ADDR", "localhost:50051")
	catalogAddr := env("CATALOG_SERVICE_ADDR", "localhost:50052")
	subAddr := env("SUBSCRIPTION_SERVICE_ADDR", "localhost:50053")
	httpAddr := env("HTTP_ADDR", ":8080")
	jwtSecret := env("JWT_SECRET", "dev-secret-change-me")

	userConn := dial(userAddr)
	catalogConn := dial(catalogAddr)
	subConn := dial(subAddr)
	defer userConn.Close()
	defer catalogConn.Close()
	defer subConn.Close()

	gw := &gateway{
		users:         userv1.NewUserServiceClient(userConn),
		catalog:       catalogv1.NewCatalogServiceClient(catalogConn),
		subscriptions: subscriptionv1.NewSubscriptionServiceClient(subConn),
		jwtSecret:     jwtSecret,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth/register", gw.handleRegister)
	mux.HandleFunc("POST /api/auth/login", gw.handleLogin)
	mux.HandleFunc("GET /api/users", gw.auth(gw.handleListUsers))
	mux.HandleFunc("GET /api/users/{id}", gw.auth(gw.handleGetUser))
	mux.HandleFunc("DELETE /api/users/{id}", gw.auth(gw.handleDeleteUser))

	mux.HandleFunc("GET /api/albums", gw.auth(gw.handleListAlbums))
	mux.HandleFunc("POST /api/albums", gw.auth(gw.handleCreateAlbum))
	mux.HandleFunc("GET /api/albums/{id}", gw.auth(gw.handleGetAlbum))
	mux.HandleFunc("PUT /api/albums/{id}", gw.auth(gw.handleUpdateAlbum))
	mux.HandleFunc("DELETE /api/albums/{id}", gw.auth(gw.handleDeleteAlbum))

	mux.HandleFunc("GET /api/performers", gw.auth(gw.handleListPerformers))
	mux.HandleFunc("POST /api/performers", gw.auth(gw.handleCreatePerformer))
	mux.HandleFunc("GET /api/performers/{id}", gw.auth(gw.handleGetPerformer))
	mux.HandleFunc("PUT /api/performers/{id}", gw.auth(gw.handleUpdatePerformer))
	mux.HandleFunc("DELETE /api/performers/{id}", gw.auth(gw.handleDeletePerformer))

	mux.HandleFunc("GET /api/tracks", gw.auth(gw.handleListTracks))
	mux.HandleFunc("POST /api/tracks", gw.auth(gw.handleCreateTrack))
	mux.HandleFunc("GET /api/tracks/{id}", gw.auth(gw.handleGetTrack))
	mux.HandleFunc("PUT /api/tracks/{id}", gw.auth(gw.handleUpdateTrack))
	mux.HandleFunc("DELETE /api/tracks/{id}", gw.auth(gw.handleDeleteTrack))

	mux.HandleFunc("GET /api/subscriptions", gw.auth(gw.handleListSubscriptions))
	mux.HandleFunc("POST /api/subscriptions", gw.auth(gw.handleCreateSubscription))
	mux.HandleFunc("GET /api/subscriptions/{userId}/{trackId}", gw.auth(gw.handleGetSubscription))
	mux.HandleFunc("PUT /api/subscriptions/{userId}/{trackId}", gw.auth(gw.handleUpdateSubscription))
	mux.HandleFunc("DELETE /api/subscriptions/{userId}/{trackId}", gw.auth(gw.handleDeleteSubscription))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Printf("gateway listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("http: %v", err)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(grpcjson.Name)),
	)
	if err != nil {
		log.Fatalf("dial %s: %v", addr, err)
	}
	return conn
}

func (gw *gateway) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		if _, err := jwtutil.Parse(gw.jwtSecret, strings.TrimPrefix(header, "Bearer ")); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		next(w, r)
	}
}

func (gw *gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req userv1.RegisterRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.users.Register(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req userv1.LoginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.users.Login(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleListUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.users.ListUsers(r.Context(), &userv1.ListUsersRequest{})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleGetUser(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.users.GetUser(r.Context(), &userv1.GetUserRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.users.DeleteUser(r.Context(), &userv1.DeleteUserRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleListAlbums(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.ListAlbums(r.Context(), &catalogv1.ListAlbumsRequest{})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleCreateAlbum(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.CreateAlbumRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.catalog.CreateAlbum(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleGetAlbum(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.GetAlbum(r.Context(), &catalogv1.GetAlbumRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleUpdateAlbum(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.UpdateAlbumRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Id = r.PathValue("id")
	resp, err := gw.catalog.UpdateAlbum(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleDeleteAlbum(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.DeleteAlbum(r.Context(), &catalogv1.DeleteAlbumRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleListPerformers(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.ListPerformers(r.Context(), &catalogv1.ListPerformersRequest{})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleCreatePerformer(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.CreatePerformerRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.catalog.CreatePerformer(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleGetPerformer(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.GetPerformer(r.Context(), &catalogv1.GetPerformerRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleUpdatePerformer(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.UpdatePerformerRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Id = r.PathValue("id")
	resp, err := gw.catalog.UpdatePerformer(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleDeletePerformer(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.DeletePerformer(r.Context(), &catalogv1.DeletePerformerRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleListTracks(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.ListTracks(r.Context(), &catalogv1.ListTracksRequest{})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleCreateTrack(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.CreateTrackRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.catalog.CreateTrack(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleGetTrack(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.GetTrack(r.Context(), &catalogv1.GetTrackRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleUpdateTrack(w http.ResponseWriter, r *http.Request) {
	var req catalogv1.UpdateTrackRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Id = r.PathValue("id")
	resp, err := gw.catalog.UpdateTrack(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleDeleteTrack(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.catalog.DeleteTrack(r.Context(), &catalogv1.DeleteTrackRequest{Id: r.PathValue("id")})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.subscriptions.ListSubscriptions(r.Context(), &subscriptionv1.ListSubscriptionsRequest{})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionv1.CreateSubscriptionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	resp, err := gw.subscriptions.CreateSubscription(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.subscriptions.GetSubscription(r.Context(), &subscriptionv1.GetSubscriptionRequest{
		UsersId: r.PathValue("userId"), TrackListId: r.PathValue("trackId"),
	})
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionv1.UpdateSubscriptionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.UsersId = r.PathValue("userId")
	req.TrackListId = r.PathValue("trackId")
	resp, err := gw.subscriptions.UpdateSubscription(r.Context(), &req)
	writeGRPC(w, resp, err)
}

func (gw *gateway) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	resp, err := gw.subscriptions.DeleteSubscription(r.Context(), &subscriptionv1.DeleteSubscriptionRequest{
		UsersId: r.PathValue("userId"), TrackListId: r.PathValue("trackId"),
	})
	writeGRPC(w, resp, err)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return false
	}
	return true
}

func writeGRPC(w http.ResponseWriter, resp any, err error) {
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
