package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticNTSPackPricer struct {
}

// NewDomesticNTSPackPricer creates a new pricer for the domestic NTS pack service
func NewDomesticNTSPackPricer() services.DomesticNTSPackPricer {
	return &domesticNTSPackPricer{}
}

// Price determines the price for a domestic NTS pack service
func (p domesticNTSPackPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, nil, errors.New("RequestedPickupDate is required")
	}
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if servicesScheduleOrigin == 0 {
		return 0, nil, errors.New("Service schedule is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	domOtherPrice, err := fetchDomOtherPrice(appCtx, contractCode, models.ReServiceCodeDPK, servicesScheduleOrigin, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup domestic other price: %w", err)
	}

	var contractYear models.ReContractYear
	err = appCtx.DB().Where("contract_id = $1", domOtherPrice.ContractID).
		Where("$2 between start_date and end_date", requestedPickupDate).
		First(&contractYear)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	// Get NTS packing factor
	shipmentTypePrice, err := fetchShipmentTypePrice(appCtx, contractCode, models.ReServiceCodeDNPK, models.MarketConus)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup shipment type price: %w", err)
	}

	basePrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	factoredPrice := escalatedPrice * shipmentTypePrice.Factor

	totalCost := unit.Cents(math.Round(factoredPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domOtherPrice.PriceCents),
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
		{
			Key:   models.ServiceItemParamNameNTSPackingFactor,
			Value: FormatFloat(shipmentTypePrice.Factor, 2),
		},
	}

	return totalCost, displayParams, nil
}

// PriceUsingParams determines the price for a domestic NTS pack service given PaymentServiceItemParams
func (p domesticNTSPackPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	servicesScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilled), servicesScheduleOrigin)
}
