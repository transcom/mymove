package ordersapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumSuccess() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", models.IssuerAirForce, "8675309", models.ElectronicOrdersAffiliationAirForce)
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

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumReadPerms() {
	testCases := map[string]struct {
		cert   models.ClientCert
		issuer models.Issuer
		affl   models.ElectronicOrdersAffiliation
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
			models.IssuerArmy,
			models.ElectronicOrdersAffiliationArmy,
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
			models.IssuerNavy,
			models.ElectronicOrdersAffiliationNavy,
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
			models.IssuerMarineCorps,
			models.ElectronicOrdersAffiliationMarineCorps,
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
			models.IssuerCoastGuard,
			models.ElectronicOrdersAffiliationCoastGuard,
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
			models.IssuerAirForce,
			models.ElectronicOrdersAffiliationAirForce,
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
	issuer := models.IssuerArmy
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
