-- +goose Up
-- +goose StatementBegin

CREATE TABLE Posts (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    FOREIGN KEY (created_by) REFERENCES Users(username) ON DELETE CASCADE
);

CREATE TABLE Threads (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    parent_id INTEGER,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    FOREIGN KEY (post_id) REFERENCES Posts(ID) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES Threads(ID) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES Users(username) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
