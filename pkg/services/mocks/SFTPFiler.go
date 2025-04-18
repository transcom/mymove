// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// SFTPFiler is an autogenerated mock type for the SFTPFiler type
type SFTPFiler struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *SFTPFiler) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteTo provides a mock function with given fields: w
func (_m *SFTPFiler) WriteTo(w io.Writer) (int64, error) {
	ret := _m.Called(w)

	if len(ret) == 0 {
		panic("no return value specified for WriteTo")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(io.Writer) (int64, error)); ok {
		return rf(w)
	}
	if rf, ok := ret.Get(0).(func(io.Writer) int64); ok {
		r0 = rf(w)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(io.Writer) error); ok {
		r1 = rf(w)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSFTPFiler creates a new instance of SFTPFiler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSFTPFiler(t interface {
	mock.TestingT
	Cleanup(func())
}) *SFTPFiler {
	mock := &SFTPFiler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
