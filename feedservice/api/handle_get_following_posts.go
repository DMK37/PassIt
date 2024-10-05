package api

import (
	"log/slog"
	"net/http"
)

func (s *server) handleGetFollowingPosts(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	posts, users, err := s.postAccessor.GetFollowingPosts(userId)
	if err != nil {
		slog.Error("could not get following posts", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get following posts"})
		return
	}

	responsePosts := make([]*ResponsePost, len(posts))

	for i, post := range posts {
		user := users[post.UserId]
		responsePosts[i] = mapPostToResponsePost(post, user)
	}

	WriteJSON(w, http.StatusOK, responsePosts)
}
