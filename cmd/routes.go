package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (s *APIServer) home (w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.getPosts()
	if err != nil {
		http.Error(w, "post error", 500)
	}
	tmpl, err := template.ParseFiles("cmd/static/home.html")
	if err != nil {
		http.Error(w, "Error in parsing html file", 500)
		return
	}
	if err := tmpl.Execute(w, posts); err != nil {
		panic(err)
	}
}

func(s *APIServer) post (w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		postPage, err := s.db.getPostTree(id)		

		if err != nil {
			http.Error(w, "invalid post", 500)
		}

		tmpl, err := template.ParseFiles("cmd/static/post.html")
		if err != nil {
			http.Error(w, "Error in parsing html file", 500)
			return
		}

		tmpl.Execute(w, postPage)
}

func(s *APIServer) create (w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("cmd/static/create.html")

		if err != nil {
			http.Error(w, "Error in parsing html file", 500)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
}

func (s *APIServer) replyPost (w http.ResponseWriter, r *http.Request) {
		post_id := r.PathValue("post_id")
		post, _ := s.db.getPost(post_id)

		reply := Reply{
			Post_id: post_id,
			Title: post.Title,
			Content: post.Content,
			Created_at: post.Created_at,
		}

		tmpl, err := template.ParseFiles("cmd/static/reply.html")

		if err != nil {
			http.Error(w, "Error in parsing html file", 500)
			return
		}

		tmpl.Execute(w, reply)
}

func (s *APIServer) postReply (w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()

		post_id :=  r.PostForm.Get("post_id")
		thread_id := r.PostForm.Get("thread_id")
		content := r.PostForm.Get("content")

		s.db.createThread(post_id, content, thread_id)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%s", post_id), http.StatusPermanentRedirect)
}

func (s *APIServer) replyThread (w http.ResponseWriter, r *http.Request) {

		post_id := r.PathValue("post_id")
		thread_id := r.PathValue("thread_id")

		parent, _ := s.db.getThread(thread_id)
				
		reply := Reply{
			Post_id: post_id,
			Thread_id: &thread_id,
			Content: parent.Content,
			Created_at: parent.Created_at,
		}

		tmpl, err := template.ParseFiles("cmd/static/reply.html")

		if err != nil {
			http.Error(w, "Error in parsing html file", 500)
			return
		}

		tmpl.Execute(w, reply)

	}

func (s *APIServer) postCreate (w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		title :=  r.PostForm.Get("title")
		body := r.PostForm.Get("body")

		s.db.createPost(title, body)

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}