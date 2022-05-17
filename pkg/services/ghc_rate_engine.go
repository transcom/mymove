package services

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// ServiceItemPricer prices a generic payment service item for a GHC move
//go:generate mockery --name ServiceItemPricer --disable-version-string
type ServiceItemPricer interface {
	PriceServiceItem(appCtx appcontext.AppContext, item models.PaymentServiceItem) (unit.Cents, models.PaymentServiceItemParams, error)
}

// PricingDisplayParam represents a parameter (key/value pair) returned from a pricer
type PricingDisplayParam struct {
	Key   models.ServiceItemParamName
	Value string
}

// PricingDisplayParams represents a slice of pricing parameters
type PricingDisplayParams []PricingDisplayParam

// ParamsPricer is an interface that all param-aware pricers implement
type ParamsPricer interface {
	PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, PricingDisplayParams, error)
}

// ManagementServicesPricer prices management services for a GHC move
//go:generate mockery --name ManagementServicesPricer --disable-version-string
type ManagementServicesPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// CounselingServicesPricer prices counseling services for a GHC move
//go:generate mockery --name CounselingServicesPricer --disable-version-string
type CounselingServicesPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticLinehaulPricer prices domestic linehaul for a GHC move
//go:generate mockery --name DomesticLinehaulPricer --disable-version-string
type DomesticLinehaulPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string, isPPM bool) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticShorthaulPricer prices the domestic shorthaul for a GHC Move
//go:generate mockery --name DomesticShorthaulPricer --disable-version-string
type DomesticShorthaulPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticOriginPricer prices the domestic origin for a GHC Move
//go:generate mockery --name DomesticOriginPricer --disable-version-string
type DomesticOriginPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticDestinationPricer prices the domestic destination price for a GHC Move
//go:generate mockery --name DomesticDestinationPricer --disable-version-string
type DomesticDestinationPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, isPPM bool) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticOriginShuttlingPricer prices the domestic origin shuttling service for a GHC Move
//go:generate mockery --name DomesticOriginShuttlingPricer --disable-version-string
type DomesticOriginShuttlingPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticDestinationShuttlingPricer prices the domestic origin shuttling service for a GHC Move
//go:generate mockery --name DomesticDestinationShuttlingPricer --disable-version-string
type DomesticDestinationShuttlingPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticCratingPricer prices the domestic crating service for a GHC Move
//go:generate mockery --name DomesticCratingPricer --disable-version-string
type DomesticCratingPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, servicesScheduleOrigin int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticUncratingPricer prices the domestic uncrating service for a GHC Move
//go:generate mockery --name DomesticUncratingPricer --disable-version-string
type DomesticUncratingPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, servicesScheduleDest int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticPackPricer prices the domestic packing for a GHC Move
//go:generate mockery --name DomesticPackPricer --disable-version-string
type DomesticPackPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticNTSPackPricer prices the domestic packing for an NTS shipment of a GHC Move
//go:generate mockery --name DomesticNTSPackPricer --disable-version-string
type DomesticNTSPackPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticUnpackPricer prices the domestic unpacking for a GHC Move
//go:generate mockery --name DomesticUnpackPricer --disable-version-string
type DomesticUnpackPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// FuelSurchargePricer prices the fuel surcharge price for a GHC Move
//go:generate mockery --name FuelSurchargePricer --disable-version-string
type FuelSurchargePricer interface {
	Price(appCtx appcontext.AppContext, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, fscWeightBasedDistanceMultiplier float64, eiaFuelPrice unit.Millicents) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticOriginFirstDaySITPricer prices domestic origin first day SIT for a GHC move
//go:generate mockery --name DomesticOriginFirstDaySITPricer --disable-version-string
type DomesticOriginFirstDaySITPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticDestinationFirstDaySITPricer prices domestic destination first day SIT for a GHC move
//go:generate mockery --name DomesticDestinationFirstDaySITPricer --disable-version-string
type DomesticDestinationFirstDaySITPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticOriginAdditionalDaysSITPricer prices domestic origin additional days SIT for a GHC move
//go:generate mockery --name DomesticOriginAdditionalDaysSITPricer --disable-version-string
type DomesticOriginAdditionalDaysSITPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticDestinationAdditionalDaysSITPricer prices domestic destination additional days SIT for a GHC move
//go:generate mockery --name DomesticDestinationAdditionalDaysSITPricer --disable-version-string
type DomesticDestinationAdditionalDaysSITPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticOriginSITPickupPricer prices domestic origin SIT pickup for a GHC move
//go:generate mockery --name DomesticOriginSITPickupPricer --disable-version-string
type DomesticOriginSITPickupPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipSITOriginOriginal string, zipSITOriginActual string, distance unit.Miles) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}

// DomesticDestinationSITDeliveryPricer prices domestic destination SIT delivery for a GHC move
//go:generate mockery --name DomesticDestinationSITDeliveryPricer --disable-version-string
type DomesticDestinationSITDeliveryPricer interface {
	Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipDest string, zipSITDest string, distance unit.Miles) (unit.Cents, PricingDisplayParams, error)
	ParamsPricer
}
