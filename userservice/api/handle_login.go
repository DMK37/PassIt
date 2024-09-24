package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not decode login request"})
		return
	}

	user, err := s.userAccessor.GetUserByUsername(loginReq.Username)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get user"})
		return
	}

	if user == nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
		return
	}

	token, err := createToken(user.Id, user.Username)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
