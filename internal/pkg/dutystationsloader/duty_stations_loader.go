package dutystationsloader

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

var splitRepl = regexp.MustCompile(" |,|-")

// A reverse-mapping from string value back to internalmessages.Affiliation
var affiliationMap = map[string]internalmessages.Affiliation{
	"ARMY":        internalmessages.AffiliationARMY,
	"NAVY":        internalmessages.AffiliationNAVY,
	"MARINES":     internalmessages.AffiliationMARINES,
	"AIR_FORCE":   internalmessages.AffiliationAIRFORCE,
	"COAST_GUARD": internalmessages.AffiliationCOASTGUARD,
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

// DutyStationWrapper wraps DutyStation data but retains TransportationOfficeName for pairing with an office
type DutyStationWrapper struct {
	TransportationOfficeName string
	models.DutyStation
}

// TransportationOfficeWrapper wraps TransportationOffice data but retains the original TransportationOfficeName for pairing with a station
type TransportationOfficeWrapper struct {
	TransportationOfficeName string
	models.TransportationOffice
}

// StationOfficePair pairs a DutyStationRow to a TransportationOfficeRow
type StationOfficePair struct {
	models.DutyStation
	models.TransportationOffice
}

// column pairs a column name with its INSERT query-stringified value
type column struct {
	name  string
	value string
}

// Gotta have a stringPointer function. Returns nil if empty string
func stringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// floatFormatter reformats a float string as an int
func floatFormatter(f string) string {
	float, err := strconv.ParseFloat(f, 64)
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf("%.0f", float)
}

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
}

func similarityPattern(s string) string {
	// "Some name" -> "%(some|name)% for SIMILAR TO lookups in postgres"
	return "%(" + splitRepl.ReplaceAllString(strings.ToLower(s), "|") + ")%"
}

// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func (b *MigrationBuilder) parseStations(path string) ([]DutyStationWrapper, error) {
	var stations []DutyStationWrapper

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return stations, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for _, row := range dataRows {
		parsed := DutyStationWrapper{
			TransportationOfficeName: getCell(row.Cells, 8),
			DutyStation: models.DutyStation{
				Name:        getCell(row.Cells, 0),
				Affiliation: affiliationMap[getCell(row.Cells, 1)],
				Address: models.Address{
					StreetAddress1: getCell(row.Cells, 2),
					StreetAddress2: stringPointer(getCell(row.Cells, 3)),
					StreetAddress3: stringPointer(getCell(row.Cells, 4)),
					City:           getCell(row.Cells, 5),
					State:          getCell(row.Cells, 6),
					PostalCode:     floatFormatter(getCell(row.Cells, 7)),
					Country:        stringPointer("United States"),
				},
			},
		}
		stations = append(stations, parsed)
	}

	return stations, nil
}

// Creates an OfficePhoneLine model given a row and column numbers
func (b *MigrationBuilder) parsePhoneData(row *xlsx.Row, namecol, labelcol, dsncol, typecol int) models.OfficePhoneLine {
	return models.OfficePhoneLine{
		Number:      getCell(row.Cells, namecol),
		Label:       stringPointer(getCell(row.Cells, labelcol)),
		IsDsnNumber: getCell(row.Cells, dsncol) != "FALSE",
		Type:        getCell(row.Cells, typecol),
	}
}

// Creates an OfficeEmail model given a row and column numbers
func (b *MigrationBuilder) parseEmailData(row *xlsx.Row, emailcol, labelcol int) models.OfficeEmail {
	return models.OfficeEmail{
		Email: getCell(row.Cells, emailcol),
		Label: stringPointer(getCell(row.Cells, labelcol)),
	}
}

// ParseOffices parses a spreadsheet of transportation offices into TransportationOfficeRow structs
func (b *MigrationBuilder) parseOffices(path string) ([]models.TransportationOffice, error) {
	var offices []models.TransportationOffice

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return offices, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for _, row := range dataRows {
		var phones []models.OfficePhoneLine
		if phone := b.parsePhoneData(row, 10, 11, 12, 13); phone.Number != "" {
			phones = append(phones, phone)
		}
		if phone := b.parsePhoneData(row, 14, 15, 16, 17); phone.Number != "" {
			phones = append(phones, phone)
		}
		if phone := b.parsePhoneData(row, 18, 19, 20, 21); phone.Number != "" {
			phones = append(phones, phone)
		}

		var emails []models.OfficeEmail
		if email := b.parseEmailData(row, 22, 23); email.Email != "" {
			emails = append(emails, email)
		}
		if email := b.parseEmailData(row, 24, 25); email.Email != "" {
			emails = append(emails, email)
		}
		if email := b.parseEmailData(row, 26, 27); email.Email != "" {
			emails = append(emails, email)
		}

		parsed := models.TransportationOffice{
			Name: getCell(row.Cells, 0),
			Address: models.Address{
				StreetAddress1: getCell(row.Cells, 1),
				StreetAddress2: stringPointer(getCell(row.Cells, 2)),
				StreetAddress3: stringPointer(getCell(row.Cells, 3)),
				City:           getCell(row.Cells, 4),
				State:          getCell(row.Cells, 5),
				PostalCode:     floatFormatter(getCell(row.Cells, 6)),
				Country:        stringPointer("United States"),
			},
			Latitude:   float32(0),
			Longitude:  float32(0),
			Hours:      stringPointer(getCell(row.Cells, 8)),
			Services:   stringPointer(getCell(row.Cells, 9)),
			PhoneLines: phones,
			Emails:     emails,
		}
		offices = append(offices, parsed)
	}

	return offices, nil
}

// Wraps a string value in single-quotes to play nice with psql syntax
func quoter(s string) string {
	return fmt.Sprintf("'%s'", s)
}

// Transforms a reflect.Value into a stringified value good for insertion
func (b *MigrationBuilder) insertionString(v reflect.Value) string {
	switch v := v.Interface().(type) {
	case string:
		return quoter(v)
	case *string:
		if v == nil {
			return ""
		}
		return quoter(*v)
	case uuid.UUID:
		return quoter(v.String())
	case *uuid.UUID:
		if v == nil {
			return "NULL"
		}
		return quoter(v.String())
	case internalmessages.Affiliation:
		return quoter(string(v))
	case bool:
		return strconv.FormatBool(v)
	case float32:
		// The only float field used is lat + long
		return fmt.Sprintf("%.4f", v)
	case time.Time:
		// Making a strong assumption that the only time fields are created_at and updated_at
		return "now()"
	default:
		return ""
	}
}

// Takes a model and a table name and creates an INSERT query
func (b *MigrationBuilder) createInsertQuery(m interface{}, model pop.TableNameAble) string {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)

	var cols []string
	var vals []string
	for i := 0; i < t.NumField(); i++ {
		// Grab the column name from the db tag and pair it with the new value
		field := t.Field(i)
		tag := field.Tag.Get("db")

		fieldVal := v.FieldByName(field.Name)

		val := b.insertionString(fieldVal)
		if tag != "" && val != "" {
			cols = append(cols, tag)
			vals = append(vals, val)
		}
	}

	// Values here should only come from trusted spreadsheets, and migration should be inspected after before being run
	// #nosec G201
	return fmt.Sprintf("INSERT into %v (%v) VALUES (%v);\n", model.TableName(), strings.Join(cols, ", "), strings.Join(vals, ", "))
}

// Extract our station/office information into mymove models and generate INSERT statements
func (b *MigrationBuilder) generateInsertionBlock(pair StationOfficePair) string {
	var query strings.Builder
	station := pair.DutyStation
	office := pair.TransportationOffice

	// New office, need to add IDs and INSERT
	if office.Name != "" && office.ID == uuid.Nil {
		// We need to generate IDs before we create the INSERT statements
		office.Address.ID = uuid.Must(uuid.NewV4())
		query.WriteString(b.createInsertQuery(office.Address, &pop.Model{Value: models.Address{}}))
		office.ID = uuid.Must(uuid.NewV4())
		office.AddressID = office.Address.ID
		query.WriteString(b.createInsertQuery(office, &pop.Model{Value: models.TransportationOffice{}}))
		for _, p := range office.PhoneLines {
			p.ID = uuid.Must(uuid.NewV4())
			p.TransportationOfficeID = office.ID
			query.WriteString(b.createInsertQuery(p, &pop.Model{Value: models.OfficePhoneLine{}}))
		}
		for _, e := range office.Emails {
			e.ID = uuid.Must(uuid.NewV4())
			e.TransportationOfficeID = office.ID
			query.WriteString(b.createInsertQuery(e, &pop.Model{Value: models.OfficeEmail{}}))
		}
	}

	if station.Name != "" {
		// If we have a valid office, set up the relationship
		if office.ID != uuid.Nil {
			station.TransportationOfficeID = &office.ID
		}

		// We'll only ever have new DutyStations here, no need to check
		station.Address.ID = uuid.Must(uuid.NewV4())
		query.WriteString(b.createInsertQuery(station.Address, &pop.Model{Value: models.Address{}}))
		station.ID = uuid.Must(uuid.NewV4())
		station.AddressID = station.Address.ID
		query.WriteString(b.createInsertQuery(station, &pop.Model{Value: models.DutyStation{}}))
	}

	return query.String()
}

// Given a list of DutyStation data, separate into two lists: new data, and data that already exists in our db
func (b *MigrationBuilder) separateExistingStations(stations []DutyStationWrapper) ([]DutyStationWrapper, []DutyStationWrapper, error) {
	var new []DutyStationWrapper
	var existing []DutyStationWrapper
	for _, newStation := range stations {
		var existingStation models.DutyStation
		query := b.db.Q().Eager().
			LeftJoin("addresses", "addresses.id=duty_stations.address_id").
			Where(`
				lower(name)=lower($1)
			OR
				lower(name) SIMILAR TO $2
				AND
				postal_code=$3
			`, newStation.DutyStation.Name, similarityPattern(newStation.DutyStation.Name), newStation.DutyStation.Address.PostalCode)
			// Where("lower(name) SIMILAR TO $1", similarityPattern(newStation.DutyStation.Name)).
			// Where("postal_code = $2", newStation.DutyStation.Address.PostalCode)
		err := query.First(&existingStation)
		if err == nil {
			b.logger.Debug("Found existing duty station in db", zap.String("New station name", newStation.Name), zap.String("Existing station", existingStation.Name))
			existing = append(existing, DutyStationWrapper{
				TransportationOfficeName: newStation.TransportationOfficeName,
				DutyStation:              existingStation,
			})
		} else if errors.Cause(err) == sql.ErrNoRows {
			new = append(new, newStation)
		} else {
			return new, existing, err
		}
	}

	return new, existing, nil
}

// Given a list of TransportationOffice data, separate into two lists: new data, and data that already exists in our db
func (b *MigrationBuilder) separateExistingOffices(offices []models.TransportationOffice) ([]TransportationOfficeWrapper, []TransportationOfficeWrapper, error) {
	var new []TransportationOfficeWrapper
	var existing []TransportationOfficeWrapper
	for _, newOffice := range offices {
		var existingOffice models.TransportationOffice
		query := b.db.Q().Eager().
			LeftJoin("addresses", "addresses.id=transportation_offices.address_id").
			Where(`
				lower(name)=lower($1)
			OR
				lower(name) SIMILAR TO $2
				AND
				postal_code=$3
			`, newOffice.Name, similarityPattern(newOffice.Name), newOffice.Address.PostalCode)
		err := query.First(&existingOffice)
		if err == nil {
			b.logger.Debug("Found existing transportation office in db", zap.String("New office name", newOffice.Name), zap.String("Existing office", existingOffice.Name))
			existing = append(existing, TransportationOfficeWrapper{
				TransportationOfficeName: newOffice.Name,
				TransportationOffice:     existingOffice,
			})
		} else if errors.Cause(err) == sql.ErrNoRows {
			new = append(new, TransportationOfficeWrapper{
				TransportationOfficeName: newOffice.Name,
				TransportationOffice:     newOffice,
			})
		} else {
			return new, existing, err
		}
	}

	return new, existing, nil
}

// Given a new DutyStation, try searching the db for a matching transportation office
func (b *MigrationBuilder) findMatchingOffice(station DutyStationWrapper) (models.TransportationOffice, error) {
	var office models.TransportationOffice
	query := b.db.Q().Eager().
		LeftJoin("addresses", "addresses.id=transportation_offices.address_id").
		Where(`
				lower(name)=lower($1)
			OR
				lower(name) SIMILAR TO $2
				AND
				postal_code=$3
			`, station.DutyStation.Name, similarityPattern(station.DutyStation.Name), station.DutyStation.Address.PostalCode)
	err := query.First(&office)

	return office, err
}

// PairStationsAndOffices creates pairs of duty stations and transportation offices using the dutyStation.TransportationOffice field
func (b *MigrationBuilder) pairOfficesToStations(stations []DutyStationWrapper, offices []TransportationOfficeWrapper) []StationOfficePair {
	// Create a map of office name to TransportationOfficeRow so we can easily look it up using station.TransportationOffice
	officesByName := make(map[string]TransportationOfficeWrapper)
	// We'll delete from this as we pair offices so we know what's left unpaired
	unpairedOfficeNames := make(map[string]bool)
	for _, t := range offices {
		officesByName[t.TransportationOfficeName] = t
		unpairedOfficeNames[t.TransportationOfficeName] = true
	}

	var pairs []StationOfficePair
	for _, s := range stations {
		if office, ok := officesByName[s.TransportationOfficeName]; ok {
			// Try to find a matching office using local data
			pairs = append(pairs, StationOfficePair{
				DutyStation:          s.DutyStation,
				TransportationOffice: office.TransportationOffice,
			})
			delete(unpairedOfficeNames, s.TransportationOfficeName)
		} else if dbOffice, err := b.findMatchingOffice(s); err == nil {
			// Else use a matched office from the database
			b.logger.Debug("Found matching transportation office in db", zap.String("Station Name", s.DutyStation.Name), zap.String("Existing Office", dbOffice.Name))
			pairs = append(pairs, StationOfficePair{
				DutyStation:          s.DutyStation,
				TransportationOffice: dbOffice,
			})
		} else {
			b.logger.Debug("Can't find a matching office for duty station", zap.String("DutyStation office name", s.TransportationOfficeName))
			// If there's no office we still want to insert the station
			pairs = append(pairs, StationOfficePair{
				DutyStation: s.DutyStation,
			})
		}
	}

	// Append left over unpaired offices
	for n := range unpairedOfficeNames {
		office := officesByName[n]
		// Existing offices will have a populated ID field, we don't need to append them since there's no station to worry about
		if office.ID == uuid.Nil {
			b.logger.Debug("New Office has no matching duty station", zap.String("TransportationOffice Name", n))
			pairs = append(pairs, StationOfficePair{
				TransportationOffice: officesByName[n].TransportationOffice,
			})
		}
	}

	return pairs
}

// Build orchestrates building the contents of a DutyStation INSERT migration
func (b *MigrationBuilder) Build(stationsFilePath, officesFilePath string) (string, error) {
	// Parse raw data from spreadsheets
	rawStations, err := b.parseStations(stationsFilePath)
	if err != nil {
		return "", err
	}

	rawOffices, err := b.parseOffices(officesFilePath)
	if err != nil {
		return "", err
	}

	// Separate data into lists of new data and existing data
	newStations, existingStations, err := b.separateExistingStations(rawStations)
	if err != nil {
		return "", err
	}
	newOffices, existingOffices, err := b.separateExistingOffices(rawOffices)
	if err != nil {
		return "", err
	}

	b.logger.Info("Separated stations",
		zap.Int("Number of new stations", len(newStations)),
		zap.Int("Number of existing stations", len(existingStations)))
	b.logger.Info("Separated offices",
		zap.Int("Number of new offices", len(newOffices)),
		zap.Int("Number of existing offices", len(existingOffices)))

	// We only care about new duty stations
	stations := newStations
	// We might need to link new duty stations to existing offices, so combine them here
	// Existing offices will be recognizable by having ID populated
	offices := append(newOffices, existingOffices...)

	// Pairs stations/offices by name, looks for matching offices in db
	pairs := b.pairOfficesToStations(stations, offices)

	var migration strings.Builder
	for _, pair := range pairs {
		migration.WriteString(b.generateInsertionBlock(pair))
		migration.WriteString("\n")
	}

	return migration.String(), nil
}
