package ordersapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestPostRevisionNoApiPerm() {
	req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirDashForce),
		MemberID:    "1234567890",
		OrdersNum:   "8675309",
	}

	handler := PostRevisionHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPostRevisionWritePerms() {
	testCases := map[string]struct {
		cert   *models.ClientCert
		issuer ordersmessages.Issuer
	}{
		"Army": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerArmy,
		},
		"Navy": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerNavy,
		},
		"MarineCorps": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerMarineDashCorps,
		},
		"CoastGuard": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerCoastDashGuard,
		},
		"AirForce": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerAirDashForce,
		},
	}
	testCases["Army"].cert.AllowArmyOrdersWrite = false
	testCases["Navy"].cert.AllowNavyOrdersWrite = false
	testCases["MarineCorps"].cert.AllowMarineCorpsOrdersWrite = false
	testCases["CoastGuard"].cert.AllowCoastGuardOrdersWrite = false
	testCases["AirForce"].cert.AllowAirForceOrdersWrite = false

	for name, testCase := range testCases {
		suite.T().Run(name, func(_ *testing.T) {
			req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.PostRevisionParams{
				HTTPRequest: req,
				Issuer:      string(testCase.issuer),
				MemberID:    "1234567890",
				OrdersNum:   "8675309",
			}

			handler := PostRevisionHandler{suite.NewHandlerConfig()}
			response := handler.Handle(params)

			suite.IsType(&handlers.ErrResponse{}, response)
			errResponse, ok := response.(*handlers.ErrResponse)
			if !ok {
				return
			}
			suite.Equal(http.StatusForbidden, errResponse.Code)
		})
	}
}
