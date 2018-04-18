package scenario

import (
	"log"

	"github.com/gobuffalo/pop"
	. "github.com/transcom/mymove/pkg/models"
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

func RunRateEngineScenario1(db *pop.Connection) {
	zip3_321 := Tariff400ngZip3{
		Zip3:          "321",
		BasepointCity: "Crescent City",
		State:         "FL",
		ServiceArea:   184,
		RateArea:      "ZIP",
		Region:        13,
	}
	mustSave(db, &zip3_321)

	zip5_32168 := Tariff400ngZip5RateArea{
		Zip5:     "32168",
		RateArea: "4964400",
	}
	mustSave(db, &zip5_32168)

	zip3_294 := Tariff400ngZip3{
		Zip3:          "294",
		BasepointCity: "Moncks Corner",
		State:         "SC",
		ServiceArea:   692,
		RateArea:      "44",
		Region:        12,
	}
	mustSave(db, &zip3_294)

	tsp := TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
		Name: "Standard Moving",
	}
	mustSave(db, &tsp)

	tdl := TrafficDistributionList{
		SourceRateArea:    "4964400",
		DestinationRegion: "12",
		CodeOfService:     "2",
	}
	mustSave(db, &tdl)
}
