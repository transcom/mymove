package report

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mocks "github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *ReportServiceSuite) TestReportFetcher() {
	ppmEstimator := mocks.PPMEstimator{}

	reportListFetcher := NewReportListFetcher(&ppmEstimator)
	defaultSearchParams := services.MoveTaskOrderFetcherParams{}

	appCtx := suite.AppContextForTest()

	// Setup data
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				MiddleName:  models.StringPointer("O"),
				Affiliation: (*models.ServiceMemberAffiliation)(internalmessages.AffiliationNAVY.Pointer()),
			},
		},
	}, nil)
	orders := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model: models.Order{
				TAC: models.StringPointer("CACI"),
			},
		},
	}, nil)
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
	}, nil)

	pr := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				Status:          models.PaymentRequestStatusReviewed,
				MoveTaskOrderID: move.ID,
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PaymentRequestID: pr.ID,
			},
		},
	}, nil)

	// Add TAC/LOA records with fully filled out LOA fields
	loa := factory.BuildFullLineOfAccounting(nil, []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaInstlAcntgActID: models.StringPointer("123"),
			},
		},
	}, nil)
	tac := factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TAC:          *move.Orders.TAC,
				TacFnBlModCd: models.StringPointer("W"),
				LoaSysID:     loa.LoaSysID,
			},
		},
		{
			Model: loa,
		},
	}, nil)

	factory.BuildMove(suite.DB(), nil, nil)

	suite.Run("successfully return only navy moves with an approved payment request", func() {
		reports, err := reportListFetcher.BuildReportFromMoves(appCtx, &defaultSearchParams)
		suite.NoError(err)

		suite.Equal(1, len(reports))
		suite.Equal(tac.TAC, *reports[0].TAC)
	})
}
