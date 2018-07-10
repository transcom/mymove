package testdatagen

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeShipment creates a single shipment record.
func MakeShipment(db *pop.Connection, requestedPickup time.Time,
	pickup time.Time, delivery time.Time,
	tdl models.TrafficDistributionList, sourceGBLOC string, market *string) (models.Shipment, error) {

	shipment := models.Shipment{
		TrafficDistributionListID: tdl.ID,
		PickupDate:                pickup,
		RequestedPickupDate:       requestedPickup,
		DeliveryDate:              delivery,
		BookDate:                  DateInsidePerformancePeriod,
		SourceGBLOC:               sourceGBLOC,
		Market:                    market,
		Status:                    "DEFAULT",
		// TODO add default values for new fields here
	}

	verrs, err := db.ValidateAndSave(&shipment)
	if verrs.HasAny() {
		err = fmt.Errorf("shipment validation errors: %v", verrs)
	}

	return shipment, err
}

// MakeShipmentData creates three shipment records
func MakeShipmentData(db *pop.Connection) {
	// Grab three UUIDs for individual TDLs
	// TODO: should this query be made in main, between creation functions,
	// and then sourced from one central place?
	tdlList := []models.TrafficDistributionList{}
	err := db.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	// Add three shipment table records using UUIDs from TDLs
	now := time.Now()
	thirtyMin, _ := time.ParseDuration("30m")
	oneWeek, _ := time.ParseDuration("7d")
	market := "dHHG"
	sourceGBLOC := "OHAI"

	MakeShipment(db, now, now, now.Add(thirtyMin), tdlList[0], sourceGBLOC, &market)
	MakeShipment(db, now.Add(oneWeek), now.Add(oneWeek), now.Add(oneWeek).Add(thirtyMin), tdlList[1], sourceGBLOC, &market)
	MakeShipment(db, now.Add(oneWeek*2), now.Add(oneWeek*2), now.Add(oneWeek*2).Add(thirtyMin), tdlList[2], sourceGBLOC, &market)
}
