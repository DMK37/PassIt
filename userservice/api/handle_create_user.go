package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DMK37/PassIt/userservice/db"
)

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("could not decode user", "error", err.Error())
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not decode user"})
		return
	}
	fmt.Println("Created user: ", user)
	createdUser := db.NewUser(user.Username, user.Email, user.Password, user.FirstName, user.LastName)

	if err := s.userAccessor.CreateUser(createdUser); err != nil {
		slog.Error("could not create user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	slog.Info("user created", "id", createdUser.Id)

	token, err := createToken(createdUser.Id, createdUser.Username)
	if err != nil {
		slog.Error("could not generate token", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"token": token, "user": createdUser.String()})
}
