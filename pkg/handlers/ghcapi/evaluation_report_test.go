package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	evaluationreportservice "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetEvaluationReportsHandler() {
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
		handler := GetEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: fetcher,
		}
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/evaluation-reports/", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveEvaluationReportsParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveEvaluationReportsOK{}, response)
		suite.NoError(response.(*moveop.GetMoveEvaluationReportsOK).Payload.Validate(strfmt.Default))
		suite.Len(response.(*moveop.GetMoveEvaluationReportsOK).Payload, 1)
	})
	suite.Run("Request error", func() {
		officeUser, move, handlerConfig := setupTestData()
		mockFetcher := mocks.EvaluationReportListFetcher{}
		handler := GetEvaluationReportsHandler{
			HandlerConfig:               handlerConfig,
			EvaluationReportListFetcher: &mockFetcher,
		}
		mockFetcher.On("FetchEvaluationReports",
			mock.AnythingOfType("*appcontext.appContext"),
			move.ID,
			officeUser.ID,
		).Return(nil, apperror.QueryError{})

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/evaluation-reports/", move.ID), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := moveop.GetMoveEvaluationReportsParams{
			HTTPRequest: request,
			MoveID:      *handlers.FmtUUID(move.ID),
		}
		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveEvaluationReportsInternalServerError{}, response)
	})
}
