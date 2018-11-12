package dbfmt

import (
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type DBFmtSuite struct {
	suite.Suite
	db *pop.Connection
}

func (suite *DBFmtSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *DBFmtSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func TestDBFmtSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	hs := &DBFmtSuite{db: db}
	suite.Run(t, hs)
}

func sameStrings(a []string, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}

func (suite *DBFmtSuite) TestTheDBFmt() {

	move := testdatagen.MakeDefaultMove(suite.db)
	suite.mustSave(&move)
	moveID := move.ID

	move = models.Move{}
	err := suite.db.Eager("Orders.Moves").Find(&move, moveID.String())
	suite.Nil(err)

	littlemove := move.Orders.Moves[0]
	move.Orders.Moves = append(move.Orders.Moves, littlemove)

	Println(move)

	// Uncomment this to work on formatting
	// suite.Fail("NONO")

}
