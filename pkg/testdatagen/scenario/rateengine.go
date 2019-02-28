package scenario

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// RunRateEngineScenario1 runs... scenario 1.
func RunRateEngineScenario1(db *pop.Connection) error {
	zip3_321 := models.Tariff400ngZip3{
		Zip3:          "321",
		BasepointCity: "Crescent City",
		State:         "FL",
		ServiceArea:   "184",
		RateArea:      "ZIP",
		Region:        "13",
	}
	if err := save(db, &zip3_321); err != nil {
		return err
	}

	zip5_32168 := models.Tariff400ngZip5RateArea{
		Zip5:     "32168",
		RateArea: "US4964400",
	}
	if err := save(db, &zip5_32168); err != nil {
		return err
	}

	zip3_294 := models.Tariff400ngZip3{
		Zip3:          "294",
		BasepointCity: "Moncks Corner",
		State:         "SC",
		ServiceArea:   "692",
		RateArea:      "US44",
		Region:        "12",
	}
	if err := save(db, &zip3_294); err != nil {
		return err
	}

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	if err := save(db, &tsp); err != nil {
		return err
	}

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "US4964400",
		DestinationRegion: "12",
		CodeOfService:     "2",
	}
	if err := save(db, &tdl); err != nil {
		return err
	}

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Orlando, FL",
		ServiceArea:        "184",
		ServicesSchedule:   2,
		LinehaulFactor:     unit.Cents(60),
		ServiceChargeCents: unit.Cents(361),
		SIT185ARateCents:   unit.Cents(1691),
		SIT185BRateCents:   unit.Cents(65),
		SITPDSchedule:      3,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &originServiceArea); err != nil {
		return err
	}

	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Charleston, SC",
		ServiceArea:        "692",
		ServicesSchedule:   2,
		LinehaulFactor:     unit.Cents(43),
		ServiceChargeCents: unit.Cents(431),
		SIT185ARateCents:   unit.Cents(1378),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      2,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &destinationServiceArea); err != nil {
		return err
	}

	linehaulRate := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 1000,
		Type:               "ConusLinehaul",
		WeightLbsLower:     4000,
		WeightLbsUpper:     4200,
		RateCents:          unit.Cents(458300),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &linehaulRate); err != nil {
		return err
	}

	itemRate210A := models.Tariff400ngItemRate{
		Code:               "210A",
		Schedule:           &destinationServiceArea.SITPDSchedule,
		WeightLbsLower:     linehaulRate.WeightLbsLower,
		WeightLbsUpper:     linehaulRate.WeightLbsUpper,
		RateCents:          unit.Cents(57600),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &itemRate210A); err != nil {
		return err
	}

	itemRate225A := models.Tariff400ngItemRate{
		Code:               "225A",
		Schedule:           &destinationServiceArea.ServicesSchedule,
		WeightLbsLower:     linehaulRate.WeightLbsLower,
		WeightLbsUpper:     linehaulRate.WeightLbsUpper,
		RateCents:          unit.Cents(9900),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &itemRate225A); err != nil {
		return err
	}

	shorthaulRate := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      0,
		CwtMilesUpper:      16001,
		RateCents:          32834,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &shorthaulRate); err != nil {
		return err
	}

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           2,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(16001),
		RateCents:          unit.Cents(6130),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &fullPackRate); err != nil {
		return err
	}

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           2,
		RateMillicents:     643650,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &fullUnpackRate); err != nil {
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
		LinehaulRate:                    0.67,
		SITRate:                         0.5,
	}

	return save(db, &tspp)
}

// RunRateEngineScenario2 runs... scenario 2.
func RunRateEngineScenario2(db *pop.Connection) error {
	zip3_945 := models.Tariff400ngZip3{
		Zip3:          "945",
		BasepointCity: "Walnut Creek",
		State:         "CA",
		ServiceArea:   "80",
		RateArea:      "US87",
		Region:        "2",
	}
	if err := save(db, &zip3_945); err != nil {
		return err
	}

	zip3_786 := models.Tariff400ngZip3{
		Zip3:          "786",
		BasepointCity: "Austin",
		State:         "TX",
		ServiceArea:   "744",
		RateArea:      "ZIP",
		Region:        "6",
	}
	if err := save(db, &zip3_786); err != nil {
		return err
	}

	zip5_78626 := models.Tariff400ngZip5RateArea{
		Zip5:     "78626",
		RateArea: "US68",
	}
	if err := save(db, &zip5_78626); err != nil {
		return err
	}

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	if err := save(db, &tsp); err != nil {
		return err
	}

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "US87",
		DestinationRegion: "6",
		CodeOfService:     "2",
	}
	if err := save(db, &tdl); err != nil {
		return err
	}

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "San Francisco, CA",
		ServiceArea:        "80",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(263),
		ServiceChargeCents: unit.Cents(489),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
		SIT185ARateCents:   unit.Cents(1447),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      3,
	}
	if err := save(db, &originServiceArea); err != nil {
		return err
	}

	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Austin, TX",
		ServiceArea:        "744",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(78),
		ServiceChargeCents: unit.Cents(452),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
		SIT185ARateCents:   unit.Cents(1642),
		SIT185BRateCents:   unit.Cents(70),
		SITPDSchedule:      3,
	}
	if err := save(db, &destinationServiceArea); err != nil {
		return err
	}

	linehaulRate1 := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1601,
		DistanceMilesUpper: 1801,
		Type:               "ConusLinehaul",
		WeightLbsLower:     7400,
		WeightLbsUpper:     7600,
		RateCents:          unit.Cents(1277900),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &linehaulRate1); err != nil {
		return err
	}

	item1Rate210A := models.Tariff400ngItemRate{
		Code:               "210A",
		Schedule:           &destinationServiceArea.SITPDSchedule,
		WeightLbsLower:     linehaulRate1.WeightLbsLower,
		WeightLbsUpper:     linehaulRate1.WeightLbsUpper,
		RateCents:          unit.Cents(57600),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &item1Rate210A); err != nil {
		return err
	}

	item1Rate225A := models.Tariff400ngItemRate{
		Code:               "225A",
		Schedule:           &destinationServiceArea.ServicesSchedule,
		WeightLbsLower:     linehaulRate1.WeightLbsLower,
		WeightLbsUpper:     linehaulRate1.WeightLbsUpper,
		RateCents:          unit.Cents(9900),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &item1Rate225A); err != nil {
		return err
	}

	linehaulRate2 := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1601,
		DistanceMilesUpper: 1701,
		Type:               "ConusLinehaul",
		WeightLbsLower:     1000,
		WeightLbsUpper:     1400,
		RateCents:          unit.Cents(1277900),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &linehaulRate2); err != nil {
		return err
	}

	item2Rate210A := models.Tariff400ngItemRate{
		Code:               "210A",
		Schedule:           &destinationServiceArea.SITPDSchedule,
		WeightLbsLower:     linehaulRate2.WeightLbsLower,
		WeightLbsUpper:     linehaulRate2.WeightLbsUpper,
		RateCents:          unit.Cents(35050),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &item2Rate210A); err != nil {
		return err
	}

	item2Rate225A := models.Tariff400ngItemRate{
		Code:               "225A",
		Schedule:           &destinationServiceArea.ServicesSchedule,
		WeightLbsLower:     linehaulRate2.WeightLbsLower,
		WeightLbsUpper:     linehaulRate2.WeightLbsUpper,
		RateCents:          unit.Cents(7700),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &item2Rate225A); err != nil {
		return err
	}

	shorthaulRate := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      96001,
		CwtMilesUpper:      128001,
		RateCents:          18242,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &shorthaulRate); err != nil {
		return err
	}

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           3,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(16001),
		RateCents:          unit.Cents(6714),
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &fullPackRate); err != nil {
		return err
	}

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           3,
		RateMillicents:     704970,
		EffectiveDateLower: May15_2018,
		EffectiveDateUpper: May14_2019,
	}
	if err := save(db, &fullUnpackRate); err != nil {
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
		LinehaulRate:                    0.67,
		SITRate:                         0.6,
	}

	return save(db, &tspp)
}
