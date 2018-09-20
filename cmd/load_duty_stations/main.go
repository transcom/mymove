package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/namsral/flag"
	"github.com/tealeg/xlsx"

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

// Parses a spreadsheet of duty stations into DutyStationRow structs
func parseStations(path string) ([]DutyStationRow, error) {
	var rows []DutyStationRow

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return rows, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for _, row := range dataRows {
		parsed := DutyStationRow{
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

// Parses a spreadsheet of transportation offices into TransportationOfficeRow structs
func parseOffices(path string) ([]TransportationOfficeRow, error) {
	var rows []TransportationOfficeRow

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return rows, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for _, row := range dataRows {
		parsed := TransportationOfficeRow{
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
		return "0"
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

// Searches local database for similarly named/located duty stations and transportation offices
func checkDatabaseForDuplicates(db *pop.Connection, stationRows []DutyStationRow, officeRows []TransportationOfficeRow) {
	var stationDupes []string
	for _, row := range stationRows {
		var stations []models.DutyStation
		query := db.Q().Eager().
			Where("name ILIKE $1", "%"+strings.Replace(row.Name, " ", "%", -1)+"%").
			Where("postal_code = $1", row.PostalCode)
		query.All(&stations)
		for _, s := range stations {
			stationDupes = append(stationDupes, fmt.Sprintf("Existing: %v, New: %v", s.Name, row.Name))
		}
	}
	fmt.Printf("Found %v duty station duplicates!\n", len(stationDupes))
	for _, d := range stationDupes {
		fmt.Println(d)
	}

	var officeDupes []string
	for _, row := range officeRows {
		var offices []models.TransportationOffice
		query := db.Q().Eager().
			Where("name ILIKE $1", "%"+strings.Replace(row.Name, " ", "%", -1)+"%").
			Where("postal_code = $1", row.PostalCode)
		// fmt.Println(query.ToSQL(&pop.Model{Value: models.TransportationOffice{}}))
		query.All(&offices)
		for _, s := range offices {
			officeDupes = append(officeDupes, fmt.Sprintf("Existing: %v, New: %v", s.Name, row.Name))
		}
	}
	fmt.Printf("Found %v transportation office duplicates!\n", len(officeDupes))
	for _, d := range officeDupes {
		fmt.Println(d)
	}
}

// Creates pairs of duty stations and transportation offices using the dutyStation.TransportationOffice field
func pairStationsAndOffices(stationRows []DutyStationRow, officeRows []TransportationOfficeRow) []StationOfficePair {
	officesByName := make(map[string]TransportationOfficeRow)
	officesUnpaired := make(map[string]bool)
	for _, t := range officeRows {
		officesByName[t.Name] = t
		officesUnpaired[t.Name] = true
	}

	var pairs []StationOfficePair
	for _, s := range stationRows {
		if office, ok := officesByName[s.TransportationOffice]; ok {
			pairs = append(pairs, StationOfficePair{
				DutyStationRow:          s,
				TransportationOfficeRow: office,
			})
			delete(officesUnpaired, s.TransportationOffice)
		} else {
			fmt.Printf("Couldn't find a matching transportation office for %v: wanted \"%v\"\n", s.Name, s.TransportationOffice)
		}
	}

	fmt.Printf("After pairing, there are %v transportation offices left unpaired!\n", len(officesUnpaired))
	for n := range officesUnpaired {
		fmt.Printf("Unpaired office: %v\n", n)
	}

	return pairs
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	validate := flag.Bool("validate", false, "Only run file validations")
	output := flag.String("output", "", "Where to output the migration file")
	stationsPath := flag.String("stations", "", "Input file for duty stations")
	officesPath := flag.String("offices", "", "Input file for transportation offices")
	flag.Parse()

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Panic(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// Parse duty stations
	stationRows, err := parseStations(*stationsPath)
	if err != nil {
		log.Panic(err)
	}

	// Parse transportation offices
	officeRows, err := parseOffices(*officesPath)
	if err != nil {
		log.Panic(err)
	}

	// Searches the database for existing stations/offices
	checkDatabaseForDuplicates(db, stationRows, officeRows)

	// Attempts to pair transportation offices with duty stations
	pairs := pairStationsAndOffices(stationRows, officeRows)
	fmt.Printf("Found %v pairs!\n", len(pairs))

	// If we just want to validate files we can exit
	if validate != nil && *validate {
		os.Exit(0)
	}

	// For each station/office pair, create a block of INSERT queries and append to migration file
	var migration strings.Builder
	migration.WriteString("-- Migration generated using cmd/load_duty_stations\n")
	migration.WriteString(fmt.Sprintf("-- Duty stations file: %v\n", *stationsPath))
	migration.WriteString(fmt.Sprintf("-- Transportation offices file: %v\n", *officesPath))
	migration.WriteString("\n")
	for _, pair := range pairs {
		migration.WriteString(generateInsertionBlock(pair))
	}

	f, err := os.OpenFile(*output, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	_, err = f.WriteString(migration.String())
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Complete! Migration written to %v\n", *output)
}
