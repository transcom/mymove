package rateengine

import (
	"math"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// NonLinehaulCostComputation represents the results of a computation.
type NonLinehaulCostComputation struct {
	OriginServiceFee      unit.Cents
	DestinationServiceFee unit.Cents
	PackFee               unit.Cents
	UnpackFee             unit.Cents
}

// Scale scales a cost computation by a multiplicative factor
func (c *NonLinehaulCostComputation) Scale(factor float64) {
	c.OriginServiceFee = c.OriginServiceFee.MultiplyFloat64(factor)
	c.DestinationServiceFee = c.DestinationServiceFee.MultiplyFloat64(factor)
	c.PackFee = c.PackFee.MultiplyFloat64(factor)
	c.UnpackFee = c.UnpackFee.MultiplyFloat64(factor)
}

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

// SitCharge calculates the SIT charge based on various factors.
// If `isPPM` (Personally Procured Move) is True we do not apply the first-day
// storage fees, 185A, to the total.
func (re *RateEngine) SitCharge(cwt unit.CWT, daysInSIT int, zip3 string, date time.Time, isPPM bool) (unit.Cents, error) {
	if daysInSIT == 0 {
		return 0, nil
	} else if daysInSIT < 0 {
		return 0, errors.New("requested SitCharge for negative days in SIT")
	}

	sa, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return 0, err
	}

	var sitTotal unit.Cents

	if isPPM {
		sitTotal = sa.SIT185BRateCents.Multiply(daysInSIT).Multiply(cwt.Int())
	} else {
		sitTotal = sa.SIT185ARateCents.Multiply(cwt.Int())
		additionalDays := daysInSIT - 1
		if additionalDays > 0 {
			sitTotal = sitTotal.AddCents(sa.SIT185BRateCents.Multiply(additionalDays).Multiply(cwt.Int()))
		}
	}
	re.logger.Info("sit calculation",
		zap.Int("cwt", cwt.Int()),
		zap.Int("185A", sa.SIT185ARateCents.Int()),
		zap.Int("185B", sa.SIT185BRateCents.Int()),
		zap.Int("days", daysInSIT),
		zap.Int("total", sitTotal.Int()))

	return sitTotal, err
}

func (re *RateEngine) nonLinehaulChargeComputation(weight unit.Pound, originZip5 string, destinationZip5 string, date time.Time) (cost NonLinehaulCostComputation, err error) {
	cwt := weight.ToCWT()
	originZip3 := Zip5ToZip3(originZip5)
	destinationZip3 := Zip5ToZip3(destinationZip5)
	cost.OriginServiceFee, err = re.serviceFeeCents(cwt, originZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine origin service fee")
	}
	cost.DestinationServiceFee, err = re.serviceFeeCents(cwt, destinationZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine destination service fee")
	}
	cost.PackFee, err = re.fullPackCents(cwt, originZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine full pack cost")
	}
	cost.UnpackFee, err = re.fullUnpackCents(cwt, destinationZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine full unpack cost")
	}

	re.logger.Info("Non-Linehaul charge total calculated",
		zap.Int("origin service fee", cost.OriginServiceFee.Int()),
		zap.Int("destination service fee", cost.DestinationServiceFee.Int()),
		zap.Int("pack fee", cost.PackFee.Int()),
		zap.Int("unpack fee", cost.UnpackFee.Int()))

	return cost, nil
}
