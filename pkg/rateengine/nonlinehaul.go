package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
)

func (re *RateEngine) serviceFeeCents(cwt int, zip3 int) (int, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}
	return cwt * serviceArea.ServiceChargeCents, nil
}

func (re *RateEngine) fullPackCents(cwt int, zip3 int) (int, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}

	fullPackRate, err := models.FetchTariff400ngFullPackRateCents(re.db, cwt, serviceArea.ServicesSchedule)
	if err != nil {
		return 0, err
	}

	return cwt * fullPackRate, nil
}

func (re *RateEngine) fullUnpackCents(cwt int, zip3 int) (int, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}

	fullUnpackRate, err := models.FetchTariff400ngFullUnpackRateMillicents(re.db, serviceArea.ServicesSchedule)
	if err != nil {
		return 0, err
	}

	return cwt * fullUnpackRate / 1000, nil
}

func (re *RateEngine) nonLinehaulChargeTotalCents(weight int, originZip int, destinationZip int) (int, error) {
	cwt := re.determineCWT(weight)
	originServiceFee, err := re.serviceFeeCents(cwt, originZip)
	destinationServiceFee, err := re.serviceFeeCents(cwt, destinationZip)
	pack, err := re.fullPackCents(cwt, originZip)
	unpack, err := re.fullUnpackCents(cwt, destinationZip)
	if err != nil {
		return 0, err
	}
	subTotal := originServiceFee + destinationServiceFee + pack + unpack
	return subTotal, nil
}
