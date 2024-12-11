-- +goose Up
-- +goose StatementBegin
-- CREATE TABLE Users (
--     ID INT PRIMARY KEY,
--     name VARCHAR(200),
--     profileID INTEGER UNIQUE
-- );

CREATE TABLE Posts (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200),
    FOREIGN KEY (created_by) REFERENCES Users(ID) ON DELETE CASCADE
);

CREATE TABLE Threads (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    thread_id INTEGER,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200),
    FOREIGN KEY (post_id) REFERENCES Posts(ID) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES Threads(ID) ON DELETE SET NULL
    FOREIGN KEY (created_by) REFERENCES Users(ID) ON DELETE CASCADE
);

CREATE TABLE Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(200) NOT NULL UNIQUE,
    password VARCHAR(200) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS
-- +goose StatementEnd
