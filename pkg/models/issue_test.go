package models_test

import (
	"github.com/satori/go.uuid"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestOptionalProperty() {
	t := suite.T()

	reporterName := "Janice Doe"

	hasReporter := Issue{
		Description:  "this describes an issue with a reporter",
		ReporterName: &reporterName,
	}

	if err := suite.db.Create(&hasReporter); err != nil {
		t.Fatal("Didn't write it to the db")
	}

	if hasReporter.ID == uuid.Nil {
		t.Error("didn't get an ID back")
	}

	if hasReporter.ReporterName == nil || *hasReporter.ReporterName != reporterName {
		t.Error("didn't get the reporter name back right.")
	}

	sansReporter := Issue{
		Description: "This describes an issue without a reporter",
	}

	if err := suite.db.Create(&sansReporter); err != nil {
		t.Fatal("Didn't write sans to the db")
	}

	if sansReporter.ReporterName != nil {
		t.Error("Somehow got a valid name back")
	}
}
