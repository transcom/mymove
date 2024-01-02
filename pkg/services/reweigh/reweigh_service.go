package reweigh

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/unit"
)

func getDivertedShipmentIDsAndLowestWeight(appCtx appcontext.AppContext, reweigh *models.Reweigh) ([]uuid.UUID, *unit.Pound, error) {
	var shipmentIDs []uuid.UUID
	var lowestWeight *unit.Pound
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	reweighFetcher := NewReweighFetcher()

	// Find out if this shipment is part of a diversion chain, and if so gather all associated shipments
	associatedDivertedShipments, err := shipmentFetcher.GetDiversionChain(appCtx, reweigh.ShipmentID)
	if err != nil {
		return nil, nil, err
	}

	// Create an array of UUIDs, dereferenced from associatedDivertedShipments
	for _, shipment := range *associatedDivertedShipments {
		shipmentIDs = append(shipmentIDs, shipment.ID)
	}
	reweighs, err := reweighFetcher.ListReweighsByShipmentIDs(appCtx, shipmentIDs)
	if err != nil {
		return nil, nil, err
	}

	// Find the lowest reweigh
	for _, tempReweighVar := range reweighs {
		// Check for valid reweigh weight, if we have set a lowest weight yet, and if so then if the new weight is lower than our lowest
		if tempReweighVar.Weight != nil && (lowestWeight == nil || *tempReweighVar.Weight < *lowestWeight) {
			// If valid and lower than our lowest, we set the new weight here
			lowestWeight = tempReweighVar.Weight
			//shipmentIdWithLowestWeight = shipmentID
		}
	}
	return shipmentIDs, lowestWeight, nil
}
