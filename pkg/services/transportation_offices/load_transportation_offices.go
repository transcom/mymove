package transportationoffices

import (
	"bufio"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofrs/uuid"

	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

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

type OfficeDataGeoLocation struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type OfficeDataLocation struct {
	Locality          string                `json:"locality"`
	CountryCode       string                `json:"country_code"`
	PostalCode        string                `json:"postal_code"`
	AddressLine1      string                `json:"address_line1"`
	AddressLine2      string                `json:"address_line2"`
	AdminstrativeArea string                `json:"administrative_area"`
	GeoLocation       OfficeDataGeoLocation `json:"geolocation"`
}

type OfficeShippingOffice struct {
	Type     string             `json:"type"`
	Title    string             `json:"title"`
	Location OfficeDataLocation `json:"location"`
}

type OfficeDataStruct struct {
	Type           string               `json:"type"`
	Title          string               `json:"title"`
	Location       OfficeDataLocation   `json:"location"`
	ShippingOffice OfficeShippingOffice `json:"shipping_office"`
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

func isUSAndConusFilter(o OfficeDataStruct) bool {
	return o.Location.CountryCode == "US" && o.Location.AdminstrativeArea != "AK" &&
		o.Location.AdminstrativeArea != "HI" && o.Type == "Transportation Office"
}

func isUSFilter(o OfficeDataStruct) bool {
	return o.Location.CountryCode == "US"
}

func isCONUSFilter(o OfficeDataStruct) bool {
	return o.Location.AdminstrativeArea != "AK" &&
		o.Location.AdminstrativeArea != "HI"
}

func (b *MigrationBuilder) filterOffice(os map[string]OfficeDataStruct, test func(OfficeDataStruct) bool) map[string]OfficeDataStruct {
	filtered := make(map[string]OfficeDataStruct)
	for key, o := range os {
		if test(o) {
			filtered[key] = o
		}
	}
	return filtered
}

// FindConusOffices find conus offices without shipping offices
func (b *MigrationBuilder) FindConusOffices(o OfficeDataStruct) models.TransportationOffices {
	dbOs, err := models.FetchTransportationOfficesByPostalCode(b.db, o.Location.PostalCode)
	if err != nil {
		fmt.Println(err)
	}

	// filter out the shipping offices
	var list models.TransportationOffices
	for _, office := range dbOs {
		//fmt.Println(office)
		if office.ShippingOfficeID == nil {
			continue
		}

		list = append(list, office)
	}

	return list
}

// FindShippingOffices finds the shipping offices in the transportation offices table
func (b *MigrationBuilder) FindShippingOffices(o OfficeShippingOffice) models.TransportationOffices {
	dbOs, err := models.FetchTransportationOfficesByPostalCode(b.db, o.Location.PostalCode)
	if err != nil {
		fmt.Println(err)
	}

	// filter out the shipping offices
	var list models.TransportationOffices
	for _, office := range dbOs {
		//fmt.Println(office)
		if office.ShippingOfficeID != nil {
			continue
		}

		list = append(list, office)
	}

	return list
}

func (b *MigrationBuilder) normalizeName(name string) string {
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
		n = b.convertAbbr(n)
		n = strings.Title(n)
		normalized = append(normalized, n)
	}

	return strings.Join(normalized, " ")
}

func (b *MigrationBuilder) convertAbbr(s string) string {
	for k, v := range abbrs {
		if k == s {
			return v
		}
	}
	return s
}

func (b *MigrationBuilder) getLocationsList() map[string]OfficeDataStruct {
	// message := map[string]interface{}{
	// 	"query": "55407",
	// }

	// bytesRepresentation, err := json.Marshal(message)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// resp, err := http.Post("https://move.mil/parser/locator-maps", "application/json", bytes.NewBuffer(bytesRepresentation))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	var result map[string]map[string]OfficeDataStruct
	officesJSON, _ := os.Open("./cmd/load_transportation_offices/data/transportation_offices.json")
	defer officesJSON.Close()
	json.NewDecoder(officesJSON).Decode(&result)

	// for key, element := range result["offices"] {
	// 	log.Println(key)
	// 	log.Println(element)
	// 	log.Println("")
	// }

	//log.Println(result["offices"])
	//log.Println(result["data"])
	return result["offices"]
}

func (b *MigrationBuilder) getPPSOGbloc() map[string]string {
	ppsoWithGbloc := make(map[string]string)
	csvFile, _ := os.Open("./cmd/load_transportation_offices/data/ppso_org_gbloc.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		//log.Println(line)
		// key: title of office, value: gbloc
		ppsoWithGbloc[strings.ToLower(strings.Trim(line[0], " "))] = line[2]
	}

	return ppsoWithGbloc
}

func (b *MigrationBuilder) Build() (string, error) {
	ppsoWithGloc := b.getPPSOGbloc()
	offices := b.getLocationsList()

	fmt.Printf("# total offices: %d\n", len(offices))

	usAndConusOffices := b.filterOffice(offices, isUSAndConusFilter)
	fmt.Printf("# us and conus only offices: %d\n", len(usAndConusOffices))

	fmt.Println("# of missing transportation offices: ")
	var missingOffices []OfficeDataStruct
	var foundOffices []OfficeDataStruct
	for _, o := range usAndConusOffices {
		//b.WriteXMLLine(o, w)

		// List the missing transportation office
		offices := b.FindConusOffices(o)
		if len(offices) == 0 {
			missingOffices = append(missingOffices, o)
			fmt.Println()
			fmt.Printf("%+v\n", o)
			fmt.Println()
		}

		// Filter matches of location
		for _, officeByZip := range offices {
			// we don't want shipping offices
			if officeByZip.ShippingOfficeID == nil {
				continue
			}

			//check by geo location
			// long is negative
			// lat is positive
			long, _ := strconv.ParseFloat(o.Location.GeoLocation.Lng, 64)
			lat, _ := strconv.ParseFloat(o.Location.GeoLocation.Lat, 64)

			fromLong := float64(officeByZip.Longitude - .2)
			toLong := float64(officeByZip.Longitude + .2)
			fromLat := float64(officeByZip.Latitude - .2)
			toLat := float64(officeByZip.Latitude + .2)

			if strings.EqualFold(officeByZip.Name, o.Title) {
				foundOffices = append(foundOffices, o)
			} else if strings.EqualFold(officeByZip.Address.StreetAddress1, o.Location.AddressLine1) {
				foundOffices = append(foundOffices, o)
			} else if fromLong <= long && toLong >= long && fromLat <= lat && toLat >= lat {
				foundOffices = append(foundOffices, o)
			} else {
				missingOffices = append(missingOffices, o)
				fmt.Println()
				fmt.Printf("%+v\n", o)
				fmt.Println()
				fmt.Println("Office by zip geo location difference: ")
				fmt.Println("Lat: ", float64(officeByZip.Latitude))
				fmt.Println("Long: ", float64(officeByZip.Longitude))
				fmt.Println("Office by json geo location difference: ")
				fmt.Println("Lat: ", o.Location.GeoLocation.Lat)
				fmt.Println("Long: ", o.Location.GeoLocation.Lng)
			}
		}
	}
	fmt.Println("# of found offices:")
	fmt.Println(len(foundOffices))
	fmt.Println("# of missing offices:")
	fmt.Println(len(missingOffices))

	//Match all missing transportation offices with a shipping office
	fmt.Println("Sql script to add new transportation offices: ")
	//var officesMissingShipping []OfficeDataStruct
	for _, missingOffice := range missingOffices {
		shippingOffices := b.FindShippingOffices(missingOffice.ShippingOffice)

		if len(shippingOffices) == 0 {
			shippingOfficeUUID, _ := uuid.NewV4()
			shippingAddressUUID, _ := uuid.NewV4()
			gbloc := ppsoWithGloc[missingOffice.ShippingOffice.Title]
			//Shipping office
			fmt.Println(fmt.Sprintf(`INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('%s', '%s', '%s', '%s', '%s', '%s', now(), now(), 'United States');`,
				shippingAddressUUID, missingOffice.Location.AddressLine1, missingOffice.Location.AddressLine2, missingOffice.Location.Locality, missingOffice.Location.AdminstrativeArea, missingOffice.Location.PostalCode))
			fmt.Println(fmt.Sprintf(`INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('%s', '%s', '%s', '%s', %s, %s, '%s', now(), now());`,
				shippingOfficeUUID, b.normalizeName(missingOffice.Title), gbloc, "NULL", missingOffice.Location.GeoLocation.Lat, missingOffice.Location.GeoLocation.Lng, shippingOfficeUUID))

			officeUUID, _ := uuid.NewV4()
			officeAddressUUID, _ := uuid.NewV4()
			//Transportation office
			fmt.Println(fmt.Sprintf(`INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('%s', '%s', '%s', '%s', '%s', '%s', now(), now(), 'United States');`,
				officeAddressUUID, missingOffice.Location.AddressLine1, missingOffice.Location.AddressLine2, missingOffice.Location.Locality, missingOffice.Location.AdminstrativeArea, missingOffice.Location.PostalCode))
			fmt.Println(fmt.Sprintf(`INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('%s', '%s', '%s', '%s', %s, %s, '%s', now(), now());`,
				officeUUID, b.normalizeName(missingOffice.Title), gbloc, officeAddressUUID, missingOffice.Location.GeoLocation.Lat, missingOffice.Location.GeoLocation.Lng, shippingOfficeUUID))

		}

		for _, shippingOffice := range shippingOffices {
			long, _ := strconv.ParseFloat(missingOffice.ShippingOffice.Location.GeoLocation.Lng, 64)
			fromLong := float64(shippingOffice.Longitude - .2)
			toLong := float64(shippingOffice.Longitude + .2)
			lat, _ := strconv.ParseFloat(missingOffice.ShippingOffice.Location.GeoLocation.Lat, 64)
			fromLat := float64(shippingOffice.Latitude - .2)
			toLat := float64(shippingOffice.Latitude + .2)

			// add missing transportation offices with known shipping office
			if strings.EqualFold(shippingOffice.Name, missingOffice.ShippingOffice.Title) ||
				strings.EqualFold(shippingOffice.Address.StreetAddress1, missingOffice.ShippingOffice.Location.AddressLine1) ||
				(fromLong <= long && toLong >= long && fromLat <= lat && toLat >= lat) {

				officeUUID, _ := uuid.NewV4()
				addressUUID, _ := uuid.NewV4()
				//Write sql script here
				fmt.Println(fmt.Sprintf(`INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('%s', '%s', '%s', '%s', '%s', '%s', now(), now(), 'United States');`,
					addressUUID, missingOffice.Location.AddressLine1, missingOffice.Location.AddressLine2, missingOffice.Location.Locality, missingOffice.Location.AdminstrativeArea, missingOffice.Location.PostalCode))
				fmt.Println(fmt.Sprintf(`INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('%s', '%s', '%s', '%s', %s, %s, '%s', now(), now());`,
					officeUUID, b.normalizeName(missingOffice.Title), shippingOffice.Gbloc, addressUUID, missingOffice.Location.GeoLocation.Lat, missingOffice.Location.GeoLocation.Lng, shippingOffice.ID))
			} else {
				// add missing transportation offices with new address
				//officesMissingShipping = append(officesMissingShipping, missingOffice)

				shippingOfficeUUID, _ := uuid.NewV4()
				shippingAddressUUID, _ := uuid.NewV4()
				gbloc := ppsoWithGloc[missingOffice.ShippingOffice.Title]
				//Shipping office
				fmt.Println(fmt.Sprintf(`INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('%s', '%s', '%s', '%s', '%s', '%s', now(), now(), 'United States');`,
					shippingAddressUUID, missingOffice.Location.AddressLine1, missingOffice.Location.AddressLine2, missingOffice.Location.Locality, missingOffice.Location.AdminstrativeArea, missingOffice.Location.PostalCode))
				fmt.Println(fmt.Sprintf(`INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('%s', '%s', '%s', '%s', %s, %s, '%s', now(), now());`,
					shippingOfficeUUID, b.normalizeName(missingOffice.Title), gbloc, "NULL", missingOffice.Location.GeoLocation.Lat, missingOffice.Location.GeoLocation.Lng, shippingOffice.ID))

				officeUUID, _ := uuid.NewV4()
				officeAddressUUID, _ := uuid.NewV4()
				//Transportation office
				fmt.Println(fmt.Sprintf(`INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('%s', '%s', '%s', '%s', '%s', '%s', now(), now(), 'United States');`,
					officeAddressUUID, missingOffice.Location.AddressLine1, missingOffice.Location.AddressLine2, missingOffice.Location.Locality, missingOffice.Location.AdminstrativeArea, missingOffice.Location.PostalCode))
				fmt.Println(fmt.Sprintf(`INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('%s', '%s', '%s', '%s', %s, %s, '%s', now(), now());`,
					officeUUID, b.normalizeName(missingOffice.Title), gbloc, officeAddressUUID, missingOffice.Location.GeoLocation.Lat, missingOffice.Location.GeoLocation.Lng, shippingOffice.ID))

			}
		}
	}

	// fmt.Println("# of missing offices missing shipping offices: ")
	// fmt.Println(len(officesMissingShipping))

	return "abc", nil
}
