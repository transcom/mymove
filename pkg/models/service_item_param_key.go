package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

//ServiceItemParamName is the name of service item parameter
type ServiceItemParamName string

func (s ServiceItemParamName) String() string {
	return string(s)
}

const (
	ServiceItemParamNameRequestedPickupDate            ServiceItemParamName = "RequestedPickupDate"
	ServiceItemParamNameWeightBilledActual             ServiceItemParamName = "WeightBilledActual"
	ServiceItemParamNameWeightActual                   ServiceItemParamName = "WeightActual"
	ServiceItemParamNameWeightEstimated                ServiceItemParamName = "WeightEstimated"
	ServiceItemParamNameDistanceZip5                   ServiceItemParamName = "DistanceZip5"
	ServiceItemParamNameDistanceZip3                   ServiceItemParamName = "DistanceZip3"
	ServiceItemParamNameZipPickupAddress               ServiceItemParamName = "ZipPickupAddress"
	ServiceItemParamNameZipDestAddress                 ServiceItemParamName = "ZipDestAddress"
	ServiceItemParamNameServiceAreaOrigin              ServiceItemParamName = "ServiceAreaOrigin"
	ServiceItemParamNameNumberDaysSIT                  ServiceItemParamName = "NumberDaysSIT"
	ServiceItemParamNameMarketOrigin                   ServiceItemParamName = "MarketOrigin"
	ServiceItemParamNameMarketDest                     ServiceItemParamName = "MarketDest"
	ServiceItemParamNameCanStandAlone                  ServiceItemParamName = "CanStandAlone"
	ServiceItemParamNameCubicFeetBilled                ServiceItemParamName = "CubicFeetBilled"
	ServiceItemParamNameCubicFeetCrating               ServiceItemParamName = "CubicFeetCrating"
	ServiceItemParamNameDistanceZip5SITOrigin          ServiceItemParamName = "DistanceZip5SITOrigin"
	ServiceItemParamNameDistanceZip5SITDest            ServiceItemParamName = "DistanceZip5SITDest"
	ServiceItemParamNameEIAFuelPrice                   ServiceItemParamName = "EIAFuelPrice"
	ServiceItemParamNameServiceAreaDest                ServiceItemParamName = "ServiceAreaDest"
	ServiceItemParamNameSITScheduleOrigin              ServiceItemParamName = "SITScheduleOrigin"
	ServiceItemParamNameSITScheduleDest                ServiceItemParamName = "SITScheduleDest"
	ServiceItemParamNameServicesScheduleOrigin         ServiceItemParamName = "ServicesScheduleOrigin"
	ServiceItemParamNameServicesScheduleDest           ServiceItemParamName = "ServicesScheduleDest"
	ServiceItemParamNamePriceAreaOrigin                ServiceItemParamName = "PriceAreaOrigin"
	ServiceItemParamNamePriceAreaDest                  ServiceItemParamName = "PriceAreaDest"
	ServiceItemParamNamePriceAreaIntlOrigin            ServiceItemParamName = "PriceAreaIntlOrigin"
	ServiceItemParamNamePriceAreaIntlDest              ServiceItemParamName = "PriceAreaIntlDest"
	ServiceItemParamNameRateAreaNonStdOrigin           ServiceItemParamName = "RateAreaNonStdOrigin"
	ServiceItemParamNameRateAreaNonStdDest             ServiceItemParamName = "RateAreaNonStdDest"
	ServiceItemParamNamePSILinehaulDom                 ServiceItemParamName = "PSI_LinehaulDom"
	ServiceItemParamNamePSILinehaulDomPrice            ServiceItemParamName = "PSI_LinehaulDomPrice"
	ServiceItemParamNamePSILinehaulShort               ServiceItemParamName = "PSI_LinehaulShort"
	ServiceItemParamNamePSILinehaulShortPrice          ServiceItemParamName = "PSI_LinehaulShortPrice"
	ServiceItemParamNamePSIPriceDomOrigin              ServiceItemParamName = "PSI_PriceDomOrigin"
	ServiceItemParamNamePSIPriceDomOriginPrice         ServiceItemParamName = "PSI_PriceDomOriginPrice"
	ServiceItemParamNamePSIPriceDomDest                ServiceItemParamName = "PSI_PriceDomDest"
	ServiceItemParamNamePSIPriceDomDestPrice           ServiceItemParamName = "PSI_PriceDomDestPrice"
	ServiceItemParamNamePSIShippingLinehaulIntlCO      ServiceItemParamName = "PSI_ShippingLinehaulIntlCO"
	ServiceItemParamNamePSIShippingLinehaulIntlCOPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlCOPrice"
	ServiceItemParamNamePSIShippingLinehaulIntlOC      ServiceItemParamName = "PSI_ShippingLinehaulIntlOC"
	ServiceItemParamNamePSIShippingLinehaulIntlOCPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlOCPrice"
	ServiceItemParamNamePSIShippingLinehaulIntlOO      ServiceItemParamName = "PSI_ShippingLinehaulIntlOO"
	ServiceItemParamNamePSIShippingLinehaulIntlOOPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlOOPrice"
	ServiceItemParamNamePSIPackingDom                  ServiceItemParamName = "PSI_PackingDom"
	ServiceItemParamNamePSIPackingDomPrice             ServiceItemParamName = "PSI_PackingDomPrice"
	ServiceItemParamNamePSIPackingHHGIntl              ServiceItemParamName = "PSI_PackingHHGIntl"
	ServiceItemParamNamePSIPackingHHGIntlPrice         ServiceItemParamName = "PSI_PackingHHGIntlPrice"
)

// ServiceItemParamType is a type of service item parameter
type ServiceItemParamType string

// String is a string representation of a ServiceItemParamType
func (s ServiceItemParamType) String() string {
	return string(s)
}

const (
	// ServiceItemParamTypeString is a string
	ServiceItemParamTypeString ServiceItemParamType = "STRING"
	// ServiceItemParamTypeDate is a date
	ServiceItemParamTypeDate ServiceItemParamType = "DATE"
	// ServiceItemParamTypeInteger is an integer
	ServiceItemParamTypeInteger ServiceItemParamType = "INTEGER"
	// ServiceItemParamTypeDecimal is a decimal
	ServiceItemParamTypeDecimal ServiceItemParamType = "DECIMAL"
)

// ServiceItemParamOrigin is a type of service item parameter origin
type ServiceItemParamOrigin string

// String is a string representation of a ServiceItemParamOrigin
func (s ServiceItemParamOrigin) String() string {
	return string(s)
}

const (
	// ServiceItemParamOriginPrime is the Prime origin
	ServiceItemParamOriginPrime ServiceItemParamOrigin = "PRIME"
	// ServiceItemParamOriginSystem is the System origin
	ServiceItemParamOriginSystem ServiceItemParamOrigin = "SYSTEM"
)

// ValidServiceItemParamName lists all valid service item param key names
var ValidServiceItemParamName = []string{
	string(ServiceItemParamNameRequestedPickupDate),
	string(ServiceItemParamNameWeightBilledActual),
	string(ServiceItemParamNameWeightActual),
	string(ServiceItemParamNameWeightEstimated),
	string(ServiceItemParamNameDistanceZip5),
	string(ServiceItemParamNameDistanceZip3),
	string(ServiceItemParamNameZipPickupAddress),
	string(ServiceItemParamNameZipDestAddress),
	string(ServiceItemParamNameServiceAreaOrigin),
	string(ServiceItemParamNameNumberDaysSIT),
	string(ServiceItemParamNameMarketOrigin),
	string(ServiceItemParamNameMarketDest),
	string(ServiceItemParamNameCanStandAlone),
	string(ServiceItemParamNameCubicFeetBilled),
	string(ServiceItemParamNameCubicFeetCrating),
	string(ServiceItemParamNameDistanceZip5SITOrigin),
	string(ServiceItemParamNameDistanceZip5SITDest),
	string(ServiceItemParamNameEIAFuelPrice),
	string(ServiceItemParamNameServiceAreaDest),
	string(ServiceItemParamNameSITScheduleOrigin),
	string(ServiceItemParamNameSITScheduleDest),
	string(ServiceItemParamNameServicesScheduleOrigin),
	string(ServiceItemParamNameServicesScheduleDest),
	string(ServiceItemParamNamePriceAreaOrigin),
	string(ServiceItemParamNamePriceAreaDest),
	string(ServiceItemParamNamePriceAreaIntlOrigin),
	string(ServiceItemParamNamePriceAreaIntlDest),
	string(ServiceItemParamNameRateAreaNonStdOrigin),
	string(ServiceItemParamNameRateAreaNonStdDest),
	string(ServiceItemParamNamePSILinehaulDom),
	string(ServiceItemParamNamePSILinehaulDomPrice),
	string(ServiceItemParamNamePSILinehaulShort),
	string(ServiceItemParamNamePSILinehaulShortPrice),
	string(ServiceItemParamNamePSIPriceDomOrigin),
	string(ServiceItemParamNamePSIPriceDomOriginPrice),
	string(ServiceItemParamNamePSIPriceDomDest),
	string(ServiceItemParamNamePSIPriceDomDestPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlCO),
	string(ServiceItemParamNamePSIShippingLinehaulIntlCOPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOC),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOCPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOO),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOOPrice),
	string(ServiceItemParamNamePSIPackingDom),
	string(ServiceItemParamNamePSIPackingDomPrice),
	string(ServiceItemParamNamePSIPackingHHGIntl),
	string(ServiceItemParamNamePSIPackingHHGIntlPrice),
}

var validServiceItemParamType = []string{
	string(ServiceItemParamTypeString),
	string(ServiceItemParamTypeDate),
	string(ServiceItemParamTypeInteger),
	string(ServiceItemParamTypeDecimal),
}

var validServiceItemParamOrigin = []string{
	string(ServiceItemParamOriginPrime),
	string(ServiceItemParamOriginSystem),
}

// ServiceItemParamKey is a key for a Service Item Param
type ServiceItemParamKey struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Key         ServiceItemParamName   `json:"key" db:"key"`
	Description string                 `json:"description" db:"description"`
	Type        ServiceItemParamType   `json:"type" db:"type"`
	Origin      ServiceItemParamOrigin `json:"origin" db:"origin"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}

// ServiceItemParamKeys is not required by pop and may be deleted
type ServiceItemParamKeys []ServiceItemParamKey

// Validate validates a ServiceItemParamKey
func (s *ServiceItemParamKey) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Key.String(), Name: "Key"},
		&validators.StringIsPresent{Field: s.Description, Name: "Description"},
		&validators.StringIsPresent{Field: string(s.Type), Name: "Type"},
		&validators.StringIsPresent{Field: string(s.Origin), Name: "Origin"},
		&validators.StringInclusion{Field: s.Type.String(), Name: "Type", List: validServiceItemParamType},
		&validators.StringInclusion{Field: s.Origin.String(), Name: "Origin", List: validServiceItemParamOrigin},
	), nil
}
