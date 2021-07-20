package payloads

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
)

func (suite *PayloadsSuite) TestMTOServiceItemModel() {
	moveTaskOrderIDField, _ := uuid.NewV4()
	moveTaskOrderIDString := strfmt.UUID(moveTaskOrderIDField.String())
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := strfmt.UUID(mtoShipmentIDField.String())

	basicServieItem := &primemessages.MTOServiceItemBasic{
		ReServiceCode: "FSC",
	}

	basicServieItem.SetMoveTaskOrderID(&moveTaskOrderIDString)
	basicServieItem.SetMoveTaskOrderID(&mtoShipmentIDString)

	suite.T().Run("Successfully returns a basic service item model", func(t *testing.T) {
		model, err := MTOServiceItemModel(basicServieItem)

		suite.NoError(err)
		suite.NotEqual(uuid.Nil, model.ID)
		suite.Equal(moveTaskOrderIDField, model.MoveTaskOrderID)
		suite.Equal(mtoShipmentIDField, model.MTOShipmentID)
		suite.Equal(basicServieItem.ReServiceCode, model.ReService.Code)
	})
}
