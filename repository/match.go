package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-test/models"
)

// MatchRepository is a repository interface for matches
type MatchRepository interface {
	GetMatchList(ctx context.Context) ([]models.MatchView, error)
}

// NewMatchRepository returns an instance of a match repository
func NewMatchRepository(db *pgxpool.Pool) MatchRepository {
	return &matchRepository{db: db}
}

type matchRepository struct {
	db *pgxpool.Pool
}

// GetMatchList retrieves the list of matches and returns them as MatchView structs
func (r *matchRepository) GetMatchList(ctx context.Context) ([]models.MatchView, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p1.username AS player1, p2.username AS player2, 
		       CASE WHEN m.winner_id = m.player1_id THEN p1.username ELSE p2.username END AS winner
		FROM matches m
		JOIN players p1 ON m.player1_id = p1.id
		JOIN players p2 ON m.player2_id = p2.id
		ORDER BY m.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.MatchView
	for rows.Next() {
		var match models.MatchView
		if err := rows.Scan(&match.Player1, &match.Player2, &match.Winner); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, rows.Err()
}
