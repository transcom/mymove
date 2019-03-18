package dutystationsloader

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
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

	builder := NewMigrationBuilder(suite.db, suite.logger)

	stationRows, err := builder.parseStations(stationsPath)
	suite.NoError(err)
	suite.Len(stationRows, 1)
	suite.Equal(stationRows[0].TransportationOfficeName, "Fort McFort PPPO")
	suite.Equal(stationRows[0].DutyStation.Name, "Fort McFort")
	suite.Equal(stationRows[0].DutyStation.Address.PostalCode, "79916")

	officesRows, err := builder.parseOffices(officesPath)
	suite.NoError(err)
	suite.Len(officesRows, 1)
	suite.Equal(officesRows[0].Name, "Fort McFort PPPO")
	suite.Equal(officesRows[0].Address.PostalCode, "79916")
	suite.Equal(officesRows[0].PhoneLines[0].Type, "Voice")
}

func (suite *DutyStationsLoaderSuite) TestInsertionString() {
	builder := NewMigrationBuilder(suite.db, suite.logger)

	something := "something"
	zeroID := uuid.Nil
	somethingID := uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f0"))
	suite.Equal("'something'", builder.insertionString(reflect.ValueOf(something)))
	suite.Equal("'something'", builder.insertionString(reflect.ValueOf(&something)))
	suite.Equal("'00000000-0000-0000-0000-000000000000'", builder.insertionString(reflect.ValueOf(zeroID)))
	suite.Equal("'00000000-0000-0000-0000-000000000000'", builder.insertionString(reflect.ValueOf(&zeroID)))
	suite.Equal("'cd40c92e-7c8a-4da4-ad58-4480df84b3f0'", builder.insertionString(reflect.ValueOf(somethingID)))
	suite.Equal("'cd40c92e-7c8a-4da4-ad58-4480df84b3f0'", builder.insertionString(reflect.ValueOf(&somethingID)))
	suite.Equal("'ARMY'", builder.insertionString(reflect.ValueOf(internalmessages.AffiliationARMY)))
	suite.Equal("false", builder.insertionString(reflect.ValueOf(false)))
	suite.Equal("5.6148", builder.insertionString(reflect.ValueOf(float32(5.61482))))
	suite.Equal("now()", builder.insertionString(reflect.ValueOf(time.Time{})))
}

func (suite *DutyStationsLoaderSuite) TestCreateInsertQuery() {
	builder := NewMigrationBuilder(suite.db, suite.logger)

	model := models.User{
		ID:            uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f0")),
		LoginGovUUID:  uuid.Must(uuid.FromString("cd40c92e-7c8a-4da4-ad58-4480df84b3f1")),
		LoginGovEmail: "email@example.com",
	}

	query := builder.createInsertQuery(model, &pop.Model{Value: models.User{}})

	suite.Equal(
		"INSERT into users (id, created_at, updated_at, login_gov_uuid, login_gov_email, disabled) VALUES ('cd40c92e-7c8a-4da4-ad58-4480df84b3f0', now(), now(), 'cd40c92e-7c8a-4da4-ad58-4480df84b3f1', 'email@example.com', false);\n",
		query)
}

func (suite *DutyStationsLoaderSuite) TestSeparateStations() {
	postalCode1 := "00001"
	address1 := models.Address{
		StreetAddress1: "something",
		City:           "something",
		State:          "CA",
		PostalCode:     postalCode1,
	}
	suite.mustSave(&address1)

	savedName := "Saved!"
	saved := models.DutyStation{
		AddressID:   address1.ID,
		Address:     address1,
		Name:        savedName,
		Affiliation: internalmessages.AffiliationARMY,
	}
	suite.mustSave(&saved)

	postalCode2 := "00002"
	address2 := models.Address{
		StreetAddress1: "something",
		City:           "something",
		State:          "CA",
		PostalCode:     postalCode2,
	}

	notSavedName := "Not saved sadface"
	notSaved := models.DutyStation{
		AddressID:   address2.ID,
		Address:     address2,
		Name:        notSavedName,
		Affiliation: internalmessages.AffiliationARMY,
	}

	builder := NewMigrationBuilder(suite.db, suite.logger)
	new, existing, err := builder.separateExistingStations([]DutyStationWrapper{
		DutyStationWrapper{
			TransportationOfficeName: "Some name",
			DutyStation:              saved,
		},
		DutyStationWrapper{
			TransportationOfficeName: "Some other name",
			DutyStation:              notSaved,
		},
	})

	suite.NoError(err)
	if suite.Len(new, 1) {
		suite.Equal(notSavedName, new[0].DutyStation.Name)
	}
	if suite.Len(existing, 1) {
		suite.Equal(savedName, existing[0].DutyStation.Name)
	}
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

	savedName := "Some Office"
	saved := models.TransportationOffice{
		AddressID: address.ID,
		Address:   address,
		Name:      savedName,
	}
	suite.mustSave(&saved)

	postalCode2 := "00002"
	address2 := models.Address{
		StreetAddress1: "something",
		City:           "something",
		State:          "CA",
		PostalCode:     postalCode2,
	}

	notSavedName := "This isn't saved"
	notSaved := models.TransportationOffice{
		AddressID: address2.ID,
		Address:   address2,
		Name:      notSavedName,
	}

	builder := NewMigrationBuilder(suite.db, suite.logger)
	new, existing, err := builder.separateExistingOffices([]models.TransportationOffice{saved, notSaved})

	suite.NoError(err)
	if suite.Len(new, 1) {
		suite.Equal(notSavedName, new[0].Name)
	}
	if suite.Len(existing, 1) {
		suite.Equal(savedName, existing[0].Name)
	}
}

func (suite *DutyStationsLoaderSuite) TestPairStationsOffices() {
	officeName1 := "oname1"
	stationName1 := "sname1"
	station1 := DutyStationWrapper{
		TransportationOfficeName: officeName1,
		DutyStation: models.DutyStation{
			Name: stationName1,
		},
	}
	office1 := TransportationOfficeWrapper{
		TransportationOfficeName: officeName1,
		TransportationOffice: models.TransportationOffice{
			Name: officeName1,
		},
	}

	officeName2 := "oname2"
	stationName2 := "sname2"
	station2 := DutyStationWrapper{
		TransportationOfficeName: officeName2,
		DutyStation: models.DutyStation{
			Name: stationName2,
		},
	}
	office2 := TransportationOfficeWrapper{
		TransportationOfficeName: officeName2,
		TransportationOffice: models.TransportationOffice{
			Name: officeName2,
		},
	}

	// This DutyStation has no pair, should still come out the other side though
	officeName3 := "oname3"
	stationName3 := "sname3"
	station3 := DutyStationWrapper{
		TransportationOfficeName: officeName3,
		DutyStation: models.DutyStation{
			Name: stationName3,
		},
	}

	builder := NewMigrationBuilder(suite.db, suite.logger)
	pairs := builder.pairOfficesToStations(
		[]DutyStationWrapper{station1, station2, station3},
		[]TransportationOfficeWrapper{office1, office2})

	suite.Len(pairs, 3)
	for _, p := range pairs {
		if p.DutyStation.Name == stationName1 {
			suite.Equal(officeName1, p.TransportationOffice.Name)
		}
		if p.DutyStation.Name == stationName2 {
			suite.Equal(officeName2, p.TransportationOffice.Name)
		}
		if p.DutyStation.Name == stationName3 {
			// Will be a blank model since it had no pair
			suite.Equal("", p.TransportationOffice.Name)
		}
	}
}
