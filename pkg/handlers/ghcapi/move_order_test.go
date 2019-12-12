package ghcapi

import (
	"net/http/httptest"

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
	suite.Equal(uuidTostrfmtUUID(moveOrder.EntitlementID), moveOrdersPayload.Entitlement.ID)
	suite.NotZero(moveOrder.Entitlement)
	suite.Equal(uuidTostrfmtUUID(moveOrder.OriginDutyStation.ID), moveOrdersPayload.OriginDutyStation.ID)
	suite.NotZero(moveOrder.OriginDutyStation)
}
