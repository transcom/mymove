// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// TransportaionOfficeAssignmentUpdater is an autogenerated mock type for the TransportaionOfficeAssignmentUpdater type
type TransportaionOfficeAssignmentUpdater struct {
	mock.Mock
}

// UpdateTransportationOfficeAssignments provides a mock function with given fields: appCtx, officeUserId, transportationOfficeAssignments
func (_m *TransportaionOfficeAssignmentUpdater) UpdateTransportationOfficeAssignments(appCtx appcontext.AppContext, officeUserId uuid.UUID, transportationOfficeAssignments models.TransportationOfficeAssignments) (models.TransportationOfficeAssignments, error) {
	ret := _m.Called(appCtx, officeUserId, transportationOfficeAssignments)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTransportationOfficeAssignments")
	}

	var r0 models.TransportationOfficeAssignments
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, models.TransportationOfficeAssignments) (models.TransportationOfficeAssignments, error)); ok {
		return rf(appCtx, officeUserId, transportationOfficeAssignments)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, models.TransportationOfficeAssignments) models.TransportationOfficeAssignments); ok {
		r0 = rf(appCtx, officeUserId, transportationOfficeAssignments)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.TransportationOfficeAssignments)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID, models.TransportationOfficeAssignments) error); ok {
		r1 = rf(appCtx, officeUserId, transportationOfficeAssignments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransportaionOfficeAssignmentUpdater creates a new instance of TransportaionOfficeAssignmentUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransportaionOfficeAssignmentUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransportaionOfficeAssignmentUpdater {
	mock := &TransportaionOfficeAssignmentUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
