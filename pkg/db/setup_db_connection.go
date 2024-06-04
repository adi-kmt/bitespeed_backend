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

func InitPool(config *DbConfig) *DbConn {
	connectionString := fmt.Sprintf("postgresql://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
		config.Username, config.Password, config.Port, config.Database)

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
