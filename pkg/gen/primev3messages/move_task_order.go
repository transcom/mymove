// Code generated by go-swagger; DO NOT EDIT.

package primev3messages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MoveTaskOrder move task order
//
// swagger:model MoveTaskOrder
type MoveTaskOrder struct {

	// approved at
	// Read Only: true
	// Format: date-time
	ApprovedAt *strfmt.DateTime `json:"approvedAt,omitempty"`

	// available to prime at
	// Read Only: true
	// Format: date-time
	AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

	// contract number
	// Read Only: true
	ContractNumber string `json:"contractNumber,omitempty"`

	// created at
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// e tag
	// Read Only: true
	ETag string `json:"eTag,omitempty"`

	// excess weight acknowledged at
	// Read Only: true
	// Format: date-time
	ExcessWeightAcknowledgedAt *strfmt.DateTime `json:"excessWeightAcknowledgedAt"`

	// excess weight qualified at
	// Read Only: true
	// Format: date-time
	ExcessWeightQualifiedAt *strfmt.DateTime `json:"excessWeightQualifiedAt"`

	// excess weight upload Id
	// Read Only: true
	// Format: uuid
	ExcessWeightUploadID *strfmt.UUID `json:"excessWeightUploadId"`

	// id
	// Example: a502b4f1-b9c4-4faf-8bdd-68292501bf26
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// move code
	// Example: HYXFJF
	// Read Only: true
	MoveCode string `json:"moveCode,omitempty"`

	mtoServiceItemsField []MTOServiceItem

	// mto shipments
	// Required: true
	MtoShipments MTOShipmentsWithoutServiceObjects `json:"mtoShipments"`

	// order
	Order *Order `json:"order,omitempty"`

	// order ID
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	OrderID strfmt.UUID `json:"orderID,omitempty"`

	// payment requests
	// Required: true
	PaymentRequests PaymentRequests `json:"paymentRequests"`

	// ppm estimated weight
	PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

	// ppm type
	// Enum: [PARTIAL FULL]
	PpmType string `json:"ppmType,omitempty"`

	// prime counseling completed at
	// Read Only: true
	// Format: date-time
	PrimeCounselingCompletedAt *strfmt.DateTime `json:"primeCounselingCompletedAt,omitempty"`

	// reference Id
	// Example: 1001-3456
	ReferenceID string `json:"referenceId,omitempty"`

	// updated at
	// Read Only: true
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
}

// MtoServiceItems gets the mto service items of this base type
func (m *MoveTaskOrder) MtoServiceItems() []MTOServiceItem {
	return m.mtoServiceItemsField
}

// SetMtoServiceItems sets the mto service items of this base type
func (m *MoveTaskOrder) SetMtoServiceItems(val []MTOServiceItem) {
	m.mtoServiceItemsField = val
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *MoveTaskOrder) UnmarshalJSON(raw []byte) error {
	var data struct {
		ApprovedAt *strfmt.DateTime `json:"approvedAt,omitempty"`

		AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

		ContractNumber string `json:"contractNumber,omitempty"`

		CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

		ETag string `json:"eTag,omitempty"`

		ExcessWeightAcknowledgedAt *strfmt.DateTime `json:"excessWeightAcknowledgedAt"`

		ExcessWeightQualifiedAt *strfmt.DateTime `json:"excessWeightQualifiedAt"`

		ExcessWeightUploadID *strfmt.UUID `json:"excessWeightUploadId"`

		ID strfmt.UUID `json:"id,omitempty"`

		MoveCode string `json:"moveCode,omitempty"`

		MtoServiceItems json.RawMessage `json:"mtoServiceItems"`

		MtoShipments MTOShipmentsWithoutServiceObjects `json:"mtoShipments"`

		Order *Order `json:"order,omitempty"`

		OrderID strfmt.UUID `json:"orderID,omitempty"`

		PaymentRequests PaymentRequests `json:"paymentRequests"`

		PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

		PpmType string `json:"ppmType,omitempty"`

		PrimeCounselingCompletedAt *strfmt.DateTime `json:"primeCounselingCompletedAt,omitempty"`

		ReferenceID string `json:"referenceId,omitempty"`

		UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	propMtoServiceItems, err := UnmarshalMTOServiceItemSlice(bytes.NewBuffer(data.MtoServiceItems), runtime.JSONConsumer())
	if err != nil && err != io.EOF {
		return err
	}

	var result MoveTaskOrder

	// approvedAt
	result.ApprovedAt = data.ApprovedAt

	// availableToPrimeAt
	result.AvailableToPrimeAt = data.AvailableToPrimeAt

	// contractNumber
	result.ContractNumber = data.ContractNumber

	// createdAt
	result.CreatedAt = data.CreatedAt

	// eTag
	result.ETag = data.ETag

	// excessWeightAcknowledgedAt
	result.ExcessWeightAcknowledgedAt = data.ExcessWeightAcknowledgedAt

	// excessWeightQualifiedAt
	result.ExcessWeightQualifiedAt = data.ExcessWeightQualifiedAt

	// excessWeightUploadId
	result.ExcessWeightUploadID = data.ExcessWeightUploadID

	// id
	result.ID = data.ID

	// moveCode
	result.MoveCode = data.MoveCode

	// mtoServiceItems
	result.mtoServiceItemsField = propMtoServiceItems

	// mtoShipments
	result.MtoShipments = data.MtoShipments

	// order
	result.Order = data.Order

	// orderID
	result.OrderID = data.OrderID

	// paymentRequests
	result.PaymentRequests = data.PaymentRequests

	// ppmEstimatedWeight
	result.PpmEstimatedWeight = data.PpmEstimatedWeight

	// ppmType
	result.PpmType = data.PpmType

	// primeCounselingCompletedAt
	result.PrimeCounselingCompletedAt = data.PrimeCounselingCompletedAt

	// referenceId
	result.ReferenceID = data.ReferenceID

	// updatedAt
	result.UpdatedAt = data.UpdatedAt

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m MoveTaskOrder) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {
		ApprovedAt *strfmt.DateTime `json:"approvedAt,omitempty"`

		AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

		ContractNumber string `json:"contractNumber,omitempty"`

		CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

		ETag string `json:"eTag,omitempty"`

		ExcessWeightAcknowledgedAt *strfmt.DateTime `json:"excessWeightAcknowledgedAt"`

		ExcessWeightQualifiedAt *strfmt.DateTime `json:"excessWeightQualifiedAt"`

		ExcessWeightUploadID *strfmt.UUID `json:"excessWeightUploadId"`

		ID strfmt.UUID `json:"id,omitempty"`

		MoveCode string `json:"moveCode,omitempty"`

		MtoShipments MTOShipmentsWithoutServiceObjects `json:"mtoShipments"`

		Order *Order `json:"order,omitempty"`

		OrderID strfmt.UUID `json:"orderID,omitempty"`

		PaymentRequests PaymentRequests `json:"paymentRequests"`

		PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

		PpmType string `json:"ppmType,omitempty"`

		PrimeCounselingCompletedAt *strfmt.DateTime `json:"primeCounselingCompletedAt,omitempty"`

		ReferenceID string `json:"referenceId,omitempty"`

		UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
	}{

		ApprovedAt: m.ApprovedAt,

		AvailableToPrimeAt: m.AvailableToPrimeAt,

		ContractNumber: m.ContractNumber,

		CreatedAt: m.CreatedAt,

		ETag: m.ETag,

		ExcessWeightAcknowledgedAt: m.ExcessWeightAcknowledgedAt,

		ExcessWeightQualifiedAt: m.ExcessWeightQualifiedAt,

		ExcessWeightUploadID: m.ExcessWeightUploadID,

		ID: m.ID,

		MoveCode: m.MoveCode,

		MtoShipments: m.MtoShipments,

		Order: m.Order,

		OrderID: m.OrderID,

		PaymentRequests: m.PaymentRequests,

		PpmEstimatedWeight: m.PpmEstimatedWeight,

		PpmType: m.PpmType,

		PrimeCounselingCompletedAt: m.PrimeCounselingCompletedAt,

		ReferenceID: m.ReferenceID,

		UpdatedAt: m.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		MtoServiceItems []MTOServiceItem `json:"mtoServiceItems"`
	}{

		MtoServiceItems: m.mtoServiceItemsField,
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this move task order
func (m *MoveTaskOrder) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApprovedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAvailableToPrimeAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExcessWeightAcknowledgedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExcessWeightQualifiedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExcessWeightUploadID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMtoServiceItems(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMtoShipments(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOrder(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOrderID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePaymentRequests(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePpmType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrimeCounselingCompletedAt(formats); err != nil {
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

func (m *MoveTaskOrder) validateApprovedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ApprovedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("approvedAt", "body", "date-time", m.ApprovedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateAvailableToPrimeAt(formats strfmt.Registry) error {
	if swag.IsZero(m.AvailableToPrimeAt) { // not required
		return nil
	}

	if err := validate.FormatOf("availableToPrimeAt", "body", "date-time", m.AvailableToPrimeAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateExcessWeightAcknowledgedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ExcessWeightAcknowledgedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("excessWeightAcknowledgedAt", "body", "date-time", m.ExcessWeightAcknowledgedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateExcessWeightQualifiedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ExcessWeightQualifiedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("excessWeightQualifiedAt", "body", "date-time", m.ExcessWeightQualifiedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateExcessWeightUploadID(formats strfmt.Registry) error {
	if swag.IsZero(m.ExcessWeightUploadID) { // not required
		return nil
	}

	if err := validate.FormatOf("excessWeightUploadId", "body", "uuid", m.ExcessWeightUploadID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateMtoServiceItems(formats strfmt.Registry) error {

	if err := validate.Required("mtoServiceItems", "body", m.MtoServiceItems()); err != nil {
		return err
	}

	for i := 0; i < len(m.MtoServiceItems()); i++ {

		if err := m.mtoServiceItemsField[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mtoServiceItems" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("mtoServiceItems" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *MoveTaskOrder) validateMtoShipments(formats strfmt.Registry) error {

	if err := validate.Required("mtoShipments", "body", m.MtoShipments); err != nil {
		return err
	}

	if err := m.MtoShipments.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("mtoShipments")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("mtoShipments")
		}
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateOrder(formats strfmt.Registry) error {
	if swag.IsZero(m.Order) { // not required
		return nil
	}

	if m.Order != nil {
		if err := m.Order.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("order")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("order")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) validateOrderID(formats strfmt.Registry) error {
	if swag.IsZero(m.OrderID) { // not required
		return nil
	}

	if err := validate.FormatOf("orderID", "body", "uuid", m.OrderID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validatePaymentRequests(formats strfmt.Registry) error {

	if err := validate.Required("paymentRequests", "body", m.PaymentRequests); err != nil {
		return err
	}

	if err := m.PaymentRequests.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("paymentRequests")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("paymentRequests")
		}
		return err
	}

	return nil
}

var moveTaskOrderTypePpmTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["PARTIAL","FULL"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		moveTaskOrderTypePpmTypePropEnum = append(moveTaskOrderTypePpmTypePropEnum, v)
	}
}

const (

	// MoveTaskOrderPpmTypePARTIAL captures enum value "PARTIAL"
	MoveTaskOrderPpmTypePARTIAL string = "PARTIAL"

	// MoveTaskOrderPpmTypeFULL captures enum value "FULL"
	MoveTaskOrderPpmTypeFULL string = "FULL"
)

// prop value enum
func (m *MoveTaskOrder) validatePpmTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, moveTaskOrderTypePpmTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *MoveTaskOrder) validatePpmType(formats strfmt.Registry) error {
	if swag.IsZero(m.PpmType) { // not required
		return nil
	}

	// value enum
	if err := m.validatePpmTypeEnum("ppmType", "body", m.PpmType); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validatePrimeCounselingCompletedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.PrimeCounselingCompletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("primeCounselingCompletedAt", "body", "date-time", m.PrimeCounselingCompletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this move task order based on the context it is used
func (m *MoveTaskOrder) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApprovedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateAvailableToPrimeAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateContractNumber(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateETag(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExcessWeightAcknowledgedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExcessWeightQualifiedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExcessWeightUploadID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMoveCode(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMtoServiceItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMtoShipments(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOrder(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePaymentRequests(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePrimeCounselingCompletedAt(ctx, formats); err != nil {
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

func (m *MoveTaskOrder) contextValidateApprovedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "approvedAt", "body", m.ApprovedAt); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateAvailableToPrimeAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "availableToPrimeAt", "body", m.AvailableToPrimeAt); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateContractNumber(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "contractNumber", "body", string(m.ContractNumber)); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdAt", "body", strfmt.DateTime(m.CreatedAt)); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateETag(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "eTag", "body", string(m.ETag)); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateExcessWeightAcknowledgedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "excessWeightAcknowledgedAt", "body", m.ExcessWeightAcknowledgedAt); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateExcessWeightQualifiedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "excessWeightQualifiedAt", "body", m.ExcessWeightQualifiedAt); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateExcessWeightUploadID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "excessWeightUploadId", "body", m.ExcessWeightUploadID); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateMoveCode(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "moveCode", "body", string(m.MoveCode)); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateMtoServiceItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.MtoServiceItems()); i++ {

		if swag.IsZero(m.mtoServiceItemsField[i]) { // not required
			return nil
		}

		if err := m.mtoServiceItemsField[i].ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mtoServiceItems" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("mtoServiceItems" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *MoveTaskOrder) contextValidateMtoShipments(ctx context.Context, formats strfmt.Registry) error {

	if err := m.MtoShipments.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("mtoShipments")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("mtoShipments")
		}
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateOrder(ctx context.Context, formats strfmt.Registry) error {

	if m.Order != nil {

		if swag.IsZero(m.Order) { // not required
			return nil
		}

		if err := m.Order.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("order")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("order")
			}
			return err
		}
	}

	return nil
}

func (m *MoveTaskOrder) contextValidatePaymentRequests(ctx context.Context, formats strfmt.Registry) error {

	if err := m.PaymentRequests.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("paymentRequests")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("paymentRequests")
		}
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidatePrimeCounselingCompletedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "primeCounselingCompletedAt", "body", m.PrimeCounselingCompletedAt); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateUpdatedAt(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "updatedAt", "body", strfmt.DateTime(m.UpdatedAt)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *MoveTaskOrder) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MoveTaskOrder) UnmarshalBinary(b []byte) error {
	var res MoveTaskOrder
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}