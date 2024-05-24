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
		ReServiceCode: &dcrtCode,
		Reason:        &reason,
		Description:   &description,
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
			ReServiceCode: &dcrtCode,
			Reason:        &reason,
			Description:   &description,
		}
		badDCRTServiceItem.Item.MTOServiceItemDimension = *item
		badDCRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badDCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badDCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		returnedModel, verrs := MTOServiceItemModel(badDCRTServiceItem)

		suite.True(verrs.HasAny(), fmt.Sprintf("invalid crate dimensions for %s service item", models.ReServiceCodeDCRT))
		suite.Nil(returnedModel, "returned a model when erroring")

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

func (suite *PayloadsSuite) TestSITAddressUpdateModel() {
	contractorRemark := "I must update the final address please"
	city := "Beverly Hills"
	state := "CA"
	postalCode := "90210"
	street := "123 Rodeo Dr."
	newAddress := primev3messages.Address{
		City:           &city,
		State:          &state,
		PostalCode:     &postalCode,
		StreetAddress1: &street,
	}

	suite.Run("Success - Returns a SITAddressUpdate model as expected", func() {
		sitAddressUpdate := primev3messages.CreateSITAddressUpdateRequest{
			MtoServiceItemID:  strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			NewAddress:        &newAddress,
			ContractorRemarks: &contractorRemark,
		}

		model := SITAddressUpdateModel(&sitAddressUpdate)

		suite.Equal(model.MTOServiceItemID.String(), sitAddressUpdate.MtoServiceItemID.String())
		suite.NotNil(model.NewAddressID.String())
		suite.Equal(model.NewAddress.City, *sitAddressUpdate.NewAddress.City)
		suite.Equal(model.NewAddress.State, *sitAddressUpdate.NewAddress.State)
		suite.Equal(model.NewAddress.PostalCode, *sitAddressUpdate.NewAddress.PostalCode)
		suite.Equal(model.NewAddress.StreetAddress1, *sitAddressUpdate.NewAddress.StreetAddress1)
		suite.Equal(*model.ContractorRemarks, *sitAddressUpdate.ContractorRemarks)
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

	var pickupAddress primev3messages.Address
	var secondaryPickupAddress primev3messages.Address
	var destinationAddress primev3messages.Address
	var secondaryDestinationAddress primev3messages.Address

	pickupAddress = primev3messages.Address{
		City:           &address.City,
		Country:        address.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = primev3messages.Address{
		City:           &address.City,
		Country:        address.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	secondaryPickupAddress = primev3messages.Address{
		City:           &address2.City,
		Country:        address2.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	secondaryDestinationAddress = primev3messages.Address{
		City:           &address2.City,
		Country:        address2.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}

	ppmShipment := primev3messages.CreatePPMShipment{
		ExpectedDepartureDate:       expectedDepartureDate,
		PickupAddress:               struct{ primev3messages.Address }{pickupAddress},
		SecondaryPickupAddress:      struct{ primev3messages.Address }{secondaryPickupAddress},
		DestinationAddress:          struct{ primev3messages.Address }{destinationAddress},
		SecondaryDestinationAddress: struct{ primev3messages.Address }{secondaryDestinationAddress},
		SitExpected:                 &sitExpected,
		EstimatedWeight:             &estimatedWeight,
		HasProGear:                  &hasProGear,
		ProGearWeight:               &proGearWeight,
		SpouseProGearWeight:         &spouseProGearWeight,
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
	suite.NotNil(model)
}
