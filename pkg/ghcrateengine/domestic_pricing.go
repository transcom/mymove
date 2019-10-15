package ghcrateengine

import (
	"fmt"
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
	ServiceCode   string // may change to Service when model is available
}

func lookupDomesticLinehaulRate(db *pop.Connection, d DomesticServicePricingData) (rate unit.Millicents, err error) {
	// TODO: check/correct syntax && implement when models are created
	// query := db.Where(
	// 	"is_peak_period = d.IsPeakPeriod").Join(
	// 	"service_areas sa", "sa.id=re_domestic_linehaul_prices.service_area_id").Join(
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

func lookupDomesticServiceAreaPrice (db *pop.Connection, pricingData DomesticServicePricingData) (cost unit.Cents, err error) {
	// select cost from re_domestic_rate_area_prices
		// joined with re_service_area, re_contracts, re_contract_years, re_services
		// where:
			// re_service_area.id = pricingData.ServiceAreaID
			// re_services.code = pricingData.ServiceCode
			// re_contracts.code = pricingData.ContractCode
			// is_peak_period = pricingData.IsPeakPeriod
			// move date is between start and end dates of contract year

	stubbedRate, err := unit.Cents(689), nil
	return stubbedRate, err
}

func lookupDomesticOtherPrice (db *pop.Connection, pricingData DomesticServicePricingData) (cost unit.Cents, err error) {
	serviceCode := "DSIT" // stubbed service code
	var stubbedRate unit.Cents = 0
	if pricingData.ServiceCode == serviceCode {
		// select cost from re_domestic_other_prices
		// joined with re_service_area, re_contracts, re_contract_years, re_services
		// where:
		// re_domestic_service_areas.s_i_t_p_d_schedule = schedule
		// re_services.code = pricingData.ServiceCode
		// re_contracts.code = pricingData.ContractCode
		// is_peak_period = pricingData.IsPeakPeriod
		// move date is between start and end dates of contract year
		stubbedRate, err = 21796, nil
	} else {
		// select cost from re_domestic_other_prices
		// joined with re_service_area, re_contracts, re_contract_years, re_services
		// where:
		// re_domestic_service_areas.services_schedule = schedule
		// re_services.code = pricingData.ServiceCode
		// re_contracts.code = pricingData.ContractCode
		// is_peak_period = pricingData.IsPeakPeriod
		// move date is between start and end dates of contract year
		stubbedRate, err = 6333, nil
	}

	return stubbedRate, err
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
		zap.Int("calculated base cost (millicents)", cost.Int()),
	)

	return cost, err
}

// CalculateBaseDomesticPerWeightCost calculates the cost based on service performed and returns the cost in cents
// This function is used to calculate
	// domestic prices: origin and destination service area, SIT day 1, SIT days-1,
	// domestic other prices: pack, unpack, and sit p/d costs
func (gre *GHCRateEngine) CalculateBaseDomesticPerWeightServiceCost (d DomesticServicePricingData, isDomesticOtherService bool) (cost unit.Cents, err error) {
	var rate unit.Cents
	if isDomesticOtherService {
		rate, err = lookupDomesticOtherPrice(gre.db, d)
		if err != nil {
			return cost, errors.Wrap(err, fmt.Sprintf("Lookup of domestic service %s failed", d.ServiceCode))
		}
	} else {
		rate, err = lookupDomesticServiceAreaPrice(gre.db, d)
		if err != nil {
			return cost, errors.Wrap(err, fmt.Sprintf("Lookup of domestic service %s failed", d.ServiceCode))
		}
	}

	cost = rate.MultiplyCWTFloat(d.Weight)

	gre.logger.Info(fmt.Sprintf("%s calculated", d.ServiceCode), // May change to use ServiceName
		zap.String("service code", d.ServiceCode),
		zap.Time("move date", d.MoveDate),
		zap.String("service area ID", d.ServiceAreaID.String()),
		zap.Float64("centiweight", float64(d.Weight)),
		zap.Bool("is peak period", d.IsPeakPeriod),
		zap.String("contract code", d.ContractCode),
		zap.Int("base rate (cents)", rate.Int()),
		zap.Int("calculated base cost (cents)", cost.Int()),
	)

	return cost, err
}


// CalculateBaseDomesticShorthaulCost calculates the cost based on service performed and returns the cost in cents
func (gre *GHCRateEngine) CalculateBaseDomesticShorthaulCost (d DomesticServicePricingData) (cost unit.Cents, err error) {
	rate, err := lookupDomesticServiceAreaPrice(gre.db, d)
	if err != nil {
		return cost, errors.Wrap(err, fmt.Sprintf("Lookup of domestic service %s failed", d.ServiceCode))
	}
	costPerWeight := rate.MultiplyCWTFloat(d.Weight)
	cost = costPerWeight.MultiplyMiles(d.Distance)

	gre.logger.Info(fmt.Sprintf("%s calculated", d.ServiceCode), // May change to use ServiceName
		zap.String("service code", d.ServiceCode),
		zap.Time("move date", d.MoveDate),
		zap.String("service area ID", d.ServiceAreaID.String()),
		zap.String("distance in miles", d.Distance.String()),
		zap.Float64("centiweight", float64(d.Weight)),
		zap.Bool("is peak period", d.IsPeakPeriod),
		zap.String("contract code", d.ContractCode),
		zap.Int("base rate (cents)", rate.Int()),
		zap.Int("calculated cost (cents)", cost.Int()),
	)

	return cost, err
}
