package bootstrap

import (
	"ChainServer/internal/common/env"
	"ChainServer/internal/db"
	"database/sql"
)

type BootstrapDB struct {
	Psql *sql.DB
}

func StartConnDB() *BootstrapDB {
	psql := db.NewPostgres(env.Cfg.DB_URL)

	return &BootstrapDB{
		Psql: psql,
	}
}
