// Code generated by go-swagger; DO NOT EDIT.

package order

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/strfmt"
)

// AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL generates an URL for the acknowledge excess unaccompanied baggage weight risk operation
type AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL struct {
	OrderID strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) WithBasePath(bp string) *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/orders/{orderID}/acknowledge-excess-unaccompanied-baggage-weight-risk"

	orderID := o.OrderID.String()
	if orderID != "" {
		_path = strings.Replace(_path, "{orderID}", orderID, -1)
	} else {
		return nil, errors.New("orderId is required on AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/ghc/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *AcknowledgeExcessUnaccompaniedBaggageWeightRiskURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}