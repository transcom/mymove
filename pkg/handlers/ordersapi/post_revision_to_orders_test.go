package ordersapi

import (
	"fmt"
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

func (suite *HandlerSuite) TestPostRevisionToOrdersNoApiPerm() {
	id, _ := uuid.NewV4()
	req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", id.String()), nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.PostRevisionToOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(id.String()),
	}

	handler := PostRevisionToOrdersHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPostRevisionToOrdersWritePerms() {
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
	testCases["Army"].cert.AllowArmyOrdersWrite = false
	testCases["Navy"].cert.AllowNavyOrdersWrite = false
	testCases["MarineCorps"].cert.AllowMarineCorpsOrdersWrite = false
	testCases["CoastGuard"].cert.AllowCoastGuardOrdersWrite = false
	testCases["AirForce"].cert.AllowAirForceOrdersWrite = false

	for name, testCase := range testCases {
		suite.T().Run(name, func(_ *testing.T) {
			// prime the DB with an order with 1 revision
			assertions := testdatagen.Assertions{
				ElectronicOrder: models.ElectronicOrder{
					Edipi:  testCase.edipi,
					Issuer: testCase.issuer,
				},
				ElectronicOrdersRevision: models.ElectronicOrdersRevision{
					Affiliation: testCase.affl,
				},
			}
			origOrder := testdatagen.MakeElectronicOrder(suite.DB(), assertions)
			req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", origOrder.ID.String()), nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.PostRevisionToOrdersParams{
				HTTPRequest: req,
				UUID:        strfmt.UUID(origOrder.ID.String()),
			}

			handler := PostRevisionToOrdersHandler{suite.NewHandlerConfig()}
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
