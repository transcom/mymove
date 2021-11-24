package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importRETaskOrderFees(appCtx appcontext.AppContext) error {
	//tab 4a) Mgmt., Coun., Trans. Prices
	var shipmentManagementServices []models.StageShipmentManagementServicesPrice
	err := appCtx.DB().All(&shipmentManagementServices)
	if err != nil {
		return fmt.Errorf("could not read staged shipment management service prices: %w", err)
	}

	//loop through the shipment management service data, pull data for management services and save in db
	for _, stageShipmentManagementServicePrice := range shipmentManagementServices {
		shipmentManagementService, foundService := gre.serviceToIDMap[models.ReServiceCodeMS]
		if !foundService {
			return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeMS)
		}

		contractYear := stageShipmentManagementServicePrice.ContractYear
		contractYearID, found := gre.contractYearToIDMap[contractYear]
		if !found {
			return fmt.Errorf("could not find contract year %s in map", stageShipmentManagementServicePrice.ContractYear)
		}

		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageShipmentManagementServicePrice.PricePerTaskOrder)
		if err != nil {
			return fmt.Errorf("could not process shipment management service price [%s]: %w", stageShipmentManagementServicePrice.PricePerTaskOrder, err)
		}

		taskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYearID,
			ServiceID:      shipmentManagementService,
			PriceCents:     unit.Cents(perUnitCentsService),
		}

		verrs, dbErr := appCtx.DB().ValidateAndSave(&taskOrderFee)
		if dbErr != nil {
			return fmt.Errorf("error saving ReTaskOrderFees: %+v with error: %w", taskOrderFee, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReTaskOrderFees: %+v with validation errors: %w", taskOrderFee, verrs)
		}
	}

	var shipmentCounselingServices []models.StageCounselingServicesPrice
	err = appCtx.DB().All(&shipmentCounselingServices)
	if err != nil {
		return fmt.Errorf("could not read staged shipment counseling service prices: %w", err)
	}

	//loop through the shipment management service data, pull data for counseling services and save in db
	for _, stageShipmentCounselingServicePrice := range shipmentCounselingServices {
		shipmentCounselingService, foundService := gre.serviceToIDMap[models.ReServiceCodeCS]
		if !foundService {
			return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeCS)
		}

		contractYear := stageShipmentCounselingServicePrice.ContractYear
		contractYearID, found := gre.contractYearToIDMap[contractYear]
		if !found {
			return fmt.Errorf("could not find contract year %s in map", stageShipmentCounselingServicePrice.ContractYear)
		}

		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageShipmentCounselingServicePrice.PricePerTaskOrder)
		if err != nil {
			return fmt.Errorf("could not process shipment counseling service price [%s]: %w", stageShipmentCounselingServicePrice.PricePerTaskOrder, err)
		}

		taskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYearID,
			ServiceID:      shipmentCounselingService,
			PriceCents:     unit.Cents(perUnitCentsService),
		}

		verrs, dbErr := appCtx.DB().ValidateAndSave(&taskOrderFee)
		if dbErr != nil {
			return fmt.Errorf("error saving ReTaskOrderFees: %+v with error: %w", taskOrderFee, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReTaskOrderFees: %+v with validation errors: %w", taskOrderFee, verrs)
		}
	}

	return nil
}
