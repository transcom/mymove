package services

import "github.com/transcom/mymove/pkg/appcontext"

// PostalCodeType is initialized as a string type
type PostalCodeType string

// PostalCodeTypes
const (
	Origin      PostalCodeType = "origin"
	Destination PostalCodeType = "destination"
)

// PostalCodeValidator is the service object interface for ValidatePostalCode
//go:generate mockery --name PostalCodeValidator --disable-version-string
type PostalCodeValidator interface {
	ValidatePostalCode(appCtx appcontext.AppContext, postalCode string, postalCodeType PostalCodeType) (bool, error)
}
