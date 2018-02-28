package handlers

import (
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestIndexShipmentsHandler() {
	t := suite.T()

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "scac",
		Name: "Transportation Service Provider 1",
	}
	suite.mustSave(&tsp)

	tdl := models.TrafficDistributionList{
		CodeOfService:     "cos",
		DestinationRegion: "dr",
		SourceRateArea:    "sra",
	}
	suite.mustSave(&tdl)

	avs := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	suite.mustSave(&avs)

	aws := models.Shipment{
		TrafficDistributionListID: tdl.ID,
	}
	suite.mustSave(&aws)

	award := models.ShipmentAward{
		ShipmentID:                      aws.ID,
		TransportationServiceProviderID: tsp.ID,
	}
	suite.mustSave(&award)

	params := shipmentop.NewIndexShipmentsParams()
	indexHandler := NewIndexShipmentsHandler(suite.db, suite.logger)
	indexResponse := indexHandler.Handle(params)

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
