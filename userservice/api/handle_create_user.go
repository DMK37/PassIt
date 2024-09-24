package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"userservice/db"

	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("could not decode user", "error", err.Error())
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not decode user"})
		return
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        slog.Error("could not hash password", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
    }

	hashedPassword := string(hashedPasswordBytes)

	createdUser := db.NewUser(user.Username, hashedPassword, user.FirstName, user.LastName)

	if err := s.userAccessor.CreateUser(createdUser); err != nil {
		slog.Error("could not create user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	slog.Info("user created", "id", createdUser.Id)

	WriteJSON(w, http.StatusCreated, createdUser)
}
