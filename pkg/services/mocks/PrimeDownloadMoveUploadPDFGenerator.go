// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	afero "github.com/spf13/afero"
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"
)

// PrimeDownloadMoveUploadPDFGenerator is an autogenerated mock type for the PrimeDownloadMoveUploadPDFGenerator type
type PrimeDownloadMoveUploadPDFGenerator struct {
	mock.Mock
}

// GenerateDownloadMoveUserUploadPDF provides a mock function with given fields: appCtx, moveOrderUploadType, move, addBookmarks
func (_m *PrimeDownloadMoveUploadPDFGenerator) GenerateDownloadMoveUserUploadPDF(appCtx appcontext.AppContext, moveOrderUploadType services.MoveOrderUploadType, move models.Move, addBookmarks bool) (afero.File, error) {
	ret := _m.Called(appCtx, moveOrderUploadType, move, addBookmarks)

	if len(ret) == 0 {
		panic("no return value specified for GenerateDownloadMoveUserUploadPDF")
	}

	var r0 afero.File
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, services.MoveOrderUploadType, models.Move, bool) (afero.File, error)); ok {
		return rf(appCtx, moveOrderUploadType, move, addBookmarks)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, services.MoveOrderUploadType, models.Move, bool) afero.File); ok {
		r0 = rf(appCtx, moveOrderUploadType, move, addBookmarks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(afero.File)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, services.MoveOrderUploadType, models.Move, bool) error); ok {
		r1 = rf(appCtx, moveOrderUploadType, move, addBookmarks)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPrimeDownloadMoveUploadPDFGenerator creates a new instance of PrimeDownloadMoveUploadPDFGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPrimeDownloadMoveUploadPDFGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *PrimeDownloadMoveUploadPDFGenerator {
	mock := &PrimeDownloadMoveUploadPDFGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
