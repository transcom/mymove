// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"
)

// PaymentRequestCreator is an autogenerated mock type for the PaymentRequestCreator type
type PaymentRequestCreator struct {
	mock.Mock
}

// CreatePaymentRequestCheck provides a mock function with given fields: appCtx, paymentRequest
func (_m *PaymentRequestCreator) CreatePaymentRequestCheck(appCtx appcontext.AppContext, paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error) {
	ret := _m.Called(appCtx, paymentRequest)

	if len(ret) == 0 {
		panic("no return value specified for CreatePaymentRequestCheck")
	}

	var r0 *models.PaymentRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.PaymentRequest) (*models.PaymentRequest, error)); ok {
		return rf(appCtx, paymentRequest)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.PaymentRequest) *models.PaymentRequest); ok {
		r0 = rf(appCtx, paymentRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PaymentRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.PaymentRequest) error); ok {
		r1 = rf(appCtx, paymentRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPaymentRequestCreator creates a new instance of PaymentRequestCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentRequestCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentRequestCreator {
	mock := &PaymentRequestCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
