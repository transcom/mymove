package paymentrequest

import (
	"testing"

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
