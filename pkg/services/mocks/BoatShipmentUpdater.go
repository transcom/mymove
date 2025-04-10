// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// BoatShipmentUpdater is an autogenerated mock type for the BoatShipmentUpdater type
type BoatShipmentUpdater struct {
	mock.Mock
}

// UpdateBoatShipmentWithDefaultCheck provides a mock function with given fields: appCtx, boatshipment, mtoShipmentID
func (_m *BoatShipmentUpdater) UpdateBoatShipmentWithDefaultCheck(appCtx appcontext.AppContext, boatshipment *models.BoatShipment, mtoShipmentID uuid.UUID) (*models.BoatShipment, error) {
	ret := _m.Called(appCtx, boatshipment, mtoShipmentID)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBoatShipmentWithDefaultCheck")
	}

	var r0 *models.BoatShipment
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.BoatShipment, uuid.UUID) (*models.BoatShipment, error)); ok {
		return rf(appCtx, boatshipment, mtoShipmentID)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.BoatShipment, uuid.UUID) *models.BoatShipment); ok {
		r0 = rf(appCtx, boatshipment, mtoShipmentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.BoatShipment)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.BoatShipment, uuid.UUID) error); ok {
		r1 = rf(appCtx, boatshipment, mtoShipmentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBoatShipmentUpdater creates a new instance of BoatShipmentUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBoatShipmentUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *BoatShipmentUpdater {
	mock := &BoatShipmentUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
