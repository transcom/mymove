package report

import (
	"time"

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
	now := time.Now()

	// Setup data
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:   models.StringPointer("John"),
				LastName:    models.StringPointer("Doe"),
				MiddleName:  models.StringPointer("Test"),
				Edipi:       models.StringPointer("1234567890"),
				Telephone:   models.StringPointer("555-555-5555"),
				Affiliation: (*models.ServiceMemberAffiliation)(internalmessages.AffiliationNAVY.Pointer()),
			},
		},
	}, nil)

	hasDependents := true
	orders := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model: models.Order{
				IssueDate:               now,
				TAC:                     models.StringPointer("CACI"),
				OrdersType:              internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
				OrdersNumber:            models.StringPointer("123456"),
				HasDependents:           hasDependents,
				OriginDutyLocationGBLOC: models.StringPointer("XYZ"),
			},
		},
	}, nil)

	ordersIssueDate := time.Now()
	endDate := ordersIssueDate.AddDate(1, 0, 0)
	dptId := "1"

	pr := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: &now,
			},
		},
	}, nil)

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model:    pr,
			LinkOnly: true,
		},
	}, nil)

	// Add TAC/LOA records with fully filled out LOA fields
	loa := factory.BuildFullLineOfAccounting(nil, []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaInstlAcntgActID: models.StringPointer("123"),
				LoaDptID:           &dptId,
			},
		},
	}, nil)

	factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TAC:               *move.Orders.TAC,
				TacFnBlModCd:      models.StringPointer("W"),
				TrnsprtnAcntBgnDt: &ordersIssueDate,
				TrnsprtnAcntEndDt: &endDate,
				LoaSysID:          loa.LoaSysID,
			},
		},
		{
			Model: loa,
		},
	}, nil)

	factory.BuildMove(suite.DB(), nil, nil)

	suite.Run("successfully return only navy moves with an approved payment request", func() {
		reports, err := reportListFetcher.BuildReportFromMoves(appCtx, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Len(reports, 1)
	})
}
