package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeTariff400ngServiceArea finds or makes a single service_area record
func MakeTariff400ngServiceArea(db *pop.Connection, assertions Assertions) models.Tariff400ngServiceArea {
	serviceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        DefaultServiceArea,
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: PeakRateCycleStart,
		EffectiveDateUpper: NonPeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	mergeModels(&serviceArea, assertions.Tariff400ngServiceArea)

	mustCreate(db, &serviceArea, assertions.Stub)

	return serviceArea
}
