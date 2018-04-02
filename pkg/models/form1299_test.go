package models_test

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicForm1299Instantiation() {
	t := suite.T()

	// Given: an instance of form 1299 with some optional values
	var mobileHomeHeightft int64 = 12

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
	suite.mustSave(&newForm)

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
	err1 := suite.db.Update(&newForm)
	if err1 != nil {
		t.Fatal("Didn't update entry.")
	}

	// Then: Values are updated accordingly
	if (oldUpdated.After(newForm.UpdatedAt)) ||
		(*newForm.ServiceMemberFirstName != "Bob") {
		t.Error("Values weren't updated.")
	}
}

func (suite *ModelSuite) TestFetchAllForm1299s() {
	t := suite.T()
	// Need to know how many we start with until we start using transactions
	initialSet := []Form1299{}
	suite.db.All(&initialSet)
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

	suite.mustSave(&form1)
	suite.mustSave(&form2)

	// When: Fetch all Form1299s is called
	form1299s, _ := FetchAllForm1299s(suite.db)

	// Then: Two records are returned as expected
	if length := len(form1299s); length != numRecords {
		t.Fatal(fmt.Sprintf("Returned %d records instead of %d", length, numRecords))
	}
}

func (suite *ModelSuite) TestFetchForm1299ByID() {
	t := suite.T()
	// Given: A 1299 form
	serviceMemberFirstName1 := "Janine"
	form := Form1299{
		ServiceMemberFirstName: &serviceMemberFirstName1,
	}
	suite.mustSave(&form)

	// When: Fetch form 1299 by ID is called
	id := strfmt.UUID(form.ID.String())
	returnedForm, err := FetchForm1299ByID(suite.db, id)

	// Then: The specified record is returned
	if err != nil {
		t.Fatal("No record found for that ID")
	}

	if *form.ServiceMemberFirstName != *returnedForm.ServiceMemberFirstName {
		t.Fatal(fmt.Sprintf("The returned form's contents don't match expected: %s vs %s",
			*form.ServiceMemberFirstName, *returnedForm.ServiceMemberFirstName))
	}
}

func (suite *ModelSuite) TestFetchForm1299ByIDEagerLoads() {
	t := suite.T()
	// Given: A 1299 form
	address := Address{
		StreetAddress1: "123 My Way",
		City:           "Seattle",
		State:          "NY",
		PostalCode:     "12345",
	}
	suite.mustSave(&address)
	form := Form1299{
		OriginOfficeAddressID: &address.ID,
	}
	suite.mustSave(&form)

	// When: Fetch form 1299 by ID is called
	id := strfmt.UUID(form.ID.String())
	returnedForm, err := FetchForm1299ByID(suite.db, id)

	// Then: The specified record is returned and the address model is populated
	if err != nil {
		t.Fatal("No record found for that ID")
	}

	if returnedForm.OriginOfficeAddress == nil || returnedForm.OriginOfficeAddress.ID != *returnedForm.OriginOfficeAddressID {
		t.Fatal("Address model wasn't populated onto form 1299")
	}
}

func (suite *ModelSuite) TestFetchForm1299ByIDReturnsError() {
	t := suite.T()
	// Given: A fake ID
	unknownID := "2400c3c5-019d-4031-9c27-8a553e022297"

	// When: Fetch form 1299 by ID is called
	_, err := FetchForm1299ByID(suite.db, strfmt.UUID(unknownID))

	// Then: No record is returned
	if err == nil || err.Error() != "sql: no rows in result set" {
		t.Fatal("There should be no record for that ID")
	}
}

func (suite *ModelSuite) TestCreateForm1299WithAddressesSavesAddresses() {
	t := suite.T()

	// Given: A form1299 model with an address struct
	address := Address{
		StreetAddress1: "123 My Way",
		City:           "Seattle",
		State:          "NY",
		PostalCode:     "12345",
	}
	form := Form1299{
		OriginOfficeAddress: &address,
	}

	// When: CreateForm1299WithAddressesSavesAddresses is called on the form
	verrs, err := CreateForm1299WithAddresses(suite.db, &form)

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

func (suite *ModelSuite) TestCreateForm1299WithAddressesReturnsErrors() {
	t := suite.T()

	// Given: A form1299 model with a blank StreetAddress1 field
	address := Address{
		City:       "Seattle",
		State:      "NY",
		PostalCode: "12345",
	}
	form := Form1299{
		OriginOfficeAddress: &address,
	}

	// When: CreateForm1299WithAddressesSavesAddresses is called on the form
	verrs, _ := CreateForm1299WithAddresses(suite.db, &form)

	// Then: The address and form should not be saved to DB, ID left blank
	if form.OriginOfficeAddressID != nil {
		t.Fatal("ID field should not be populated")
	}

	// And: There should be validation errors
	if !verrs.HasAny() {
		t.Fatal("There should be a validation error")
	}
}
