package main

import (
	"context"
	"database/sql"
	"discuz/database/models"
	"discuz/utils"
	"fmt"
	"log"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"

	"golang.org/x/crypto/bcrypt"
)

type dbDriver struct {
	db *sql.DB
}

func connectDB() (*dbDriver, error) {
	driver, err := sql.Open("sqlite", "./app.db")

	if err != nil {
		log.Fatal("Failed to connect", err)
	}

	err = driver.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Println("Database connected successfully!")

	return &dbDriver{
		db: driver,
	}, nil
}

func (driver *dbDriver) getPosts() ([]models.Posts, error) {

	rows, err := driver.db.QueryContext(context.Background(), `SELECT * FROM Posts`)

	if err != nil {
		return nil, err
	}

	posts := []models.Posts{}

	for rows.Next() {
		post := models.Posts{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at, &post.Created_by)
		utils.ParseDateString(&post.Created_at)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (driver *dbDriver) getPostTree(id string) (models.PostPage, error) {

	postPage := models.PostPage{
		Post: models.Posts{ID: id},
	}

	row := driver.db.QueryRowContext(context.Background(), `SELECT title, content, created_at, created_by from Posts WHERE ID = ?`, id)

	threadRows, err := driver.db.QueryContext(context.Background(), `SELECT ID, parent_id, content, created_at, created_by from Threads WHERE post_id = ?`, id)

	for threadRows.Next() {
		thread := models.Threads{Post_id: id}
		threadRows.Scan(&thread.ID, &thread.Parent_id, &thread.Content, &thread.Created_at, &thread.Created_by)
		utils.ParseDateString(&thread.Created_at)
		postPage.Threads = append(postPage.Threads, &thread)
	}

	tree := utils.BuildTree(postPage.Threads)
	postPage.Threads = tree

	if err != nil {
		return postPage, err
	}

	err = row.Scan(&postPage.Post.Title, &postPage.Post.Content, &postPage.Post.Created_at, &postPage.Post.Created_by)

	if err != nil {
		return postPage, err
	}

	utils.ParseDateString(&postPage.Post.Created_at)

	return postPage, nil
}

func (driver *dbDriver) getThread(id string) (models.Threads, error) {
	thread := models.Threads{}

	row := driver.db.QueryRowContext(context.Background(), `SELECT content, created_at, created_by from Threads WHERE ID = ?`, id)

	row.Scan(&thread.Content, &thread.Created_at, &thread.Created_by)

	utils.ParseDateString(&thread.Created_at)

	return thread, nil
}

func (driver *dbDriver) getPost(id string) (models.Posts, error) {
	post := models.Posts{}

	row := driver.db.QueryRowContext(context.Background(), `SELECT title, content, created_at, created_by from Posts WHERE ID = ?`, id)

	row.Scan(&post.Title, &post.Content, &post.Created_at, &post.Created_by)

	utils.ParseDateString(&post.Created_at)

	return post, nil
}

func (driver *dbDriver) createPost(title, content, created_by string) error {
	statement := `INSERT INTO Posts (title, content, created_at, created_by) values (?, ?, datetime('now'), ?)`

	_, err := driver.db.ExecContext(context.Background(), statement, title, content, created_by)

	return err
}

func (driver *dbDriver) createThread(post_id, content, parent_id, replied_by string) error {

	topLevelThread := `INSERT INTO Threads (post_id, parent_id, content, created_at, created_by) values (?, NULL, ?, datetime('now'), ?)`
	followThread := `INSERT INTO Threads (post_id, parent_id, content, created_at, created_by) values (?, ?, ?, datetime('now'), ?)`

	if parent_id == "<nil>" {
		_, err := driver.db.ExecContext(context.Background(), topLevelThread, post_id, content, replied_by)
		return err
	} else {
		_, err := driver.db.ExecContext(context.Background(), followThread, post_id, parent_id, content, replied_by)
		return err
	}
}

func (driver *dbDriver) createUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)

	if err != nil {
		return err
	}

	_, err2 := driver.db.ExecContext(context.Background(), `INSERT INTO Users (username, password) VALUES (?, ?)`, username, hash)

	if err2 != nil {
		return err
	}

	return nil
}

func (driver *dbDriver) authUser(username, password string) (string, error) {

	var hashed []byte

	row := driver.db.QueryRowContext(context.Background(), `SELECT password FROM Users WHERE username = ?`, username)
	err := row.Scan(&hashed)

	if err != nil {
		fmt.Println("ERROR IN SCAN")
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))

	if err != nil {
		fmt.Println("ERROR IN COMPARE")
		return "", err
	}

	return username, nil
}

func (driver *dbDriver) getUserBySessionToken(token string) (models.Session, error) {
	var session models.Session
	var expiryString string

	row := driver.db.QueryRowContext(context.Background(), `SELECT token, user, expires_at FROM Sessions WHERE token = ?`, token)

	err := row.Scan(&session.Token, &session.User, &expiryString)

	if err != nil {
		log.Fatal("Error retrieving session:", err)
		return models.Session{}, err
	}

	expiry, err := time.Parse(time.RFC3339, expiryString)

	if err != nil {
		log.Fatal("Error parsing time:", err)
		return models.Session{}, err
	}

	session.Expires_at = expiry

	return session, nil
}

func (driver *dbDriver) createSession(token, username, expiry string) error {

	_, err := driver.db.ExecContext(context.Background(), `INSERT INTO Sessions (token, user, expires_at) VALUES (?, ?, ?)`, token, username, expiry)

	if err != nil {
		return err
	}

	return nil
}

func (driver *dbDriver) deleteSession(username string) bool {

	_, err := driver.db.ExecContext(context.Background(), `DELETE FROM Sessions WHERE user = ?`, username)

	if err != nil {
		fmt.Println(err)
	}

	return err == nil
}
