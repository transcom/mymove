// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"
)

// MTOAgentUpdater is an autogenerated mock type for the MTOAgentUpdater type
type MTOAgentUpdater struct {
	mock.Mock
}

// UpdateMTOAgentBasic provides a mock function with given fields: appCtx, mtoAgent, eTag
func (_m *MTOAgentUpdater) UpdateMTOAgentBasic(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	ret := _m.Called(appCtx, mtoAgent, eTag)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMTOAgentBasic")
	}

	var r0 *models.MTOAgent
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent, string) (*models.MTOAgent, error)); ok {
		return rf(appCtx, mtoAgent, eTag)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent, string) *models.MTOAgent); ok {
		r0 = rf(appCtx, mtoAgent, eTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOAgent)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.MTOAgent, string) error); ok {
		r1 = rf(appCtx, mtoAgent, eTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMTOAgentPrime provides a mock function with given fields: appCtx, mtoAgent, eTag
func (_m *MTOAgentUpdater) UpdateMTOAgentPrime(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error) {
	ret := _m.Called(appCtx, mtoAgent, eTag)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMTOAgentPrime")
	}

	var r0 *models.MTOAgent
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent, string) (*models.MTOAgent, error)); ok {
		return rf(appCtx, mtoAgent, eTag)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOAgent, string) *models.MTOAgent); ok {
		r0 = rf(appCtx, mtoAgent, eTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOAgent)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.MTOAgent, string) error); ok {
		r1 = rf(appCtx, mtoAgent, eTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMTOAgentUpdater creates a new instance of MTOAgentUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMTOAgentUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *MTOAgentUpdater {
	mock := &MTOAgentUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
