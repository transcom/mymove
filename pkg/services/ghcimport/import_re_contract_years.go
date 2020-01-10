package ghcimport

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREContractYears(dbTx *pop.Connection) error {
	// populate contractYearsToIDMap
	var priceEscalationDiscounts []models.StagePriceEscalationDiscount
	err := dbTx.All(&priceEscalationDiscounts)
	if err != nil {
		return fmt.Errorf("could not read staged price escalation discounts: %w", err)
	}

	gre.contractYearToIDMap = make(map[string]uuid.UUID)
	incrementYear := 0
	compoundedEscalation := 1.00000

	//loop through the price escalation discounts data and pull contract year and escalations
	for _, stagePriceEscalationDiscount := range priceEscalationDiscounts {
		escalation, err := strconv.ParseFloat(stagePriceEscalationDiscount.ForecastingAdjustment, 64)
		if err != nil {
			return fmt.Errorf("could not process forecast adjustment [%s]: %w", stagePriceEscalationDiscount.ForecastingAdjustment, err)
		}

		escalationCompounded, err := strconv.ParseFloat(stagePriceEscalationDiscount.PriceEscalation, 64)
		if err != nil {
			return fmt.Errorf("could not process price escalation [%s]: %w", stagePriceEscalationDiscount.PriceEscalation, err)
		}
		compoundedEscalation *= escalationCompounded

		startDate := time.Date(2018, time.June, 01, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2019, time.May, 31, 0, 0, 0, 0, time.UTC)
		incrementYear++

		contractYear := models.ReContractYear{
			ContractID:           gre.contractID,
			Name:                 stagePriceEscalationDiscount.ContractYear,
			StartDate:            startDate.AddDate(incrementYear, 0, 0),
			EndDate:              endDate.AddDate(incrementYear, 0, 0),
			Escalation:           escalation,
			EscalationCompounded: compoundedEscalation,
		}

		verrs, dbErr := dbTx.ValidateAndSave(&contractYear)
		if dbErr != nil {
			return fmt.Errorf("error saving ReContractYears: %+v with error: %w", contractYear, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReContractYears: %+v with validation errors: %w", contractYear, verrs)
		}

		//add to map
		gre.contractYearToIDMap[contractYear.Name] = contractYear.ID
	}

	return nil
}
