package ghcapi

import (
	"net/http/httptest"

	moveorder "github.com/transcom/mymove/pkg/services/move_order"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetMoveOrderHandlerIntegration() {
	moveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	request := httptest.NewRequest("GET", "/move-orders/{moveOrderID}", nil)
	params := moveorderop.GetMoveOrderParams{
		HTTPRequest: request,
		MoveOrderID: strfmt.UUID(moveOrder.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveOrdersHandler{
		context,
		moveorder.NewMoveOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveOrderOK := response.(*moveorderop.GetMoveOrderOK)
	moveOrdersPayload := moveOrderOK.Payload

	suite.Assertions.IsType(&moveorderop.GetMoveOrderOK{}, response)
	suite.Equal(moveOrder.ID.String(), moveOrdersPayload.ID.String())
	suite.Equal(moveOrder.ServiceMemberID.String(), moveOrdersPayload.CustomerID.String())
	suite.Equal(moveOrder.NewDutyStationID.String(), moveOrdersPayload.DestinationDutyStation.ID.String())
	suite.NotNil(moveOrder.NewDutyStation)
	payloadEntitlement := moveOrdersPayload.Entitlement
	suite.Equal((*moveOrder.EntitlementID).String(), payloadEntitlement.ID.String())
	moveOrderEntitlement := moveOrder.Entitlement
	suite.NotNil(moveOrderEntitlement)
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
	suite.Equal(int64(moveOrderEntitlement.WeightAllotment().TotalWeightSelf), payloadEntitlement.TotalWeight)
	suite.Equal(int64(*moveOrderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	suite.Equal(moveOrder.OriginDutyStation.ID.String(), moveOrdersPayload.OriginDutyStation.ID.String())
	suite.NotZero(moveOrder.OriginDutyStation)
	suite.NotZero(moveOrdersPayload.DateIssued)
}
