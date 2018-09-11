package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
	"go.uber.org/zap"
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

	err = db.Eager(
		"Move.Orders",
		"PickupAddress",
		"DeliveryAddress",
		"ServiceMember",
	).Where("shipment_offers.accepted=true").
		Where("move_id = $1", &moveID).
		Join("shipment_offers", "shipment_offers.shipment_id = shipments.id").
		All(&shipments)
	if err != nil {
		log.Fatal(err)
	}
	if len(shipments) == 0 {
		log.Fatal("No accepted shipments found")
	}
	var logger = zap.NewNop()

	var costsByShipments []ediinvoice.CostByShipment

	engine := rateengine.NewRateEngine(db, logger, route.NewTestingPlanner(362)) //TODO: create the propper route/planner
	for _, shipment := range shipments {
		costByShipment, err := HandleRunRateEngineOnShipment(shipment, engine)
		if err != nil {
			log.Fatal(err)
		}
		costsByShipments = append(costsByShipments, costByShipment)
	}
	edi, err := ediinvoice.Generate858C(costsByShipments, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(edi)
}

// HandleRunRateEngineOnShipment runs the rate engine on a shipment and returns the shipment and cost
func HandleRunRateEngineOnShipment(shipment models.Shipment, engine *rateengine.RateEngine) (ediinvoice.CostByShipment, error) {
	// Apply rate engine to shipment
	var shipmentCost ediinvoice.CostByShipment
	cost, err := engine.ComputeShipment(unit.Pound(*shipment.WeightEstimate),
		shipment.PickupAddress.PostalCode,
		shipment.DeliveryAddress.PostalCode,
		time.Time(*shipment.PickupDate),
		0,  // We don't want any SIT charges
		.4, // TODO: placeholder: story to get actual linehaul discount
		0.0,
	)
	if err != nil {
		return ediinvoice.CostByShipment{}, err
	}

	shipmentCost = ediinvoice.CostByShipment{
		Shipment: shipment,
		Cost:     cost,
	}
	return shipmentCost, err
}
