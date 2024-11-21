package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	evaluationReportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	evaluationreportservice "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetShipmentEvaluationReportsHandler() {
	setupTestData := func() (models.OfficeUser, models.Move, handlers.HandlerConfig) {
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		handlerConfig := suite.createS3HandlerConfig()
		return officeUser, move, handlerConfig
	}

	suite.Run("Successful list fetch", func() {
		officeUser, move, handlerConfig := setupTestData()
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type: models.EvaluationReportTypeShipment,
				},
			},
		}, nil)

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
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		handlerConfig := suite.HandlerConfig()
		return officeUser, move, handlerConfig
	}

	suite.Run("Successful list fetch", func() {
		officeUser, move, handlerConfig := setupTestData()
		factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
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
		move := factory.BuildMove(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		fetcher := evaluationreportservice.NewEvaluationReportFetcher()

		evaluationReport := factory.BuildEvaluationReport(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

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
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
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
		move := factory.BuildMove(suite.DB(), nil, nil)

		handler := CreateEvaluationReportHandler{
			HandlerConfig:           handlerConfig,
			EvaluationReportCreator: creator,
		}

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

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
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		reportID := report.ID

		updater := &mocks.EvaluationReportUpdater{}
		handlerConfig := suite.HandlerConfig()
		handler := SaveEvaluationReportHandler{handlerConfig, updater}
		requestUser := factory.BuildUser(nil, nil, nil) //Build stubbed user

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		now := strfmt.Date(time.Now())
		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				InspectionDate:                     &now,
				InspectionType:                     ghcmessages.EvaluationReportInspectionTypePHYSICAL.Pointer(),
				Location:                           ghcmessages.EvaluationReportLocationOTHER.Pointer(),
				LocationDescription:                models.StringPointer("location description"),
				ObservedShipmentDeliveryDate:       handlers.FmtDate(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)),
				ObservedShipmentPhysicalPickupDate: handlers.FmtDate(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)), Remarks: models.StringPointer("new remarks"),
				SeriousIncident:     handlers.FmtBool(true),
				SeriousIncidentDesc: models.StringPointer("serious incident description"),
				ViolationsObserved:  handlers.FmtBool(false),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		requestUser := factory.BuildUser(nil, nil, nil)

		request := httptest.NewRequest("PUT", fmt.Sprintf("/evaluation-reports/%s", reportID), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)

		params := evaluationReportop.SaveEvaluationReportParams{
			HTTPRequest: request,
			Body: &ghcmessages.EvaluationReport{
				Remarks: models.StringPointer("new remarks"),
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
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
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
		requestUser := factory.BuildUser(nil, nil, nil)

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
		requestUser := factory.BuildUser(nil, nil, nil)

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
		requestUser := factory.BuildUser(nil, nil, nil)

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
		requestUser := factory.BuildUser(nil, nil, nil)

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

func (suite *HandlerSuite) TestAddAppealToViolationHandler() {
	suite.Run("Successful Add Appeal to Violation", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.ReportViolationsAddAppeal{}
		handler := AddAppealToViolationHandler{
			HandlerConfig:             handlerConfig,
			ReportViolationsAddAppeal: mockService,
		}
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})
		body := ghcmessages.CreateAppeal{AppealStatus: "sustained", Remarks: "Appeal submitted"}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/%s/appeal/add", report.ID, violation.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := evaluationReportop.AddAppealToViolationParams{
			HTTPRequest:       request,
			ReportID:          strfmt.UUID(report.ID.String()),
			ReportViolationID: strfmt.UUID(violation.ID.String()),
			Body:              &body,
		}

		mockService.On("AddAppealToViolation", mock.Anything, report.ID, violation.ID, officeUser.ID, body.Remarks, body.AppealStatus).Return(models.GsrAppeal{}, nil).Once()

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToViolationNoContent{}, response)
	})

	suite.Run("Unsuccessful Add Appeal to Violation - Missing Data", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.ReportViolationsAddAppeal{}
		handler := AddAppealToViolationHandler{
			HandlerConfig:             handlerConfig,
			ReportViolationsAddAppeal: mockService,
		}
		reportID := uuid.Must(uuid.NewV4())
		violationID := uuid.Must(uuid.NewV4())
		body := ghcmessages.CreateAppeal{AppealStatus: "", Remarks: ""}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/%s/appeal/add", reportID, violationID), nil)
		params := evaluationReportop.AddAppealToViolationParams{
			HTTPRequest:       request,
			ReportID:          strfmt.UUID(reportID.String()),
			ReportViolationID: strfmt.UUID(violationID.String()),
			Body:              &body,
		}

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToViolationForbidden{}, response)
	})

	suite.Run("Unsuccessful Add Appeal to Violation - Service Error", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.ReportViolationsAddAppeal{}
		handler := AddAppealToViolationHandler{
			HandlerConfig:             handlerConfig,
			ReportViolationsAddAppeal: mockService,
		}
		reportID := uuid.Must(uuid.NewV4())
		violationID := uuid.Must(uuid.NewV4())
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		body := ghcmessages.CreateAppeal{AppealStatus: "sustained", Remarks: "Appeal submitted"}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/%s/appeal/add", reportID, violationID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := evaluationReportop.AddAppealToViolationParams{
			HTTPRequest:       request,
			ReportID:          strfmt.UUID(reportID.String()),
			ReportViolationID: strfmt.UUID(violationID.String()),
			Body:              &body,
		}

		mockService.On("AddAppealToViolation", mock.Anything, reportID, violationID, mock.Anything, body.Remarks, body.AppealStatus).Return(models.GsrAppeal{}, apperror.NewForbiddenError("Forbidden")).Once()

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToViolationForbidden{}, response)
	})
}

func (suite *HandlerSuite) TestAddAppealToSeriousIncidentHandler() {
	suite.Run("Successful Add Appeal to Serious Incident", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.SeriousIncidentAddAppeal{}
		handler := AddAppealToSeriousIncidentHandler{
			HandlerConfig:            handlerConfig,
			SeriousIncidentAddAppeal: mockService,
		}
		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeGSR})
		body := ghcmessages.CreateAppeal{AppealStatus: "sustained", Remarks: "Appeal submitted for serious incident"}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/appeal/add", report.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := evaluationReportop.AddAppealToSeriousIncidentParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(report.ID.String()),
			Body:        &body,
		}

		mockService.On("AddAppealToSeriousIncident", mock.Anything, report.ID, officeUser.ID, body.Remarks, body.AppealStatus).Return(models.GsrAppeal{}, nil).Once()

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToSeriousIncidentNoContent{}, response)
	})

	suite.Run("Unsuccessful Add Appeal to Serious Incident - Missing Data", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.SeriousIncidentAddAppeal{}
		handler := AddAppealToSeriousIncidentHandler{
			HandlerConfig:            handlerConfig,
			SeriousIncidentAddAppeal: mockService,
		}
		reportID := uuid.Must(uuid.NewV4())
		body := ghcmessages.CreateAppeal{AppealStatus: "", Remarks: ""}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/appeal/add", reportID), nil)
		params := evaluationReportop.AddAppealToSeriousIncidentParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(reportID.String()),
			Body:        &body,
		}

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToSeriousIncidentForbidden{}, response)
	})

	suite.Run("Unsuccessful Add Appeal to Serious Incident - Service Error", func() {
		handlerConfig := suite.HandlerConfig()
		mockService := &mocks.SeriousIncidentAddAppeal{}
		handler := AddAppealToSeriousIncidentHandler{
			HandlerConfig:            handlerConfig,
			SeriousIncidentAddAppeal: mockService,
		}
		reportID := uuid.Must(uuid.NewV4())
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})
		body := ghcmessages.CreateAppeal{AppealStatus: "pending", Remarks: "Appeal submitted"}
		request := httptest.NewRequest("POST", fmt.Sprintf("/evaluation-reports/%s/appeal/add", reportID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := evaluationReportop.AddAppealToSeriousIncidentParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(reportID.String()),
			Body:        &body,
		}

		mockService.On("AddAppealToSeriousIncident", mock.Anything, reportID, mock.Anything, body.Remarks, body.AppealStatus).Return(models.GsrAppeal{}, apperror.NewForbiddenError("Forbidden")).Once()

		response := handler.Handle(params)
		suite.IsType(&evaluationReportop.AddAppealToSeriousIncidentForbidden{}, response)
	})
}
