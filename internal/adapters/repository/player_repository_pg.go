package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"log"
)

type PlayerRepositoryPg struct {
	db *pgxpool.Pool
}

func NewPlayerRepositoryPg(db *pgxpool.Pool) repository.PlayerRepository {
	return &PlayerRepositoryPg{db: db}
}

func (r *PlayerRepositoryPg) GetPlayerByID(ctx context.Context, id int) (entities.Player, error) {
	var player entities.Player
	query := `
		SELECT id, username, level, total_matches, total_wins, last_login
		FROM players 
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(&player.ID, &player.Username, &player.Level, &player.TotalMatches, &player.TotalWins, &player.LastLogin)
	if err != nil {
		log.Println("Error fetching player by ID:", err)
		return entities.Player{}, err
	}
	return player, nil
}

func (r *PlayerRepositoryPg) UpdatePlayer(ctx context.Context, player entities.Player) error {
	query := `
		UPDATE players
		SET username = $1, level = $2, total_matches = $3, total_wins = $4, last_login = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, player.Username, player.Level, player.TotalMatches, player.TotalWins, player.LastLogin, player.ID)
	if err != nil {
		log.Println("Error updating player:", err)
		return err
	}
	return nil
}

func (r *PlayerRepositoryPg) DeletePlayer(ctx context.Context, id int) error {
	query := `
		DELETE FROM players
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error deleting player:", err)
		return err
	}
	return nil
}

func (r *PlayerRepositoryPg) GetAllPlayers(ctx context.Context) ([]entities.Player, error) {
	var players []entities.Player
	query := `
		SELECT id, username, level, total_matches, total_wins, last_login
		FROM players
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Println("Error fetching all players:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player entities.Player
		err := rows.Scan(&player.ID, &player.Username, &player.Level, &player.TotalMatches, &player.TotalWins, &player.LastLogin)
		if err != nil {
			log.Println("Error scanning player data:", err)
			return nil, err
		}
		players = append(players, player)
	}

	if rows.Err() != nil {
		log.Println("Error iterating over rows:", rows.Err())
		return nil, rows.Err()
	}
	return players, nil
}

func (r *PlayerRepositoryPg) CreatePlayer(ctx context.Context, player entities.Player) (int, error) {
	var id int
	query := `
        INSERT INTO players (username, level, total_matches, total_wins, last_login)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query, player.Username, player.Level, player.TotalMatches, player.TotalWins, player.LastLogin).Scan(&id)
	if err != nil {
		log.Println("Error creating player:", err)
		return 0, err
	}
	return id, nil
}
