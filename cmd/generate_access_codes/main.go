package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gobuffalo/pop"

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
	hhg := flag.Int("hhg", 0, "The number of access codes for HHG moves to create")
	ppm := flag.Int("ppm", 0, "The number of access codes for PPM moves to create")

	flag.Parse()

	// DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	if *hhg == 0 && *ppm == 0 {
		log.Fatal("Usage: generate_access_codes -ppm 1000 -hhg 4000")
	}
	// go run cmd/generate_access_codes/main.go -ppm 2000 -hhg 500
	for i := 0; i < *hhg; i++ {
		selectedMoveType := models.SelectedMoveTypeHHG

		accessCode := models.AccessCode{
			Code:     models.GenerateLocator(),
			MoveType: &selectedMoveType,
		}

		mustSave(db, &accessCode)
	}

	for i := 0; i < *ppm; i++ {
		selectedMoveType := models.SelectedMoveTypePPM

		accessCode := models.AccessCode{
			Code:     models.GenerateLocator(),
			MoveType: &selectedMoveType,
		}

		mustSave(db, &accessCode)
	}

	fmt.Println("Completed generate_access_codes")
}
