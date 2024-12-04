package sequence

import (
	"github.com/transcom/mymove/pkg/testingsuite"
)

type SequenceSuite struct {
	*testingsuite.PopTestSuite
}

// func (suite *SequenceSuite) SetupTest() {
// 	err := suite.DB().RawQuery("SELECT setval($1, 1);", testSequence).Exec()
// 	suite.NoError(err, "Error resetting sequence")
// }

// func TestSequenceSuite(t *testing.T) {

// 	hs := &SequenceSuite{
// 		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
// 	}
// 	suite.Run(t, hs)
// 	hs.PopTestSuite.TearDown()
// }
