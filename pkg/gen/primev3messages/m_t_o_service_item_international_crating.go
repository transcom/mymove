// Code generated by go-swagger; DO NOT EDIT.

package primev3messages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// MTOServiceItemInternationalCrating Describes a international crating/uncrating service item subtype of a MTOServiceItem.
//
// swagger:model MTOServiceItemInternationalCrating
type MTOServiceItemInternationalCrating struct {
	eTagField string

	idField strfmt.UUID

	lockedPriceCentsField *int64

	moveTaskOrderIdField *strfmt.UUID

	mtoShipmentIdField strfmt.UUID

	reServiceNameField string

	rejectionReasonField *string

	serviceRequestDocumentsField ServiceRequestDocuments

	statusField MTOServiceItemStatus

	// The dimensions for the crate the item will be shipped in.
	// Required: true
	Crate struct {
		MTOServiceItemDimension
	} `json:"crate"`

	// A description of the item being crated.
	// Example: Decorated horse head to be crated.
	// Required: true
	Description *string `json:"description"`

	// external crate
	ExternalCrate *bool `json:"externalCrate,omitempty"`

	// The dimensions of the item being crated.
	// Required: true
	Item struct {
		MTOServiceItemDimension
	} `json:"item"`

	// To identify whether the service was provided within (CONUS) or (OCONUS)
	// Example: CONUS
	// Enum: [CONUS OCONUS]
	Market string `json:"market,omitempty"`

	// A unique code for the service item. Indicates if the service is for crating (ICRT) or uncrating (IUCRT).
	// Required: true
	// Enum: [ICRT IUCRT]
	ReServiceCode *string `json:"reServiceCode"`

	// The contractor's explanation for why an item needed to be crated or uncrated. Used by the TOO while deciding to approve or reject the service item.
	//
	// Example: Storage items need to be picked up
	Reason *string `json:"reason"`

	// standalone crate
	StandaloneCrate *bool `json:"standaloneCrate,omitempty"`
}

// ETag gets the e tag of this subtype
func (m *MTOServiceItemInternationalCrating) ETag() string {
	return m.eTagField
}

// SetETag sets the e tag of this subtype
func (m *MTOServiceItemInternationalCrating) SetETag(val string) {
	m.eTagField = val
}

// ID gets the id of this subtype
func (m *MTOServiceItemInternationalCrating) ID() strfmt.UUID {
	return m.idField
}

// SetID sets the id of this subtype
func (m *MTOServiceItemInternationalCrating) SetID(val strfmt.UUID) {
	m.idField = val
}

// LockedPriceCents gets the locked price cents of this subtype
func (m *MTOServiceItemInternationalCrating) LockedPriceCents() *int64 {
	return m.lockedPriceCentsField
}

// SetLockedPriceCents sets the locked price cents of this subtype
func (m *MTOServiceItemInternationalCrating) SetLockedPriceCents(val *int64) {
	m.lockedPriceCentsField = val
}

// ModelType gets the model type of this subtype
func (m *MTOServiceItemInternationalCrating) ModelType() MTOServiceItemModelType {
	return "MTOServiceItemInternationalCrating"
}

// SetModelType sets the model type of this subtype
func (m *MTOServiceItemInternationalCrating) SetModelType(val MTOServiceItemModelType) {
}

// MoveTaskOrderID gets the move task order ID of this subtype
func (m *MTOServiceItemInternationalCrating) MoveTaskOrderID() *strfmt.UUID {
	return m.moveTaskOrderIdField
}

// SetMoveTaskOrderID sets the move task order ID of this subtype
func (m *MTOServiceItemInternationalCrating) SetMoveTaskOrderID(val *strfmt.UUID) {
	m.moveTaskOrderIdField = val
}

// MtoShipmentID gets the mto shipment ID of this subtype
func (m *MTOServiceItemInternationalCrating) MtoShipmentID() strfmt.UUID {
	return m.mtoShipmentIdField
}

// SetMtoShipmentID sets the mto shipment ID of this subtype
func (m *MTOServiceItemInternationalCrating) SetMtoShipmentID(val strfmt.UUID) {
	m.mtoShipmentIdField = val
}

// ReServiceName gets the re service name of this subtype
func (m *MTOServiceItemInternationalCrating) ReServiceName() string {
	return m.reServiceNameField
}

// SetReServiceName sets the re service name of this subtype
func (m *MTOServiceItemInternationalCrating) SetReServiceName(val string) {
	m.reServiceNameField = val
}

// RejectionReason gets the rejection reason of this subtype
func (m *MTOServiceItemInternationalCrating) RejectionReason() *string {
	return m.rejectionReasonField
}

// SetRejectionReason sets the rejection reason of this subtype
func (m *MTOServiceItemInternationalCrating) SetRejectionReason(val *string) {
	m.rejectionReasonField = val
}

// ServiceRequestDocuments gets the service request documents of this subtype
func (m *MTOServiceItemInternationalCrating) ServiceRequestDocuments() ServiceRequestDocuments {
	return m.serviceRequestDocumentsField
}

// SetServiceRequestDocuments sets the service request documents of this subtype
func (m *MTOServiceItemInternationalCrating) SetServiceRequestDocuments(val ServiceRequestDocuments) {
	m.serviceRequestDocumentsField = val
}

// Status gets the status of this subtype
func (m *MTOServiceItemInternationalCrating) Status() MTOServiceItemStatus {
	return m.statusField
}

// SetStatus sets the status of this subtype
func (m *MTOServiceItemInternationalCrating) SetStatus(val MTOServiceItemStatus) {
	m.statusField = val
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *MTOServiceItemInternationalCrating) UnmarshalJSON(raw []byte) error {
	var data struct {

		// The dimensions for the crate the item will be shipped in.
		// Required: true
		Crate struct {
			MTOServiceItemDimension
		} `json:"crate"`

		// A description of the item being crated.
		// Example: Decorated horse head to be crated.
		// Required: true
		Description *string `json:"description"`

		// external crate
		ExternalCrate *bool `json:"externalCrate,omitempty"`

		// The dimensions of the item being crated.
		// Required: true
		Item struct {
			MTOServiceItemDimension
		} `json:"item"`

		// To identify whether the service was provided within (CONUS) or (OCONUS)
		// Example: CONUS
		// Enum: [CONUS OCONUS]
		Market string `json:"market,omitempty"`

		// A unique code for the service item. Indicates if the service is for crating (ICRT) or uncrating (IUCRT).
		// Required: true
		// Enum: [ICRT IUCRT]
		ReServiceCode *string `json:"reServiceCode"`

		// The contractor's explanation for why an item needed to be crated or uncrated. Used by the TOO while deciding to approve or reject the service item.
		//
		// Example: Storage items need to be picked up
		Reason *string `json:"reason"`

		// standalone crate
		StandaloneCrate *bool `json:"standaloneCrate,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		ETag string `json:"eTag,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		LockedPriceCents *int64 `json:"lockedPriceCents,omitempty"`

		ModelType MTOServiceItemModelType `json:"modelType"`

		MoveTaskOrderID *strfmt.UUID `json:"moveTaskOrderID"`

		MtoShipmentID strfmt.UUID `json:"mtoShipmentID,omitempty"`

		ReServiceName string `json:"reServiceName,omitempty"`

		RejectionReason *string `json:"rejectionReason,omitempty"`

		ServiceRequestDocuments ServiceRequestDocuments `json:"serviceRequestDocuments,omitempty"`

		Status MTOServiceItemStatus `json:"status,omitempty"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result MTOServiceItemInternationalCrating

	result.eTagField = base.ETag

	result.idField = base.ID

	result.lockedPriceCentsField = base.LockedPriceCents

	if base.ModelType != result.ModelType() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid modelType value: %q", base.ModelType)
	}
	result.moveTaskOrderIdField = base.MoveTaskOrderID

	result.mtoShipmentIdField = base.MtoShipmentID

	result.reServiceNameField = base.ReServiceName

	result.rejectionReasonField = base.RejectionReason

	result.serviceRequestDocumentsField = base.ServiceRequestDocuments

	result.statusField = base.Status

	result.Crate = data.Crate
	result.Description = data.Description
	result.ExternalCrate = data.ExternalCrate
	result.Item = data.Item
	result.Market = data.Market
	result.ReServiceCode = data.ReServiceCode
	result.Reason = data.Reason
	result.StandaloneCrate = data.StandaloneCrate

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m MTOServiceItemInternationalCrating) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// The dimensions for the crate the item will be shipped in.
		// Required: true
		Crate struct {
			MTOServiceItemDimension
		} `json:"crate"`

		// A description of the item being crated.
		// Example: Decorated horse head to be crated.
		// Required: true
		Description *string `json:"description"`

		// external crate
		ExternalCrate *bool `json:"externalCrate,omitempty"`

		// The dimensions of the item being crated.
		// Required: true
		Item struct {
			MTOServiceItemDimension
		} `json:"item"`

		// To identify whether the service was provided within (CONUS) or (OCONUS)
		// Example: CONUS
		// Enum: [CONUS OCONUS]
		Market string `json:"market,omitempty"`

		// A unique code for the service item. Indicates if the service is for crating (ICRT) or uncrating (IUCRT).
		// Required: true
		// Enum: [ICRT IUCRT]
		ReServiceCode *string `json:"reServiceCode"`

		// The contractor's explanation for why an item needed to be crated or uncrated. Used by the TOO while deciding to approve or reject the service item.
		//
		// Example: Storage items need to be picked up
		Reason *string `json:"reason"`

		// standalone crate
		StandaloneCrate *bool `json:"standaloneCrate,omitempty"`
	}{

		Crate: m.Crate,

		Description: m.Description,

		ExternalCrate: m.ExternalCrate,

		Item: m.Item,

		Market: m.Market,

		ReServiceCode: m.ReServiceCode,

		Reason: m.Reason,

		StandaloneCrate: m.StandaloneCrate,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		ETag string `json:"eTag,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		LockedPriceCents *int64 `json:"lockedPriceCents,omitempty"`

		ModelType MTOServiceItemModelType `json:"modelType"`

		MoveTaskOrderID *strfmt.UUID `json:"moveTaskOrderID"`

		MtoShipmentID strfmt.UUID `json:"mtoShipmentID,omitempty"`

		ReServiceName string `json:"reServiceName,omitempty"`

		RejectionReason *string `json:"rejectionReason,omitempty"`

		ServiceRequestDocuments ServiceRequestDocuments `json:"serviceRequestDocuments,omitempty"`

		Status MTOServiceItemStatus `json:"status,omitempty"`
	}{

		ETag: m.ETag(),

		ID: m.ID(),

		LockedPriceCents: m.LockedPriceCents(),

		ModelType: m.ModelType(),

		MoveTaskOrderID: m.MoveTaskOrderID(),

		MtoShipmentID: m.MtoShipmentID(),

		ReServiceName: m.ReServiceName(),

		RejectionReason: m.RejectionReason(),

		ServiceRequestDocuments: m.ServiceRequestDocuments(),

		Status: m.Status(),
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this m t o service item international crating
func (m *MTOServiceItemInternationalCrating) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMoveTaskOrderID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMtoShipmentID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateServiceRequestDocuments(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCrate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateItem(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMarket(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReServiceCode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MTOServiceItemInternationalCrating) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID()) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID().String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateMoveTaskOrderID(formats strfmt.Registry) error {

	if err := validate.Required("moveTaskOrderID", "body", m.MoveTaskOrderID()); err != nil {
		return err
	}

	if err := validate.FormatOf("moveTaskOrderID", "body", "uuid", m.MoveTaskOrderID().String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateMtoShipmentID(formats strfmt.Registry) error {

	if swag.IsZero(m.MtoShipmentID()) { // not required
		return nil
	}

	if err := validate.FormatOf("mtoShipmentID", "body", "uuid", m.MtoShipmentID().String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateServiceRequestDocuments(formats strfmt.Registry) error {

	if swag.IsZero(m.ServiceRequestDocuments()) { // not required
		return nil
	}

	if err := m.ServiceRequestDocuments().Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("serviceRequestDocuments")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("serviceRequestDocuments")
		}
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status()) { // not required
		return nil
	}

	if err := m.Status().Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("status")
		}
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateCrate(formats strfmt.Registry) error {

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateDescription(formats strfmt.Registry) error {

	if err := validate.Required("description", "body", m.Description); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) validateItem(formats strfmt.Registry) error {

	return nil
}

var mTOServiceItemInternationalCratingTypeMarketPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["CONUS","OCONUS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mTOServiceItemInternationalCratingTypeMarketPropEnum = append(mTOServiceItemInternationalCratingTypeMarketPropEnum, v)
	}
}

// property enum
func (m *MTOServiceItemInternationalCrating) validateMarketEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, mTOServiceItemInternationalCratingTypeMarketPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *MTOServiceItemInternationalCrating) validateMarket(formats strfmt.Registry) error {

	if swag.IsZero(m.Market) { // not required
		return nil
	}

	// value enum
	if err := m.validateMarketEnum("market", "body", m.Market); err != nil {
		return err
	}

	return nil
}

var mTOServiceItemInternationalCratingTypeReServiceCodePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ICRT","IUCRT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mTOServiceItemInternationalCratingTypeReServiceCodePropEnum = append(mTOServiceItemInternationalCratingTypeReServiceCodePropEnum, v)
	}
}

// property enum
func (m *MTOServiceItemInternationalCrating) validateReServiceCodeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, mTOServiceItemInternationalCratingTypeReServiceCodePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *MTOServiceItemInternationalCrating) validateReServiceCode(formats strfmt.Registry) error {

	if err := validate.Required("reServiceCode", "body", m.ReServiceCode); err != nil {
		return err
	}

	// value enum
	if err := m.validateReServiceCodeEnum("reServiceCode", "body", *m.ReServiceCode); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this m t o service item international crating based on the context it is used
func (m *MTOServiceItemInternationalCrating) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateETag(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateReServiceName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRejectionReason(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateServiceRequestDocuments(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCrate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateItem(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateETag(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "eTag", "body", string(m.ETag())); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", strfmt.UUID(m.ID())); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateModelType(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ModelType().ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("modelType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("modelType")
		}
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateReServiceName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "reServiceName", "body", string(m.ReServiceName())); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateRejectionReason(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "rejectionReason", "body", m.RejectionReason()); err != nil {
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateServiceRequestDocuments(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ServiceRequestDocuments().ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("serviceRequestDocuments")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("serviceRequestDocuments")
		}
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Status()) { // not required
		return nil
	}

	if err := m.Status().ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("status")
		}
		return err
	}

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateCrate(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *MTOServiceItemInternationalCrating) contextValidateItem(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *MTOServiceItemInternationalCrating) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MTOServiceItemInternationalCrating) UnmarshalBinary(b []byte) error {
	var res MTOServiceItemInternationalCrating
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}