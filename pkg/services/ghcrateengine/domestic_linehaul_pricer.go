package ghcrateengine

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dlhPricerMinimumWeight   = 500
	dlhPricerMinimumDistance = 50
)

type domesticLinehaulPricer struct {
	db *pop.Connection
}

// NewDomesticLinehaulPricer creates a new pricer for domestic linehaul services
func NewDomesticLinehaulPricer(db *pop.Connection) services.DomesticLinehaulPricer {
	return &domesticLinehaulPricer{
		db: db,
	}
}

// Price determines the price for a domestic linehaul
func (p domesticLinehaulPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, distance int, weightBilledActual int, serviceArea string) (unit.Cents, error) {
	priceAndEscalation, err := fetchDomesticLinehaulPrice(p.db, contractCode, requestedPickupDate, isPeakPeriod, distance, weightBilledActual, serviceArea)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
	}

	weightPounds := unit.Pound(weightBilledActual)
	distanceMiles := unit.Miles(distance)
	baseTotalPrice := weightPounds.ToCWTFloat64() * distanceMiles.Float64() * priceAndEscalation.PriceMillicents.Float64()
	escalatedTotalPrice := priceAndEscalation.EscalationCompounded * baseTotalPrice

	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	return totalPriceCents, nil

}

// PriceUsingParams determines the price for a domestic linehaul given PaymentServiceItemParams
func (p domesticLinehaulPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), err
	}

	distanceZip3, err := getParamInt(params, models.ServiceItemParamNameDistanceZip3)
	if err != nil {
		return unit.Cents(0), err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), err
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	return p.Price(contractCode, requestedPickupDate, isPeakPeriod, distanceZip3, weightBilledActual, serviceAreaOrigin)
}

func fetchDomesticLinehaulPrice(db *pop.Connection, contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, distance int, weight int, serviceArea string) (milliCentPriceAndEscalation, error) {
	// Validate parameters
	if requestedPickupDate.IsZero() {
		return milliCentPriceAndEscalation{}, errors.New("MoveDate is required")
	}
	if distance < dlhPricerMinimumDistance {
		return milliCentPriceAndEscalation{}, fmt.Errorf("distance must be at least %d", dlhPricerMinimumDistance)
	}
	if weight < dlhPricerMinimumWeight {
		return milliCentPriceAndEscalation{}, fmt.Errorf("weight must be at least %d", dlhPricerMinimumWeight)
	}
	if len(serviceArea) == 0 {
		return milliCentPriceAndEscalation{}, errors.New("ServiceArea is required")
	}

	var mpe milliCentPriceAndEscalation
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
	err := db.RawQuery(
		query,
		contractCode,
		requestedPickupDate,
		isPeakPeriod,
		weight,
		distance,
		serviceArea).First(&mpe)
	if err != nil {
		return mpe, errors.Wrap(err, "Lookup of domestic linehaul rate failed")
	}

	return mpe, nil

}

// priceAndEscalation is used to hold data returned by the database query
type milliCentPriceAndEscalation struct {
	PriceMillicents      unit.Millicents `db:"price_millicents"`
	EscalationCompounded float64         `db:"escalation_compounded"`
}

func (p milliCentPriceAndEscalation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("PriceMillicents", p.PriceMillicents.Int())
	encoder.AddFloat64("EscalationCompounded", p.EscalationCompounded)
	return nil
}
