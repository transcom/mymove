package scenario

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// RunPPMSITEstimateScenario1 runs... scenario 1.
func RunPPMSITEstimateScenario1(db *pop.Connection) error {
	originZip5_779 := models.Tariff400ngZip5RateArea{
		Zip5:     "77901",
		RateArea: "US68",
	}
	if err := save(db, &originZip5_779); err != nil {
		return err
	}

	destZip5_674 := models.Tariff400ngZip5RateArea{
		Zip5:     "67401",
		RateArea: "US58",
	}
	if err := save(db, &destZip5_674); err != nil {
		return err
	}

	originZip3_779 := models.Tariff400ngZip3{
		Zip3:          "779",
		RateArea:      "US68",
		BasepointCity: "Victoria",
		State:         "TX",
		ServiceArea:   "748",
		Region:        "6",
	}
	if err := save(db, &originZip3_779); err != nil {
		return err
	}

	destZip3_674 := models.Tariff400ngZip3{
		Zip3:          "674",
		Region:        "5",
		BasepointCity: "Salina",
		State:         "KS",
		RateArea:      "US58",
		ServiceArea:   "320",
	}
	if err := save(db, &destZip3_674); err != nil {
		return err
	}

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	if err := save(db, &tsp); err != nil {
		return err
	}

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "US68",
		DestinationRegion: "5",
		CodeOfService:     "D",
	}
	if err := save(db, &tdl); err != nil {
		return err
	}

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		LinehaulFactor:     unit.Cents(39),
		ServiceChargeCents: unit.Cents(350),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	if err := save(db, &originServiceArea); err != nil {
		return err
	}

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		LinehaulFactor:     unit.Cents(43),
		ServiceChargeCents: unit.Cents(350),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	if err := save(db, &destServiceArea); err != nil {
		return err
	}

	band := 1
	tspp := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          Oct1_2018,
		PerformancePeriodEnd:            Dec31_2018,
		RateCycleStart:                  Oct1_2018,
		RateCycleEnd:                    May14_2019,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     &band,
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}

	return save(db, &tspp)
}
