// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// RejectedOfficeUserFetcherPop is an autogenerated mock type for the RejectedOfficeUserFetcherPop type
type RejectedOfficeUserFetcherPop struct {
	mock.Mock
}

// FetchRejectedOfficeUserByID provides a mock function with given fields: appCtx, id
func (_m *RejectedOfficeUserFetcherPop) FetchRejectedOfficeUserByID(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error) {
	ret := _m.Called(appCtx, id)

	if len(ret) == 0 {
		panic("no return value specified for FetchRejectedOfficeUserByID")
	}

	var r0 models.OfficeUser
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) (models.OfficeUser, error)); ok {
		return rf(appCtx, id)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) models.OfficeUser); ok {
		r0 = rf(appCtx, id)
	} else {
		r0 = ret.Get(0).(models.OfficeUser)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRejectedOfficeUserFetcherPop creates a new instance of RejectedOfficeUserFetcherPop. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRejectedOfficeUserFetcherPop(t interface {
	mock.TestingT
	Cleanup(func())
}) *RejectedOfficeUserFetcherPop {
	mock := &RejectedOfficeUserFetcherPop{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
