-- +goose Up
CREATE TABLE leaderboard (
                             id SERIAL PRIMARY KEY,
                             player_id INT REFERENCES players(id) ON DELETE CASCADE,
                             score INT NOT NULL,
                             rank INT NOT NULL,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE leaderboard;