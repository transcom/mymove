package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	evaluationReportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	evaluationreportservice "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetShipmentEvaluationReportsHandler() {
	setupTestData := func() (models.OfficeUser, models.Move, handlers.HandlerConfig) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		handlerConfig := suite.createS3HandlerConfig()
		return officeUser, move, handlerConfig
	}

	suite.Run("Successful list fetch", func() {
		officeUser, move, handlerConfig := setupTestData()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			OfficeUser:  officeUser,
			Move:        move,
			MTOShipment: shipment,
		})

		fetcher := evaluationreportservice.NewEvaluationReportFetcher()
		handler := GetShipmentEvaluationReportsHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: fetcher,
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/shipment-evaluation-reports-list", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveShipmentEvaluationReportsListParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveShipmentEvaluationReportsListOK{}, response)
		payload := response.(*moveop.GetMoveShipmentEvaluationReportsListOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload, 1)
	})

	suite.Run("Request error", func() {
		officeUser, move, handlerConfig := setupTestData()
		mockFetcher := mocks.EvaluationReportFetcher{}
		handler := GetShipmentEvaluationReportsHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: &mockFetcher,
		}
		mockFetcher.On("FetchEvaluationReports",
			mock.AnythingOfType("*appcontext.appContext"),
			models.EvaluationReportTypeShipment,
			move.ID,
			officeUser.ID,
		).Return(nil, apperror.QueryError{})

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/shipment-evaluation-reports-list", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveShipmentEvaluationReportsListParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveShipmentEvaluationReportsListInternalServerError{}, response)
		payload := response.(*moveop.GetMoveShipmentEvaluationReportsListInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestGetCounselingEvaluationReportsHandler() {
	setupTestData := func() (models.OfficeUser, models.Move, handlers.HandlerConfig) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		handlerConfig := suite.HandlerConfig()
		return officeUser, move, handlerConfig
	}

	suite.Run("Successful list fetch", func() {
		officeUser, move, handlerConfig := setupTestData()
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			OfficeUser: officeUser,
			Move:       move,
		})

		fetcher := evaluationreportservice.NewEvaluationReportFetcher()
		handler := GetCounselingEvaluationReportsHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: fetcher,
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/counseling-evaluation-reports-list", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveCounselingEvaluationReportsListParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveCounselingEvaluationReportsListOK{}, response)
		payload := response.(*moveop.GetMoveCounselingEvaluationReportsListOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload, 1)
	})

	suite.Run("Request error", func() {
		officeUser, move, handlerConfig := setupTestData()
		mockFetcher := mocks.EvaluationReportFetcher{}
		handler := GetCounselingEvaluationReportsHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: &mockFetcher,
		}
		mockFetcher.On("FetchEvaluationReports",
			mock.AnythingOfType("*appcontext.appContext"),
			models.EvaluationReportTypeCounseling,
			move.ID,
			officeUser.ID,
		).Return(nil, apperror.QueryError{})

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/counseling-evaluation-reports-list", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveCounselingEvaluationReportsListParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveCounselingEvaluationReportsListInternalServerError{}, response)
		payload := response.(*moveop.GetMoveCounselingEvaluationReportsListInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestGetEvaluationReportByIDHandler() {
	// 200 response
	suite.Run("Successful fetch (integration) test", func() {
		handlerConfig := suite.HandlerConfig()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := evaluationreportservice.NewEvaluationReportFetcher()

		evaluationReport := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				OfficeUserID: officeUser.ID,
				MoveID:       move.ID,
			}})

		handler := GetEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: fetcher,
		}

		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s",
			evaluationReport.ID.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := evaluationReportop.GetEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(evaluationReport.ID.String()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.GetEvaluationReportOK{}, response)
		payload := response.(*evaluationReportop.GetEvaluationReportOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// 404 response
	suite.Run("404 response when service returns not found", func() {
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuidForReport, _ := uuid.NewV4()
		handlerConfig := suite.HandlerConfig()
		mockFetcher := mocks.EvaluationReportFetcher{}
		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s", uuidForReport.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := evaluationReportop.GetEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(uuidForReport.String()),
		}
		mockFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			uuidForReport,
			officeUser.ID,
		).Return(nil, apperror.NewNotFoundError(uuidForReport, "while looking for evaluation report"))

		handler := GetEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: &mockFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.GetEvaluationReportNotFound{}, response)
		payload := response.(*evaluationReportop.GetEvaluationReportNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// 403 response
	suite.Run("403 response when service returns forbidden", func() {
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuidForReport, _ := uuid.NewV4()
		handlerConfig := suite.HandlerConfig()
		mockFetcher := mocks.EvaluationReportFetcher{}
		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s", uuidForReport.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := evaluationReportop.GetEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(uuidForReport.String()),
		}
		mockFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			uuidForReport,
			officeUser.ID,
		).Return(nil, apperror.NewForbiddenError("Draft evaluation reports are viewable only by their owner/creator."))

		handler := GetEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: &mockFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.GetEvaluationReportForbidden{}, response)
		payload := response.(*evaluationReportop.GetEvaluationReportForbidden).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestCreateEvaluationReportHandler() {

	var officeUser models.OfficeUser

	suite.Run("Successful POST", func() {

		handlerConfig := suite.HandlerConfig()

		creator := &mocks.EvaluationReportCreator{}
		move := testdatagen.MakeDefaultMove(suite.DB())

		handler := CreateEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportCreator: creator,
		}

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: move.ID,
			},
		})

		body := ghcmessages.CreateEvaluationReport{ShipmentID: strfmt.UUID(shipment.ID.String())}
		request := httptest.NewRequest("POST", "/moves/shipment-evaluation-reports/", nil)

		params := evaluationReportop.CreateEvaluationReportParams{
			HTTPRequest: request,
			Body:        &body,
		}

		returnReport := models.EvaluationReport{
			ID:           uuid.Must(uuid.NewV4()),
			MoveID:       move.ID,
			Move:         move,
			ShipmentID:   &shipment.ID,
			Shipment:     &shipment,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Type:         models.EvaluationReportTypeShipment,
			OfficeUser:   officeUser,
			OfficeUserID: officeUser.ID,
		}

		creator.On("CreateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("string"),
		).Return(&returnReport, nil).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.CreateEvaluationReportOK{}, response)
		payload := response.(*evaluationReportop.CreateEvaluationReportOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful POST", func() {

		handlerConfig := suite.HandlerConfig()

		creator := &mocks.EvaluationReportCreator{}
		handler := CreateEvaluationReportHandler{handlerConfig, creator}

		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: move.ID,
			},
		})
		body := ghcmessages.CreateEvaluationReport{ShipmentID: strfmt.UUID(shipment.ID.String())}
		request := httptest.NewRequest("POST", "/moves/shipment-evaluation-reports/", nil)

		params := evaluationReportop.CreateEvaluationReportParams{
			HTTPRequest: request,
			Body:        &body,
		}

		creator.On("CreateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("string"),
		).Return(nil, fmt.Errorf("error")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.CreateEvaluationReportInternalServerError{}, response)
		payload := response.(*evaluationReportop.CreateEvaluationReportInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestDeleteEvaluationReportHandler() {

	suite.Run("Successful DELETE", func() {
		reportID := uuid.Must(uuid.NewV4())

		deleter := &mocks.EvaluationReportDeleter{}
		handlerConfig := suite.HandlerConfig()
		handler := DeleteEvaluationReportHandler{handlerConfig, deleter}

		request := httptest.NewRequest("DELETE", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)

		params := evaluationReportop.DeleteEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    *handlers.FmtUUID(reportID),
		}

		deleter.On("DeleteEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DeleteEvaluationReportNoContent{}, response)

		// Validate outgoing payload: no payload
	})
}

func (suite *HandlerSuite) TestSubmitEvaluationReportHandler() {
	suite.Run("Successful save", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(nil).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportNoContent{}, response)

		// Validate outgoing payload: no payload
	})

	suite.Run("Precondition failed", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewPreconditionFailedError(reportID, nil)).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportPreconditionFailed{}, response)
		payload := response.(*evaluationReportop.SubmitEvaluationReportPreconditionFailed).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Not found error", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewNotFoundError(reportID, "message")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportNotFound{}, response)
		payload := response.(*evaluationReportop.SubmitEvaluationReportNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Invalid input", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewInvalidInputError(reportID, nil, nil, "message")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportUnprocessableEntity{}, response)
		payload := response.(*evaluationReportop.SubmitEvaluationReportUnprocessableEntity).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Forbidden error", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewForbiddenError("message")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportForbidden{}, response)
		payload := response.(*evaluationReportop.SubmitEvaluationReportForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Internal server error", func() {
		updater := &mocks.EvaluationReportUpdater{}

		reportID := uuid.Must(uuid.NewV4())
		handlerConfig := suite.HandlerConfig()
		handler := SubmitEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedOfficeUser(suite.DB())

		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/submit", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := evaluationReportop.SubmitEvaluationReportParams{
			HTTPRequest: request,
			IfMatch:     "",
			ReportID:    *handlers.FmtUUID(reportID),
		}

		updater.On("SubmitEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewInternalServerError("message")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Assertions.IsType(&evaluationReportop.SubmitEvaluationReportInternalServerError{}, response)
		payload := response.(*evaluationReportop.SubmitEvaluationReportInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestSaveEvaluationReportHandler() {

	suite.Run("Successful save", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		now := strfmt.Date(time.Now())
		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				EvaluationLengthMinutes:       handlers.FmtInt64(45),
				InspectionDate:                &now,
				InspectionType:                ghcmessages.EvaluationReportInspectionTypePHYSICAL.Pointer().Pointer(),
				Location:                      ghcmessages.EvaluationReportLocationORIGIN.Pointer(),
				LocationDescription:           swag.String("location description"),
				ObservedClaimsResponseDate:    handlers.FmtDate(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				ObservedDate:                  handlers.FmtDate(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
				ObservedPickupDate:            handlers.FmtDate(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)),
				ObservedPickupSpreadStartDate: handlers.FmtDate(time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC)),
				ObservedPickupSpreadEndDate:   handlers.FmtDate(time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC)),
				ObservedDeliveryDate:          handlers.FmtDate(time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC)),
				Remarks:                       swag.String("new remarks"),
				SeriousIncident:               handlers.FmtBool(true),
				SeriousIncidentDesc:           swag.String("serious incident description"),
				TravelTimeMinutes:             handlers.FmtInt64(30),
				ViolationsObserved:            handlers.FmtBool(true),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(nil).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportNoContent{}, response)

		// Validate outgoing payload: no payload
	})

	suite.Run("Not found error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewNotFoundError(reportID, "message")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportNotFound{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Invalid input error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewInvalidInputError(reportID, nil, nil, "message")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportUnprocessableEntity{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportUnprocessableEntity).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Precondition failed error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewPreconditionFailedError(reportID, nil)).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportPreconditionFailed{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportPreconditionFailed).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Forbidden error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewForbiddenError("")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportForbidden{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Conflict error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(apperror.NewConflictError(reportID, "")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportConflict{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportConflict).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unknown error", func() {
		reportID := uuid.Must(uuid.NewV4())

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: swag.String("new remarks"),
			},
			ReportID: *handlers.FmtUUID(reportID),
		}

		updater.On("UpdateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("string"),
		).Return(fmt.Errorf("this is some sort of error")).Once()

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.SaveEvaluationReportInternalServerError{}, response)
		payload := response.(*evaluationReportop.SaveEvaluationReportInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestDownloadEvaluationReportHandler() {

	suite.Run("Successful download", func() {
		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		reportID := report.ID

		reportFetcher := &mocks.EvaluationReportFetcher{}
		shipmentFetcher := &mocks.MTOShipmentFetcher{}
		orderFetcher := &mocks.OrderFetcher{}
		violationsFetcher := &mocks.ReportViolationFetcher{}
		handlerConfig := suite.HandlerConfig()
		handler := DownloadEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: reportFetcher,
			MTOShipmentFetcher:      shipmentFetcher,
			OrderFetcher:            orderFetcher,
			ReportViolationFetcher:  violationsFetcher,
		}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s/download", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := evaluationReportop.DownloadEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    *handlers.FmtUUID(reportID),
		}

		reportFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(&report, nil).Once()

		orderFetcher.On("FetchOrder",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(&report.Move.Orders, nil)

		violationsFetcher.On("FetchReportViolationsByReportID",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(models.ReportViolations{}, nil)

		shipmentFetcher.On("ListMTOShipments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return([]models.MTOShipment{}, nil)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DownloadEvaluationReportOK{}, response)
		payload := response.(*evaluationReportop.DownloadEvaluationReportOK).Payload

		// Validate outgoing payload: payload should be an instance of a ReadCloser interface
		suite.NotNil(payload)
	})

	suite.Run("Not found error", func() {
		reportID := uuid.Must(uuid.NewV4())

		reportFetcher := &mocks.EvaluationReportFetcher{}
		handlerConfig := suite.HandlerConfig()
		handler := DownloadEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: reportFetcher,
		}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s/download", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := evaluationReportop.DownloadEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    *handlers.FmtUUID(reportID),
		}

		reportFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, apperror.NewNotFoundError(uuid.Nil, "not found")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DownloadEvaluationReportNotFound{}, response)
		payload := response.(*evaluationReportop.DownloadEvaluationReportNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Query error should result in 500", func() {
		reportID := uuid.Must(uuid.NewV4())

		reportFetcher := &mocks.EvaluationReportFetcher{}
		handlerConfig := suite.HandlerConfig()
		handler := DownloadEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: reportFetcher,
		}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s/download", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := evaluationReportop.DownloadEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    *handlers.FmtUUID(reportID),
		}

		reportFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, apperror.NewQueryError("", nil, "")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DownloadEvaluationReportInternalServerError{}, response)
		payload := response.(*evaluationReportop.DownloadEvaluationReportInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unknown error should result in 500", func() {
		reportID := uuid.Must(uuid.NewV4())

		reportFetcher := &mocks.EvaluationReportFetcher{}
		handlerConfig := suite.HandlerConfig()
		handler := DownloadEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportFetcher: reportFetcher,
		}
		requestUser := testdatagen.MakeStubbedUser(suite.DB())

		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s/download", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := evaluationReportop.DownloadEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    *handlers.FmtUUID(reportID),
		}

		reportFetcher.On("FetchEvaluationReportByID",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, fmt.Errorf("an error")).Once()

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DownloadEvaluationReportInternalServerError{}, response)
		payload := response.(*evaluationReportop.DownloadEvaluationReportInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
