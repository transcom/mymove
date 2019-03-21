package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags

	"github.com/transcom/mymove/pkg/models"
)

/* load-office-data is a tool to load Transportation Office data into a local data base.
From there it will be exported
as part of a migration
*/

// OfficeLocation is the form of the `location` object in json representation of PPPO's
type OfficeLocation struct {
	Address1    string  `json:"street_address"`
	Address2    string  `json:"extended_address"`
	Locality    string  `json:"locality"`
	State       string  `json:"region"`
	StateCode   string  `json:"region_code"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
}

// JSONEmail is the form of the items in the email_addresses list
type JSONEmail struct {
	EmailAddress string  `json:"email_address"`
	Note         *string `json:"note"`
}

// JSONPhoneNumber is the form of the items in the phone_numbers list
type JSONPhoneNumber struct {
	Number string `json:"phone_number"`
	Type   string `json:"phone_type"`
	DSN    bool   `json:"dsn"`
}

// JSONOffice is the form of the objects listed in the PPPO json file
type JSONOffice struct {
	Name               string            `json:"name"`
	ShippingOfficeName *string           `json:"shipping_office_name"`
	Location           OfficeLocation    `json:"location"`
	Hours              string            `json:"hours"`
	Services           []string          `json:"services"`
	EmailAddresses     []JSONEmail       `json:"email_addresses"`
	PhoneNumbers       []JSONPhoneNumber `json:"phone_numbers"`
}

/* This is code used to parse the XML data from the DPS database.
We are not (yet) using that as it seems dirty, e.g. there are two entries for the Fort Stewart PPPO each with different
PPSO's listed.
This code will need to come back i once we have some clarity around this and need to HOOK up PPSOs
// XMLOffice contains the fields for a PPPO or PPPSO in the xlsx file
type XMLOffice struct {
	Name           string
	Address1       string
	Address2       string
	City           string
	State          string
	Zip            string
	Country        string
	ShippingOffice *XMLOffice
}

func xmlOfficeFromStringSlice(slice []string) *XMLOffice {
	return &XMLOffice{
		Name:     strings.TrimSpace(slice[0]),
		Address1: strings.TrimSpace(slice[1]),
		Address2: strings.TrimSpace(slice[2]),
		City:     strings.TrimSpace(slice[3]),
		State:    strings.TrimSpace(slice[4]),
		Zip:      strings.TrimSpace(slice[5]),
		Country:  strings.TrimSpace(slice[6]),
	}
}
*/
func validated(verrs *validate.Errors, err error) (retError error) {
	if verrs.HasAny() {
		log.Printf("Validation errors:- %v", verrs)
		retError = verrs
	}
	if err != nil {
		log.Printf("Creation error:- %v", err)
		retError = err
	}
	return
}

func splitOnNewLineAndCommaThenJoin(s string) string {
	var all []string
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		partsOfLine := strings.Split(line, ",")
		for _, part := range partsOfLine {
			all = append(all, strings.TrimSpace(part))
		}
	}
	return strings.Join(all, "; ")
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	jsonFile := flag.String("jsonFile", "", "Json File name to load - should be a *.json file containing PPO data")
	// xmlFile := flag.String("xmlFile", "", "XML File name to load - should be a *.xlsx file containing PPPO-PPSO relationships")
	flag.Parse()

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}
	if *jsonFile == "" { //|| *xmlFile == "" {
		log.Fatal("usage: load-office-data -jsonFile <filename> ") //  -xmlFile <other_filename>")
	}

	var offices []JSONOffice
	jf, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		log.Fatalf("Error %v opening JSON %s.", err, *jsonFile)
	}
	err = json.Unmarshal(jf, &offices)
	if err != nil {
		log.Fatalf("Error %v parsing json", err)
	}

	err = db.Transaction(func(connection *pop.Connection) error {
		for _, office := range offices {
			if office.Location.CountryCode != "US" {
				continue
			}

			// Address
			address := models.Address{
				StreetAddress1: office.Location.Address1,
				StreetAddress2: &office.Location.Address2,
				City:           office.Location.Locality,
				State:          office.Location.StateCode,
				PostalCode:     office.Location.PostalCode,
				Country:        &office.Location.Country,
			}
			err = validated(db.ValidateAndCreate(&address))
			if err != nil {
				log.Printf("%v - address.", office.Name)
				return err
			}

			// Office
			hours := splitOnNewLineAndCommaThenJoin(office.Hours)
			services := splitOnNewLineAndCommaThenJoin(strings.Join(office.Services, ", "))
			transportationOffice := models.TransportationOffice{
				Name:      office.Name,
				AddressID: address.ID,
				Latitude:  office.Location.Latitude,
				Longitude: office.Location.Longitude,
				Hours:     &hours,
				Services:  &services,
			}
			err = validated(db.ValidateAndCreate(&transportationOffice))
			if err != nil {
				log.Printf("%v - office.", office.Name)
				return err
			}

			// Phone numbers
			for _, jsonPhone := range office.PhoneNumbers {
				// Some numbers are listed as alternates, e.g. "123 555-2323 / 123 555-2324"
				for _, num := range strings.Split(jsonPhone.Number, "/") {
					phone := models.OfficePhoneLine{
						TransportationOfficeID: transportationOffice.ID,
						Number:                 strings.TrimSpace(num),
						Type:                   jsonPhone.Type,
						IsDsnNumber:            jsonPhone.DSN,
					}
					err = validated(db.ValidateAndCreate(&phone))
					if err != nil {
						log.Printf("%v - office, %v - number", office.Name, phone.Number)
						return err
					}
				}
			}

			for _, jsonEmail := range office.EmailAddresses {
				email := models.OfficeEmail{
					TransportationOfficeID: transportationOffice.ID,
					Email:                  strings.TrimSpace(jsonEmail.EmailAddress),
					Label:                  jsonEmail.Note,
				}
				err = validated(db.ValidateAndCreate(&email))
				if err != nil {
					log.Printf("%v - office, %v - email", office.Name, email)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Tranaction failed:- %v", err)
	}

	/*	xf, err := excelize.OpenFile(*xmlFile)
		if err != nil {
			log.Fatalf("Error %v opening XLSX %s.", err, *xmlFile)
		}

		ppoByZip := make(map[string][]*XMLOffice)
		for idx, row := range xf.GetRows(xf.GetSheetName(1)) {
			if idx > 0 {
				pppo := xmlOfficeFromStringSlice(row)
				list, ok := ppoByZip[pppo.Zip]
				if ok {
					ppoByZip[pppo.Zip] = append(list, pppo)
				} else {
					ppoByZip[pppo.Zip] = []*XMLOffice{pppo}
				}
				pppo.ShippingOffice = xmlOfficeFromStringSlice(row[7:])

			}
		}

		officeMap := make(map[string]JSONOffice)
		for _, office := range offices {
			if office.Location.CountryCode != "US" {
				continue
			}
			officeMap[office.Name] = office
			pppos, ok := ppoByZip[office.Location.PostalCode]
			if !ok {
				log.Printf("Couldn't find pppos for %s", office.Name)
			}
			if len(pppos) > 1 {
				log.Printf("Ambiguous ppos for %s", office.Name)
				for _, p := range pppos {
					log.Printf("\t%s", p.Name)
				}
			}
		}
	*/
	log.Printf("Loaded %d offices", len(offices))
}
