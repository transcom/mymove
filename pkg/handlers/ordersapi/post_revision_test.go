package ordersapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestPostRevisionNew() {
	req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:           true,
		AllowAirForceOrdersWrite: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	rev := ordersmessages.Revision{
		SeqNum: 0,
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
		HasDependents: true,
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

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		MemberID:    "1234567890",
		OrdersNum:   "8675309",
		Revision:    &rev,
	}

	handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionCreated{}, response)

	createdResponse, ok := response.(*ordersoperations.PostRevisionCreated)
	if !ok {
		return
	}

	id, err := uuid.FromString(createdResponse.Payload.UUID.String())
	suite.Assertions.NoError(err)

	// check that the order and its revision are actually in the DB
	order, err := models.FetchElectronicOrderByID(suite.DB(), id)
	suite.NoError(err)
	suite.NotNil(order)
	suite.Equal(params.Issuer, string(order.Issuer))
	suite.Equal(params.MemberID, order.Edipi)
	suite.Equal(params.OrdersNum, order.OrdersNumber)
	suite.Len(order.Revisions, 1)
	storedRev := order.Revisions[0]
	suite.EqualValues(rev.SeqNum, storedRev.SeqNum)
	suite.Equal(rev.Member.GivenName, storedRev.GivenName)
	suite.Equal(rev.Member.FamilyName, storedRev.FamilyName)
	suite.Equal(string(rev.Member.Rank), string(storedRev.Paygrade))
	suite.Equal(rev.PcsAccounting.Tac, storedRev.HhgTAC)
	suite.Equal(string(rev.Status), string(storedRev.Status))
	suite.Equal(string(rev.TourType), string(storedRev.TourType))
	suite.Equal(string(rev.OrdersType), string(storedRev.OrdersType))
	suite.Equal(rev.HasDependents, storedRev.HasDependents)
	suite.Equal(rev.NoCostMove, storedRev.NoCostMove)
	suite.Equal(rev.LosingUnit.Uic, storedRev.LosingUIC)
	suite.Equal(rev.LosingUnit.Name, storedRev.LosingUnitName)
	suite.Equal(rev.LosingUnit.City, storedRev.LosingUnitCity)
	suite.Equal(rev.LosingUnit.Locality, storedRev.LosingUnitLocality)
	suite.Equal(rev.LosingUnit.PostalCode, storedRev.LosingUnitPostalCode)
}

func (suite *HandlerSuite) TestPostRevisionNewAmendment() {
	// prime the DB with an order with 1 revision
	origOrder := testdatagen.MakeElectronicOrder(suite.DB(), "1234567890", models.IssuerAirForce, "8675309", models.ElectronicOrdersAffiliationAirForce)

	req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:           true,
		AllowAirForceOrdersWrite: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

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
		HasDependents: true,
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

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		MemberID:    origOrder.Edipi,
		OrdersNum:   origOrder.OrdersNumber,
		Revision:    &rev,
	}

	handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionCreated{}, response)

	createdResponse, ok := response.(*ordersoperations.PostRevisionCreated)
	if !ok {
		return
	}

	id, err := uuid.FromString(createdResponse.Payload.UUID.String())
	suite.Assertions.NoError(err)

	// check that the order and its new revision are actually in the DB
	order, err := models.FetchElectronicOrderByID(suite.DB(), id)
	suite.NoError(err)
	suite.NotNil(order)
	suite.Equal(params.Issuer, string(order.Issuer))
	suite.Equal(params.MemberID, order.Edipi)
	suite.Equal(params.OrdersNum, order.OrdersNumber)
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
	suite.Equal(rev.HasDependents, storedRev.HasDependents)
	suite.Equal(rev.NoCostMove, storedRev.NoCostMove)
	suite.Equal(rev.LosingUnit.Uic, storedRev.LosingUIC)
	suite.Equal(rev.LosingUnit.Name, storedRev.LosingUnitName)
	suite.Equal(rev.LosingUnit.City, storedRev.LosingUnitCity)
	suite.Equal(rev.LosingUnit.Locality, storedRev.LosingUnitLocality)
	suite.Equal(rev.LosingUnit.PostalCode, storedRev.LosingUnitPostalCode)
}

func (suite *HandlerSuite) TestPostRevisionNoApiPerm() {
	req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
	clientCert := models.ClientCert{}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		MemberID:    "1234567890",
		OrdersNum:   "8675309",
	}

	handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&ordersoperations.PostRevisionForbidden{}, response)
}

func (suite *HandlerSuite) TestPostRevisionWritePerms() {
	testCases := map[string]struct {
		cert   models.ClientCert
		issuer ordersmessages.Issuer
	}{
		"Army": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerArmy,
		},
		"Navy": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
			},
			ordersmessages.IssuerNavy,
		},
		"MarineCorps": {
			models.ClientCert{
				AllowOrdersAPI:             true,
				AllowAirForceOrdersWrite:   true,
				AllowArmyOrdersWrite:       true,
				AllowCoastGuardOrdersWrite: true,
				AllowNavyOrdersWrite:       true,
			},
			ordersmessages.IssuerMarineCorps,
		},
		"CoastGuard": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowAirForceOrdersWrite:    true,
				AllowArmyOrdersWrite:        true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerCoastGuard,
		},
		"AirForce": {
			models.ClientCert{
				AllowOrdersAPI:              true,
				AllowArmyOrdersWrite:        true,
				AllowCoastGuardOrdersWrite:  true,
				AllowMarineCorpsOrdersWrite: true,
				AllowNavyOrdersWrite:        true,
			},
			ordersmessages.IssuerAirForce,
		},
	}

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
			req = suite.AuthenticateClientCertRequest(req, &testCase.cert)

			params := ordersoperations.PostRevisionParams{
				HTTPRequest: req,
				Issuer:      string(testCase.issuer),
				MemberID:    "1234567890",
				OrdersNum:   "8675309",
			}

			handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
			response := handler.Handle(params)

			suite.Assertions.IsType(&ordersoperations.PostRevisionForbidden{}, response)
		})
	}
}
