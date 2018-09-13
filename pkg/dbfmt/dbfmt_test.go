package dbfmt_test

import (
	"log"
	"testing"

	"fmt"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"reflect"
	"sort"

	. "github.com/transcom/mymove/pkg/dbfmt"
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

func (suite *DBFmtSuite) verifyValidationErrors(model models.ValidateableModel, exp map[string][]string) {
	t := suite.T()
	t.Helper()

	verrs, err := model.Validate(suite.db)
	if err != nil {
		t.Fatal(err)
	}

	if verrs.Count() != len(exp) {
		t.Errorf("expected %d errors, got %d", len(exp), verrs.Count())
	}

	var expKeys []string
	for key, errors := range exp {
		e := verrs.Get(key)
		expKeys = append(expKeys, key)
		if !sameStrings(e, errors) {
			t.Errorf("expected errors on %s to be %v, got %v", key, errors, e)
		}
	}

	for _, key := range verrs.Keys() {
		if !sliceContains(key, expKeys) {
			errors := verrs.Get(key)
			t.Errorf("unexpected validation errors on %s: %v", key, errors)
		}
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

func sliceContains(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
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

	goodString := PrettyString(move)
	fmt.Println(goodString)

	suite.Fail("NONO")

}
