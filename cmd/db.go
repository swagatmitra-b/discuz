package main

import (
	"database/sql"
	"discuz/database/models"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type dbDriver struct {
	db *sql.DB
}

func connectDB() (*dbDriver, error) {
	driver, err := sql.Open("sqlite3", "./app.db")

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

	if err != nil {
		return nil, err
	}

	posts := []models.Posts{}

	for rows.Next() {
		post := models.Posts{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_at, &post.Created_by)
		parseDateString(&post.Created_at)

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

	row := driver.db.QueryRow(`SELECT title, content, created_at from Posts WHERE ID = ?`, id)

	threadRows, err := driver.db.Query(`SELECT ID, parent_id, content, created_at from Threads WHERE post_id = ?`, id)

	for threadRows.Next() {
		thread := models.Threads{Post_id: id}
		threadRows.Scan(&thread.ID, &thread.Parent_id, &thread.Content, &thread.Created_at)
		parseDateString(&thread.Created_at)
		postPage.Threads = append(postPage.Threads, &thread)
	}

	tree := buildTree(postPage.Threads)
	postPage.Threads = tree

	if err != nil {
		return postPage, err
	}

	err = row.Scan(&postPage.Post.Title, &postPage.Post.Content, &postPage.Post.Created_at)
	parseDateString(&postPage.Post.Created_at)

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
	statement := `INSERT INTO Posts (title, content, created_at, created_by) values (?, ?, datetime('now'), 1)`

	_, err := driver.db.Exec(statement, title, content)

	return err
}

func (driver *dbDriver) createThread(post_id, content, parent_id string) error {

	topLevelThread := `INSERT INTO Threads (post_id, parent_id, content, created_at) values (?, NULL, ?, datetime('now'))`
	followThread := `INSERT INTO Threads (post_id, parent_id, content, created_at) values (?, ?, ?, datetime('now'))`

	if parent_id == "<nil>" {
		fmt.Printf("top level")
		_, err := driver.db.Exec(topLevelThread, post_id, content)
		return err
	} else {
		fmt.Printf("follow")
		_, err := driver.db.Exec(followThread, post_id, parent_id, content)
		return err
	}
}

func (driver *dbDriver) createUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)

	if err != nil {
		return err
	}

	_, err2 := driver.db.Exec(`INSERT INTO Users (username, password) VALUES (?, ?)`, username, hash)

	if err2 != nil {
		return err
	}

	return nil
}

func (driver *dbDriver) authUser(username, password string) (string, error) {

	var hashed []byte

	row := driver.db.QueryRow(`SELECT password FROM Users WHERE username = ?`, username)
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

func (driver *dbDriver) getUser(username string) (models.User, error) {
	var user models.User

	row := driver.db.QueryRow(`SELECT id, username, password FROM Users WHERE username = ?`, username)

	err := row.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func buildTree(threads []*models.Threads) []*models.Threads {

	treeMap := []*models.Threads{}

	for _, thread := range threads {
		if !thread.Parent_id.Valid {
			treeMap = append(treeMap, thread)
		} else {
			searchTree(treeMap, thread)
		}
	}

	for _, thread := range treeMap {
		thread.Root_id = thread.ID
		assignRootThread(thread, thread.ID)
	}

	return treeMap
}

func searchTree(treeMap []*models.Threads, target *models.Threads) {
	for _, thread := range treeMap {
		if thread.ID == target.Parent_id.String {
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

func assignRootThread(thread *models.Threads, id string) {
	for _, thread := range thread.Children {
		thread.Root_id = id
		if len(thread.Children) != 0 {
			assignRootThread(thread, id)
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

// func printTree(tree []*models.Threads) {
// 	for _, thread := range tree {
// 		fmt.Println(thread)
// 		if len(thread.Children) == 0 {
// 			fmt.Println(thread)
// 			return
// 		} else {
// 			printTree(thread.Children)
// 		}
// 	}
// }
