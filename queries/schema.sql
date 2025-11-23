-- Dummy schema file for sqlc (since we're using an external database)
-- This helps sqlc understand the table structure

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
