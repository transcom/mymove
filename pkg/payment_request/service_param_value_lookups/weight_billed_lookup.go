package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightBilledLookup does lookup on weight billed
type WeightBilledLookup struct {
	MTOShipment models.MTOShipment
}

func (r WeightBilledLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	var estimatedWeight *unit.Pound
	var originalWeight *unit.Pound
	var value string

	switch keyData.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDOSHUT,
		models.ReServiceCodeDDSHUT,
		models.ReServiceCodeIOSHUT,
		models.ReServiceCodeIDSHUT:
		estimatedWeight = keyData.MTOServiceItem.EstimatedWeight

		// Check both the service item weight and if it can't find that then check the shipment's weight
		if keyData.MTOServiceItem.ActualWeight == nil {
			originalWeight = r.MTOShipment.PrimeActualWeight
			if originalWeight == nil {
				return "", fmt.Errorf("could not find actual weight for MTOServiceItemID [%s] or for MTOShipmentID [%s]", keyData.MTOServiceItem.ID, r.MTOShipment.ID)
			}
		} else {
			originalWeight = keyData.MTOServiceItem.ActualWeight
		}

		if estimatedWeight != nil {
			estimatedWeightCap := math.Round(float64(*estimatedWeight) * 1.10)
			if float64(*originalWeight) > estimatedWeightCap {
				value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(estimatedWeightCap))
			} else {
				value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(*originalWeight))
			}
		} else {
			value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(*originalWeight))
		}
		return value, nil
	case models.ReServiceCodeDDSFSC,
		models.ReServiceCodeDOSFSC,
		models.ReServiceCodeFSC:

		var weightBilled string

		// Check if a value is in WeightBilled
		query := `select psip.value
			from payment_service_item_params psip
				join payment_service_items psi
					on psip.payment_service_item_id = psi.id
				join mto_service_items msi
					on msi.id = psi.mto_service_item_id
				join re_services rs
					on rs.id = msi.re_service_id
				join payment_requests pr
					on psi.payment_request_id = pr.id
				join service_item_param_keys sipk
					on sipk.id = psip.service_item_param_key_id
			where sipk.key = 'WeightBilled' and psi.payment_request_id = $1 and rs.code = $2`

		err := appCtx.DB().RawQuery(query, keyData.PaymentRequestID, keyData.MTOServiceItem.ReService.Code).First(&weightBilled)

		if err != nil && err != sql.ErrNoRows {
			return "", err
		}

		if len(weightBilled) > 0 {
			return weightBilled, nil
		}
		estimatedWeight = r.MTOShipment.PrimeEstimatedWeight

		originalWeight = r.MTOShipment.PrimeActualWeight

		if originalWeight == nil {
			// TODO: Do we need a different error -- is this a "normal" scenario?
			return "", fmt.Errorf("could not find actual weight for MTOServiceItemID [%s]", keyData.MTOServiceItem.ID)
		}

		// Make sure the reweigh (if any) is loaded since that's expected by the calculate shipment billable weight service.
		err = appCtx.DB().Load(&r.MTOShipment, "Reweigh")
		if err != nil {
			return "", err
		}

		// Check if the shipment has a reweigh
		noReweigh := r.MTOShipment.Reweigh.ID.IsNil()

		// If shipment has an estimatedWeight and is not a Reweigh. Else If the Shipment does have an estimatedWeight and has a Reweigh then call calculateMinimumBillableWeight(). Else return value from applyMinimum.
		if estimatedWeight != nil && noReweigh {
			estimatedWeightCap := math.Round(float64(*estimatedWeight) * 1.10)
			if float64(*originalWeight) > estimatedWeightCap {
				value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(estimatedWeightCap))
			} else {
				value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(*originalWeight))
			}
		} else if estimatedWeight != nil && !noReweigh {
			value, err = calculateMinimumBillableWeight(appCtx, r.MTOShipment, keyData)
			if err != nil {
				return "", err
			}
		} else {
			value = applyMinimum(keyData.MTOServiceItem.ReService.Code, r.MTOShipment.ShipmentType, int(*originalWeight))
		}
		return value, nil

	default:
		// Shipments that are a diversion must utilize the lowest weight that can be found
		// in their "diverted shipment chain". Diverted shipments are tied together by "divertedFromShipmentId"s after the implementation
		// of the createMTOShipment V2.
		// Only diverted shipments created utilizing the `prime/v2/createMTOShipment` endpoint will be able to get identified as "chains"
		// Shipments created with the V1 endpoint will not be referenced by a divertedFromShipmentId and hence will just use the lowest weight assigned
		// to the shipment as is.
		if r.MTOShipment.Diversion {
			// Identify diversion chain for weight calculations.
			mtoShipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
			diversionChain, err := mtoShipmentFetcher.GetDiversionChain(appCtx, r.MTOShipment.ID)
			if err != nil {
				return "", err
			}

			// Initialize to maximum int value of 32. This is done to replicate `Number.MAX_SAFE_INTEGER` and comparing down like it was
			// done on the frontend with JavaScript
			var lowestWeight = math.MaxInt32
			var shipmentWithLowestWeight *models.MTOShipment
			for _, divertedShipment := range *diversionChain {
				if divertedShipment.PrimeActualWeight == nil {
					// ! Payments should never be created for a diverted shipment that has a nil PrimeActualWeight inside the chain
					// ! Unless it is a partial payment
					return "", fmt.Errorf("all shipments in the diversion chain must have a `PrimeActualWeight` field if you are not creating a partial payment. Please update the shipment prior to creating the payment request. Shipment ID: %s", divertedShipment.ID)
				}
				// Calculate the billable weight for each shipment in the chain
				billableWeight, err := calculateMinimumBillableWeight(appCtx, divertedShipment, keyData)
				if err != nil {
					return "", err
				}
				// Convert ascii to int so we can compare to our current lowest weight
				billableWeightInt, err := strconv.Atoi(billableWeight)
				if err != nil {
					return "", err
				}
				// Update the lowest weight if the current shipment's weight is lower
				if billableWeightInt < lowestWeight {
					lowestWeight = billableWeightInt
					newDivertedShipmentMemoryRef := divertedShipment
					shipmentWithLowestWeight = &newDivertedShipmentMemoryRef
				}
			}
			if shipmentWithLowestWeight == nil || lowestWeight == math.MaxInt32 {
				return "", fmt.Errorf("unexpected error when calculating the minimum billable weight for a chain of diverted shipments, a lowest weight could not be identified")
			}

			return calculateMinimumBillableWeight(appCtx, *shipmentWithLowestWeight, keyData)
		}

		// If not a diversion, proceed with calculations normally
		return calculateMinimumBillableWeight(appCtx, r.MTOShipment, keyData)
	}
}

func calculateMinimumBillableWeight(appCtx appcontext.AppContext, shipment models.MTOShipment, keyData *ServiceItemParamKeyData) (string, error) {
	originalWeight := shipment.PrimeActualWeight
	// Make sure there's an actual weight since that's nullable but required for pricing
	if originalWeight == nil {
		// TODO: Do we need a different error -- is this a "normal" scenario?
		return "", fmt.Errorf("could not find actual weight for MTOShipmentID [%s]", shipment.ID)
	}

	// Make sure the reweigh (if any) is loaded since that's expected by the calculate shipment billable weight service.
	err := appCtx.DB().Load(&shipment, "Reweigh")
	if err != nil {
		return "", err
	}

	calculator := mtoshipment.NewShipmentBillableWeightCalculator()
	billableWeightInputs := calculator.CalculateShipmentBillableWeight(&shipment)
	if billableWeightInputs.CalculatedBillableWeight == nil {
		return "", fmt.Errorf("got a nil calculated billable weight from service for MTOShipmentID [%s]", shipment.ID)
	}

	return applyMinimum(keyData.MTOServiceItem.ReService.Code, shipment.ShipmentType, int(*billableWeightInputs.CalculatedBillableWeight)), nil

}

// Looks at code and applies minimum if necessary, otherwise returns weight passed in
func applyMinimum(code models.ReServiceCode, shipmentType models.MTOShipmentType, weight int) string {
	result := weight
	switch shipmentType {
	case models.MTOShipmentTypeUnaccompaniedBaggage:
		switch code {
		case models.ReServiceCodeIOSHUT,
			models.ReServiceCodeIDSHUT:
			if weight < 300 {
				result = 300
			}
		}
	case models.MTOShipmentTypePPM:
		result = weight
	default:
		switch code {
		case models.ReServiceCodeDLH,
			models.ReServiceCodeDSH,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDOFSIT,
			models.ReServiceCodeDDFSIT,
			models.ReServiceCodeDOASIT,
			models.ReServiceCodeDOSFSC,
			models.ReServiceCodeDDASIT,
			models.ReServiceCodeDOPSIT,
			models.ReServiceCodeDDDSIT,
			models.ReServiceCodeDDSFSC,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDNPK,
			models.ReServiceCodeDUPK,
			models.ReServiceCodeDOSHUT,
			models.ReServiceCodeDDSHUT,
			models.ReServiceCodeISLH,
			models.ReServiceCodeIHPK,
			models.ReServiceCodeIHUPK,
			models.ReServiceCodeIOFSIT,
			models.ReServiceCodeIDFSIT,
			models.ReServiceCodeIOASIT,
			models.ReServiceCodeIDASIT,
			models.ReServiceCodeIOPSIT,
			models.ReServiceCodeIDDSIT,
			models.ReServiceCodeIOSHUT,
			models.ReServiceCodeIDSHUT,
			models.ReServiceCodeFSC:
			if weight < 500 {
				result = 500
			}
		case models.ReServiceCodeUBP,
			models.ReServiceCodeIUBPK,
			models.ReServiceCodeIUBUPK:
			if weight < 300 {
				result = 300
			}
		}
	}
	return fmt.Sprintf("%d", result)
}
