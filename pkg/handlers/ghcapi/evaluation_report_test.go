package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	evaluationReportop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
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

		fetcher := evaluationreportservice.NewEvaluationReportListFetcher()
		handler := GetShipmentEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: fetcher,
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
		mockFetcher := mocks.EvaluationReportListFetcher{}
		handler := GetShipmentEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: &mockFetcher,
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

		fetcher := evaluationreportservice.NewEvaluationReportListFetcher()
		handler := GetCounselingEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: fetcher,
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
		mockFetcher := mocks.EvaluationReportListFetcher{}
		handler := GetCounselingEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: &mockFetcher,
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

func (suite *HandlerSuite) TestCreateEvaluationReportHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	suite.Run("Successful POST", func() {

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		creator := &mocks.EvaluationReportCreator{}
		handler := CreateEvaluationReportHandler{handlerConfig, creator}

		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: move.ID,
			},
		})
		body := ghcmessages.CreateShipmentEvaluationReport{ShipmentID: handlers.FmtUUID(shipment.ID)}
		request := httptest.NewRequest("POST", "/moves/shipment-evaluation-reports/", nil)

		params := evaluationReportop.CreateEvaluationReportForShipmentParams{
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
		).Return(&returnReport, nil).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.CreateEvaluationReportForShipmentOK{}, response)
	})

	suite.Run("Unsuccessful POST", func() {

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		creator := &mocks.EvaluationReportCreator{}
		handler := CreateEvaluationReportHandler{handlerConfig, creator}

		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: move.ID,
			},
		})
		body := ghcmessages.CreateShipmentEvaluationReport{ShipmentID: handlers.FmtUUID(shipment.ID)}
		request := httptest.NewRequest("POST", "/moves/shipment-evaluation-reports/", nil)

		params := evaluationReportop.CreateEvaluationReportForShipmentParams{
			HTTPRequest: request,
			Body:        &body,
		}

		creator.On("CreateEvaluationReport",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.EvaluationReport"),
		).Return(nil, fmt.Errorf("error")).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.CreateEvaluationReportForShipmentInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteEvaluationReportHandler() {

	suite.Run("Successful DELETE", func() {
		reportID := uuid.Must(uuid.NewV4())

		deleter := &mocks.EvaluationReportDeleter{}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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

		response := handler.Handle(params)

		suite.Assertions.IsType(&evaluationReportop.DeleteEvaluationReportNoContent{}, response)
	})
}
