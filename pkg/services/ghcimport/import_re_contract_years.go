package ghcimport

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/pop/v5"
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

	basePeriodEndDateForPrimeContract1 := gre.ContractStartDate.AddDate(1, 0, -1)

	//loop through the price escalation discounts data and pull contract year and escalations
	for _, stagePriceEscalationDiscount := range priceEscalationDiscounts {
		escalation, err := strconv.ParseFloat(stagePriceEscalationDiscount.PriceEscalation, 64)
		if err != nil {
			return fmt.Errorf("could not process price escalation [%s]: %w", stagePriceEscalationDiscount.PriceEscalation, err)
		}
		compoundedEscalation *= escalation

		contractYear := models.ReContractYear{
			ContractID:           gre.ContractID,
			Name:                 stagePriceEscalationDiscount.ContractYear,
			StartDate:            gre.ContractStartDate.AddDate(incrementYear, 0, 0),
			EndDate:              basePeriodEndDateForPrimeContract1.AddDate(incrementYear, 0, 0),
			Escalation:           escalation,
			EscalationCompounded: compoundedEscalation,
		}
		incrementYear++

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
