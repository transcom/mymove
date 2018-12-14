package main

import (
	"fmt"
	"log"
	"os"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/gex"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/service/invoice"
)

// Call this from command line with go run cmd/generate_shipment_edi/main.go -moveID <UUID>
func main() {
	shipmentIDString := flag.String("shipmentID", "", "The ID of the shipment to invoice")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	sendToGex := flag.Bool("gex", false, "Choose to send the file to gex")
	transactionName := flag.String("transactionName", "test", "The required name sent in the url of the gex api request")
	flag.Parse()

	if *shipmentIDString == "" {
		log.Fatal("Usage: go run cmd/generate_shipment_edi/main.go --shipmentID <29cb984e-c70d-46f0-926d-cd89e07a6ec3> --gex false")
	}

	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	shipmentID := uuid.Must(uuid.FromString(*shipmentIDString))
	shipment, err := invoice.FetchShipmentForInvoice{db}.Call(shipmentID)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	invoice858C, err := ediinvoice.Generate858C(shipment, db, false, clock.New())
	if err != nil {
		log.Fatal(err)
	}

	if *sendToGex == true {
		fmt.Println("Sending to GEX. . .")
		invoice858CString, err := invoice858C.EDIString()
		if err != nil {
			log.Fatal(err)
		}
		statusCode, err := gex.SendInvoiceToGex(logger, invoice858CString, *transactionName)
		fmt.Printf("status code: %v, error: %v", statusCode, err)
	} else {
		ediWriter := edi.NewWriter(os.Stdout)
		ediWriter.WriteAll(invoice858C.Segments())
	}
}
