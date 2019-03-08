package ordersapi

import (
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/gen/ordersmessages"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersSuccess() {
	order := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", ordersmessages.IssuerAirForce, "8675309", ordersmessages.AffiliationAirForce)
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

func (suite *HandlerSuite) TestGetOrdersReadPerms() {
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
			req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)
			req = suite.AuthenticateClientCertRequest(req, &testCase.cert)

			params := ordersoperations.GetOrdersParams{
				HTTPRequest: req,
				UUID:        strfmt.UUID(order.ID.String()),
			}

			handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.GetOrdersForbidden{}, response)
		})
	}
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
