package models

import (
	"database/sql"
)

type Posts struct {
	ID         string
	Title      string
	Content    string
	Created_at string
}

type Threads struct {
	ID         string
	Post_id    string
	Thread_id  sql.NullString
	Content    string
	Created_at string
	Children   []*Threads
}