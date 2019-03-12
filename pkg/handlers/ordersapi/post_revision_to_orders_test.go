package ordersapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestPostRevisionToOrdersNewAmendment() {
	// prime the DB with an order with 1 revision
	origOrder := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", models.IssuerAirForce, "8675309", models.ElectronicOrdersAffiliationAirForce)

	req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", origOrder.ID), nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:           true,
		AllowAirForceOrdersWrite: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	hasDependents := true
	rev := ordersmessages.Revision{
		SeqNum: 1,
		Member: &ordersmessages.Member{
			GivenName:   "First",
			FamilyName:  "Last",
			Affiliation: ordersmessages.AffiliationAirForce,
			Rank:        ordersmessages.RankW1,
		},
		Status:        ordersmessages.StatusAuthorized,
		DateIssued:    handlers.FmtDateTime(time.Now()),
		NoCostMove:    false,
		TdyEnRoute:    false,
		TourType:      ordersmessages.TourTypeAccompanied,
		OrdersType:    ordersmessages.OrdersTypeSeparation,
		HasDependents: &hasDependents,
		LosingUnit: &ordersmessages.Unit{
			Uic:        handlers.FmtString("FFFS00"),
			Name:       handlers.FmtString("SPC721 COMMUNICATIONS SQ"),
			City:       handlers.FmtString("CHEYENNE MTN"),
			Locality:   handlers.FmtString("CO"),
			PostalCode: handlers.FmtString("80914"),
		},
		PcsAccounting: &ordersmessages.Accounting{
			Tac: handlers.FmtString("F67C"),
		},
	}

	params := ordersoperations.PostRevisionToOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(origOrder.ID.String()),
		Revision:    &rev,
	}

	handler := PostRevisionToOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionToOrdersCreated{}, response)

	createdResponse, ok := response.(*ordersoperations.PostRevisionToOrdersCreated)
	if !ok {
		return
	}

	id, err := uuid.FromString(createdResponse.Payload.UUID.String())
	suite.Assertions.NoError(err)

	// check that the order and its new revision are actually in the DB
	order, err := models.FetchElectronicOrderByID(suite.DB(), id)
	suite.NoError(err)
	suite.NotNil(order)
	suite.Len(order.Revisions, 2)
	storedRev := order.Revisions[1]
	suite.EqualValues(rev.SeqNum, storedRev.SeqNum)
	suite.Equal(rev.Member.GivenName, storedRev.GivenName)
	suite.Equal(rev.Member.FamilyName, storedRev.FamilyName)
	suite.Equal(string(rev.Member.Rank), string(storedRev.Paygrade))
	suite.Equal(rev.PcsAccounting.Tac, storedRev.HhgTAC)
	suite.Equal(string(rev.Status), string(storedRev.Status))
	suite.Equal(string(rev.TourType), string(storedRev.TourType))
	suite.Equal(string(rev.OrdersType), string(storedRev.OrdersType))
	suite.Equal(*rev.HasDependents, storedRev.HasDependents)
	suite.Equal(rev.NoCostMove, storedRev.NoCostMove)
	suite.Equal(rev.LosingUnit.Uic, storedRev.LosingUIC)
	suite.Equal(rev.LosingUnit.Name, storedRev.LosingUnitName)
	suite.Equal(rev.LosingUnit.City, storedRev.LosingUnitCity)
	suite.Equal(rev.LosingUnit.Locality, storedRev.LosingUnitLocality)
	suite.Equal(rev.LosingUnit.PostalCode, storedRev.LosingUnitPostalCode)
}

func (suite *HandlerSuite) TestPostRevisionToOrdersNoApiPerm() {
	id, _ := uuid.NewV4()
	req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", id.String()), nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.PostRevisionToOrdersParams{
		HTTPRequest: req,
		UUID:        strfmt.UUID(id.String()),
	}

	handler := PostRevisionToOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionToOrdersForbidden{}, response)
}

func (suite *HandlerSuite) TestPostRevisionToOrdersWritePerms() {
	testCases := map[string]struct {
		cert   models.ClientCert
		issuer models.Issuer
		affl   models.ElectronicOrdersAffiliation
		edipi  string
	}{
		"Army": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			models.IssuerArmy,
			models.ElectronicOrdersAffiliationArmy,
			"1234567890",
		},
		"Navy": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
			},
			models.IssuerNavy,
			models.ElectronicOrdersAffiliationNavy,
			"1234567891",
		},
		"MarineCorps": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersWrite:   true,
				AllowArmyOrdersWrite:       true,
				AllowCoastGuardOrdersWrite: true,
				AllowNavyOrdersWrite:       true,
			},
			models.IssuerMarineCorps,
			models.ElectronicOrdersAffiliationMarineCorps,
			"1234567892",
		},
		"CoastGuard": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			models.IssuerCoastGuard,
			models.ElectronicOrdersAffiliationCoastGuard,
			"1234567893",
		},
		"AirForce": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			models.IssuerAirForce,
			models.ElectronicOrdersAffiliationAirForce,
			"1234567894",
		},
	}

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			// prime the DB with an order with 1 revision
			origOrder := testdatagen.MakeElectronicOrder(suite.DB(), testCase.edipi, testCase.issuer, "8675309", testCase.affl)
			req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", origOrder.ID.String()), nil)
			req = suite.AuthenticateClientCertRequest(req, &testCase.cert)

			params := ordersoperations.PostRevisionToOrdersParams{
				HTTPRequest: req,
				UUID:        strfmt.UUID(origOrder.ID.String()),
			}

			handler := PostRevisionToOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.PostRevisionToOrdersForbidden{}, response)
		})
	}
}
