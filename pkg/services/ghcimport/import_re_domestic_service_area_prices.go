package ghcimport

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func correctServiceAreaNumber(serviceAreaNumber string) string {
	num, err := strconv.ParseFloat(serviceAreaNumber, 64)
	if err != nil {
		return serviceAreaNumber
	}
	return fmt.Sprintf("%03d", int(num))
}

// borrowed from https://play.golang.org/p/AwXg4_jp8lo
func parseCents(s string) (int64, error) {
	n := strings.SplitN(s, ".", 3)
	if len(n) != 2 || len(n[1]) != 2 {
		err := fmt.Errorf("format error: %s", s)
		return 0, err
	}
	d, err := strconv.ParseInt(n[0], 10, 56)
	if err != nil {
		return 0, err
	}
	c, err := strconv.ParseUint(n[1], 10, 8)
	if err != nil {
		return 0, err
	}
	if d < 0 {
		c = -c
	}
	return d*100 + int64(c), nil
}

func stringDollarsToCents(dollars string) (unit.Cents, error) {
	priceCents, err := parseCents(strings.Replace(dollars, "$", "", 1))
	if err != nil {
		return 0, err
	}
	return unit.Cents(int(priceCents)), nil
}

func (gre *GHCRateEngineImporter) importREDomesticServiceAreaPrices(db *pop.Connection) error {
	var stageDomPricingModels []models.StageDomesticServiceAreaPrice

	if err := db.All(&stageDomPricingModels); err != nil {
		return fmt.Errorf("Error looking up StageDomesticServiceAreaPrice data: %w", err)
	}

	for _, stageDomPricingModel := range stageDomPricingModels {
		var domPricingModels models.ReDomesticServiceAreaPrices

		var isPeakPeriod bool
		if stageDomPricingModel.Season == "Peak" {
			isPeakPeriod = true
		}

		//ServiceAreaNumber                     string `db:"service_area_number" csv:"service_area_number"`
		servicesSchedule, ssErr := strconv.Atoi(stageDomPricingModel.ServicesSchedule)
		if ssErr != nil {
			return fmt.Errorf("Failed to parse ServicesSchedule for %+v: %w", stageDomPricingModel, ssErr)
		}
		sITPDSchedule, spErr := strconv.Atoi(stageDomPricingModel.SITPickupDeliverySchedule)
		if spErr != nil {
			return fmt.Errorf("Failed to parse SITPickupDeliverySchedule for %+v: %w", stageDomPricingModel, spErr)
		}
		serviceAreaNumber := correctServiceAreaNumber(stageDomPricingModel.ServiceAreaNumber)
		var serviceArea models.ReDomesticServiceArea
		err := db.Where("service_area = $1 and services_schedule = $2 and sit_pd_schedule = $3", serviceAreaNumber, servicesSchedule, sITPDSchedule).First(&serviceArea)
		if err != nil || serviceArea.ID == uuid.Nil {
			return fmt.Errorf("Cannot find service area number '%s' with services schedule '%d' and SITPickupDeliverySchedule '%d': %w", serviceAreaNumber, servicesSchedule, sITPDSchedule, err)
		}

		//DSH - ShorthaulPrice
		service, err := models.FetchReServiceItem(db, "DSH")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DSH: %w", err)
		}

		cents, convErr := stringDollarsToCents(stageDomPricingModel.ShorthaulPrice)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for Shorthaul data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            cents,
		}

		domPricingModels = append(domPricingModels, domPricingModelDSH)

		//DODP - OriginDestinationPrice
		service, err = models.FetchReServiceItem(db, "DODP")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DODP: %w", err)
		}

		cents, convErr = stringDollarsToCents(stageDomPricingModel.OriginDestinationPrice)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationPrice data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDODP := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            cents,
		}

		domPricingModels = append(domPricingModels, domPricingModelDODP)

		//DFSIT - OriginDestinationSITFirstDayWarehouse
		service, err = models.FetchReServiceItem(db, "DFSIT")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DFSIT: %w", err)
		}

		cents, convErr = stringDollarsToCents(stageDomPricingModel.OriginDestinationSITFirstDayWarehouse)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationSITFirstDayWarehouse data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDFSIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            cents,
		}

		domPricingModels = append(domPricingModels, domPricingModelDFSIT)

		//DASIT - OriginDestinationSITAddlDays
		service, err = models.FetchReServiceItem(db, "DASIT")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DASIT: %w", err)
		}

		cents, convErr = stringDollarsToCents(stageDomPricingModel.OriginDestinationSITAddlDays)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationSITAddlDays data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDASIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            cents,
		}

		domPricingModels = append(domPricingModels, domPricingModelDASIT)

		for _, model := range domPricingModels {
			verrs, err := db.ValidateAndSave(&model)
			if err != nil {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with error: %w", model, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}
