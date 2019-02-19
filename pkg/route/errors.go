package route

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// ErrorCode contains error codes for the route package
type ErrorCode string

const (
	// UnsupportedPostalCode happens when we can't map a ZIP5 to a set of Lat/Long
	UnsupportedPostalCode = "UNSUPPORTED_POSTAL_CODE"
	// UnroutableRoute happens when a valid route can't be calculated between two locations
	UnroutableRoute = "UNROUTABLE_ROUTE"
	// AddressLookupError happens when doing a LatLong lookup of an address
	AddressLookupError = "ADDRESS_LOOKUP_ERROR"
	// GeocodeResponseDecodingError happens when attempting to decode a geocode response
	GeocodeResponseDecodingError = "GEOCODE_RESPONSE_DECODE_ERROR"
	// RoutingResponseDecodingError happens when attempting to decode a routing response
	RoutingResponseDecodingError = "ROUTING_RESPONSE_DECODE_ERROR"
	// UnknownError is for when the cause of the error can't be ascertained
	UnknownError = "UNKNOWN_ERROR"
)

// Error is used for handling errors from the Route package
type Error interface {
	error
	Code() ErrorCode
}

// BaseError contains basic route error functionality
type BaseError struct {
	code ErrorCode
}

// Code returns the error code enum
func (b *BaseError) Code() ErrorCode {
	return b.code
}

type unsupportedPostalCode struct {
	BaseError
	postalCode string
}

// NewUnsupportedPostalCodeError creates a new UnsupportedPostalCode error.
func NewUnsupportedPostalCodeError(postalCode string) Error {
	return &unsupportedPostalCode{
		BaseError{UnsupportedPostalCode},
		postalCode,
	}
}

func (e *unsupportedPostalCode) Error() string {
	return fmt.Sprintf("Unsupported postal code lookup (%s)", e.postalCode)
}

type responseError struct {
	BaseError
	statusCode  int
	routingInfo string
}

func (e *responseError) Error() string {
	return fmt.Sprintf("Error when communicating with HERE server: (error_code: (%s), status_code: %d, routing_info: %s)", e.code, e.statusCode, e.routingInfo)
}

// NewUnroutableRouteError creates a new responseError error.
func NewUnroutableRouteError(statusCode int, source LatLong, dest LatLong) Error {
	return &responseError{
		BaseError{UnroutableRoute},
		statusCode,
		fmt.Sprintf("source: (%s), dest: (%s", source.Coords(), dest.Coords()),
	}
}

// NewUnknownRoutingError returns an error for failed postal code lookups
func NewUnknownRoutingError(statusCode int, source LatLong, dest LatLong) Error {
	return &responseError{
		BaseError{UnknownError},
		statusCode,
		fmt.Sprintf("source: (%s), dest: (%s", source.Coords(), dest.Coords()),
	}
}

// NewAddressLookupError returns a known error for failed address lookups
func NewAddressLookupError(statusCode int, a *models.Address) Error {
	return &responseError{
		BaseError{AddressLookupError},
		statusCode,
		fmt.Sprintf(a.LineFormat()),
	}
}

// NewUnknownAddressLookupError returns an unknown error for failed address lookups
func NewUnknownAddressLookupError(statusCode int, a *models.Address) Error {
	return &responseError{
		BaseError{UnknownError},
		statusCode,
		fmt.Sprintf(a.LineFormat()),
	}
}

// NewPostalCodeLookupError returns an error for failed postal code lookups
func NewPostalCodeLookupError(statusCode int, postalCode string) Error {
	return &responseError{
		BaseError{UnsupportedPostalCode},
		statusCode,
		postalCode,
	}
}

type geocodeResponseDecodingError struct {
	BaseError
	response GeocodeResponseBody
}

func (e *geocodeResponseDecodingError) Error() string {
	return fmt.Sprintf("Error trying to decode GeocodeResponse: %+v", e.response)
}

// NewGeocodeResponseDecodingError creates a new geocodeResponseDecodingError error.
func NewGeocodeResponseDecodingError(r GeocodeResponseBody) Error {
	return &geocodeResponseDecodingError{
		BaseError{GeocodeResponseDecodingError},
		r,
	}
}

type routingResponseDecodingError struct {
	BaseError
	response RoutingResponseBody
}

func (e *routingResponseDecodingError) Error() string {
	return fmt.Sprintf("Error trying to decode RoutingResponseBody: %+v", e.response)
}

// NewRoutingResponseDecodingError creates a new routingResponseDecodingError error.
func NewRoutingResponseDecodingError(r RoutingResponseBody) Error {
	return &routingResponseDecodingError{
		BaseError{RoutingResponseDecodingError},
		r,
	}
}
