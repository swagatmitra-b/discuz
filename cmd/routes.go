package main

import (
	"discuz/utils"
	"fmt"
	"net/http"
)

func (s *APIServer) home(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.getPosts()
	if err != nil {
		http.Error(w, "post error", 500)
	}

	utils.RenderTemplate(w, "home", posts)
}

func (s *APIServer) post(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	postPage, err := s.db.getPostTree(id)

	if err != nil {
		http.Error(w, "invalid post", 500)
	}

	utils.RenderTemplate(w, "post", postPage)
}

func (s *APIServer) create(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "create", nil)
}

func (s *APIServer) replyPost(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	post, _ := s.db.getPost(post_id)

	reply := Reply{
		Post_id:    post_id,
		Title:      post.Title,
		Content:    post.Content,
		Created_at: post.Created_at,
	}

	utils.RenderTemplate(w, "reply", reply)
}

func (s *APIServer) postReply(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	post_id := r.PostForm.Get("post_id")
	thread_id := r.PostForm.Get("thread_id")
	content := r.PostForm.Get("content")

	s.db.createThread(post_id, content, thread_id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%s", post_id), http.StatusPermanentRedirect)
}

func (s *APIServer) replyThread(w http.ResponseWriter, r *http.Request) {

	post_id := r.PathValue("post_id")
	thread_id := r.PathValue("thread_id")

	parent, _ := s.db.getThread(thread_id)

	reply := Reply{
		Post_id:    post_id,
		Thread_id:  &thread_id,
		Content:    parent.Content,
		Created_at: parent.Created_at,
	}

	utils.RenderTemplate(w, "reply", reply)
}

func (s *APIServer) postCreate(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	s.db.createPost(title, content)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (s *APIServer) createUser(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "user", nil)
}

func (s *APIServer) postCreateUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	s.db.createUser(username, password)

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func (s *APIServer) loginUser(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login", nil)
}

func (s *APIServer) login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	user, err := s.db.authUser(username, password)

	if err != nil {
		http.Error(w, "Authentication Error", 500)
		return
	}

	fmt.Println(user)

	// session, _ := s.session.Get(r, "auth")
	// session.Values["user"] = user
	// err = session.Save(r, w)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
