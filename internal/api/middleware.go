package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ctu-ikz/schedule-be/internal/util"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		log.Printf("%s %s %d %s [%s]",
			r.Method,
			r.URL.Path,
			ww.statusCode,
			r.RemoteAddr,
			duration,
		)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JSON(w, http.StatusUnauthorized, map[string]string{"status": "unauthorized"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			JSON(w, http.StatusUnauthorized, map[string]string{"status": "unauthorized"})
			return
		}

		token := parts[1]
		_, err := util.ParseAccessToken(token)
		if err != nil {
			JSON(w, http.StatusUnauthorized, map[string]string{"status": "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
