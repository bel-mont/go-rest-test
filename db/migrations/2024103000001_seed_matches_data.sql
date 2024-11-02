-- +goose Up

-- Match 1: player1 vs player2, player2 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (1, 2, 2);

-- Match 2: player1 vs player3, player3 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (1, 3, 3);

-- Match 3: player2 vs player3, player3 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (2, 3, 3);

-- Match 4: player1 vs player2, player1 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (1, 2, 1);

-- Match 5: player2 vs player3, player2 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (2, 3, 2);

-- Match 6: player3 vs player1, player1 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (3, 1, 1);

-- Match 7: player1 vs player3, player3 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (1, 3, 3);

-- Match 8: player2 vs player1, player2 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (2, 1, 2);

-- Match 9: player3 vs player2, player3 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (3, 2, 3);

-- Match 10: player3 vs player1, player3 wins
INSERT INTO public.matches (player1_id, player2_id, winner_id)
VALUES (3, 1, 3);

-- +goose Down

DELETE FROM public.matches
WHERE (player1_id = 1 AND player2_id = 2 AND winner_id = 2)
   OR (player1_id = 1 AND player2_id = 3 AND winner_id = 3)
   OR (player1_id = 2 AND player2_id = 3 AND winner_id = 3)
   OR (player1_id = 1 AND player2_id = 2 AND winner_id = 1)
   OR (player1_id = 2 AND player2_id = 3 AND winner_id = 2)
   OR (player1_id = 3 AND player2_id = 1 AND winner_id = 1)
   OR (player1_id = 1 AND player2_id = 3 AND winner_id = 3)
   OR (player1_id = 2 AND player2_id = 1 AND winner_id = 2)
   OR (player1_id = 3 AND player2_id = 2 AND winner_id = 3)
   OR (player1_id = 3 AND player2_id = 1 AND winner_id = 3);
