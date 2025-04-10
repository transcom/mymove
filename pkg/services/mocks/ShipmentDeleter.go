// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	uuid "github.com/gofrs/uuid"
)

// ShipmentDeleter is an autogenerated mock type for the ShipmentDeleter type
type ShipmentDeleter struct {
	mock.Mock
}

// DeleteShipment provides a mock function with given fields: appCtx, shipmentID
func (_m *ShipmentDeleter) DeleteShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (uuid.UUID, error) {
	ret := _m.Called(appCtx, shipmentID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteShipment")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) (uuid.UUID, error)); ok {
		return rf(appCtx, shipmentID)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) uuid.UUID); ok {
		r0 = rf(appCtx, shipmentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, shipmentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShipmentDeleter creates a new instance of ShipmentDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShipmentDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShipmentDeleter {
	mock := &ShipmentDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
