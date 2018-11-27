package testdatagen

import (
	"github.com/gobuffalo/pop"

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

	mustCreate(db, &serviceArea)

	return serviceArea
}

// MakeDefaultTariff400ngServiceArea makes a ServiceArea with default values
func MakeDefaultTariff400ngServiceArea(db *pop.Connection) models.Tariff400ngServiceArea {
	return MakeTariff400ngServiceArea(db, Assertions{})
}

// GeoModelReturn contains origin and destination zip3 and service area models
type GeoModelReturn struct {
	originZip3        models.Tariff400ngZip3
	destZip3          models.Tariff400ngZip3
	originServiceArea models.Tariff400ngServiceArea
	destServiceArea   models.Tariff400ngServiceArea
}

// MakeTariff400ngGeoModelsForShipment makes zip3 and service area records for a shipment's origin and destination addresses
func MakeTariff400ngGeoModelsForShipment(db *pop.Connection, shipment models.Shipment) GeoModelReturn {
	var result GeoModelReturn

	originAddress := shipment.PickupAddress
	originAssertions := Assertions{}
	originAssertions.Tariff400ngZip3.Zip3 = zip5ToZip3(originAddress.PostalCode)
	result.originZip3 = FetchOrMakeTariff400ngZip3(db, originAssertions)
	originAssertions.Tariff400ngServiceArea.ServiceArea = result.originZip3.ServiceArea
	result.originServiceArea = MakeTariff400ngServiceArea(db, originAssertions)

	destAddress := shipment.Move.Orders.NewDutyStation.Address
	destAssertions := Assertions{}
	destAssertions.Tariff400ngZip3.Zip3 = zip5ToZip3(destAddress.PostalCode)
	result.destZip3 = FetchOrMakeTariff400ngZip3(db, destAssertions)
	destAssertions.Tariff400ngServiceArea.ServiceArea = result.destZip3.ServiceArea
	result.destServiceArea = MakeTariff400ngServiceArea(db, destAssertions)

	return result
}
