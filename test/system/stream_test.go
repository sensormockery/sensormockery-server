package system_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	apiV1 "github.com/sensormockery/sensormockery-server/pkg/api/v1"
	"github.com/sensormockery/sensormockery-server/pkg/db"
	"github.com/sensormockery/sensormockery-server/pkg/db/dto"
)

var _ = Describe("Stream", func() {
	var (
		apiV1URL string
		dbConn   *sql.DB
		err      error
	)

	BeforeEach(func() {
		apiV1URL = fmt.Sprintf("%s%s", os.Getenv(APIURL), apiV1.APIPrefix)

		dbConn, err = db.GetDBConn()
		if err != nil {
			Fail("Fail obtaining db connection: " + err.Error())
		}
	})

	Describe("Calling create stream endpoint", func() {
		BeforeEach(func() {
			cleanUpTable(dbConn, dto.StreamsTable)
			resetTableSequence(dbConn, dto.StreamsTable, "id")
		})

		AfterEach(func() {
			cleanUpTable(dbConn, dto.StreamsTable)
			resetTableSequence(dbConn, dto.StreamsTable, "id")
		})

		Context("with a valid request", func() {
			It("should create a stream", func() {
				stream := &apiV1.Stream{
					WaveType:   "sine",
					Sensor:     "accelorometer",
					NoiseCoeff: 0.25,
				}

				req, err := json.Marshal(stream)
				if err != nil {
					Fail(err.Error())
				}

				resp, err := http.Post(apiV1URL+apiV1.CreateStreamPath, "application/json", bytes.NewBuffer(req))
				if err != nil {
					Fail(err.Error())
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					Fail(err.Error())
				}
				defer resp.Body.Close()

				createStreamResp := &apiV1.CreateStreamResp{}
				if err := json.Unmarshal(body, createStreamResp); err != nil {
					Fail(err.Error())
				}

				Expect(createStreamResp.ID).To(Equal(1))

				rowsCnt, err := numTableRows(dbConn, dto.StreamsTable)
				if err != nil {
					Fail(err.Error())
				}
				Expect(rowsCnt).To(Equal(1))
			})
		})
	})
})

func numTableRows(dbConn *sql.DB, table string) (count int, err error) {
	rows, err := dbConn.Query("SELECT COUNT(*) as count FROM " + table)
	if err != nil {
		Fail(err.Error())
	}

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}
