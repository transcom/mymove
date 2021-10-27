package dbfmt

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBFmtSuite struct {
	testingsuite.PopTestSuite
}

func TestDBFmtSuite(t *testing.T) {
	hs := &DBFmtSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *DBFmtSuite) TestTheDBFmt() {

	move := testdatagen.MakeDefaultMove(suite.DB())
	suite.MustSave(&move)
	moveID := move.ID

	move = models.Move{}
	err := suite.DB().Eager("Orders.Moves").Find(&move, moveID.String())
	suite.NoError(err)

	littlemove := move.Orders.Moves[0]
	move.Orders.Moves = append(move.Orders.Moves, littlemove)

	Println(move)

	// Uncomment this to work on formatting
	// suite.Fail("NONO")

}
