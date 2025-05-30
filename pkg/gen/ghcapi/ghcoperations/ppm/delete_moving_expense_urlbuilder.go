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

// DeleteMovingExpenseURL generates an URL for the delete moving expense operation
type DeleteMovingExpenseURL struct {
	MovingExpenseID strfmt.UUID
	PpmShipmentID   strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteMovingExpenseURL) WithBasePath(bp string) *DeleteMovingExpenseURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteMovingExpenseURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *DeleteMovingExpenseURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/ppm-shipments/{ppmShipmentId}/moving-expenses/{movingExpenseId}"

	movingExpenseID := o.MovingExpenseID.String()
	if movingExpenseID != "" {
		_path = strings.Replace(_path, "{movingExpenseId}", movingExpenseID, -1)
	} else {
		return nil, errors.New("movingExpenseId is required on DeleteMovingExpenseURL")
	}

	ppmShipmentID := o.PpmShipmentID.String()
	if ppmShipmentID != "" {
		_path = strings.Replace(_path, "{ppmShipmentId}", ppmShipmentID, -1)
	} else {
		return nil, errors.New("ppmShipmentId is required on DeleteMovingExpenseURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/ghc/v1"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *DeleteMovingExpenseURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *DeleteMovingExpenseURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *DeleteMovingExpenseURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on DeleteMovingExpenseURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on DeleteMovingExpenseURL")
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
func (o *DeleteMovingExpenseURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
