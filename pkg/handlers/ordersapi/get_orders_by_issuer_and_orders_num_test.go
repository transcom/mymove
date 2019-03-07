package ordersapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumSuccess() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerAirForce, ordersmessages.AffiliationAirForce)
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

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoArmyPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerArmy, ordersmessages.AffiliationArmy)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoNavyPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerNavy, ordersmessages.AffiliationNavy)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoMarineCorpsPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerMarineCorps, ordersmessages.AffiliationMarineCorps)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:            true,
		AllowAirForceOrdersRead:   true,
		AllowArmyOrdersRead:       true,
		AllowCoastGuardOrdersRead: true,
		AllowNavyOrdersRead:       true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoAirForcePerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerAirForce, ordersmessages.AffiliationAirForce)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersByIssuerAndOrdersNumNoCoastGuardPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerCoastGuard, ordersmessages.AffiliationCoastGuard)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/issuers/%s/orders/%s", string(order.Issuer), order.OrdersNumber), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersByIssuerAndOrdersNumParams{
		HTTPRequest: req,
		Issuer:      string(order.Issuer),
		OrdersNum:   order.OrdersNumber,
	}

	handler := GetOrdersByIssuerAndOrdersNumHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersByIssuerAndOrdersNumForbidden{}, response)
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
