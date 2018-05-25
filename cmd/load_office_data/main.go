package main

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"io/ioutil"
	"log"
	"strings"
)

// OfficeLocation is the form of the `location` object in json representation of PPPO's
type OfficeLocation struct {
	Address1    string  `json:"street_address"`
	Address2    string  `json:"extended_address"`
	Locality    string  `json:"locality"`
	State       string  `json:"region"`
	StateCode   string  `json:"TX"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
}

// JSONOffice is the form of the objects listed in the PPPO json file
type JSONOffice struct {
	Name               string         `json:"name"`
	ShippingOfficeName *string        `json:"shipping_office_name"`
	Location           OfficeLocation `json:"location"`
	parent             *JSONOffice
}

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

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	// env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	jsonFile := flag.String("jsonFile", "", "Json File name to load - should be a *.json file containing PPO data")
	xmlFile := flag.String("xmlFile", "", "XML File name to load - should be a *.xlsx file containing PPPO-PPSO relationships")
	flag.Parse()

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	/*	db, err := pop.Connect(*env)
		if err != nil {
			log.Fatal(err)
		}
	*/
	if *jsonFile == "" || *xmlFile == "" {
		log.Fatal("usage: load-office-data -jsonFile <filename> -xmlFile <other_filename>")
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

	xf, err := excelize.OpenFile(*xmlFile)
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
	log.Printf("Loaded %d offices", len(offices))
}
