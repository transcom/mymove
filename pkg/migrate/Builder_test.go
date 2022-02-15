package migrate

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MigrateSuite) TestBuilderCompile() {

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

	// Create the builder and point to the fixture path
	uri := "file://./fixtures/loop.sql"
	m := pop.Match{
		// Version MUST BE UNIQUE for this test to work
		Version:   "20190715140000",
		Name:      "loop",
		DBType:    "all",
		Direction: "up",
		Type:      "sql",
	}
	builder := &Builder{Match: &m, Path: uri}

	// Compile the migration
	wait := 10 * time.Millisecond
	migration, errCompile := builder.Compile(nil, wait, suite.Logger())
	suite.NoError(errCompile)
	suite.NotNil(migration)

	// Create a migrator and add the migration to it
	migrator := pop.NewMigrator(suite.DB())
	migrator.Migrations[migration.Direction] = append(migrator.Migrations[migration.Direction], *migration)

	// Migrate to use the Runner
	errUp := migrator.Up()
	suite.NoError(errUp)
}

func (suite *MigrateSuite) TestBuilderCompileInvalidPath() {

	// Create the builder and point to the fixture path
	uri := "invalid_path"
	m := pop.Match{
		Version:   "20190715144534",
		Name:      "invalid_path",
		DBType:    "all",
		Direction: "up",
		Type:      "sql",
	}
	builder := &Builder{Match: &m, Path: uri}

	wait := 10 * time.Millisecond
	migration, err := builder.Compile(nil, wait, suite.Logger())
	suite.NotNil(err)
	suite.Nil(migration)
}

func (suite *MigrateSuite) TestBuilderCompileBadType() {

	// Create the builder and point to the fixture path
	uri := "file://./fixtures/loop.sql"
	m := pop.Match{
		Version:   "20190715144534",
		Name:      "loop",
		DBType:    "all",
		Direction: "up",
		Type:      "bad_type",
	}
	builder := &Builder{Match: &m, Path: uri}

	wait := 10 * time.Millisecond
	migration, err := builder.Compile(nil, wait, suite.Logger())
	suite.NotNil(err)
	suite.Nil(migration)
}

func (suite *MigrateSuite) TestBuilderCompileInvalidDirection() {

	// Create the builder and point to the fixture path
	uri := "file://./fixtures/loop.sql"
	m := pop.Match{
		Version:   "20190715144534",
		Name:      "loop",
		DBType:    "all",
		Direction: "bad_direction",
		Type:      "sql",
	}
	builder := &Builder{Match: &m, Path: uri}

	wait := 10 * time.Millisecond
	migration, err := builder.Compile(nil, wait, suite.Logger())
	suite.NotNil(err)
	suite.Nil(migration)
}

func (suite *MigrateSuite) TestBuilderCompileUnsupportedDialect() {

	// Create the builder and point to the fixture path
	uri := "file://./fixtures/loop.sql"
	m := pop.Match{
		Version:   "20190715144534",
		Name:      "loop",
		DBType:    "bad_dialect",
		Direction: "up",
		Type:      "sql",
	}
	builder := &Builder{Match: &m, Path: uri}

	wait := 10 * time.Millisecond
	migration, err := builder.Compile(nil, wait, suite.Logger())
	suite.NotNil(err)
	suite.Nil(migration)
}

func (suite *MigrateSuite) TestBuilderCompileUpdateFromSetSQL() {

	// Create the builder and point to the fixture path
	uri := "file://./fixtures/update_from_set.sql"
	m := pop.Match{
		Version:   "20190715144534",
		Name:      "update_from_set",
		DBType:    "all",
		Direction: "up",
		Type:      "sql",
	}
	builder := &Builder{Match: &m, Path: uri}

	wait := 10 * time.Millisecond
	migration, err := builder.Compile(nil, wait, suite.Logger())
	suite.Nil(err)
	suite.NotNil(migration)

	// Migrate to use the Runner
	migrator := pop.NewMigrator(suite.DB())
	migrator.Migrations[migration.Direction] = append(migrator.Migrations[migration.Direction], *migration)
	err = migrator.Up()
	suite.Nil(err)
}
