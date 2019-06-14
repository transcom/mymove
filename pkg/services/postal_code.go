package services

// PostalCodeType is initialized as a string type
type PostalCodeType string

// PostalCodeTypes
const (
	Origin      PostalCodeType = "origin"
	Destination PostalCodeType = "destination"
)

// PostalCodeValidator is the service object interface for ValidatePostalCode
//go:generate mockery -name PostalCodeValidator
type PostalCodeValidator interface {
	ValidatePostalCode(postalCode string, postalCodeType PostalCodeType) (bool, error)
}
