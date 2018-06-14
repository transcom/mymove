package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/namsral/flag"

	"github.com/transcom/mymove/pkg/paperwork"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")

	moveID := flag.String("move", "", "The move ID to generate advance paperwork for")
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

	if *moveID == "" {
		log.Fatal("Usage: paperwork -move <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	id := uuid.Must(uuid.FromString(*moveID))
	if err = paperwork.GenerateAdvancePaperwork(db, id); err != nil {
		log.Fatal(err)
	}
}
