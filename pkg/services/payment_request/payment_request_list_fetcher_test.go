package paymentrequest

import (
	"fmt"
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
	hhgMoveType := models.SelectedMoveTypeHHG
	branchNavy := models.AffiliationNAVY
	//
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Lena"),
			LastName:  models.StringPointer("Spacemen"),
			Edipi: models.StringPointer("AZFG"),

		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
			Locator:          "ZZZZ",
		},
	})
	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		PaymentRequest: models.PaymentRequest{
			CreatedAt: prevCreatedAt,
		},
	})

	paymentRequest2 := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Edipi: models.StringPointer("EZFG"),
			LastName:  models.StringPointer("Spacemen"),
			FirstName: models.StringPointer("Leo"),
			Affiliation: &branchNavy,
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusAPPROVED,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: paymentRequest2.MoveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	//
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	// Sort by service member name
	params := &services.FetchPaymentRequestListParams{Page: swag.Int64(1), Sort: swag.String("last_name"), Order: swag.String("desc")}
	expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, params)
	paymentRequests := *expectedPaymentRequests

	fmt.Printf(" first re %+v", paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
	fmt.Printf("second re %+v", paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
	suite.NoError(err)
	suite.Equal(2, len(paymentRequests))
	suite.Equal("Lena", *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName )
}
