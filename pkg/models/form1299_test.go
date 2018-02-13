package models

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

func TestBasicForm1299Instantiation(t *testing.T) {

	// Given: an instance of form 1299 with some optional values
	var (
		mobileHomeHeightft int64 = 12
	)
	stationOrdersDate := time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)
	serviceMemberFirstName := "Jane"
	serviceMemberLastName := "Goodall"

	newForm := Form1299{
		MobileHomeHeightFt:     &mobileHomeHeightft,
		ServiceMemberFirstName: &serviceMemberFirstName,
		ServiceMemberLastName:  &serviceMemberLastName,
		StationOrdersDate:      &stationOrdersDate,
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
	if (*newForm.MobileHomeHeightFt != 12) ||
		(*newForm.ServiceMemberFirstName != "Jane") ||
		(*newForm.ServiceMemberLastName != "Goodall") ||
		(*newForm.StationOrdersDate != stationOrdersDate) {
		t.Fatal("Value not properly set")
	}

	// When: column values are updated
	oldUpdated := newForm.UpdatedAt
	serviceMemberFirstName = "Bob"
	err1 := dbConnection.Update(&newForm)
	if err1 != nil {
		t.Fatal("Didn't update entry.")
	}

	// Then: Values are updated accordingly
	if (oldUpdated.After(newForm.UpdatedAt)) ||
		(*newForm.ServiceMemberFirstName != "Bob") {
		t.Error("Values weren't updated.")
	}
}
