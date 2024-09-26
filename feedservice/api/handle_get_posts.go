package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) handleGetPosts(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["user_id"]

	posts, err := s.postAccessor.GetPosts(userId, 10)
	if err != nil {
		slog.Error("could not get posts with userId: "+userId, "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get posts"})
		return
	}

	WriteJSON(w, http.StatusOK, posts)
}
