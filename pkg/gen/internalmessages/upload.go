// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Upload An uploaded file.
//
// swagger:model Upload
type Upload struct {

	// bytes
	// Required: true
	// Read Only: true
	Bytes int64 `json:"bytes"`

	// content type
	// Example: application/pdf
	// Required: true
	// Read Only: true
	ContentType string `json:"contentType"`

	// created at
	// Required: true
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt"`

	// deleted at
	// Read Only: true
	// Format: date-time
	DeletedAt *strfmt.DateTime `json:"deletedAt,omitempty"`

	// filename
	// Example: filename.pdf
	// Required: true
	// Read Only: true
	Filename string `json:"filename"`

	// id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id"`

	// is weight ticket
	IsWeightTicket bool `json:"isWeightTicket,omitempty"`

	// rotation
	// Example: 2
	Rotation int64 `json:"rotation,omitempty"`

	// status
	// Read Only: true
	// Enum: [INFECTED CLEAN PROCESSING]
	Status string `json:"status,omitempty"`

	// updated at
	// Required: true
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt"`

	// upload type
	// Example: OFFICE
	// Read Only: true
	// Enum: [USER PRIME OFFICE]
	UploadType string `json:"uploadType,omitempty"`

	// url
	// Example: https://uploads.domain.test/dir/c56a4180-65aa-42ec-a945-5fd21dec0538
	// Required: true
	// Read Only: true
	// Format: uri
	URL strfmt.URI `json:"url"`
}

// Validate validates this upload
func (m *Upload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBytes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateContentType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDeletedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFilename(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUploadType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURL(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Upload) validateBytes(formats strfmt.Registry) error {

	if err := validate.Required("bytes", "body", int64(m.Bytes)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateContentType(formats strfmt.Registry) error {

	if err := validate.RequiredString("contentType", "body", m.ContentType); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateCreatedAt(formats strfmt.Registry) error {

	if err := validate.Required("createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateDeletedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.DeletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("deletedAt", "body", "date-time", m.DeletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateFilename(formats strfmt.Registry) error {

	if err := validate.RequiredString("filename", "body", m.Filename); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

var uploadTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["INFECTED","CLEAN","PROCESSING"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		uploadTypeStatusPropEnum = append(uploadTypeStatusPropEnum, v)
	}
}

const (

	// UploadStatusINFECTED captures enum value "INFECTED"
	UploadStatusINFECTED string = "INFECTED"

	// UploadStatusCLEAN captures enum value "CLEAN"
	UploadStatusCLEAN string = "CLEAN"

	// UploadStatusPROCESSING captures enum value "PROCESSING"
	UploadStatusPROCESSING string = "PROCESSING"
)

// prop value enum
func (m *Upload) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, uploadTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Upload) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

var uploadTypeUploadTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["USER","PRIME","OFFICE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		uploadTypeUploadTypePropEnum = append(uploadTypeUploadTypePropEnum, v)
	}
}

const (

	// UploadUploadTypeUSER captures enum value "USER"
	UploadUploadTypeUSER string = "USER"

	// UploadUploadTypePRIME captures enum value "PRIME"
	UploadUploadTypePRIME string = "PRIME"

	// UploadUploadTypeOFFICE captures enum value "OFFICE"
	UploadUploadTypeOFFICE string = "OFFICE"
)

// prop value enum
func (m *Upload) validateUploadTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, uploadTypeUploadTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Upload) validateUploadType(formats strfmt.Registry) error {
	if swag.IsZero(m.UploadType) { // not required
		return nil
	}

	// value enum
	if err := m.validateUploadTypeEnum("uploadType", "body", m.UploadType); err != nil {
		return err
	}

	return nil
}

func (m *Upload) validateURL(formats strfmt.Registry) error {

	if err := validate.Required("url", "body", strfmt.URI(m.URL)); err != nil {
		return err
	}

	if err := validate.FormatOf("url", "body", "uri", m.URL.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this upload based on the context it is used
func (m *Upload) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBytes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateContentType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDeletedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateFilename(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUpdatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUploadType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateURL(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Upload) contextValidateBytes(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "bytes", "body", int64(m.Bytes)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateContentType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "contentType", "body", string(m.ContentType)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateDeletedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "deletedAt", "body", m.DeletedAt); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateFilename(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "filename", "body", string(m.Filename)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "status", "body", string(m.Status)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateUploadType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "uploadType", "body", string(m.UploadType)); err != nil {
		return err
	}

	return nil
}

func (m *Upload) contextValidateURL(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "url", "body", strfmt.URI(m.URL)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Upload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Upload) UnmarshalBinary(b []byte) error {
	var res Upload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
