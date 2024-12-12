-- +goose Up
-- +goose StatementBegin
CREATE TABLE Sessions (
    token VARCHAR(255) PRIMARY KEY,
    user VARCHAR(255) NOT NULL UNIQUE,
    expires_at VARCHAR(255) NOT NULL,
    FOREIGN KEY (user) REFERENCES Users(username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM Sessions;
-- +goose StatementEnd
