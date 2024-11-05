-- +goose Up
INSERT INTO players (username, level, total_matches, total_wins) VALUES
                                                                     ('player1', 1, 5, 3),
                                                                     ('player2', 2, 10, 5),
                                                                     ('player3', 3, 15, 8);

INSERT INTO auth (username, password_hash) VALUES
                                               ('player1', 'hashed_password_1'),
                                               ('player2', 'hashed_password_2'),
                                               ('player3', 'hashed_password_3');

-- +goose Down
DELETE FROM auth WHERE username IN ('player1', 'player2', 'player3');
DELETE FROM players WHERE username IN ('player1', 'player2', 'player3');
