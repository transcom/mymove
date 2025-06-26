package report

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/entitlements"
	mocks "github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ReportServiceSuite) TestReportFetcher() {
	ppmEstimator := mocks.PPMEstimator{}
	moveFetcher := mocks.MoveFetcher{}
	tacFetcher := mocks.TransportationAccountingCodeFetcher{}
	loaFetcher := mocks.LineOfAccountingFetcher{}
	waf := entitlements.NewWeightAllotmentFetcher()

	reportListFetcher := NewPPTASReportListFetcher(&ppmEstimator, &moveFetcher, &tacFetcher, &loaFetcher, waf)

	appCtx := suite.AppContextForTest()

	// Setup data
	testDate := time.Now()
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
				TAC:       models.StringPointer("E12A"),
				IssueDate: testDate,
			},
		},
	}, nil)
	reweighedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	reweighWeight := unit.Pound(2399)
	reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
		Reweigh: models.Reweigh{
			Weight: &reweighWeight,
		},
		MTOShipment: reweighedShipment,
	})

	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ServiceCounselingCompletedAt: &testDate,
			},
		},
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusApproved,
				PrimeEstimatedWeight: models.PoundPointer(1000),
			},
		},
	}, nil)

	move.MTOShipments[0].Reweigh = &reweigh

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

	beginDate := time.Now().AddDate(0, 0, -10)
	endDate := time.Now().AddDate(0, 0, 10)
	hsgdscd := models.LineOfAccountingHouseholdGoodsCodeEnlisted
	loa := factory.BuildFullLineOfAccounting(nil, []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaSysID:           models.StringPointer("ooga booga"),
				LoaInstlAcntgActID: models.StringPointer("123"),
				LoaBgnDt:           &beginDate,
				LoaEndDt:           &endDate,
				LoaHsGdsCd:         &hsgdscd,
			},
		},
	}, nil)
	tac := factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TAC:               "E12A",
				TacFnBlModCd:      models.StringPointer("W"),
				LoaSysID:          loa.LoaSysID,
				TrnsprtnAcntBgnDt: &beginDate,
				TrnsprtnAcntEndDt: &endDate,
			},
		},
		{
			Model:    loa,
			LinkOnly: false,
		},
	}, nil)

	var movesForReport models.Moves
	movesForReport = append(movesForReport, move)

	suite.Run("successfully create a report", func() {
		tacFetcher.On("FetchOrderTransportationAccountingCodes",
			mock.Anything,
			mock.Anything,
			"E12A",
			mock.AnythingOfType("*appcontext.appContext"),
		).Return(nil, nil)

		reports, err := reportListFetcher.BuildPPTASReportsFromMoves(appCtx, movesForReport)
		suite.NoError(err)

		suite.Equal(1, len(reports))
		suite.Equal(tac.TAC, *reports[0].TAC)

		// 110% of prime estimated weight
		maxBillableWeight := move.MTOShipments[0].PrimeEstimatedWeight.Float64() * 1.1
		suite.Equal(int(maxBillableWeight), reports[0].MaxBillableWeight.Int())
	})
}
