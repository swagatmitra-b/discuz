package main

import (
	"context"
	"discuz/database/models"
	"net/http"
)

type ContextKey string

const contextKey = ContextKey("user")

func (s *APIServer) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized: Missing session token", http.StatusUnauthorized)
			return
		}

		sessionToken := cookie.Value
		session, err := s.db.getUserBySessionToken(sessionToken)
		if err != nil || session == (models.Session{}) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey, session.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *APIServer) UserContextMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		sessionToken := cookie.Value
		session, err := s.db.getUserBySessionToken(sessionToken)
		if err != nil || session == (models.Session{}) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey, session.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}