package main

// import (
// 	"context"
// 	"net/http"
// )

// type ContextKey string

// const contextKey = ContextKey("user")

// func (s *APIServer) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		session, _ := s.session.Get(r, "auth")
// 		userId, exists := session.Values["user"].(string)
// 		http.SetCookie(w, &http.Cookie{
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		if !exists {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		user, err := s.db.getUser(userId)

// 		if err != nil {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), contextKey, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
