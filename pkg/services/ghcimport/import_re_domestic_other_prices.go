package ghcimport

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// DomOtherPriceToInsert is the domestic other price to insert
type DomOtherPriceToInsert struct {
	model   models.ReDomesticOtherPrice
	message string
}

func importPackUnpackPrices(appCtx appcontext.AppContext, serviceToIDMap map[models.ReServiceCode]uuid.UUID, contractID uuid.UUID) ([]DomOtherPriceToInsert, error) {
	var stagePackPrices []models.StageDomesticOtherPackPrice
	var modelsToSave []DomOtherPriceToInsert

	if err := appCtx.DB().All(&stagePackPrices); err != nil {
		return nil, fmt.Errorf("error looking up StageDomesticOtherPackPrice data: %w", err)
	}

	// DPK Dom. Packing
	packServiceID, found := serviceToIDMap[models.ReServiceCodeDPK]
	if !found {
		return nil, fmt.Errorf("missing service [%s] in map of services", models.ReServiceCodeDPK)
	}

	// DUPK Dom. Unpacking
	unpackServiceID, found := serviceToIDMap[models.ReServiceCodeDUPK]
	if !found {
		return nil, fmt.Errorf("missing service [%s] in map of services", models.ReServiceCodeDUPK)
	}

	for _, stagePackPrice := range stagePackPrices {
		peakCents, err := priceToCents(stagePackPrice.PeakPricePerCwt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price for service code %s: %+v error: %w", models.ReServiceCodeDPK, stagePackPrice.PeakPricePerCwt, err)
		}

		nonPeakCents, err := priceToCents(stagePackPrice.NonPeakPricePerCwt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price for service code %s: %+v error: %w", models.ReServiceCodeDUPK, stagePackPrice.NonPeakPricePerCwt, err)
		}

		servicesSchedule, err := stringToInteger(stagePackPrice.ServicesSchedule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ServicesSchedule for pack/unpack: %+v error: %w", stagePackPrice.ServicesSchedule, err)
		}

		packNonPeakPriceModel := models.ReDomesticOtherPrice{
			ContractID:   contractID,
			Schedule:     servicesSchedule,
			IsPeakPeriod: false,
			PriceCents:   unit.Cents(nonPeakCents),
		}
		packPeakPriceModel := models.ReDomesticOtherPrice{
			ContractID:   contractID,
			Schedule:     servicesSchedule,
			IsPeakPeriod: true,
			PriceCents:   unit.Cents(peakCents),
		}

		if stagePackPrice.ServiceProvided == "Packing (per cwt)" {
			packNonPeakPriceModel.ServiceID = packServiceID
			packPeakPriceModel.ServiceID = packServiceID
		} else if stagePackPrice.ServiceProvided == "Unpack (per cwt)" {
			packNonPeakPriceModel.ServiceID = unpackServiceID
			packPeakPriceModel.ServiceID = unpackServiceID
		} else {
			return nil, fmt.Errorf("failed to import pack/unpack prices receieved unexpected ServiceProvided: %s in %+v", stagePackPrice.ServiceProvided, stagePackPrice)
		}

		modelsToSave = append(modelsToSave, DomOtherPriceToInsert{model: packNonPeakPriceModel, message: "Non-Peak Pack/Unpack"})
		modelsToSave = append(modelsToSave, DomOtherPriceToInsert{model: packPeakPriceModel, message: "Peak Pack/Unpack"})
	}

	return modelsToSave, nil
}

func importSitPrices(appCtx appcontext.AppContext, serviceToIDMap map[models.ReServiceCode]uuid.UUID, contractID uuid.UUID) ([]DomOtherPriceToInsert, error) {
	var stageSitPrices []models.StageDomesticOtherSitPrice
	var modelsToSave []DomOtherPriceToInsert

	if err := appCtx.DB().All(&stageSitPrices); err != nil {
		return nil, fmt.Errorf("error looking up StageDomesticOtherSitPrice data: %w", err)
	}

	// DOPSIT Dom. Origin SIT Pickup
	originSitPickupID, found := serviceToIDMap[models.ReServiceCodeDOPSIT]
	if !found {
		return nil, fmt.Errorf("missing service [%s] in map of services", models.ReServiceCodeDOPSIT)
	}

	// DDDSIT Dom. Destination SIT Delivery
	destSitDeliveryID, found := serviceToIDMap[models.ReServiceCodeDDDSIT]
	if !found {
		return nil, fmt.Errorf("missing service [%s] in map of services", models.ReServiceCodeDDDSIT)
	}

	for _, stageSitPrice := range stageSitPrices {
		peakCents, err := priceToCents(stageSitPrice.PeakPricePerCwt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price for service code %s: %+v error: %w", models.ReServiceCodeDOPSIT, stageSitPrice.PeakPricePerCwt, err)
		}

		nonPeakCents, err := priceToCents(stageSitPrice.NonPeakPricePerCwt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price for service code %s: %+v error: %w", models.ReServiceCodeDDDSIT, stageSitPrice.NonPeakPricePerCwt, err)
		}

		schedule, err := stringToInteger(stageSitPrice.SITPickupDeliverySchedule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SITPickupDeliverySchedule: %+v error: %w", stageSitPrice.SITPickupDeliverySchedule, err)
		}

		modelsToSave = append(
			modelsToSave,
			DomOtherPriceToInsert{model: models.ReDomesticOtherPrice{
				ContractID:   contractID,
				ServiceID:    originSitPickupID,
				Schedule:     schedule,
				IsPeakPeriod: false,
				PriceCents:   unit.Cents(nonPeakCents),
			}, message: "SIT Non Peak Pickup"})
		modelsToSave = append(
			modelsToSave,
			DomOtherPriceToInsert{model: models.ReDomesticOtherPrice{
				ContractID:   contractID,
				ServiceID:    destSitDeliveryID,
				Schedule:     schedule,
				IsPeakPeriod: false,
				PriceCents:   unit.Cents(nonPeakCents),
			}, message: "SIT Non Peak Delivery"})
		modelsToSave = append(
			modelsToSave,
			DomOtherPriceToInsert{model: models.ReDomesticOtherPrice{
				ContractID:   contractID,
				ServiceID:    originSitPickupID,
				Schedule:     schedule,
				IsPeakPeriod: true,
				PriceCents:   unit.Cents(peakCents),
			}, message: "SIT Peak Pickup"})
		modelsToSave = append(
			modelsToSave,
			DomOtherPriceToInsert{model: models.ReDomesticOtherPrice{
				ContractID:   contractID,
				ServiceID:    destSitDeliveryID,
				Schedule:     schedule,
				IsPeakPeriod: true,
				PriceCents:   unit.Cents(peakCents),
			}, message: "SIT Peak Delivery"})
	}

	return modelsToSave, nil
}

func saveModel(appCtx appcontext.AppContext, message string, model *models.ReDomesticOtherPrice) error {
	verrs, err := appCtx.DB().ValidateAndSave(model)
	if verrs.HasAny() {
		return fmt.Errorf("error saving ReDomesticOtherPrice %s: %+v with validation errors: %w", message, model, verrs)
	}
	if err != nil {
		return fmt.Errorf("error saving ReDomesticOtherPrice %s: %+v with error: %w", message, model, err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) importREDomesticOtherPrices(appCtx appcontext.AppContext) error {

	var modelsToSavePack []DomOtherPriceToInsert
	var modelsToSaveSit []DomOtherPriceToInsert

	var err error
	modelsToSavePack, err = importPackUnpackPrices(appCtx, gre.serviceToIDMap, gre.ContractID)
	if err != nil {
		return err
	}

	modelsToSaveSit, err = importSitPrices(appCtx, gre.serviceToIDMap, gre.ContractID)
	if err != nil {
		return err
	}

	modelsToSave := append(modelsToSavePack, modelsToSaveSit...)
	for _, modelToSave := range modelsToSave {
		//nolint:gosec
		if err := saveModel(appCtx, modelToSave.message, &modelToSave.model); err != nil {
			return err
		}
	}

	return nil
}
