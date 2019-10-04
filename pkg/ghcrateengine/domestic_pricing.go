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
	CalcCWTWeight unit.CWTFloat // record this here as 5.00 if actualWt less than minimum of 5.00 cwt (500lb)
	IsPeakPeriod  bool
	ContractCode  string
}

func lookupDomesticLinehaulRate(db *pop.Connection, d DomesticServicePricingData) unit.Millicents {
	// TODO: To be implemented when models are created
	//query := db.Where(
	//	"isPeakPeriod = d.isPeakPeriod").Join(
	//	"serviceAreas sa", "sa.id=re_domestic_linehaul_prices.service_area_id").Join(
	//	"re_contracts c", "re_contract.id=re_domestic_linehaul_prices.contract_id").Join(
	//	"re_contract_years cy", "re_contract.id=cy.contract_id").Where(
	//		"weight gte weight_lower").Where(
	//			"weight lte weight_upper").Where(
	//				"distance gte miles_lower").Where(
	//					"distance lte miles_upper").Where(
	//										"re_contracts.code=?", d.ContractCode).Where(
	//											"d.MoveDate gte cy.start_date").Where(
	//												"d.MoveDate lte cy.end_date")

	//domesticLinehaulPricing := models.DomesticLinehaulPricing{}
	//err := db.query.All(&domesticLinehaulPricing)
	//if err != nil {
	//	return err
	//}

	var stubPrice unit.Millicents = 272700
	return stubPrice
}

func lookupContractYearEscalation(db *pop.Connection, moveDate time.Time, contractCode string) float64 {
	// TODO: look up contract using contractCode and move Date
	// select escalation from re_contracts innerjoin re_contract_years on re_contract_years.id = re_contract_years.contract_id
	// where contract_code = contractCode and moveDate is between contract_years.start_date and contract_years.end_date
	stubEscalation := 1.02
	return stubEscalation
}

// Calculation Functions
// CalculateDemesticLinehaul calculates the cost domestic linehaul and returns the cost in millicents
func (gre *GHCRateEngine) CalculateDomesticLinehaul(d DomesticServicePricingData) unit.Millicents {
	rate := lookupDomesticLinehaulRate(gre.db, d)
	escalation := lookupContractYearEscalation(gre.db, d.MoveDate, d.ContractCode) // type float64
	escalatedRate := rate.MultiplyFloat64(escalation)

	return escalatedRate.MultiplyFloat64(float64(d.CalcCWTWeight))
}
