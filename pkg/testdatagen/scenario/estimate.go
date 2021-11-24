package scenario

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

// RunPPMSITEstimateScenario1 runs... scenario 1.
func RunPPMSITEstimateScenario1(appCtx appcontext.AppContext) error {
	originZip5_779 := models.Tariff400ngZip5RateArea{
		Zip5:     "77901",
		RateArea: "US68",
	}
	if err := testdatagen.Save(appCtx.DB(), &originZip5_779); err != nil {
		return err
	}

	destZip5_674 := models.Tariff400ngZip5RateArea{
		Zip5:     "67401",
		RateArea: "US58",
	}
	if err := testdatagen.Save(appCtx.DB(), &destZip5_674); err != nil {
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
	if err := testdatagen.Save(appCtx.DB(), &originZip3_779); err != nil {
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
	if err := testdatagen.Save(appCtx.DB(), &destZip3_674); err != nil {
		return err
	}

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	if err := testdatagen.Save(appCtx.DB(), &tsp); err != nil {
		return err
	}

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "US68",
		DestinationRegion: "5",
		CodeOfService:     "D",
	}
	if err := testdatagen.Save(appCtx.DB(), &tdl); err != nil {
		return err
	}

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(39),
		ServiceChargeCents: unit.Cents(350),
		EffectiveDateLower: May15TestYear,
		EffectiveDateUpper: May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	if err := testdatagen.Save(appCtx.DB(), &originServiceArea); err != nil {
		return err
	}

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(43),
		ServiceChargeCents: unit.Cents(350),
		EffectiveDateLower: May15TestYear,
		EffectiveDateUpper: May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	if err := testdatagen.Save(appCtx.DB(), &destServiceArea); err != nil {
		return err
	}

	band := 1
	tspp := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          Oct1TestYear,
		PerformancePeriodEnd:            Dec31TestYear,
		RateCycleStart:                  Oct1TestYear,
		RateCycleEnd:                    May14FollowingYear,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     &band,
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}

	return testdatagen.Save(appCtx.DB(), &tspp)
}
