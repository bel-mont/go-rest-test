-- +goose Up
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    player1_id INT REFERENCES players(id) ON DELETE CASCADE,
    player2_id INT REFERENCES players(id) ON DELETE CASCADE,
    winner_id INT REFERENCES players(id) ON DELETE CASCADE,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE matches;