package services

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"

	"github.com/transcom/mymove/pkg/unit"
)

// ServiceItemPricer prices a generic payment service item for a GHC move
//go:generate mockery -name ServiceItemPricer
type ServiceItemPricer interface {
	PriceServiceItem(item models.PaymentServiceItem) (unit.Cents, error)
	UsingConnection(db *pop.Connection) ServiceItemPricer
}

// ParamsPricer is an interface that all param-aware pricers implement
type ParamsPricer interface {
	PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error)
}

// ManagementServicesPricer prices management services for a GHC move
//go:generate mockery -name ManagementServicesPricer
type ManagementServicesPricer interface {
	Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, error)
	ParamsPricer
}

// CounselingServicesPricer prices counseling services for a GHC move
//go:generate mockery -name CounselingServicesPricer
type CounselingServicesPricer interface {
	Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, error)
	ParamsPricer
}

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery -name DomesticLinehaulPricer
type DomesticLinehaulPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, distance int, weightBilledActual int, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticShorthaulPricer prices the domestic shorthaul for a GHC Move
//go:generate mockery -name DomesticShorthaulPricer
type DomesticShorthaulPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticOriginPricer prices the domestic origin for a GHC Move
//go:generate mockery -name DomesticOriginPricer
type DomesticOriginPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticDestinationPricer prices the domestic destination price for a GHC Move
//go:generate mockery -name DomesticDestinationPricer
type DomesticDestinationPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticPackPricer prices the domestic packing and unpacking for a GHC Move
//go:generate mockery -name DomesticPackPricer
type DomesticPackPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, error)
	ParamsPricer
}

// DomesticUnpackPricer prices the domestic unpacking for a GHC Move
//go:generate mockery -name DomesticUnpackPricer
type DomesticUnpackPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (unit.Cents, error)
	ParamsPricer
}

// FuelSurchargePricer prices the fuel surcharge price for a GHC Move
//go:generate mockery -name FuelSurchargePricer
type FuelSurchargePricer interface {
	Price(contractCode string, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, weightBasedDistanceMultiplier float64, fuelPrice unit.Millicents) (unit.Cents, error)
	ParamsPricer
}

// DomesticOriginFirstDaySITPricer prices domestic origin first day SIT for a GHC move
//go:generate mockery -name DomesticOriginFirstDaySITPricer
type DomesticOriginFirstDaySITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticDestinationFirstDaySITPricer prices domestic destination first day SIT for a GHC move
//go:generate mockery -name DomesticDestinationFirstDaySITPricer
type DomesticDestinationFirstDaySITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string) (unit.Cents, error)
	ParamsPricer
}

// DomesticDestinationAdditionalDaysSITPricer prices domestic destination additional days SIT for a GHC move
//go:generate mockery -name DomesticDestinationAdditionalDaysSITPricer
type DomesticDestinationAdditionalDaysSITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, error)
	ParamsPricer
}

// Older pricers below (pre-dates payment requests)

// DomesticServiceAreaPricer domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
//go:generate mockery -name DomesticServiceAreaPricer
type DomesticServiceAreaPricer interface {
	PriceDomesticServiceArea(moveDate time.Time, weight unit.Pound, serviceArea string, servicesCode string) (unit.Cents, error)
}

//DomesticFuelSurchargePricer prices fuel surcharge for domestic GHC moves
//go:generate mockery -name DomesticFuelSurchargePricer
type DomesticFuelSurchargePricer interface {
	PriceDomesticFuelSurcharge(moveDate time.Time, planner route.Planner, weight unit.Pound, source string, destination string) (unit.Cents, error)
}
