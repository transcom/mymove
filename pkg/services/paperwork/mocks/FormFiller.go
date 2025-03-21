// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
	paperwork "github.com/transcom/mymove/pkg/paperwork"
)

// FormFiller is an autogenerated mock type for the FormFiller type
type FormFiller struct {
	mock.Mock
}

// AppendPage provides a mock function with given fields: _a0, _a1, _a2
func (_m *FormFiller) AppendPage(_a0 io.ReadSeeker, _a1 map[string]paperwork.FieldPos, _a2 interface{}) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for AppendPage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.ReadSeeker, map[string]paperwork.FieldPos, interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Output provides a mock function with given fields: _a0
func (_m *FormFiller) Output(_a0 io.Writer) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Output")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFormFiller creates a new instance of FormFiller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFormFiller(t interface {
	mock.TestingT
	Cleanup(func())
}) *FormFiller {
	mock := &FormFiller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
