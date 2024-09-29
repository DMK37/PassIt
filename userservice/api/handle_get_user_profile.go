package api

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {

	username := mux.Vars(r)["username"]

	user, err := s.userAccessor.GetUserByUsername(username)
	if err != nil {
		slog.Error("could not get user", "error", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get user"})
		return
	}
	if user == nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	userProfile := &UserProfile{
		Id:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Followers: user.Followers,
		Following: user.Following,
		Avatar:    user.Avatar,
	}

	WriteJSON(w, http.StatusOK, userProfile)
}
