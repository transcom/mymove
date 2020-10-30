package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gobuffalo/pop/v5"

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

// DomesticServiceAreaPricer is a service object to price domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days (does not include shorthaul)
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
}

func (dsa *domesticServiceAreaPricer) PriceDomesticServiceArea(moveDate time.Time, weight unit.Pound, serviceArea string, serviceCode string) (cost unit.Cents, err error) {
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
		return 0, errors.New("ServiceCode is required")
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
	totalCost := unit.Cents(math.Round(escalatedTotalPrice))

	dsa.logger.Info(fmt.Sprintf("%s calculated", serviceCode), // May change to use ServiceName
		zap.String("contractCode:", dsa.contractCode),
		zap.String("serviceCode:", serviceCode),
		zap.Time("moveDate:", moveDate),
		zap.String("serviceArea:", serviceArea),
		zap.Int("weight (lb):", weight.Int()),
		zap.Int("effectiveWeight (lb):", effectiveWeight.Int()),
		zap.Bool("isPeakPeriod: ", isPeakPeriod),
		zap.Object("centPriceAndEscalation:", priceAndEscalation),
		zap.Float64("baseCost (cents):", baseTotalPrice),
		zap.Float64("escalatedCost (cents):", escalatedTotalPrice),
		zap.Int("totalCost (cents):", totalCost.Int()),
	)

	return totalCost, err
}
