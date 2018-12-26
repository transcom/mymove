package main

import (
	"fmt"
	"github.com/transcom/mymove/pkg/edi/gex"
	"github.com/transcom/mymove/pkg/models"
	"log"
	"os"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/namsral/flag"

	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/service/invoice"
)

// Call this from command line with go run cmd/generate_shipment_edi/main.go -shipmentID <UUID> --approver <email>
func main() {
	shipmentIDString := flag.String("shipmentID", "", "The ID of the shipment to invoice")
	approverEmail := flag.String("approver", "", "The office approver e-mail")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	sendToGex := flag.Bool("gex", false, "Choose to send the file to gex")
	transactionName := flag.String("transactionName", "test", "The required name sent in the url of the gex api request")
	flag.Parse()

	if *shipmentIDString == "" || *approverEmail == "" {
		log.Fatal("Usage: go run cmd/generate_shipment_edi/main.go --shipmentID <29cb984e-c70d-46f0-926d-cd89e07a6ec3> --approver <officeuser1@example.com> --gex false")
	}

	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	shipmentID := uuid.Must(uuid.FromString(*shipmentIDString))
	shipment, err := invoice.FetchShipmentForInvoice{DB: db}.Call(shipmentID)
	if err != nil {
		log.Fatal(err)
	}

	approver, err := models.FetchOfficeUserByEmail(db, *approverEmail)
	if err != nil {
		log.Fatalf("Could not fetch office user with e-mail %s: %v", *approverEmail, err)
	}

	var invoiceModel models.Invoice
	verrs, err := invoice.CreateInvoice{DB: db, Clock: clock.New()}.Call(*approver, &invoiceModel, shipment)
	if err != nil {
		log.Fatal(err)
	}
	if verrs.HasAny() {
		log.Fatal(verrs)
	}

	invoice858C, err := ediinvoice.Generate858C(shipment, invoiceModel, db, false, clock.New())
	if err != nil {
		log.Fatal(err)
	}

	if *sendToGex == true {
		fmt.Println("Sending to GEX. . .")
		invoice858CString, err := invoice858C.EDIString()
		if err != nil {
			log.Fatal(err)
		}
		resp, err := gex.SendInvoiceToGex(invoice858CString, *transactionName)
		fmt.Printf("status code: %v, error: %v", resp.StatusCode, err)

		if resp.StatusCode != 200 {
			// Update invoice record as failed
			invoiceModel.Status = models.InvoiceStatusSUBMISSIONFAILURE
			verrs, err := db.ValidateAndSave(&invoiceModel)
			if err != nil {
				log.Fatal(err)
			}
			if verrs.HasAny() {
				log.Fatal(verrs)
			}
		} else {
			// Update invoice record as submitted
			shipmentLineItems := shipment.ShipmentLineItems
			verrs, err = invoice.UpdateInvoiceSubmitted{DB: db}.Call(&invoiceModel, shipmentLineItems)
			if err != nil {
				log.Fatal(err)
			}
			if verrs.HasAny() {
				log.Fatal(verrs)
			}
		}
	} else {
		ediWriter := edi.NewWriter(os.Stdout)
		ediWriter.WriteAll(invoice858C.Segments())
	}
}
