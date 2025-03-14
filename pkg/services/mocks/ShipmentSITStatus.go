// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"
)

// ShipmentSITStatus is an autogenerated mock type for the ShipmentSITStatus type
type ShipmentSITStatus struct {
	mock.Mock
}

// CalculateShipmentSITAllowance provides a mock function with given fields: appCtx, shipment
func (_m *ShipmentSITStatus) CalculateShipmentSITAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment) (int, error) {
	ret := _m.Called(appCtx, shipment)

	if len(ret) == 0 {
		panic("no return value specified for CalculateShipmentSITAllowance")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) (int, error)); ok {
		return rf(appCtx, shipment)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) int); ok {
		r0 = rf(appCtx, shipment)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, models.MTOShipment) error); ok {
		r1 = rf(appCtx, shipment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CalculateShipmentSITStatus provides a mock function with given fields: appCtx, shipment
func (_m *ShipmentSITStatus) CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*services.SITStatus, models.MTOShipment, error) {
	ret := _m.Called(appCtx, shipment)

	if len(ret) == 0 {
		panic("no return value specified for CalculateShipmentSITStatus")
	}

	var r0 *services.SITStatus
	var r1 models.MTOShipment
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) (*services.SITStatus, models.MTOShipment, error)); ok {
		return rf(appCtx, shipment)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) *services.SITStatus); ok {
		r0 = rf(appCtx, shipment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.SITStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, models.MTOShipment) models.MTOShipment); ok {
		r1 = rf(appCtx, shipment)
	} else {
		r1 = ret.Get(1).(models.MTOShipment)
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, models.MTOShipment) error); ok {
		r2 = rf(appCtx, shipment)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CalculateShipmentsSITStatuses provides a mock function with given fields: appCtx, shipments
func (_m *ShipmentSITStatus) CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]services.SITStatus {
	ret := _m.Called(appCtx, shipments)

	if len(ret) == 0 {
		panic("no return value specified for CalculateShipmentsSITStatuses")
	}

	var r0 map[string]services.SITStatus
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, []models.MTOShipment) map[string]services.SITStatus); ok {
		r0 = rf(appCtx, shipments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]services.SITStatus)
		}
	}

	return r0
}

// RetrieveShipmentSIT provides a mock function with given fields: appCtx, shipment
func (_m *ShipmentSITStatus) RetrieveShipmentSIT(appCtx appcontext.AppContext, shipment models.MTOShipment) (models.SITServiceItemGroupings, error) {
	ret := _m.Called(appCtx, shipment)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveShipmentSIT")
	}

	var r0 models.SITServiceItemGroupings
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) (models.SITServiceItemGroupings, error)); ok {
		return rf(appCtx, shipment)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.MTOShipment) models.SITServiceItemGroupings); ok {
		r0 = rf(appCtx, shipment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.SITServiceItemGroupings)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, models.MTOShipment) error); ok {
		r1 = rf(appCtx, shipment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShipmentSITStatus creates a new instance of ShipmentSITStatus. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShipmentSITStatus(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShipmentSITStatus {
	mock := &ShipmentSITStatus{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
