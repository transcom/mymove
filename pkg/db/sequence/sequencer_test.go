package sequence

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type SequenceSuite struct {
	testingsuite.PopTestSuite
}

// AppContextForTest returns the AppContext for the test suite
func (suite *SequenceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), nil, nil)
}

func (suite *SequenceSuite) SetupTest() {
	err := suite.DB().RawQuery("CREATE SEQUENCE IF NOT EXISTS test_sequence;").Exec()
	suite.NoError(err, "Error creating test sequence")
	err = suite.DB().RawQuery("SELECT setval($1, 1);", testSequence).Exec()
	suite.NoError(err, "Error resetting sequence")
}

func TestSequenceSuite(t *testing.T) {

	hs := &SequenceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
