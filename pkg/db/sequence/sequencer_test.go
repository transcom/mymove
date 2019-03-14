package sequence

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SequenceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *SequenceSuite) SetupTest() {
	suite.DB().TruncateAll()
	err := suite.DB().RawQuery("CREATE SEQUENCE IF NOT EXISTS test_sequence;").Exec()
	suite.NoError(err, "Error creating test sequence")
	err = suite.DB().RawQuery("SELECT setval($1, 1);", testSequence).Exec()
	suite.NoError(err, "Error resetting sequence")
}

func TestSequenceSuite(t *testing.T) {

	hs := &SequenceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}
