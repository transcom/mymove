package models

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/satori/go.uuid"
)

var dbConnection *pop.Connection

func TestOptionalProperty(t *testing.T) {
	reporterName := "Janice Doe"

	hasReporter := Issue{
		Description:  "this describes an issue with a reporter",
		ReporterName: nulls.NewString(reporterName),
	}

	if err := dbConnection.Create(&hasReporter); err != nil {
		t.Fatal("Didn't write it to the db")
	}

	if hasReporter.ID == uuid.Nil {
		t.Error("didn't get an ID back")
	}

	if !hasReporter.ReporterName.Valid || hasReporter.ReporterName.String != reporterName {
		t.Error("didn't get the reporter name back right.")
	}

	sansReporter := Issue{
		Description: "This describes an issue without a reporter",
	}

	if err := dbConnection.Create(&sansReporter); err != nil {
		t.Fatal("Didn't write sans to the db")
	}

	if sansReporter.ReporterName.Valid {
		t.Error("Somehow got a valid name back")
	}
}

func setupDBConnection() {

	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	dbConnection = conn

}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
