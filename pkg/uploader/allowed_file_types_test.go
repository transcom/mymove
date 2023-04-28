package uploader_test

import (
	"fmt"

	"github.com/transcom/mymove/pkg/uploader"
)

type fileTypeTestCase struct {
	fileType string
	allowed  bool
}

var fileTypeTestCasesAny = []fileTypeTestCase{
	{"image/pdf", true},
	{"image/painting", true},
	{"application/zip", true},
	{uploader.FileTypeText, true},
	{"image/gif", true},
}

var fileTypeTestCasesServiceMember = []fileTypeTestCase{
	{uploader.FileTypeJPEG, true},
	{uploader.FileTypePNG, true},
	{"image/gif", false},
	{uploader.FileTypePDF, true},
	{"application/zip", false},
	{uploader.FileTypeText, false},
}

var fileTypeTestCasesPPM = []fileTypeTestCase{
	{uploader.FileTypeJPEG, true},
	{uploader.FileTypePNG, true},
	{"image/gif", false},
	{uploader.FileTypePDF, true},
	{uploader.FileTypeExcel, true},
	{uploader.FileTypeExcelXLSX, true},
	{"application/zip", false},
	{uploader.FileTypeText, false},
}

var fileTypeTestCasesPaymentRequest = []fileTypeTestCase{
	{uploader.FileTypeJPEG, true},
	{uploader.FileTypePNG, true},
	{"image/gif", false},
	{uploader.FileTypePDF, true},
	{"application/zip", false},
	{uploader.FileTypeText, false},
}

var fileTypeTestCasesText = []fileTypeTestCase{
	{uploader.FileTypeJPEG, false},
	{uploader.FileTypePNG, false},
	{uploader.FileTypePDF, false},
	{"application/zip", false},
	{uploader.FileTypeText, true},
}

var fileTypeTestCasesPDF = []fileTypeTestCase{
	{uploader.FileTypeJPEG, false},
	{uploader.FileTypePNG, false},
	{uploader.FileTypePDF, true},
	{"application/zip", false},
	{uploader.FileTypeText, false},
}

func (suite *UploaderSuite) verifyFileTypes(name string, allowedFileTypes uploader.AllowedFileTypes, cases []fileTypeTestCase) {
	for _, tc := range cases {
		suite.Run(fmt.Sprintf("%s/%s", name, tc.fileType), func() {
			suite.Equal(tc.allowed, allowedFileTypes.Contains(tc.fileType), "%s.Contains(%s) should be %v", name, tc.fileType, tc.allowed)
		})
	}

}

func (suite *UploaderSuite) TestAllowedFileTypes() {
	suite.verifyFileTypes("AllowedTypesServiceMember", uploader.AllowedTypesServiceMember, fileTypeTestCasesServiceMember)
	suite.verifyFileTypes("AllowedTypesPPMDocuments", uploader.AllowedTypesPPMDocuments, fileTypeTestCasesPPM)
	suite.verifyFileTypes("AllowedTypesPaymentRequest", uploader.AllowedTypesPaymentRequest, fileTypeTestCasesPaymentRequest)
	suite.verifyFileTypes("AllowedTypesAny", uploader.AllowedTypesAny, fileTypeTestCasesAny)
	suite.verifyFileTypes("AllowedTypesText", uploader.AllowedTypesText, fileTypeTestCasesText)
	suite.verifyFileTypes("AllowedTypesPDF", uploader.AllowedTypesPDF, fileTypeTestCasesPDF)
}
