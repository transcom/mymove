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

func (suite *PayloadsSuite) TestReweigh() {
	id, _ := uuid.NewV4()
	shipmentID, _ := uuid.NewV4()
	requestedAt := time.Now()
	createdAt := time.Now()
	updatedAt := time.Now()

	reweigh := models.Reweigh{
		ID:          id,
		ShipmentID:  shipmentID,
		RequestedAt: requestedAt,
		RequestedBy: models.ReweighRequesterTOO,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	suite.T().Run("Success - Returns a basic rewweigh payload without optional fields", func(t *testing.T) {
		returnedPayload := Reweigh(&reweigh)

		suite.IsType(&primemessages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(returnedPayload.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(returnedPayload.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(returnedPayload.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primemessages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
		suite.Equal(strfmt.DateTime(returnedPayload.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(returnedPayload.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Nil(returnedPayload.Weight)
		suite.Nil(returnedPayload.VerificationReason)
		suite.Nil(returnedPayload.VerificationProvidedAt)
		suite.NotNil(returnedPayload.ETag)

	})

	suite.T().Run("Success - Returns a basic rewweigh payload with optional fields", func(t *testing.T) {
		// Set optional fields
		weight := int64(2000)
		reweigh.Weight = handlers.PoundPtrFromInt64Ptr(&weight)

		verificationProvidedAt := time.Now()
		reweigh.VerificationProvidedAt = &verificationProvidedAt

		verificationReason := "Because I said so"
		reweigh.VerificationReason = &verificationReason

		// Send model through func
		returnedPayload := Reweigh(&reweigh)

		suite.IsType(&primemessages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(returnedPayload.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(returnedPayload.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(returnedPayload.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primemessages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
		suite.Equal(strfmt.DateTime(returnedPayload.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(returnedPayload.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Equal(handlers.FmtPoundPtr(reweigh.Weight), returnedPayload.Weight)
		suite.Equal(handlers.FmtStringPtr(reweigh.VerificationReason), returnedPayload.VerificationReason)
		suite.Equal(handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt), returnedPayload.VerificationProvidedAt)
		suite.NotNil(returnedPayload.ETag)

	})
}
