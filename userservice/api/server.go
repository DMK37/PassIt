package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"userservice/db"

	"github.com/gorilla/mux"
)

type server struct {
	listenAddr   string
	userAccessor db.UserAccessor
}

func NewServer(listenAddr string) *server {

	userAccessor, err := db.NewDynamoDBUserAccessor()
	if err != nil {
		slog.Error("could not create user accessor", "error", err.Error())
		return nil
	}
	return &server{
		listenAddr:   listenAddr,
		userAccessor: userAccessor,
	}
}

func (s *server) Start() {

	r := mux.NewRouter()

	r.HandleFunc("/users", s.handleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)

	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(authMiddleware)

	privateRouter.HandleFunc("/users", s.handleGetUser).Methods(http.MethodGet)

	slog.Info("UserService server starting", "addr", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		slog.Error("could not start server", "error", err.Error())
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
