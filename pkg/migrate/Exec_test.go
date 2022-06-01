package migrate

import (
	"os"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MigrateSuite) TestExecWithDeleteUsersSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/deleteUsers.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
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
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
}

func (suite *MigrateSuite) TestExecWithCopyFromStdinMultipleSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/copyFromStdinMultiple.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
}

func (suite *MigrateSuite) TestExecWithCopyFromStdinTypes() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/copyFromStdinTypes.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		innerErr := Exec(f, tx, wait, suite.Logger())
		suite.NoError(innerErr)
		return innerErr
	})
	suite.NoError(errTransaction)

	// Verify we didn't lose anything when inserting into database
	// (some fields were previously converted from string to float/int as part of the exec)

	// String
	var contract models.ReContract
	err = suite.DB().Find(&contract, "8ef44ca4-c589-4c39-93d8-3410a762ec6c")
	suite.NoError(err)
	suite.Equal("Pricing", contract.Code)

	// Float
	var contractYear models.ReContractYear
	err = suite.DB().Find(&contractYear, "9ca0c8d2-3b14-4a49-9709-1f710383642f")
	suite.NoError(err)
	suite.InDelta(1.01970, contractYear.Escalation, 0.0001)

	// Integer (but should store as string due to leading zeros)
	var domesticServiceArea models.ReDomesticServiceArea
	err = suite.DB().Find(&domesticServiceArea, "ea949687-431e-4c00-bbe7-646158471d4e")
	suite.NoError(err)
	suite.Equal("004", domesticServiceArea.ServiceArea)

	// null/nil
	var zip3 models.ReZip3
	err = suite.DB().Find(&zip3, "7eee5a4f-c457-49eb-bd6a-216fb00ab43c")
	suite.NoError(err)
	suite.Nil(zip3.RateAreaID)
}

func (suite *MigrateSuite) TestExecWithLoopSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/loop.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
}

func (suite *MigrateSuite) TestExecWithUpdateFromSetSQL() {

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/update_from_set.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
}

func (suite *MigrateSuite) TestExecWithInsertConflictSQL() {

	// Create Transportation Office
	testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			ID: uuid.Must(uuid.FromString("c219d9e5-2659-427d-be33-bf439251b7f3")),
		},
	})

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/insert_conflict.sql")
	suite.NoError(err)

	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
		wait := 10 * time.Millisecond
		err := Exec(f, tx, wait, suite.Logger())
		suite.NoError(err)
		return err
	})
	suite.NoError(errTransaction)
}
