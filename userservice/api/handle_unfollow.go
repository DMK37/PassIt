package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) handleUnfollow(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	followUserId := mux.Vars(r)["user_id"]

	if err := s.userAccessor.UnfollowUser(userId, followUserId); err != nil {
		slog.Error("could not follow user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not follow user"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "user followed"})
}