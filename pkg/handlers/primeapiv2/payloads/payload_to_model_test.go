package payloads

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primev2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestMTOServiceItemModel() {
	moveTaskOrderIDField, _ := uuid.NewV4()
	mtoShipmentIDField, _ := uuid.NewV4()
	mtoShipmentIDString := handlers.FmtUUID(mtoShipmentIDField)

	// Basic Service Item
	basicServiceItem := &primev2messages.MTOServiceItemBasic{
		ReServiceCode: primev2messages.NewReServiceCode(primev2messages.ReServiceCode(models.ReServiceCodeFSC)),
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

	item := &primev2messages.MTOServiceItemDimension{
		Height: &itemMeasurement,
		Width:  &itemMeasurement,
		Length: &itemMeasurement,
	}

	crate := &primev2messages.MTOServiceItemDimension{
		Height: &crateMeasurement,
		Width:  &crateMeasurement,
		Length: &crateMeasurement,
	}

	DCRTServiceItem := &primev2messages.MTOServiceItemDomesticCrating{
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
	sitFinalDestAddress := primev2messages.Address{
		City:           &destCity,
		PostalCode:     &destPostalCode,
		StreetAddress1: &destStreet,
	}

	destServiceItem := &primev2messages.MTOServiceItemDestSIT{
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
		badCrate := &primev2messages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badDCRTServiceItem := &primev2messages.MTOServiceItemDomesticCrating{
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

	suite.Run("Success - Returns SIT destination service item model", func() {
		destSITServiceItem := &primev2messages.MTOServiceItemDestSIT{
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
		destSITServiceItem := &primev2messages.MTOServiceItemDestSIT{
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

	reweigh := &primev2messages.UpdateReweigh{
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

	sitExtension := &primev2messages.CreateSITExtension{
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
		mtoAgentMsg := &primev2messages.MTOAgent{
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
		mtoAgentsMsg := &primev2messages.MTOAgents{
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
		mtoShipment := &primev2messages.CreateMTOShipment{}

		serviceItemsList, verrs := MTOServiceItemModelListFromCreate(mtoShipment)

		suite.Nil(verrs)
		suite.NotNil(serviceItemsList)
		suite.Len(serviceItemsList, len(mtoShipment.MtoServiceItems()))
	})

	suite.Run("successful multiple items", func() {
		mtoShipment := &primev2messages.CreateMTOShipment{}

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
		mtoShipment := &primev2messages.UpdateMTOShipment{}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		actualWeight := int64(1000)
		ntsRecordedWeight := int64(2000)
		estimatedWeight := int64(1500)
		actualProGearWeight := int64(1000)
		actualSpouseProGearWeight := int64(250)
		mtoShipment := &primev2messages.UpdateMTOShipment{
			PrimeActualWeight:         &actualWeight,
			NtsRecordedWeight:         &ntsRecordedWeight,
			PrimeEstimatedWeight:      &estimatedWeight,
			ActualProGearWeight:       &actualProGearWeight,
			ActualSpouseProGearWeight: &actualSpouseProGearWeight,
		}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model.PrimeActualWeight)
		suite.NotNil(model.NTSRecordedWeight)
		suite.NotNil(model.PrimeEstimatedWeight)
		suite.Equal(model.ActualProGearWeight, handlers.PoundPtrFromInt64Ptr(&actualProGearWeight))
		suite.Equal(model.ActualSpouseProGearWeight, handlers.PoundPtrFromInt64Ptr(&actualSpouseProGearWeight))
	})

	suite.Run("ppm", func() {
		mtoShipment := &primev2messages.UpdateMTOShipment{
			PpmShipment: &primev2messages.UpdatePPMShipment{},
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
		updateMTOServiceItemSIT := primev2messages.UpdateMTOServiceItemSIT{
			ReServiceCode: reServiceCode,
		}

		model, _ := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceItemSIT)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		estimatedWeight := int64(5000)
		actualWeight := int64(4500)
		updateMTOServiceItemShuttle := primev2messages.UpdateMTOServiceItemShuttle{
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
		mtoServiceItemOriginSIT := primev2messages.MTOServiceItemOriginSIT{
			Reason: &reason,
		}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.False(verrs.HasAny())
	})

	suite.Run("No reason provided", func() {
		mtoServiceItemOriginSIT := primev2messages.MTOServiceItemOriginSIT{}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.True(verrs.HasAny())
	})
}

func (suite *PayloadsSuite) TestShipmentAddressUpdateModel() {
	shipmentID := uuid.Must(uuid.NewV4())
	contractorRemarks := ""
	newAddress := primev2messages.Address{
		City:           handlers.FmtString(""),
		State:          handlers.FmtString(""),
		PostalCode:     handlers.FmtString(""),
		StreetAddress1: handlers.FmtString(""),
	}

	nonSITAddressUpdate := primev2messages.UpdateShipmentDestinationAddress{
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

	ppmShipment := primev2messages.CreatePPMShipment{
		ExpectedDepartureDate:        expectedDepartureDate,
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
	suite.True(*model.IsActualExpenseReimbursement)
}

func (suite *PayloadsSuite) TestPPMShipmentModelFromUpdate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)
	estimatedWeight := int64(5000)
	proGearWeight := int64(500)
	spouseProGearWeight := int64(50)

	ppmShipment := primev2messages.UpdatePPMShipment{
		ExpectedDepartureDate:        expectedDepartureDate,
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
	suite.True(*model.IsActualExpenseReimbursement)
	suite.NotNil(model)
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
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithValidInput() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	mtoShipment := primev2messages.CreateMTOShipment{
		MoveTaskOrderID: &moveTaskOrderID,
	}

	result, verrs := MTOShipmentModelFromCreate(&mtoShipment)

	suite.Nil(verrs)
	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Nil(result.PrimeEstimatedWeight)
	suite.Nil(result.PickupAddress)
	suite.Nil(result.DestinationAddress)
	suite.Empty(result.MTOAgents)
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithOptionalFields() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	divertedFromShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	primeEstimatedWeight := int64(3000)
	requestedPickupDate := strfmt.Date(time.Now())

	var pickupAddress primev2messages.Address
	var destinationAddress primev2messages.Address

	pickupAddress = primev2messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("123 Main St"),
	}
	destinationAddress = primev2messages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("456 Main St"),
	}

	remarks := "customer wants fast delivery"
	mtoShipment := &primev2messages.CreateMTOShipment{
		MoveTaskOrderID:        &moveTaskOrderID,
		CustomerRemarks:        &remarks,
		DivertedFromShipmentID: divertedFromShipmentID,
		CounselorRemarks:       handlers.FmtString("Approved for special handling"),
		PrimeEstimatedWeight:   &primeEstimatedWeight,
		RequestedPickupDate:    &requestedPickupDate,
		PickupAddress:          struct{ primev2messages.Address }{pickupAddress},
		DestinationAddress:     struct{ primev2messages.Address }{destinationAddress},
	}

	result, verrs := MTOShipmentModelFromCreate(mtoShipment)

	suite.Nil(verrs)
	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Equal(*mtoShipment.CustomerRemarks, *result.CustomerRemarks)
	suite.NotNil(result.DivertedFromShipmentID)
	suite.Equal(mtoShipment.DivertedFromShipmentID.String(), result.DivertedFromShipmentID.String())

	suite.NotNil(result.PrimeEstimatedWeight)
	suite.Equal(unit.Pound(primeEstimatedWeight), *result.PrimeEstimatedWeight)
	suite.NotNil(result.PrimeEstimatedWeightRecordedDate)
	suite.WithinDuration(time.Now(), *result.PrimeEstimatedWeightRecordedDate, time.Second)

	suite.NotNil(result.PickupAddress)
	suite.Equal("123 Main St", result.PickupAddress.StreetAddress1)
	suite.NotNil(result.DestinationAddress)
	suite.Equal("456 Main St", result.DestinationAddress.StreetAddress1)
}
