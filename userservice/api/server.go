package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"userservice/internal/models"

	"github.com/gorilla/mux"
)

type Server struct {
	
}

func (s *Server) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.createUser).Methods(http.MethodPost)

	if err := http.ListenAndServe(":8080", nil); err != nil {
        slog.Error("could not start server", "error", err.Error())
	}
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        slog.Error("could not decode user", "error", err.Error())
		http.Error(w, "could not decode user", http.StatusBadRequest)
		return
	}

	// save user
	fmt.Println("user", user)
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}