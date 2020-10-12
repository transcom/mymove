package paymentrequest

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	suite.T().Run("Returns payment requests matching office user GBLOC", func(t *testing.T) {
		// The default GBLOC is "LKNQ" for office users and payment requests
		testdatagen.MakeDefaultPaymentRequest(suite.DB())
		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			TransportationOffice: models.TransportationOffice{
				Gbloc: "ABCD",
			},
		})

		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID)

		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
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

		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestList(officeUser.ID)

		suite.NoError(err)
		suite.Equal(0, len(*expectedPaymentRequests))
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListFailure() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	suite.T().Run("Error when office user ID does not exist", func(t *testing.T) {
		nonexistentOfficeUserID := uuid.Must(uuid.NewV4())
		_, err := paymentRequestListFetcher.FetchPaymentRequestList(nonexistentOfficeUserID)

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching transportationOffice for officeUserID")
		suite.Contains(err.Error(), nonexistentOfficeUserID.String())
	})
}
