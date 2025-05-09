// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// PPMDocumentFetcher is an autogenerated mock type for the PPMDocumentFetcher type
type PPMDocumentFetcher struct {
	mock.Mock
}

// GetPPMDocuments provides a mock function with given fields: appCtx, mtoShipmentID
func (_m *PPMDocumentFetcher) GetPPMDocuments(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.PPMDocuments, error) {
	ret := _m.Called(appCtx, mtoShipmentID)

	if len(ret) == 0 {
		panic("no return value specified for GetPPMDocuments")
	}

	var r0 *models.PPMDocuments
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) (*models.PPMDocuments, error)); ok {
		return rf(appCtx, mtoShipmentID)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID) *models.PPMDocuments); ok {
		r0 = rf(appCtx, mtoShipmentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PPMDocuments)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID) error); ok {
		r1 = rf(appCtx, mtoShipmentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPPMDocumentFetcher creates a new instance of PPMDocumentFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPPMDocumentFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *PPMDocumentFetcher {
	mock := &PPMDocumentFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
