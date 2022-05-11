package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ServiceItemParamName is the name of service item parameter
type ServiceItemParamName string

func (s ServiceItemParamName) String() string {
	return string(s)
}

const (
	// ServiceItemParamNameActualPickupDate is the param key name ActualPickupDate
	ServiceItemParamNameActualPickupDate ServiceItemParamName = "ActualPickupDate"
	// ServiceItemParamNameContractCode is the param key name ContractCode
	ServiceItemParamNameContractCode ServiceItemParamName = "ContractCode"
	// ServiceItemParamNameContractYearName is the param key name ContractYearName
	ServiceItemParamNameContractYearName ServiceItemParamName = "ContractYearName"
	// ServiceItemParamNameCubicFeetBilled is the param key name CubicFeetBilled
	ServiceItemParamNameCubicFeetBilled ServiceItemParamName = "CubicFeetBilled"
	// ServiceItemParamNameCubicFeetCrating is the param key name CubicFeetCrating
	ServiceItemParamNameCubicFeetCrating ServiceItemParamName = "CubicFeetCrating"
	// ServiceItemParamNameDimensionHeight is the param key name DimensionHeight
	ServiceItemParamNameDimensionHeight ServiceItemParamName = "DimensionHeight"
	// ServiceItemParamNameDimensionLength is the param key name DimensionLength
	ServiceItemParamNameDimensionLength ServiceItemParamName = "DimensionLength"
	// ServiceItemParamNameDimensionWidth is the param key name DimensionWidth
	ServiceItemParamNameDimensionWidth ServiceItemParamName = "DimensionWidth"
	// ServiceItemParamNameDistanceZip3 is the param key name DistanceZip3
	ServiceItemParamNameDistanceZip3 ServiceItemParamName = "DistanceZip3"
	// ServiceItemParamNameDistanceZip5 is the param key name DistanceZip5
	ServiceItemParamNameDistanceZip5 ServiceItemParamName = "DistanceZip5"
	// ServiceItemParamNameDistanceZipSITDest is the param key name DistanceZipSITDest
	ServiceItemParamNameDistanceZipSITDest ServiceItemParamName = "DistanceZipSITDest"
	// ServiceItemParamNameDistanceZipSITOrigin is the param key name DistanceZipSITOrigin
	ServiceItemParamNameDistanceZipSITOrigin ServiceItemParamName = "DistanceZipSITOrigin"
	// ServiceItemParamNameEIAFuelPrice is the param key name EIAFuelPrice
	ServiceItemParamNameEIAFuelPrice ServiceItemParamName = "EIAFuelPrice"
	// ServiceItemParamNameEscalationCompounded is the param key name EscalationCompounded
	ServiceItemParamNameEscalationCompounded ServiceItemParamName = "EscalationCompounded"
	// ServiceItemParamNameFSCMultiplier is the param key name FSCMultiplier
	ServiceItemParamNameFSCMultiplier ServiceItemParamName = "FSCMultiplier"
	// ServiceItemParamNameFSCPriceDifferenceInCents is the param key name FSCPriceDifferenceInCents
	ServiceItemParamNameFSCPriceDifferenceInCents ServiceItemParamName = "FSCPriceDifferenceInCents"
	// ServiceItemParamNameFSCWeightBasedDistanceMultiplier is the param key name FSCWeightBasedDistanceMultiplier
	ServiceItemParamNameFSCWeightBasedDistanceMultiplier ServiceItemParamName = "FSCWeightBasedDistanceMultiplier"
	// ServiceItemParamNameIsPeak is the param key name IsPeak
	ServiceItemParamNameIsPeak ServiceItemParamName = "IsPeak"
	// ServiceItemParamNameMarketDest is the param key name MarketDest
	ServiceItemParamNameMarketDest ServiceItemParamName = "MarketDest"
	// ServiceItemParamNameMarketOrigin is the param key name MarketOrigin
	ServiceItemParamNameMarketOrigin ServiceItemParamName = "MarketOrigin"
	// ServiceItemParamNameMTOAvailableToPrimeAt is the param key name MTOAvailableToPrimeAt
	ServiceItemParamNameMTOAvailableToPrimeAt ServiceItemParamName = "MTOAvailableToPrimeAt"
	// ServiceItemParamNameNTSPackingFactor is the param key name NTSPackingFactor
	ServiceItemParamNameNTSPackingFactor ServiceItemParamName = "NTSPackingFactor"
	// ServiceItemParamNameNumberDaysSIT is the param key name NumberDaysSIT
	ServiceItemParamNameNumberDaysSIT ServiceItemParamName = "NumberDaysSIT"
	// ServiceItemParamNamePriceAreaDest is the param key name PriceAreaDest
	ServiceItemParamNamePriceAreaDest ServiceItemParamName = "PriceAreaDest"
	// ServiceItemParamNamePriceAreaIntlDest is the param key name PriceAreaIntlDest
	ServiceItemParamNamePriceAreaIntlDest ServiceItemParamName = "PriceAreaIntlDest"
	// ServiceItemParamNamePriceAreaIntlOrigin is the param key name PriceAreaIntlOrigin
	ServiceItemParamNamePriceAreaIntlOrigin ServiceItemParamName = "PriceAreaIntlOrigin"
	// ServiceItemParamNamePriceAreaOrigin is the param key name PriceAreaOrigin
	ServiceItemParamNamePriceAreaOrigin ServiceItemParamName = "PriceAreaOrigin"
	// ServiceItemParamNamePriceRateOrFactor is the param key name PriceRateOrFactor
	ServiceItemParamNamePriceRateOrFactor ServiceItemParamName = "PriceRateOrFactor"
	// ServiceItemParamNamePSILinehaulDom is the param key name PSI_LinehaulDom
	ServiceItemParamNamePSILinehaulDom ServiceItemParamName = "PSI_LinehaulDom"
	// ServiceItemParamNamePSILinehaulDomPrice is the param key name PSI_LinehaulDomPrice
	ServiceItemParamNamePSILinehaulDomPrice ServiceItemParamName = "PSI_LinehaulDomPrice"
	// ServiceItemParamNamePSILinehaulShort is the param key name PSI_LinehaulShort
	ServiceItemParamNamePSILinehaulShort ServiceItemParamName = "PSI_LinehaulShort"
	// ServiceItemParamNamePSILinehaulShortPrice is the param key name PSI_LinehaulShortPrice
	ServiceItemParamNamePSILinehaulShortPrice ServiceItemParamName = "PSI_LinehaulShortPrice"
	// ServiceItemParamNamePSIPriceDomDest is the param key name PSI_PriceDomDest
	ServiceItemParamNamePSIPriceDomDest ServiceItemParamName = "PSI_PriceDomDest"
	// ServiceItemParamNamePSIPriceDomDestPrice is the param key name PSI_PriceDomDestPrice
	ServiceItemParamNamePSIPriceDomDestPrice ServiceItemParamName = "PSI_PriceDomDestPrice"
	// ServiceItemParamNamePSIPriceDomOrigin is the param key name PSI_PriceDomOrigin
	ServiceItemParamNamePSIPriceDomOrigin ServiceItemParamName = "PSI_PriceDomOrigin"
	// ServiceItemParamNamePSIPriceDomOriginPrice is the param key name PSI_PriceDomOriginPrice
	ServiceItemParamNamePSIPriceDomOriginPrice ServiceItemParamName = "PSI_PriceDomOriginPrice"
	// ServiceItemParamNamePSIShippingLinehaulIntlCO is the param key name PSI_ShippingLinehaulIntlCO
	ServiceItemParamNamePSIShippingLinehaulIntlCO ServiceItemParamName = "PSI_ShippingLinehaulIntlCO"
	// ServiceItemParamNamePSIShippingLinehaulIntlCOPrice is the param key name PSI_ShippingLinehaulIntlCOPrice
	ServiceItemParamNamePSIShippingLinehaulIntlCOPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlCOPrice"
	// ServiceItemParamNamePSIShippingLinehaulIntlOC is the param key name PSI_ShippingLinehaulIntlOC
	ServiceItemParamNamePSIShippingLinehaulIntlOC ServiceItemParamName = "PSI_ShippingLinehaulIntlOC"
	// ServiceItemParamNamePSIShippingLinehaulIntlOCPrice is the param key name PSI_ShippingLinehaulIntlOCPrice
	ServiceItemParamNamePSIShippingLinehaulIntlOCPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlOCPrice"
	// ServiceItemParamNamePSIShippingLinehaulIntlOO is the param key name PSI_ShippingLinehaulIntlOO
	ServiceItemParamNamePSIShippingLinehaulIntlOO ServiceItemParamName = "PSI_ShippingLinehaulIntlOO"
	// ServiceItemParamNamePSIShippingLinehaulIntlOOPrice is the param key name PSI_ShippingLinehaulIntlOOPrice
	ServiceItemParamNamePSIShippingLinehaulIntlOOPrice ServiceItemParamName = "PSI_ShippingLinehaulIntlOOPrice"
	// ServiceItemParamNameRateAreaNonStdDest is the param key name RateAreaNonStdDest
	ServiceItemParamNameRateAreaNonStdDest ServiceItemParamName = "RateAreaNonStdDest"
	// ServiceItemParamNameRateAreaNonStdOrigin is the param key name RateAreaNonStdOrigin
	ServiceItemParamNameRateAreaNonStdOrigin ServiceItemParamName = "RateAreaNonStdOrigin"
	// ServiceItemParamNameReferenceDate is the param key name ReferenceDate
	ServiceItemParamNameReferenceDate ServiceItemParamName = "ReferenceDate"
	// ServiceItemParamNameRequestedPickupDate is the param key name RequestedPickupDate
	ServiceItemParamNameRequestedPickupDate ServiceItemParamName = "RequestedPickupDate"
	// ServiceItemParamNameServiceAreaDest is the param key name ServiceAreaDest
	ServiceItemParamNameServiceAreaDest ServiceItemParamName = "ServiceAreaDest"
	// ServiceItemParamNameServiceAreaOrigin is the param key name ServiceAreaOrigin
	ServiceItemParamNameServiceAreaOrigin ServiceItemParamName = "ServiceAreaOrigin"
	// ServiceItemParamNameServicesScheduleDest is the param key name ServicesScheduleDest
	ServiceItemParamNameServicesScheduleDest ServiceItemParamName = "ServicesScheduleDest"
	// ServiceItemParamNameServicesScheduleOrigin is the param key name ServicesScheduleOrigin
	ServiceItemParamNameServicesScheduleOrigin ServiceItemParamName = "ServicesScheduleOrigin"
	// ServiceItemParamNameSITPaymentRequestEnd is the param key name SITPaymentRequestEnd
	ServiceItemParamNameSITPaymentRequestEnd ServiceItemParamName = "SITPaymentRequestEnd"
	// ServiceItemParamNameSITPaymentRequestStart is the param key name SITPaymentRequestStart
	ServiceItemParamNameSITPaymentRequestStart ServiceItemParamName = "SITPaymentRequestStart"
	// ServiceItemParamNameSITScheduleDest is the param key name SITScheduleDest
	ServiceItemParamNameSITScheduleDest ServiceItemParamName = "SITScheduleDest"
	// ServiceItemParamNameSITScheduleOrigin is the param key name SITScheduleOrigin
	ServiceItemParamNameSITScheduleOrigin ServiceItemParamName = "SITScheduleOrigin"
	// ServiceItemParamNameWeightAdjusted is the param key name WeightAdjusted
	ServiceItemParamNameWeightAdjusted ServiceItemParamName = "WeightAdjusted"
	// ServiceItemParamNameWeightBilled is the param key name WeightBilled
	ServiceItemParamNameWeightBilled ServiceItemParamName = "WeightBilled"
	// ServiceItemParamNameWeightEstimated is the param key name WeightEstimated
	ServiceItemParamNameWeightEstimated ServiceItemParamName = "WeightEstimated"
	// ServiceItemParamNameWeightOriginal is the param key name WeightOriginal
	ServiceItemParamNameWeightOriginal ServiceItemParamName = "WeightOriginal"
	// ServiceItemParamNameWeightReweigh is the param key name WeightReweigh
	ServiceItemParamNameWeightReweigh ServiceItemParamName = "WeightReweigh"
	// ServiceItemParamNameZipDestAddress is the param key name ZipDestAddress
	ServiceItemParamNameZipDestAddress ServiceItemParamName = "ZipDestAddress"
	// ServiceItemParamNameZipPickupAddress is the param key name ZipPickupAddress
	ServiceItemParamNameZipPickupAddress ServiceItemParamName = "ZipPickupAddress"
	// ServiceItemParamNameZipSITDestHHGFinalAddress is the param key name ZipSITDestHHGFinalAddress
	ServiceItemParamNameZipSITDestHHGFinalAddress ServiceItemParamName = "ZipSITDestHHGFinalAddress"
	// ServiceItemParamNameZipSITOriginHHGActualAddress is the param key name ZipSITOriginHHGActualAddress
	ServiceItemParamNameZipSITOriginHHGActualAddress ServiceItemParamName = "ZipSITOriginHHGActualAddress"
	// ServiceItemParamNameZipSITOriginHHGOriginalAddress is the param key name ZipSITOriginHHGOriginalAddress
	ServiceItemParamNameZipSITOriginHHGOriginalAddress ServiceItemParamName = "ZipSITOriginHHGOriginalAddress"
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
	// ServiceItemParamTypeTimestamp is a timestamp
	ServiceItemParamTypeTimestamp ServiceItemParamType = "TIMESTAMP"
	// ServiceItemParamTypePaymentServiceItemUUID is a UUID
	ServiceItemParamTypePaymentServiceItemUUID ServiceItemParamType = "PaymentServiceItemUUID"
	// ServiceItemParamTypeBoolean is a boolean
	ServiceItemParamTypeBoolean ServiceItemParamType = "BOOLEAN"
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
	// ServiceItemParamOriginPricer is the Pricer origin
	ServiceItemParamOriginPricer ServiceItemParamOrigin = "PRICER"
	// ServiceItemParamOriginPaymentRequest is the PaymentRequest origin
	ServiceItemParamOriginPaymentRequest ServiceItemParamOrigin = "PAYMENT_REQUEST"
)

// ValidServiceItemParamNames lists all valid service item param key names
var ValidServiceItemParamNames = []ServiceItemParamName{
	ServiceItemParamNameActualPickupDate,
	ServiceItemParamNameContractCode,
	ServiceItemParamNameContractYearName,
	ServiceItemParamNameCubicFeetBilled,
	ServiceItemParamNameCubicFeetCrating,
	ServiceItemParamNameDimensionHeight,
	ServiceItemParamNameDimensionLength,
	ServiceItemParamNameDimensionWidth,
	ServiceItemParamNameDistanceZip3,
	ServiceItemParamNameDistanceZip5,
	ServiceItemParamNameDistanceZipSITDest,
	ServiceItemParamNameDistanceZipSITOrigin,
	ServiceItemParamNameEIAFuelPrice,
	ServiceItemParamNameEscalationCompounded,
	ServiceItemParamNameFSCMultiplier,
	ServiceItemParamNameFSCPriceDifferenceInCents,
	ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
	ServiceItemParamNameIsPeak,
	ServiceItemParamNameMarketDest,
	ServiceItemParamNameMarketOrigin,
	ServiceItemParamNameMTOAvailableToPrimeAt,
	ServiceItemParamNameNTSPackingFactor,
	ServiceItemParamNameNumberDaysSIT,
	ServiceItemParamNamePriceAreaDest,
	ServiceItemParamNamePriceAreaIntlDest,
	ServiceItemParamNamePriceAreaIntlOrigin,
	ServiceItemParamNamePriceAreaOrigin,
	ServiceItemParamNamePriceRateOrFactor,
	ServiceItemParamNamePSILinehaulDom,
	ServiceItemParamNamePSILinehaulDomPrice,
	ServiceItemParamNamePSILinehaulShort,
	ServiceItemParamNamePSILinehaulShortPrice,
	ServiceItemParamNamePSIPriceDomDest,
	ServiceItemParamNamePSIPriceDomDestPrice,
	ServiceItemParamNamePSIPriceDomOrigin,
	ServiceItemParamNamePSIPriceDomOriginPrice,
	ServiceItemParamNamePSIShippingLinehaulIntlCO,
	ServiceItemParamNamePSIShippingLinehaulIntlCOPrice,
	ServiceItemParamNamePSIShippingLinehaulIntlOC,
	ServiceItemParamNamePSIShippingLinehaulIntlOCPrice,
	ServiceItemParamNamePSIShippingLinehaulIntlOO,
	ServiceItemParamNamePSIShippingLinehaulIntlOOPrice,
	ServiceItemParamNameRateAreaNonStdDest,
	ServiceItemParamNameRateAreaNonStdOrigin,
	ServiceItemParamNameReferenceDate,
	ServiceItemParamNameRequestedPickupDate,
	ServiceItemParamNameServiceAreaDest,
	ServiceItemParamNameServiceAreaOrigin,
	ServiceItemParamNameServicesScheduleDest,
	ServiceItemParamNameServicesScheduleOrigin,
	ServiceItemParamNameSITPaymentRequestEnd,
	ServiceItemParamNameSITPaymentRequestStart,
	ServiceItemParamNameSITScheduleDest,
	ServiceItemParamNameSITScheduleOrigin,
	ServiceItemParamNameWeightAdjusted,
	ServiceItemParamNameWeightBilled,
	ServiceItemParamNameWeightEstimated,
	ServiceItemParamNameWeightOriginal,
	ServiceItemParamNameWeightReweigh,
	ServiceItemParamNameZipDestAddress,
	ServiceItemParamNameZipPickupAddress,
	ServiceItemParamNameZipSITDestHHGFinalAddress,
	ServiceItemParamNameZipSITOriginHHGActualAddress,
	ServiceItemParamNameZipSITOriginHHGOriginalAddress,
}

// ValidServiceItemParamNameStrings lists all valid service item param key names
var ValidServiceItemParamNameStrings = []string{
	string(ServiceItemParamNameActualPickupDate),
	string(ServiceItemParamNameContractCode),
	string(ServiceItemParamNameContractYearName),
	string(ServiceItemParamNameCubicFeetBilled),
	string(ServiceItemParamNameCubicFeetCrating),
	string(ServiceItemParamNameDimensionHeight),
	string(ServiceItemParamNameDimensionLength),
	string(ServiceItemParamNameDimensionWidth),
	string(ServiceItemParamNameDistanceZip3),
	string(ServiceItemParamNameDistanceZip5),
	string(ServiceItemParamNameDistanceZipSITDest),
	string(ServiceItemParamNameDistanceZipSITOrigin),
	string(ServiceItemParamNameEIAFuelPrice),
	string(ServiceItemParamNameEscalationCompounded),
	string(ServiceItemParamNameFSCMultiplier),
	string(ServiceItemParamNameFSCPriceDifferenceInCents),
	string(ServiceItemParamNameFSCWeightBasedDistanceMultiplier),
	string(ServiceItemParamNameIsPeak),
	string(ServiceItemParamNameMarketDest),
	string(ServiceItemParamNameMarketOrigin),
	string(ServiceItemParamNameMTOAvailableToPrimeAt),
	string(ServiceItemParamNameNTSPackingFactor),
	string(ServiceItemParamNameNumberDaysSIT),
	string(ServiceItemParamNamePriceAreaDest),
	string(ServiceItemParamNamePriceAreaIntlDest),
	string(ServiceItemParamNamePriceAreaIntlOrigin),
	string(ServiceItemParamNamePriceAreaOrigin),
	string(ServiceItemParamNamePriceRateOrFactor),
	string(ServiceItemParamNamePSILinehaulDom),
	string(ServiceItemParamNamePSILinehaulDomPrice),
	string(ServiceItemParamNamePSILinehaulShort),
	string(ServiceItemParamNamePSILinehaulShortPrice),
	string(ServiceItemParamNamePSIPriceDomDest),
	string(ServiceItemParamNamePSIPriceDomDestPrice),
	string(ServiceItemParamNamePSIPriceDomOrigin),
	string(ServiceItemParamNamePSIPriceDomOriginPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlCO),
	string(ServiceItemParamNamePSIShippingLinehaulIntlCOPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOC),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOCPrice),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOO),
	string(ServiceItemParamNamePSIShippingLinehaulIntlOOPrice),
	string(ServiceItemParamNameRateAreaNonStdDest),
	string(ServiceItemParamNameRateAreaNonStdOrigin),
	string(ServiceItemParamNameReferenceDate),
	string(ServiceItemParamNameRequestedPickupDate),
	string(ServiceItemParamNameServiceAreaDest),
	string(ServiceItemParamNameServiceAreaOrigin),
	string(ServiceItemParamNameServicesScheduleDest),
	string(ServiceItemParamNameServicesScheduleOrigin),
	string(ServiceItemParamNameSITPaymentRequestEnd),
	string(ServiceItemParamNameSITPaymentRequestStart),
	string(ServiceItemParamNameSITScheduleDest),
	string(ServiceItemParamNameSITScheduleOrigin),
	string(ServiceItemParamNameWeightAdjusted),
	string(ServiceItemParamNameWeightBilled),
	string(ServiceItemParamNameWeightEstimated),
	string(ServiceItemParamNameWeightOriginal),
	string(ServiceItemParamNameWeightReweigh),
	string(ServiceItemParamNameZipDestAddress),
	string(ServiceItemParamNameZipPickupAddress),
	string(ServiceItemParamNameZipSITDestHHGFinalAddress),
	string(ServiceItemParamNameZipSITOriginHHGActualAddress),
	string(ServiceItemParamNameZipSITOriginHHGOriginalAddress),
}

// ValidServiceItemParamTypes lists all valid service item param types
var ValidServiceItemParamTypes = []string{
	string(ServiceItemParamTypeString),
	string(ServiceItemParamTypeDate),
	string(ServiceItemParamTypeInteger),
	string(ServiceItemParamTypeDecimal),
	string(ServiceItemParamTypeTimestamp),
	string(ServiceItemParamTypePaymentServiceItemUUID),
	string(ServiceItemParamTypeBoolean),
}

// ValidServiceItemParamOrigins lists all valid service item param origins
var ValidServiceItemParamOrigins = []string{
	string(ServiceItemParamOriginPrime),
	string(ServiceItemParamOriginSystem),
	string(ServiceItemParamOriginPricer),
	string(ServiceItemParamOriginPaymentRequest),
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
		&validators.StringInclusion{Field: s.Key.String(), Name: "Key", List: ValidServiceItemParamNameStrings},
		&validators.StringIsPresent{Field: s.Description, Name: "Description"},
		&validators.StringIsPresent{Field: string(s.Type), Name: "Type"},
		&validators.StringInclusion{Field: s.Type.String(), Name: "Type", List: ValidServiceItemParamTypes},
		&validators.StringIsPresent{Field: string(s.Origin), Name: "Origin"},
		&validators.StringInclusion{Field: s.Origin.String(), Name: "Origin", List: ValidServiceItemParamOrigins},
	), nil
}
