package ordersapi

import (
	"fmt"
	"net/http"
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

func (suite *HandlerSuite) TestPostRevisionToOrders() {
	// prime the DB with an order with 1 revision
	origOrder := testdatagen.MakeDefaultElectronicOrder(suite.DB())

	req := httptest.NewRequest("POST", fmt.Sprintf("/orders/v1/orders/%s", origOrder.ID), nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:           true,
		AllowAirForceOrdersWrite: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	seqNum := int64(1)
	hasDependents := true
	rev := ordersmessages.Revision{
		SeqNum: &seqNum,
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
	suite.T().Run("Success", func(t *testing.T) {
		response := handler.Handle(params)

		suite.IsType(&ordersoperations.PostRevisionToOrdersCreated{}, response)
		createdResponse, ok := response.(*ordersoperations.PostRevisionToOrdersCreated)
		if !ok {
			return
		}
		id, err := uuid.FromString(createdResponse.Payload.UUID.String())
		suite.NoError(err)

		// check that the order and its new revision are actually in the DB
		order, err := models.FetchElectronicOrderByID(suite.DB(), id)
		suite.NoError(err)
		suite.NotNil(order)
		suite.Len(order.Revisions, 2)
		storedRev := order.Revisions[1]
		suite.EqualValues(*rev.SeqNum, storedRev.SeqNum)
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
	})

	suite.T().Run("SeqNumConflict", func(t *testing.T) {
		// Sending the amendment again should result in a conflict because the SeqNum will be taken
		response := handler.Handle(params)
		suite.IsType(&handlers.ErrResponse{}, response)
		errResponse, ok := response.(*handlers.ErrResponse)
		if !ok {
			return
		}
		suite.Equal(http.StatusConflict, errResponse.Code)
	})
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
		suite.T().Run(name, func(t *testing.T) {
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

			handler := PostRevisionToOrdersHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
