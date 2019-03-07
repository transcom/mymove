package ordersapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/ordersmessages"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersSuccess() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerAirForce, ordersmessages.AffiliationAirForce)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:          true,
		AllowAirForceOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersOK{}, response)
	okResponse := response.(*ordersoperations.GetOrdersOK)
	suite.Equal(strfmt.UUID(order.ID.String()), okResponse.Payload.UUID)
}

func (suite *HandlerSuite) TestGetOrdersNoApiPerm() {
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	uuid, _ := uuid.NewV4()
	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(uuid.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersNoArmyPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerArmy, ordersmessages.AffiliationArmy)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersNoNavyPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerNavy, ordersmessages.AffiliationNavy)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersNoMarineCorpsPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerMarineCorps, ordersmessages.AffiliationMarineCorps)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:            true,
		AllowAirForceOrdersRead:   true,
		AllowArmyOrdersRead:       true,
		AllowCoastGuardOrdersRead: true,
		AllowNavyOrdersRead:       true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersNoAirForcePerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerAirForce, ordersmessages.AffiliationAirForce)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersNoCoastGuardPerm() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB(), ordersmessages.IssuerCoastGuard, ordersmessages.AffiliationCoastGuard)
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)

	clientCert := models.ClientCert{
		AllowOrdersAPI:             true,
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(order.ID.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestGetOrdersMissingUUID() {
	req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	uuid, _ := uuid.NewV4()
	params := ordersoperations.GetOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(uuid.String()),
	}

	handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.GetOrdersNotFound{}, response)
}
