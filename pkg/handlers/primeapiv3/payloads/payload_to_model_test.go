package payloads

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestMTOServiceItemModel() {
	moveTaskOrderIDField, _ := uuid.NewV4()
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	// Basic Service Item
	basicServiceItem := &primev3messages.MTOServiceItemBasic{
		ReServiceCode: primev3messages.NewReServiceCode(primev3messages.ReServiceCode(models.ReServiceCodeFSC)),
	}

	basicServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	basicServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	// DCRT Service Item
	itemMeasurement := int32(1100)
	crateMeasurement := int32(1200)
	dcrtCode := models.ReServiceCodeDCRT.String()
	reason := "Reason"
	description := "Description"
	standaloneCrate := false

	item := &primev3messages.MTOServiceItemDimension{
		Height: &itemMeasurement,
		Width:  &itemMeasurement,
		Length: &itemMeasurement,
	}

	crate := &primev3messages.MTOServiceItemDimension{
		Height: &crateMeasurement,
		Width:  &crateMeasurement,
		Length: &crateMeasurement,
	}

	DCRTServiceItem := &primev3messages.MTOServiceItemDomesticCrating{
		ReServiceCode:   &dcrtCode,
		Reason:          &reason,
		Description:     &description,
		StandaloneCrate: &standaloneCrate,
	}
	DCRTServiceItem.Item.MTOServiceItemDimension = *item
	DCRTServiceItem.Crate.MTOServiceItemDimension = *crate

	DCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	DCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	destReason := "service member will pick up from storage at destination"
	destServiceCode := models.ReServiceCodeDDFSIT.String()
	destDate := strfmt.Date(time.Now())
	destTime := "1400Z"
	destCity := "Beverly Hills"
	destPostalCode := "90210"
	destStreet := "123 Rodeo Dr."
	sitFinalDestAddress := primev3messages.Address{
		City:           &destCity,
		PostalCode:     &destPostalCode,
		StreetAddress1: &destStreet,
	}

	destServiceItem := &primev3messages.MTOServiceItemDestSIT{
		ReServiceCode:               &destServiceCode,
		FirstAvailableDeliveryDate1: &destDate,
		FirstAvailableDeliveryDate2: &destDate,
		DateOfContact1:              &destDate,
		DateOfContact2:              &destDate,
		TimeMilitary1:               &destTime,
		TimeMilitary2:               &destTime,
		SitDestinationFinalAddress:  &sitFinalDestAddress,
		Reason:                      &destReason,
	}

	destServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	destServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	suite.Run("Success - Returns a basic service item model", func() {
		returnedModel, verrs := MTOServiceItemModel(basicServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeFSC, returnedModel.ReService.Code)
	})

	suite.Run("Success - Returns a DCRT service item model", func() {
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
		suite.Equal(models.ReServiceCodeDCRT, returnedModel.ReService.Code)
		suite.Equal(DCRTServiceItem.Reason, returnedModel.Reason)
		suite.Equal(DCRTServiceItem.Description, returnedModel.Description)
		suite.Equal(unit.ThousandthInches(*DCRTServiceItem.Item.Length), returnedItem.Length)
		suite.Equal(unit.ThousandthInches(*DCRTServiceItem.Crate.Length), returnedCrate.Length)
	})

	suite.Run("Fail -  Returns error for DCRT service item because of validation error", func() {
		badCrateMeasurement := int32(200)
		badCrate := &primev3messages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badDCRTServiceItem := &primev3messages.MTOServiceItemDomesticCrating{
			ReServiceCode:   &dcrtCode,
			Reason:          &reason,
			Description:     &description,
			StandaloneCrate: &standaloneCrate,
		}
		badDCRTServiceItem.Item.MTOServiceItemDimension = *item
		badDCRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badDCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badDCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		returnedModel, verrs := MTOServiceItemModel(badDCRTServiceItem)

		suite.True(verrs.HasAny(), fmt.Sprintf("invalid crate dimensions for %s service item", models.ReServiceCodeDCRT))
		suite.Nil(returnedModel, "returned a model when erroring")

	})

	suite.Run("Success - Returns a ICRT/IUCRT service item model", func() {
		// ICRT
		icrtCode := models.ReServiceCodeICRT.String()
		externalCrate := false
		ICRTServiceItem := &primev3messages.MTOServiceItemInternationalCrating{
			ReServiceCode:   &icrtCode,
			Reason:          &reason,
			Description:     &description,
			StandaloneCrate: &standaloneCrate,
			ExternalCrate:   &externalCrate,
		}
		ICRTServiceItem.Item.MTOServiceItemDimension = *item
		ICRTServiceItem.Crate.MTOServiceItemDimension = *crate

		ICRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		ICRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		returnedModel, verrs := MTOServiceItemModel(ICRTServiceItem)

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
		suite.Equal(models.ReServiceCodeICRT, returnedModel.ReService.Code)
		suite.Equal(ICRTServiceItem.Reason, returnedModel.Reason)
		suite.Equal(ICRTServiceItem.Description, returnedModel.Description)
		suite.Equal(ICRTServiceItem.StandaloneCrate, returnedModel.StandaloneCrate)
		suite.Equal(ICRTServiceItem.ExternalCrate, returnedModel.ExternalCrate)
		suite.Equal(unit.ThousandthInches(*ICRTServiceItem.Item.Length), returnedItem.Length)
		suite.Equal(unit.ThousandthInches(*ICRTServiceItem.Crate.Length), returnedCrate.Length)

		// IUCRT
		iucrtCode := models.ReServiceCodeIUCRT.String()
		IUCRTServiceItem := &primev3messages.MTOServiceItemInternationalCrating{
			ReServiceCode: &iucrtCode,
			Reason:        &reason,
			Description:   &description,
		}
		IUCRTServiceItem.Item.MTOServiceItemDimension = *item
		IUCRTServiceItem.Crate.MTOServiceItemDimension = *crate

		IUCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		IUCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		iucrtReturnedModel, verrs := MTOServiceItemModel(IUCRTServiceItem)

		var icurtReturnedItem, icurtReturnedCrate models.MTOServiceItemDimension
		for _, dimension := range iucrtReturnedModel.Dimensions {
			if dimension.Type == models.DimensionTypeItem {
				icurtReturnedItem = dimension
			} else {
				icurtReturnedCrate = dimension
			}
		}

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), iucrtReturnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), iucrtReturnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeIUCRT, iucrtReturnedModel.ReService.Code)
		suite.Equal(IUCRTServiceItem.Reason, iucrtReturnedModel.Reason)
		suite.Equal(IUCRTServiceItem.Description, iucrtReturnedModel.Description)
		suite.Equal(unit.ThousandthInches(*ICRTServiceItem.Item.Length), icurtReturnedItem.Length)
		suite.Equal(unit.ThousandthInches(*ICRTServiceItem.Crate.Length), icurtReturnedCrate.Length)
	})

	suite.Run("Fail -  Returns error for ICRT/IUCRT service item because of validation error", func() {
		// ICRT
		icrtCode := models.ReServiceCodeICRT.String()
		externalCrate := false
		badCrateMeasurement := int32(200)
		badCrate := &primev3messages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badICRTServiceItem := &primev3messages.MTOServiceItemInternationalCrating{
			ReServiceCode:   &icrtCode,
			Reason:          &reason,
			Description:     &description,
			StandaloneCrate: &standaloneCrate,
			ExternalCrate:   &externalCrate,
		}
		badICRTServiceItem.Item.MTOServiceItemDimension = *item
		badICRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badICRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badICRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		returnedModel, verrs := MTOServiceItemModel(badICRTServiceItem)

		suite.True(verrs.HasAny(), fmt.Sprintf("invalid crate dimensions for %s service item", models.ReServiceCodeICRT))
		suite.Nil(returnedModel, "returned a model when erroring")

		// IUCRT
		iucrtCode := models.ReServiceCodeIUCRT.String()

		badIUCRTServiceItem := &primev3messages.MTOServiceItemInternationalCrating{
			ReServiceCode:   &iucrtCode,
			Reason:          &reason,
			Description:     &description,
			StandaloneCrate: &standaloneCrate,
			ExternalCrate:   &externalCrate,
		}
		badIUCRTServiceItem.Item.MTOServiceItemDimension = *item
		badIUCRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badIUCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badIUCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		iucrtReturnedModel, verrs := MTOServiceItemModel(badIUCRTServiceItem)

		suite.True(verrs.HasAny(), fmt.Sprintf("invalid crate dimensions for %s service item", models.ReServiceCodeIUCRT))
		suite.Nil(iucrtReturnedModel, "returned a model when erroring")
	})

	suite.Run("Success - Returns SIT destination service item model", func() {
		destSITServiceItem := &primev3messages.MTOServiceItemDestSIT{
			ReServiceCode:               &destServiceCode,
			FirstAvailableDeliveryDate1: &destDate,
			FirstAvailableDeliveryDate2: &destDate,
			DateOfContact1:              &destDate,
			DateOfContact2:              &destDate,
			TimeMilitary1:               &destTime,
			TimeMilitary2:               &destTime,
			SitDestinationFinalAddress:  &sitFinalDestAddress,
			Reason:                      &destReason,
		}

		destSITServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		destSITServiceItem.SetMtoShipmentID(*mtoShipmentIDString)
		returnedModel, verrs := MTOServiceItemModel(destSITServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDDFSIT, returnedModel.ReService.Code)
		suite.Equal(destPostalCode, returnedModel.SITDestinationFinalAddress.PostalCode)
		suite.Equal(destStreet, returnedModel.SITDestinationFinalAddress.StreetAddress1)
	})

	suite.Run("Success - Returns SIT destination service item model without customer contact fields", func() {
		destSITServiceItem := &primev3messages.MTOServiceItemDestSIT{
			ReServiceCode:              &destServiceCode,
			SitDestinationFinalAddress: &sitFinalDestAddress,
			Reason:                     &destReason,
		}

		destSITServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		destSITServiceItem.SetMtoShipmentID(*mtoShipmentIDString)
		returnedModel, verrs := MTOServiceItemModel(destSITServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDDFSIT, returnedModel.ReService.Code)
		suite.Equal(destPostalCode, returnedModel.SITDestinationFinalAddress.PostalCode)
		suite.Equal(destStreet, returnedModel.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(destReason, *returnedModel.Reason)
	})
}

func (suite *PayloadsSuite) TestReweighModelFromUpdate() {
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	idField, _ := uuid.NewV4()
	idString := handlers.FmtUUID(idField)

	verificationReason := "Because I said so"
	weight := int64(2000)

	reweigh := &primev3messages.UpdateReweigh{
		VerificationReason: &verificationReason,
		Weight:             &weight,
	}

	suite.Run("Success - Returns a reweigh model", func() {
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

	sitExtension := &primev3messages.CreateSITExtension{
		RequestedDays:     &daysRequested,
		ContractorRemarks: &remarks,
		RequestReason:     &reason,
	}

	suite.Run("Success - Returns a sit extension model", func() {
		returnedModel := SITExtensionModel(sitExtension, *mtoShipmentIDString)

		suite.Equal(mtoShipmentIDField, returnedModel.MTOShipmentID)
		suite.Equal(int(daysRequested), returnedModel.RequestedDays)
		suite.Equal(models.SITExtensionRequestReasonAwaitingCompletionOfResidence, returnedModel.RequestReason)
		suite.Equal(sitExtension.ContractorRemarks, returnedModel.ContractorRemarks)
	})

}

func (suite *PayloadsSuite) TestMTOAgentModel() {
	suite.Run("success", func() {
		mtoAgentMsg := &primev3messages.MTOAgent{
			ID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
		}

		mtoAgentModel := MTOAgentModel(mtoAgentMsg)

		suite.NotNil(mtoAgentModel)
	})

	suite.Run("unsuccessful", func() {
		mtoAgentModel := MTOAgentModel(nil)
		suite.Nil(mtoAgentModel)
	})
}

func (suite *PayloadsSuite) TestMTOAgentsModel() {
	suite.Run("success", func() {
		mtoAgentsMsg := &primev3messages.MTOAgents{
			{
				ID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			},
			{
				ID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			},
		}

		mtoAgentsModel := MTOAgentsModel(mtoAgentsMsg)

		suite.NotNil(mtoAgentsModel)
		suite.Len(*mtoAgentsModel, len(*mtoAgentsMsg))

		for i, agentModel := range *mtoAgentsModel {
			agentMsg := (*mtoAgentsMsg)[i]
			suite.Equal(agentMsg.ID.String(), agentModel.ID.String())
		}
	})

	suite.Run("unsuccessful", func() {
		mtoAgentsModel := MTOAgentsModel(nil)
		suite.Nil(mtoAgentsModel)
	})
}

func (suite *PayloadsSuite) TestMTOServiceItemModelListFromCreate() {
	suite.Run("successful", func() {
		mtoShipment := &primev3messages.CreateMTOShipment{}

		serviceItemsList, verrs := MTOServiceItemModelListFromCreate(mtoShipment)

		suite.Nil(verrs)
		suite.NotNil(serviceItemsList)
		suite.Len(serviceItemsList, len(mtoShipment.MtoServiceItems()))
	})

	suite.Run("successful multiple items", func() {
		mtoShipment := &primev3messages.CreateMTOShipment{}

		serviceItemsList, verrs := MTOServiceItemModelListFromCreate(mtoShipment)

		suite.Nil(verrs)
		suite.NotNil(serviceItemsList)
		suite.Len(serviceItemsList, len(mtoShipment.MtoServiceItems()))
	})

	suite.Run("unsuccessful", func() {
		serviceItemsList, verrs := MTOServiceItemModelListFromCreate(nil)
		suite.Nil(verrs)
		suite.Nil(serviceItemsList)
	})
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromUpdate() {
	suite.Run("nil", func() {
		model := MTOShipmentModelFromUpdate(nil, strfmt.UUID(uuid.Must(uuid.NewV4()).String()))
		suite.Nil(model)
	})

	suite.Run("notnil", func() {
		mtoShipment := &primev3messages.UpdateMTOShipment{}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		actualWeight := int64(1000)
		ntsRecordedWeight := int64(2000)
		estimatedWeight := int64(1500)
		mtoShipment := &primev3messages.UpdateMTOShipment{
			PrimeActualWeight:    &actualWeight,
			NtsRecordedWeight:    &ntsRecordedWeight,
			PrimeEstimatedWeight: &estimatedWeight,
		}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model.PrimeActualWeight)
		suite.NotNil(model.NTSRecordedWeight)
		suite.NotNil(model.PrimeEstimatedWeight)
	})

	suite.Run("ppm", func() {
		mtoShipment := &primev3messages.UpdateMTOShipment{
			PpmShipment: &primev3messages.UpdatePPMShipment{},
		}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model.PPMShipment)
	})
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate() {
	suite.Run("nil", func() {
		model, err := MTOShipmentModelFromCreate(nil)
		suite.Nil(model)
		suite.NotNil(err)
	})

	suite.Run("empty but not nil", func() {
		mtoShipment := &primev3messages.CreateMTOShipment{}
		model, err := MTOShipmentModelFromCreate(mtoShipment)

		suite.Nil(model)
		suite.NotNil(err)
	})

	pickupAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "CA",
		PostalCode:     "90210",
	}

	pickupAddressMessage := primev3messages.Address{
		StreetAddress1: &pickupAddress.StreetAddress1,
		City:           &pickupAddress.City,
		State:          &pickupAddress.State,
		PostalCode:     &pickupAddress.PostalCode,
	}

	destinationAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "IL",
		PostalCode:     "62225",
	}

	destinationAddressMessage := primev3messages.Address{
		StreetAddress1: &destinationAddress.StreetAddress1,
		City:           &destinationAddress.City,
		State:          &destinationAddress.State,
		PostalCode:     &destinationAddress.PostalCode,
	}

	agentEmail := "test@test.gov"
	firstName := "John"
	lastName := "Doe"
	phone := "123-456-7890"
	agent := primev3messages.MTOAgent{
		AgentType: "type",
		Email:     &agentEmail,
		FirstName: &firstName,
		LastName:  &lastName,
		Phone:     &phone,
	}

	agents := primev3messages.MTOAgents{
		&agent,
	}

	moveUuid, _ := uuid.NewV4()
	moveUuidString := moveUuid.String()
	divertedUuid, _ := uuid.NewV4()
	divertedUuidString := divertedUuid.String()

	createMTOShipmentMessage := &primev3messages.CreateMTOShipment{
		MoveTaskOrderID:        (*strfmt.UUID)(&moveUuidString),
		Agents:                 agents,
		CustomerRemarks:        nil,
		PointOfContact:         "John Doe",
		PrimeEstimatedWeight:   handlers.FmtInt64(1200),
		RequestedPickupDate:    handlers.FmtDatePtr(models.TimePointer(time.Now())),
		ShipmentType:           primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeBOATHAULAWAY),
		PickupAddress:          struct{ primev3messages.Address }{pickupAddressMessage},
		DestinationAddress:     struct{ primev3messages.Address }{destinationAddressMessage},
		DivertedFromShipmentID: (strfmt.UUID)(divertedUuidString),
	}

	suite.Run("with actual payload", func() {
		model, err := MTOShipmentModelFromCreate(createMTOShipmentMessage)

		suite.Nil(err)

		suite.NotNil(model.PickupAddress)
		suite.NotNil(model.DestinationAddress)
		suite.NotNil(model.ShipmentType)
		suite.NotNil(model.PrimeEstimatedWeight)
	})
}

func (suite *PayloadsSuite) TestServiceRequestDocumentUploadModel() {
	upload := models.Upload{
		Bytes:       0,
		ContentType: "",
		Filename:    "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := ServiceRequestDocumentUploadModel(upload)

	suite.Equal(upload.Bytes, *result.Bytes)
	suite.Equal(upload.ContentType, *result.ContentType)
	suite.Equal(upload.Filename, *result.Filename)
	suite.Equal((strfmt.DateTime)(upload.CreatedAt), result.CreatedAt)
	suite.Equal((strfmt.DateTime)(upload.UpdatedAt), result.UpdatedAt)
}

func (suite *PayloadsSuite) TestMTOServiceItemModelFromUpdate() {
	suite.Run("DDDSIT", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		reServiceCode := string(models.ReServiceCodeDDDSIT)
		updateMTOServiceItemSIT := primev3messages.UpdateMTOServiceItemSIT{
			ReServiceCode: reServiceCode,
		}

		model, _ := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceItemSIT)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		estimatedWeight := int64(5000)
		actualWeight := int64(4500)
		updateMTOServiceItemShuttle := primev3messages.UpdateMTOServiceItemShuttle{
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		}

		model, _ := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceItemShuttle)

		suite.NotNil(model)
	})
}

func (suite *PayloadsSuite) TestValidateReasonOriginSIT() {
	suite.Run("Reason provided", func() {
		reason := "reason"
		mtoServiceItemOriginSIT := primev3messages.MTOServiceItemOriginSIT{
			Reason: &reason,
		}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.False(verrs.HasAny())
	})

	suite.Run("No reason provided", func() {
		mtoServiceItemOriginSIT := primev3messages.MTOServiceItemOriginSIT{}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.True(verrs.HasAny())
	})
}

func (suite *PayloadsSuite) TestShipmentAddressUpdateModel() {
	shipmentID := uuid.Must(uuid.NewV4())
	contractorRemarks := ""
	newAddress := primev3messages.Address{
		City:           handlers.FmtString(""),
		State:          handlers.FmtString(""),
		PostalCode:     handlers.FmtString(""),
		StreetAddress1: handlers.FmtString(""),
	}

	nonSITAddressUpdate := primev3messages.UpdateShipmentDestinationAddress{
		ContractorRemarks: &contractorRemarks,
		NewAddress:        &newAddress,
	}

	model := ShipmentAddressUpdateModel(&nonSITAddressUpdate, shipmentID)

	suite.Equal(shipmentID, model.ShipmentID)
	suite.Equal(contractorRemarks, model.ContractorRemarks)
	suite.NotNil(model.NewAddress)
	suite.Equal(*newAddress.City, model.NewAddress.City)
	suite.Equal(*newAddress.State, model.NewAddress.State)
	suite.Equal(*newAddress.PostalCode, model.NewAddress.PostalCode)
	suite.Equal(*newAddress.StreetAddress1, model.NewAddress.StreetAddress1)
}

func (suite *PayloadsSuite) TestPPMShipmentModelFromCreate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)
	sitExpected := true
	estimatedWeight := int64(5000)
	hasProGear := true
	proGearWeight := int64(500)
	spouseProGearWeight := int64(50)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	address2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "11111",
	}
	address3 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "54321",
	}

	var pickupAddress primev3messages.Address
	var secondaryPickupAddress primev3messages.Address
	var tertiaryPickupAddress primev3messages.Address
	var destinationAddress primev3messages.PPMDestinationAddress
	var secondaryDestinationAddress primev3messages.Address
	var tertiaryDestinationAddress primev3messages.Address

	pickupAddress = primev3messages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = primev3messages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	secondaryPickupAddress = primev3messages.Address{
		City:           &address2.City,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	secondaryDestinationAddress = primev3messages.Address{
		City:           &address2.City,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	tertiaryPickupAddress = primev3messages.Address{
		City:           &address3.City,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}
	tertiaryDestinationAddress = primev3messages.Address{
		City:           &address3.City,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}

	ppmShipment := primev3messages.CreatePPMShipment{
		ExpectedDepartureDate:  expectedDepartureDate,
		PickupAddress:          struct{ primev3messages.Address }{pickupAddress},
		SecondaryPickupAddress: struct{ primev3messages.Address }{secondaryPickupAddress},
		TertiaryPickupAddress:  struct{ primev3messages.Address }{tertiaryPickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
		SecondaryDestinationAddress:  struct{ primev3messages.Address }{secondaryDestinationAddress},
		TertiaryDestinationAddress:   struct{ primev3messages.Address }{tertiaryDestinationAddress},
		SitExpected:                  &sitExpected,
		EstimatedWeight:              &estimatedWeight,
		HasProGear:                   &hasProGear,
		ProGearWeight:                &proGearWeight,
		SpouseProGearWeight:          &spouseProGearWeight,
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}

	model := PPMShipmentModelFromCreate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(models.PPMShipmentStatusSubmitted, model.Status)
	suite.True(*model.SITExpected)
	suite.Equal(unit.Pound(estimatedWeight), *model.EstimatedWeight)
	suite.True(*model.HasProGear)
	suite.Equal(unit.Pound(proGearWeight), *model.ProGearWeight)
	suite.Equal(unit.Pound(spouseProGearWeight), *model.SpouseProGearWeight)
	suite.True(*model.HasSecondaryPickupAddress)
	suite.True(*model.HasSecondaryDestinationAddress)
	suite.True(*model.HasTertiaryPickupAddress)
	suite.True(*model.HasTertiaryDestinationAddress)
	suite.True(*model.IsActualExpenseReimbursement)
	suite.NotNil(model)
}

func (suite *PayloadsSuite) TestPPMShipmentModelWithOptionalDestinationStreet1FromCreate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)

	country := models.Country{
		Country:     "US",
		CountryName: "United States",
	}

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
		Country:        &country,
	}

	var pickupAddress primev3messages.Address
	var destinationAddress primev3messages.PPMDestinationAddress

	pickupAddress = primev3messages.Address{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = primev3messages.PPMDestinationAddress{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model := PPMShipmentModelFromUpdate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model2 := PPMShipmentModelFromUpdate(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 2
	streetAddress1 := "1234 Street"
	destinationAddress.StreetAddress1 = &streetAddress1
	ppmShipmentValidDestinatonStreet1 := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model3 := PPMShipmentModelFromUpdate(&ppmShipmentValidDestinatonStreet1)
	suite.Equal(model3.DestinationAddress.StreetAddress1, streetAddress1)
}

func (suite *PayloadsSuite) TestPPMShipmentModelFromUpdate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)
	estimatedWeight := int64(5000)
	proGearWeight := int64(500)
	spouseProGearWeight := int64(50)

	country := models.Country{
		Country:     "US",
		CountryName: "United States",
	}

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
		Country:        &country,
	}
	address2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "11111",
		Country:        &country,
	}
	address3 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "54321",
		Country:        &country,
	}

	var pickupAddress primev3messages.Address
	var secondaryPickupAddress primev3messages.Address
	var tertiaryPickupAddress primev3messages.Address
	var destinationAddress primev3messages.PPMDestinationAddress
	var secondaryDestinationAddress primev3messages.Address
	var tertiaryDestinationAddress primev3messages.Address

	pickupAddress = primev3messages.Address{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = primev3messages.PPMDestinationAddress{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	secondaryPickupAddress = primev3messages.Address{
		City:           &address2.City,
		Country:        &address2.Country.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	secondaryDestinationAddress = primev3messages.Address{
		City:           &address2.City,
		Country:        &address2.Country.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	tertiaryPickupAddress = primev3messages.Address{
		City:           &address3.City,
		Country:        &address3.Country.Country,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}
	tertiaryDestinationAddress = primev3messages.Address{
		City:           &address3.City,
		Country:        &address3.Country.Country,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}

	ppmShipment := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate:  expectedDepartureDate,
		PickupAddress:          struct{ primev3messages.Address }{pickupAddress},
		SecondaryPickupAddress: struct{ primev3messages.Address }{secondaryPickupAddress},
		TertiaryPickupAddress:  struct{ primev3messages.Address }{tertiaryPickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
		SecondaryDestinationAddress:  struct{ primev3messages.Address }{secondaryDestinationAddress},
		TertiaryDestinationAddress:   struct{ primev3messages.Address }{tertiaryDestinationAddress},
		SitExpected:                  models.BoolPointer(true),
		EstimatedWeight:              &estimatedWeight,
		HasProGear:                   models.BoolPointer(true),
		ProGearWeight:                &proGearWeight,
		SpouseProGearWeight:          &spouseProGearWeight,
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}

	model := PPMShipmentModelFromUpdate(&ppmShipment)

	suite.NotNil(model)
	suite.True(*model.SITExpected)
	suite.Equal(unit.Pound(estimatedWeight), *model.EstimatedWeight)
	suite.True(*model.HasProGear)
	suite.Equal(unit.Pound(proGearWeight), *model.ProGearWeight)
	suite.Equal(unit.Pound(spouseProGearWeight), *model.SpouseProGearWeight)
	suite.Nil(model.HasSecondaryPickupAddress)
	suite.Nil(model.HasSecondaryDestinationAddress)
	suite.Nil(model.HasTertiaryPickupAddress)
	suite.Nil(model.HasTertiaryDestinationAddress)
	suite.True(*model.IsActualExpenseReimbursement)
	suite.NotNil(model)
}

func (suite *PayloadsSuite) TestMobileHomeShipmentModelFromCreate() {
	make := "BrandA"
	model := "ModelX"
	year := int64(2024)
	lengthInInches := int64(60)
	heightInInches := int64(13)
	widthInInches := int64(10)

	expectedMobileHome := models.MobileHome{
		Make:           models.StringPointer(make),
		Model:          models.StringPointer(model),
		Year:           models.IntPointer(int(year)),
		LengthInInches: models.IntPointer(int(lengthInInches)),
		HeightInInches: models.IntPointer(int(heightInInches)),
		WidthInInches:  models.IntPointer(int(widthInInches)),
	}

	suite.Run("Success - Complete input", func() {
		mobileHomeShipment := &primev3messages.CreateMobileHomeShipment{
			Make:           models.StringPointer(make),
			Model:          models.StringPointer(model),
			Year:           &year,
			LengthInInches: &lengthInInches,
			HeightInInches: &heightInInches,
			WidthInInches:  &widthInInches,
		}

		moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		mtoShipment := primev3messages.CreateMTOShipment{
			MoveTaskOrderID:    &moveTaskOrderID,
			ShipmentType:       primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeMOBILEHOME),
			MobileHomeShipment: mobileHomeShipment,
		}

		returnedMobileHome, _ := MobileHomeShipmentModelFromCreate(&mtoShipment)

		suite.IsType(&models.MobileHome{}, returnedMobileHome)
		suite.Equal(expectedMobileHome.Make, returnedMobileHome.Make)
		suite.Equal(expectedMobileHome.Model, returnedMobileHome.Model)
		suite.Equal(expectedMobileHome.Year, returnedMobileHome.Year)
		suite.Equal(expectedMobileHome.LengthInInches, returnedMobileHome.LengthInInches)
		suite.Equal(expectedMobileHome.HeightInInches, returnedMobileHome.HeightInInches)
		suite.Equal(expectedMobileHome.WidthInInches, returnedMobileHome.WidthInInches)
	})
}

func (suite *PayloadsSuite) TestBoatShipmentModelFromCreate() {
	make := "BrandA"
	model := "ModelX"
	year := int64(2024)
	lengthInInches := int64(60)
	heightInInches := int64(13)
	widthInInches := int64(10)
	hasTrailer := true
	isRoadworthy := true

	expectedBoatHaulAway := models.BoatShipment{
		Make:           models.StringPointer(make),
		Model:          models.StringPointer(model),
		Year:           models.IntPointer(int(year)),
		LengthInInches: models.IntPointer(int(lengthInInches)),
		HeightInInches: models.IntPointer(int(heightInInches)),
		WidthInInches:  models.IntPointer(int(widthInInches)),
		HasTrailer:     &hasTrailer,
		IsRoadworthy:   &isRoadworthy,
	}

	boatShipment := &primev3messages.CreateBoatShipment{
		Make:           models.StringPointer(make),
		Model:          models.StringPointer(model),
		Year:           &year,
		LengthInInches: &lengthInInches,
		HeightInInches: &heightInInches,
		WidthInInches:  &widthInInches,
		HasTrailer:     &hasTrailer,
		IsRoadworthy:   &isRoadworthy,
	}
	suite.Run("Success - Complete input for MTOShipmentTypeBOATHAULAWAY", func() {
		moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		mtoShipment := primev3messages.CreateMTOShipment{
			MoveTaskOrderID: &moveTaskOrderID,
			ShipmentType:    primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeBOATHAULAWAY),
			BoatShipment:    boatShipment,
		}

		returnedBoatHaulAway, _ := BoatShipmentModelFromCreate(&mtoShipment)

		suite.IsType(&models.BoatShipment{}, returnedBoatHaulAway)

		suite.Equal(expectedBoatHaulAway.Make, returnedBoatHaulAway.Make)
		suite.Equal(expectedBoatHaulAway.Model, returnedBoatHaulAway.Model)
		suite.Equal(expectedBoatHaulAway.Year, returnedBoatHaulAway.Year)
		suite.Equal(expectedBoatHaulAway.LengthInInches, returnedBoatHaulAway.LengthInInches)
		suite.Equal(expectedBoatHaulAway.HeightInInches, returnedBoatHaulAway.HeightInInches)
		suite.Equal(expectedBoatHaulAway.WidthInInches, returnedBoatHaulAway.WidthInInches)
		suite.Equal(expectedBoatHaulAway.HasTrailer, returnedBoatHaulAway.HasTrailer)
		suite.Equal(expectedBoatHaulAway.IsRoadworthy, returnedBoatHaulAway.IsRoadworthy)
	})

	suite.Run("Success - Complete input for MTOShipmentTypeBOATTOWAWAY", func() {
		hasTrailer = false
		isRoadworthy = false

		expectedBoatTowAway := models.BoatShipment{
			Make:           models.StringPointer(make),
			Model:          models.StringPointer(model),
			Year:           models.IntPointer(int(year)),
			LengthInInches: models.IntPointer(int(lengthInInches)),
			HeightInInches: models.IntPointer(int(heightInInches)),
			WidthInInches:  models.IntPointer(int(widthInInches)),
			HasTrailer:     &hasTrailer,
			IsRoadworthy:   &isRoadworthy,
		}

		boatShipment := &primev3messages.CreateBoatShipment{
			Make:           models.StringPointer(make),
			Model:          models.StringPointer(model),
			Year:           &year,
			LengthInInches: &lengthInInches,
			HeightInInches: &heightInInches,
			WidthInInches:  &widthInInches,
			HasTrailer:     &hasTrailer,
			IsRoadworthy:   &isRoadworthy,
		}
		moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		mtoShipment := primev3messages.CreateMTOShipment{
			MoveTaskOrderID: &moveTaskOrderID,
			ShipmentType:    primev3messages.NewMTOShipmentType(primev3messages.MTOShipmentTypeBOATTOWAWAY),
			BoatShipment:    boatShipment,
		}

		returnedBoatTowAway, _ := BoatShipmentModelFromCreate(&mtoShipment)

		suite.IsType(&models.BoatShipment{}, returnedBoatTowAway)

		suite.Equal(expectedBoatTowAway.Make, returnedBoatTowAway.Make)
		suite.Equal(expectedBoatTowAway.Model, returnedBoatTowAway.Model)
		suite.Equal(expectedBoatTowAway.Year, returnedBoatTowAway.Year)
		suite.Equal(expectedBoatTowAway.LengthInInches, returnedBoatTowAway.LengthInInches)
		suite.Equal(expectedBoatTowAway.HeightInInches, returnedBoatTowAway.HeightInInches)
		suite.Equal(expectedBoatTowAway.WidthInInches, returnedBoatTowAway.WidthInInches)
		suite.Equal(expectedBoatTowAway.HasTrailer, returnedBoatTowAway.HasTrailer)
		suite.Equal(expectedBoatTowAway.IsRoadworthy, returnedBoatTowAway.IsRoadworthy)
	})
}

func (suite *PayloadsSuite) TestCountryModel_WithValidCountry() {
	countryName := "US"
	result := CountryModel(&countryName)

	suite.NotNil(result)
	suite.Equal(countryName, result.Country)
}

func (suite *PayloadsSuite) TestCountryModel_WithNilCountry() {
	var countryName *string = nil
	result := CountryModel(countryName)

	suite.Nil(result)
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithNilInput() {
	result, verrs := MTOShipmentModelFromCreate(nil)
	suite.Nil(result)
	suite.NotNil(verrs)
	suite.True(verrs.HasAny())
	suite.Contains(verrs.Keys(), "mtoShipment")
	suite.Equal("mtoShipment object is nil.", verrs.Get("mtoShipment")[0])
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithValidInput() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	mtoShipment := primev3messages.CreateMTOShipment{
		MoveTaskOrderID: &moveTaskOrderID,
	}

	result, _ := MTOShipmentModelFromCreate(&mtoShipment)

	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Nil(result.PrimeEstimatedWeight)
	suite.Nil(result.PickupAddress)
	suite.Nil(result.DestinationAddress)
	suite.Nil(result.SecondaryPickupAddress)
	suite.Nil(result.TertiaryPickupAddress)
	suite.Nil(result.SecondaryDeliveryAddress)
	suite.Nil(result.TertiaryDeliveryAddress)
	suite.Empty(result.MTOAgents)
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithOptionalFields() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	divertedFromShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	primeEstimatedWeight := int64(3000)
	requestedPickupDate := strfmt.Date(time.Now())

	var pickupAddress, secondaryPickupAddress, destinationAddress, tertiaryDestinationAddress primev3messages.Address

	pickupAddress = primev3messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("123 Main St"),
	}

	secondaryPickupAddress = primev3messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("74103"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("789 Elm St"),
	}

	destinationAddress = primev3messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("456 Main St"),
	}

	tertiaryDestinationAddress = primev3messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("74104"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("1010 Oak St"),
	}

	remarks := "customer wants fast delivery"
	mtoShipment := &primev3messages.CreateMTOShipment{
		MoveTaskOrderID:            &moveTaskOrderID,
		CustomerRemarks:            &remarks,
		DivertedFromShipmentID:     divertedFromShipmentID,
		CounselorRemarks:           handlers.FmtString("Approved for special handling"),
		PrimeEstimatedWeight:       &primeEstimatedWeight,
		RequestedPickupDate:        &requestedPickupDate,
		PickupAddress:              struct{ primev3messages.Address }{pickupAddress},
		SecondaryPickupAddress:     struct{ primev3messages.Address }{secondaryPickupAddress},
		DestinationAddress:         struct{ primev3messages.Address }{destinationAddress},
		TertiaryDestinationAddress: struct{ primev3messages.Address }{tertiaryDestinationAddress},
	}

	result, _ := MTOShipmentModelFromCreate(mtoShipment)

	// Check the main fields
	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Equal(*mtoShipment.CustomerRemarks, *result.CustomerRemarks)
	suite.NotNil(result.DivertedFromShipmentID)
	suite.Equal(mtoShipment.DivertedFromShipmentID.String(), result.DivertedFromShipmentID.String())

	// Check weight and recorded date
	suite.NotNil(result.PrimeEstimatedWeight)
	suite.Equal(unit.Pound(primeEstimatedWeight), *result.PrimeEstimatedWeight)
	suite.NotNil(result.PrimeEstimatedWeightRecordedDate)
	suite.WithinDuration(time.Now(), *result.PrimeEstimatedWeightRecordedDate, time.Second)

	// Check pickup and delivery addresses
	suite.NotNil(result.PickupAddress)
	suite.Equal("123 Main St", result.PickupAddress.StreetAddress1)
	suite.NotNil(result.SecondaryPickupAddress)
	suite.Equal("789 Elm St", result.SecondaryPickupAddress.StreetAddress1)
	suite.NotNil(result.DestinationAddress)
	suite.Equal("456 Main St", result.DestinationAddress.StreetAddress1)
	suite.NotNil(result.TertiaryDeliveryAddress)
	suite.Equal("1010 Oak St", result.TertiaryDeliveryAddress.StreetAddress1)
}

func (suite *PayloadsSuite) TestPPMShipmentModelWithOptionalDestinationStreet1FromUpdate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}

	var pickupAddress primev3messages.Address
	var destinationAddress primev3messages.PPMDestinationAddress

	pickupAddress = primev3messages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = primev3messages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model := PPMShipmentModelFromUpdate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model2 := PPMShipmentModelFromUpdate(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 2
	streetAddress1 := "1234 Street"
	destinationAddress.StreetAddress1 = &streetAddress1
	ppmShipmentValidDestinatonStreet1 := primev3messages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ primev3messages.Address }{pickupAddress},
		DestinationAddress: struct {
			primev3messages.PPMDestinationAddress
		}{destinationAddress},
	}

	model3 := PPMShipmentModelFromUpdate(&ppmShipmentValidDestinatonStreet1)
	suite.Equal(model3.DestinationAddress.StreetAddress1, streetAddress1)
}
