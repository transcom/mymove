// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"

	time "time"

	unit "github.com/transcom/mymove/pkg/unit"
)

// DomesticUncratingPricer is an autogenerated mock type for the DomesticUncratingPricer type
type DomesticUncratingPricer struct {
	mock.Mock
}

// Price provides a mock function with given fields: appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest
func (_m *DomesticUncratingPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, servicesScheduleDest int) (unit.Cents, services.PricingDisplayParams, error) {
	ret := _m.Called(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest)

	if len(ret) == 0 {
		panic("no return value specified for Price")
	}

	var r0 unit.Cents
	var r1 services.PricingDisplayParams
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) (unit.Cents, services.PricingDisplayParams, error)); ok {
		return rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) unit.Cents); ok {
		r0 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest)
	} else {
		r0 = ret.Get(0).(unit.Cents)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) services.PricingDisplayParams); ok {
		r1 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(services.PricingDisplayParams)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) error); ok {
		r2 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleDest)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PriceUsingParams provides a mock function with given fields: appCtx, params
func (_m *DomesticUncratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	ret := _m.Called(appCtx, params)

	if len(ret) == 0 {
		panic("no return value specified for PriceUsingParams")
	}

	var r0 unit.Cents
	var r1 services.PricingDisplayParams
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error)); ok {
		return rf(appCtx, params)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.PaymentServiceItemParams) unit.Cents); ok {
		r0 = rf(appCtx, params)
	} else {
		r0 = ret.Get(0).(unit.Cents)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, models.PaymentServiceItemParams) services.PricingDisplayParams); ok {
		r1 = rf(appCtx, params)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(services.PricingDisplayParams)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, models.PaymentServiceItemParams) error); ok {
		r2 = rf(appCtx, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewDomesticUncratingPricer creates a new instance of DomesticUncratingPricer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDomesticUncratingPricer(t interface {
	mock.TestingT
	Cleanup(func())
}) *DomesticUncratingPricer {
	mock := &DomesticUncratingPricer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
