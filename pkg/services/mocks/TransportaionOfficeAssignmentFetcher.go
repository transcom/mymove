// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// TransportaionOfficeAssignmentFetcher is an autogenerated mock type for the TransportaionOfficeAssignmentFetcher type
type TransportaionOfficeAssignmentFetcher struct {
	mock.Mock
}

// FetchTransportaionOfficeAssignmentsByOfficeUserID provides a mock function with given fields: appCtx, officeUserId
func (_m *TransportaionOfficeAssignmentFetcher) FetchTransportaionOfficeAssignmentsByOfficeUserID(appCtx appcontext.AppContext, officeUserId uuid.UUID) (models.TransportationOfficeAssignments, error) {
	ret := _m.Called(appCtx, officeUserId)

	if len(ret) == 0 {
		panic("no return value specified for FetchTransportaionOfficeAssignmentsByOfficeUserID")
	}

	var r0 models.TransportationOfficeAssignments
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) (models.TransportationOfficeAssignments, error)); ok {
		return rf(appCtx, officeUserId)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) models.TransportationOfficeAssignments); ok {
		r0 = rf(appCtx, officeUserId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.TransportationOfficeAssignments)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, officeUserId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransportaionOfficeAssignmentFetcher creates a new instance of TransportaionOfficeAssignmentFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransportaionOfficeAssignmentFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransportaionOfficeAssignmentFetcher {
	mock := &TransportaionOfficeAssignmentFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
