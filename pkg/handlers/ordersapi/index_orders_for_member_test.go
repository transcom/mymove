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

func (suite *HandlerSuite) TestIndexOrdersForMemberNumSuccess() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", models.IssuerAirForce, "8675309", models.ElectronicOrdersAffiliationAirForce)
	order2 := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", models.IssuerAirForce, "8675310", models.ElectronicOrdersAffiliationAirForce)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:          true,
		AllowAirForceOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.IndexOrdersForMemberParams{
		HTTPRequest: req,
		Edipi:       order.Edipi,
	}

	handler := IndexOrdersForMemberHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.IndexOrdersForMemberOK{}, response)
	okResponse := response.(*ordersoperations.IndexOrdersForMemberOK)
	suite.Len(okResponse.Payload, 2)
	suite.Equal(order.Edipi, okResponse.Payload[0].Edipi)
	suite.Equal(string(order.Issuer), string(okResponse.Payload[0].Issuer))
	suite.Equal(order.Edipi, okResponse.Payload[1].Edipi)
	suite.Equal(string(order.Issuer), string(okResponse.Payload[1].Issuer))
	suite.Contains([]string{order.OrdersNumber, order2.OrdersNumber}, okResponse.Payload[0].OrdersNum)
	suite.Contains([]string{order.OrdersNumber, order2.OrdersNumber}, okResponse.Payload[1].OrdersNum)
	suite.NotEqual(okResponse.Payload[0].OrdersNum, okResponse.Payload[1].OrdersNum)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNumNoApiPerm() {
	req := httptest.NewRequest("GET", "/orders/v1/edipis/1234567890/orders", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.IndexOrdersForMemberParams{
		HTTPRequest: req,
		Edipi:       "1234567890",
	}

	handler := IndexOrdersForMemberHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.IndexOrdersForMemberForbidden{}, response)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNumNoReadPerms() {
	req := httptest.NewRequest("GET", "/orders/v1/edipis/1234567890/orders", nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.IndexOrdersForMemberParams{
		HTTPRequest: req,
		Edipi:       "1234567890",
	}

	handler := IndexOrdersForMemberHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.IndexOrdersForMemberForbidden{}, response)
}

func (suite *HandlerSuite) TestIndexOrderForMemberReadPerms() {
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
			order := testdatagen.MakeElectronicOrder(suite.DB(), testCase.edipi, testCase.issuer, "8675309", testCase.affl)
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.IndexOrdersForMemberParams{
				HTTPRequest: req,
				Edipi:       order.Edipi,
			}

			handler := IndexOrdersForMemberHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.IndexOrdersForMemberOK{}, response)
			okResponse := response.(*ordersoperations.IndexOrdersForMemberOK)
			suite.Len(okResponse.Payload, 0)
		})
	}
}
