package transportationoffices

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

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

func ReadXMLFile(file string) []byte {
	xmlFile, err := os.Open("cmd/load_transportation_offices/data/" + file)
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

func Filter(os []Office, test func(Office) bool) []Office {
	var filtered []Office
	for _, o := range os {
		if test(o) {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

func CheckDbForConusOffices(db *pop.Connection, o Office) models.TransportationOffices {
	fmt.Printf("name: %s\n", o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME)
	dbOs, err := models.FetchTransportationOfficesByPostalCode(db, o.LISTGCNSLINFO.GCNSLINFO.CNSLZIP)
	if err != nil {
		fmt.Println(err)
	}
	return dbOs

}

func OutputResults(o Office, dbO models.TransportationOffices, w io.Writer) {
	fmt.Fprintf(w, "\nname: %s\n", o.LISTGCNSLINFO.GCNSLINFO.CNSLNAME)
	fmt.Fprintf(w, "city: %s | state: %s | zip: %s \n", o.LISTGCNSLINFO.GCNSLINFO.CNSLCITY, o.LISTGCNSLINFO.GCNSLINFO.CNSLSTATE, o.LISTGCNSLINFO.GCNSLINFO.CNSLZIP)
	printNames(dbO, w)
}

func printNames(ts models.TransportationOffices, w io.Writer) {
	if len(ts) == 0 {
		fmt.Printf("*** NOT FOUND\n")
		fmt.Fprintf(w, "*** NOT FOUND\n")
	}
	for _, t := range ts {
		_, _ = fmt.Fprintf(w, "\tdb: %v\n", t.Name)
	}
}

func OpenWriteFile() {

}
