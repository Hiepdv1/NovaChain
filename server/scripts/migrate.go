package scripts

import (
	"database/sql"
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

var (
	_, file, _, _ = runtime.Caller(0)
	root          = path.Dir(path.Join(file, "../"))
)

func AutoMigrate(db *sql.DB) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Panicf("❌ Failed to init migration driver: %v", err)
	}

	migrationPath := fmt.Sprintf("file://%s", filepath.ToSlash(filepath.Join(root, "migrations")))

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)

	if err != nil {
		log.Panicf("❌ Failed to create migrate instance: %v", err)
	}

	err = m.Up()

	switch err {
	case nil:
		log.Info("✅ Migration complete.")
	case migrate.ErrNoChange:
		log.Info("ℹ️ No new migrations.")

	default:
		log.Panicf("❌ Migration failed: %v", err)

	}
}
