package api

import "net/http"

func healthHandler(w http.ResponseWriter, r *http.Request) {
    WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}