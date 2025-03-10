// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// MTOShipmentFetcher is an autogenerated mock type for the MTOShipmentFetcher type
type MTOShipmentFetcher struct {
	mock.Mock
}

// GetDiversionChain provides a mock function with given fields: appCtx, shipmentID
func (_m *MTOShipmentFetcher) GetDiversionChain(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*[]models.MTOShipment, error) {
	ret := _m.Called(appCtx, shipmentID)

	if len(ret) == 0 {
		panic("no return value specified for GetDiversionChain")
	}

	var r0 *[]models.MTOShipment
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) (*[]models.MTOShipment, error)); ok {
		return rf(appCtx, shipmentID)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) *[]models.MTOShipment); ok {
		r0 = rf(appCtx, shipmentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.MTOShipment)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, shipmentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShipment provides a mock function with given fields: appCtx, shipmentID, eagerAssociations
func (_m *MTOShipmentFetcher) GetShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eagerAssociations ...string) (*models.MTOShipment, error) {
	_va := make([]interface{}, len(eagerAssociations))
	for _i := range eagerAssociations {
		_va[_i] = eagerAssociations[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, appCtx, shipmentID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetShipment")
	}

	var r0 *models.MTOShipment
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, ...string) (*models.MTOShipment, error)); ok {
		return rf(appCtx, shipmentID, eagerAssociations...)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, ...string) *models.MTOShipment); ok {
		r0 = rf(appCtx, shipmentID, eagerAssociations...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOShipment)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID, ...string) error); ok {
		r1 = rf(appCtx, shipmentID, eagerAssociations...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListMTOShipments provides a mock function with given fields: appCtx, moveID
func (_m *MTOShipmentFetcher) ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) ([]models.MTOShipment, error) {
	ret := _m.Called(appCtx, moveID)

	if len(ret) == 0 {
		panic("no return value specified for ListMTOShipments")
	}

	var r0 []models.MTOShipment
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) ([]models.MTOShipment, error)); ok {
		return rf(appCtx, moveID)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) []models.MTOShipment); ok {
		r0 = rf(appCtx, moveID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.MTOShipment)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, moveID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMTOShipmentFetcher creates a new instance of MTOShipmentFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMTOShipmentFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MTOShipmentFetcher {
	mock := &MTOShipmentFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
