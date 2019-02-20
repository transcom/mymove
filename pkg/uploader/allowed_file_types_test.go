package uploader_test

import "github.com/transcom/mymove/pkg/uploader"

func (suite *UploaderSuite) TestAllowsAny() {
	suite.True(uploader.AllowedTypesAny.AllowsAny(), "AllowedTypesAny should allow any")
	suite.False(uploader.AllowedTypesPDF.AllowsAny(), "AllowedTypesPDF should not allow any")
	suite.False(uploader.AllowedTypesText.AllowsAny(), "AllowedTypesText should not allow any")
	suite.False(uploader.AllowedTypesServiceMember.AllowsAny(), "AllowedTypesServiceMember should not allow any")
}
