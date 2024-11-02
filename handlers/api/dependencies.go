package api

import "github.com/jackc/pgx/v5/pgxpool"

var db *pgxpool.Pool

func SetDB(pool *pgxpool.Pool) {
	db = pool
}
