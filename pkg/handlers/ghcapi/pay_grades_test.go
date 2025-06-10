package ghcapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	ordersop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/orders"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestGetPayGradesHandler() {
	suite.Run("successful returns pay grades", func() {
		affiliation := models.AffiliationAIRFORCE.String()

		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/paygrade/"+affiliation, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := ordersop.GetPayGradesParams{
			HTTPRequest: req,
			Affiliation: affiliation,
		}

		handler := GetPayGradesHandler{
			HandlerConfig: suite.NewHandlerConfig()}

		response := handler.Handle(params)
		suite.Assertions.IsType(&ordersop.GetPayGradesOK{}, response)
		responsePayload := response.(*ordersop.GetPayGradesOK)
		suite.Equal(26, len(responsePayload.Payload))
	})
}
