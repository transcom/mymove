package paymentrequest

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListbyMove() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	expectedMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{Locator: "ABC123"},
	})

	// we need a mapping for the pickup address postal code to our user's gbloc
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	// We need a payment request with a move that has a shipment that's within the GBLOC
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrderID: expectedMove.ID,
			MoveTaskOrder:   expectedMove,
		},
	})

	// Hidden move should not be returned
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Show: swag.Bool(false),
		},
	})

	suite.T().Run("Only returns visible (where Move.Show is not false) payment requests matching office user GBLOC", func(t *testing.T) {
		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestListByMove(suite.AppContextForTest(), officeUser.ID, "ABC123")

		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequestsToCheck := *expectedPaymentRequests
		suite.Equal(paymentRequest.ID, paymentRequestsToCheck[0].ID)
	})

}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	expectedMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	// we need a mapping for the pickup address postal code to our user's gbloc
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	// We need a payment request with a move that has a shipment that's within the GBLOC
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrderID: expectedMove.ID,
			MoveTaskOrder:   expectedMove,
		},
	})

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
		OriginDutyLocation: models.DutyLocation{
			Name: "KJKJKJKJKJK",
		},
	})

	// Hidden move should not be returned
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Show: swag.Bool(false),
		},
	})
	// Marine Corps payment requests should be excluded even if in the same GBLOC
	marines := models.AffiliationMARINES
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "LKNQ",
			ID:    uuid.Must(uuid.NewV4()),
		},
		ServiceMember: models.ServiceMember{Affiliation: &marines},
	})

	suite.T().Run("Only returns visible (where Move.Show is not false) payment requests matching office user GBLOC", func(t *testing.T) {
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))

		paymentRequestsForComparison := *expectedPaymentRequests
		suite.Equal(paymentRequest.ID, paymentRequestsForComparison[0].ID)
	})

	suite.T().Run("Returns payment request matching an arbitrary filter", func(t *testing.T) {
		// Locator
		locator := paymentRequest.MoveTaskOrder.Locator
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2), Locator: &locator})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests := *expectedPaymentRequests
		suite.Equal(paymentRequest.MoveTaskOrder.Locator, paymentRequests[0].MoveTaskOrder.Locator)

		// Branch
		serviceMember := paymentRequest.MoveTaskOrder.Orders.ServiceMember
		affiliation := models.AffiliationAIRFORCE
		serviceMember.Affiliation = &affiliation
		err = suite.DB().Save(&serviceMember)
		suite.NoError(err)

		branch := serviceMember.Affiliation.String()
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2), Branch: &branch})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests = *expectedPaymentRequests
		suite.Equal(models.AffiliationAIRFORCE, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})

	suite.T().Run("Returns payment request matching the originDutyLocation filter", func(t *testing.T) {
		locationName := paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.Name

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2), OriginDutyLocation: &locationName})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests := *expectedPaymentRequests
		suite.Equal(locationName, paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name)

	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListStatusFilter() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove3 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove4 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove5 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove6 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove1.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusReviewed,
			MoveTaskOrderID: expectedMove1.ID,
			MoveTaskOrder:   expectedMove1,
		},
	})

	rejectedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			MoveTaskOrderID: expectedMove2.ID,
			MoveTaskOrder:   expectedMove2,
		},
	})

	sentToGexPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusSentToGex,
			MoveTaskOrderID: expectedMove3.ID,
			MoveTaskOrder:   expectedMove3,
		},
	})
	recByGexPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusReceivedByGex,
			MoveTaskOrderID: expectedMove4.ID,
			MoveTaskOrder:   expectedMove4,
		},
	})
	paidPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusPaid,
			MoveTaskOrderID: expectedMove5.ID,
			MoveTaskOrder:   expectedMove5,
		},
	})

	pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusPending,
			MoveTaskOrderID: expectedMove6.ID,
			MoveTaskOrder:   expectedMove6,
		},
	})

	allPaymentRequests := []models.PaymentRequest{pendingPaymentRequest, reviewedPaymentRequest, rejectedPaymentRequest, sentToGexPaymentRequest, recByGexPaymentRequest, paidPaymentRequest}

	suite.T().Run("Returns all payment requests when no status filter is specified", func(t *testing.T) {
		_, actualCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{})
		suite.NoError(err)
		suite.Equal(len(allPaymentRequests), actualCount)
	})

	suite.T().Run("Returns all payment requests when all status filters are selected", func(t *testing.T) {
		_, actualCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Payment requested", "Reviewed", "Rejected", "Paid"}})
		suite.NoError(err)
		suite.Equal(len(allPaymentRequests), actualCount)
	})

	suite.T().Run("Returns only those payment requests with the exact status", func(t *testing.T) {
		pendingPaymentRequests, pendingCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Payment requested"}})
		pending := *pendingPaymentRequests
		suite.NoError(err)
		suite.Equal(1, pendingCount)
		suite.Equal(pendingPaymentRequest.ID, pending[0].ID)

		reviewedPaymentRequests, reviewedCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Reviewed"}})
		reviewed := *reviewedPaymentRequests
		suite.NoError(err)
		suite.Equal(3, reviewedCount)

		reviewedIDs := []uuid.UUID{reviewedPaymentRequest.ID, sentToGexPaymentRequest.ID, recByGexPaymentRequest.ID}
		for _, pr := range reviewed {
			suite.Contains(reviewedIDs, pr.ID)
		}

		rejectedPaymentRequests, rejectedCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Rejected"}})
		rejected := *rejectedPaymentRequests
		suite.NoError(err)
		suite.Equal(1, rejectedCount)
		suite.Equal(rejectedPaymentRequest.ID, rejected[0].ID)

		paidPaymentRequests, paidCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Paid"}})
		paid := *paidPaymentRequests
		suite.NoError(err)
		suite.Equal(1, paidCount)
		suite.Equal(paidPaymentRequest.ID, paid[0].ID)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListUSMCGBLOC() {
	officeUUID, _ := uuid.NewV4()
	marines := models.AffiliationMARINES
	army := models.AffiliationARMY
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	expectedMoveNotUSMC := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMoveNotUSMC.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	paymentRequestUSMC := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "LKNQ",
			ID:    officeUUID,
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		ServiceMember: models.ServiceMember{Affiliation: &marines},
	})

	paymentRequestUSMC2 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			SequenceNumber: 2,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "LKNQ",
			ID:    officeUUID,
		},
		Move:          paymentRequestUSMC.MoveTaskOrder,
		ServiceMember: models.ServiceMember{Affiliation: &marines},
	})

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusPending,
			MoveTaskOrderID: expectedMoveNotUSMC.ID,
			MoveTaskOrder:   expectedMoveNotUSMC,
		},
		ServiceMember: models.ServiceMember{Affiliation: &army},
	})

	officeUserUSMC := testdatagen.MakeOfficeUserWithUSMCGBLOC(suite.DB())

	suite.T().Run("returns USMC payment requests", func(t *testing.T) {
		paymentRequestListFetcher := NewPaymentRequestListFetcher()
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUserUSMC.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(models.AffiliationMARINES, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
		suite.Equal(models.AffiliationMARINES, *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation)
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})
		paymentRequests = *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(1, len(paymentRequests))
		suite.Equal(models.AffiliationARMY, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})

	suite.T().Run("returns USMC payment requests for move", func(t *testing.T) {
		paymentRequestListFetcher := NewPaymentRequestListFetcher()
		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestListByMove(suite.AppContextForTest(), officeUserUSMC.ID, paymentRequestUSMC.MoveTaskOrder.Locator)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(paymentRequestUSMC.ID, paymentRequests[0].ID)
		suite.Equal(paymentRequestUSMC2.ID, paymentRequests[1].ID)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListNoGBLOCMatch() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	suite.T().Run("No results when GBLOC does not match", func(t *testing.T) {
		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				Gbloc: "EFGH",
			},
		})
		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				Gbloc: "ABCD",
			},
		})

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.NoError(err)
		suite.Equal(0, len(*expectedPaymentRequests))
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListFailure() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()

	suite.T().Run("Error when office user ID does not exist", func(t *testing.T) {
		nonexistentOfficeUserID := uuid.Must(uuid.NewV4())
		_, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), nonexistentOfficeUserID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching transportationOffice for officeUserID")
		suite.Contains(err.Error(), nonexistentOfficeUserID.String())
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListWithPagination() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	expectedMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	_ = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusPending,
			MoveTaskOrderID: expectedMove1.ID,
			MoveTaskOrder:   expectedMove1,
		},
	})
	_ = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusPending,
			MoveTaskOrderID: expectedMove2.ID,
			MoveTaskOrder:   expectedMove2,
		},
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove1.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	expectedPaymentRequests, count, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(1)})

	suite.NoError(err)
	suite.Equal(1, len(*expectedPaymentRequests))
	suite.Equal(2, count)

}

func (suite *PaymentRequestServiceSuite) TestListPaymentRequestWithSortOrder() {

	var expectedNameOrder []string
	var expectedDodIDOrder []string
	var expectedStatusOrder []string
	var expectedCreatedAtOrder []time.Time
	var expectedLocatorOrder []string
	var expectedBranchOrder []string
	var expectedOriginDutyLocation []string

	hhgMoveType := models.SelectedMoveTypeHHG
	branchNavy := models.AffiliationNAVY
	//
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	originDutyLocation1 := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			Name: "Applewood, CA 99999",
		},
	})

	originDutyLocation2 := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			Name: "Scott AFB",
		},
	})

	expectedMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Edipi:       models.StringPointer("EZFG"),
			LastName:    models.StringPointer("Spacemen"),
			FirstName:   models.StringPointer("Lena"),
			Affiliation: &branchNavy,
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Locator:          "AAAA",
		},
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusPaid,
		},
		Order: models.Order{
			OriginDutyLocationID: &originDutyLocation1.ID,
			OriginDutyLocation:   &originDutyLocation1,
		},
	})

	expectedMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Leo"),
			LastName:  models.StringPointer("Spacemen"),
			Edipi:     models.StringPointer("AZFG"),
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Locator:          "ZZZZ",
		},
		Order: models.Order{
			OriginDutyLocationID: &originDutyLocation2.ID,
			OriginDutyLocation:   &originDutyLocation2,
		},
	})

	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))
	paymentRequest1 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrderID: expectedMove1.ID,
			MoveTaskOrder:   expectedMove1,
			Status:          models.PaymentRequestStatusPending,
			CreatedAt:       prevCreatedAt,
		},
	})
	paymentRequest2 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status:          models.PaymentRequestStatusReviewed,
			MoveTaskOrderID: expectedMove2.ID,
			MoveTaskOrder:   expectedMove2,
		},
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		expectedMove1.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: paymentRequest2.MoveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	expectedNameOrder = append(expectedNameOrder, *paymentRequest1.MoveTaskOrder.Orders.ServiceMember.FirstName, *paymentRequest2.MoveTaskOrder.Orders.ServiceMember.FirstName)
	expectedDodIDOrder = append(expectedDodIDOrder, *paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Edipi, *paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Edipi)
	expectedStatusOrder = append(expectedStatusOrder, string(paymentRequest1.Status), string(paymentRequest2.Status))
	expectedCreatedAtOrder = append(expectedCreatedAtOrder, paymentRequest1.CreatedAt, paymentRequest2.CreatedAt)
	expectedLocatorOrder = append(expectedLocatorOrder, paymentRequest1.MoveTaskOrder.Locator, paymentRequest2.MoveTaskOrder.Locator)
	expectedBranchOrder = append(expectedBranchOrder, string(*paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Affiliation), string(*paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Affiliation))
	expectedOriginDutyLocation = append(expectedOriginDutyLocation, string(paymentRequest1.MoveTaskOrder.Orders.OriginDutyLocation.Name), string(paymentRequest2.MoveTaskOrder.Orders.OriginDutyLocation.Name))

	paymentRequestListFetcher := NewPaymentRequestListFetcher()

	suite.T().Run("Sort by service member name ASC", func(t *testing.T) {
		sort.Strings(expectedNameOrder)

		params := services.FetchPaymentRequestListParams{Sort: swag.String("lastName"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedNameOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
		suite.Equal(expectedNameOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
	})

	suite.T().Run("Sort by service member name DESC", func(t *testing.T) {
		sort.Strings(expectedNameOrder)

		// Sort by service member name
		params := services.FetchPaymentRequestListParams{Sort: swag.String("lastName"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedNameOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
		suite.Equal(expectedNameOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
	})

	suite.T().Run("Sort by dodID ASC", func(t *testing.T) {
		sort.Strings(expectedDodIDOrder)

		// Sort by dodID
		params := services.FetchPaymentRequestListParams{Sort: swag.String("dodID"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedDodIDOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)
		suite.Equal(expectedDodIDOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)
	})

	suite.T().Run("Sort by dodID DESC", func(t *testing.T) {
		sort.Strings(expectedDodIDOrder)

		params := services.FetchPaymentRequestListParams{Sort: swag.String("dodID"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedDodIDOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)
		suite.Equal(expectedDodIDOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)
	})

	suite.T().Run("Sort by status ASC", func(t *testing.T) {
		params := services.FetchPaymentRequestListParams{Sort: swag.String("status"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedStatusOrder[0], string(paymentRequests[0].Status))
		suite.Equal(expectedStatusOrder[1], string(paymentRequests[1].Status))
	})

	suite.T().Run("Sort by status DESC", func(t *testing.T) {
		params := services.FetchPaymentRequestListParams{Sort: swag.String("status"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedStatusOrder[0], string(paymentRequests[1].Status))
		suite.Equal(expectedStatusOrder[1], string(paymentRequests[0].Status))
	})

	suite.T().Run("Sort by age ASC", func(t *testing.T) {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: swag.String("age"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.T().Run("Sort by age DESC", func(t *testing.T) {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: swag.String("age"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.T().Run("Sort by submittedAt ASC", func(t *testing.T) {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: swag.String("submittedAt"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.T().Run("Sort by submittedAt DESC", func(t *testing.T) {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: swag.String("submittedAt"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.T().Run("Sort by locator ASC", func(t *testing.T) {
		sort.Strings(expectedLocatorOrder)
		params := services.FetchPaymentRequestListParams{Sort: swag.String("locator"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedLocatorOrder[0], strings.TrimSpace(paymentRequests[0].MoveTaskOrder.Locator))
		suite.Equal(expectedLocatorOrder[1], strings.TrimSpace(paymentRequests[1].MoveTaskOrder.Locator))
	})

	suite.T().Run("Sort by locator DESC", func(t *testing.T) {
		sort.Strings(expectedLocatorOrder)

		params := services.FetchPaymentRequestListParams{Sort: swag.String("locator"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedLocatorOrder[0], strings.TrimSpace(paymentRequests[1].MoveTaskOrder.Locator))
		suite.Equal(expectedLocatorOrder[1], strings.TrimSpace(paymentRequests[0].MoveTaskOrder.Locator))
	})

	suite.T().Run("Sort by branch ASC", func(t *testing.T) {
		sort.Strings(expectedBranchOrder)
		params := services.FetchPaymentRequestListParams{Sort: swag.String("branch"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedBranchOrder[0], string(*paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation))
		suite.Equal(expectedBranchOrder[1], string(*paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation))
	})

	suite.T().Run("Sort by branch DESC", func(t *testing.T) {
		sort.Strings(expectedBranchOrder)
		params := services.FetchPaymentRequestListParams{Sort: swag.String("branch"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedBranchOrder[0], string(*paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation))
		suite.Equal(expectedBranchOrder[1], string(*paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation))
	})

	suite.T().Run("Sort by originDutyLocation ASC", func(t *testing.T) {
		sort.Strings(expectedOriginDutyLocation)
		params := services.FetchPaymentRequestListParams{Sort: swag.String("originDutyLocation"), Order: swag.String("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		suite.NoError(err)

		paymentRequests := *expectedPaymentRequests
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedOriginDutyLocation[0], string(paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name))
		suite.Equal(expectedOriginDutyLocation[1], string(paymentRequests[1].MoveTaskOrder.Orders.OriginDutyLocation.Name))
	})

	suite.T().Run("Sort by originDutyLocation DESC", func(t *testing.T) {
		sort.Strings(expectedOriginDutyLocation)

		params := services.FetchPaymentRequestListParams{Sort: swag.String("originDutyLocation"), Order: swag.String("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedOriginDutyLocation[0], string(paymentRequests[1].MoveTaskOrder.Orders.OriginDutyLocation.Name))
		suite.Equal(expectedOriginDutyLocation[1], string(paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name))
	})
}
