package ghcapi

import (
	"net/http/httptest"

	moveorder "github.com/transcom/mymove/pkg/services/move_order"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	moveorderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/handlers"
)

func uuidTostrfmtUUID(id interface{}) strfmt.UUID {
	if s, ok := id.(uuid.UUID); ok {
		return strfmt.UUID(s.String())
	}
	return ""
}

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
	suite.Equal(uuidTostrfmtUUID(moveOrder.ID), moveOrdersPayload.ID)
	suite.Equal(uuidTostrfmtUUID(moveOrder.CustomerID), moveOrdersPayload.CustomerID)
	suite.Equal(uuidTostrfmtUUID(moveOrder.DestinationDutyStationID), moveOrdersPayload.DestinationDutyStation.ID)
	suite.NotZero(moveOrder.DestinationDutyStation)
	payloadEntitlement := moveOrdersPayload.Entitlement
	suite.Equal(uuidTostrfmtUUID(moveOrder.EntitlementID), payloadEntitlement.ID)
	moveOrderEntitlement := moveOrder.Entitlement
	suite.NotNil(moveOrderEntitlement)
	suite.Equal(int64(moveOrderEntitlement.ProGearWeight), payloadEntitlement.ProGearWeight)
	suite.Equal(int64(moveOrderEntitlement.ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
	suite.Equal(int64(moveOrderEntitlement.TotalWeightSelf), payloadEntitlement.TotalWeight)
	suite.Equal(uuidTostrfmtUUID(moveOrder.OriginDutyStation.ID), moveOrdersPayload.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}
