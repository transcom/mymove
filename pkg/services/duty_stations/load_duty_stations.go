package dutystations

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
)

const hereRequestTimeout = time.Duration(15) * time.Second

const (
	// InsertTemplate is the query insert template for duty stations
	InsertTemplate string = `
	{{range .}}
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at, country) VALUES ('{{.AddressID}}', 'N/A', '{{.Address.City}}', '{{.Address.State}}', '{{.Address.PostalCode}}', now(), now(), 'United States');
INSERT INTO duty_stations (id, name, affiliation, address_id, created_at, updated_at, transportation_office_id) VALUES ('{{.DutyStationID}}', '{{.Stations.Name}}', 'MARINES', '{{.AddressID}}', now(), now(), '{{.To.ID}}');
	{{end}}`
)

// DutyStationMigration represents a duty station migration
type DutyStationMigration struct {
	Address       models.Address
	To            models.TransportationOffice
	Stations      StationData
	AddressID     uuid.UUID
	DutyStationID uuid.UUID
}

// StationData represents Duty Station data
type StationData struct {
	Unit string
	Name string
	Zip  string
}

// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func (b MigrationBuilder) ParseStations(filename string) ([]StationData, error) {
	var stations []StationData

	csvFile, err := os.Open(filepath.Clean(filename))
	if err != nil {
		fmt.Println(err)
		return stations, err
	}
	r := csv.NewReader(csvFile)

	// Skip the first header row
	dataRows, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return stations, err
	}
	for _, row := range dataRows[1:] {
		parsed := StationData{
			Unit: row[0],
			Name: row[1],
			Zip:  row[2],
		}
		if parsed.Name == "" {
			continue
		}
		stations = append(stations, parsed)
	}

	return stations, nil
}

// MigrationBuilder has methods that assist in building a DutyStation INSERT migration
type MigrationBuilder struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewMigrationBuilder returns a new instance of a MigrationBuilder
func NewMigrationBuilder(db *pop.Connection, logger *zap.Logger) MigrationBuilder {
	return MigrationBuilder{
		db,
		logger,
	}
}

func (b *MigrationBuilder) filterMarines(dss models.DutyStations) models.DutyStations {
	var filtered []models.DutyStation
	for _, ds := range dss {
		if ds.Affiliation == internalmessages.AffiliationMARINES {
			filtered = append(filtered, ds)
		}
	}
	return filtered
}

func (b *MigrationBuilder) findDutyStations(s StationData) models.DutyStations {
	zip := s.Zip
	stations, err := models.FetchDutyStationsByPostalCode(b.db, zip)
	if err != nil {
		fmt.Println(err)
	}
	filteredStations := b.filterMarines(stations)
	return filteredStations
}

func (b *MigrationBuilder) addressLatLong(address models.Address) (route.LatLong, error) {
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	p := route.NewHEREPlannerHP(b.logger, hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)

	plannerType := reflect.TypeOf(p)
	for i := 0; i < plannerType.NumMethod(); i++ {
		method := plannerType.Method(i)
		fmt.Println(method.Name)
	}

	return p.GetAddressLatLong(&address)
}

func getCityState(unit string) (string, string) {
	lst := strings.Split(unit, " ")
	if len(lst[len(lst)-1]) != 2 {
		fmt.Println("Misformatted unit: ", unit)
	}
	return strings.Join(lst[:len(lst)-1], " "), lst[len(lst)-1]
}

func (b *MigrationBuilder) nearestTransportationOffice(address models.Address) (models.TransportationOffice, error) {
	latLong, err := b.addressLatLong(address)
	if err != nil {
		return models.TransportationOffice{}, err
	}
	to, err := models.FetchNearestTransportationOffice(b.db, latLong.Longitude, latLong.Latitude)
	if err != nil {
		return to, err
	}
	return to, nil
}

// Build builds a migration for loading duty stations
func (b *MigrationBuilder) Build(dutyStationsFilePath string) ([]DutyStationMigration, error) {
	stations, err := b.ParseStations(dutyStationsFilePath)
	if err != nil {
		return nil, err
	}

	var DutyStationMigrations []DutyStationMigration
	for _, s := range stations {
		city, state := getCityState(s.Unit)
		address := models.Address{
			City:       city,
			State:      state,
			PostalCode: s.Zip,
		}
		if state == "HI" || state == "AK" {
			continue
		}

		dbDutyStations := b.findDutyStations(s)
		if len(dbDutyStations) == 0 {
			//fmt.Println("*** missing... add?? ***")
			to, err := b.nearestTransportationOffice(address)
			if err != nil {
				fmt.Println("Error encountered finding nearest transportation office: ", err)
				continue
			}
			DutyStationMigrations = append(DutyStationMigrations, DutyStationMigration{address, to, s, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())})
		}
	}
	return DutyStationMigrations, nil
}
