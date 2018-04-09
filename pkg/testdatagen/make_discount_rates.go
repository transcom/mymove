package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDiscountRate creates a single DiscountRate.
func MakeDiscountRate(db *pop.Connection, tsp *models.TransportationServiceProvider) (models.DiscountRate, error) {
	if tsp == nil {
		newTSP, err := MakeTSP(db, "Very Good TSP", "NINO")
		if err != nil {
			log.Panic(err)
		}
		tsp = &newTSP
	}

	discountRate := models.DiscountRate{
		PeakRateCycle:            true,
		Origin:                   "US11",
		Destination:              "REGION 4",
		CodeOfService:            "D",
		StandardCarrierAlphaCode: tsp.StandardCarrierAlphaCode,
		LinehaulRate:             40.10,
		SITRate:                  60,
		EffectiveDateLower:       RateEngineEffectiveDateStart,
		EffectiveDateUpper:       RateEngineEffectiveDateEnd,
	}

	verrs, err := db.ValidateAndSave(&discountRate)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return discountRate, err
}
