package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
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

func (suite *PayloadsSuite) TestReport() {
	appCtx := suite.AppContextForTest()
	now := time.Now()

	// Setup data
	report := models.PPTASReport{
		FirstName:     models.StringPointer("John"),
		LastName:      models.StringPointer("Doe"),
		MiddleInitial: models.StringPointer("A"),
		Edipi:         models.StringPointer("1234567890"),
		PhonePrimary:  models.StringPointer("555-555-5555"),
		EmailPrimary:  models.StringPointer("john.doe@example.com"),
		Address: &models.Address{
			StreetAddress1: "123 Main St",
			City:           "Some City",
			State:          "NY",
			PostalCode:     "10001",
			County:         "Some County",
		},
		OrdersDate:   &now,
		TAC:          models.StringPointer("CACI"),
		OrdersType:   internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber: models.StringPointer("123456"),
	}

	suite.Run("valid report", func() {
		payload := PPTASReport(appCtx, &report)

		suite.NotNil(payload)
		suite.Equal(*report.FirstName, payload.FirstName)
		suite.Equal(*report.LastName, payload.LastName)
		suite.Equal("A", *payload.MiddleInitial)
		suite.Equal(*report.Edipi, payload.Edipi)
		suite.Equal(*report.PhonePrimary, payload.PhonePrimary)
		suite.Equal(*report.EmailPrimary, payload.EmailPrimary)
		suite.Equal(*report.OrdersNumber, payload.OrdersNumber)
		suite.Equal(int64(report.ShipmentNum), payload.ShipmentNum)
	})

	suite.Run("nil report", func() {
		payload := PPTASReport(appCtx, nil)

		suite.Nil(payload)
	})
}
