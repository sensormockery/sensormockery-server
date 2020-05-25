package system_test

import (
	"database/sql"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const APIURL = "API_URL"

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Suite")
}

func cleanUpTable(dbConn *sql.DB, table string) {
	dbConn.Exec("DELETE FROM " + table)
}

func resetTableSequence(dbConn *sql.DB, table, seqCol string) {
	dbConn.Exec("ALTER SEQUENCE " + table + "_" + seqCol + "_seq RESTART WITH 1")
}
