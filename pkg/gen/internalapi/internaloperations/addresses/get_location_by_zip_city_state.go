// Code generated by go-swagger; DO NOT EDIT.

package addresses

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetLocationByZipCityStateHandlerFunc turns a function with the right signature into a get location by zip city state handler
type GetLocationByZipCityStateHandlerFunc func(GetLocationByZipCityStateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLocationByZipCityStateHandlerFunc) Handle(params GetLocationByZipCityStateParams) middleware.Responder {
	return fn(params)
}

// GetLocationByZipCityStateHandler interface for that can handle valid get location by zip city state params
type GetLocationByZipCityStateHandler interface {
	Handle(GetLocationByZipCityStateParams) middleware.Responder
}

// NewGetLocationByZipCityState creates a new http.Handler for the get location by zip city state operation
func NewGetLocationByZipCityState(ctx *middleware.Context, handler GetLocationByZipCityStateHandler) *GetLocationByZipCityState {
	return &GetLocationByZipCityState{Context: ctx, Handler: handler}
}

/*
	GetLocationByZipCityState swagger:route GET /addresses/zip-city-lookup/{search} addresses getLocationByZipCityState

Returns city, state, postal code, and county associated with the specified full/partial postal code or city state string

Find by API using full/partial postal code or city name that returns an us_post_region_cities json object containing city, state, county and postal code.
*/
type GetLocationByZipCityState struct {
	Context *middleware.Context
	Handler GetLocationByZipCityStateHandler
}

func (o *GetLocationByZipCityState) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetLocationByZipCityStateParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
