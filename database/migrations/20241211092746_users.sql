-- +goose Up
-- +goose StatementBegin

CREATE TABLE Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
