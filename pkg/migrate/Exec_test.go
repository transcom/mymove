package migrate

import (
	"os"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MigrateSuite) TestExecWithDeleteUsersSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/deleteUsers.sql")
	suite.Nil(err)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait)
		suite.Nil(err)
		return nil
	})
}

func (suite *MigrateSuite) TestExecWithCopyFromStdinSQL() {

	// Create common TSP
	testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			ID: uuid.Must(uuid.FromString("231a7b21-346c-4e94-b6bc-672413733f77")),
		},
	})

	// Create TDLs
	testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			ID: uuid.Must(uuid.FromString("27f1fbeb-090c-4a91-955c-67899de4d6d6")),
		},
	})

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/copyFromStdin.sql")
	suite.Nil(err)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait)
		suite.Nil(err)
		return nil
	})
}

func (suite *MigrateSuite) TestExecWithLoopSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/loop.sql")
	suite.Nil(err)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait)
		suite.Nil(err)
		return nil
	})
}

func (suite *MigrateSuite) TestExecWithUpdateFromSetSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/update_from_set.sql")
	suite.Nil(err)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait)
		suite.Nil(err)
		return nil
	})
}
