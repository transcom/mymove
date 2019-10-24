package ghcrateengine

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// NewDomesticServiceAreaPricer is the public constructor for a DomesticRateAreaPricer using Pop
func NewDomesticServiceAreaPricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticServiceAreaPricer {
	return &domesticServiceAreaPricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// DomesticServiceAreaPricer is a service object to price domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
// domestic other prices: pack, unpack, and sit p/d costs for a GHC move
type domesticServiceAreaPricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

func lookupDomesticServiceAreaPrice(db *pop.Connection, moveDate time.Time, serviceArea string, serviceCode string, contractCode string, isPeakPeriod bool) (pe centPriceAndEscalation, err error) {
	pe = centPriceAndEscalation{}

	query := `
		select price_cents, escalation_compounded from re_domestic_service_area_prices dsap
		inner join re_domestic_service_areas sa on sa.id = dsap.domestic_service_area_id
		inner join re_contracts on re_contracts.id = dsap.contract_id
		inner join re_contract_years on re_contracts.id = re_contract_years.contract_id
		inner join re_services on re_services.id = dsap.service_id
		where sa.service_area = $1
		and re_services.code = $2
		and re_contracts.code = $3
		and dsap.is_peak_period = $4
		and $5 between re_contract_years.start_date and re_contract_years.end_date;
	`
	err = db.RawQuery(
		query, serviceArea, serviceCode, contractCode, isPeakPeriod, moveDate).First(
		&pe)
	if err != nil {
		return pe, errors.Wrap(err, "Fetch domestic service area price failed")
	}
	return pe, err

	//stubbedRate, err := unit.Cents(689), nil
	//return stubbedRate, err
}

func (dsa *domesticServiceAreaPricer) PriceDomesticServiceArea (moveDate time.Time, weight unit.Pound, serviceArea string, serviceCode string) (cost unit.Cents, err error) {
	// Validate parameters
	if moveDate.IsZero() {
		return 0, errors.New("MoveDate is required")
	}
	if weight <= 0 {
		return 0, errors.New("Weight must be greater than 0")
	}
	if len(serviceArea) == 0 {
		return 0, errors.New("ServiceArea is required")
	}
	if len(serviceCode) == 0 {
		return 0, errors.New("ServicesCode is required")
	}

	isPeakPeriod := IsPeakPeriod(moveDate)
	priceAndEscalation, err := lookupDomesticServiceAreaPrice(dsa.db, moveDate, serviceArea, serviceCode, dsa.contractCode, isPeakPeriod)
	if err != nil {
		return cost, errors.Wrap(err, fmt.Sprintf("Lookup of domestic service %s failed", serviceCode))
	}

	effectiveWeight := weight
	if weight <= minDomesticWeight {
		effectiveWeight = minDomesticWeight
	}

	baseTotalPrice := priceAndEscalation.PriceCents.Float64() * effectiveWeight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * priceAndEscalation.EscalationCompounded
	totalCost := unit.Cents(escalatedTotalPrice)

	dsa.logger.Info(fmt.Sprintf("%s calculated", serviceCode), // May change to use ServiceName
		zap.String("contract code", dsa.contractCode),
		zap.String("service code", serviceCode),
		zap.Time("move date", moveDate),
		zap.String("service area", serviceArea),
		zap.Float64("weight lb", float64(weight)),
		zap.Float64("effective weight lb", float64(effectiveWeight)),
		zap.Bool("is peak period", isPeakPeriod),
		zap.Object("centPriceAndEscalation", priceAndEscalation),
		zap.Float64("base cost (cents)", baseTotalPrice),
		zap.Float64("escalated cost (cents)", baseTotalPrice),
		zap.Int("calculated cost (cents)", cost.Int()),
	)

	return totalCost, err
}
