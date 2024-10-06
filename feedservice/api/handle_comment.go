package api

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleComment(w http.ResponseWriter, r *http.Request) {
	
	userId := r.Context().Value("user_id").(string)

	var commentReq CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not decode comment request"})
		return
	}

	err := s.postAccessor.CommentPost(userId, commentReq.PostId, commentReq.OwnerId, commentReq.Text)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not comment post"})
		return
	}
}