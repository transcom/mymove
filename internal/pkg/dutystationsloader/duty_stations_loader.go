package dutystationsloader

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/tealeg/xlsx"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// A reverse-mapping from string value back to internalmessages.Affiliation
var affiliationMap = map[string]internalmessages.Affiliation{
	"ARMY":        internalmessages.AffiliationARMY,
	"NAVY":        internalmessages.AffiliationNAVY,
	"MARINES":     internalmessages.AffiliationMARINES,
	"AIR_FORCE":   internalmessages.AffiliationAIRFORCE,
	"COAST_GUARD": internalmessages.AffiliationCOASTGUARD,
}

// DutyStationRow contains a single row of a duty station spreadsheet
type DutyStationRow struct {
	RowNum               int
	Name                 string
	Affiliation          string
	StreetAddress1       string
	StreetAddress2       string
	StreetAddress3       string
	City                 string
	State                string
	PostalCode           string
	TransportationOffice string
}

// PhoneInfo contains phone-related fields from a transportation office spreadsheet
type PhoneInfo struct {
	Number string
	Label  string
	IsDSN  bool
	Type   string
}

// EmailInfo contains email-related fields from a transportation office spreadsheet
type EmailInfo struct {
	Email string
	Label string
}

// TransportationOfficeRow contains a single row of a transportation office spreadsheet
type TransportationOfficeRow struct {
	RowNum         int
	Name           string
	StreetAddress1 string
	StreetAddress2 string
	StreetAddress3 string
	City           string
	State          string
	PostalCode     string
	LatLong        string
	Hours          string
	Services       string
	Phone1         PhoneInfo
	Phone2         PhoneInfo
	Phone3         PhoneInfo
	Email1         EmailInfo
	Email2         EmailInfo
	Email3         EmailInfo
}

// StationOfficePair pairs a DutyStationRow to a TransportationOfficeRow
type StationOfficePair struct {
	DutyStationRow
	TransportationOfficeRow
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

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
}

// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func ParseStations(path string) ([]DutyStationRow, error) {
	var rows []DutyStationRow

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return rows, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for i, row := range dataRows {
		parsed := DutyStationRow{
			RowNum:               i,
			Name:                 getCell(row.Cells, 0),
			Affiliation:          getCell(row.Cells, 1),
			StreetAddress1:       getCell(row.Cells, 2),
			StreetAddress2:       getCell(row.Cells, 3),
			StreetAddress3:       getCell(row.Cells, 4),
			City:                 getCell(row.Cells, 5),
			State:                getCell(row.Cells, 6),
			PostalCode:           getCell(row.Cells, 7),
			TransportationOffice: getCell(row.Cells, 8),
		}
		rows = append(rows, parsed)
	}

	return rows, nil
}

// ParseOffices parses a spreadsheet of transportation offices into TransportationOfficeRow structs
func ParseOffices(path string) ([]TransportationOfficeRow, error) {
	var rows []TransportationOfficeRow

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return rows, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for i, row := range dataRows {
		parsed := TransportationOfficeRow{
			RowNum:         i,
			Name:           getCell(row.Cells, 0),
			StreetAddress1: getCell(row.Cells, 1),
			StreetAddress2: getCell(row.Cells, 2),
			StreetAddress3: getCell(row.Cells, 3),
			City:           getCell(row.Cells, 4),
			State:          getCell(row.Cells, 5),
			PostalCode:     getCell(row.Cells, 6),
			LatLong:        getCell(row.Cells, 7),
			Hours:          getCell(row.Cells, 8),
			Services:       getCell(row.Cells, 9),
			Phone1: PhoneInfo{
				Number: getCell(row.Cells, 10),
				Label:  getCell(row.Cells, 11),
				IsDSN:  getCell(row.Cells, 12) != "FALSE",
				Type:   getCell(row.Cells, 13),
			},
			Phone2: PhoneInfo{
				Number: getCell(row.Cells, 14),
				Label:  getCell(row.Cells, 15),
				IsDSN:  getCell(row.Cells, 16) != "FALSE",
				Type:   getCell(row.Cells, 17),
			},
			Phone3: PhoneInfo{
				Number: getCell(row.Cells, 18),
				Label:  getCell(row.Cells, 19),
				IsDSN:  getCell(row.Cells, 20) != "FALSE",
				Type:   getCell(row.Cells, 21),
			},
			Email1: EmailInfo{
				Email: getCell(row.Cells, 22),
				Label: getCell(row.Cells, 23),
			},
			Email2: EmailInfo{
				Email: getCell(row.Cells, 24),
				Label: getCell(row.Cells, 25),
			},
			Email3: EmailInfo{
				Email: getCell(row.Cells, 26),
				Label: getCell(row.Cells, 27),
			},
		}
		rows = append(rows, parsed)
	}

	return rows, nil
}

// Wraps a string value in single-quotes to play nice with psql syntax
func quoter(s string) string {
	return fmt.Sprintf("'%v'", s)
}

// Transforms a reflect.Value into a stringified value good for insertion
func insertionString(v reflect.Value) string {
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
func createInsertQuery(m interface{}, model pop.TableNameAble) string {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)

	var cols []string
	var vals []string
	for i := 0; i < t.NumField(); i++ {
		// Grab the column name from the db tag and pair it with the new value
		field := t.Field(i)
		tag := field.Tag.Get("db")

		fieldVal := v.FieldByName(field.Name)

		val := insertionString(fieldVal)
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
func generateInsertionBlock(pair StationOfficePair) string {
	stationRow := pair.DutyStationRow
	officeRow := pair.TransportationOfficeRow

	//
	// Transportation office models
	//
	officeAddress := models.Address{
		ID:             uuid.Must(uuid.NewV4()),
		StreetAddress1: officeRow.StreetAddress1,
		StreetAddress2: stringPointer(officeRow.StreetAddress2),
		StreetAddress3: stringPointer(officeRow.StreetAddress3),
		City:           officeRow.City,
		State:          officeRow.State,
		PostalCode:     officeRow.PostalCode,
		Country:        stringPointer("United States"),
	}
	office := models.TransportationOffice{
		ID:        uuid.Must(uuid.NewV4()),
		Name:      officeRow.Name,
		AddressID: officeAddress.ID,
		Hours:     stringPointer(officeRow.Hours),
		Services:  stringPointer(officeRow.Services),
	}
	phones := []PhoneInfo{
		pair.TransportationOfficeRow.Phone1,
		pair.TransportationOfficeRow.Phone2,
		pair.TransportationOfficeRow.Phone3,
	}
	var phoneModels []models.OfficePhoneLine
	for _, p := range phones {
		if p.Number == "" {
			continue
		}

		model := models.OfficePhoneLine{
			ID: uuid.Must(uuid.NewV4()),
			TransportationOfficeID: office.ID,
			Number:                 p.Number,
			Label:                  stringPointer(p.Label),
			IsDsnNumber:            p.IsDSN,
			Type:                   p.Type,
		}
		phoneModels = append(phoneModels, model)
	}

	emails := []EmailInfo{
		pair.TransportationOfficeRow.Email1,
		pair.TransportationOfficeRow.Email2,
		pair.TransportationOfficeRow.Email3,
	}
	var emailModels []models.OfficeEmail
	for _, e := range emails {
		if e.Email == "" {
			continue
		}

		model := models.OfficeEmail{
			ID: uuid.Must(uuid.NewV4()),
			TransportationOfficeID: office.ID,
			Email: e.Email,
			Label: stringPointer(e.Label),
		}
		emailModels = append(emailModels, model)
	}

	//
	// Duty station models
	//
	stationAddress := models.Address{
		ID:             uuid.Must(uuid.NewV4()),
		StreetAddress1: stationRow.StreetAddress1,
		StreetAddress2: stringPointer(stationRow.StreetAddress2),
		StreetAddress3: stringPointer(stationRow.StreetAddress3),
		City:           stationRow.City,
		State:          stationRow.State,
		PostalCode:     stationRow.PostalCode,
		Country:        stringPointer("United States"),
	}

	station := models.DutyStation{
		ID: uuid.Must(uuid.NewV4()),
		TransportationOfficeID: &office.ID,
		Name:        stationRow.Name,
		Affiliation: affiliationMap[stationRow.Affiliation],
		AddressID:   stationAddress.ID,
	}

	// Finally, build our block of INSERT queries. Order is important for foreign key relationships.
	var query strings.Builder
	query.WriteString(createInsertQuery(officeAddress, &pop.Model{Value: models.Address{}}))
	query.WriteString(createInsertQuery(office, &pop.Model{Value: models.TransportationOffice{}}))
	for _, p := range phoneModels {
		query.WriteString(createInsertQuery(p, &pop.Model{Value: models.OfficePhoneLine{}}))
	}
	for _, e := range emailModels {
		query.WriteString(createInsertQuery(e, &pop.Model{Value: models.OfficeEmail{}}))
	}
	query.WriteString(createInsertQuery(stationAddress, &pop.Model{Value: models.Address{}}))
	query.WriteString(createInsertQuery(station, &pop.Model{Value: models.DutyStation{}}))
	query.WriteString("\n")

	return query.String()
}

// CheckDatabaseForDuplicates searches local database for similarly named/located duty stations and transportation offices
// Searches for matching words in `name` field combined with matching `postal_code`
func CheckDatabaseForDuplicates(db *pop.Connection, stationRows []DutyStationRow, officeRows []TransportationOfficeRow) ([]DutyStationRow, []TransportationOfficeRow, error) {
	var err error
	var stationDupes []DutyStationRow
	var officeDupes []TransportationOfficeRow
	for _, row := range stationRows {
		var stations []models.DutyStation
		query := db.Q().Eager().
			LeftJoin("addresses", "addresses.id=duty_stations.address_id").
			Where("name ILIKE $1", "%"+strings.Replace(row.Name, " ", "%", -1)+"%").
			Where("postal_code = $2", row.PostalCode)
		err = query.All(&stations)
		if len(stations) > 0 {
			stationDupes = append(stationDupes, row)
		}
	}
	if err != nil {
		return stationDupes, officeDupes, err
	}

	for _, row := range officeRows {
		var offices []models.TransportationOffice
		query := db.Q().Eager().
			LeftJoin("addresses", "addresses.id=transportation_offices.address_id").
			Where("name ILIKE $1", "%"+strings.Replace(row.Name, " ", "%", -1)+"%").
			Where("postal_code = $2", row.PostalCode)
		err = query.All(&offices)
		if len(offices) > 0 {
			officeDupes = append(officeDupes, row)
		}
	}
	if err != nil {
		return stationDupes, officeDupes, err
	}

	return stationDupes, officeDupes, nil
}

// PairStationsAndOffices creates pairs of duty stations and transportation offices using the dutyStation.TransportationOffice field
func PairStationsAndOffices(stationRows []DutyStationRow, officeRows []TransportationOfficeRow) ([]StationOfficePair, []DutyStationRow, []TransportationOfficeRow) {
	// Create a map of office name to TransportationOfficeRow so we can easily look it up using station.TransportationOffice
	officesByName := make(map[string]TransportationOfficeRow)
	// We'll delete from this as we pair offices so we know what's left unpaired
	officesUnpaired := make(map[string]bool)
	for _, t := range officeRows {
		officesByName[t.Name] = t
		officesUnpaired[t.Name] = true
	}

	var pairs []StationOfficePair
	var unpairedStations []DutyStationRow
	for _, s := range stationRows {
		// Either we find a matching TransportationOffice name or it's unpaired
		if office, ok := officesByName[s.TransportationOffice]; ok {
			pairs = append(pairs, StationOfficePair{
				DutyStationRow:          s,
				TransportationOfficeRow: office,
			})
			delete(officesUnpaired, s.TransportationOffice)
		} else {
			unpairedStations = append(unpairedStations, s)
		}
	}

	var unpairedOffices []TransportationOfficeRow
	for n := range officesUnpaired {
		unpairedOffices = append(unpairedOffices, officesByName[n])
	}

	return pairs, unpairedStations, unpairedOffices
}

// GenerateMigrationString generates the contents of a migration to INSERT the supplied station/office pairs
func GenerateMigrationString(pairs []StationOfficePair) string {
	// For each station/office pair, create a block of INSERT queries and append to migration file
	var migration strings.Builder
	for _, pair := range pairs {
		migration.WriteString(generateInsertionBlock(pair))
	}

	return migration.String()
}
