package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	locationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
)

func (suite *HandlerSuite) TestSearchDutyLocationHandler() {
	t := suite.T()

	// Need a logged in user
	lgu := uuid.Must(uuid.NewV4()).String()
	user := models.User{
		OktaID:    lgu,
		OktaEmail: "email@example.com",
	}
	suite.MustSave(&user)

	newAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "CA",
		PostalCode:     "12345",
		County:         "County",
	}
	factory.FetchOrBuildCountry(suite.AppContextForTest().DB(), nil, nil)
	addressCreator := address.NewAddressCreator()
	createdAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAddress)
	suite.NoError(err)

	location1 := models.DutyLocation{
		Name:        "First Location",
		AddressID:   createdAddress.ID,
		Affiliation: internalmessages.NewAffiliation(internalmessages.AffiliationAIRFORCE),
	}
	suite.MustSave(&location1)

	location2 := models.DutyLocation{
		Name:        "Second Location",
		AddressID:   createdAddress.ID,
		Affiliation: internalmessages.NewAffiliation(internalmessages.AffiliationAIRFORCE),
	}
	suite.MustSave(&location2)

	req := httptest.NewRequest("GET", "/duty_locations", nil)

	// Make sure the context contains the auth values
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, session)

	newSearchParams := locationop.SearchDutyLocationsParams{
		HTTPRequest: req.WithContext(ctx),
		Search:      "first",
	}

	handler := SearchDutyLocationsHandler{suite.HandlerConfig()}
	response := handler.Handle(newSearchParams)

	// Assert we got back the 201 response
	searchResponse := response.(*locationop.SearchDutyLocationsOK)
	locationPayloads := searchResponse.Payload

	suite.NoError(locationPayloads.Validate(strfmt.Default))

	if len(locationPayloads) != 1 {
		t.Errorf("Should have 1 responses, got %v", len(locationPayloads))
	}

	if *locationPayloads[0].Name != "First Location" {
		t.Errorf("Location name should have been \"First Location \", got %v", locationPayloads[0].Name)
	}

}
