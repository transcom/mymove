package dbfmt

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type DBFmtSuite struct {
	testingsuite.PopTestSuite
}

func (suite *DBFmtSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDBFmtSuite(t *testing.T) {
	hs := &DBFmtSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}

func sameStrings(a []string, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}

func (suite *DBFmtSuite) TestTheDBFmt() {

	move := testdatagen.MakeDefaultMove(suite.DB())
	suite.MustSave(&move)
	moveID := move.ID

	move = models.Move{}
	err := suite.DB().Eager("Orders.Moves").Find(&move, moveID.String())
	suite.Nil(err)

	littlemove := move.Orders.Moves[0]
	move.Orders.Moves = append(move.Orders.Moves, littlemove)

	Println(move)

	// Uncomment this to work on formatting
	// suite.Fail("NONO")

}
