package database

import (
	"database/sql"
	"fmt"
	"time"

	
    _ "github.com/jackc/pgx/v5/stdlib"
)



type Config struct {
	Driver		string // "postgres"
	Port		int
	Host		string
	DBName		string
	Username	string
	Password	string
	SSLMode		string
}

type DB struct {
	*sql.DB
}

func NewConnection(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Health() error {
	return db.Ping()
}