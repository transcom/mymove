package main

import (
	"log"

	"github.com/gobuffalo/pop"
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
		log.Fatal("Usage: make_dps_user -email <my_email@example.com>")
	}

	newUser := models.DpsUser{
		LoginGovEmail: *email,
	}

	mustSave(db, &newUser)
}
