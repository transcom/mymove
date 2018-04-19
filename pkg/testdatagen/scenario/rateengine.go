package scenario

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Fatalf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		log.Fatalf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

// RunRateEngineScenario1 runs... scenario 1.
func RunRateEngineScenario1(db *pop.Connection) {
	may15_2018 := time.Date(2018, time.May, 15, 0, 0, 0, 0, time.UTC)
	october15_2018 := time.Date(2018, time.October, 15, 0, 0, 0, 0, time.UTC)
	may15_2019 := time.Date(2019, time.May, 15, 0, 0, 0, 0, time.UTC)

	zip3_321 := models.Tariff400ngZip3{
		Zip3:          "321",
		BasepointCity: "Crescent City",
		State:         "FL",
		ServiceArea:   184,
		RateArea:      "ZIP",
		Region:        13,
	}
	mustSave(db, &zip3_321)

	zip5_32168 := models.Tariff400ngZip5RateArea{
		Zip5:     "32168",
		RateArea: "4964400",
	}
	mustSave(db, &zip5_32168)

	zip3_294 := models.Tariff400ngZip3{
		Zip3:          "294",
		BasepointCity: "Moncks Corner",
		State:         "SC",
		ServiceArea:   692,
		RateArea:      "44",
		Region:        12,
	}
	mustSave(db, &zip3_294)

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
		Name: "Standard Moving",
	}
	mustSave(db, &tsp)

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "4964400",
		DestinationRegion: "12",
		CodeOfService:     "2",
	}
	mustSave(db, &tdl)

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Orlando, FL",
		ServiceArea:        184,
		ServicesSchedule:   2,
		LinehaulFactor:     unit.Cents(60),
		ServiceChargeCents: unit.Cents(361),
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &originServiceArea)

	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Charleston, SC",
		ServiceArea:        692,
		ServicesSchedule:   2,
		LinehaulFactor:     unit.Cents(43),
		ServiceChargeCents: unit.Cents(431),
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &destinationServiceArea)

	linehaulRate := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 1000,
		Type:               "ConusLinehaul",
		WeightLbsLower:     4000,
		WeightLbsUpper:     4200,
		RateCents:          unit.Cents(458300),
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &linehaulRate)

	shorthaulRate := models.Tariff400ngShorthaulRate{
		CwtMilesLower:      0,
		CwtMilesUpper:      16001,
		RateCents:          32834,
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &shorthaulRate)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           2,
		WeightLbsLower:     unit.Pound(0),
		WeightLbsUpper:     unit.Pound(16001),
		RateCents:          unit.Cents(6130),
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &fullPackRate)

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           2,
		RateMillicents:     643650,
		EffectiveDateLower: may15_2018,
		EffectiveDateUpper: may15_2019,
	}
	mustSave(db, &fullUnpackRate)

	band := 1
	tspp := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          may15_2018,
		PerformancePeriodEnd:            october15_2018,
		RateCycleStart:                  may15_2018,
		RateCycleEnd:                    october15_2018,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     &band,
		BestValueScore:                  90,
		LinehaulRate:                    0.67,
		SITRate:                         0.5,
	}

	mustSave(db, &tspp)
}
