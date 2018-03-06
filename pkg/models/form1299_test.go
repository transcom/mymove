package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"

	. "github.com/transcom/mymove/pkg/models"
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

func TestFetchAllForm1299s(t *testing.T) {
	// Need to know how many we start with until we start using transactions
	initialSet := []Form1299{}
	dbConnection.All(&initialSet)
	initialLength := len(initialSet)

	// Given: A couple 1299 forms
	numRecords := 2 + initialLength
	serviceMemberFirstName1 := "Jane"
	form1 := Form1299{
		ServiceMemberFirstName: &serviceMemberFirstName1,
	}
	serviceMemberFirstName2 := "Joe"
	form2 := Form1299{
		ServiceMemberFirstName: &serviceMemberFirstName2,
	}

	err1 := dbConnection.Create(&form1)
	err2 := dbConnection.Create(&form2)
	if err1 != nil || err2 != nil {
		t.Fatal("Didn't write to the db.")
	}

	// When: Fetch all Form1299s is called
	form1299s, _ := FetchAllForm1299s(dbConnection)

	// Then: Two records are returned as expected
	if length := len(form1299s); length != numRecords {
		t.Fatal(fmt.Sprintf("Returned %d records instead of %d", length, numRecords))
	}
}

func TestFetchForm1299ByID(t *testing.T) {
	// Given: A 1299 form
	serviceMemberFirstName1 := "Janine"
	form := Form1299{
		ServiceMemberFirstName: &serviceMemberFirstName1,
	}

	if err := dbConnection.Create(&form); err != nil {
		t.Fatal("Didn't write to the db.")
	}

	// When: Fetch form 1299 by ID is called
	id := strfmt.UUID(form.ID.String())
	returnedForm, err := FetchForm1299ByID(dbConnection, id)

	// Then: The specified record is returned
	if err != nil {
		t.Fatal("No record found for that ID")
	}

	if *form.ServiceMemberFirstName != *returnedForm.ServiceMemberFirstName {
		t.Fatal(fmt.Sprintf("The returned form's contents don't match expected: %s vs %s",
			*form.ServiceMemberFirstName, *returnedForm.ServiceMemberFirstName))
	}
}

func TestFetchForm1299ByIDReturnsError(t *testing.T) {
	// Given: A fake ID
	unknownID := "2400c3c5-019d-4031-9c27-8a553e022297"

	// When: Fetch form 1299 by ID is called
	_, err := FetchForm1299ByID(dbConnection, strfmt.UUID(unknownID))

	// Then: No record is returned
	if err == nil || err.Error() != "sql: no rows in result set" {
		t.Fatal("There should be no record for that ID")
	}
}

func TestPopulateAddressesMethod(t *testing.T) {
	// Given: A form that should have an address but doesn't
	address := Address{
		StreetAddress1: "123 My Way",
		City:           "Seattle",
		State:          "NY",
		Zip:            "12345",
	}
	dbConnection.Create(&address)

	sentForm := Form1299{
		OriginOfficeAddress:   &address,
		OriginOfficeAddressID: &address.ID,
	}
	dbConnection.Create(&sentForm)
	returnedForm := Form1299{}
	dbConnection.Find(&returnedForm, sentForm.ID)

	if returnedForm.OriginOfficeAddress != nil {
		t.Fatal("Form should not have an address yet")
	}

	// When: PopulateAddresses is called
	returnedForm.PopulateAddresses(dbConnection)

	// Then: Addresses are populated as expected
	if returnedForm.OriginOfficeAddress == nil {
		t.Fatal("Form should have an address populated")
	}
}

func TestCreateForm1299WithAddressesSavesAddresses(t *testing.T) {
	// Given: A form1299 model with an address struct
	address := Address{
		StreetAddress1: "123 My Way",
		City:           "Seattle",
		State:          "NY",
		Zip:            "12345",
	}
	form := Form1299{
		OriginOfficeAddress: &address,
	}

	// When: CreateForm1299WithAddressesSavesAddresses is called on the form
	verrs, err := CreateForm1299WithAddresses(dbConnection, &form)

	// Then: The address and form should be saved to DB, ID populated
	blankUUID := uuid.UUID{}
	if address.ID != blankUUID && *form.OriginOfficeAddressID != address.ID {
		t.Fatal("Address ID should match saved ID")
	}

	// And: There should be no errors
	if verrs.HasAny() || err != nil {
		t.Fatal("There was an error while saving form")
	}
}

func TestCreateForm1299WithAddressesReturnsErrors(t *testing.T) {
	// Given: A form1299 model with a blank StreetAddress1 field
	address := Address{
		City:  "Seattle",
		State: "NY",
		Zip:   "12345",
	}
	form := Form1299{
		OriginOfficeAddress: &address,
	}

	// When: CreateForm1299WithAddressesSavesAddresses is called on the form
	verrs, _ := CreateForm1299WithAddresses(dbConnection, &form)

	// Then: The address and form should not be saved to DB, ID left blank
	if form.OriginOfficeAddressID != nil {
		t.Fatal("ID field should not be populated")
	}

	// And: There should be validation errors
	if !verrs.HasAny() {
		t.Fatal("There should be a validation error")
	}
}
