// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	afero "github.com/spf13/afero"
	mock "github.com/stretchr/testify/mock"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

	services "github.com/transcom/mymove/pkg/services"
)

// WeightTicketGenerator is an autogenerated mock type for the WeightTicketGenerator type
type WeightTicketGenerator struct {
	mock.Mock
}

// CleanupFile provides a mock function with given fields: weightFile
func (_m *WeightTicketGenerator) CleanupFile(weightFile afero.File) error {
	ret := _m.Called(weightFile)

	if len(ret) == 0 {
		panic("no return value specified for CleanupFile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(afero.File) error); ok {
		r0 = rf(weightFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FillWeightEstimatorPDFForm provides a mock function with given fields: PageValues, fileName
func (_m *WeightTicketGenerator) FillWeightEstimatorPDFForm(PageValues services.WeightEstimatorPages, fileName string) (afero.File, *pdfcpu.PDFInfo, error) {
	ret := _m.Called(PageValues, fileName)

	if len(ret) == 0 {
		panic("no return value specified for FillWeightEstimatorPDFForm")
	}

	var r0 afero.File
	var r1 *pdfcpu.PDFInfo
	var r2 error
	if rf, ok := ret.Get(0).(func(services.WeightEstimatorPages, string) (afero.File, *pdfcpu.PDFInfo, error)); ok {
		return rf(PageValues, fileName)
	}
	if rf, ok := ret.Get(0).(func(services.WeightEstimatorPages, string) afero.File); ok {
		r0 = rf(PageValues, fileName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(afero.File)
		}
	}

	if rf, ok := ret.Get(1).(func(services.WeightEstimatorPages, string) *pdfcpu.PDFInfo); ok {
		r1 = rf(PageValues, fileName)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*pdfcpu.PDFInfo)
		}
	}

	if rf, ok := ret.Get(2).(func(services.WeightEstimatorPages, string) error); ok {
		r2 = rf(PageValues, fileName)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewWeightTicketGenerator creates a new instance of WeightTicketGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWeightTicketGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *WeightTicketGenerator {
	mock := &WeightTicketGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
