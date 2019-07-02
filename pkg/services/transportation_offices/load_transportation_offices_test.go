package transportationoffices

import (
	"fmt"
	"log"

	// "reflect"
	"testing"
	// "time"

	"github.com/gobuffalo/pop"
	// "github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	// "github.com/transcom/mymove/pkg/gen/internalmessages"
	// "github.com/transcom/mymove/pkg/models"
)

// run tests: go test -v  ./pkg/services/transportation_offices/
type DutyStationsLoaderSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *DutyStationsLoaderSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *DutyStationsLoaderSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func TestTransportationOfficesSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &DutyStationsLoaderSuite{
		db:     db,
		logger: logger,
	}

	suite.Run(t, hs)
}

func (suite *DutyStationsLoaderSuite) TestParsingFunctions() {
	officesPath := "./testdata/transportation_offices.xml"
	builder := NewMigrationBuilder(suite.db, suite.logger)

	officesRows, err := builder.parseOffices(officesPath)
	suite.NoError(err)
	suite.Len(officesRows, 5)
	officeRow := officesRows[0].LISTGCNSLINFO.GCNSLINFO
	suite.Equal(officeRow.CNSLNAME, "ALTUS AFB, OK")
	suite.Equal(officeRow.CNSLCITY, "ALTUS")
	suite.Equal(officeRow.CNSLSTATE, "OK")
	suite.Equal(officeRow.CNSLZIP, "73523")
	suite.Equal(officeRow.PPSONAME, "JPPSO SOUTH CENTRAL ")
	suite.Equal(officeRow.PPSOZIP, "78236")
}

func (suite *DutyStationsLoaderSuite) TestIsUsFilter() {
	officesPath := "./testdata/transportation_offices.xml"
	builder := NewMigrationBuilder(suite.db, suite.logger)

	officesRows, err := builder.parseOffices(officesPath)
	suite.NoError(err)
	suite.Len(officesRows, 5)

	usOffices := builder.isUS(officesRows)
	suite.Len(usOffices, 4)
}

func (suite *DutyStationsLoaderSuite) TestIsConusFilter() {
	officesPath := "./testdata/transportation_offices.xml"
	builder := NewMigrationBuilder(suite.db, suite.logger)

	officesRows, err := builder.parseOffices(officesPath)
	suite.NoError(err)
	suite.Len(officesRows, 5)

	usOffices := builder.isConus(officesRows)
	suite.Len(usOffices, 3)
}

func (suite *DutyStationsLoaderSuite) TestNormalizeName() {
	builder := NewMigrationBuilder(suite.db, suite.logger)
	var nameTests = []struct {
		name       string
		normalized string
	}{
		{"PPPO FORT MEADE, MD", "PPPO Fort Meade, MD"},
		{"PPPO GROTON/-NEW LONDON CT", "PPPO Groton/-New London CT"},
		{"COLUMBUS AFB, MS", "Columbus AFB, MS"},
		{"PPPO - USMA WEST POINT", "PPPO - USMA West Point"},
		{"PPPO Ft Sill", "PPPO Fort Sill"},
		{"PPPO - USNA Annapolis", "PPPO - USNA Annapolis"},
		{"KIRTLAND PPPO", "Kirtland PPPO"},
		{"PPPO - MCB QUANTICO", "PPPO - Marine Corp Base Quantico"},
		{"PPPO - USCG DIST WASHINGTON DC", "PPPO - USCG DIST Washington DC"},
		{"PPPO - JB ANDREWS-NAF", "PPPO - JB Andrews-NAF"},
		// PPPO, FLCJ, NAS, KEY WEST, FL

	}

	for _, n := range nameTests {
		fmt.Printf("\n\n")
		suite.Equal(builder.normalizeName(n.name), n.normalized)
	}
	// officesPath := "./testdata/transportation_offices.xml"
	// builder := NewMigrationBuilder(suite.db, suite.logger)

	// officesRows, err := builder.parseOffices(officesPath)
	// suite.NoError(err)
	// suite.Len(officesRows, 5)

	// usOffices := builder.isConus(officesRows)
	// suite.Len(usOffices, 3)
}
