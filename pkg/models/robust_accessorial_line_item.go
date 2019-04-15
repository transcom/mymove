package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

// BaseShipmentLineItemParams holds the basic parameters for a ShipmentLineItem
type BaseShipmentLineItemParams struct {
	Tariff400ngItemID   uuid.UUID
	Tariff400ngItemCode string
	Quantity1           *unit.BaseQuantity
	Quantity2           *unit.BaseQuantity
	Location            string
	Notes               *string
}

// AdditionalShipmentLineItemParams holds any additional parameters for a ShipmentLineItem
type AdditionalShipmentLineItemParams struct {
	Description         *string
	ItemDimensions      *AdditionalLineItemDimensions
	CrateDimensions     *AdditionalLineItemDimensions
	Reason              *string
	EstimateAmountCents *unit.Cents
	ActualAmountCents   *unit.Cents
	Date                *time.Time
	Time                *string
	Address             *Address
}

// AdditionalLineItemDimensions holds the length, width and height that will be converted to inches
type AdditionalLineItemDimensions struct {
	Length unit.ThousandthInches
	Width  unit.ThousandthInches
	Height unit.ThousandthInches
}

// upsertItemCodeDependency applies specific validation, creates or updates additional objects/fields for item codes.
// Mutates the shipmentLineItem passed in.
func upsertItemCodeDependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	itemCode := baseParams.Tariff400ngItemCode

	// Backwards compatible with "Old school" base quantity field
	if is105Item(itemCode, additionalParams) {
		return upsertItemCode105Dependency(db, baseParams, additionalParams, shipmentLineItem)
	} else if is35AItem(itemCode, additionalParams) {
		return upsertItemCode35ADependency(db, baseParams, additionalParams, shipmentLineItem)
	} else if is226AItem(itemCode, additionalParams) {
		return upsertItemCode226ADependency(db, baseParams, additionalParams, shipmentLineItem)
	} else if is125Item(itemCode, additionalParams) {
		return upsertItemCode125Dependency(db, baseParams, additionalParams, shipmentLineItem)
	}

	return upsertDefaultDependency(db, baseParams, additionalParams, shipmentLineItem)
}

// createShipmentLineItemDimensions creates new item and crate dimensions for shipment line item
func createShipmentLineItemDimensions(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	// save dimensions to shipmentLineItem
	shipmentLineItem.ItemDimensions = ShipmentLineItemDimensions{
		Length: unit.ThousandthInches(additionalParams.ItemDimensions.Length),
		Width:  unit.ThousandthInches(additionalParams.ItemDimensions.Width),
		Height: unit.ThousandthInches(additionalParams.ItemDimensions.Height),
	}
	verrs, err := db.ValidateAndCreate(&shipmentLineItem.ItemDimensions)
	if verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error creating item dimensions for shipment line item")
		return responseVErrors, responseError
	}

	shipmentLineItem.CrateDimensions = ShipmentLineItemDimensions{
		Length: unit.ThousandthInches(additionalParams.CrateDimensions.Length),
		Width:  unit.ThousandthInches(additionalParams.CrateDimensions.Width),
		Height: unit.ThousandthInches(additionalParams.CrateDimensions.Height),
	}
	verrs, err = db.ValidateAndCreate(&shipmentLineItem.CrateDimensions)
	if verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error creating crate dimensions for shipment line item")
		return responseVErrors, responseError
	}

	shipmentLineItem.ItemDimensionsID = &shipmentLineItem.ItemDimensions.ID
	shipmentLineItem.CrateDimensionsID = &shipmentLineItem.CrateDimensions.ID

	return responseVErrors, responseError
}

// updateShipmentLineItemDimensions updates existing shipment line item dimensions
func updateShipmentLineItemDimensions(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	// save dimensions to shipmentLineItem
	shipmentLineItem.ItemDimensions.Length = unit.ThousandthInches(additionalParams.ItemDimensions.Length)
	shipmentLineItem.ItemDimensions.Width = unit.ThousandthInches(additionalParams.ItemDimensions.Width)
	shipmentLineItem.ItemDimensions.Height = unit.ThousandthInches(additionalParams.ItemDimensions.Height)

	verrs, err := db.ValidateAndUpdate(&shipmentLineItem.ItemDimensions)
	if verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error updating item dimensions for shipment line item")
		return responseVErrors, responseError
	}

	shipmentLineItem.CrateDimensions.Length = unit.ThousandthInches(additionalParams.CrateDimensions.Length)
	shipmentLineItem.CrateDimensions.Width = unit.ThousandthInches(additionalParams.CrateDimensions.Width)
	shipmentLineItem.CrateDimensions.Height = unit.ThousandthInches(additionalParams.CrateDimensions.Height)

	verrs, err = db.ValidateAndUpdate(&shipmentLineItem.CrateDimensions)
	if verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error updating crate dimensions for shipment line item")
		return responseVErrors, responseError
	}

	return responseVErrors, responseError
}

// is105Item determines whether the shipment line item is a new (robust) 105B/E item.
func is105Item(itemCode string, additionalParams *AdditionalShipmentLineItemParams) bool {
	hasDimension := additionalParams.ItemDimensions != nil || additionalParams.CrateDimensions != nil
	if (itemCode == "105B" || itemCode == "105E") && hasDimension {
		return true
	}
	return false
}

// upsertItemCode105Dependency specifically upserts item code 105B/E for shipmentLineItem passed in
func upsertItemCode105Dependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	//Additional validation check if item and crate dimensions exist
	if additionalParams.ItemDimensions == nil || additionalParams.CrateDimensions == nil {
		responseError = errors.New("Must have both item and crate dimensions params")
		return responseVErrors, responseError
	}

	if shipmentLineItem.ItemDimensions.ID == uuid.Nil || shipmentLineItem.CrateDimensions.ID == uuid.Nil {
		verrs, err := createShipmentLineItemDimensions(db, baseParams, additionalParams, shipmentLineItem)
		if verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating shipment line item dimensions")
			return responseVErrors, responseError
		}
	} else {
		verrs, err := updateShipmentLineItemDimensions(db, baseParams, additionalParams, shipmentLineItem)
		if verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error updating shipment line item dimensions")
			return responseVErrors, responseError
		}
	}

	// get crate volume in cubic feet
	crateVolume, err := unit.DimensionToCubicFeet(additionalParams.CrateDimensions.Length, additionalParams.CrateDimensions.Width, additionalParams.CrateDimensions.Height)
	if err != nil {
		return nil, errors.Wrap(err, "Dimension units must be greater than 0")
	}

	// format value to base quantity i.e. times 10,000
	formattedQuantity1 := unit.BaseQuantityFromFloat(float32(crateVolume))
	shipmentLineItem.Quantity1 = formattedQuantity1
	shipmentLineItem.Description = additionalParams.Description

	return responseVErrors, responseError
}

// is35AItem determines whether the shipment line item is a new (robust) 35A item.
func is35AItem(itemCode string, additionalParams *AdditionalShipmentLineItemParams) bool {
	isRobustItem := additionalParams.Reason != nil || additionalParams.EstimateAmountCents != nil || additionalParams.ActualAmountCents != nil
	if itemCode == "35A" && isRobustItem {
		return true
	}
	return false
}

// upsertItemCode35ADependency specifically upserts item code 35A for shipmentLineItem passed in
func upsertItemCode35ADependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	if shipmentLineItem.ID != uuid.Nil && shipmentLineItem.Status == ShipmentLineItemStatusAPPROVED {
		// Line item exists
		// Only update the ActualAmountCents
		shipmentLineItem.ActualAmountCents = additionalParams.ActualAmountCents
	} else if additionalParams.Description == nil || additionalParams.Reason == nil || additionalParams.EstimateAmountCents == nil {
		// Required to create 35A line item
		// Description, Reason and EstimateAmounCents
		responseError = errors.New("Must have Description, Reason and EstimateAmountCents params")
		return responseVErrors, responseError
	} else {
		shipmentLineItem.Description = additionalParams.Description
		shipmentLineItem.Reason = additionalParams.Reason
		shipmentLineItem.EstimateAmountCents = additionalParams.EstimateAmountCents
		shipmentLineItem.ActualAmountCents = additionalParams.ActualAmountCents
	}

	if shipmentLineItem.ActualAmountCents != nil {
		if *shipmentLineItem.ActualAmountCents <= *shipmentLineItem.EstimateAmountCents {
			shipmentLineItem.Quantity1 = unit.BaseQuantityFromCents(*shipmentLineItem.ActualAmountCents)
		} else {
			shipmentLineItem.Quantity1 = unit.BaseQuantityFromCents(*shipmentLineItem.EstimateAmountCents)
		}
	} else {
		// If ActualAmountCents is unset, set base quantity to 0.
		quantity1 := unit.BaseQuantityFromInt(0)
		shipmentLineItem.Quantity1 = quantity1
	}
	return responseVErrors, responseError
}

// is226AItem determines whether the shipment line item is a new (robust) 226A item.
func is226AItem(itemCode string, additionalParams *AdditionalShipmentLineItemParams) bool {
	isRobustItem := additionalParams.Description != nil || additionalParams.Reason != nil || additionalParams.ActualAmountCents != nil
	if itemCode == "226A" && isRobustItem {
		return true
	}
	return false
}

// upsertItemCode226ADependency specifically upserts item code 226A for shipmentLineItem passed in
func upsertItemCode226ADependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	// Required to create 226A line item
	// Description, Reason and EstimateAmounCents
	if additionalParams.Description == nil || additionalParams.Reason == nil || additionalParams.ActualAmountCents == nil {
		responseError = errors.New("Must have Description, Reason and ActualAmountCents params")
		return responseVErrors, responseError
	}

	shipmentLineItem.Description = additionalParams.Description
	shipmentLineItem.Reason = additionalParams.Reason
	shipmentLineItem.ActualAmountCents = additionalParams.ActualAmountCents
	shipmentLineItem.Quantity1 = unit.BaseQuantityFromCents(*shipmentLineItem.ActualAmountCents)

	return responseVErrors, responseError
}

// is125Item determines whether the shipment line item is a new (robust) 125 item.
func is125Item(itemCode string, additionalParams *AdditionalShipmentLineItemParams) bool {
	isRobustItem := additionalParams.Reason != nil || additionalParams.Date != nil || additionalParams.Time != nil || additionalParams.Address != nil
	if strings.HasPrefix(itemCode, "125") && isRobustItem {
		return true
	}
	return false
}

// upsertItemCode125Dependency specifically upserts item code 125 for shipmentLineItem passed in
func upsertItemCode125Dependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	// Required to create 125A/B/C/D line item
	if additionalParams.Reason == nil || additionalParams.Date == nil || additionalParams.Address == nil {
		responseError = errors.New("Must have Reason, Date and Address params")
		return responseVErrors, responseError
	}

	shipmentLineItem.Address = *additionalParams.Address
	// Create address if it doesn't exist
	if shipmentLineItem.AddressID == nil {
		verrs, err := db.ValidateAndCreate(&shipmentLineItem.Address)
		if verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating shipment line item address")
			return responseVErrors, responseError
		}

		shipmentLineItem.AddressID = &shipmentLineItem.Address.ID
	} else {
		//otherwise, update the address
		shipmentLineItem.Address.ID = *shipmentLineItem.AddressID
		verrs, err := db.ValidateAndUpdate(&shipmentLineItem.Address)
		if verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error updating shipment line item address")
			return responseVErrors, responseError
		}
	}

	shipmentLineItem.Reason = additionalParams.Reason
	shipmentLineItem.Date = additionalParams.Date
	shipmentLineItem.Time = additionalParams.Time
	shipmentLineItem.Quantity1 = unit.BaseQuantityFromInt(1) // flat rate, set to base quantity 1

	return responseVErrors, responseError
}

// upsertDefaultDependency upserts standard shipmentLineItem passed in
func upsertDefaultDependency(db *pop.Connection, baseParams *BaseShipmentLineItemParams, additionalParams *AdditionalShipmentLineItemParams, shipmentLineItem *ShipmentLineItem) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()
	if baseParams.Quantity1 == nil {
		// General pre-approval request
		// Check if base quantity is filled out
		responseError = errors.New("Quantity1 required for tariff400ngItemCode: " + baseParams.Tariff400ngItemCode)
	} else {
		// Good to fill out quantity1
		shipmentLineItem.Quantity1 = *baseParams.Quantity1
	}

	return responseVErrors, responseError
}
