// Code generated by go-swagger; DO NOT EDIT.

package ghcmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// SignedCertificationType The type of signed certification:
//   - PPM_PAYMENT: This is used when the customer has a PPM shipment that they have uploaded their documents for and are
//     ready to submit their documentation for review. When they submit, they will be asked to sign certifying the
//     information is correct.
//   - SHIPMENT: This is used when a customer submits their move with their shipments to be reviewed by office users.
//   - PRE_CLOSEOUT_REVIEWED_PPM_PAYMENT: This is used when a move has a PPM shipment and is set to
//     service-counseling-completed "Submit move details" by service counselor.
//   - CLOSEOUT_REVIEWED_PPM_PAYMENT: This is used when a PPM shipment is reviewed by counselor in close out queue.
//
// swagger:model SignedCertificationType
type SignedCertificationType string

func NewSignedCertificationType(value SignedCertificationType) *SignedCertificationType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated SignedCertificationType.
func (m SignedCertificationType) Pointer() *SignedCertificationType {
	return &m
}

const (

	// SignedCertificationTypePPMPAYMENT captures enum value "PPM_PAYMENT"
	SignedCertificationTypePPMPAYMENT SignedCertificationType = "PPM_PAYMENT"

	// SignedCertificationTypeSHIPMENT captures enum value "SHIPMENT"
	SignedCertificationTypeSHIPMENT SignedCertificationType = "SHIPMENT"

	// SignedCertificationTypePRECLOSEOUTREVIEWEDPPMPAYMENT captures enum value "PRE_CLOSEOUT_REVIEWED_PPM_PAYMENT"
	SignedCertificationTypePRECLOSEOUTREVIEWEDPPMPAYMENT SignedCertificationType = "PRE_CLOSEOUT_REVIEWED_PPM_PAYMENT"

	// SignedCertificationTypeCLOSEOUTREVIEWEDPPMPAYMENT captures enum value "CLOSEOUT_REVIEWED_PPM_PAYMENT"
	SignedCertificationTypeCLOSEOUTREVIEWEDPPMPAYMENT SignedCertificationType = "CLOSEOUT_REVIEWED_PPM_PAYMENT"
)

// for schema
var signedCertificationTypeEnum []interface{}

func init() {
	var res []SignedCertificationType
	if err := json.Unmarshal([]byte(`["PPM_PAYMENT","SHIPMENT","PRE_CLOSEOUT_REVIEWED_PPM_PAYMENT","CLOSEOUT_REVIEWED_PPM_PAYMENT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		signedCertificationTypeEnum = append(signedCertificationTypeEnum, v)
	}
}

func (m SignedCertificationType) validateSignedCertificationTypeEnum(path, location string, value SignedCertificationType) error {
	if err := validate.EnumCase(path, location, value, signedCertificationTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this signed certification type
func (m SignedCertificationType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateSignedCertificationTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this signed certification type based on the context it is used
func (m SignedCertificationType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := validate.ReadOnly(ctx, "", "body", SignedCertificationType(m)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}