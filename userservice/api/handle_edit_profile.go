package api

import (
	"log/slog"
	"net/http"
)

func (s *server) handleEditProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("could not parse multipart form", "error", err.Error())
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "could not parse multipart form"})
		return
	}

	// Extract form fields
	username := r.FormValue("username")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	path := ""
	_, handler, err := r.FormFile("image")
	if err == nil {
		path, err = s.storage.UploadImage(handler, userId)
		if err != nil {
			slog.Error("could not upload image", "error", err.Error())
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not upload image"})
			return
		}
	}

	if err := s.userAccessor.UpdateUser(userId, username, firstName, lastName, path); err != nil {
		slog.Error("could not edit profile", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not edit profile"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "profile edited"})
}
