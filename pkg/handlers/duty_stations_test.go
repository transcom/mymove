package handlers

import (
	"net/http/httptest"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestSearchDutyStationHandler() {
	t := suite.T()

	// Need a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.mustSave(&address)

	station1 := models.DutyStation{
		Name:      "First Station",
		Branch:    internalmessages.MilitaryBranchARMY,
		AddressID: address.ID,
	}
	suite.mustSave(&station1)

	station2 := models.DutyStation{
		Name:      "Second Station",
		Branch:    internalmessages.MilitaryBranchARMY,
		AddressID: address.ID,
	}
	suite.mustSave(&station2)

	req := httptest.NewRequest("GET", "/duty_stations", nil)
	newSearchParams := stationop.SearchDutyStationsParams{
		HTTPRequest: req,
		Branch:      string(internalmessages.MilitaryBranchARMY),
		Search:      "first",
	}

	// Make sure the context contains the auth values
	ctx := newSearchParams.HTTPRequest.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	newSearchParams.HTTPRequest = newSearchParams.HTTPRequest.WithContext(ctx)

	handler := SearchDutyStationsHandler(NewHandlerContext(suite.db, suite.logger, nil))
	response := handler.Handle(newSearchParams)

	// Assert we got back the 201 response
	searchResponse := response.(*stationop.SearchDutyStationsOK)
	stationPayloads := searchResponse.Payload

	if len(stationPayloads) != 1 {
		t.Errorf("Should have only got 1 response, got %v", len(stationPayloads))
	}

	if *stationPayloads[0].Name != "First Station" {
		t.Errorf("Station name should have been \"First Station \", got %v", stationPayloads[0].Name)
	}

	if *stationPayloads[0].Address.City != "city" {
		t.Error("Address should have been loaded")
	}
}
