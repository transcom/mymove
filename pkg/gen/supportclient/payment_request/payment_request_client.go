// Code generated by go-swagger; DO NOT EDIT.

package payment_request

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new payment request API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for payment request API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetPaymentRequestEDI(params *GetPaymentRequestEDIParams, opts ...ClientOption) (*GetPaymentRequestEDIOK, error)

	ListMTOPaymentRequests(params *ListMTOPaymentRequestsParams, opts ...ClientOption) (*ListMTOPaymentRequestsOK, error)

	ProcessReviewedPaymentRequests(params *ProcessReviewedPaymentRequestsParams, opts ...ClientOption) (*ProcessReviewedPaymentRequestsOK, error)

	RecalculatePaymentRequest(params *RecalculatePaymentRequestParams, opts ...ClientOption) (*RecalculatePaymentRequestCreated, error)

	UpdatePaymentRequestStatus(params *UpdatePaymentRequestStatusParams, opts ...ClientOption) (*UpdatePaymentRequestStatusOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
	GetPaymentRequestEDI gets payment request e d i

	Returns the EDI (Electronic Data Interchange) message for the payment request identified

by the given payment request ID. Note that the EDI returned in the JSON payload will have where there
would normally be line breaks (due to JSON not allowing line breaks in a string).

This is a support endpoint and will not be available in production.
*/
func (a *Client) GetPaymentRequestEDI(params *GetPaymentRequestEDIParams, opts ...ClientOption) (*GetPaymentRequestEDIOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPaymentRequestEDIParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getPaymentRequestEDI",
		Method:             "GET",
		PathPattern:        "/payment-requests/{paymentRequestID}/edi",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetPaymentRequestEDIReader{formats: a.formats},
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
	success, ok := result.(*GetPaymentRequestEDIOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getPaymentRequestEDI: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	ListMTOPaymentRequests lists m t o payment requests

	### Functionality

This endpoint lists all PaymentRequests associated with a given MoveTaskOrder.

This is a support endpoint and is not available in production.
*/
func (a *Client) ListMTOPaymentRequests(params *ListMTOPaymentRequestsParams, opts ...ClientOption) (*ListMTOPaymentRequestsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListMTOPaymentRequestsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listMTOPaymentRequests",
		Method:             "GET",
		PathPattern:        "/move-task-orders/{moveTaskOrderID}/payment-requests",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListMTOPaymentRequestsReader{formats: a.formats},
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
	success, ok := result.(*ListMTOPaymentRequestsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for listMTOPaymentRequests: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	ProcessReviewedPaymentRequests processes reviewed payment requests

	Updates the status of reviewed payment requests and sends PRs to Syncada if

the SendToSyncada flag is set

This is a support endpoint and will not be available in production.
*/
func (a *Client) ProcessReviewedPaymentRequests(params *ProcessReviewedPaymentRequestsParams, opts ...ClientOption) (*ProcessReviewedPaymentRequestsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewProcessReviewedPaymentRequestsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "processReviewedPaymentRequests",
		Method:             "PATCH",
		PathPattern:        "/payment-requests/process-reviewed",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ProcessReviewedPaymentRequestsReader{formats: a.formats},
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
	success, ok := result.(*ProcessReviewedPaymentRequestsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for processReviewedPaymentRequests: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	RecalculatePaymentRequest recalculates payment request

	Recalculates an existing pending payment request by creating a new payment request for the same service

items but is priced based on the current inputs (weights, dates, etc.). The previously existing payment
request is then deprecated. A link is made between the new and existing payment requests.

This is a support endpoint and will not be available in production.
*/
func (a *Client) RecalculatePaymentRequest(params *RecalculatePaymentRequestParams, opts ...ClientOption) (*RecalculatePaymentRequestCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRecalculatePaymentRequestParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "recalculatePaymentRequest",
		Method:             "POST",
		PathPattern:        "/payment-requests/{paymentRequestID}/recalculate",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &RecalculatePaymentRequestReader{formats: a.formats},
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
	success, ok := result.(*RecalculatePaymentRequestCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for recalculatePaymentRequest: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
	UpdatePaymentRequestStatus updates payment request status

	Updates status of a payment request to REVIEWED, SENT_TO_GEX, TPPS_RECEIVED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, PAID, EDI_ERROR, or DEPRECATED.

A status of REVIEWED can optionally have a `rejectionReason`.

This is a support endpoint and is not available in production.
*/
func (a *Client) UpdatePaymentRequestStatus(params *UpdatePaymentRequestStatusParams, opts ...ClientOption) (*UpdatePaymentRequestStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdatePaymentRequestStatusParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "updatePaymentRequestStatus",
		Method:             "PATCH",
		PathPattern:        "/payment-requests/{paymentRequestID}/status",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdatePaymentRequestStatusReader{formats: a.formats},
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
	success, ok := result.(*UpdatePaymentRequestStatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for updatePaymentRequestStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}