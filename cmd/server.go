package main

import (
	"log"
	"net/http"
)

type APIServer struct {
	address string
	db      *dbDriver
	// session *sessions.CookieStore
}

type Reply struct {
	Post_id    string
	Thread_id  *string
	Title      string
	Content    string
	Created_at string
}

func createAPIServer(addr string, db *dbDriver) APIServer {
	// session := sessions.NewCookieStore([]byte("r80ds5cq9sfwahld2paksdjd07shwl"))
	// session.Options.HttpOnly = true
	// session.Options.SameSite = http.SameSiteLaxMode
	return APIServer{
		address: addr,
		db:      db,
		// session: session,
	}
}

func (s *APIServer) launch() error {
	router := http.NewServeMux()
	router.HandleFunc("/", s.home)
	router.HandleFunc("/post/{id}", s.post)
	router.HandleFunc("/create", s.create)
	router.HandleFunc("POST /create", s.postCreate)
	router.HandleFunc("/reply/post/{post_id}", s.replyPost)
	router.HandleFunc("/reply/thread/{post_id}/{thread_id}", s.replyThread)
	router.HandleFunc("POST /reply", s.postReply)
	router.HandleFunc("/user", s.createUser)
	router.HandleFunc("POST /user", s.postCreateUser)
	router.HandleFunc("/login", s.loginUser)
	router.HandleFunc("POST /login", s.login)

	server := http.Server{
		Addr:    s.address,
		Handler: router,
	}

	log.Printf("Server running on port %s", s.address)

	return server.ListenAndServe()
}
