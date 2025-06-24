package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestGetPayGradesHandler() {
	suite.Run("successful returns pay grades", func() {
		affiliation := models.AffiliationCOASTGUARD.String()

		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/paygrade/"+affiliation, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := ordersop.GetPayGradesParams{
			HTTPRequest: req,
			Affiliation: affiliation,
		}

		handler := GetPayGradesHandler{
			HandlerConfig: suite.HandlerConfig()}

		response := handler.Handle(params)
		suite.Assertions.IsType(&ordersop.GetPayGradesOK{}, response)
		responsePayload := response.(*ordersop.GetPayGradesOK)
		suite.Equal(24, len(responsePayload.Payload))
	})
}
