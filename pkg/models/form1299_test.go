package models_test

import (
	"log"
	"testing"

	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection

func TestBasicInstantiation(t *testing.T) {

	// Given: an instance of form 1299
	newForm := models.Form1299{
		MobileHomeHeightFt:     12,
		ServiceMemberFirstName: "Jane",
		ServiceMemberLastName:  "Goodall",
		HhgProgearPounds:       12,
	}

	// When: A Form entry is created in the db
	err := dbConnection.Create(&newForm)
	if err != nil {
		t.Fatal("Didn't write to the db.")
	}

	// Then: assert that the ID has a value of type uuid
	if newForm.ID == uuid.Nil {
		t.Error("No UUID was set by writing to the DB.")
	}

	// And: assert that the form contains the expected values
	if (newForm.MobileHomeHeightFt != 12) &&
		(newForm.ServiceMemberFirstName != "Jane") &&
		(newForm.ServiceMemberLastName != "Goodall") &&
		(newForm.HhgProgearPounds != 12) {
		t.Error("Value not properly set")
	}

	// When: column value is updated
	oldUpdated := newForm.UpdatedAt
	newForm.ServiceMemberFirstName = "Bob"
	err1 := dbConnection.Update(&newForm)
	if err1 != nil {
		t.Fatal("Didn't update entry.")
	}

	// Then: Values are updated accordingly
	if (oldUpdated.After(newForm.UpdatedAt)) &&
		(newForm.ServiceMemberFirstName != "Bob") {
		t.Error("Values weren't updated.")
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
