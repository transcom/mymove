package ordersapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumSuccess() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerAirForce, "8675309", ordersmessages.AffiliationAirForce)
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
	okResponse := response.(*ordersoperations.GetOrdersByIssuerAndOrdersNumOK)
	suite.Equal(order.Issuer, okResponse.Payload.Issuer)
	suite.Equal(order.OrdersNumber, okResponse.Payload.OrdersNum)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoApiPerm() {
	req := httptest.NewRequest("GET", "/orders/v1/issuers/air-force/orders/8675309", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		OrdersNum:   "8675309",
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumReadPerms() {
	testCases := map[string]struct {
		cert   models.ClientCert
		issuer ordersmessages.Issuer
		affl   ordersmessages.Affiliation
		edipi  string
	}{
		"Army": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersRead:    true,
				AllowCoastGuardOrdersRead:  true,
				AllowMarineCorpsOrdersRead: true,
				AllowNavyOrdersRead:        true,
			},
			ordersmessages.IssuerArmy,
			ordersmessages.AffiliationArmy,
			"1234567890",
		},
		"Navy": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersRead:    true,
				AllowArmyOrdersRead:        true,
				AllowCoastGuardOrdersRead:  true,
				AllowMarineCorpsOrdersRead: true,
			},
			ordersmessages.IssuerNavy,
			ordersmessages.AffiliationNavy,
			"1234567891",
		},
		"MarineCorps": {
			models.ClientCert{
				AllowOrdersAPI:            true,
				AllowAirForceOrdersRead:   true,
				AllowArmyOrdersRead:       true,
				AllowCoastGuardOrdersRead: true,
				AllowNavyOrdersRead:       true,
			},
			ordersmessages.IssuerMarineCorps,
			ordersmessages.AffiliationMarineCorps,
			"1234567892",
		},
		"CoastGuard": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersRead:    true,
				AllowArmyOrdersRead:        true,
				AllowMarineCorpsOrdersRead: true,
				AllowNavyOrdersRead:        true,
			},
			ordersmessages.IssuerCoastGuard,
			ordersmessages.AffiliationCoastGuard,
			"1234567893",
		},
		"AirForce": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowArmyOrdersRead:        true,
				AllowCoastGuardOrdersRead:  true,
				AllowMarineCorpsOrdersRead: true,
				AllowNavyOrdersRead:        true,
			},
			ordersmessages.IssuerAirForce,
			ordersmessages.AffiliationAirForce,
			"1234567894",
		},
	}

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			order := testdatagen.MakeElectronicOrder(suite.DB(), testCase.edipi, testCase.issuer, "8675309", testCase.affl)
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)
			req = suite.AuthenticateClientCertRequest(req, &testCase.cert)

			params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
				HTTPRequest: req,
				Issuer:      string(order.Issuer),
				OrdersNum:   order.OrdersNumber,
			}

			handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
		})
	}
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNotFound() {
	issuer := ordersmessages.IssuerArmy
	ordersNum := "notfound"
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(issuer), ordersNum), nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(issuer),
		OrdersNum:   ordersNum,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumNotFound{}, response)
}
