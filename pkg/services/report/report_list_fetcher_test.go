package report

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *ReportServiceSuite) TestReportFetcher() {
	reportListFetcher := NewReportListFetcher()
	defaultSearchParams := services.MoveFetcherParams{}

	suite.Run("successfully return only navy moves with an approved payment request", func() {
		nonNavyMove := factory.BuildMove(suite.DB(), nil, nil)

		origAffiliation := models.AffiliationNAVY
		navyMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &origAffiliation,
				},
			},
		}, nil)
		pr := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    navyMove,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    pr,
				LinkOnly: true,
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
		}, nil)

		actualMove, err := reportListFetcher.FetchMovesForReports(suite.AppContextForTest(), &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Len(actualMove, 1)
		suite.NotEqual(nonNavyMove.ID, actualMove[0].ID)
		suite.Equal(navyMove.ID, actualMove[0].ID)
	})
}
