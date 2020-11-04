package edisegment

import (
	"fmt"
	"testing"
)

func (suite *SegmentSuite) TestValidateGS() {
	validGS := GS{
		FunctionalIdentifierCode: "SI",
		ApplicationSendersCode:   fmt.Sprintf("%-9s", "MYMOVE"),
		ApplicationReceiversCode: "8004171844",
		Date:                     "20190903",
		Time:                     "1617",
		GroupControlNumber:       1,
		ResponsibleAgencyCode:    "X",
		Version:                  "004010",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validGS)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		gs := GS{
			FunctionalIdentifierCode: "XX",                         // eq
			ApplicationSendersCode:   fmt.Sprintf("%-9s", "XXXXX"), // eq
			ApplicationReceiversCode: "123456789",                  // eq
			Date:                     "20190945",                   // datetime
			Time:                     "2517",                       // datetime
			GroupControlNumber:       0,                            // min
			ResponsibleAgencyCode:    "Y",                          // eq
			Version:                  "123456",                     // eq
		}

		err := suite.validator.Struct(gs)
		suite.ValidateError(err, "FunctionalIdentifierCode", "eq")
		suite.ValidateError(err, "ApplicationSendersCode", "eq")
		suite.ValidateError(err, "ApplicationReceiversCode", "eq")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateError(err, "Time", "datetime")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateError(err, "ResponsibleAgencyCode", "eq")
		suite.ValidateError(err, "Version", "eq")
		suite.ValidateErrorLen(err, 8)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		gs := validGS
		gs.GroupControlNumber = 1000000000 // max

		err := suite.validator.Struct(gs)
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
