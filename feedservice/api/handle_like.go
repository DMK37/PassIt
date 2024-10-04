package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *server) handleLike(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("user_id").(string)

	var likeReq LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&likeReq); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not decode like request"})
		return
	}

	err := s.postAccessor.LikePost(userId, likeReq.PostId, likeReq.OwnerId)

	if err != nil {
		slog.Error("could not like post with postId: "+likeReq.PostId, "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not like post"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "post liked"})
}
