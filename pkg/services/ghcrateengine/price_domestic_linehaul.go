package ghcrateengine

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// priceAndEscalation is used to hold data returned by the database query
type priceAndEscalation struct {
	PriceMillicents      unit.Millicents `db:"price_millicents"`
	EscalationCompounded float64         `db:"escalation_compounded"`
}

// MarshalLogObject allows priceAndEscalation to be logged by zap
func (p priceAndEscalation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("PriceMillicents", p.PriceMillicents.Int())
	encoder.AddFloat64("EscalationCompounded", p.EscalationCompounded)
	return nil
}

// PriceDomesticLinehaul produces the price in cents for the linehaul charge for the given move parameters
func (p domesticLinehaulPricer) PriceDomesticLinehaul(data services.DomesticServicePricingData) (unit.Cents, error) {
	// TODO: Validate params

	// Minimum weight is 500 pounds
	effectiveWeight := data.Weight
	if data.Weight < minDomesticWeight {
		effectiveWeight = minDomesticWeight
	}

	isPeakPeriod := IsPeakPeriod(data.MoveDate)

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
		data.MoveDate,
		isPeakPeriod,
		effectiveWeight,
		data.Distance,
		data.ServiceArea).First(&pe)
	if err != nil {
		return 0, errors.Wrap(err, "Lookup of domestic linehaul rate failed")
	}

	baseTotalPrice := effectiveWeight.ToCWTFloat64() * data.Distance.Float64() * pe.PriceMillicents.Float64()
	escalatedTotalPrice := pe.EscalationCompounded * baseTotalPrice

	// TODO: Round or truncate?
	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	p.logger.Info("Base domestic linehaul calculated",
		zap.String("contractCode", p.contractCode),
		zap.Object("input", data),
		zap.Bool("isPeakPeriod", isPeakPeriod),
		zap.Int("effectiveWeight", effectiveWeight.Int()),
		zap.Object("priceAndEscalation", pe),
		zap.Float64("baseTotalPrice", baseTotalPrice),
		zap.Float64("escalatedTotalPrice", escalatedTotalPrice),
		zap.Int("totalPriceCents", totalPriceCents.Int()),
	)

	return totalPriceCents, err
}
