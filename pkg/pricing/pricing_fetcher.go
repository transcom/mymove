package pricing

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/unit"
)

func FetchServiceItemPrice(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, mtoShipment models.MTOShipment, planner route.Planner) (unit.Cents, error) {
	// No reason to attempt to calculate a price without an estimated weight.
	min_weight := unit.Pound(0)
	if mtoShipment.PrimeEstimatedWeight == nil || mtoShipment.PrimeEstimatedWeight == &min_weight {
		return 0, nil
	}
	isPPM := (mtoShipment.ShipmentType == models.MTOShipmentTypePPM)

	requestedPickupDate := *mtoShipment.RequestedPickupDate
	currTime := time.Now()
	primeEstimatedWeight := mtoShipment.PrimeEstimatedWeight

	if mtoShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
		newWeight := int(primeEstimatedWeight.Float64() * 1.1)
		primeEstimatedWeight = (*unit.Pound)(&newWeight)
	}

	contractCode, err := FetchContractCode(appCtx, currTime)
	if err != nil {
		contractCode, err = FetchContractCode(appCtx, requestedPickupDate)
		if err != nil {
			return 0, err
		}
	}
	calculatedDistance := 0
	if isServiceItemCodeDistanceBased(serviceItem.ReService.Code) {
		if mtoShipment.PickupAddress != nil && mtoShipment.DestinationAddress != nil && planner != nil {
			calculatedDistance, err = planner.ZipTransitDistance(appCtx, mtoShipment.PickupAddress.PostalCode, mtoShipment.DestinationAddress.PostalCode)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("invalid address or missing planner provided")
		}
	}

	var serviceAreaPostalCode string
	if UseDestinationAddressServiceArea(serviceItem.ReService.Code) {
		serviceAreaPostalCode = mtoShipment.DestinationAddress.PostalCode
	} else {
		serviceAreaPostalCode = mtoShipment.PickupAddress.PostalCode
	}

	domesticServiceArea, err := fetchDomesticServiceArea(appCtx, contractCode, serviceAreaPostalCode)
	if err != nil {
		return 0, err
	}

	var price unit.Cents
	switch serviceItem.ReService.Code {
	case models.ReServiceCodeDOP:
		originPricer := ghcrateengine.NewDomesticOriginPricer()

		price, _, err = originPricer.Price(appCtx, contractCode, requestedPickupDate, *primeEstimatedWeight, domesticServiceArea.ServiceArea, isPPM)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeDPK:
		packPricer := ghcrateengine.NewDomesticPackPricer()

		servicesScheduleOrigin := domesticServiceArea.ServicesSchedule

		price, _, err = packPricer.Price(appCtx, contractCode, requestedPickupDate, *primeEstimatedWeight, servicesScheduleOrigin, isPPM)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeDDP:
		destinationPricer := ghcrateengine.NewDomesticDestinationPricer()

		price, _, err = destinationPricer.Price(appCtx, contractCode, requestedPickupDate, *primeEstimatedWeight, domesticServiceArea.ServiceArea, isPPM)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeDUPK:
		unpackPricer := ghcrateengine.NewDomesticUnpackPricer()

		serviceScheduleDestination := domesticServiceArea.ServicesSchedule

		price, _, err = unpackPricer.Price(appCtx, contractCode, requestedPickupDate, *primeEstimatedWeight, serviceScheduleDestination, isPPM)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeDLH:
		linehaulPricer := ghcrateengine.NewDomesticLinehaulPricer()

		price, _, err = linehaulPricer.Price(appCtx, contractCode, requestedPickupDate, unit.Miles(calculatedDistance), *primeEstimatedWeight, domesticServiceArea.ServiceArea, isPPM)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeDSH:
		shorthaulPricer := ghcrateengine.NewDomesticShorthaulPricer()

		price, _, err = shorthaulPricer.Price(appCtx, contractCode, requestedPickupDate, unit.Miles(calculatedDistance), *primeEstimatedWeight, domesticServiceArea.ServiceArea)
		if err != nil {
			return 0, err
		}
	case models.ReServiceCodeFSC:
		fuelSurchargePricer := ghcrateengine.NewFuelSurchargePricer()
		var pickupDateForFSC time.Time

		// actual pickup date likely won't exist at the time of service item creation, but it could
		// use requested pickup date if no actual date exists
		if mtoShipment.ActualPickupDate != nil {
			pickupDateForFSC = *mtoShipment.ActualPickupDate
		} else {
			pickupDateForFSC = requestedPickupDate
		}

		fscWeightBasedDistanceMultiplier := LookupFSCWeightBasedDistanceMultiplier(appCtx, *primeEstimatedWeight)

		fscWeightBasedDistanceMultiplierFloat, err := strconv.ParseFloat(fscWeightBasedDistanceMultiplier, 64)
		if err != nil {
			return 0, err
		}
		eiaFuelPrice, err := LookupEIAFuelPrice(appCtx, pickupDateForFSC)
		if err != nil {
			return 0, err
		}
		price, _, err = fuelSurchargePricer.Price(appCtx, pickupDateForFSC, unit.Miles(calculatedDistance), *primeEstimatedWeight, fscWeightBasedDistanceMultiplierFloat, eiaFuelPrice, isPPM)
		if err != nil {
			return 0, err
		}
	default:
		// this is an invalid code, return an error
		return 0, apperror.NewNotFoundError(uuid.Nil, "invalid service item code was provided")
	}
	return price, nil
}

func FetchContractCode(appCtx appcontext.AppContext, date time.Time) (string, error) {
	var contractYear models.ReContractYear
	err := appCtx.DB().EagerPreload("Contract").Where("? between start_date and end_date", date).
		First(&contractYear)
	if err != nil || contractYear.Contract.Code == uuid.Nil.String() {
		return "", apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("no contract year found for %s", date.String()))
	}

	contract := contractYear.Contract

	contractCode := contract.Code
	return contractCode, nil
}

func isServiceItemCodeDistanceBased(serviceCode models.ReServiceCode) bool {
	return (serviceCode == models.ReServiceCodeFSC || serviceCode == models.ReServiceCodeDSH || serviceCode == models.ReServiceCodeDLH)
}

func UseDestinationAddressServiceArea(serviceCode models.ReServiceCode) bool {
	return (serviceCode == models.ReServiceCodeDDP || serviceCode == models.ReServiceCodeDUPK)
}

func fetchDomesticServiceArea(appCtx appcontext.AppContext, contractCode string, shipmentPostalCode string) (models.ReDomesticServiceArea, error) {
	// find the service area by querying for the service area associated with the zip3
	zip := shipmentPostalCode
	zip3 := zip[0:3]
	var domesticServiceArea models.ReDomesticServiceArea
	err := appCtx.DB().Q().
		Join("re_zip3s", "re_zip3s.domestic_service_area_id = re_domestic_service_areas.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_areas.contract_id").
		Where("re_zip3s.zip3 = ?", zip3).
		Where("re_contracts.code = ?", contractCode).
		First(&domesticServiceArea)
	if err != nil {
		return domesticServiceArea, fmt.Errorf("unable to find domestic service area for %s under contract code %s", zip3, contractCode)
	}

	return domesticServiceArea, nil
}

func LookupFSCWeightBasedDistanceMultiplier(appCtx appcontext.AppContext, primeEstimatedWeight unit.Pound) string {
	weight := primeEstimatedWeight.Int()
	const weightBasedDistanceMultiplierLevelOne = "0.000417"
	const weightBasedDistanceMultiplierLevelTwo = "0.0006255"
	const weightBasedDistanceMultiplierLevelThree = "0.000834"
	const weightBasedDistanceMultiplierLevelFour = "0.00139"

	if weight <= 5000 {
		return weightBasedDistanceMultiplierLevelOne
	} else if weight <= 10000 {
		return weightBasedDistanceMultiplierLevelTwo
	} else if weight <= 24000 {
		return weightBasedDistanceMultiplierLevelThree
		//nolint:revive
	} else {
		return weightBasedDistanceMultiplierLevelFour
	}
}

func LookupEIAFuelPrice(appCtx appcontext.AppContext, pickupDate time.Time) (unit.Millicents, error) {
	db := appCtx.DB()

	// Find the GHCDieselFuelPrice object with the closest prior PublicationDate to the ActualPickupDate of the MTOShipment in question
	var ghcDieselFuelPrice models.GHCDieselFuelPrice
	err := db.Where("publication_date <= ?", pickupDate).Order("publication_date DESC").Last(&ghcDieselFuelPrice)
	if err != nil {
		return 0, apperror.NewNotFoundError(uuid.Nil, "Unable to find GHCDieselFuelPrice")
	}
	return ghcDieselFuelPrice.FuelPriceInMillicents, nil
}
