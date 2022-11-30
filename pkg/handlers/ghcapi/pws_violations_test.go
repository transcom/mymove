package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	pwsviolationsop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/pws_violations"
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
		suite.IsType(&pwsviolationsop.GetPWSViolationsOK{}, response)
		payload := response.(*pwsviolationsop.GetPWSViolationsOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
