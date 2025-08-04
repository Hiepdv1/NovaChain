package config

import (
	"ChainServer/internal/common/env"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func InitPostgres() *sql.DB {
	envConfig := env.New()
	dbURL := envConfig.DB_URL

	if dbURL == "" {
		log.Panic("DATABASE_URL is not set in environment")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Panicf("❌ Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	DB = db

	log.Info("✅ Connected to PostgreSQL via sql.DB (pgx driver)")

	return db
}
