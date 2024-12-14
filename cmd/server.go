package main

import (
	"log"
	"net/http"
)

type APIServer struct {
	address string
	db      *dbDriver
}

type Reply struct {
	Post_id    string
	Thread_id  *string
	Title      *string
	Content    string
	Created_at string
	Created_by string
	Replied_by string
	XsrfToken  string
}

func createAPIServer(addr string, db *dbDriver) APIServer {
	return APIServer{
		address: addr,
		db:      db,
	}
}

func (s *APIServer) launch() error {
	router := http.NewServeMux()
	router.HandleFunc("/", s.UserContextMiddleware(s.home))
	router.HandleFunc("/post/{id}", s.UserContextMiddleware(s.post))
	router.HandleFunc("/create", s.UserContextMiddleware(s.AuthMiddleware(s.create)))
	router.HandleFunc("POST /create", s.UserContextMiddleware(s.AuthMiddleware(s.postCreate)))
	router.HandleFunc("/reply/post/{post_id}", s.UserContextMiddleware(s.AuthMiddleware(s.replyPost)))
	router.HandleFunc("/reply/thread/{post_id}/{thread_id}", s.UserContextMiddleware(s.AuthMiddleware(s.replyThread)))
	router.HandleFunc("POST /reply", s.UserContextMiddleware(s.AuthMiddleware(s.postReply)))
	router.HandleFunc("/user", s.UserContextMiddleware(s.createUser))
	router.HandleFunc("POST /user", s.UserContextMiddleware(s.postCreateUser))
	router.HandleFunc("/login", s.UserContextMiddleware(s.loginUser))
	router.HandleFunc("POST /login", s.UserContextMiddleware(s.login))
	router.HandleFunc("/logout", s.UserContextMiddleware(s.logout))

	server := http.Server{
		Addr:    s.address,
		Handler: router,
	}

	log.Printf("Server running on port %s", s.address)

	return server.ListenAndServe()
}
