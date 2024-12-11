package main

import (
	"database/sql"
	"discuz/database/models"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type dbDriver struct {
	db *sql.DB
}

type PostPage struct {
	Post models.Posts
	Threads []*models.Threads
}

func connectDB() (*dbDriver, error) {
	driver, err := sql.Open("sqlite3", "./database/app.db")

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

	rows, err := driver.db.Query(`SELECT * FROM Posts`)

	if err != nil  {
		return nil, err	
	}

	posts := []models.Posts{}

	for rows.Next() {
		post := models.Posts{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at)
		parseDateString(&post.Created_at)

		if err != nil  {
			return nil, err	
		}

		posts = append(posts, post)
	}

	return posts, nil	
}

func (driver *dbDriver) getPostTree(id string) (PostPage, error) {

	postPage := PostPage{
		Post: models.Posts{ID: id},
	}

	row := driver.db.QueryRow(`SELECT title, content, created_at from Posts WHERE ID = ?`, id)

	threadRows, err := driver.db.Query(`SELECT ID, thread_id, content, created_at from Threads WHERE post_id = ?`, id)

	for threadRows.Next() {
		thread := models.Threads{Post_id: id}
		threadRows.Scan(&thread.ID, &thread.Thread_id, &thread.Content, &thread.Created_at)
		parseDateString(&thread.Created_at)
		postPage.Threads = append(postPage.Threads, &thread)
	}

	tree := buildTree(postPage.Threads)
	postPage.Threads = tree

	if err != nil {
		return postPage, err
	}

	err2 := row.Scan(&postPage.Post.Title, &postPage.Post.Content, &postPage.Post.Created_at)
	parseDateString(&postPage.Post.Created_at)

	if err2 != nil {
		return postPage, err2
	}

	return postPage, nil 
}

func (driver *dbDriver) getThread(id string) (models.Threads, error) {
	thread := models.Threads{}

	row := driver.db.QueryRow(`SELECT content, created_at from Threads WHERE ID = ?`, id)

	row.Scan(&thread.Content, &thread.Created_at)

	return thread, nil
}

func (driver *dbDriver) getPost(id string) (models.Posts, error) {
	post := models.Posts{}

	row := driver.db.QueryRow(`SELECT title, content, created_at from Posts WHERE ID = ?`, id)

	row.Scan(&post.Title, &post.Content, &post.Created_at)

	parseDateString(&post.Created_at)

	return post, nil
}

func (driver *dbDriver) createPost(title string, content string) error {
	statement := `INSERT INTO Posts (title, content, created_at) values (?, ?, datetime('now'))`

	_, err2 := driver.db.Exec(statement, title, content)

	return err2
}

func (driver *dbDriver) createThread(post_id, content, thread_id string) error {

	topLevelThread := `INSERT INTO Threads (post_id, thread_id, content, created_at) values (?, NULL, ?, datetime('now'))`
	followThread := `INSERT INTO Threads (post_id, thread_id, content, created_at) values (?, ?, ?, datetime('now'))`

	if thread_id == "<nil>" {
		fmt.Printf("top level")
		_, err := driver.db.Exec(topLevelThread, post_id, content)
		return err
	} else {
		fmt.Printf("follow")
		_, err := driver.db.Exec(followThread, post_id, thread_id, content)
		return err
	}
}

func buildTree(threads []*models.Threads) []*models.Threads {

	treeMap := []*models.Threads{}

	for _, thread := range threads {
		if !thread.Thread_id.Valid {
			treeMap = append(treeMap, thread)
		} else {
			searchTree(treeMap, thread)
		}
	}

	return treeMap
}

func searchTree(treeMap []*models.Threads, target *models.Threads) {
	for _, thread := range treeMap {
		if thread.ID == target.Thread_id.String {
			thread.Children = append(thread.Children, target)
			return
		}
		if len(thread.Children) == 0 {
			continue
		} else {
			searchTree(thread.Children, target)
		}
	}
}

func parseDateString(datestring *string) {
	t, err := time.Parse("2006-01-02T15:04:05Z", *datestring)
	if err != nil {
		return
	}
	*datestring = fmt.Sprintf("%d %s %d, %d:%d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
}
