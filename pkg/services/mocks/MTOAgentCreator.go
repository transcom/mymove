// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"
)

// MTOAgentCreator is an autogenerated mock type for the MTOAgentCreator type
type MTOAgentCreator struct {
	mock.Mock
}

// CreateMTOAgentPrime provides a mock function with given fields: appCtx, mtoAgent
func (_m *MTOAgentCreator) CreateMTOAgentPrime(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	ret := _m.Called(appCtx, mtoAgent)

	if len(ret) == 0 {
		panic("no return value specified for CreateMTOAgentPrime")
	}

	var r0 *models.MTOAgent
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent) (*models.MTOAgent, error)); ok {
		return rf(appCtx, mtoAgent)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent) *models.MTOAgent); ok {
		r0 = rf(appCtx, mtoAgent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOAgent)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.MTOAgent) error); ok {
		r1 = rf(appCtx, mtoAgent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMTOAgentCreator creates a new instance of MTOAgentCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMTOAgentCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MTOAgentCreator {
	mock := &MTOAgentCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
