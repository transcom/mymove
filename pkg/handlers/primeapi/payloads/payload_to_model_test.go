package payloads

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
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
	basicServiceItem := &primemessages.MTOServiceItemBasic{
		ReServiceCode: primemessages.NewReServiceCode(primemessages.ReServiceCode(models.ReServiceCodeFSC)),
	}

	basicServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	basicServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	// DCRT Service Item
	itemMeasurement := int32(1100)
	crateMeasurement := int32(1200)
	estimatedWeight := int64(1000)
	actualWeight := int64(1000)
	dcrtCode := models.ReServiceCodeDCRT.String()
	ddshutCode := models.ReServiceCodeDDSHUT.String()
	doshutCode := models.ReServiceCodeDOSHUT.String()
	idshutCode := models.ReServiceCodeIDSHUT.String()
	ioshutCode := models.ReServiceCodeIOSHUT.String()
	reason := "Reason"
	description := "Description"
	standaloneCrate := false

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
		ReServiceCode:   &dcrtCode,
		Reason:          &reason,
		Description:     &description,
		StandaloneCrate: &standaloneCrate,
	}

	DCRTServiceItem.Item.MTOServiceItemDimension = *item
	DCRTServiceItem.Crate.MTOServiceItemDimension = *crate

	DCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	DCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	DDSHUTServiceItem := &primemessages.MTOServiceItemDomesticShuttle{
		ReServiceCode:   &ddshutCode,
		Reason:          &reason,
		EstimatedWeight: &estimatedWeight,
		ActualWeight:    &actualWeight,
	}
	DDSHUTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	DDSHUTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	DOSHUTServiceItem := &primemessages.MTOServiceItemDomesticShuttle{
		ReServiceCode:   &doshutCode,
		Reason:          &reason,
		EstimatedWeight: &estimatedWeight,
		ActualWeight:    &actualWeight,
	}
	DOSHUTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	DOSHUTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	IDSHUTServiceItem := &primemessages.MTOServiceItemInternationalShuttle{
		ReServiceCode:   &idshutCode,
		Reason:          &reason,
		EstimatedWeight: &estimatedWeight,
		ActualWeight:    &actualWeight,
	}
	IDSHUTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	IDSHUTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	IOSHUTServiceItem := &primemessages.MTOServiceItemInternationalShuttle{
		ReServiceCode:   &ioshutCode,
		Reason:          &reason,
		EstimatedWeight: &estimatedWeight,
		ActualWeight:    &actualWeight,
	}
	IOSHUTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
	IOSHUTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

	originReason := "storage at origin"
	originServiceCode := models.ReServiceCodeDOFSIT.String()
	originSITEntryDate := strfmt.Date(time.Now())
	originSITDepartureDate := strfmt.Date(time.Now())
	originState := "TN"
	originCity := "Beverly Hills"
	originPostalCode := "90210"
	originStreet1 := "123 Rodeo Dr."
	originCounty1 := "LOS ANGELES"
	originUSPRCID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	sitHHGActualOriginAddress := primemessages.Address{
		State:                &originState,
		City:                 &originCity,
		PostalCode:           &originPostalCode,
		StreetAddress1:       &originStreet1,
		County:               &originCounty1,
		UsPostRegionCitiesID: originUSPRCID,
	}

	destReason := "service member will pick up from storage at destination"
	destServiceCode := models.ReServiceCodeDDFSIT.String()
	destDate := strfmt.Date(time.Now())
	destTime := "1400Z"
	destCity := "Beverly Hills"
	destPostalCode := "90210"
	destCounty := "LOS ANGELES"
	destStreet := "123 Rodeo Dr."
	destUSPRCID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	sitFinalDestAddress := primemessages.Address{
		City:                 &destCity,
		PostalCode:           &destPostalCode,
		StreetAddress1:       &destStreet,
		County:               &destCounty,
		UsPostRegionCitiesID: destUSPRCID,
	}

	destServiceItem := &primemessages.MTOServiceItemDestSIT{
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
		badCrate := &primemessages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badDCRTServiceItem := &primemessages.MTOServiceItemDomesticCrating{
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
		ICRTServiceItem := &primemessages.MTOServiceItemInternationalCrating{
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
		IUCRTServiceItem := &primemessages.MTOServiceItemInternationalCrating{
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

	suite.Run("Success - Returns a DDSHUT service item model", func() {
		returnedModel, verrs := MTOServiceItemModel(DDSHUTServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDDSHUT, returnedModel.ReService.Code)
		suite.Equal(DDSHUTServiceItem.Reason, returnedModel.Reason)
	})

	suite.Run("Success - Returns a DOSHUT service item model", func() {
		returnedModel, verrs := MTOServiceItemModel(DOSHUTServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDOSHUT, returnedModel.ReService.Code)
		suite.Equal(DOSHUTServiceItem.Reason, returnedModel.Reason)
	})

	suite.Run("Success - Returns a IOSHUT service item model", func() {
		returnedModel, verrs := MTOServiceItemModel(IOSHUTServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeIOSHUT, returnedModel.ReService.Code)
		suite.Equal(IOSHUTServiceItem.Reason, returnedModel.Reason)
	})

	suite.Run("Success - Returns a IDSHUT service item model", func() {
		returnedModel, verrs := MTOServiceItemModel(IDSHUTServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeIDSHUT, returnedModel.ReService.Code)
		suite.Equal(IDSHUTServiceItem.Reason, returnedModel.Reason)
	})

	suite.Run("Fail -  Returns error for ICRT/IUCRT service item because of validation error", func() {
		// ICRT
		icrtCode := models.ReServiceCodeICRT.String()
		externalCrate := false
		badCrateMeasurement := int32(200)
		badCrate := &primemessages.MTOServiceItemDimension{
			Height: &badCrateMeasurement,
			Width:  &badCrateMeasurement,
			Length: &badCrateMeasurement,
		}

		badICRTServiceItem := &primemessages.MTOServiceItemInternationalCrating{
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

		badIUCRTServiceItem := &primemessages.MTOServiceItemInternationalCrating{
			ReServiceCode: &iucrtCode,
			Reason:        &reason,
			Description:   &description,
		}
		badIUCRTServiceItem.Item.MTOServiceItemDimension = *item
		badIUCRTServiceItem.Crate.MTOServiceItemDimension = *badCrate

		badIUCRTServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		badIUCRTServiceItem.SetMtoShipmentID(*mtoShipmentIDString)

		iucrtReturnedModel, verrs := MTOServiceItemModel(badIUCRTServiceItem)

		suite.True(verrs.HasAny(), fmt.Sprintf("invalid crate dimensions for %s service item", models.ReServiceCodeIUCRT))
		suite.Nil(iucrtReturnedModel, "returned a model when erroring")
	})

	suite.Run("Success - Returns SIT origin service item model", func() {
		originSITServiceItem := &primemessages.MTOServiceItemOriginSIT{
			ReServiceCode:      &originServiceCode,
			SitEntryDate:       &originSITEntryDate,
			SitDepartureDate:   &originSITDepartureDate,
			SitHHGActualOrigin: &sitHHGActualOriginAddress,
			Reason:             &originReason,
		}

		originSITServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		originSITServiceItem.SetMtoShipmentID(*mtoShipmentIDString)
		returnedModel, verrs := MTOServiceItemModel(originSITServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDOFSIT, returnedModel.ReService.Code)
		suite.Equal(originStreet1, returnedModel.SITOriginHHGActualAddress.StreetAddress1)
		suite.Equal(originCity, returnedModel.SITOriginHHGActualAddress.City)
		suite.Equal(originState, returnedModel.SITOriginHHGActualAddress.State)
		suite.Equal(originPostalCode, returnedModel.SITOriginHHGActualAddress.PostalCode)
		suite.Equal(originSITEntryDate, *handlers.FmtDatePtr(returnedModel.SITEntryDate))
		suite.Equal(originSITDepartureDate, *handlers.FmtDatePtr(returnedModel.SITDepartureDate))
	})

	suite.Run("Success - Returns international SIT origin service item model", func() {
		originSITServiceItem := &primemessages.MTOServiceItemInternationalOriginSIT{
			ReServiceCode:      &originServiceCode,
			SitEntryDate:       &originSITEntryDate,
			SitDepartureDate:   &originSITDepartureDate,
			SitHHGActualOrigin: &sitHHGActualOriginAddress,
			Reason:             &originReason,
		}

		originSITServiceItem.SetMoveTaskOrderID(handlers.FmtUUID(moveTaskOrderIDField))
		originSITServiceItem.SetMtoShipmentID(*mtoShipmentIDString)
		returnedModel, verrs := MTOServiceItemModel(originSITServiceItem)

		suite.NoVerrs(verrs)
		suite.Equal(moveTaskOrderIDField.String(), returnedModel.MoveTaskOrderID.String())
		suite.Equal(mtoShipmentIDField.String(), returnedModel.MTOShipmentID.String())
		suite.Equal(models.ReServiceCodeDOFSIT, returnedModel.ReService.Code)
		suite.Equal(originStreet1, returnedModel.SITOriginHHGActualAddress.StreetAddress1)
		suite.Equal(originCity, returnedModel.SITOriginHHGActualAddress.City)
		suite.Equal(originState, returnedModel.SITOriginHHGActualAddress.State)
		suite.Equal(originPostalCode, returnedModel.SITOriginHHGActualAddress.PostalCode)
		suite.Equal(originSITEntryDate, *handlers.FmtDatePtr(returnedModel.SITEntryDate))
		suite.Equal(originSITDepartureDate, *handlers.FmtDatePtr(returnedModel.SITDepartureDate))
	})

	suite.Run("Success - Returns SIT destination service item model", func() {
		destSITServiceItem := &primemessages.MTOServiceItemDestSIT{
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
		suite.Equal(destUSPRCID.String(), returnedModel.SITDestinationFinalAddress.UsPostRegionCityID.String())
	})

	suite.Run("Success - Returns international SIT destination service item model", func() {
		destSITServiceItem := &primemessages.MTOServiceItemInternationalDestSIT{
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
		suite.Equal(destUSPRCID.String(), returnedModel.SITDestinationFinalAddress.UsPostRegionCityID.String())
	})

	suite.Run("Success - Returns SIT destination service item model without customer contact fields", func() {
		destSITServiceItem := &primemessages.MTOServiceItemDestSIT{
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
		suite.Equal(destUSPRCID.String(), returnedModel.SITDestinationFinalAddress.UsPostRegionCityID.String())
		suite.Equal(destReason, *returnedModel.Reason)
	})

	suite.Run("Success - Returns internatonal SIT destination service item model without customer contact fields", func() {
		destSITServiceItem := &primemessages.MTOServiceItemInternationalDestSIT{
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
		suite.Equal(destUSPRCID.String(), returnedModel.SITDestinationFinalAddress.UsPostRegionCityID.String())
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

	reweigh := &primemessages.UpdateReweigh{
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

	sitExtension := &primemessages.CreateSITExtension{
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
		mtoAgentMsg := &primemessages.MTOAgent{
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
		mtoAgentsMsg := &primemessages.MTOAgents{
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
		mtoShipment := &primemessages.CreateMTOShipment{}

		serviceItemsList, verrs := MTOServiceItemModelListFromCreate(mtoShipment)

		suite.Nil(verrs)
		suite.NotNil(serviceItemsList)
		suite.Len(serviceItemsList, len(mtoShipment.MtoServiceItems()))
	})

	suite.Run("successful multiple items", func() {
		mtoShipment := &primemessages.CreateMTOShipment{}

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
		mtoShipment := &primemessages.UpdateMTOShipment{}
		mtoShipmentID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
		model := MTOShipmentModelFromUpdate(mtoShipment, mtoShipmentID)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		actualWeight := int64(1000)
		ntsRecordedWeight := int64(2000)
		estimatedWeight := int64(1500)
		mtoShipment := &primemessages.UpdateMTOShipment{
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
		mtoShipment := &primemessages.UpdateMTOShipment{
			PpmShipment: &primemessages.UpdatePPMShipment{},
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
		updateMTOServiceItemSIT := primemessages.UpdateMTOServiceItemSIT{
			ReServiceCode: reServiceCode,
		}

		model, _ := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceItemSIT)

		suite.NotNil(model)
	})

	suite.Run("weight", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		estimatedWeight := int64(5000)
		actualWeight := int64(4500)
		updateMTOServiceItemShuttle := primemessages.UpdateMTOServiceItemShuttle{
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		}

		model, _ := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceItemShuttle)

		suite.NotNil(model)
	})

	suite.Run("PODFSC", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		portCode := "PDX"
		reServiceCode := string(models.ReServiceCodePODFSC)
		updateMTOServiceInternationalPortFsc := primemessages.UpdateMTOServiceItemInternationalPortFSC{
			ReServiceCode: reServiceCode,
			PortCode:      &portCode,
		}

		model, errs := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceInternationalPortFsc)

		suite.Empty(errs)
		suite.NotNil(model)
		suite.Equal(model.PODLocation.Port.PortCode, portCode)
	})

	suite.Run("POEFSC", func() {
		mtoServiceItemID := uuid.Must(uuid.NewV4()).String()
		portCode := "PDX"
		reServiceCode := string(models.ReServiceCodePOEFSC)
		updateMTOServiceInternationalPortFsc := primemessages.UpdateMTOServiceItemInternationalPortFSC{
			ReServiceCode: reServiceCode,
			PortCode:      &portCode,
		}

		model, errs := MTOServiceItemModelFromUpdate(mtoServiceItemID, &updateMTOServiceInternationalPortFsc)

		suite.Empty(errs)
		suite.NotNil(model)
		suite.Equal(model.POELocation.Port.PortCode, portCode)
	})
}

func (suite *PayloadsSuite) TestValidateReasonOriginSIT() {
	suite.Run("Reason provided", func() {
		reason := "reason"
		mtoServiceItemOriginSIT := primemessages.MTOServiceItemOriginSIT{
			Reason: &reason,
		}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.False(verrs.HasAny())
	})

	suite.Run("No reason provided", func() {
		mtoServiceItemOriginSIT := primemessages.MTOServiceItemOriginSIT{}

		verrs := validateReasonOriginSIT(mtoServiceItemOriginSIT)
		suite.True(verrs.HasAny())
	})
}

func (suite *PayloadsSuite) TestShipmentAddressUpdateModel() {
	shipmentID := uuid.Must(uuid.NewV4())
	contractorRemarks := ""
	newAddress := primemessages.Address{
		City:           handlers.FmtString(""),
		State:          handlers.FmtString(""),
		PostalCode:     handlers.FmtString(""),
		StreetAddress1: handlers.FmtString(""),
	}

	nonSITAddressUpdate := primemessages.UpdateShipmentDestinationAddress{
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

	ppmShipment := primemessages.CreatePPMShipment{
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

	ppmShipment := primemessages.UpdatePPMShipment{
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
	result := MTOShipmentModelFromCreate(nil)
	suite.Nil(result)
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithValidInput() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	mtoShipment := primemessages.CreateMTOShipment{
		MoveTaskOrderID: &moveTaskOrderID,
	}

	result := MTOShipmentModelFromCreate(&mtoShipment)

	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Nil(result.PrimeEstimatedWeight)
	suite.Nil(result.PickupAddress)
	suite.Nil(result.DestinationAddress)
	suite.Empty(result.MTOAgents)
}

func (suite *PayloadsSuite) TestMTOShipmentModelFromCreate_WithOptionalFields() {
	moveTaskOrderID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	primeEstimatedWeight := int64(3000)
	requestedPickupDate := strfmt.Date(time.Now())

	var pickupAddress primemessages.Address
	var destinationAddress primemessages.Address

	pickupAddress = primemessages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("123 Main St"),
	}
	destinationAddress = primemessages.Address{
		City:           handlers.FmtString("Tulsa"),
		PostalCode:     handlers.FmtString("90210"),
		State:          handlers.FmtString("OK"),
		StreetAddress1: handlers.FmtString("456 Main St"),
	}

	remarks := "customer wants fast delivery"
	mtoShipment := &primemessages.CreateMTOShipment{
		MoveTaskOrderID:      &moveTaskOrderID,
		CustomerRemarks:      &remarks,
		CounselorRemarks:     handlers.FmtString("Approved for special handling"),
		PrimeEstimatedWeight: &primeEstimatedWeight,
		RequestedPickupDate:  &requestedPickupDate,
		PickupAddress:        struct{ primemessages.Address }{pickupAddress},
		DestinationAddress:   struct{ primemessages.Address }{destinationAddress},
	}

	result := MTOShipmentModelFromCreate(mtoShipment)

	suite.NotNil(result)
	suite.Equal(mtoShipment.MoveTaskOrderID.String(), result.MoveTaskOrderID.String())
	suite.Equal(*mtoShipment.CustomerRemarks, *result.CustomerRemarks)

	suite.NotNil(result.PrimeEstimatedWeight)
	suite.Equal(unit.Pound(primeEstimatedWeight), *result.PrimeEstimatedWeight)
	suite.NotNil(result.PrimeEstimatedWeightRecordedDate)
	suite.WithinDuration(time.Now(), *result.PrimeEstimatedWeightRecordedDate, time.Second)

	suite.NotNil(result.PickupAddress)
	suite.Equal("123 Main St", result.PickupAddress.StreetAddress1)
	suite.NotNil(result.DestinationAddress)
	suite.Equal("456 Main St", result.DestinationAddress.StreetAddress1)
}

func (suite *PayloadsSuite) TestVLocationModel() {
	city := "LOS ANGELES"
	state := "CA"
	postalCode := "90210"
	county := "LOS ANGELES"
	usPostRegionCityId := uuid.Must(uuid.NewV4())

	vLocation := &primemessages.VLocation{
		City:                 city,
		State:                state,
		PostalCode:           postalCode,
		County:               &county,
		UsPostRegionCitiesID: strfmt.UUID(usPostRegionCityId.String()),
	}

	payload := VLocationModel(vLocation)

	suite.IsType(payload, &models.VLocation{})
	suite.Equal(usPostRegionCityId.String(), payload.UsPostRegionCitiesID.String(), "Expected UsPostRegionCitiesID to match")
	suite.Equal(city, payload.CityName, "Expected City to match")
	suite.Equal(state, payload.StateName, "Expected State to match")
	suite.Equal(postalCode, payload.UsprZipID, "Expected PostalCode to match")
	suite.Equal(county, payload.UsprcCountyNm, "Expected County to match")
}
