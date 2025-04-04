// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"
)

// PaymentRequestReviewedProcessor is an autogenerated mock type for the PaymentRequestReviewedProcessor type
type PaymentRequestReviewedProcessor struct {
	mock.Mock
}

// ProcessAndLockReviewedPR provides a mock function with given fields: appCtx, pr
func (_m *PaymentRequestReviewedProcessor) ProcessAndLockReviewedPR(appCtx appcontext.AppContext, pr models.PaymentRequest) error {
	ret := _m.Called(appCtx, pr)

	if len(ret) == 0 {
		panic("no return value specified for ProcessAndLockReviewedPR")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.PaymentRequest) error); ok {
		r0 = rf(appCtx, pr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessReviewedPaymentRequest provides a mock function with given fields: appCtx
func (_m *PaymentRequestReviewedProcessor) ProcessReviewedPaymentRequest(appCtx appcontext.AppContext) {
	_m.Called(appCtx)
}

// NewPaymentRequestReviewedProcessor creates a new instance of PaymentRequestReviewedProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentRequestReviewedProcessor(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentRequestReviewedProcessor {
	mock := &PaymentRequestReviewedProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
