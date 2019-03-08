package ordersapi

import (
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestPostRevisionNew() {
}

func (suite *HandlerSuite) TestPostRevisionNewAmendment() {
}

func (suite *HandlerSuite) TestPostRevisionNoApiPerm() {
	req := httptest.NewRequest("POST", "/orders/v1/orders/", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		MemberID:    "1234567890",
		OrdersNum:   "8675309",
	}

	handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionForbidden{}, response)
}

func (suite *HandlerSuite) TestPostRevisionWritePerms() {
	testCases := map[string]struct {
		cert   models.ClientCert
		issuer ordersmessages.Issuer
	}{
		"Army": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerArmy,
		},
		"Navy": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
			},
			ordersmessages.IssuerNavy,
		},
		"MarineCorps": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersWrite:   true,
				AllowArmyOrdersWrite:       true,
				AllowCoastGuardOrdersWrite: true,
				AllowNavyOrdersWrite:       true,
			},
			ordersmessages.IssuerMarineCorps,
		},
		"CoastGuard": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerCoastGuard,
		},
		"AirForce": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerAirForce,
		},
	}

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/orders/v1/orders/", nil)
			req = suite.AuthenticateClientCertRequest(req, &testCase.cert)

			params := ordersoperations.PostRevisionParams{
				HTTPRequest: req,
				Issuer:      string(testCase.issuer),
				MemberID:    "1234567890",
				OrdersNum:   "8675309",
			}

			handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.PostRevisionForbidden{}, response)
		})
	}
}
