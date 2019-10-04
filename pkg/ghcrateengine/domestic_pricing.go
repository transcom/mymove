package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

type DomesticServicePricingData struct {
	MoveDate      time.Time
	ServiceAreaID uuid.UUID
	Distance      unit.Miles
	Weight        unit.Pound
	IsPeakPeriod  bool
	ContractCode  string
}

func lookupDomesticLinehaulRate(db *pop.Connection, d DomesticServicePricingData) unit.Millicents {
	//To be implemented when models are created
	//query := db.Where(
	//	"isPeakPeriod = d.isPeakPeriod").Where(
	//		"weight gte weight_lower").Where(
	//			"weight lte weight_upper").Where(
	//				"distance gte miles_lower").Where(
	//					"distance lte miles_upper").Where(
	//										"re_contracts.code=?", d.ContractCode).Where(
	//											"d.MoveDate gte re_contract_years.start_date").Where(
	//												"d.MoveDate lte re_contract_years.end_date").Join(
	//						"serviceAreas sa", "sa.id=re_domestic_linehaul_prices.service_area_id").Join(
	//								"re_contracts", "re_contract.id=re_domestic_linehaul_prices.contract_id").Join(
	//										"re_contract_years cy", "re_contract.id=cy.contract_id")
	//domesticLinehaulPricing := models.DomesticLinehaulPricing{}
	//err := db.query.All(&domesticLinehaulPricing)
	//if err != nil {
	//	return err
	//}

	var stubPrice unit.Millicents = 272700
	return stubPrice
}

// Calculation Functions
// CalculateDemesticLinehaul calculates the cost domestic linehaul and returns the cost in millicents
func (gre *GHCRateEngine) CalculateDomesticLinehaul(d DomesticServicePricingData) unit.Millicents {
	rate := lookupDomesticLinehaulRate(gre.db, d)
	// TODO: look up escalation and multiply rate by escalation
	// calculate total
	return rate.Multiply(d.Weight.ToCWT().Int())
}
