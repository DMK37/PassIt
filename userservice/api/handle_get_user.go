package api

import (
	"log/slog"
	"net/http"
)

func (s *server) handleGetUser(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("user_id").(string)
    if !ok {
        slog.Error("user_id not found in context")
        WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
        return
    }

	user, err := s.userAccessor.GetUserById(userId)
	if err != nil {
		slog.Error("could not get user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get user"})
		return
	}

	if user == nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	WriteJSON(w, http.StatusOK, user)
}