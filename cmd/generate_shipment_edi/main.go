package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
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
	err = db.Where("move_id = $1", &moveID).Eager(
		"Move.Orders",
		"PickupAddress",
		"DeliveryAddress",
		"ServiceMember",
	).All(&shipments)
	if err != nil {
		log.Fatal(err)
	}

	logger := zap.NewNop()                             //??
	engine := rateengine.NewRateEngine(db, log, route) //route?

	type CostByShipment struct {
		shipment models.Shipment
		cost     CostComputation
	}

	// func HandleRunRateEngineOnShipment(shipment models.Shipment, engine *RateEngine) (CostByShipment, error) {
	// 	// Apply rate engine to shipment
	// 	var shipmentCost CostByShipment
	//
	// 	cost, err := engine.ComputeHHG(unit.Pounds(shipment.WeightEstimate),
	// 		shipment.PickupAddress.PostalCode,   // how to get nested value?
	// 		shipment.DeliveryAddress.PostalCode, // how to get nested value?
	// 		time.Time(shipment.PickupDate),
	// 		0,          // We don't want any SIT charges
	// 		lhDiscount, // where to get this? same as PPMDiscountFetch?
	// 		0.0,
	// 	)
	// 	if err != nil {
	// 		return "", err
	// 	}
	//
	// 	shipmentCost = CostbyShipment{
	// 		shipment,
	// 		cost,
	// 	}
	// 	return shipmentCost, err
	// }

	var costsByShipments []models.Shipment
	for _, shipment := range shipments {
		engine := rateengine.NewRateEngine(db, log, route) //route?
		costByShipment := HandleRunRateEngineOnShipment(shipment, engine)
		costsByShipments += costByShipment
	}

	edi, err := ediinvoice.Generate858C(costsByShipments, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(edi)
}
