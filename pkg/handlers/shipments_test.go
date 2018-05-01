package handlers

import (
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestIndexShipmentsHandler() {
	t := suite.T()

	tspName := "Transportation Service Provider 1"
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "scac",
		Name: &tspName,
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
		SourceGBLOC:               "AGFM",
	}
	suite.mustSave(&avs)

	aws := models.Shipment{
		TrafficDistributionListID: tdl.ID,
		SourceGBLOC:               "AGFM",
	}
	suite.mustSave(&aws)

	offer := models.ShipmentOffer{
		ShipmentID:                      aws.ID,
		TransportationServiceProviderID: tsp.ID,
	}
	suite.mustSave(&offer)

	params := shipmentop.NewIndexShipmentsParams()
	indexHandler := IndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	indexResponse := indexHandler.Handle(params)

	okResponse, ok := indexResponse.(*shipmentop.IndexShipmentsOK)
	if !ok {
		t.Fatalf("Request failed: %#v", indexResponse)
	}
	shipments := okResponse.Payload

	if len(shipments) != 2 {
		t.Errorf("expected %d shipments, got %d", 2, len(shipments))
	}

	offeredCount := 0
	availableCount := 0
	for _, shipment := range shipments {
		if shipment.TransportationServiceProviderID != nil {
			offeredCount++
			if shipment.TransportationServiceProviderID.String() != tsp.ID.String() {
				t.Errorf("got wrong tsp id, expected %s, got %s", tsp.ID.String(), shipment.TransportationServiceProviderID.String())

			}
		} else {
			availableCount++
		}
	}

	if offeredCount != 1 {
		t.Errorf("expected %d offered shipments, got %d", 1, offeredCount)
	}

	if availableCount != 1 {
		t.Errorf("expected %d available shipments, got %d", 1, availableCount)
	}
}
