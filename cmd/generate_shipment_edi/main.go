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
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
)

// Call this from command line with go run cmd/generate_shipment_edi/main.go -moveID <UUID>
func main() {
	moveIDString := flag.String("moveID", "", "The ID of the move where shipments are found")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	sendToGex := flag.Bool("gex", false, "Choose to send the file to gex")
	transactionName := flag.String("transactionName", "test", "The required name sent in the url of the gex api request")
	hereGeoEndpoint := flag.String("here_maps_geocode_endpoint", "", "URL for the HERE maps geocoder endpoint")
	hereRouteEndpoint := flag.String("here_maps_routing_endpoint", "", "URL for the HERE maps routing endpoint")
	hereAppID := flag.String("here_maps_app_id", "", "HERE maps App ID for this application")
	hereAppCode := flag.String("here_maps_app_code", "", "HERE maps App API code")
	flag.Parse()

	if *moveIDString == "" {
		log.Fatal("Usage: go run cmd/generate_shipment_edi/main.go --moveID <29cb984e-c70d-46f0-926d-cd89e07a6ec3> --gex false")
	}

	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	moveID := uuid.Must(uuid.FromString(*moveIDString))
	var shipments models.Shipments

	err = db.Eager(
		"PickupAddress",
		"Move.Orders.NewDutyStation.Address",
		"ServiceMember",
		"ShipmentOffers.TransportationServiceProviderPerformance",
		"ShipmentLineItems.Tariff400ngItem",
	).Where("shipment_offers.accepted=true").
		Where("move_id = $1", &moveID).
		Join("shipment_offers", "shipment_offers.shipment_id = shipments.id").
		All(&shipments)
	if err != nil {
		log.Fatal(err)
	}
	if len(shipments) == 0 {
		log.Fatal("No shipments with accepted shipment offers found")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	planner := route.NewHEREPlanner(logger, *hereGeoEndpoint, *hereRouteEndpoint, *hereAppID, *hereAppCode)
	var costsByShipments []rateengine.CostByShipment

	engine := rateengine.NewRateEngine(db, logger, planner)
	for _, shipment := range shipments {
		costByShipment, err := engine.HandleRunOnShipment(shipment)
		if err != nil {
			log.Fatal(err)
		}
		costsByShipments = append(costsByShipments, costByShipment)
	}
	invoice858C, err := ediinvoice.Generate858C(costsByShipments, db, false, clock.New())
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
