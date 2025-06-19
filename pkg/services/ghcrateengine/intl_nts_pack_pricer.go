package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type intlNTSHHGPackPricer struct {
	// INPK requires IHPK for base pricing
	// This is because iHHG -> iNTS requires
	// HHG packing, but just with the NTS market factor
	basePricer services.IntlHHGPackPricer
}

func NewIntlNTSHHGPackPricer(basePricer services.IntlHHGPackPricer) services.IntlNTSHHGPackPricer {
	return &intlNTSHHGPackPricer{
		basePricer,
	}
}

// INPK pricing uses the base IHPK logic and applies the market factor afterwards
func (p *intlNTSHHGPackPricer) Price(
	appCtx appcontext.AppContext,
	contractCode string,
	referenceDate time.Time,
	weight unit.Pound,
	perUnitCents int,
) (unit.Cents, services.PricingDisplayParams, error) {
	// While we could just call `priceIntlPackUnpack` like how
	// IHPK does, that is not future proof. We must rely on IHPK base price
	basePrice, displayParams, err := p.basePricer.Price(appCtx, contractCode, referenceDate, weight, perUnitCents)
	if err != nil {
		return 0, nil, err
	}

	// Now we get the info needed for the INPK market factor

	contract, err := fetchContractByContractCode(appCtx, contractCode)
	if err != nil {
		return 0, nil, fmt.Errorf("could not find contract with code: %s: %w", contractCode, err)
	}

	inpk, err := models.FetchReServiceByCode(appCtx.DB(), models.ReServiceCodeINPK)
	if err != nil {
		return 0, nil, err
	}

	// Now we get the factor itself
	// Params have not been created yet so we need to find and append the market factor param
	factor, err := models.FetchMarketFactor(appCtx, contract.ID, inpk.ID, models.MarketOconus.String())
	if err != nil {
		return 0, nil, err
	}

	// Now we multiply the IHPK base price by the NTS factor
	finalPrice := unit.Cents(math.Round(float64(basePrice) * factor))

	// Append the factor to the params
	factorParam := services.PricingDisplayParam{
		Key:   models.ServiceItemParamNameNTSPackingFactor,
		Value: FormatFloat(factor, -1),
	}
	displayParams = append(displayParams, factorParam)

	return finalPrice, displayParams, nil
}

// PriceUsingParams calls through to the existing IHPK PriceUsingParams logic, then
// multiplies the result by the INPK multiplier.
func (p *intlNTSHHGPackPricer) PriceUsingParams(
	appCtx appcontext.AppContext,
	params models.PaymentServiceItemParams,
) (unit.Cents, services.PricingDisplayParams, error) {
	basePrice, displayParams, err := p.basePricer.PriceUsingParams(appCtx, params)
	if err != nil {
		return 0, nil, err
	}

	factor, err := getParamFloat(params, models.ServiceItemParamNameNTSPackingFactor)
	if err != nil {
		return 0, nil, err
	}

	// Now we multiply the IHPK base price by the NTS factor
	finalPrice := unit.Cents(math.Round(float64(basePrice) * factor))

	return finalPrice, displayParams, nil
}
