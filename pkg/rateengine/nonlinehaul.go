package rateengine

import (
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func (re *RateEngine) serviceFeeCents(cwt int, zip3 string) (int, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}
	return cwt * serviceArea.ServiceChargeCents, nil
}

func (re *RateEngine) fullPackCents(cwt int, zip3 string) (int, error) {
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

func (re *RateEngine) fullUnpackCents(cwt int, zip3 string) (int, error) {
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

func (re *RateEngine) nonLinehaulChargeTotalCents(weight int, originZip string, destinationZip string) (int, error) {
	cwt := re.determineCWT(weight)
	originServiceFee, err := re.serviceFeeCents(cwt, originZip)
	if err != nil {
		return 0, err
	}
	destinationServiceFee, err := re.serviceFeeCents(cwt, destinationZip)
	if err != nil {
		return 0, err
	}
	pack, err := re.fullPackCents(cwt, originZip)
	if err != nil {
		return 0, err
	}
	unpack, err := re.fullUnpackCents(cwt, destinationZip)
	if err != nil {
		return 0, err
	}
	subTotal := originServiceFee + destinationServiceFee + pack + unpack

	re.logger.Info("Non-Linehaul charge total calculated",
		zap.Int("origin service fee", originServiceFee),
		zap.Int("destination service fee", destinationServiceFee),
		zap.Int("pack fee", pack),
		zap.Int("unpack fee", unpack))

	return subTotal, nil
}
