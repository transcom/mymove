package ordersapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestPostRevision() {
	req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
	clientCert := models.ClientCert{
		AllowOrdersAPI:           true,
		AllowAirForceOrdersWrite: true,
	}
	req = suite.AuthenticateClientCertRequest(req, &clientCert)

	seqNum := int64(0)
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

	params := ordersoperations.PostRevisionParams{
		HTTPRequest: req,
		Issuer:      string(ordersmessages.IssuerAirForce),
		MemberID:    "1234567890",
		OrdersNum:   "8675309",
		Revision:    &rev,
	}

	handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	suite.T().Run("NewSuccess", func(t *testing.T) {
		response := handler.Handle(params)

		suite.IsType(&ordersoperations.PostRevisionCreated{}, response)
		createdResponse, ok := response.(*ordersoperations.PostRevisionCreated)
		if !ok {
			return
		}
		id, err := uuid.FromString(createdResponse.Payload.UUID.String())
		suite.NoError(err)

		// check that the order and its revision are actually in the DB
		order, err := models.FetchElectronicOrderByID(suite.DB(), id)
		suite.NoError(err)
		suite.NotNil(order)
		suite.Equal(params.Issuer, string(order.Issuer))
		suite.Equal(params.MemberID, order.Edipi)
		suite.Equal(params.OrdersNum, order.OrdersNumber)
		suite.Len(order.Revisions, 1)
		storedRev := order.Revisions[0]
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

	suite.T().Run("AmendmentSuccess", func(t *testing.T) {
		seqNum = int64(1)
		response := handler.Handle(params)

		suite.IsType(&ordersoperations.PostRevisionCreated{}, response)
		createdResponse, ok := response.(*ordersoperations.PostRevisionCreated)
		if !ok {
			return
		}
		id, err := uuid.FromString(createdResponse.Payload.UUID.String())
		suite.NoError(err)

		// check that the order and its new revision are actually in the DB
		order, err := models.FetchElectronicOrderByID(suite.DB(), id)
		suite.NoError(err)
		suite.NotNil(order)
		suite.Equal(params.Issuer, string(order.Issuer))
		suite.Equal(params.MemberID, order.Edipi)
		suite.Equal(params.OrdersNum, order.OrdersNumber)
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

	suite.T().Run("EdipiConflict", func(t *testing.T) {
		params.MemberID = "9999999999"
		seqNum = int64(99999)
		response := handler.Handle(params)
		suite.IsType(&handlers.ErrResponse{}, response)
		errResponse, ok := response.(*handlers.ErrResponse)
		if !ok {
			return
		}
		suite.Equal(http.StatusConflict, errResponse.Code)
	})
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

	suite.IsType(&handlers.ErrResponse{}, response)
	errResponse, ok := response.(*handlers.ErrResponse)
	if !ok {
		return
	}
	suite.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPostRevisionWritePerms() {
	testCases := map[string]struct {
		cert   *models.ClientCert
		issuer ordersmessages.Issuer
	}{
		"Army": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerArmy,
		},
		"Navy": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerNavy,
		},
		"MarineCorps": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerMarineCorps,
		},
		"CoastGuard": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerCoastGuard,
		},
		"AirForce": {
			makeAllPowerfulClientCert(),
			ordersmessages.IssuerAirForce,
		},
	}
	testCases["Army"].cert.AllowArmyOrdersWrite = false
	testCases["Navy"].cert.AllowNavyOrdersWrite = false
	testCases["MarineCorps"].cert.AllowMarineCorpsOrdersWrite = false
	testCases["CoastGuard"].cert.AllowCoastGuardOrdersWrite = false
	testCases["AirForce"].cert.AllowAirForceOrdersWrite = false

	for name, testCase := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/orders/v1/orders", nil)
			req = suite.AuthenticateClientCertRequest(req, testCase.cert)

			params := ordersoperations.PostRevisionParams{
				HTTPRequest: req,
				Issuer:      string(testCase.issuer),
				MemberID:    "1234567890",
				OrdersNum:   "8675309",
			}

			handler := PostRevisionHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
