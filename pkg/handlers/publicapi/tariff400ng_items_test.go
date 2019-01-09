package publicapi

import (
	"net/http/httptest"

	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) getTariff400ngItemsParams(tspUser models.TspUser, requiresPreApproval bool) accessorialop.GetTariff400ngItemsParams {
	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/tariff400ng_items", nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	return accessorialop.GetTariff400ngItemsParams{
		HTTPRequest:         req,
		RequiresPreApproval: handlers.FmtBool(requiresPreApproval),
	}
}

func (suite *HandlerSuite) TestGetTariff400ngItemsHandler() {
	// Does not require pre-approval
	testdatagen.MakeDefaultTariff400ngItem(suite.DB())
	// Does require pre-approval
	item2 := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                "9000",
			RequiresPreApproval: true,
		},
	})

	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	// Test that only pre-approval items are returned
	requireParams := suite.getTariff400ngItemsParams(tspUser, true)
	handler := GetTariff400ngItemsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(requireParams)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.GetTariff400ngItemsOK{}, response)
	okResponse := response.(*accessorialop.GetTariff400ngItemsOK)

	// And: Payload returned is the one requiring pre-approval
	suite.Len(okResponse.Payload, 1)
	suite.Equal(okResponse.Payload[0].Code, item2.Code)

	// Test that all items are returned
	requireParams = suite.getTariff400ngItemsParams(tspUser, false)
	handler = GetTariff400ngItemsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(requireParams)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&accessorialop.GetTariff400ngItemsOK{}, response)
	okResponse = response.(*accessorialop.GetTariff400ngItemsOK)

	// And: Test that both items are returned
	suite.Len(okResponse.Payload, 2)
}
