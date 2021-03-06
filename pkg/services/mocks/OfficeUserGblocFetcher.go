// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/gofrs/uuid"
)

// OfficeUserGblocFetcher is an autogenerated mock type for the OfficeUserGblocFetcher type
type OfficeUserGblocFetcher struct {
	mock.Mock
}

// FetchGblocForOfficeUser provides a mock function with given fields: id
func (_m *OfficeUserGblocFetcher) FetchGblocForOfficeUser(id uuid.UUID) (string, error) {
	ret := _m.Called(id)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
