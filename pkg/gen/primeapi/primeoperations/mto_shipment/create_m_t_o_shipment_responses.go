// Code generated by go-swagger; DO NOT EDIT.

package mto_shipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// CreateMTOShipmentGoneCode is the HTTP code returned for type CreateMTOShipmentGone
const CreateMTOShipmentGoneCode int = 410

/*
CreateMTOShipmentGone This endpoint is deprecated. Please use `/prime/v3/createMTOShipment` instead.

swagger:response createMTOShipmentGone
*/
type CreateMTOShipmentGone struct {
}

// NewCreateMTOShipmentGone creates CreateMTOShipmentGone with default headers values
func NewCreateMTOShipmentGone() *CreateMTOShipmentGone {

	return &CreateMTOShipmentGone{}
}

// WriteResponse to the client
func (o *CreateMTOShipmentGone) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(410)
}