// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"
)

// UserFetcher is an autogenerated mock type for the UserFetcher type
type UserFetcher struct {
	mock.Mock
}

// FetchUser provides a mock function with given fields: appCtx, filters
func (_m *UserFetcher) FetchUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.User, error) {
	ret := _m.Called(appCtx, filters)

	if len(ret) == 0 {
		panic("no return value specified for FetchUser")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, []services.QueryFilter) (models.User, error)); ok {
		return rf(appCtx, filters)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, []services.QueryFilter) models.User); ok {
		r0 = rf(appCtx, filters)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, []services.QueryFilter) error); ok {
		r1 = rf(appCtx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserFetcher creates a new instance of UserFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserFetcher {
	mock := &UserFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
