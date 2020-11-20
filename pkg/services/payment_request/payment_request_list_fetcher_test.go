package paymentrequest

import (
	"sort"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// The default GBLOC is "LKNQ" for office users and payment requests
	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
	})
	testdatagen.MakeDefaultPaymentRequest(suite.DB())

	suite.T().Run("Returns payment requests matching office user GBLOC", func(t *testing.T) {
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.NoError(err)
		suite.Equal(2, len(*expectedPaymentRequests))
	})

	suite.T().Run("Returns payment request matching an arbitrary filter", func(t *testing.T) {
		// Locator
		moveID := paymentRequest.MoveTaskOrder.Locator
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2), MoveID: &moveID})
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
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2), Branch: &branch})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests = *expectedPaymentRequests
		suite.Equal(models.AffiliationAIRFORCE, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListUSMCGBLOC() {
	suite.T().Run("returns USMC payment requests", func(t *testing.T) {
		officeUUID, _ := uuid.NewV4()
		marines := models.AffiliationMARINES
		army := models.AffiliationARMY

		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			TransportationOffice: models.TransportationOffice{
				Gbloc: "USMC",
				ID:    officeUUID,
			},
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
			ServiceMember: models.ServiceMember{Affiliation: &marines},
		})

		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
			ServiceMember: models.ServiceMember{Affiliation: &army},
		})

		officeUserOooRah := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{OfficeUser: models.OfficeUser{TransportationOfficeID: officeUUID}})
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUserOooRah.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(1, len(paymentRequests))
		suite.Equal(models.AffiliationMARINES, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)

		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})
		paymentRequests = *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(1, len(paymentRequests))
		suite.Equal(models.AffiliationARMY, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListNoGBLOCMatch() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
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

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.NoError(err)
		suite.Equal(0, len(*expectedPaymentRequests))
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListFailure() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	suite.T().Run("Error when office user ID does not exist", func(t *testing.T) {
		nonexistentOfficeUserID := uuid.Must(uuid.NewV4())
		_, _, err := paymentRequestListFetcher.FetchPaymentRequestList(nonexistentOfficeUserID,
			&services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(2)})

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching transportationOffice for officeUserID")
		suite.Contains(err.Error(), nonexistentOfficeUserID.String())
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListWithPagination() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	for i := 0; i < 2; i++ {
		testdatagen.MakeDefaultPaymentRequest(suite.DB())
	}

	expectedPaymentRequests, count, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &services.FetchPaymentRequestListParams{Page: swag.Int64(1), PerPage: swag.Int64(1)})

	suite.NoError(err)
	suite.Equal(1, len(*expectedPaymentRequests))
	suite.Equal(2, count)
}

func (suite *PaymentRequestServiceSuite) TestListPaymentRequestWithSortOrder() {

	var expectedNameOrder []string
	var expectedDodIDOrder []string
	var expectedStatusOrder []string
	var expectedCreatedAtOrder []time.Time
	var expectedMoveIDOrder []string
	var expectedBranchOrder []string

	hhgMoveType := models.SelectedMoveTypeHHG
	branchNavy := models.AffiliationNAVY
	//
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Leo"),
			LastName:  models.StringPointer("Spacemen"),
			Edipi:     models.StringPointer("AZFG"),
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Locator:          "ZZZZ",
		},
	})
	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))

	paymentRequest1 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		PaymentRequest: models.PaymentRequest{
			Status:    models.PaymentRequestStatusReviewed,
			CreatedAt: prevCreatedAt,
		},
	})

	paymentRequest2 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Edipi:       models.StringPointer("EZFG"),
			LastName:    models.StringPointer("Spacemen"),
			FirstName:   models.StringPointer("Lena"),
			Affiliation: &branchNavy,
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusAPPROVED,
			Locator:          "AAAA",
		},
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusPaid,
		},
	})

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
	expectedMoveIDOrder = append(expectedMoveIDOrder, paymentRequest1.MoveTaskOrder.Locator, paymentRequest2.MoveTaskOrder.Locator)
	expectedBranchOrder = append(expectedBranchOrder, string(*paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Affiliation), string(*paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Affiliation))

	sort.Strings(expectedNameOrder)
	sort.Strings(expectedDodIDOrder)
	sort.Strings(expectedStatusOrder)
	sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
	sort.Strings(expectedMoveIDOrder)
	sort.Strings(expectedBranchOrder)

	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	// Sort by service member name
	params := services.FetchPaymentRequestListParams{Sort: swag.String("service_members.last_name"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests := *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedNameOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
	suite.Equal(expectedNameOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)

	// Sort by service member name
	params = services.FetchPaymentRequestListParams{Sort: swag.String("service_members.last_name"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedNameOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
	suite.Equal(expectedNameOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)

	// Sort by dodID
	params = services.FetchPaymentRequestListParams{Sort: swag.String("service_members.edipi"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedDodIDOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)
	suite.Equal(expectedDodIDOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)

	params = services.FetchPaymentRequestListParams{Sort: swag.String("service_members.edipi"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedDodIDOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)
	suite.Equal(expectedDodIDOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)

	// Sort by status
	params = services.FetchPaymentRequestListParams{Sort: swag.String("payment_requests.status"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedStatusOrder[0], paymentRequests[0].Status)
	suite.Equal(expectedStatusOrder[1], paymentRequests[1].Status)

	params = services.FetchPaymentRequestListParams{Sort: swag.String("payment_requests.status"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedStatusOrder[0], paymentRequests[1].Status)
	suite.Equal(expectedStatusOrder[1], paymentRequests[0].Status)

	// Sort by Age & Date Submitted
	params = services.FetchPaymentRequestListParams{Sort: swag.String("payment_requests.created_at"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedStatusOrder[0], paymentRequests[0].Status)
	suite.Equal(expectedStatusOrder[1], paymentRequests[1].Status)

	params = services.FetchPaymentRequestListParams{Sort: swag.String("payment_requests.created_at"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedStatusOrder[0], paymentRequests[1].Status)
	suite.Equal(expectedStatusOrder[1], paymentRequests[0].Status)

	// Sort by Move ID
	params = services.FetchPaymentRequestListParams{Sort: swag.String("moves.locator"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedMoveIDOrder[0], paymentRequests[0].Status)
	suite.Equal(expectedMoveIDOrder[1], paymentRequests[1].Status)

	params = services.FetchPaymentRequestListParams{Sort: swag.String("moves.locator"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedMoveIDOrder[0], paymentRequests[1].Status)
	suite.Equal(expectedMoveIDOrder[1], paymentRequests[0].Status)

	// Sort by Branch
	params = services.FetchPaymentRequestListParams{Sort: swag.String("service_members.affiliation"), Order: swag.String("asc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedBranchOrder[0], paymentRequests[0].Status)
	suite.Equal(expectedBranchOrder[1], paymentRequests[1].Status)

	params = services.FetchPaymentRequestListParams{Sort: swag.String("service_members.affiliation"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, &params)
	paymentRequests = *expectedPaymentRequests

	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal(expectedBranchOrder[0], paymentRequests[1].Status)
	suite.Equal(expectedBranchOrder[1], paymentRequests[0].Status)

}
