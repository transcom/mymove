// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"
)

// PortLocationFetcher is an autogenerated mock type for the PortLocationFetcher type
type PortLocationFetcher struct {
	mock.Mock
}

// FetchPortLocationByPortCode provides a mock function with given fields: appCtx, portCode
func (_m *PortLocationFetcher) FetchPortLocationByPortCode(appCtx appcontext.AppContext, portCode string) (*models.PortLocation, error) {
	ret := _m.Called(appCtx, portCode)

	if len(ret) == 0 {
		panic("no return value specified for FetchPortLocationByPortCode")
	}

	var r0 *models.PortLocation
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string) (*models.PortLocation, error)); ok {
		return rf(appCtx, portCode)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string) *models.PortLocation); ok {
		r0 = rf(appCtx, portCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PortLocation)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, string) error); ok {
		r1 = rf(appCtx, portCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPortLocationFetcher creates a new instance of PortLocationFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPortLocationFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *PortLocationFetcher {
	mock := &PortLocationFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
