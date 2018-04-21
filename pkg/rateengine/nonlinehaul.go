package rateengine

import (
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (re *RateEngine) serviceFeeCents(cwt unit.CWT, zip3 string, date time.Time) (unit.Cents, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}
	return serviceArea.ServiceChargeCents.Multiply(cwt.Int()), nil
}

func (re *RateEngine) fullPackCents(cwt unit.CWT, zip3 string, date time.Time) (unit.Cents, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}

	fullPackRate, err := models.FetchTariff400ngFullPackRateCents(re.db, cwt.ToPounds(), serviceArea.ServicesSchedule, date)
	if err != nil {
		return 0, err
	}

	return fullPackRate.Multiply(cwt.Int()), nil
}

func (re *RateEngine) fullUnpackCents(cwt unit.CWT, zip3 string, date time.Time) (unit.Cents, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}

	fullUnpackRate, err := models.FetchTariff400ngFullUnpackRateMillicents(re.db, serviceArea.ServicesSchedule, date)
	if err != nil {
		return 0, err
	}

	return unit.Cents(math.Round(float64(cwt.Int()*fullUnpackRate) / 1000.0)), nil
}

// sitCharge calculates the SIT charge based on various factors.
// If `isPPM` (Personally Procured Move) is True we do not apply the first-day
// storage fees, 185A, to the total.
func (re *RateEngine) sitCharge(cwt unit.CWT, daysInSIT int, zip3 string, date time.Time, isPPM bool) (unit.Cents, error) {
	if daysInSIT <= 0 {
		re.logger.Info("requested sitCharge for zero or less days in SIT?")
		return 0, nil
	}

	sa, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}

	var sitTotal unit.Cents

	if isPPM {
		sitTotal = unit.Cents(cwt.Int() * sa.SIT185BRateCents.Int() * daysInSIT)
	} else {
		sitTotal = unit.Cents(cwt.Int() * sa.SIT185ARateCents.Int())
		daysInSIT--
		if daysInSIT > 0 {
			sitTotal += unit.Cents(cwt.Int() * sa.SIT185BRateCents.Int() * daysInSIT)
		}
	}
	re.logger.Info("sit calculation", zap.Int("cwt", cwt.Int()), zap.Int("185B", sa.SIT185BRateCents.Int()), zap.Int("days", daysInSIT), zap.Int("total", sitTotal.Int()))

	return sitTotal, err
}

func (re *RateEngine) nonLinehaulChargeTotalCents(weight unit.Pound, originZip5 string,
	destinationZip5 string, date time.Time, daysInSIT int, isPPM bool) (unit.Cents, error) {

	originZip3, destinationZip3 := re.zip5ToZip3(originZip5, destinationZip5)
	originServiceFee, err := re.serviceFeeCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		return 0, err
	}
	destinationServiceFee, err := re.serviceFeeCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		return 0, err
	}
	pack, err := re.fullPackCents(weight.ToCWT(), originZip3, date)
	if err != nil {
		return 0, err
	}
	unpack, err := re.fullUnpackCents(weight.ToCWT(), destinationZip3, date)
	if err != nil {
		return 0, err
	}
	sit, err := re.sitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, isPPM)
	if err != nil {
		return 0, err
	}
	subTotal := originServiceFee + destinationServiceFee + pack + unpack + sit

	re.logger.Info("Non-Linehaul charge total calculated",
		zap.Int("origin service fee", originServiceFee.Int()),
		zap.Int("destination service fee", destinationServiceFee.Int()),
		zap.Int("pack fee", pack.Int()),
		zap.Int("unpack fee", unpack.Int()),
		zap.Int("SIT fee", sit.Int()))

	return subTotal, nil
}
