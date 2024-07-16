package report

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mocks "github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ReportServiceSuite) TestReportFetcher() {
	ppmEstimator := mocks.PPMEstimator{}

	reportListFetcher := NewReportListFetcher(&ppmEstimator)
	defaultSearchParams := services.MoveTaskOrderFetcherParams{}

	appCtx := suite.AppContextForTest()
	now := time.Now()

	// Setup data
	serviceMember := models.ServiceMember{
		FirstName:      models.StringPointer("John"),
		LastName:       models.StringPointer("Doe"),
		MiddleName:     models.StringPointer("A"),
		Edipi:          models.StringPointer("1234567890"),
		Telephone:      models.StringPointer("555-555-5555"),
		PersonalEmail:  models.StringPointer("john.doe@example.com"),
		BackupContacts: []models.BackupContact{{Email: "backup@example.com"}},
		ResidentialAddress: &models.Address{
			StreetAddress1: "123 Main St",
			City:           "Some City",
			State:          "NY",
			PostalCode:     "10001",
			County:         "Some County",
		},
	}

	hasDependents := true
	orders := models.Order{
		ServiceMember:           serviceMember,
		IssueDate:               now,
		TAC:                     models.StringPointer("CACI"),
		OrdersType:              internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:            models.StringPointer("123456"),
		HasDependents:           hasDependents,
		Entitlement:             &models.Entitlement{},
		NewDutyLocation:         models.DutyLocation{},
		OriginDutyLocationGBLOC: models.StringPointer("XYZ"),
	}

	shipmentId, _ := uuid.NewV4()
	shipment2Id, _ := uuid.NewV4()
	primeWeight := unit.Pound(5000)
	estimatedWeight := unit.Pound(4500)
	distance := unit.Miles(300)
	move := models.Move{
		Orders: orders,
		MTOShipments: models.MTOShipments{
			{
				ID:                   shipmentId,
				PrimeActualWeight:    &primeWeight,
				PrimeEstimatedWeight: &estimatedWeight,
				Distance:             &distance,
			},
			{
				ID: shipment2Id,
			},
		},
		PaymentRequests: []models.PaymentRequest{
			{
				Status:     models.PaymentRequestStatusReviewed,
				ReviewedAt: &now,
			},
		},
		ServiceCounselingCompletedAt: &now,
	}

	ordersIssueDate := time.Now()
	endDate := ordersIssueDate.AddDate(1, 0, 0)
	dptId := "1"

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

	suite.Run("successfully return only navy moves with an approved payment request", func() {
		nonNavyMove := factory.BuildMove(suite.DB(), nil, nil)
		reports, err := reportListFetcher.BuildReportFromMoves(appCtx, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Len(reports, 1)
		suite.NotEqual(reports[0].ShipmentId, nonNavyMove.ID)
	})
}
