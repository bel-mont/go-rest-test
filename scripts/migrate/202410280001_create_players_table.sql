-- +goose Up
CREATE TABLE players (
                         id SERIAL PRIMARY KEY,
                         username VARCHAR(100) UNIQUE NOT NULL,
                         level INT DEFAULT 1,
                         total_matches INT DEFAULT 0,
                         total_wins INT DEFAULT 0,
                         last_login TIMESTAMP
);

-- +goose Down
DROP TABLE players;
