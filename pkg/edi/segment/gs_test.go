package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateGS() {
	validGS := GS{
		FunctionalIdentifierCode: "SI",
		ApplicationSendersCode:   "MILMOVE",
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

	altValidGS := GS{
		FunctionalIdentifierCode: "AG",
		ApplicationSendersCode:   "8004171844",
		ApplicationReceiversCode: "MILMOVE",
		Date:                     "20190903",
		Time:                     "1617",
		GroupControlNumber:       1,
		ResponsibleAgencyCode:    "X",
		Version:                  "004010",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(altValidGS)
		suite.NoError(err)
	})

	suite.T().Run("validate failure 1", func(t *testing.T) {
		gs := GS{
			FunctionalIdentifierCode: "XX",        // oneof
			ApplicationSendersCode:   "XXXXX",     // oneof
			ApplicationReceiversCode: "123456789", // oneof
			Date:                     "20190945",  // datetime
			Time:                     "2517",      // datetime
			GroupControlNumber:       0,           // min
			ResponsibleAgencyCode:    "Y",         // eq
			Version:                  "123456",    // eq
		}

		err := suite.validator.Struct(gs)
		suite.ValidateError(err, "FunctionalIdentifierCode", "oneof")
		suite.ValidateError(err, "ApplicationSendersCode", "oneof")
		suite.ValidateError(err, "ApplicationReceiversCode", "oneof")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateError(err, "Time", "datetime=1504|datetime=150405")
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
