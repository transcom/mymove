// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SignedCertification Signed certification
//
// swagger:model SignedCertification
type SignedCertification struct {

	// Full text that the customer agreed to and signed.
	// Required: true
	CertificationText *string `json:"certificationText"`

	// certification type
	// Required: true
	CertificationType SignedCertificationType `json:"certificationType"`

	// created at
	// Required: true
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt"`

	// Date that the customer signed the certification.
	// Required: true
	// Format: date
	Date *strfmt.Date `json:"date"`

	// A hash that should be used as the "If-Match" header for any updates.
	// Required: true
	// Read Only: true
	ETag string `json:"eTag"`

	// The ID of the signed certification.
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id"`

	// The ID of the move associated with this signed certification.
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Read Only: true
	// Format: uuid
	MoveID strfmt.UUID `json:"moveId"`

	// The ID of the PPM shipment associated with this signed certification, if any.
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Read Only: true
	// Format: uuid
	PpmID *strfmt.UUID `json:"ppmId"`

	// The signature that the customer provided.
	// Required: true
	Signature *string `json:"signature"`

	// The ID of the user that signed.
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Read Only: true
	// Format: uuid
	SubmittingUserID strfmt.UUID `json:"submittingUserId"`

	// updated at
	// Required: true
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt"`
}

// Validate validates this signed certification
func (m *SignedCertification) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCertificationText(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCertificationType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateETag(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMoveID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePpmID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSignature(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubmittingUserID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SignedCertification) validateCertificationText(formats strfmt.Registry) error {

	if err := validate.Required("certificationText", "body", m.CertificationText); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateCertificationType(formats strfmt.Registry) error {

	if err := validate.Required("certificationType", "body", SignedCertificationType(m.CertificationType)); err != nil {
		return err
	}

	if err := m.CertificationType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("certificationType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("certificationType")
		}
		return err
	}

	return nil
}

func (m *SignedCertification) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateDate(formats strfmt.Registry) error {

	if err := validate.Required("date", "body", m.Date); err != nil {
		return err
	}

	if err := validate.FormatOf("date", "body", "date", m.Date.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateETag(formats strfmt.Registry) error {

	if err := validate.RequiredString("eTag", "body", m.ETag); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateMoveID(formats strfmt.Registry) error {

	if err := validate.Required("moveId", "body", strfmt.UUID(m.MoveID)); err != nil {
		return err
	}

	if err := validate.FormatOf("moveId", "body", "uuid", m.MoveID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validatePpmID(formats strfmt.Registry) error {
	if swag.IsZero(m.PpmID) { // not required
		return nil
	}

	if err := validate.FormatOf("ppmId", "body", "uuid", m.PpmID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateSignature(formats strfmt.Registry) error {

	if err := validate.Required("signature", "body", m.Signature); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateSubmittingUserID(formats strfmt.Registry) error {

	if err := validate.Required("submittingUserId", "body", strfmt.UUID(m.SubmittingUserID)); err != nil {
		return err
	}

	if err := validate.FormatOf("submittingUserId", "body", "uuid", m.SubmittingUserID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this signed certification based on the context it is used
func (m *SignedCertification) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCertificationType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateETag(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMoveID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePpmID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSubmittingUserID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUpdatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SignedCertification) contextValidateCertificationType(ctx context.Context, formats strfmt.Registry) error {

	if err := m.CertificationType.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("certificationType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("certificationType")
		}
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateETag(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "eTag", "body", string(m.ETag)); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateMoveID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "moveId", "body", strfmt.UUID(m.MoveID)); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidatePpmID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "ppmId", "body", m.PpmID); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateSubmittingUserID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "submittingUserId", "body", strfmt.UUID(m.SubmittingUserID)); err != nil {
		return err
	}

	return nil
}

func (m *SignedCertification) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SignedCertification) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SignedCertification) UnmarshalBinary(b []byte) error {
	var res SignedCertification
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}