package ghcrateengine

import (
	"time"

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

func (d *domesticShorthaulPricer) PriceDomesticShorthaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (totalCost unit.Cents, err error) {

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
	err = d.db.RawQuery(
		query, serviceArea, shorthaulServiceCode, d.contractCode, isPeakPeriod, moveDate).First(
		&pe)
	// multiply by distance and rate and weight
	basePrice := pe.PriceCents.Float64() * distance.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * pe.EscalationCompounded
	totalCost = unit.Cents(escalatedPrice)

	return totalCost, err
}
