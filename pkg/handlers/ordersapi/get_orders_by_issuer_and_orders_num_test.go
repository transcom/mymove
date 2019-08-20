package ordersapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumSuccess() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:          true,
		AllowAirForceOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumOK{}, response)
	okResponse, ok := response.(*ordersoperations.GetOrdersByIssuerAndOrdersNumOK)
	if !ok {
		return
	}
	suite.Equal(string(order.Issuer), string(okResponse.Payload.Issuer))
	suite.Equal(order.OrdersNumber, okResponse.Payload.OrdersNum)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoApiPerm() {
	req := httptest.NewRequest("GET", "/orders/v1/issuers/air-force/orders/8675309", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(models.IssuerAirForce),
		OrdersNum:   "8675309",
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumReadPerms() {
	testCases := map[string]struct {
		cert   *models.ClientCert
		issuer models.Issuer
		affl   models.ElectronicOrdersAffiliation
		edipi  string
	}{
		"Army": {
			makeAllPowerfulClientCert(),
			models.IssuerArmy,
			models.ElectronicOrdersAffiliationArmy,
			"1234567890",
		},
		"Navy": {
			makeAllPowerfulClientCert(),
			models.IssuerNavy,
			models.ElectronicOrdersAffiliationNavy,
			"1234567891",
		},
		"MarineCorps": {
			makeAllPowerfulClientCert(),
			models.IssuerMarineCorps,
			models.ElectronicOrdersAffiliationMarineCorps,
			"1234567892",
		},
		"CoastGuard": {
			makeAllPowerfulClientCert(),
			models.IssuerCoastGuard,
			models.ElectronicOrdersAffiliationCoastGuard,
			"1234567893",
		},
		"AirForce": {
			makeAllPowerfulClientCert(),
			models.IssuerAirForce,
			models.ElectronicOrdersAffiliationAirForce,
			"1234567894",
		},
	}
	testCases["Army"].cert.AllowArmyOrdersRead = false
	testCases["Navy"].cert.AllowNavyOrdersRead = false
	testCases["MarineCorps"].cert.AllowMarineCorpsOrdersRead = false
	testCases["CoastGuard"].cert.AllowCoastGuardOrdersRead = false
	testCases["AirForce"].cert.AllowAirForceOrdersRead = false

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			assertions := testdatagen.Assertions{
				ElectronicOrder: models.ElectronicOrder{
					Edipi:  testCase.edipi,
					Issuer: testCase.issuer,
				},
				ElectronicOrdersRevision: models.ElectronicOrdersRevision{
					Affiliation: testCase.affl,
				},
			}
			order := testdatagen.MakeElectronicOrder(suite.DB(), assertions)
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
				HTTPRequest: req,
				Issuer:      string(order.Issuer),
				OrdersNum:   order.OrdersNumber,
			}

			handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNotFound() {
	issuer := models.IssuerArmy
	ordersNum := "notfound"
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(issuer), ordersNum), nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:      true,
		AllowArmyOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(issuer),
		OrdersNum:   ordersNum,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusNotFound, errResponse.Code)
}
