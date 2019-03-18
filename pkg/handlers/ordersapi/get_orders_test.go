package ordersapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrdersSuccess() {
	order := testdatagen.MakeDefaultElectronicOrder(suite.DB())
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

	suite.IsType(&ordersoperations.GetOrdersOK{}, response)
	okResponse, ok := response.(*ordersoperations.GetOrdersOK)
	if !ok {
		return
	}
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

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestGetOrdersReadPerms() {
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
			req := httptest.NewRequest("GET", "/orders/v1/orders/", nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.GetOrdersParams{
				HTTPRequest: req,
				UUID:        strfmt.UUID(order.ID.String()),
			}

			handler := GetOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusNotFound, errResponse.Code)
}
