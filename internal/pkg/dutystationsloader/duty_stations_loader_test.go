package dutystationsloader

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

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

func TestDutyStationsLoaderSuite(t *testing.T) {
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
	stationsPath := "./testdata/stations.xlsx"
	officesPath := "./testdata/offices.xlsx"

	stationRows, err := ParseStations(stationsPath)
	suite.NoError(err)
	suite.Len(stationRows, 1)
	suite.Equal(stationRows[0].Name, "Fort McFort")
	suite.Equal(stationRows[0].PostalCode, "79916")

	officesRows, err := ParseOffices(officesPath)
	suite.NoError(err)
	suite.Len(officesRows, 1)
	suite.Equal(officesRows[0].Name, "Fort McFort PPPO")
	suite.Equal(officesRows[0].PostalCode, "79916")
	suite.Equal(officesRows[0].Phone1.Type, "Voice")
}

func (suite *DutyStationsLoaderSuite) TestInsertionString() {
	something := "something"
	zeroID := uuid.Nil
	somethingID := uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f0"))
	suite.Equal("'something'", insertionString(reflect.ValueOf(something)))
	suite.Equal("'something'", insertionString(reflect.ValueOf(&something)))
	suite.Equal("'00000000-0000-0000-0000-000000000000'", insertionString(reflect.ValueOf(zeroID)))
	suite.Equal("'00000000-0000-0000-0000-000000000000'", insertionString(reflect.ValueOf(&zeroID)))
	suite.Equal("'cd40c92e-7c8a-4da4-ad58-4480df84b3f0'", insertionString(reflect.ValueOf(somethingID)))
	suite.Equal("'cd40c92e-7c8a-4da4-ad58-4480df84b3f0'", insertionString(reflect.ValueOf(&somethingID)))
	suite.Equal("'ARMY'", insertionString(reflect.ValueOf(internalmessages.AffiliationARMY)))
	suite.Equal("false", insertionString(reflect.ValueOf(false)))
	suite.Equal("5.6148", insertionString(reflect.ValueOf(float32(5.61482))))
	suite.Equal("now()", insertionString(reflect.ValueOf(time.Time{})))
}

func (suite *DutyStationsLoaderSuite) TestCreateInsertQuery() {
	model := models.User{
		ID:            uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f0")),
		LoginGovUUID:  uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f1")),
		LoginGovEmail: "email@example.com",
	}

	query := createInsertQuery(model, &pop.Model{Value: models.User{}})

	suite.Equal(
		"INSERT into users (id, created_at, updated_at, login_gov_uuid, login_gov_email) VALUES ('cd40c92e-7c8a-4da4-ad58-4480df84b3f0', now(), now(), 'cd40c92e-7c8a-4da4-ad58-4480df84b3f1', 'email@example.com');\n",
		query)
}

func (suite *DutyStationsLoaderSuite) TestCheckForDuplicates() {
	postalCode := "00001"
	address := models.Address{
		StreetAddress1: "something",
		City:           "something",
		State:          "CA",
		PostalCode:     postalCode,
	}
	suite.mustSave(&address)

	stationName := "Some Station"
	station := models.DutyStation{
		AddressID:   address.ID,
		Name:        stationName,
		Affiliation: internalmessages.AffiliationARMY,
	}
	suite.mustSave(&station)

	officeName := "Some Office"
	office := models.TransportationOffice{
		AddressID: address.ID,
		Name:      officeName,
	}
	suite.mustSave(&office)

	stationRows := []DutyStationRow{
		DutyStationRow{
			Name:       stationName,
			PostalCode: postalCode,
		},
		DutyStationRow{
			Name:       "Something new",
			PostalCode: "00002",
		},
	}

	officeRows := []TransportationOfficeRow{
		TransportationOfficeRow{
			Name:       officeName,
			PostalCode: postalCode,
		},
	}

	stationDupes, officeDupes, err := CheckDatabaseForDuplicates(suite.db, stationRows, officeRows)
	suite.NoError(err)
	suite.Len(stationDupes, 1)
	suite.Len(officeDupes, 1)
}
