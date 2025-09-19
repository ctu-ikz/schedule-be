package api

import (
	"net/http"

	"github.com/ctu-ikz/schedule-be/internal/api/handler"
	"github.com/gorilla/mux"
)

func Router(authHandler *handler.AuthHandler) http.Handler {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.Handle("/health",
		authenticationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})),
	).Methods(http.MethodGet)

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", authHandler.Refresh).Methods(http.MethodPost)

	return r
}
