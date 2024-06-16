package db

import (
	"context"
	"fmt"

	db "github.com/adi-kmt/bitespeed_backend/db/sqlc"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConn struct {
	DbPool    *pgxpool.Pool
	DbQueries *db.Queries
}

func InitPool(config *dbConfig) *DbConn {
	connectionString := fmt.Sprintf("postgresql://%s:%s@db:%s/%s?sslmode=disable",
		config.username, config.password, config.database, config.port)

	ctx := context.Background()

	d, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Errorf("error opening database: %v", err)
	}

	if err := d.Ping(ctx); err != nil {
		log.Errorf("error pinging database: %v", err)
	}

	return &DbConn{
		DbPool:    d,
		DbQueries: db.New(d),
	}
}

func (db *DbConn) Close() error {
	if db.DbPool != nil {
		db.DbPool.Close()
	}
	return nil
}
