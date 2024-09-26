package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DMK37/PassIt/feedservice/db"
	"github.com/DMK37/PassIt/feedservice/storage"
	"github.com/gorilla/mux"
)

type server struct {
	listenAddr   string
	postAccessor db.PostAccessor
	imageStorage storage.ImageStorage
}

func NewServer(listenAddr string) *server {

	postAccessor, err := db.NewDynamoDBPostAccessor()
	if err != nil {
		slog.Error("could not create post accessor", "error", err.Error())
		return nil
	}

	imageStorage, err := storage.NewS3ImageStorage()
	if err != nil {
		slog.Error("could not create image storage", "error", err.Error())
		return nil
	}
	return &server{
		listenAddr:   listenAddr,
		postAccessor: postAccessor,
		imageStorage: imageStorage,
	}
}

func (s *server) Start() {

	r := mux.NewRouter()

	r.HandleFunc("/users/{user_id}/posts/{post_id}", s.handleGetPost).Methods(http.MethodGet)
	r.HandleFunc("/users/{user_id}/posts", s.handleGetPosts).Methods(http.MethodGet)

	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(authMiddleware)

	privateRouter.HandleFunc("/posts", s.handleCreatePost).Methods(http.MethodPost)

	slog.Info("FeedService server starting", "addr", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		slog.Error("could not start server", "error", err.Error())
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}