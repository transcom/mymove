package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/namsral/flag"

	"github.com/transcom/mymove/pkg/models"
)

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if verrs.HasAny() {
		log.Fatalf("validation Errors %v", verrs)
	}
	if err != nil {
		log.Fatalf("Failed to save %v", err)
	}
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	email := flag.String("email", "", "The email of the office user to create")
	firstName := flag.String("first_name", "Testy", "First name of the office user to create")
	lastName := flag.String("last_name", "McTester", "Last name of the office user to create")
	number := flag.String("number", "415-555-1212", "Phone number of the office user to create")
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

	if *email == "" {
		log.Fatal("Usage: make_office_user -email <my_email@example.com>")
	}

	// Attempt to load an existing user.
	var user models.User
	db.Where("login_gov_email = $1", *email).Last(&user)

	// Now create the Truss JPPSO
	address := models.Address{
		StreetAddress1: "1333 Minna St",
		City:           "San Francisco",
		State:          "CA",
		PostalCode:     "94115",
	}
	mustSave(db, &address)
	office := models.TransportationOffice{
		Name:      "Truss",
		AddressID: address.ID,
		Latitude:  37.7678355,
		Longitude: -122.4199298,
		Hours:     models.StringPointer("0900-1800 Mon-Sat"),
	}
	mustSave(db, &office)

	newUser := models.OfficeUser{
		FirstName:              *firstName,
		LastName:               *lastName,
		Telephone:              *number,
		TransportationOfficeID: office.ID,
		Email:                  *email,
	}
	if user.ID != uuid.Nil {
		newUser.UserID = &user.ID
	}

	mustSave(db, &newUser)
}
