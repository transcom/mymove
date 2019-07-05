package dutystations

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"

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

type StationData struct {
	Unit string
	Name string
	Zip  string
}

// A safe way to get a cell from a slice of cells, returning empty string if not found
func getCell(cells []*xlsx.Cell, i int) string {
	if len(cells) > i {
		return cells[i].String()
	}

	return ""
}

// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func (b MigrationBuilder) ParseStations(filename string) ([]StationData, error) {
	var stations []StationData

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return stations, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[1].Rows[1:]
	// dataRows := xlFile.Sheets[1].Rows[1:245]
	for _, row := range dataRows {
		parsed := StationData{
			Unit: getCell(row.Cells, 0),
			Name: getCell(row.Cells, 1),
			Zip:  getCell(row.Cells, 2),
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

//func isUSFilter(o Office) bool {
//	return o.LISTGCNSLINFO.GCNSLINFO.CNSLCOUNTRY == "US"
//}
//
//func isCONUSFilter(o Office) bool {
//	return o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "AK" &&
//		o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "HI"
//}
//
//func isAFBFilter(o Office) bool {
//	return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "AFB")
//}
//
//func isArmyFilter(o Office) bool {
//	return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "FORT") ||
//		strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "FT") ||
//		strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "USMA")
//}
//
//func isNavyFilter(o Office) bool {
//	return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "NAV") ||
//		strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "NAVSUP") ||
//		strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "USNA")
//}
//
//func isOtherFilter(o Office) bool {
//	return !(isAFBFilter(o) || isArmyFilter(o) || isNavyFilter(o))
//}
//
//func (b *MigrationBuilder) isUS(offices []Office) []Office {
//	filter := func(o Office) bool {
//		return o.LISTGCNSLINFO.GCNSLINFO.CNSLCOUNTRY == "US"
//	}
//	return b.filterOffice(offices, filter)
//}
//
//func (b *MigrationBuilder) isConus(offices []Office) []Office {
//	filter := func(o Office) bool {
//		return o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "AK" &&
//			o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "HI"
//	}
//	return b.filterOffice(offices, filter)
//}
//
//func (b *MigrationBuilder) isAFB(offices []Office) []Office {
//	filter := func(o Office) bool {
//		return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "AFB")
//	}
//	return b.filterOffice(offices, filter)
//}
//
//func (b *MigrationBuilder) isArmy(offices []Office) []Office {
//	filter := func(o Office) bool {
//		return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "FORT") ||
//			strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "FT") ||
//			strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "USMA")
//	}
//	return b.filterOffice(offices, filter)
//}
//
//func (b *MigrationBuilder) isNav(offices []Office) []Office {
//	filter := func(o Office) bool {
//		return strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "NAV") ||
//			strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "NAVSUP") ||
//			strings.Contains(strings.ToUpper(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME), "USNA")
//
//	}
//	return b.filterOffice(offices, filter)
//}
//
//func (b *MigrationBuilder) filterOffice(os []Office, test func(Office) bool) []Office {
//	var filtered []Office
//	for _, o := range os {
//		if test(o) {
//			filtered = append(filtered, o)
//		}
//	}
//	return filtered
//}

func FilterTransportationOffices(os models.TransportationOffices, test func(models.TransportationOffice) bool) models.TransportationOffices {
	var filtered models.TransportationOffices
	for _, o := range os {
		if test(o) {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

func (b *MigrationBuilder) FindTransportationOffice(s StationData) models.TransportationOffices {
	zip := s.Zip

	dbOs, err := models.FetchTransportationOfficesByPostalCode(b.db, zip)
	if err != nil {
		fmt.Println(err)
	}

	if len(dbOs) == 0 {
		partialZip := zip[:len(zip)-1] + "%"
		//fmt.Fprintf(w, "*** partialZip: %s \n", partialZip)
		dbOs, err = models.FetchTransportationOfficesByPostalCode(b.db, partialZip)
		if err != nil {
			fmt.Println(err)
		}
	}

	return dbOs
}

func (b *MigrationBuilder) WriteLine(s StationData, row *[]string) {
	name := b.normalizeName(s.Name)

	//fmt.Printf("\nname: %s  | zip: %s \n", name, s.Zip)
	//fmt.Fprintf(w, "\nname: %s  | zip: %s \n", name, s.Zip)
	newRow := append(*row, name, s.Zip)
	*row = newRow
}

func (b *MigrationBuilder) WriteDbRecs(ts models.TransportationOffices, row *[]string) int {
	var newRow []string
	if len(ts) == 0 {
		// b.logger.Debug("*** NONE FOUND... BLAH")
		//fmt.Printf("*** %s NOT FOUND\n", officeType)
		//fmt.Fprintf(w, "*** %s NOT FOUND\n", officeType)
		newRow = append(*row, "*** office NOT FOUND")
		*row = newRow
		return 1
	}
	for _, t := range ts {
		//fmt.Fprintf(w, "\t%s: %s\n", officeType, t.Name)
		newRow = append(*row, t.Name)
		*row = newRow
	}
	return 0
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

func (b *MigrationBuilder) Build(dutyStationsFilePath string, outputFilePath string) (string, error) {
	// Parse raw data from xml
	stations, err := b.ParseStations(dutyStationsFilePath)
	if err != nil {
		return "", err
	}
	fmt.Printf("# total stations: %d\n", len(stations))
	f, err := os.Create(outputFilePath)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write([]string{"Duty Station Name", "ZIP", "Transportation Offices -->"})
	for _, s := range stations {
		var row []string
		fmt.Printf("%v\n", s)
		dbOffices := b.FindTransportationOffice(s)
		if len(dbOffices) > 0 {
			b.WriteLine(s, &row)
			b.WriteDbRecs(dbOffices, &row)
			w.Write(row)
		}
	}
	w.Flush()
	return "abc", nil
}
