package payloads

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestMTOServiceItemModel() {
	moveTaskOrderIDField, _ := uuid.NewV4()
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	// Basic Service Item
	basicServieItem := &primemessages.MTOServiceItemBasic{
		ReServiceCode: primemessages.NewReServiceCode("FSC"),
	}

	basicServieItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	basicServieItem.SetMtoShipmentID(*mtoShipmentIDString)

	// DCRT Service Item
	itemMeasurement := int32(1100)
	crateMeasurement := int32(1200)
	dcrtCode := "DCRT"
	reason := "Reason"
	description := "Description"

	item := &primemessages.MTOServiceItemDimension{
		Height: &itemMeasurement,
		Width:  &itemMeasurement,
		Length: &itemMeasurement,
	}

	crate := &primemessages.MTOServiceItemDimension{
		Height: &crateMeasurement,
		Width:  &crateMeasurement,
		Length: &crateMeasurement,
	}

	DCRTServiceItem := &primemessages.MTOServiceItemDomesticCrating{
		ReServiceCode: &dcrtCode,
		Reason:        &reason,
		Description:   &description,
	}
	DCRTServiceItem.Item.MTOServiceItemDimension = *item
	DCRTServiceItem.Crate.MTOServiceItemDimension = *crate

	DCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	DCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	suite.T().Run("Success - Returns a basic service item model", func(t *testing.T) {
		returnedModel, verrs := MTOServiceItemModel(basicServieItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCode("FSC"), returnedModel.ReService.Code)
	})

	suite.T().Run("Success - Returns a DCRT service item model", func(t *testing.T) {
		returnedModel, verrs := MTOServiceItemModel(DCRTServiceItem)

		var returnedItem, returnedCrate models.MTOServiceItemDimension
		for _, dimension := range returnedModel.Dimensions {
			if dimension.Type == models.DimensionTypeItem {
				returnedItem = dimension
			} else {
				returnedCrate = dimension
			}
		}

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCode("DCRT"), returnedModel.ReService.Code)
		suite.Equal(DCRTServiceItem.Reason, returnedModel.Reason)
		suite.Equal(DCRTServiceItem.Description, returnedModel.Description)
		suite.Equal(unit.ThousandthInches(*DCRTServiceItem.Item.Length), returnedItem.Length)
		suite.Equal(unit.ThousandthInches(*DCRTServiceItem.Crate.Length), returnedCrate.Length)
	})

	suite.T().Run("Fail -  Returns error for DCRT service item because of validation error", func(t *testing.T) {
		badCrateMeasurement := int32(200)
		badCrate := &primemessages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badDCRTServiceItem := &primemessages.MTOServiceItemDomesticCrating{
			ReServiceCode: &dcrtCode,
			Reason:        &reason,
			Description:   &description,
		}
		badDCRTServiceItem.Item.MTOServiceItemDimension = *item
		badDCRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badDCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badDCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		returnedModel, verrs := MTOServiceItemModel(badDCRTServiceItem)

		suite.True(verrs.HasAny(), "invalid crate dimensions for DCRT service item")
		suite.Nil(returnedModel, "returned a model when erroring")

	})
}

func (suite *PayloadsSuite) TestReweighModelFromUpdate() {
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	idField, _ := uuid.NewV4()
	idString := handlers.FmtUUID(idField)

	verificationReason := "Because I said so"
	weight := int64(2000)

	reweigh := &primemessages.UpdateReweigh{
		VerificationReason: &verificationReason,
		Weight:             &weight,
	}

	suite.T().Run("Success - Returns a reweigh model", func(t *testing.T) {
		returnedModel := ReweighModelFromUpdate(reweigh, *idString, *mtoShipmentIDString)

		suite.Equal(idField.String(), returnedModel.ID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.ShipmentID.String())
		suite.Equal(handlers.PoundPtrFromInt64Ptr(reweigh.Weight), returnedModel.Weight)
		suite.Equal(reweigh.VerificationReason, returnedModel.VerificationReason)
	})

}

func (suite *PayloadsSuite) TestSITExtensionModel() {
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	daysRequested := int64(30)
	remarks := "We need an extension"
	reason := "AWAITING_COMPLETION_OF_RESIDENCE"

	sitExtension := &primemessages.CreateSITExtension{
		RequestedDays:     &daysRequested,
		ContractorRemarks: &remarks,
		RequestReason:     &reason,
	}

	suite.T().Run("Success - Returns a sit extension model", func(t *testing.T) {
		returnedModel := SITExtensionModel(sitExtension, *mtoShipmentIDString)

		suite.Equal(mtoShipmentIDField, returnedModel.MTOShipmentID)
		suite.Equal(int(daysRequested), returnedModel.RequestedDays)
		suite.Equal(models.SITExtensionRequestReasonAwaitingCompletionOfResidence, returnedModel.RequestReason)
		suite.Equal(sitExtension.ContractorRemarks, returnedModel.ContractorRemarks)
	})

}
