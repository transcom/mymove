package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	addressop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func fakeAddressPayload() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: swag.String("An address"),
		StreetAddress2: swag.String("Apt. 2"),
		StreetAddress3: swag.String("address line 3"),
		City:           swag.String("Happytown"),
		State:          swag.String("AL"),
		PostalCode:     swag.String("01234"),
	}
}

func (suite *HandlerSuite) TestShowAddressHandler() {

	suite.Run("successful lookup", func() {
		address := models.Address{
			StreetAddress1: "some address",
			City:           "city",
			State:          "state",
			PostalCode:     "12345",
		}
		suite.MustSave(&address)

		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		fakeUUID, _ := uuid.FromString("not-valid-uuid")

		tests := []struct {
			ID        uuid.UUID
			hasResult bool
			resultID  string
		}{
			{ID: address.ID, hasResult: true, resultID: address.ID.String()},
			{ID: fakeUUID, hasResult: false, resultID: ""},
		}

		for _, ts := range tests {
			req := httptest.NewRequest("GET", "/addresses/"+ts.ID.String(), nil)
			req = suite.AuthenticateUserRequest(req, requestUser)

			params := addressop.ShowAddressParams{
				HTTPRequest: req,
				AddressID:   *handlers.FmtUUID(ts.ID),
			}

			handler := ShowAddressHandler{suite.HandlerConfig()}
			res := handler.Handle(params)

			response := res.(*addressop.ShowAddressOK)
			payload := response.Payload

			if ts.hasResult {
				suite.NotNil(payload, "Should have address record")
				suite.Equal(payload.ID.String(), ts.resultID, "Address ID doest match")
			} else {
				suite.Nil(payload, "Should not have address record")
			}
		}
	})

}
