package paymentrequest

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/gobuffalo/pop/v5"

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
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, nil, nil)

		suite.NoError(err)
		suite.Equal(2, len(*expectedPaymentRequests))
	})

	suite.T().Run("Returns payment request matching an arbitrary filter", func(t *testing.T) {
		// Locator
		moveID := paymentRequest.MoveTaskOrder.Locator
		moveIDQuery := moveIDFilter(&moveID)
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, nil, nil, moveIDQuery)
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
		branchQuery := branchFilter(&branch)
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, nil, nil, branchQuery)
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests = *expectedPaymentRequests
		suite.Equal(models.AffiliationAIRFORCE, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
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

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, nil, nil)

		suite.NoError(err)
		suite.Equal(0, len(*expectedPaymentRequests))
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListFailure() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	suite.T().Run("Error when office user ID does not exist", func(t *testing.T) {
		nonexistentOfficeUserID := uuid.Must(uuid.NewV4())
		_, _, err := paymentRequestListFetcher.FetchPaymentRequestList(nonexistentOfficeUserID, nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching transportationOffice for officeUserID")
		suite.Contains(err.Error(), nonexistentOfficeUserID.String())
	})
}

type FilterOption func(*pop.Query)

func moveIDFilter(moveID *string) FilterOption {
	return func(query *pop.Query) {
		if moveID != nil {
			query = query.Where("moves.locator = ?", *moveID)
		}
	}
}
func branchFilter(branch *string) FilterOption {
	return func(query *pop.Query) {
		if branch != nil {
			query = query.Where("service_members.affiliation = ?", *branch)
		}
	}
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListWithPagination() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	for i := 0; i < 2; i++ {
		testdatagen.MakeDefaultPaymentRequest(suite.DB())
	}

	expectedPaymentRequests, count, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID, swag.Int(1), swag.Int(1))

	suite.NoError(err)
	suite.Equal(1, len(*expectedPaymentRequests))
	suite.Equal(2, count)

}
