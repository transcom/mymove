package rateengine

import (
	"math"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// FeeAndRate holds the rate lookup and calculated fee (non-discounted)
type FeeAndRate struct {
	Fee  unit.Cents
	Rate unit.Millicents
}

// NonLinehaulCostComputation represents the results of a computation.
type NonLinehaulCostComputation struct {
	OriginService      FeeAndRate
	DestinationService FeeAndRate
	Pack               FeeAndRate
	Unpack             FeeAndRate
}

// SITComputation represents the parts of the SIT calculation (needs to be separable to apply correct discount rates).
type SITComputation struct {
	SITPart            unit.Cents
	LinehaulPart       unit.Cents
	NonDiscountedTotal unit.Cents
}

// Scale scales a cost computation by a multiplicative factor
func (c *NonLinehaulCostComputation) Scale(factor float64) {
	c.OriginService.Fee = c.OriginService.Fee.MultiplyFloat64(factor)
	c.DestinationService.Fee = c.DestinationService.Fee.MultiplyFloat64(factor)
	c.Pack.Fee = c.Pack.Fee.MultiplyFloat64(factor)
	c.Unpack.Fee = c.Unpack.Fee.MultiplyFloat64(factor)
}

// ApplyDiscount will apply the linehaul and SIT discounts to the appropriate parts of the SIT computation.
func (s SITComputation) ApplyDiscount(linehaulDiscount unit.DiscountRate, sitDiscount unit.DiscountRate) unit.Cents {
	return sitDiscount.Apply(s.SITPart).AddCents(linehaulDiscount.Apply(s.LinehaulPart))
}

// MarshalLogObject allows SITComputation to be logged by Zap.
func (s SITComputation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("SITPart", s.SITPart.Int())
	encoder.AddInt("LinehaulPart", s.LinehaulPart.Int())
	encoder.AddInt("NonDiscountedTotal", s.NonDiscountedTotal.Int())

	return nil
}

// serviceFeeCents returns the NON-DISCOUNTED rate in millicents with the fee
func (re *RateEngine) serviceFeeCents(cwt unit.CWT, zip3 string, date time.Time) (FeeAndRate, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return FeeAndRate{}, err
	}
	rateCents := serviceArea.ServiceChargeCents
	feeCents := rateCents.Multiply(cwt.Int())
	return FeeAndRate{Fee: feeCents, Rate: rateCents.ToMillicents()}, nil
}

// fullPackCents Returns the NON-DISCOUNTED rate in millicents with the fee
func (re *RateEngine) fullPackCents(cwt unit.CWT, zip3 string, date time.Time) (FeeAndRate, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return FeeAndRate{}, err
	}

	fullPackRate, err := models.FetchTariff400ngFullPackRateCents(re.db, cwt.ToPounds(), serviceArea.ServicesSchedule, date)
	if err != nil {
		return FeeAndRate{}, err
	}

	return FeeAndRate{Fee: fullPackRate.Multiply(cwt.Int()), Rate: fullPackRate.ToMillicents()}, nil
}

// fullUnpackCents Returns the NON-DISCOUNTED rate in millicents with the fee
func (re *RateEngine) fullUnpackCents(cwt unit.CWT, zip3 string, date time.Time) (FeeAndRate, error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return FeeAndRate{}, err
	}

	fullUnpackRate, err := models.FetchTariff400ngFullUnpackRateMillicents(re.db, serviceArea.ServicesSchedule, date)
	if err != nil {
		return FeeAndRate{}, err
	}

	return FeeAndRate{Fee: unit.Cents(math.Round(float64(cwt.Int()*fullUnpackRate) / 1000.0)), Rate: unit.Millicents(fullUnpackRate)}, nil
}

// SitCharge calculates the SIT charge based on various factors.
// Note: Assumes the caller will apply any SIT discount rate as needed (no discounts applied here).
func (re *RateEngine) SitCharge(cwt unit.CWT, daysInSIT int, zip3 string, date time.Time, isPPM bool) (SITComputation, error) {
	if daysInSIT == 0 {
		return SITComputation{}, nil
	} else if daysInSIT < 0 {
		return SITComputation{}, errors.New("requested SitCharge for negative days in SIT")
	}

	effectiveCWT := cwt
	if !isPPM {
		// An HHG uses a minimum weight of 1000 pounds.
		// TODO: If an HHG shipment is delivered partially out of SIT (split deliveries), 1000 lb min does not apply.
		minCWT := unit.Pound(1000).ToCWT()
		if cwt < minCWT {
			effectiveCWT = minCWT
		}
	}

	sa, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3, date)
	if err != nil {
		return SITComputation{}, err
	}

	// Both PPMs and HHGs use 185A and 185B in the same way.
	sitPart := sa.SIT185ARateCents.Multiply(effectiveCWT.Int())
	additionalDays := daysInSIT - 1
	if additionalDays > 0 {
		sitPart = sitPart.AddCents(sa.SIT185BRateCents.Multiply(additionalDays).Multiply(effectiveCWT.Int()))
	}

	zapFields := []zap.Field{
		zap.Int("cwt", cwt.Int()),
		zap.Int("days", daysInSIT),
		zap.String("zip3", zip3),
		zap.Time("date", date),
		zap.Bool("isPPM", isPPM),
		zap.Int("effectiveCWT", effectiveCWT.Int()),
		zap.Int("servicesSchedule", sa.ServicesSchedule),
		zap.Int("sitPDSchedule", sa.SITPDSchedule),
		zap.Int("185A", sa.SIT185ARateCents.Int()),
		zap.Int("185B", sa.SIT185BRateCents.Int()),
	}

	var linehaulPart unit.Cents
	if isPPM {
		// PPM SIT formula:
		//   (185A SIT first day rate * CWT) +
		//   (185B SIT additional day rate * additional days * CWT) +
		//   210A SIT PD 30 miles or less for SIT PD schedule of service area +
		//   225A SIT PD Self/Mini Storage for services schedule of service area
		rate210A, err := models.FetchTariff400ngItemRate(re.db, "210A", sa.SITPDSchedule, effectiveCWT.ToPounds(), date)
		if err != nil {
			return SITComputation{}, errors.Wrapf(err, "No 210A rate found for schedule %v, %v pounds, date %v", sa.SITPDSchedule, effectiveCWT.ToPounds(), date)
		}
		sitPart = sitPart.AddCents(rate210A.RateCents)

		rate225A, err := models.FetchTariff400ngItemRate(re.db, "225A", sa.ServicesSchedule, effectiveCWT.ToPounds(), date)
		if err != nil {
			return SITComputation{}, errors.Wrapf(err, "No 225A rate found for schedule %v, %v pounds, date %v", sa.ServicesSchedule, effectiveCWT.ToPounds(), date)
		}
		linehaulPart = rate225A.RateCents

		zapFields = append(zapFields,
			zap.Int("210A", rate210A.RateCents.Int()),
			zap.Int("225A", rate225A.RateCents.Int()))
	} else {
		// Just return 185A and 185B parts of HHG for now.  Full implementation in later story.

		// TODO: The rest of the HHG scenarios are as follows (to be added to the 185A and 185B parts):
		//   * 30 miles or less from original delivery address to final delivery address (block 18 on GBL):
		//       (185A SIT first day rate * CWT) +
		//       (185B SIT additional day rate * additional days * CWT)
		//       210A SIT PD 30 miles or less for SIT PD schedule of service area
		//   * Between 31 and 50 miles from original delivery address to final delivery address (block 18 on GBL):
		//       (185A SIT first day rate * CWT) +
		//	     (185B SIT additional day rate * additional days * CWT)
		//       210A SIT PD 30 miles or less for SIT PD schedule of service area +
		//       210B SIT PD 30 to 50 miles SIT PD schedule of service area
		//   * Over 50 miles from original delivery address to final delivery address (block 18 on GBL):
		//       (185A SIT first day rate * CWT) +
		//	     (185B SIT additional day rate * additional days * CWT)
		//       210C SIT PD over 50 miles SIT PD schedule of service area
	}

	sitComputation := SITComputation{
		SITPart:            sitPart,
		LinehaulPart:       linehaulPart,
		NonDiscountedTotal: sitPart.AddCents(linehaulPart),
	}

	zapFields = append(zapFields, zap.Object("sit computation", sitComputation))
	re.logger.Info("sit calculation", zapFields...)

	return sitComputation, err
}

func (re *RateEngine) nonLinehaulChargeComputation(weight unit.Pound, originZip5 string, destinationZip5 string, date time.Time) (cost NonLinehaulCostComputation, err error) {
	cwt := weight.ToCWT()
	originZip3 := Zip5ToZip3(originZip5)
	destinationZip3 := Zip5ToZip3(destinationZip5)
	cost.OriginService, err = re.serviceFeeCents(cwt, originZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine origin service fee")
	}
	cost.DestinationService, err = re.serviceFeeCents(cwt, destinationZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine destination service fee")
	}
	cost.Pack, err = re.fullPackCents(cwt, originZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine full pack cost")
	}
	cost.Unpack, err = re.fullUnpackCents(cwt, destinationZip3, date)
	if err != nil {
		return cost, errors.Wrap(err, "Failed to  determine full unpack cost")
	}

	re.logger.Info("Non-Linehaul charge total calculated",
		zap.Int("origin service fee", cost.OriginService.Fee.Int()),
		zap.Int("destination service fee", cost.DestinationService.Fee.Int()),
		zap.Int("pack fee", cost.Pack.Fee.Int()),
		zap.Int("unpack fee", cost.Unpack.Fee.Int()))

	return cost, nil
}
