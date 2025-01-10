package move

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveServiceSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *MoveServiceSuite) SetupSuite() {
	suite.PreloadData(func() {
		err := factory.DeleteAllotmentsFromDatabase(suite.DB())
		suite.FatalNoError(err)
		factory.SetupDefaultAllotments(suite.DB())
	})
}

func TestMoveServiceSuite(t *testing.T) {

	hs := &MoveServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
