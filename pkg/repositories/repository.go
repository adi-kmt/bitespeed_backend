package repositories

import (
	db "github.com/adi-kmt/bitespeed_backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewRepository(q *db.Queries, conn *pgxpool.Pool) *Repository {
	return &Repository{
		q:    q,
		conn: conn,
	}
}
