package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	reportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	evaluationreportservice "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetShipmentEvaluationReportsHandler() {
	setupTestData := func() (models.OfficeUser, models.Move, handlers.HandlerConfig) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveShipmentEvaluationReportsListOK{}, response)
		suite.NoError(response.(*moveop.GetMoveShipmentEvaluationReportsListOK).Payload.Validate(strfmt.Default))
		suite.Len(response.(*moveop.GetMoveShipmentEvaluationReportsListOK).Payload, 1)
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
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveShipmentEvaluationReportsListInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestGetCounselingEvaluationReportsHandler() {
	setupTestData := func() (models.OfficeUser, models.Move, handlers.HandlerConfig) {
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveCounselingEvaluationReportsListOK{}, response)
		suite.NoError(response.(*moveop.GetMoveCounselingEvaluationReportsListOK).Payload.Validate(strfmt.Default))
		suite.Len(response.(*moveop.GetMoveCounselingEvaluationReportsListOK).Payload, 1)
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
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveCounselingEvaluationReportsListInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestGetEvaluationReportByIDHandler() {
	// 200 response
	suite.Run("Successful fetch (integration) test", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
		params := reportop.GetEvaluationReportParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(evaluationReport.ID.String()),
		}
		response := handler.Handle(params)
		suite.IsType(&reportop.GetEvaluationReportOK{}, response)
	})

	// 404 response
	suite.Run("404 response when service returns not found", func() {
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuidForReport, _ := uuid.NewV4()
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		mockFetcher := mocks.EvaluationReportFetcher{}
		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s", uuidForReport.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := reportop.GetEvaluationReportParams{
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

		response := handler.Handle(params)
		suite.IsType(&reportop.GetEvaluationReportNotFound{}, response)
	})

	// 403 response
	suite.Run("403 response when service returns forbidden", func() {
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuidForReport, _ := uuid.NewV4()
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		mockFetcher := mocks.EvaluationReportFetcher{}
		request := httptest.NewRequest("GET", fmt.Sprintf("/evaluation-reports/%s", uuidForReport.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := reportop.GetEvaluationReportParams{
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

		response := handler.Handle(params)
		suite.IsType(&reportop.GetEvaluationReportForbidden{}, response)
	})
}
