package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"testing"

	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/operations/shipments"
	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/testing"
)

func TestIndexShipmentsHandler(t *testing.T) {
	tx, rollback := StartTransaction(t, DB)
	defer rollback()

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "scac",
		Name: "Transportation Service Provider 1",
	}
	MustSave(t, tx, &tsp)

	tdl := models.TrafficDistributionList{
		CodeOfService:     "cos",
		DestinationRegion: "dr",
		SourceRateArea:    "sra",
	}
	MustSave(t, tx, &tdl)

	avs := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	MustSave(t, tx, &avs)

	aws := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	MustSave(t, tx, &aws)

	award := models.ShipmentAward{
		ShipmentID:                      aws.ID,
		TransportationServiceProviderID: tsp.ID,
	}
	MustSave(t, tx, &award)

	var indexResponse middleware.Responder
	stubDB(tx, func() {
		params := shipmentop.NewIndexShipmentsParams()
		indexResponse = IndexShipmentsHandler(params)
	})

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
