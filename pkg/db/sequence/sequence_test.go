package sequence

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

const testSequence = "test_sequence"

func (suite *SequenceSuite) TestSetVal() {
	err := SetVal(suite.db, testSequence, 30)
	suite.NoError(err, "Error setting value of sequence")

	var nextVal int64
	err = suite.db.RawQuery("SELECT nextval($1);", testSequence).First(&nextVal)
	suite.NoError(err, "Error getting current value of sequence")
	assert.Equal(suite.T(), nextVal, int64(31))
}

func (suite *SequenceSuite) TestNextVal() {
	actual, err := NextVal(suite.db, testSequence)
	if suite.NoError(err) {
		assert.Equal(suite.T(), actual, int64(2))
	}
}

type SequenceSuite struct {
	testingsuite.BaseTestSuite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *SequenceSuite) SetupTest() {
	suite.db.TruncateAll()
	err := suite.db.RawQuery("CREATE SEQUENCE IF NOT EXISTS test_sequence;").Exec()
	suite.NoError(err, "Error creating test sequence")
	err = suite.db.RawQuery("SELECT setval($1, 1);", testSequence).Exec()
	suite.NoError(err, "Error resetting sequence")
}

func TestSequenceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &SequenceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
