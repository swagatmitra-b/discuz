package main

import (
	"discuz/database/models"
	"discuz/utils"
	"fmt"
	"net/http"
	"time"
)

type PageData struct {
	Posts []models.Posts
	User  string
}

type XsrfCreate struct {
	User      string
	XsrfToken string
}

func (s *APIServer) home(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.getPosts()

	if err != nil {
		http.Error(w, "Error in Fetching Posts", 500)
		return
	}
	data := PageData{Posts: posts}

	user, ok := r.Context().Value(contextKey).(string)

	if !ok {
		data.User = ""
	} else {
		data.User = user
	}

	utils.RenderTemplate(w, "home", data)
}

func (s *APIServer) post(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	postPage, err := s.db.getPostTree(id)

	if err != nil {
		http.Error(w, "Invalid post", 500)
	}

	utils.RenderTemplate(w, "post", postPage)
}

func (s *APIServer) create(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(contextKey).(string)

	if !ok {
		http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
		return
	}

	xsrfCookie, err := r.Cookie("xsrf_token")
	if err != nil || xsrfCookie.Value == "" {
		http.Error(w, "Unauthorized: Please ensure an active session.", http.StatusUnauthorized)
		return
	}

	utils.RenderTemplate(w, "create", XsrfCreate{
		User:      user,
		XsrfToken: xsrfCookie.Value,
	})
}

func (s *APIServer) replyPost(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(contextKey).(string)

	if !ok {
		http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
		return
	}

	xsrfCookie, err := r.Cookie("xsrf_token")
	if err != nil || xsrfCookie.Value == "" {
		http.Error(w, "Unauthorized: Please ensure an active session", http.StatusUnauthorized)
		return
	}

	post_id := r.PathValue("post_id")
	post, _ := s.db.getPost(post_id)

	reply := Reply{
		Post_id:    post_id,
		Title:      &post.Title,
		Content:    post.Content,
		Created_at: post.Created_at,
		Created_by: post.Created_by,
		Replied_by: user,
		XsrfToken:  xsrfCookie.Value,
	}

	utils.RenderTemplate(w, "reply", reply)
}

func (s *APIServer) postReply(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	post_id := r.PostForm.Get("post_id")
	thread_id := r.PostForm.Get("thread_id")
	content := r.PostForm.Get("content")
	replied_by := r.PostForm.Get("replied_by")
	xsrf_token := r.PostForm.Get("xsrf_token")

	s.db.createThread(post_id, content, thread_id, replied_by)

	xsrfCookie, err := r.Cookie("xsrf_token")
	if err != nil || xsrfCookie.Value == "" {
		http.Error(w, "Unauthorized: Please ensure an active session", http.StatusUnauthorized)
		return
	}

	if xsrfCookie.Value != xsrf_token {
		http.Error(w, "Oops! Something went wrong. Please try again", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%s", post_id), http.StatusPermanentRedirect)
}

func (s *APIServer) replyThread(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(contextKey).(string)

	if !ok {
		http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
		return
	}

	post_id := r.PathValue("post_id")
	thread_id := r.PathValue("thread_id")

	parent, _ := s.db.getThread(thread_id)

	xsrfCookie, err := r.Cookie("xsrf_token")
	if err != nil || xsrfCookie.Value == "" {
		http.Error(w, "Unauthorized: Please ensure an active session", http.StatusUnauthorized)
		return
	}

	reply := Reply{
		Post_id:    post_id,
		Thread_id:  &thread_id,
		Content:    parent.Content,
		Created_at: parent.Created_at,
		Created_by: parent.Created_by,
		Replied_by: user,
		XsrfToken:  xsrfCookie.Value,
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
	creator := r.PostForm.Get("created_by")
	xsrf_token := r.PostForm.Get("xsrf_token")

	if len(title) == 0 || len(content) == 0 {
		http.Error(w, "Title or content cannot be empty", http.StatusUnauthorized)
		return
	}

	xsrfCookie, err := r.Cookie("xsrf_token")
	if err != nil || xsrfCookie.Value == "" {
		http.Error(w, "Unauthorized: Please ensure an active session", http.StatusUnauthorized)
		return
	}

	if xsrfCookie.Value != xsrf_token {
		http.Error(w, "Oops! Something went wrong. Please try again", http.StatusUnauthorized)
		return
	}

	s.db.createPost(title, content, creator)

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

	if len(password) == 0 {
		http.Error(w, "Password cannot be empty", http.StatusUnauthorized)
		return
	}

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

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token",
		Value:    csrfToken,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
	})

	err = s.db.createSession(sessionToken, user, time.Now().Add(24*time.Hour).Format(time.RFC3339))

	if err != nil {
		http.Error(w, "Failed to create session", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (s *APIServer) logout(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(contextKey).(string)

	if !ok {
		fmt.Println("no user in logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "xsrf_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	ok = s.db.deleteSession(user)

	if !ok {
		http.Error(w, "Could not log out successfully :/", 500)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
