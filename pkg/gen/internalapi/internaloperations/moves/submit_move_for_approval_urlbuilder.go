// Code generated by go-swagger; DO NOT EDIT.

package moves

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/strfmt"
)

// SubmitMoveForApprovalURL generates an URL for the submit move for approval operation
type SubmitMoveForApprovalURL struct {
	MoveID strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *SubmitMoveForApprovalURL) WithBasePath(bp string) *SubmitMoveForApprovalURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *SubmitMoveForApprovalURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *SubmitMoveForApprovalURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/moves/{moveId}/submit"

	moveID := o.MoveID.String()
	if moveID != "" {
		_path = strings.Replace(_path, "{moveId}", moveID, -1)
	} else {
		return nil, errors.New("moveId is required on SubmitMoveForApprovalURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/internal"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *SubmitMoveForApprovalURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *SubmitMoveForApprovalURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *SubmitMoveForApprovalURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on SubmitMoveForApprovalURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on SubmitMoveForApprovalURL")
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
func (o *SubmitMoveForApprovalURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}