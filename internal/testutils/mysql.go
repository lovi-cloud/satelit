package testutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lovi-cloud/satelit/pkg/datastore"
)

const schemaDirRelativePathFormat = "%s/../mysql/%s"
const fixturesDirRelativePathFormat = "%s/../mysql/fixtures/%s"

// SetupDefaultFixtures is insert query for all test
func SetupDefaultFixtures() {
	_, pwd, _, _ := runtime.Caller(0)

	defaultFixtureDir := fmt.Sprintf(fixturesDirRelativePathFormat, path.Dir(pwd), "default")
	defaultFixturePathes := walkSchema(defaultFixtureDir)
	for _, dpath := range defaultFixturePathes {
		execSchema(dpath)
	}
}

func walkSchema(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func execSchema(fpath string) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalf("schema reading error: %v", err)
	}

	queries := strings.Split(string(b), ";")

	for _, query := range queries[:len(queries)-1] {
		_, err = testDB.Exec(query)
		if err != nil {
			log.Fatalf("exec schema error: %v, query: %s", err, query)
		}
	}
}

func createTablesIfNotExist() {
	_, pwd, _, _ := runtime.Caller(0)
	schemaPath := fmt.Sprintf(schemaDirRelativePathFormat, path.Dir(pwd), "schema.sql")
	execSchema(schemaPath)
}

func truncateTables() {
	rows, err := testDB.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("show tables error: %#v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Fatalf("show table error: %#v", err)
			continue
		}

		cmds := []string{
			"SET FOREIGN_KEY_CHECKS = 0",
			fmt.Sprintf("TRUNCATE %s", tableName),
			"SET FOREIGN_KEY_CHECKS = 1",
		}
		for _, cmd := range cmds {
			_, err := testDB.Exec(cmd)

			if err != nil {
				mysqlErr, ok := err.(*mysql.MySQLError)

				if ok {
					if mysqlErr.Number == 0xde2 {
						// is rejected
						continue
					}
				} else {
					log.Fatalf("truncate error: %#v", err)
					continue
				}
			}
		}
	}
}

// GetTestDatastore return pointer of datastore
func GetTestDatastore() (datastore.Datastore, func()) {
	if testDatastore == nil {
		panic("datastore is not initialized yet")
	}

	return testDatastore, func() { truncateTables() }
}

// GetTestDB return pointer of testDB
func GetTestDB() (*sqlx.DB, func()) {
	if testDB == nil {
		panic("testDB is not initialized yet")
	}

	return testDB, func() { truncateTables() }
}
