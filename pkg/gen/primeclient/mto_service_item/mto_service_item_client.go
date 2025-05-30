// Code generated by go-swagger; DO NOT EDIT.

package mto_service_item

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new mto service item API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for mto service item API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateMTOServiceItem(params *CreateMTOServiceItemParams, opts ...ClientOption) (*CreateMTOServiceItemOK, error)

	CreateServiceRequestDocumentUpload(params *CreateServiceRequestDocumentUploadParams, opts ...ClientOption) (*CreateServiceRequestDocumentUploadCreated, error)

	UpdateMTOServiceItem(params *UpdateMTOServiceItemParams, opts ...ClientOption) (*UpdateMTOServiceItemOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
	CreateMTOServiceItem creates m t o service item

	Creates one or more MTOServiceItems. Not all service items may be created, please see details below.

This endpoint supports different body definitions. In the modelType field below, select the modelType corresponding

	to the service item you wish to create and the documentation will update with the new definition.

Upon creation these items are associated with a Move Task Order and an MTO Shipment.
The request must include UUIDs for the MTO and MTO Shipment connected to this service item. Some service item types require
additional service items to be autogenerated when added - all created service items, autogenerated included,
will be returned in the response.

To update a service item, please use [updateMTOServiceItem](#operation/updateMTOServiceItem) endpoint.

---

**`MTOServiceItemOriginSIT`**

MTOServiceItemOriginSIT is a subtype of MTOServiceItem.

This model type describes a domestic origin SIT service item. Items can be created using this
model type with the following codes:

**DOFSIT**

**1st day origin SIT service item**. When a DOFSIT is requested, the API will auto-create the following group of service items:
  - DOFSIT - Domestic origin 1st day SIT
  - DOASIT - Domestic origin Additional day SIT
  - DOPSIT - Domestic origin SIT pickup
  - DOSFSC - Domestic origin SIT fuel surcharge

**DOASIT**

**Addt'l days origin SIT service item**. This represents an additional day of storage for the same item.
Additional DOASIT service items can be created and added to an existing shipment that **includes a DOFSIT service item**.

---

**`MTOServiceItemDestSIT`**

MTOServiceItemDestSIT is a subtype of MTOServiceItem.

This model type describes a domestic destination SIT service item. Items can be created using this
model type with the following codes:

**DDFSIT**

**1st day destination SIT service item**.

These additional fields are optional for creating a DDFSIT:
  - `firstAvailableDeliveryDate1`
  - string <date>
  - First available date that Prime can deliver SIT service item.
  - firstAvailableDeliveryDate1, dateOfContact1, and timeMilitary1 are required together
  - `dateOfContact1`
  - string <date>
  - Date of attempted contact by the prime corresponding to `timeMilitary1`
  - dateOfContact1, timeMilitary1, and firstAvailableDeliveryDate1 are required together
  - `timeMilitary1`
  - string\d{4}Z
  - Time of attempted contact corresponding to `dateOfContact1`, in military format.
  - timeMilitary1, dateOfContact1, and firstAvailableDeliveryDate1 are required together
  - `firstAvailableDeliveryDate2`
  - string <date>
  - Second available date that Prime can deliver SIT service item.
  - firstAvailableDeliveryDate2, dateOfContact2, and timeMilitary2 are required together
  - `dateOfContact2`
  - string <date>
  - Date of attempted contact delivery by the prime corresponding to `timeMilitary2`
  - dateOfContact2, timeMilitary2, and firstAvailableDeliveryDate2 are required together
  - `timeMilitary2`
  - string\d{4}Z
  - Time of attempted contact corresponding to `dateOfContact2`, in military format.
  - timeMilitary2, dateOfContact2, and firstAvailableDeliveryDate2 are required together

When a DDFSIT is requested, the API will auto-create the following group of service items:
  - DDFSIT - Domestic destination 1st day SIT
  - DDASIT - Domestic destination Additional day SIT
  - DDDSIT - Domestic destination SIT delivery
  - DDSFSC - Domestic destination SIT fuel surcharge

**NOTE** When providing the `sitEntryDate` value in the payload, please ensure that the date is not BEFORE
`firstAvailableDeliveryDate1` or `firstAvailableDeliveryDate2`. If it is, you will receive an error response.

**DDASIT**

**Addt'l days destination SIT service item**. This represents an additional day of storage for the same item.
Additional DDASIT service items can be created and added to an existing shipment that **includes a DDFSIT service item**.

---

**`MTOServiceItemInternationalOriginSIT`**

MTOServiceItemInternationalOriginSIT is a subtype of MTOServiceItem.

This model type describes a international origin SIT service item. Items can be created using this
model type with the following codes:

**IOFSIT**

**1st day origin SIT service item**. When a IOFSIT is requested, the API will auto-create the following group of service items:
  - IOFSIT - International origin 1st day SIT
  - IOASIT - International origin Additional day SIT
  - IOPSIT - International origin SIT pickup
  - IOSFSC - International origin SIT fuel surcharge

**IOASIT**

**Addt'l days origin SIT service item**. This represents an additional day of storage for the same item.
Additional IOASIT service items can be created and added to an existing shipment that **includes a IOFSIT service item**.

---

**`MTOServiceItemInternationalDestSIT`**

MTOServiceItemInternationalDestSIT is a subtype of MTOServiceItem.

This model type describes a international destination SIT service item. Items can be created using this
model type with the following codes:

**IDFSIT**

**1st day destination SIT service item**.

These additional fields are optional for creating a IDFSIT:
  - `firstAvailableDeliveryDate1`
  - string <date>
  - First available date that Prime can deliver SIT service item.
  - firstAvailableDeliveryDate1, dateOfContact1, and timeMilitary1 are required together
  - `dateOfContact1`
  - string <date>
  - Date of attempted contact by the prime corresponding to `timeMilitary1`
  - dateOfContact1, timeMilitary1, and firstAvailableDeliveryDate1 are required together
  - `timeMilitary1`
  - string\d{4}Z
  - Time of attempted contact corresponding to `dateOfContact1`, in military format.
  - timeMilitary1, dateOfContact1, and firstAvailableDeliveryDate1 are required together
  - `firstAvailableDeliveryDate2`
  - string <date>
  - Second available date that Prime can deliver SIT service item.
  - firstAvailableDeliveryDate2, dateOfContact2, and timeMilitary2 are required together
  - `dateOfContact2`
  - string <date>
  - Date of attempted contact delivery by the prime corresponding to `timeMilitary2`
  - dateOfContact2, timeMilitary2, and firstAvailableDeliveryDate2 are required together
  - `timeMilitary2`
  - string\d{4}Z
  - Time of attempted contact corresponding to `dateOfContact2`, in military format.
  - timeMilitary2, dateOfContact2, and firstAvailableDeliveryDate2 are required together

When a IDFSIT is requested, the API will auto-create the following group of service items:
  - IDFSIT - International destination 1st day SIT
  - IDASIT - International destination Additional day SIT
  - IDDSIT - International destination SIT delivery
  - IDSFSC - International destination SIT fuel surcharge

**NOTE** When providing the `sitEntryDate` value in the payload, please ensure that the date is not BEFORE
`firstAvailableDeliveryDate1` or `firstAvailableDeliveryDate2`. If it is, you will receive an error response.

**IDASIT**

**Addt'l days destination SIT service item**. This represents an additional day of storage for the same item.
Additional IDASIT service items can be created and added to an existing shipment that **includes a IDFSIT service item**.
*/
func (a *Client) CreateMTOServiceItem(params *CreateMTOServiceItemParams, opts ...ClientOption) (*CreateMTOServiceItemOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateMTOServiceItemParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createMTOServiceItem",
		Method:             "POST",
		PathPattern:        "/mto-service-items",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateMTOServiceItemReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateMTOServiceItemOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createMTOServiceItem: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	CreateServiceRequestDocumentUpload creates service request document upload

	### Functionality

This endpoint **uploads** a Service Request document for a
ServiceItem.

The ServiceItem should already exist.

ServiceItems are created with the
[createMTOServiceItem](#operation/createMTOServiceItem)
endpoint.
*/
func (a *Client) CreateServiceRequestDocumentUpload(params *CreateServiceRequestDocumentUploadParams, opts ...ClientOption) (*CreateServiceRequestDocumentUploadCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateServiceRequestDocumentUploadParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createServiceRequestDocumentUpload",
		Method:             "POST",
		PathPattern:        "/mto-service-items/{mtoServiceItemID}/uploads",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"multipart/form-data"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateServiceRequestDocumentUploadReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateServiceRequestDocumentUploadCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createServiceRequestDocumentUpload: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	UpdateMTOServiceItem updates m t o service item

	Updates MTOServiceItems after creation. Not all service items or fields may be updated, please see details below.

This endpoint supports different body definitions. In the modelType field below, select the modelType corresponding

	to the service item you wish to update and the documentation will update with the new definition.

* Addresses: To update a destination service item's SIT destination final address, update the shipment delivery address.
For approved shipments, please use [updateShipmentDestinationAddress](#mtoShipment/updateShipmentDestinationAddress).
For shipments not yet approved, please use [updateMTOShipmentAddress](#mtoShipment/updateMTOShipmentAddress).

* SIT Service Items: Take note that when updating `sitCustomerContacted`, `sitDepartureDate`, or `sitRequestedDelivery`, we want
those to be updated on `DOASIT` (for origin SIT) and `DDASIT` (for destination SIT). If updating those values in other service
items, the office users will not have as much attention to those values.

To create a service item, please use [createMTOServiceItem](#mtoServiceItem/createMTOServiceItem)) endpoint.

* Resubmitting rejected SIT/Accessorial service items: This endpoint will handle the logic of changing the status of rejected SIT/Accessorial service items from
REJECTED to SUBMITTED. Please provide the `requestedApprovalsRequestedStatus: true` when resubmitting as this will give attention to the TOO to
review the resubmitted SIT/Accessorial service item. Another note, `updateReason` must have a different value than the current `reason` value on the service item.
If this value is not updated, then an error will be sent back.

The following SIT service items can be resubmitted following a rejection:
- DDASIT
- DDDSIT
- DDFSIT
- DOASIT
- DOPSIT
- DOFSIT
- DDSFSC
- DOSFSC
- IDASIT
- IDDSIT
- IDFSIT
- IOASIT
- IOPSIT
- IOFSIT
- IDSFSC
- IOSFSC

The following Accessorial service items can be resubmitted following a rejection:
- IOSHUT
- IDSHUT

At a MINIMUM, the payload for resubmitting a rejected SIT/Accessorial service item must look like this:
```json

	{
	  "reServiceCode": "DDFSIT",
	  "updateReason": "A reason that differs from the previous reason",
	  "modelType": "UpdateMTOServiceItemSIT",
	  "requestApprovalsRequestedStatus": true
	}

```

The following service items allow you to update the Port that the shipment will use:
- PODFSC (Port of Debarkation can be updated)
- POEFSC (Port of Embarkation can be updated)

At a MINIMUM, the payload for updating the port should contain the reServiceCode (PODFSC or POEFSC), modelType (UpdateMTOServiceItemInternationalPortFSC), portCode, and id for the service item.
Please see the example payload below:
```json

	{
	  "id": "1ed224b6-c65e-4616-b88e-8304d26c9562",
	  "modelType": "UpdateMTOServiceItemInternationalPortFSC",
	  "portCode": "SEA",
	  "reServiceCode": "POEFSC"
	}

```

The following crating/uncrating service items can be resubmitted following a rejection:
- ICRT
- IUCRT

At a MINIMUM, the payload for resubmitting a rejected crating/uncrating service item must look like this:
```json

	{
	  "item": {
	    "length": 10000,
	    "width": 10000,
	    "height": 10000
	  },
	  "crate": {
	    "length": 20000,
	    "width": 20000,
	    "height": 20000
	  },
	  "updateReason": "A reason that differs from the previous reason",
	  "modelType": "UpdateMTOServiceItemCrating",
	  "requestApprovalsRequestedStatus": true
	}

```
*/
func (a *Client) UpdateMTOServiceItem(params *UpdateMTOServiceItemParams, opts ...ClientOption) (*UpdateMTOServiceItemOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateMTOServiceItemParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "updateMTOServiceItem",
		Method:             "PATCH",
		PathPattern:        "/mto-service-items/{mtoServiceItemID}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateMTOServiceItemReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UpdateMTOServiceItemOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for updateMTOServiceItem: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
