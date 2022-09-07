package ghcapi

import (
	"net/http/httptest"

	pwsviolationsop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	violationservice "github.com/transcom/mymove/pkg/services/pws_violation"
)

func (suite *HandlerSuite) TestGetPWSViolationsHandler() {

	suite.Run("Successful fetch", func() {
		fetcher := violationservice.NewPWSViolationsFetcher()

		request := httptest.NewRequest("GET", "/pws-violations", nil)
		params := pwsviolationsop.GetPWSViolationsParams{
			HTTPRequest: request,
		}
		handlerConfig := suite.HandlerConfig()
		handler := GetPWSViolationsHandler{
			HandlerConfig:        handlerConfig,
			PWSViolationsFetcher: fetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&pwsviolationsop.GetPWSViolationsOK{}, response)
	})

}
