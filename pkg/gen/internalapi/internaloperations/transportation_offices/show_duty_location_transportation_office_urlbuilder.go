// Code generated by go-swagger; DO NOT EDIT.

package transportation_offices

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/strfmt"
)

// ShowDutyLocationTransportationOfficeURL generates an URL for the show duty location transportation office operation
type ShowDutyLocationTransportationOfficeURL struct {
	DutyLocationID strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ShowDutyLocationTransportationOfficeURL) WithBasePath(bp string) *ShowDutyLocationTransportationOfficeURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ShowDutyLocationTransportationOfficeURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *ShowDutyLocationTransportationOfficeURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/duty_locations/{dutyLocationId}/transportation_office"

	dutyLocationID := o.DutyLocationID.String()
	if dutyLocationID != "" {
		_path = strings.Replace(_path, "{dutyLocationId}", dutyLocationID, -1)
	} else {
		return nil, errors.New("dutyLocationId is required on ShowDutyLocationTransportationOfficeURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/internal"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *ShowDutyLocationTransportationOfficeURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *ShowDutyLocationTransportationOfficeURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *ShowDutyLocationTransportationOfficeURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on ShowDutyLocationTransportationOfficeURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on ShowDutyLocationTransportationOfficeURL")
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
func (o *ShowDutyLocationTransportationOfficeURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}