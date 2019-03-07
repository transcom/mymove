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

func (suite *HandlerSuite) TestIndexOrdersForMemberNumSuccess() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerAirForce, "8675309", ordersmessages.AffiliationAirForce)
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
	suite.Len(okResponse.Payload, 1)
	suite.Equal(order.Edipi, okResponse.Payload[0].Edipi)
	suite.Equal(order.Issuer, okResponse.Payload[0].Issuer)
	suite.Equal(order.OrdersNumber, okResponse.Payload[0].OrdersNum)
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

func (suite *HandlerSuite) TestIndexOrdersForMemberNumNoArmyPerm() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerArmy, "8675309", ordersmessages.AffiliationArmy)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
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
	suite.Len(okResponse.Payload, 0)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNoNavyPerm() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerNavy, "8675309", ordersmessages.AffiliationNavy)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
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
	suite.Len(okResponse.Payload, 0)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNoMarineCorpsPerm() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerMarineCorps, "8675309", ordersmessages.AffiliationMarineCorps)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:            true,
		AllowAirForceOrdersRead:   true,
		AllowArmyOrdersRead:       true,
		AllowCoastGuardOrdersRead: true,
		AllowNavyOrdersRead:       true,
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
	suite.Len(okResponse.Payload, 0)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNoAirForcePerm() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerAirForce, "8675309", ordersmessages.AffiliationAirForce)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
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
	suite.Len(okResponse.Payload, 0)
}

func (suite *HandlerSuite) TestIndexOrdersForMemberNoCoastGuardPerm() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerCoastGuard, "8675309", ordersmessages.AffiliationCoastGuard)
	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/v1/edipis/%s/orders", order.Edipi), nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
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
	suite.Len(okResponse.Payload, 0)
}
