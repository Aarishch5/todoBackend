package migrations

import (
	"embed"
	"errors"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed 001_init.up.sql
var migrationFile embed.FS

var DB *sqlx.DB

func OpenConnection() error {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL environment variable is not set")
	}

	var err error

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	return runMigrations()
}

func runMigrations() error {
	sqlBytes, err := migrationFile.ReadFile("001_init.up.sql")
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(sqlBytes))
	return err
}
