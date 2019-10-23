package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// minDomesticWeight is the minimum weight used in domestic calculations (weights below this are upgraded to the min)
const minDomesticWeight = unit.Pound(500)

// NewDomesticLinehaulPricer is the public constructor for a DomesticLinehaulPricer using Pop
func NewDomesticLinehaulPricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticLinehaulPricer {
	return &domesticLinehaulPricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// domesticLinehaulPricer is a service object to price domestic linehaul
type domesticLinehaulPricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

// PriceDomesticLinehaul produces the price in cents for the linehaul charge for the given move parameters
func (p domesticLinehaulPricer) PriceDomesticLinehaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error) {
	// Validate parameters
	if moveDate.IsZero() {
		return 0, errors.New("MoveDate is required")
	}
	if distance <= 0 {
		return 0, errors.New("Distance must be greater than 0")
	}
	if weight <= 0 {
		return 0, errors.New("Weight must be greater than 0")
	}
	if len(serviceArea) == 0 {
		return 0, errors.New("ServiceArea is required")
	}

	// Minimum weight is 500 pounds
	effectiveWeight := weight
	if weight < minDomesticWeight {
		effectiveWeight = minDomesticWeight
	}

	isPeakPeriod := IsPeakPeriod(moveDate)

	var pe priceAndEscalation
	query :=
		`select price_millicents, escalation_compounded
         from re_domestic_linehaul_prices dlp
         inner join re_contracts c on dlp.contract_id = c.id
         inner join re_contract_years cy on c.id = cy.contract_id
         inner join re_domestic_service_areas dsa on dlp.domestic_service_area_id = dsa.id
         where c.code = $1
         and $2 between cy.start_date and cy.end_date
         and dlp.is_peak_period = $3
         and $4 between dlp.weight_lower and dlp.weight_upper
         and $5 between dlp.miles_lower and dlp.miles_upper
         and dsa.service_area = $6;`
	err := p.db.RawQuery(
		query,
		p.contractCode,
		moveDate,
		isPeakPeriod,
		effectiveWeight,
		distance,
		serviceArea).First(&pe)
	if err != nil {
		return 0, errors.Wrap(err, "Lookup of domestic linehaul rate failed")
	}

	baseTotalPrice := effectiveWeight.ToCWTFloat64() * distance.Float64() * pe.PriceMillicents.Float64()
	escalatedTotalPrice := pe.EscalationCompounded * baseTotalPrice

	// TODO: Round or truncate?
	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	p.logger.Info("Base domestic linehaul calculated",
		zap.String("contractCode", p.contractCode),
		zap.Time("moveDate", moveDate),
		zap.Int("distance", distance.Int()),
		zap.Int("weight", weight.Int()),
		zap.String("serviceArea", serviceArea),
		zap.Bool("isPeakPeriod", isPeakPeriod),
		zap.Int("effectiveWeight", effectiveWeight.Int()),
		zap.Object("priceAndEscalation", pe),
		zap.Float64("baseTotalPrice", baseTotalPrice),
		zap.Float64("escalatedTotalPrice", escalatedTotalPrice),
		zap.Int("totalPriceCents", totalPriceCents.Int()),
	)

	return totalPriceCents, err
}
