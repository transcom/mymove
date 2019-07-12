package dutystations

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/route"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const hereRequestTimeout = time.Duration(15) * time.Second

var uppercaseWords = map[string]bool{
	// seeing double w/ a comma == a hack to deal w/ commas in the office name
	"AFB":     true,
	"AFB,":    true,
	"DIST":    true,
	"DIST,":   true,
	"FLCJ":    true,
	"FLCJ,":   true,
	"JB":      true,
	"JRB":     true,
	"JRB,":    true,
	"LCR":     true,
	"LCR,":    true,
	"MCAS":    true,
	"MCAS,":   true,
	"NAVSUP":  true,
	"NAVSUP,": true,
	"NAF":     true,
	"NAF,":    true,
	"NAS":     true,
	"NAS,":    true,
	"PPPO":    true,
	"PPPO,":   true,
	"USCG":    true,
	"USCG,":   true,
	"USMA":    true,
	"USMA,":   true,
	"USNA":    true,
	"USNA,":   true,
	"HQTRS":   true,
}

var states = map[string]bool{
	"AL": true,
	"AK": true,
	"AZ": true,
	"AR": true,
	"CA": true,
	"CO": true,
	"CT": true,
	"DC": true,
	"DE": true,
	"FL": true,
	"GA": true,
	"HI": true,
	"ID": true,
	"IL": true,
	"IN": true,
	"IA": true,
	"KS": true,
	"KY": true,
	"LA": true,
	"ME": true,
	"MD": true,
	"MA": true,
	"MI": true,
	"MN": true,
	"MS": true,
	"MO": true,
	"MT": true,
	"NE": true,
	"NV": true,
	"NH": true,
	"NJ": true,
	"NM": true,
	"NY": true,
	"NC": true,
	"ND": true,
	"OH": true,
	"OK": true,
	"OR": true,
	"PA": true,
	"RI": true,
	"SC": true,
	"SD": true,
	"TN": true,
	"TX": true,
	"UT": true,
	"VT": true,
	"VA": true,
	"WA": true,
	"WV": true,
	"WI": true,
	"WY": true,
}

var abbrs = map[string]string{
	"ft":          "fort",
	"mcb":         "marine corp base",
	"andrews-naf": "Andrews-NAF",
}

type StationData struct {
	Unit string
	Name string
	Zip  string
}

// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func (b MigrationBuilder) ParseStations(filename string) ([]StationData, error) {
	var stations []StationData

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return stations, err
	}
	r := csv.NewReader(csvFile)

	// Skip the first header row
	dataRows, err := r.ReadAll()
	// dataRows := xlFile.Sheets[1].Rows[1:245]
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

func FilterTransportationOffices(os models.TransportationOffices, test func(models.TransportationOffice) bool) models.TransportationOffices {
	var filtered models.TransportationOffices
	for _, o := range os {
		if test(o) {
			filtered = append(filtered, o)
		}
	}
	return filtered
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
	//fmt.Println("unfiltered duty stations: ", len(stations))
	filteredStations := b.filterMarines(stations)
	//fmt.Println("filtered duty stations: ", len(filteredStations))
	return filteredStations
}

// func (b *MigrationBuilder) FindTransportationOffice(s StationData) models.TransportationOffices {
// 	zip := s.Zip

// 	dbOs, err := models.FetchTransportationOfficesByPostalCode(b.db, zip)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	if len(dbOs) == 0 {
// 		partialZip := zip[:len(zip)-1] + "%"
// 		//fmt.Fprintf(w, "*** partialZip: %s \n", partialZip)
// 		dbOs, err = models.FetchTransportationOfficesByPostalCode(b.db, partialZip)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}

// 	return dbOs
// }

func (b *MigrationBuilder) WriteLine(s StationData, row *[]string) {
	name := normalizeName(s.Name)
	//fmt.Printf("\nname: %s  | zip: %s \n", name, s.Zip)
	//fmt.Fprintf(w, "\nname: %s  | zip: %s \n", name, s.Zip)
	newRow := append(*row, name, s.Zip)
	*row = newRow
}

func (b *MigrationBuilder) WriteDbRecs(ts models.DutyStations) {
	for _, t := range ts {
		fmt.Println("\tdb: ", t.Name, " | ", t.Affiliation)
	}
}

func normalizeName(name string) string {
	var normalized []string
	nameSplit := strings.Fields(name)
	for _, n := range nameSplit {
		if _, exists := uppercaseWords[n]; exists {
			normalized = append(normalized, n)
			continue
		}

		if _, exists := states[n]; exists {
			normalized = append(normalized, n)
			continue
		}

		n = strings.ToLower(n)
		n = convertAbbr(n)
		n = strings.Title(n)
		normalized = append(normalized, n)
	}
	return strings.Join(normalized, " ")
}

func convertAbbr(s string) string {
	for k, v := range abbrs {
		if k == s {
			return v
		}
	}
	return s
}

func (b *MigrationBuilder) addressLatLong(address models.Address) (route.LatLong, error) {
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	hereClient := &http.Client{Timeout: hereRequestTimeout}
	p := route.NewHEREPlannerMine(b.logger, hereClient, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)

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

func createInsertAddress(address models.Address, id uuid.UUID) string {
	// nolint
	return fmt.Sprintf(`INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at, country) VALUES ('%s', 'N/A', '%s', '%s', '%s', now(), now(), 'US');`, id, address.City, address.State, address.PostalCode)
}

func createInsertDutyStations(addressID uuid.UUID, officeID uuid.UUID, stationName string) string {
	dutyStationID := uuid.Must(uuid.NewV4())
	// nolint
	return fmt.Sprintf(`INSERT INTO duty_stations (id, name, affiliation, address_id, created_at, updated_at, transportation_office_id) VALUES ('%s', '%s', 'MARINES', '%s', now(), now(), '%s');`, dutyStationID, stationName, addressID, officeID)
}

func (b *MigrationBuilder) generateInsertionBlock(address models.Address, to models.TransportationOffice, station StationData) string {
	var query strings.Builder
	addressID := uuid.Must(uuid.NewV4())

	query.WriteString(createInsertAddress(address, addressID))
	query.WriteString("\n")
	query.WriteString(createInsertDutyStations(addressID, to.ID, station.Name))
	query.WriteString("\n")

	return query.String()
}

func (b *MigrationBuilder) Build(dutyStationsFilePath string) (string, error) {
	stations, err := b.ParseStations(dutyStationsFilePath)
	if err != nil {
		return "", err
	}
	//fmt.Printf("# total stations: %d\n", len(stations))

	var migration strings.Builder
	for _, s := range stations {
		//fmt.Println("\n", s.Name, " | ", s.Zip)
		city, state := getCityState(s.Unit)
		address := models.Address{
			City:       city,
			State:      state,
			PostalCode: s.Zip,
		}
		//fmt.Println(city, " | ", state)
		if state == "HI" || state == "AK" {
			fmt.Println("\t*** skipping non-conus")
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
			migration.WriteString(b.generateInsertionBlock(address, to, s))
		}
	}
	return migration.String(), nil
}
