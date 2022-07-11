package ghcapi

import (
	"fmt"
	"net/http/httptest"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	evaluationreportservice "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetEvaluationReportsHandler() {
	suite.Run("Successful list fetch", func() {
		fetcher := evaluationreportservice.NewEvaluationReportListFetcher()
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		// TODO add some reports or something in here and validate response format

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
	})
}
