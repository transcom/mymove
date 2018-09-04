package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
)

// Call this from command line with go run cmd/generate_shipment_edi/main.go -moveID <UUID>
func main() {
	moveIDString := flag.String("moveID", "", "The ID of the move where shipments are found")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	flag.Parse()

	if *moveIDString == "" {
		log.Fatal("Usage: generate_shipment_edi -moveID <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	moveID := uuid.Must(uuid.FromString(*moveIDString))
	var shipments models.Shipments
	err = db.Where("move_id = $1", &moveID).All(&shipments)
	if err != nil {
		log.Fatal(err)
	}

	edi, err := ediinvoice.Generate858C(shipments, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(edi)
}
