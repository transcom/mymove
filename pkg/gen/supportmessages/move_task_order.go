// Code generated by go-swagger; DO NOT EDIT.

package supportmessages

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

	// Indicates this MoveTaskOrder has been approved by an office user such as the Task Ordering Officer (TOO).
	//
	// Format: date-time
	ApprovedAt *strfmt.DateTime `json:"approvedAt,omitempty"`

	// Indicates this MoveTaskOrder is available for Prime API handling.
	//
	// In production, only MoveTaskOrders for which this is set will be available to the API.
	//
	// Format: date-time
	AvailableToPrimeAt *strfmt.DateTime `json:"availableToPrimeAt,omitempty"`

	// ID associated with the contractor, in this case Prime
	//
	// Example: 5db13bb4-6d29-4bdb-bc81-262f4513ecf6
	// Required: true
	// Format: uuid
	ContractorID *strfmt.UUID `json:"contractorID"`

	// Date the MoveTaskOrder was created on.
	// Read Only: true
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// Uniquely identifies the state of the MoveTaskOrder object (but not the nested objects)
	//
	// It will change everytime the object is updated. Client should store the value.
	// Updates to this MoveTaskOrder will require that this eTag be passed in with the If-Match header.
	//
	// Read Only: true
	ETag string `json:"eTag,omitempty"`

	// ID of the MoveTaskOrder object.
	// Example: 1f2270c7-7166-40ae-981e-b200ebdf3054
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// Unique 6-character code the customer can use to refer to their move
	// Example: ABC123
	// Read Only: true
	MoveCode string `json:"moveCode,omitempty"`

	mtoServiceItemsField []MTOServiceItem

	// mto shipments
	MtoShipments MTOShipments `json:"mtoShipments,omitempty"`

	// order
	// Required: true
	Order *Order `json:"order"`

	// ID of the Order object
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	OrderID strfmt.UUID `json:"orderID,omitempty"`

	// payment requests
	PaymentRequests PaymentRequests `json:"paymentRequests,omitempty"`

	// If the move is a PPM, this is the estimated weight in lbs.
	PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

	// If the move is a PPM, indicates whether it is full or partial.
	// Enum: [FULL PARTIAL]
	PpmType string `json:"ppmType,omitempty"`

	// Unique ID associated with this Order.
	//
	// No two MoveTaskOrders may have the same ID.
	// Attempting to create a MoveTaskOrder may fail if this referenceId has been used already.
	//
	// Example: 1001-3456
	// Read Only: true
	ReferenceID string `json:"referenceId,omitempty"`

	// status
	Status MoveStatus `json:"status,omitempty"`

	// Date on which this MoveTaskOrder was last updated.
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

		ContractorID *strfmt.UUID `json:"contractorID"`

		CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

		ETag string `json:"eTag,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		MoveCode string `json:"moveCode,omitempty"`

		MtoServiceItems json.RawMessage `json:"mtoServiceItems"`

		MtoShipments MTOShipments `json:"mtoShipments,omitempty"`

		Order *Order `json:"order"`

		OrderID strfmt.UUID `json:"orderID,omitempty"`

		PaymentRequests PaymentRequests `json:"paymentRequests,omitempty"`

		PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

		PpmType string `json:"ppmType,omitempty"`

		ReferenceID string `json:"referenceId,omitempty"`

		Status MoveStatus `json:"status,omitempty"`

		UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var propMtoServiceItems []MTOServiceItem
	if string(data.MtoServiceItems) != "null" {
		mtoServiceItems, err := UnmarshalMTOServiceItemSlice(bytes.NewBuffer(data.MtoServiceItems), runtime.JSONConsumer())
		if err != nil && err != io.EOF {
			return err
		}
		propMtoServiceItems = mtoServiceItems
	}

	var result MoveTaskOrder

	// approvedAt
	result.ApprovedAt = data.ApprovedAt

	// availableToPrimeAt
	result.AvailableToPrimeAt = data.AvailableToPrimeAt

	// contractorID
	result.ContractorID = data.ContractorID

	// createdAt
	result.CreatedAt = data.CreatedAt

	// eTag
	result.ETag = data.ETag

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

	// referenceId
	result.ReferenceID = data.ReferenceID

	// status
	result.Status = data.Status

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

		ContractorID *strfmt.UUID `json:"contractorID"`

		CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

		ETag string `json:"eTag,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		MoveCode string `json:"moveCode,omitempty"`

		MtoShipments MTOShipments `json:"mtoShipments,omitempty"`

		Order *Order `json:"order"`

		OrderID strfmt.UUID `json:"orderID,omitempty"`

		PaymentRequests PaymentRequests `json:"paymentRequests,omitempty"`

		PpmEstimatedWeight int64 `json:"ppmEstimatedWeight,omitempty"`

		PpmType string `json:"ppmType,omitempty"`

		ReferenceID string `json:"referenceId,omitempty"`

		Status MoveStatus `json:"status,omitempty"`

		UpdatedAt strfmt.DateTime `json:"updatedAt,omitempty"`
	}{

		ApprovedAt: m.ApprovedAt,

		AvailableToPrimeAt: m.AvailableToPrimeAt,

		ContractorID: m.ContractorID,

		CreatedAt: m.CreatedAt,

		ETag: m.ETag,

		ID: m.ID,

		MoveCode: m.MoveCode,

		MtoShipments: m.MtoShipments,

		Order: m.Order,

		OrderID: m.OrderID,

		PaymentRequests: m.PaymentRequests,

		PpmEstimatedWeight: m.PpmEstimatedWeight,

		PpmType: m.PpmType,

		ReferenceID: m.ReferenceID,

		Status: m.Status,

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

	if err := m.validateContractorID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
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

	if err := m.validateStatus(formats); err != nil {
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

func (m *MoveTaskOrder) validateContractorID(formats strfmt.Registry) error {

	if err := validate.Required("contractorID", "body", m.ContractorID); err != nil {
		return err
	}

	if err := validate.FormatOf("contractorID", "body", "uuid", m.ContractorID.String(), formats); err != nil {
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
	if swag.IsZero(m.MtoServiceItems()) { // not required
		return nil
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
	if swag.IsZero(m.MtoShipments) { // not required
		return nil
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

	if err := validate.Required("order", "body", m.Order); err != nil {
		return err
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
	if swag.IsZero(m.PaymentRequests) { // not required
		return nil
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
	if err := json.Unmarshal([]byte(`["FULL","PARTIAL"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		moveTaskOrderTypePpmTypePropEnum = append(moveTaskOrderTypePpmTypePropEnum, v)
	}
}

const (

	// MoveTaskOrderPpmTypeFULL captures enum value "FULL"
	MoveTaskOrderPpmTypeFULL string = "FULL"

	// MoveTaskOrderPpmTypePARTIAL captures enum value "PARTIAL"
	MoveTaskOrderPpmTypePARTIAL string = "PARTIAL"
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

func (m *MoveTaskOrder) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("status")
		}
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

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateETag(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
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

	if err := m.contextValidateReferenceID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
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

func (m *MoveTaskOrder) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID)); err != nil {
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

func (m *MoveTaskOrder) contextValidateReferenceID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "referenceId", "body", string(m.ReferenceID)); err != nil {
		return err
	}

	return nil
}

func (m *MoveTaskOrder) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("status")
		}
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