-- +goose Up
CREATE TABLE matchmaking (
    id SERIAL PRIMARY KEY,
    player_id INT REFERENCES players(id) ON DELETE CASCADE,
    skill_level INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    queued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE matchmaking;