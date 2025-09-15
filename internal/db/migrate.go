package db

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

var migrationsFS embed.FS

// RunMigrations menjalankan file SQL di db/migrations
func RunMigrations(gdb *gorm.DB) {
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("sql db: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("pg driver: %v", err)
	}

	// "migrations" sesuai dengan prefix pada //go:embed
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("iofs: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Printf("migrations applied ✅")
}
