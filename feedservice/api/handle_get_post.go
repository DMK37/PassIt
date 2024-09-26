package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) handleGetPost(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["user_id"]
	postId := mux.Vars(r)["post_id"]

	post, err := s.postAccessor.GetPost(userId, postId)
	if err != nil {
		slog.Error("could not get post with userId: "+userId+" and postId: "+postId, "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get post"})
		return
	}

	if post == nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "post not found"})
		return
	}

	WriteJSON(w, http.StatusOK, post)
}
