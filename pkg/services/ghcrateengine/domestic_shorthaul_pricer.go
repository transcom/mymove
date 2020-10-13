package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DomesticShorthaulPricer is a service object to price domestic shorthaul
type domesticShorthaulPricer struct {
	db *pop.Connection
	//logger       Logger
}

// NewDomesticShorthaulPricer is the public constructor for a DomesticRateAreaPricer using Pop
func NewDomesticShorthaulPricer(db *pop.Connection) services.DomesticShorthaulPricer {
	return &domesticShorthaulPricer{
		db: db,
		//logger:       logger,
	}
}

// Price determines the price for a counseling service
func (p domesticShorthaulPricer) Price(contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (totalCost unit.Cents, err error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, errors.New("RequestedPickupDate is required")
	}
	if weight < minDomesticWeight {
		return 0, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if distance <= 0 {
		return 0, errors.New("Distance must be greater than 0")
	}
	if len(serviceArea) == 0 {
		return 0, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	// look up rate for shorthaul
	shorthaulServiceCode := "DSH"
	var contractYear models.ReContractYear
	var domServiceAreaPrice models.ReDomesticServiceAreaPrice
	err = p.db.Q().
		Join("re_domestic_service_areas sa", "domestic_service_area_id = sa.id").
		Join("re_services", "service_id = re_services.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_area_prices.contract_id").
		Where("sa.service_area = $1", serviceArea).
		Where("re_services.code = $2", shorthaulServiceCode).
		Where("re_contracts.code = $3", contractCode).
		Where("is_peak_period = $4", isPeakPeriod).
		First(&domServiceAreaPrice)
	if err != nil {
		return 0, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}
	err = p.db.Where("contract_id = $1", domServiceAreaPrice.ContractID).
		Where("$2 between start_date and end_date", requestedPickupDate).
		First(&contractYear)
	if err != nil {
		return 0, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64() * distance.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost = unit.Cents(math.Round(escalatedPrice))

	// To be fixed under this story: https://dp3.atlassian.net/browse/MB-2352
	// unable to get logger to pass in for instantiation at the pkg/handler/primeapi/api.go
	//p.logger.Info(fmt.Sprintf("%s calculated", shorthaulServiceCode), // May change to use ServiceName
	//zap.String("contractCode:", contractCode),
	//zap.String("serviceCode:", shorthaulServiceCode),
	//zap.Time("requestedPickupDate:", requestedPickupDate),
	//zap.String("serviceArea:", serviceArea),
	//zap.Int("distance (mi):", distance.Int()),
	//zap.Int("weight (lb):", weight.Int()),
	//zap.Int("effectiveWeight (lb):", effectiveWeight.Int()),
	//zap.Bool("isPeakPeriod: ", isPeakPeriod),
	//zap.Int("Dom. Service Area PriceCents: ", domServiceAreaPrice.PriceCents),
	//zap.Int("Contract Year Escalation: ", contractYear.EscalationCompounded),
	//zap.Float64("baseCost (cents):", basePrice),
	//zap.Float64("escalatedCost (cents):", escalatedPrice),
	//zap.Int("totalCost (cents):", totalCost.Int()),
	//)
	return totalCost, err
}

func (p domesticShorthaulPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), err
	}

	distanceZip5, err := getParamInt(params, models.ServiceItemParamNameDistanceZip5)
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

	total, err := p.Price(contractCode, requestedPickupDate, unit.Miles(distanceZip5), unit.Pound(weightBilledActual), serviceAreaOrigin)
	return total, err
}
