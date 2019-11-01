package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// NewDomesticServiceAreaPricer is the public constructor for a DomesticRateAreaPricer using Pop
func NewDomesticShorthaulPricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticShorthaulPricer {
	return &domesticShorthaulPricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// DomesticShorthaulPricer is a service object to price domestic shorthaul
type domesticShorthaulPricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

func (dsh *domesticShorthaulPricer) PriceDomesticShorthaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (totalCost unit.Cents, err error) {
	// Validate parameters
	if moveDate.IsZero() {
		return 0, errors.New("MoveDate is required")
	}
	if weight <= 0 {
		return 0, errors.New("Weight must be greater than 0")
	}
	if distance <= 0 {
		return 0, errors.New("Distance must be greater than 0")
	}
	if len(serviceArea) == 0 {
		return 0, errors.New("ServiceArea is required")
	}

	pe := centPriceAndEscalation{}
	isPeakPeriod := IsPeakPeriod(moveDate)
	// look up rate for shorthaul
	shorthaulServiceCode := "DSH"
	query := `
		select price_cents, escalation_compounded from re_domestic_service_area_prices dsap
		inner join re_domestic_service_areas sa on dsap.domestic_service_area_id = sa.id
		inner join re_services on dsap.service_id = re_services.id
		inner join re_contracts on re_contracts.id = dsap.contract_id
		inner join re_contract_years on re_contracts.id = re_contract_years.contract_id
		where sa.service_area = $1
		and re_services.code = $2
		and re_contracts.code = $3
		and dsap.is_peak_period = $4
		and $5 between re_contract_years.start_date and re_contract_years.end_date;
	`
	err = dsh.db.RawQuery(
		query, serviceArea, shorthaulServiceCode, dsh.contractCode, isPeakPeriod, moveDate).First(
		&pe)

	effectiveWeight := weight
	if weight <= minDomesticWeight {
		effectiveWeight = minDomesticWeight
	}

	basePrice := pe.PriceCents.Float64() * distance.Float64() * effectiveWeight.ToCWTFloat64()
	escalatedPrice := basePrice * pe.EscalationCompounded
	totalCost = unit.Cents(escalatedPrice) // TODO: truncates the price to get an integer- is that what we want?

	dsh.logger.Info(fmt.Sprintf("%s calculated", shorthaulServiceCode), // May change to use ServiceName
		zap.String("contractCode:", dsh.contractCode),
		zap.String("serviceCode:", shorthaulServiceCode),
		zap.Time("moveDate:", moveDate),
		zap.String("serviceArea:", serviceArea),
		zap.Int("distance (mi):", distance.Int()),
		zap.Int("weight (lb):", weight.Int()),
		zap.Int("effectiveWeight (lb):", effectiveWeight.Int()),
		zap.Bool("isPeakPeriod: ", isPeakPeriod),
		zap.Object("centPriceAndEscalation:", pe),
		zap.Float64("baseCost (cents):", basePrice),
		zap.Float64("escalatedCost (cents):", escalatedPrice),
		zap.Int("totalCost (cents):", totalCost.Int()),
	)
	return totalCost, err
}
