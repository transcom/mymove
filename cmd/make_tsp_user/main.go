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
	email := flag.String("email", "", "The email of the TSP user to create")
	firstName := flag.String("first_name", "Testy", "First name of the TSP user to create")
	lastName := flag.String("last_name", "McTester", "Last name of the TSP user to create")
	number := flag.String("number", "415-555-1212", "Phone number of the TSP user to create")
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
		log.Fatal("Usage: make_tsp_user -email <my_email@example.com>")
	}

	// Attempt to load an existing user.
	var user models.User
	db.Where("login_gov_email = $1", *email).Last(&user)

	// Attempt to load the Truss TSP
	var tsp models.TransportationServiceProvider
	db.Where("standard_carrier_alpha_code = $1", "TRS1").First(&tsp)
	if tsp.ID == uuid.Nil {
		// TSP not found, create one for Truss
		tsp = models.TransportationServiceProvider{
			StandardCarrierAlphaCode: "TRSS",
		}
		mustSave(db, &tsp)
	}

	newUser := models.TspUser{
		FirstName:                       *firstName,
		LastName:                        *lastName,
		Telephone:                       *number,
		TransportationServiceProviderID: tsp.ID,
		Email:                           *email,
	}
	if user.ID != uuid.Nil {
		newUser.UserID = &user.ID
	}

	mustSave(db, &newUser)
}
