package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeShipment creates a single shipment record.
func MakeShipment(db *pop.Connection, pickup time.Time, delivery time.Time,
	tdl models.TrafficDistributionList) (models.Shipment, error) {

	market := "dHHG"
	GBLOC := "AGFM"
	shipment := models.Shipment{
		TrafficDistributionListID: tdl.ID,
		PickupDate:                pickup,
		DeliveryDate:              delivery,
		GBLOC:                     &GBLOC,
		Market:                    &market,
		AwardDate:                 DateInsidePerformancePeriod,
	}

	_, err := db.ValidateAndSave(&shipment)
	if err != nil {
		log.Panic(err)
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

	MakeShipment(db, now, now.Add(thirtyMin), tdlList[0])
	MakeShipment(db, now.Add(oneWeek), now.Add(oneWeek).Add(thirtyMin), tdlList[1])
	MakeShipment(db, now.Add(oneWeek*2), now.Add(oneWeek*2).Add(thirtyMin), tdlList[2])
}
