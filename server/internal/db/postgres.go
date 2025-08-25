package db

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	log "github.com/sirupsen/logrus"
)

var (
	once sync.Once
	Psql *sql.DB
)

func NewPostgres(dbURL string) *sql.DB {
	once.Do(func() {
		if dbURL == "" {
			log.Panic("DATABASE_URL is not set in environment")
		}

		db, err := sql.Open("pgx", dbURL)
		if err != nil {
			log.Panicf("failed to open database: %v", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Panicf("failed to ping database: %v", err)
		}

		log.Info("✅ Connected to PostgreSQL via sql.DB (pgx driver)")

		Psql = db

	})

	return Psql
}
