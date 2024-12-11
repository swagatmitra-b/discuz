package models

import (
	"database/sql"
)

type Posts struct {
	ID         string
	Title      string
	Content    string
	Created_at string
	Created_by string
}

type Threads struct {
	ID         string
	Post_id    string
	Root_id    string
	Parent_id  sql.NullString
	Content    string
	Created_at string
	Created_by string
	Children   []*Threads
}

type User struct {
	ID       string
	Username string
	Password string
}

type PostPage struct {
	Post    Posts
	Threads []*Threads
}
