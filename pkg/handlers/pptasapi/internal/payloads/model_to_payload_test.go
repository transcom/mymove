package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestInternalServerError() {
	traceID, _ := uuid.NewV4()
	detail := "Err"

	noDetailError := InternalServerError(nil, traceID)
	suite.Equal(handlers.InternalServerErrMessage, *noDetailError.Title)
	suite.Equal(handlers.InternalServerErrDetail, *noDetailError.Detail)
	suite.Equal(traceID.String(), noDetailError.Instance.String())

	detailError := InternalServerError(&detail, traceID)
	suite.Equal(handlers.InternalServerErrMessage, *detailError.Title)
	suite.Equal(detail, *detailError.Detail)
	suite.Equal(traceID.String(), detailError.Instance.String())
}

func (suite *PayloadsSuite) TestListReport() {
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

	suite.Run("valid move", func() {
		payload := ListReport(appCtx, &move)

		suite.NotNil(payload)
		suite.Equal(*serviceMember.FirstName, payload.FirstName)
		suite.Equal(*serviceMember.LastName, payload.LastName)
		suite.Equal("A", payload.MiddleInitial)
		suite.Equal(*serviceMember.Edipi, payload.Edipi)
		suite.Equal(*serviceMember.Telephone, payload.PhonePrimary)
		suite.Equal(*serviceMember.PersonalEmail, payload.EmailPrimary)
		suite.Equal(serviceMember.BackupContacts[0].Email, *payload.EmailSecondary)
		suite.Equal(string(orders.OrdersType), payload.OrdersType)
		suite.Equal(*orders.OrdersNumber, payload.OrdersNumber)
		suite.Equal(strfmt.DateTime(orders.IssueDate), payload.OrdersDate)
		suite.Equal(int64(len(move.MTOShipments)), payload.ShipmentNum)
		suite.Equal(move.MTOShipments[0].PrimeEstimatedWeight.Float64(), payload.WeightEstimate)
		suite.Equal(move.MTOShipments[0].PrimeActualWeight.Float64(), payload.ActualOriginNetWeight)
		longLoa := "1*1234*20242025*1234*1234*1234*1234*12345*123456*123456789012*88888888*1234*1234567890*1*123456*1*123*123456*123456*1*12*12345678*123456789012345*12*123*123456789012345678*123*123456789012"
		suite.Equal(&longLoa, payload.Loa)
	})

	suite.Run("nil move", func() {
		payload := ListReport(appCtx, nil)

		suite.Nil(payload)
	})
}
