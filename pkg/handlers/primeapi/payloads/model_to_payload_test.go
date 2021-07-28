package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestMoveTaskOrder() {
	moveTaskOrderID, _ := uuid.NewV4()
	ordersID, _ := uuid.NewV4()
	referenceID := "testID"
	primeTime := time.Now()
	submittedAt := time.Now()
	hhgMoveType := models.SelectedMoveTypeHHG

	basicMove := models.Move{
		ID:                      moveTaskOrderID,
		Locator:                 "TESTTEST",
		CreatedAt:               time.Now(),
		AvailableToPrimeAt:      &primeTime,
		OrdersID:                ordersID,
		Orders:                  models.Order{},
		ReferenceID:             &referenceID,
		PaymentRequests:         models.PaymentRequests{},
		SubmittedAt:             &submittedAt,
		UpdatedAt:               time.Now(),
		SelectedMoveType:        &hhgMoveType,
		PersonallyProcuredMoves: models.PersonallyProcuredMoves{},
		MoveDocuments:           models.MoveDocuments{},
		Status:                  models.MoveStatusAPPROVED,
		SignedCertifications:    models.SignedCertifications{},
		MTOServiceItems:         models.MTOServiceItems{},
		MTOShipments:            models.MTOShipments{},
	}

	suite.T().Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func(t *testing.T) {
		returnedModel := MoveTaskOrder(&basicMove)

		suite.IsType(&primemessages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(referenceID, returnedModel.ReferenceID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotNil(returnedModel.ETag)
	})
}
