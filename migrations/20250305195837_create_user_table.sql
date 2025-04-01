-- +goose Up
CREATE TABLE users (
    id SERIAL primary key,
    login VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
