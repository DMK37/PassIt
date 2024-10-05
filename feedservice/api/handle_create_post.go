package api

import (
	"log/slog"
	"net/http"

	"github.com/DMK37/PassIt/feedservice/db"
)

func (s *server) handleCreatePost(w http.ResponseWriter, r *http.Request) {

	text := r.FormValue("text")
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		slog.Error("user_id not found in context")
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		slog.Error("form exceeded max size", "error", err.Error())
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "form exceeded max size"})
		return
	}

	images := r.MultipartForm.File["images"]

	// Upload images to S3 and get URLs
	imageURLs := make([]string, len(images))

	for i, image := range images {
		url, err := s.imageStorage.UploadImage(image, userId)
		if err != nil {
			slog.Error("could not upload image", "error", err.Error())
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not upload image"})
			return
		}
		imageURLs[i] = url
	}

	post := db.NewPost(userId, text, imageURLs)

	if err := s.postAccessor.CreatePost(post); err != nil {
		slog.Error("could not save post", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create post"})
		return
	}

	user, err := s.postAccessor.GetPostUser(userId)
	if err != nil {
		slog.Error("could not get user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get user"})
		return
	}

	responsePost := mapPostToResponsePost(post, user)

	WriteJSON(w, http.StatusOK, responsePost)
}
