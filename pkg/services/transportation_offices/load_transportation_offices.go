package transportationoffices

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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

type OfficeData struct {
	XMLName        xml.Name `xml:"CONTACT_INFO_INITIAL_SITES_R2NEW"`
	Text           string   `xml:",chardata"`
	Name           string   `xml:"name"`
	LISTGCNSLORGID struct {
		Text       string   `xml:",chardata"`
		GCNSLORGID []Office `xml:"G_CNSL_ORG_ID"`
	} `xml:"LIST_G_CNSL_ORG_ID"`
}

type Office struct {
	XMLName       xml.Name `xml:"G_CNSL_ORG_ID"`
	Text          string   `xml:",chardata"`
	CNSLORGID1    string   `xml:"CNSL_ORG_ID1"`
	LISTGCNSLINFO struct {
		Text      string `xml:",chardata"`
		GCNSLINFO struct {
			Text                string `xml:",chardata"`
			CNSLCTRYNM          string `xml:"CNSL_CTRY_NM"`
			CNSLSTDIV           string `xml:"CNSL_ST_DIV"`
			CNSLNAME            string `xml:"CNSL_NAME"`
			CNSLADDR1           string `xml:"CNSL_ADDR1"`
			CNSLADDR2           string `xml:"CNSL_ADDR2"`
			CNSLCITY            string `xml:"CNSL_CITY"`
			CNSLSTATE           string `xml:"CNSL_STATE"`
			CNSLZIP             string `xml:"CNSL_ZIP"`
			CNSLCOUNTRY         string `xml:"CNSL_COUNTRY"`
			PPSOORGID           string `xml:"PPSO_ORG_ID"`
			PPSOCOUNTRY         string `xml:"PPSO_COUNTRY"`
			PPSOSTATE           string `xml:"PPSO_STATE"`
			PPSOZIP             string `xml:"PPSO_ZIP"`
			PPSONAME            string `xml:"PPSO_NAME"`
			PPSOCITY            string `xml:"PPSO_CITY"`
			PPSOADDR2           string `xml:"PPSO_ADDR2"`
			PPSOADDR1           string `xml:"PPSO_ADDR1"`
			LISTGPPSOEMAILORGID struct {
				Text            string `xml:",chardata"`
				GPPSOEMAILORGID struct {
					Text           string `xml:",chardata"`
					ORGIDP         string `xml:"ORG_IDP"`
					LISTGPPSOEMAIL struct {
						Text       string `xml:",chardata"`
						GPpsoEmail []struct {
							Text          string `xml:",chardata"`
							EMAILTYPEP    string `xml:"EMAIL_TYPEP"`
							EMAILADDRESSP string `xml:"EMAIL_ADDRESSP"`
						} `xml:"G_ppso_email"`
					} `xml:"LIST_G_PPSO_EMAIL"`
				} `xml:"G_PPSO_EMAIL_ORG_ID"`
			} `xml:"LIST_G_PPSO_EMAIL_ORG_ID"`
			LISTGPPSOPHONEORGID struct {
				Text            string `xml:",chardata"`
				GPPSOPHONEORGID struct {
					Text                string `xml:",chardata"`
					PPSOORGID2          string `xml:"PPSO_ORG_ID2"`
					LISTGPPSOPHONENOTES struct {
						Text            string `xml:",chardata"`
						GPPSOPHONENOTES []struct {
							Text           string `xml:",chardata"`
							PPSOVOICEORFAX string `xml:"PPSO_VOICE_OR_FAX"`
							PPSOPHONENUM   string `xml:"PPSO_PHONE_NUM"`
							PPSODSNNUM     string `xml:"PPSO_DSN_NUM"`
							PPSOPHONETYPE  string `xml:"PPSO_PHONE_TYPE"`
							PPSOCOMMORDSN  string `xml:"PPSO_COMM_OR_DSN"`
							PPSOPHONENOTES string `xml:"PPSO_PHONE_NOTES"`
						} `xml:"G_PPSO_PHONE_NOTES"`
					} `xml:"LIST_G_PPSO_PHONE_NOTES"`
				} `xml:"G_PPSO_PHONE_ORG_ID"`
			} `xml:"LIST_G_PPSO_PHONE_ORG_ID"`
		} `xml:"G_CNSL_INFO"`
	} `xml:"LIST_G_CNSL_INFO"`
	LISTGCNSLEMAILORGID struct {
		Text            string `xml:",chardata"`
		GCNSLEMAILORGID struct {
			Text           string `xml:",chardata"`
			ORGID2         string `xml:"ORG_ID2"`
			LISTGCNSLEMAIL struct {
				Text       string `xml:",chardata"`
				GCNSLEMAIL struct {
					Text         string `xml:",chardata"`
					EMAILTYPE    string `xml:"EMAIL_TYPE"`
					EMAILADDRESS string `xml:"EMAIL_ADDRESS"`
				} `xml:"G_CNSL_EMAIL"`
			} `xml:"LIST_G_CNSL_EMAIL"`
		} `xml:"G_CNSL_EMAIL_ORG_ID"`
	} `xml:"LIST_G_CNSL_EMAIL_ORG_ID"`
	LISTGCNSLPHONEORGID struct {
		Text            string `xml:",chardata"`
		GCNSLPHONEORGID struct {
			Text                string `xml:",chardata"`
			CNSLORGID2          string `xml:"CNSL_ORG_ID2"`
			LISTGCNSLPHONENOTES struct {
				Text            string `xml:",chardata"`
				GCNSLPHONENOTES []struct {
					Text           string `xml:",chardata"`
					CNSLVOICEORFAX string `xml:"CNSL_VOICE_OR_FAX"`
					CNSLDSNNUM     string `xml:"CNSL_DSN_NUM"`
					CNSLAREACODE   string `xml:"CNSL_AREA_CODE"`
					CNSLPHONETYPE  string `xml:"CNSL_PHONE_TYPE"`
					CNSLCOMMORDSN  string `xml:"CNSL_COMM_OR_DSN"`
					CNSLPHONENUM   string `xml:"CNSL_PHONE_NUM"`
					CNSLPHONENOTES string `xml:"CNSL_PHONE_NOTES"`
				} `xml:"G_CNSL_PHONE_NOTES"`
			} `xml:"LIST_G_CNSL_PHONE_NOTES"`
		} `xml:"G_CNSL_PHONE_ORG_ID"`
	} `xml:"LIST_G_CNSL_PHONE_ORG_ID"`
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

func (b *MigrationBuilder) parseOffices(path string) ([]Office, error) {
	fileBytes := ReadXMLFile(path)
	data := UnmarshalXML(fileBytes)
	officeData := data.LISTGCNSLORGID.GCNSLORGID

	return officeData, nil

}

func ReadXMLFile(file string) []byte {
	// xmlFile, err := os.Open("cmd/load_transportation_offices/data/" + file)
	xmlFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	return byteValue
}

func UnmarshalXML(byteValue []byte) OfficeData {
	var od OfficeData
	xml.Unmarshal(byteValue, &od)
	return od
}

func (b *MigrationBuilder) isUS(offices []Office) []Office {
	filter := func(o Office) bool {
		return o.LISTGCNSLINFO.GCNSLINFO.CNSLCOUNTRY == "US"
	}
	return b.filterOffice(offices, filter)
}

func (b *MigrationBuilder) isConus(offices []Office) []Office {
	filter := func(o Office) bool {
		return o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "AK" &&
			o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE != "HI"
	}
	return b.filterOffice(offices, filter)
}

func (b *MigrationBuilder) filterOffice(os []Office, test func(Office) bool) []Office {
	var filtered []Office
	for _, o := range os {
		if test(o) {
			filtered = append(filtered, o)
		}
	}
	return filtered
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

func (b *MigrationBuilder) FindConusOffices(o Office, w io.Writer) models.TransportationOffices {
	zip := o.LISTGCNSLINFO.GCNSLINFO.CNSLZIP

	dbOs, err := models.FetchTransportationOfficesByPostalCode(b.db, zip)
	if err != nil {
		fmt.Println(err)
	}

	if len(dbOs) == 0 {
		partialZip := zip[:len(zip)-1] + "%"
		fmt.Fprintf(w, "*** partialZip: %s \n", partialZip)
		dbOs, err = models.FetchTransportationOfficesByPostalCode(b.db, partialZip)
		if err != nil {
			fmt.Println(err)
		}
	}

	return dbOs
}

func (b *MigrationBuilder) FindPPSOs(o Office) models.TransportationOffices {
	zip := o.LISTGCNSLINFO.GCNSLINFO.PPSOZIP
	dbPPSOs, _ := models.FetchTransportationOfficesByPostalCode(b.db, zip)

	JPPSOFilter := func(o models.TransportationOffice) bool {
		return strings.HasPrefix(o.Name, "JPPSO") // true
	}

	return FilterTransportationOffices(dbPPSOs, JPPSOFilter)
}

func (b *MigrationBuilder) WriteXMLLine(o Office, w io.Writer) {
	name := b.normalizeName(o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME)
	city := o.LISTGCNSLINFO.GCNSLINFO.CNSLCITY
	state := o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE
	zip := o.LISTGCNSLINFO.GCNSLINFO.CNSLZIP

	fmt.Printf("\nname: %s\n", name)
	fmt.Fprintf(w, "\nname: %s\n", name)
	fmt.Printf("city: %s | state: %s | zip: %s \n", city, state, zip)
	fmt.Fprintf(w, "city: %s | state: %s | zip: %s \n", city, state, zip)
}

func (b *MigrationBuilder) WriteDbRecs(officeType string, ts models.TransportationOffices, w io.Writer) int {
	if len(ts) == 0 {
		// b.logger.Debug("*** NONE FOUND... BLAH")
		fmt.Printf("*** %s NOT FOUND\n", officeType)
		fmt.Fprintf(w, "*** %s NOT FOUND\n", officeType)

		return 1
	}
	for _, t := range ts {
		fmt.Fprintf(w, "\t%s: %s\n", officeType, t.Name)
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

func (b *MigrationBuilder) Build(officesFilePath string, outputFilePath string) (string, error) {
	// Parse raw data from xml
	offices, err := b.parseOffices(officesFilePath)
	if err != nil {
		return "", err
	}
	fmt.Printf("# total offices: %d\n", len(offices))

	usOffices := b.isUS(offices)
	fmt.Printf("# us only offices: %d\n", len(usOffices))

	conusOffices := b.isConus(usOffices)
	fmt.Printf("# conus only offices: %d\n", len(conusOffices))

	fmt.Println(outputFilePath)
	f, err := os.Create(outputFilePath)
	defer f.Close()
	w := bufio.NewWriter(f)

	for _, o := range conusOffices {
		b.WriteXMLLine(o, w)
		dbOffices := b.FindConusOffices(o, w)
		dbPPSOs := b.FindPPSOs(o)
		b.WriteDbRecs("office", dbOffices, w)
		b.WriteDbRecs("JPPSO", dbPPSOs, w)
	}
	w.Flush()
	return "abc", nil

}
