package tests

import (
	"fmt"
	"os"
	"testing"

	"dfss/mgdb"
)

func TestMain(m *testing.M) {

	dbURI = os.Getenv("DFSS_MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost/dfss-test"
	}

	var err error
	dbManager, err = mgdb.NewManager(dbURI)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code := m.Run()
	eraseDatabase()
	os.Exit(code)
}
