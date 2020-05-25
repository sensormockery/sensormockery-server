package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	// ChangelogPath holds the path to the migration scripts.
	ChangelogPath = "file://pkg/db/migration"
)

// Migrate applies database migration scripts.
func Migrate() error {
	dbConn, err := GetDBConn()

	if err != nil {
		return fmt.Errorf("Error obtaining db connection: %s", err.Error())
	}

	driver, err := dbConnToDriver(dbConn)

	if err != nil {
		return fmt.Errorf("Error converting db connection to go-migrate driver: %s", err.Error())
	}

	migrator, err := migrate.NewWithDatabaseInstance(ChangelogPath, DBType, driver)

	if err != nil {
		return fmt.Errorf("Error obtaining go-migrate migrator: %s", err.Error())
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Error migrating db: %s", err.Error())
	}

	return nil
}

func dbConnToDriver(dbConn *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(dbConn, &postgres.Config{})
}
