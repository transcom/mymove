package sequence

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

const testSequence = "test_sequence"

func (suite *SequenceSuite) TestSetVal() {
	err := SetVal(suite.DB(), testSequence, 30)
	suite.NoError(err, "Error setting value of sequence")

	var nextVal int64
	err = suite.DB().RawQuery("SELECT nextval($1);", testSequence).First(&nextVal)
	suite.NoError(err, "Error getting current value of sequence")
	assert.Equal(suite.T(), nextVal, int64(31))
}

func (suite *SequenceSuite) TestNextVal() {
	actual, err := NextVal(suite.DB(), testSequence)
	if suite.NoError(err) {
		assert.Equal(suite.T(), actual, int64(2))
	}
}

type SequenceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *SequenceSuite) SetupTest() {
	suite.DB().TruncateAll()
	err := suite.DB().RawQuery("CREATE SEQUENCE IF NOT EXISTS test_sequence;").Exec()
	suite.NoError(err, "Error creating test sequence")
	err = suite.DB().RawQuery("SELECT setval($1, 1);", testSequence).Exec()
	suite.NoError(err, "Error resetting sequence")
}

func TestSequenceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &SequenceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
