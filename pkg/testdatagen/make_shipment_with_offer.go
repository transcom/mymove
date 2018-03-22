package testdatagen

import (
  "fmt"
  "log"
  "time"

  "github.com/markbates/pop"
  "github.com/transcom/mymove/pkg/models"
)

// MakeShipmentWithOffer creates a single ShipmentWithOffer record.
func MakeShipmentWithOffer(db *pop.Connection, book time.Time, pickup time.Time, requestedPickup time.Time,
  tdl models.TrafficDistributionList, tsp models.TransportationServiceProvider, admin *bool) (models.ShipmentWithOffer, error) {

  trueBool := true
  reason := "not feeling it"
  shipmentWithOffer := models.ShipmentWithOffer{
    BookDate:                        book,
    PickupDate:                      pickup,
    RequestedPickupDate:             requestedPickup,
    TrafficDistributionListID:       tdl.ID,
    TransportationServiceProviderID: &tsp.ID,
    Accepted:                        &trueBool,
    RejectionReason:                 &reason,
    AdministrativeShipment:          admin,
  }

  _, err := db.ValidateAndSave(&shipmentWithOffer)
  if err != nil {
    log.Panic(err)
  }

  return shipmentWithOffer, err
}

// MakeShipmentData creates three shipment records
func MakeShipmentWithOfferData(db *pop.Connection) {
  // Grab three UUIDs for individual TDLs
  tdlList := []models.TrafficDistributionList{}
  err := db.All(&tdlList)
  if err != nil {
    fmt.Println("TDL ID import failed.")
  }

  // Grab three UUIDs for individual TSPs
  tspList := []models.TransportationServiceProvider{}
  err = db.All(&tspList)
  if err != nil {
    fmt.Println("TSP ID import failed.")
  }

  // Two variables to refer to in pointers
  trueBool := true
  falseBool := false

  // Some times to create diversity of sample dates away from time.Now
  now := time.Now()
  thirtyMin, _ := time.ParseDuration("30m")
  oneWeek, _ := time.ParseDuration("7d")

  MakeShipmentWithOffer(db, now, now.Add(thirtyMin), now.Add(oneWeek), tdlList[0], tspList[0], &trueBool)
  MakeShipmentWithOffer(db, now.Add(oneWeek), now.Add(oneWeek), now.Add(oneWeek).Add(thirtyMin), tdlList[1], tspList[1], &falseBool)
  MakeShipmentWithOffer(db, now.Add(oneWeek*2), now.Add(oneWeek*3), now.Add(oneWeek*4), tdlList[2], tspList[2], &falseBool)
}
