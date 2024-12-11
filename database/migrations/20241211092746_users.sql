-- +goose Up
-- +goose StatementBegin

CREATE TABLE Users (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(200) NOT NULL UNIQUE,
    password VARCHAR(200) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
