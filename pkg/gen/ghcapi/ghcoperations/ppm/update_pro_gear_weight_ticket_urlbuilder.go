// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/strfmt"
)

// UpdateProGearWeightTicketURL generates an URL for the update pro gear weight ticket operation
type UpdateProGearWeightTicketURL struct {
	PpmShipmentID         strfmt.UUID
	ProGearWeightTicketID strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *UpdateProGearWeightTicketURL) WithBasePath(bp string) *UpdateProGearWeightTicketURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *UpdateProGearWeightTicketURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *UpdateProGearWeightTicketURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/ppm-shipments/{ppmShipmentId}/pro-gear-weight-tickets/{proGearWeightTicketId}"

	ppmShipmentID := o.PpmShipmentID.String()
	if ppmShipmentID != "" {
		_path = strings.Replace(_path, "{ppmShipmentId}", ppmShipmentID, -1)
	} else {
		return nil, errors.New("ppmShipmentId is required on UpdateProGearWeightTicketURL")
	}

	proGearWeightTicketID := o.ProGearWeightTicketID.String()
	if proGearWeightTicketID != "" {
		_path = strings.Replace(_path, "{proGearWeightTicketId}", proGearWeightTicketID, -1)
	} else {
		return nil, errors.New("proGearWeightTicketId is required on UpdateProGearWeightTicketURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/ghc/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *UpdateProGearWeightTicketURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *UpdateProGearWeightTicketURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *UpdateProGearWeightTicketURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on UpdateProGearWeightTicketURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on UpdateProGearWeightTicketURL")
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
func (o *UpdateProGearWeightTicketURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}