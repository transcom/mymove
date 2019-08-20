package internalapi

import (
	"net/http/httptest"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestSearchDutyStationHandler() {
	t := suite.T()

	// Need a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.MustSave(&user)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.MustSave(&address)

	station1 := models.DutyStation{
		Name:        "First Station",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station1)

	station2 := models.DutyStation{
		Name:        "Second Station",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station2)

	req := httptest.NewRequest("GET", "/duty_stations", nil)

	// Make sure the context contains the auth values
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, session)

	newSearchParams := stationop.SearchDutyStationsParams{
		HTTPRequest: req.WithContext(ctx),
		Search:      "first",
	}

	handler := SearchDutyStationsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
