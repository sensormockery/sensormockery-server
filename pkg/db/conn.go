package db

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	// DBUser for the db.
	DBUser = "DB_USER"
	// DBPass for the db.
	DBPass = "DB_PASS"
	// DBDomain where the db is hosted.
	DBDomain = "DB_DOMAIN"
	// DBPort on which the db serves.
	DBPort = "DB_PORT"
	// DBName of the db.
	DBName = "DB_NAME"
	// DBType of the DBMS. Example postgres, mysql, etc.
	DBType = "postgres"
)

var dbConn *sql.DB = nil

// GetDBConn returns a connection to the sensormockery db.
func GetDBConn() (*sql.DB, error) {
	if dbConn == nil {
		dbURL := getDBURL()
		db, err := sql.Open(DBType, dbURL)

		if err != nil {
			return dbConn, err
		}

		dbConn = db
	}

	return dbConn, nil
}

func getDBURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		DBType,
		os.Getenv(DBUser),
		os.Getenv(DBPass),
		os.Getenv(DBDomain),
		os.Getenv(DBPort),
		os.Getenv(DBName))
}
