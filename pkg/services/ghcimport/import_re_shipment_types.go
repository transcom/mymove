package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREShipmentTypes(dbTx *pop.Connection) error {
	// populate shipmentTypeToIDMap
	var domesticIntlAddlPrices []models.StageDomesticInternationalAdditionalPrices
	err := dbTx.All(&domesticIntlAddlPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic internatinoal additional prices: %w", err)
	}

	gre.shipmentTypeToIDMap = make(map[string]uuid.UUID)
	shipmentCodePositionInSlice := 0

	//loop through the domestic international additional prices data and pull shipment types
	for _, stageDomesticIntlAddlPrices := range domesticIntlAddlPrices {
		//slice of shipment codes for each shipment type
		shipmentCodes := []string{"DMHF", "DBTF", "DBHF", "IBTF", "IBHF", "DNPKF", "INPKF"}

		shipmentType := models.ReShipmentType{
			Code: shipmentCodes[shipmentCodePositionInSlice],
			Name: stageDomesticIntlAddlPrices.ShipmentType,
		}
		shipmentCodePositionInSlice++

		verrs, dbErr := dbTx.ValidateAndSave(&shipmentType)
		if dbErr != nil {
			return fmt.Errorf("error saving ReShipmentTypes: %+v with error: %w", shipmentType, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReShipmentTypes: %+v with validation errors: %w", shipmentType, verrs)
		}

		//add to map
		gre.shipmentTypeToIDMap[shipmentType.Name] = shipmentType.ID
	}

	return nil
}
