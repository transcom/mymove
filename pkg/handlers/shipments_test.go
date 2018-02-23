package handlers

import (
	"log"
	"testing"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func mustSave(t *testing.T, s interface{}) {
	verrs, err := dbConnection.ValidateAndSave(s)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", s, verrs)
	}
}

func TestIndexShipmentsHandler(t *testing.T) {
	// TODO: This test assumes an empty database. This can be removed once we have
	// a more centralized way of handling test setup & teardown.
	dbConnection.TruncateAll()

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "scac",
		Name: "Transportation Service Provider 1",
	}
	mustSave(t, &tsp)

	tdl := models.TrafficDistributionList{
		CodeOfService:     "cos",
		DestinationRegion: "dr",
		SourceRateArea:    "sra",
	}
	mustSave(t, &tdl)

	avs := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	mustSave(t, &avs)

	aws := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	mustSave(t, &aws)

	award := models.ShipmentAward{
		ShipmentID:                      aws.ID,
		TransportationServiceProviderID: tsp.ID,
	}
	mustSave(t, &award)

	params := shipmentop.NewIndexShipmentsParams()
	indexResponse := IndexShipmentsHandler(params)

	okResponse, ok := indexResponse.(*shipmentop.IndexShipmentsOK)
	if !ok {
		t.Fatalf("Request failed: %#v", indexResponse)
	}
	shipments := okResponse.Payload

	if len(shipments) != 2 {
		t.Errorf("expected %d shipments, got %d", 2, len(shipments))
	}

	awardedCount := 0
	availableCount := 0
	for _, shipment := range shipments {
		if shipment.TransportationServiceProviderID != nil {
			awardedCount++
			if shipment.TransportationServiceProviderID.String() != tsp.ID.String() {
				t.Errorf("got wrong tsp id, expected %s, got %s", tsp.ID.String(), shipment.TransportationServiceProviderID.String())

			}
		} else {
			availableCount++
		}
	}

	if awardedCount != 1 {
		t.Errorf("expected %d awarded shipments, got %d", 1, awardedCount)
	}

	if availableCount != 1 {
		t.Errorf("expected %d available shipments, got %d", 1, availableCount)
	}
}
