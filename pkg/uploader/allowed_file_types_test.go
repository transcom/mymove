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
	{"text/plain", true},
	{"image/gif", true},
}

var fileTypeTestCasesServiceMember = []fileTypeTestCase{
	{"image/jpeg", true},
	{"image/png", true},
	{"image/gif", false},
	{"application/pdf", true},
	{"application/zip", false},
	{"text/plain", false},
}

var fileTypeTestCasesText = []fileTypeTestCase{
	{"image/jpeg", false},
	{"image/png", false},
	{"application/pdf", false},
	{"application/zip", false},
	{"text/plain", true},
}

var fileTypeTestCasesPDF = []fileTypeTestCase{
	{"image/jpeg", false},
	{"image/png", false},
	{"application/pdf", true},
	{"application/zip", false},
	{"text/plain", false},
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
	suite.verifyFileTypes("AllowedTypesAny", uploader.AllowedTypesAny, fileTypeTestCasesAny)
	suite.verifyFileTypes("AllowedTypesText", uploader.AllowedTypesText, fileTypeTestCasesText)
	suite.verifyFileTypes("AllowedTypesPDF", uploader.AllowedTypesPDF, fileTypeTestCasesPDF)
}
