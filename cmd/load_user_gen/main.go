package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags

	"github.com/transcom/mymove/pkg/models"
)

const (
	fnameIndex = iota
	lnameIndex
	emailIndex
	baseIndex
)

func checkHeader(headers []string, index int, name string) {
	if strings.ToLower(headers[index]) != name {
		log.Fatalf("Expected column %v to be %v, got %v", index, name, headers[index])
	}
}
func checkHeaders(headers []string) {
	checkHeader(headers, fnameIndex, "fname")
	checkHeader(headers, lnameIndex, "lname")
	checkHeader(headers, emailIndex, "email")
	checkHeader(headers, baseIndex, "office")
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	test := flag.Bool("test", false, "Whether to generate testy mcTest emails")
	flag.Parse()

	r := csv.NewReader(os.Stdin)
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Reading headers - %v", err)
	}
	checkHeaders(headers)

	//DB connection
	err = pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		baseName := row[baseIndex]
		var offices []models.TransportationOffice
		err = db.Where("lower(name) like '%' || $1 || '%'", strings.ToLower(baseName)).Eager("PhoneLines").All(&offices)
		if err != nil {
			log.Fatalf("Looking up bases - %v", err)
		} else if len(offices) == 0 {
			log.Printf("Couldn't find office - %v", baseName)
			continue
		} else if len(offices) > 1 {
			log.Printf("More than one office matches - %v", baseName)
			for _, o := range offices {
				log.Printf("\t%s", o.Name)
			}
		}
		office := offices[0]
		var number string
		for _, line := range office.PhoneLines {
			if line.Type == "voice" {
				number = line.Number
				break
			}
		}
		if number == "" {
			fmt.Printf("We don't have a number for - %v", baseName)
		}
		id, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		if *test {
			fmt.Printf("INSERT INTO public.office_users VALUES ('%s', NULL, 'McTest', 'Testy-%s', NULL, 'test-%s@example.com', '(415) 555-1212', '%s', now(), now());\n",
				id.String(),
				id.String(),
				id.String(),
				offices[0].ID.String())

		} else {
			fmt.Printf("INSERT INTO public.office_users VALUES ('%s', NULL, '%s', '%s', NULL, '%s', '%s', '%s', now(), now());\n",
				id.String(),
				row[lnameIndex],
				row[fnameIndex],
				row[emailIndex],
				number,
				offices[0].ID.String())

		}
	}
}
