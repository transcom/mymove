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

func lookupDomesticLinehaulRate(db *pop.Connection, pricingData DomesticServicePricingData) (rate unit.Millicents, err error) {

	//domesticLinehaulPrice := models.ReDomesticLinehaulPrice{}
	//query := `
	//	select price_cents from re_domestic_linehaul_prices lh
	//	inner join re_domestic_service_areas sa on sa.id = lh.domestic_service_area_id
	//	inner join re_contracts on re_contracts.id = lh.contract_id
	//	inner join re_contract_years on re_contracts.id = re_contract_years.contract_id
	//	where lh.is_peak_period = $1
	//	and $2 between lh.weight_lower and lh.weight_upper
	//	and $3 between lh.miles_lower and lh.miles_upper
	//	and $4 between re_contract_year.start_date and re_contract_year.end_date
	//	and $5 re_contracts.code = $6
	//`
	//err = db.RawQuery(
	//	query, pricingData.IsPeakPeriod, pricingData.Weight, pricingData.Distance, pricingData.MoveDate, pricingData.ContractCode).First(
	//	&domesticLinehaulPrice)
	//if err != nil {
	//	return rate, errors.Wrap(err, "Fetch domestic other price failed")
	//}

	//rate = domesticLinehaulPrice.PriceMillicents TODO: uncomment and remove stubbed rate when connected to db

	//return rate, err
	stubbedRate := unit.Millicents(272700) // stubbed

	return stubbedRate, err
}

func lookupDomesticServiceAreaPrice(db *pop.Connection, pricingData DomesticServicePricingData) (rate unit.Cents, err error) {
	//domesticServiceAreaPrice := models.ReDomesticServiceAreaPrice{}

	//query := `
	//	select price_price from re_domestic_service_area_prices dsap
	//	inner join re_domestic_service_areas sa on sa.id = dsap.domestic_service_area_id
	//	inner join re_contracts on re_contracts.id = dsap.contract_id
	//	inner join re_contract_years on re_contracts.id = re_contract_years.contract_id
	//	inner join re_services on re_services.id = dsap.services_id
	//	where sa.id = $1
	//	and re_services.code = $2
	//	and re_contracts.code = $3
	//	and dsap.is_peak_period = $4
	//	and $5 between re_contract_year.start_date and re_contract_year.end_date;'
	//`
	//err = db.RawQuery(
	//	query, pricingData.ServiceAreaID, pricingData.ServiceAreaID, pricingData.ContractCode, pricingData.IsPeakPeriod, pricingData.MoveDate).First(
	//	&domesticServiceAreaPrice)
	//if err != nil {
	//	return rate, errors.Wrap(err, "Fetch domestic other price failed")
	//}
	//rate = domesticServiceAreaPrice.PriceCents
	//return rate, err TODO: uncomment and remove stubbed rate when connected to db

	stubbedRate, err := unit.Cents(689), nil
	return stubbedRate, err
}

func lookupDomesticOtherPrice(db *pop.Connection, pd DomesticServicePricingData) (rate unit.Cents, err error) {
	//var domesticOtherPrice models.ReDomesticOtherPrice
	//domesticServiceArea := models.ReDomesticServiceArea{}
	//
	//var schedule int
	//sitPDServiceCode := "SITPD" // stubbed service code
	//packServiceCode := "DPK"    // stubbed service code
	//unpackServiceCode := "DUPK" // stubbed service code
	//
	//err = db.Q().Where(fmt.Sprintf("id = %s", pd.ServiceAreaID)).Last(&domesticServiceArea)
	//if pd.ServiceCode == sitPDServiceCode {
	//	schedule = domesticServiceArea.SITPDSchedule
	//} else if pd.ServiceCode == packServiceCode || pd.ServiceCode == unpackServiceCode {
	//	schedule = domesticServiceArea.ServiceSchedule
	//} else {
	//	// throw error??
	//	return rate, errors.Wrap(err, "must be pack, unpack, or SIT P/D service")
	//}
	//
	//query := `
	//	select per_unit_price from re_domestic_other_prices
	//     inner join re_contracts on re_domestic_other_prices.contract_id = re_contracts.id
	//     inner join re_contract_years on re_contracts.id = re_contract_years.contract_id
	//	 inner join re_services on re_services.id = re_domestic_other_prices.services_id
	//     where re_domestic_other_prices.schedule = $1
	//	  and re_contracts.code = $2
	//	  and $3 between re_contract_years.start_date and re_contract_years.end_date
	//	  and re_domestic_linehaul_prices.is_peak_period = $4
	//	  and $5 between re_domestic_linehaul_prices.weight_lower and re_domestic_linehaul_prices.weight_upper
	//	  and re_domestic_service_areas.id = $6;
	//`
	//
	//err = db.RawQuery(
	//	query, schedule, pd.ContractCode, pd.MoveDate, pd.IsPeakPeriod, pd.Weight, pd.ServiceAreaID).First(
	//	&domesticOtherPrice)
	//if err != nil {
	//	return rate, errors.Wrap(err, "Fetch domestic other price failed")
	//}

	//return domesticOtherPrice.PerUnitPrice, err TODO: uncomment and remove stubbed rate when connected to db

	stubbedRate, err := unit.Cents(23440), nil

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
func (gre *GHCRateEngine) CalculateBaseDomesticPerWeightServiceCost(d DomesticServicePricingData, isDomesticOtherService bool) (cost unit.Cents, err error) {
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
func (gre *GHCRateEngine) CalculateBaseDomesticShorthaulCost(d DomesticServicePricingData) (cost unit.Cents, err error) {
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
