package services

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// ServiceItemPricer prices a generic payment service item for a GHC move
//go:generate mockery -name ServiceItemPricer
type ServiceItemPricer interface {
	PriceServiceItem(item models.PaymentServiceItem) (unit.Cents, error)
	UsingConnection(db *pop.Connection) ServiceItemPricer
}

// PricingParam represents a parameter (key/value pair) returned from a pricer
// nolint:structcheck // TODO: Remove once we start using these fields in pricers
type PricingParam struct {
	key   models.ServiceItemParamName
	value interface{}
}

// ParamsPricer is an interface that all param-aware pricers implement
type ParamsPricer interface {
	PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, []PricingParam, error)
}

// ManagementServicesPricer prices management services for a GHC move
//go:generate mockery -name ManagementServicesPricer
type ManagementServicesPricer interface {
	Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// CounselingServicesPricer prices counseling services for a GHC move
//go:generate mockery -name CounselingServicesPricer
type CounselingServicesPricer interface {
	Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery -name DomesticLinehaulPricer
type DomesticLinehaulPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticShorthaulPricer prices the domestic shorthaul for a GHC Move
//go:generate mockery -name DomesticShorthaulPricer
type DomesticShorthaulPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticOriginPricer prices the domestic origin for a GHC Move
//go:generate mockery -name DomesticOriginPricer
type DomesticOriginPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticDestinationPricer prices the domestic destination price for a GHC Move
//go:generate mockery -name DomesticDestinationPricer
type DomesticDestinationPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticPackPricer prices the domestic packing and unpacking for a GHC Move
//go:generate mockery -name DomesticPackPricer
type DomesticPackPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticUnpackPricer prices the domestic unpacking for a GHC Move
//go:generate mockery -name DomesticUnpackPricer
type DomesticUnpackPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// FuelSurchargePricer prices the fuel surcharge price for a GHC Move
//go:generate mockery -name FuelSurchargePricer
type FuelSurchargePricer interface {
	Price(contractCode string, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, weightBasedDistanceMultiplier float64, fuelPrice unit.Millicents) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticOriginFirstDaySITPricer prices domestic origin first day SIT for a GHC move
//go:generate mockery -name DomesticOriginFirstDaySITPricer
type DomesticOriginFirstDaySITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticDestinationFirstDaySITPricer prices domestic destination first day SIT for a GHC move
//go:generate mockery -name DomesticDestinationFirstDaySITPricer
type DomesticDestinationFirstDaySITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticOriginAdditionalDaysSITPricer prices domestic origin additional days SIT for a GHC move
//go:generate mockery -name DomesticOriginAdditionalDaysSITPricer
type DomesticOriginAdditionalDaysSITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticDestinationAdditionalDaysSITPricer prices domestic destination additional days SIT for a GHC move
//go:generate mockery -name DomesticDestinationAdditionalDaysSITPricer
type DomesticDestinationAdditionalDaysSITPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticOriginSITPickupPricer prices domestic origin SIT pickup for a GHC move
//go:generate mockery -name DomesticOriginSITPickupPricer
type DomesticOriginSITPickupPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipSITOriginOriginal string, zipSITOriginActual string, distance unit.Miles) (unit.Cents, []PricingParam, error)
	ParamsPricer
}

// DomesticDestinationSITDeliveryPricer prices domestic destination SIT delivery for a GHC move
//go:generate mockery -name DomesticDestinationSITDeliveryPricer
type DomesticDestinationSITDeliveryPricer interface {
	Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipDest string, zipSITDest string, distance unit.Miles) (unit.Cents, []PricingParam, error)
	ParamsPricer
}
