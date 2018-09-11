package awardqueue

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *AwardQueueSuite) Test_CheckAllTSPsBlackedOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	tsp := testdatagen.MakeDefaultTSP(suite.db)

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)

	pickupDate := blackoutStartDate.Add(time.Hour)
	deliveryDate := blackoutStartDate.Add(time.Hour * 24 * 60)
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			PickupDate:          &pickupDate,
			DeliveryDate:        &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *shipment.TrafficDistributionList

	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
		},
	})

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		RequestedPickupDate:             &pickupDate,
		PickupDate:                      &pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	expectedError := "could not find a TSP without blackout dates"
	// See if shipment was offered
	if err == nil || offer != nil {
		t.Errorf("Shipment was offered to a blacked out TSP!")
	} else if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Did not receive proper error message. Expected '%s', got '%s' instead.", expectedError, err)
	}
}

func (suite *AwardQueueSuite) Test_CheckShipmentDuringBlackOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	tsp := testdatagen.MakeDefaultTSP(suite.db)

	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.AddDate(0, 1, 0)

	blackoutPickupDate := blackoutStartDate.AddDate(0, 0, 1)
	blackoutDeliverDate := blackoutStartDate.AddDate(0, 0, 5)

	blackoutShipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &blackoutPickupDate,
			PickupDate:          &blackoutPickupDate,
			DeliveryDate:        &blackoutDeliverDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	pickupDate := blackoutEndDate.AddDate(0, 0, 1)
	deliveryDate := blackoutEndDate.AddDate(0, 0, 2)

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			PickupDate:          &pickupDate,
			DeliveryDate:        &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *blackoutShipment.TrafficDistributionList

	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &sourceGBLOC,
			Market:                          &market,
		},
	})

	// Run the Award Queue
	queue.assignShipments()

	shipmentOffer := models.ShipmentOffer{}
	query := suite.db.Where("shipment_id = $1", shipment.ID)
	if err := query.First(&shipmentOffer); err != nil {
		t.Errorf("Couldn't find shipment offer with shipment_ID: %v\n", shipment.ID)
	}

	blackoutShipmentOffer := models.ShipmentOffer{}
	blackoutQuery := suite.db.Where("shipment_id = $1", blackoutShipment.ID)
	if err := blackoutQuery.First(&blackoutShipmentOffer); err != nil {
		t.Errorf("Couldn't find shipment offer: %v", blackoutShipment.ID)
	}

	if shipmentOffer.AdministrativeShipment != false {
		t.Errorf("Shipment Awards erroneously assigned administrative status.")
	}

	if blackoutShipmentOffer.AdministrativeShipment != true {
		t.Errorf("Shipment Awards erroneously not assigned administrative status.")
	}

	// Test that shipments were awarded appropriately.
	if err := suite.db.Find(&blackoutShipment, blackoutShipment.ID); err != nil {
		t.Errorf("Couldn't find blackout shipment with ID: %v", blackoutShipment.ID)
	}
	if shipment.Status == models.ShipmentStatusAWARDED {
		t.Errorf("Blackout shipment was erroneously awarded for shipment ID: %v", shipment.ID)
	}

	if err := suite.db.Find(&shipment, shipment.ID); err != nil {
		t.Errorf("Couldn't find shipment with ID: %v", shipment.ID)
	}
	if shipment.Status != models.ShipmentStatusAWARDED {
		t.Errorf("Shipment should have been awarded for shipment ID %v but instead had status %v", shipment.ID,
			shipment.Status)
	}

	suite.verifyOfferCount(tsp, 2)
}

func (suite *AwardQueueSuite) Test_ShipmentWithinBlackoutDates() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	// Creates a TSP with a blackout date connected to both.
	testTSP1 := testdatagen.MakeDefaultTSP(suite.db)

	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	testStartDate := testdatagen.DateInsidePeakRateCycle
	testEndDate := testStartDate.Add(time.Hour * 24 * 2)

	// Two pickup times to check with ShipmentWithinBlackoutDates
	testPickupDateBetween := testStartDate.Add(time.Hour * 24)
	testPickupDateAfter := testEndDate.Add(time.Hour * 24 * 5)

	// Two shipments using these pickup dates to provide to ShipmentWithinBlackoutDates
	testShipmentBetween := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &testPickupDateBetween,
			PickupDate:          &testStartDate,
			DeliveryDate:        &testEndDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	testShipmentAfter := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &testPickupDateAfter,
			PickupDate:          &testStartDate,
			DeliveryDate:        &testEndDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *testShipmentBetween.TrafficDistributionList

	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: testTSP1.ID,
			StartBlackoutDate:               testStartDate,
			EndBlackoutDate:                 testEndDate,
			TrafficDistributionListID:       &tdl.ID,
		},
	})

	// One TSP with no blackout dates
	testTSP2 := testdatagen.MakeDefaultTSP(suite.db)

	// Two shipments with offers, using the shipments above
	shipmentWithOfferBetween := models.ShipmentWithOffer{
		ID: testShipmentBetween.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &testPickupDateBetween,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
	}

	shipmentWithOfferAfter := models.ShipmentWithOffer{
		ID: testShipmentAfter.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &testPickupDateAfter,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        &testdatagen.DateInsidePerformancePeriod,
	}

	// Checks a date that falls within the blackout date range; returns true.
	test1, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, shipmentWithOfferBetween)

	if err != nil {
		t.Fatal(err)
	} else if !test1 {
		t.Errorf("Expected true, got false instead.")
	}

	// Checks a date that falls after the blackout date range; returns false.
	test2, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, shipmentWithOfferAfter)

	if err != nil {
		t.Fatal(err)
	} else if test2 {
		t.Errorf("Expected false, got true instead.")
	}

	// Checks a TSP with no blackout dates and returns false.
	test3, err := queue.ShipmentWithinBlackoutDates(testTSP2.ID, shipmentWithOfferAfter)

	if err != nil {
		t.Fatal(err)
	} else if test3 {
		t.Errorf("Expected false, got true instead.")
	}
}

func (suite *AwardQueueSuite) Test_FindAllUnassignedShipments() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	_, err := queue.findAllUnassignedShipments()

	if err != nil {
		t.Error("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be offered, and that
// it actually gets offered.
func (suite *AwardQueueSuite) Test_OfferSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	// Make a shipment
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			PickupDate:          &pickupDate,
			DeliveryDate:        &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *shipment.TrafficDistributionList

	// Make a TSP to handle it
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &pickupDate,
		RequestedPickupDate:             &pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        shipment.BookDate,
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	// See if shipment was offered
	if err != nil {
		t.Errorf("Shipment offer expected no errors, received: %v", err)
	} else if offer == nil {
		t.Error("ShipmentOffer was not found.")
	} else {
		if err := suite.db.Find(&shipment, shipment.ID); err != nil {
			t.Errorf("Couldn't find shipment with ID: %v", shipment.ID)
		}
		if shipment.Status != models.ShipmentStatusAWARDED {
			t.Errorf("Shipment should have been awarded for shipment ID %v but instead had status %v", shipment.ID,
				shipment.Status)
		}
	}
}

// Test that we can create a shipment that should NOT be offered because it is not in a TDL
// with any TSPs, and that it doesn't get offered.
func (suite *AwardQueueSuite) Test_FailOfferingSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	// Make a shipment in a new TDL, which inherently has no TSPs
	market := "dHHG"
	sourceGBLOC := "OHAI"
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			PickupDate:          &pickupDate,
			DeliveryDate:        &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *shipment.TrafficDistributionList

	// Create a ShipmentWithOffer to feed the award queue
	shipmentWithOffer := models.ShipmentWithOffer{
		ID: shipment.ID,
		TrafficDistributionListID:       &tdl.ID,
		PickupDate:                      &pickupDate,
		RequestedPickupDate:             &pickupDate,
		BookDate:                        &pickupDate,
		TransportationServiceProviderID: nil,
		AdministrativeShipment:          swag.Bool(false),
	}

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(shipmentWithOffer)

	// See if shipment was offered
	if err == nil {
		t.Errorf("Shipment offer expected an error, received none.")
	}
	if offer != nil {
		t.Errorf("Wrong return value, expected nil, got %v", offer)
	}
}

func (suite *AwardQueueSuite) TestAssignShipmentsSingleTSP() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)

	const shipmentsToMake = 10

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make a few shipments in this TDL
	var shipments [shipmentsToMake]models.Shipment
	for i := 0; i < shipmentsToMake; i++ {
		shipments[i] = testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate: &pickupDate,
				PickupDate:          &pickupDate,
				DeliveryDate:        &deliveryDate,
				SourceGBLOC:         &sourceGBLOC,
				Market:              &market,
				Status:              models.ShipmentStatusSUBMITTED,
			},
		})
	}

	tdl := *shipments[0].TrafficDistributionList

	// Make a TSP in the same TDL to handle these shipments
	tsp := testdatagen.MakeDefaultTSP(suite.db)

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	// Run the Award Queue
	queue.assignShipments()

	// Count the number of shipments offered to our TSP
	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	offers := []models.ShipmentOffer{}
	count, err := query.Count(&offers)

	if err != nil {
		t.Errorf("Error counting shipment offers: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentOffers found. Expected %d found %d", shipmentsToMake, count)
	}

	for _, shipment := range shipments {
		if err := suite.db.Find(&shipment, shipment.ID); err != nil {
			t.Errorf("Couldn't find shipment with ID: %v", shipment.ID)
		}
		if shipment.Status != models.ShipmentStatusAWARDED {
			t.Errorf("Shipment should have been awarded for shipment ID %v but instead had status %v", shipment.ID,
				shipment.Status)
		}
	}
}

func (suite *AwardQueueSuite) TestAssignShipmentsToMultipleTSPs() {
	t := suite.T()

	suite.db.TruncateAll()

	queue := NewAwardQueue(suite.db, suite.logger)

	const shipmentsToMake = 17

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make shipments in this TDL
	var shipments [shipmentsToMake]models.Shipment
	for i := 0; i < shipmentsToMake; i++ {
		shipments[i] = testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate: &pickupDate,
				PickupDate:          &pickupDate,
				DeliveryDate:        &deliveryDate,
				SourceGBLOC:         &sourceGBLOC,
				Market:              &market,
				Status:              models.ShipmentStatusSUBMITTED,
			},
		})
	}

	tdl := *shipments[0].TrafficDistributionList

	// Make TSPs in the same TDL to handle these shipments
	tsp1 := testdatagen.MakeDefaultTSP(suite.db)
	tsp2 := testdatagen.MakeDefaultTSP(suite.db)
	tsp3 := testdatagen.MakeDefaultTSP(suite.db)
	tsp4 := testdatagen.MakeDefaultTSP(suite.db)
	tsp5 := testdatagen.MakeDefaultTSP(suite.db)

	// TSPs should be orderd by offer_count first, then BVS.
	testdatagen.MakeTSPPerformance(suite.db, tsp1, tdl, swag.Int(1), mps+5, 0, .4, .4)
	testdatagen.MakeTSPPerformance(suite.db, tsp2, tdl, swag.Int(1), mps+4, 0, .3, .3)
	testdatagen.MakeTSPPerformance(suite.db, tsp3, tdl, swag.Int(2), mps+2, 0, .2, .2)
	testdatagen.MakeTSPPerformance(suite.db, tsp4, tdl, swag.Int(3), mps+3, 0, .1, .1)
	testdatagen.MakeTSPPerformance(suite.db, tsp5, tdl, swag.Int(4), mps+1, 0, .6, .6)

	// Run the Award Queue
	queue.assignShipments()

	suite.verifyOfferCount(tsp1, 6)
	suite.verifyOfferCount(tsp2, 5)
	suite.verifyOfferCount(tsp3, 3)
	suite.verifyOfferCount(tsp4, 2)
	suite.verifyOfferCount(tsp5, 1)

	for _, shipment := range shipments {
		if err := suite.db.Find(&shipment, shipment.ID); err != nil {
			t.Errorf("Couldn't find shipment with ID: %v", shipment.ID)
		}
		if shipment.Status != models.ShipmentStatusAWARDED {
			t.Errorf("Shipment should have been awarded for shipment ID %v but instead had status %v", shipment.ID,
				shipment.Status)
		}
	}
}

func (suite *AwardQueueSuite) Test_GetTSPsPerBandWithRemainder() {
	t := suite.T()
	// Check bands should expect differing num of TSPs when not divisible by 4
	// Remaining TSPs should be divided among bands in descending order
	tspPerBandList := getTSPsPerBand(10)
	expectedBandList := []int{3, 3, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_GetTSPsPerBandNoRemainder() {
	t := suite.T()
	// Check bands should expect correct num of TSPs when num of TSPs is divisible by 4
	tspPerBandList := getTSPsPerBand(8)
	expectedBandList := []int{2, 2, 2, 2}
	if !equalSlice(tspPerBandList, expectedBandList) {
		t.Errorf("Failed to correctly divide TSP counts. Expected to find %d, found %d", expectedBandList, tspPerBandList)
	}
}

func (suite *AwardQueueSuite) Test_AssignTSPsToBands() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	tspsToMake := 5

	tdl := testdatagen.MakeDefaultTDL(suite.db)

	for i := 0; i < tspsToMake; i++ {
		tsp := testdatagen.MakeDefaultTSP(suite.db)
		score := float64(mps + i + 1)

		rate := unit.NewDiscountRateFromPercent(45.3)
		testdatagen.MakeTSPPerformance(suite.db, tsp, tdl, nil, score, 0, rate, rate)
	}

	err := queue.assignPerformanceBands()

	if err != nil {
		t.Errorf("Failed to assign to performance bands: %v", err)
	}

	perfs, err := models.FetchTSPPerformanceForQualityBandAssignment(suite.db, tdl.ID, mps)
	if err != nil {
		t.Errorf("Failed to fetch TSPPerformances: %v", err)
	}

	expectedBands := []int{1, 1, 2, 3, 4}

	for i, perf := range perfs {
		band := expectedBands[i]
		if perf.QualityBand == nil {
			t.Errorf("No quality band assigned for Peformance #%v, got nil", perf.ID)
		} else if (*perf.QualityBand) != band {
			t.Errorf("Wrong quality band: expected %v, got %v", band, *perf.QualityBand)
		}
	}
}

// Test_AwardTSPsInDifferentRateCycles ensures that TSPs that service different
// rate cycles get awarded shipments appropriately
func (suite *AwardQueueSuite) Test_AwardTSPsInDifferentRateCycles() {
	t := suite.T()
	queue := NewAwardQueue(suite.db, suite.logger)
	sm := testdatagen.MakeDefaultServiceMember(suite.db)

	twoMonths, _ := time.ParseDuration("2 months")
	twoMonthsLater := testdatagen.PerformancePeriodStart.Add(twoMonths)

	tdl := testdatagen.MakeDefaultTDL(suite.db)

	// Make Peak TSP and Shipment
	tspPeak := testdatagen.MakeDefaultTSP(suite.db)

	tspPerfPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TransportationServiceProviderID: tspPeak.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err := suite.db.ValidateAndSave(&tspPerfPeak)
	if err != nil {
		t.Error(err)
	}

	shipmentPeak := models.Shipment{
		TrafficDistributionListID: &tdl.ID,
		PickupDate:                &testdatagen.DateInsidePeakRateCycle,
		RequestedPickupDate:       &testdatagen.DateInsidePeakRateCycle,
		DeliveryDate:              &twoMonthsLater,
		BookDate:                  &testdatagen.PerformancePeriodStart,
		SourceGBLOC:               &testdatagen.DefaultSrcGBLOC,
		Market:                    &testdatagen.DefaultMarket,
		MoveID:                    testdatagen.MakeDefaultMove(suite.db).ID,
		Status:                    models.ShipmentStatusSUBMITTED,
		ServiceMemberID:           sm.ID,
	}
	_, err = suite.db.ValidateAndSave(&shipmentPeak)
	if err != nil {
		t.Error(err)
	}

	// Make Non-Peak TSP and Shipment
	tspNonPeak := testdatagen.MakeTSP(suite.db, testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: "NPEK",
		},
	})
	tspPerfNonPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.NonPeakRateCycleStart,
		RateCycleEnd:                    testdatagen.NonPeakRateCycleEnd,
		TransportationServiceProviderID: tspNonPeak.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err = suite.db.ValidateAndSave(&tspPerfNonPeak)
	if err != nil {
		t.Error(err)
	}

	shipmentNonPeak := models.Shipment{
		TrafficDistributionListID: &tdl.ID,
		PickupDate:                &testdatagen.DateInsideNonPeakRateCycle,
		RequestedPickupDate:       &testdatagen.DateInsideNonPeakRateCycle,
		DeliveryDate:              &twoMonthsLater,
		BookDate:                  &testdatagen.PerformancePeriodStart,
		SourceGBLOC:               &testdatagen.DefaultSrcGBLOC,
		Market:                    &testdatagen.DefaultMarket,
		MoveID:                    testdatagen.MakeDefaultMove(suite.db).ID,
		Status:                    models.ShipmentStatusSUBMITTED,
		ServiceMemberID:           sm.ID,
	}
	_, err = suite.db.ValidateAndSave(&shipmentNonPeak)
	if err != nil {
		t.Error(err)
	}

	queue.assignShipments()

	suite.verifyOfferCount(tspPeak, 1)
	suite.verifyOfferCount(tspNonPeak, 1)

	// Test that shipments were awarded appropriately.
	if err := suite.db.Find(&shipmentPeak, shipmentPeak.ID); err != nil {
		t.Errorf("Couldn't find peak shipment with ID: %v", shipmentPeak.ID)
	}
	if shipmentPeak.Status != models.ShipmentStatusAWARDED {
		t.Errorf("Shipment should have been awarded for peak shipment ID %v but instead had status %v",
			shipmentPeak.ID, shipmentPeak.Status)
	}

	if err := suite.db.Find(&shipmentNonPeak, shipmentNonPeak.ID); err != nil {
		t.Errorf("Couldn't find non-peak shipment with ID: %v", shipmentNonPeak.ID)
	}
	if shipmentNonPeak.Status != models.ShipmentStatusAWARDED {
		t.Errorf("Shipment should have been awarded for non-peak shipment ID %v but instead had status %v",
			shipmentNonPeak.ID, shipmentNonPeak.Status)
	}
}

func (suite *AwardQueueSuite) verifyOfferCount(tsp models.TransportationServiceProvider, expectedCount int) {
	t := suite.T()
	t.Helper()

	query := suite.db.Where("transportation_service_provider_id = $1", tsp.ID)
	offers := []models.ShipmentOffer{}
	count, err := query.Count(&offers)

	if err != nil {
		t.Fatalf("Error counting shipment offers: %v", err)
	}
	if count != expectedCount {
		t.Errorf("Wrong number of ShipmentOffers found: expected %d, got %d", expectedCount, count)
	}

	var tspPerformance models.TransportationServiceProviderPerformance
	if err := query.First(&tspPerformance); err != nil {
		t.Errorf("No TSP Performance record found with id %s", tsp.ID)
	}
	if expectedCount != tspPerformance.OfferCount {
		t.Errorf("Wrong OfferCount for TSP: expected %d, got %d", expectedCount, tspPerformance.OfferCount)
	}
}

func equalSlice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type AwardQueueSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *AwardQueueSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestAwardQueueSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &AwardQueueSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
