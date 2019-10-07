package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/unit"
)

type DomesticServicePricingData struct {
	MoveDate      time.Time
	ServiceAreaID uuid.UUID
	Distance      unit.Miles
	Weight        unit.CWTFloat // record this here as 5.00 if actualWt less than minimum of 5.00 cwt (500lb)
	IsPeakPeriod  bool
	ContractCode  string
}

func lookupDomesticLinehaulRate(db *pop.Connection, d DomesticServicePricingData) (rate unit.Millicents, err error) {
	// TODO: check/correct syntax && implement when models are created
	// query := db.Where(
	// 	"is_peak_period = d.IsPeakPeriod").Join(
	// 	"serviceAreas sa", "sa.id=re_domestic_linehaul_prices.service_area_id").Join(
	// 	"re_contracts c", "re_contract.id=re_domestic_linehaul_prices.contract_id").Join(
	// 	"re_contract_years cy", "re_contract.id=cy.contract_id").Where(
	// 		"sa.id=d.ServiceAreaID").Where(
	// 		"weight gte weight_lower").Where(
	// 			"weight lte weight_upper").Where(
	// 				"distance gte miles_lower").Where(
	// 					"distance lte miles_upper").Where(
	// 										"re_contracts.code=?", d.ContractCode).Where(
	// 											"d.MoveDate gte cy.start_date").Where(
	// 												"d.MoveDate lte cy.end_date")

	// domesticLinehaulPricing := models.DomesticLinehaulPricing{}
	// err := db.query.All(&domesticLinehaulPricing)
	// if err != nil {
	// 	return err
	// }

	rate = 272700 // stubbed

	return rate, err
}

// Calculation Functions
// CalculateBaseDomesticLinehaul calculates the cost domestic linehaul and returns the cost in millicents
func (gre *GHCRateEngine) CalculateBaseDomesticLinehaul(d DomesticServicePricingData) (cost unit.Millicents, err error) {
	rate, err := lookupDomesticLinehaulRate(gre.db, d)

	if err != nil {
		return cost, errors.Wrap(err, "Lookup of domestic linehaul rate failed")
	}

	cost = rate.MultiplyFloat64(float64(d.Weight))

	gre.logger.Info("Base domestic linehaul calculated",
		zap.Time("move date", d.MoveDate),
		zap.String("service area ID", d.ServiceAreaID.String()),
		zap.String("distance in miles", d.Distance.String()),
		zap.Float64("centiweight", float64(d.Weight)),
		zap.Bool("is peak period", d.IsPeakPeriod),
		zap.String("contract code", d.ContractCode),
		zap.Int("base rate (millicents)", rate.Int()),
		zap.Int("calculated cost (millicents)", cost.Int()),
	)

	return cost, err
}
