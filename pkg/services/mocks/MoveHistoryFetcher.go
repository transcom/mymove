// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"
)

// MoveHistoryFetcher is an autogenerated mock type for the MoveHistoryFetcher type
type MoveHistoryFetcher struct {
	mock.Mock
}

// FetchMoveHistory provides a mock function with given fields: appCtx, params, useDatabaseProcInstead
func (_m *MoveHistoryFetcher) FetchMoveHistory(appCtx appcontext.AppContext, params *services.FetchMoveHistoryParams, useDatabaseProcInstead bool) (*models.MoveHistory, int64, error) {
	ret := _m.Called(appCtx, params, useDatabaseProcInstead)

	if len(ret) == 0 {
		panic("no return value specified for FetchMoveHistory")
	}

	var r0 *models.MoveHistory
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *services.FetchMoveHistoryParams, bool) (*models.MoveHistory, int64, error)); ok {
		return rf(appCtx, params, useDatabaseProcInstead)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *services.FetchMoveHistoryParams, bool) *models.MoveHistory); ok {
		r0 = rf(appCtx, params, useDatabaseProcInstead)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MoveHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *services.FetchMoveHistoryParams, bool) int64); ok {
		r1 = rf(appCtx, params, useDatabaseProcInstead)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, *services.FetchMoveHistoryParams, bool) error); ok {
		r2 = rf(appCtx, params, useDatabaseProcInstead)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewMoveHistoryFetcher creates a new instance of MoveHistoryFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMoveHistoryFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MoveHistoryFetcher {
	mock := &MoveHistoryFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
