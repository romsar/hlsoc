package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DB struct {
	db *sql.DB
}

func Open(dsn string) (*DB, error) {
	db, err := openConnection(dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open connection for master db with dsn %s: %w", dsn, err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.Close()
}

func openConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unnable to connect to pg: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(0)

	return db, nil
}
