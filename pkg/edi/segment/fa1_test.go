package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateFA1() {
	validFA1 := FA1{
		AgencyQualifierCode: "DF",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validFA1)
		suite.NoError(err)
	})

	suite.T().Run("validate failure", func(t *testing.T) {
		fa1 := FA1{
			AgencyQualifierCode: "XX", // oneof
		}

		err := suite.validator.Struct(fa1)
		suite.ValidateError(err, "AgencyQualifierCode", "oneof")
		suite.ValidateErrorLen(err, 1)
	})
}
