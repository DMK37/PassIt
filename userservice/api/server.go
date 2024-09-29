package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DMK37/PassIt/userservice/db"
	"github.com/DMK37/PassIt/userservice/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type server struct {
	listenAddr   string
	userAccessor db.UserAccessor
	storage      storage.ImageStorage
}

func NewServer(listenAddr string) *server {

	userAccessor, err := db.NewDynamoDBUserAccessor()
	if err != nil {
		slog.Error("could not create user accessor", "error", err.Error())
		return nil
	}

	storage, err := storage.NewS3ImageStorage()
	if err != nil {
		slog.Error("could not create image storage", "error", err.Error())
		return nil
	}
	return &server{
		listenAddr:   listenAddr,
		userAccessor: userAccessor,
		storage:      storage,
	}
}

func (s *server) Start() {

	r := mux.NewRouter()

	r.HandleFunc("/users", s.handleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/users/{username}", s.handleGetUserProfile).Methods(http.MethodGet)

	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(authMiddleware)

	privateRouter.HandleFunc("/users", s.handleGetUser).Methods(http.MethodGet)
	privateRouter.HandleFunc("/follow/{user_id}", s.handleFollow).Methods(http.MethodPost)
	privateRouter.HandleFunc("/unfollow/{user_id}", s.handleUnfollow).Methods(http.MethodPost)
	privateRouter.HandleFunc("/profile/edit", s.handleEditProfile).Methods(http.MethodPut)

	slog.Info("UserService server starting", "addr", s.listenAddr)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Update with your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	if err := http.ListenAndServe(s.listenAddr, handler); err != nil {
		slog.Error("could not start server", "error", err.Error())
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
