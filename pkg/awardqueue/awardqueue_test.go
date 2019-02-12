package awardqueue

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging/hnyzap"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *AwardQueueSuite) Test_CheckAllTSPsBlackedOut() {
	t := suite.T()
	queue := NewAwardQueue(suite.DB(), suite.logger)

	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)

	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, blackoutStartDate.Add(time.Hour))
	deliveryDate := dates.NextWorkday(*calendar, blackoutStartDate.Add(time.Hour*24*60))
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			ActualDeliveryDate:  &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			BookDate:            &testdatagen.DateInsidePerformancePeriod,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *shipment.TrafficDistributionList

	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	testdatagen.MakeBlackoutDate(suite.DB(), testdatagen.Assertions{
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
	offer, err := queue.attemptShipmentOffer(context.Background(), shipment)

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
	queue := NewAwardQueue(suite.DB(), suite.logger)

	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC

	calendar := dates.NewUSCalendar()
	blackoutStartDate := testdatagen.DateInsidePeakRateCycle
	blackoutEndDate := dates.NextWorkday(*calendar, blackoutStartDate.AddDate(0, 1, 0))

	blackoutPickupDate := dates.NextWorkday(*calendar, blackoutStartDate.AddDate(0, 0, 1))
	blackoutDeliverDate := dates.NextWorkday(*calendar, blackoutStartDate.AddDate(0, 0, 5))

	blackoutShipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &blackoutPickupDate,
			ActualPickupDate:    &blackoutPickupDate,
			ActualDeliveryDate:  &blackoutDeliverDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	pickupDate := dates.NextWorkday(*calendar, blackoutEndDate.AddDate(0, 0, 2))
	deliveryDate := dates.NextWorkday(*calendar, blackoutEndDate.AddDate(0, 0, 3))

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			ActualDeliveryDate:  &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *blackoutShipment.TrafficDistributionList

	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	testdatagen.MakeBlackoutDate(suite.DB(), testdatagen.Assertions{
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
	queue.assignShipments(context.Background())

	shipmentOffer := models.ShipmentOffer{}
	query := suite.DB().Where("shipment_id = $1", shipment.ID)
	if err := query.First(&shipmentOffer); err != nil {
		t.Errorf("Couldn't find shipment offer with shipment_ID: %v\n", shipment.ID)
	}

	blackoutShipmentOffer := models.ShipmentOffer{}
	blackoutQuery := suite.DB().Where("shipment_id = $1", blackoutShipment.ID)
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
	if err := suite.DB().Find(&blackoutShipment, blackoutShipment.ID); err != nil {
		t.Errorf("Couldn't find blackout shipment with ID: %v", blackoutShipment.ID)
	}
	if shipment.Status == models.ShipmentStatusAWARDED {
		t.Errorf("Blackout shipment was erroneously awarded for shipment ID: %v", shipment.ID)
	}

	if err := suite.DB().Find(&shipment, shipment.ID); err != nil {
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
	queue := NewAwardQueue(suite.DB(), suite.logger)
	// Creates a TSP with a blackout date connected to both.
	testTSP1 := testdatagen.MakeDefaultTSP(suite.DB())

	calendar := dates.NewUSCalendar()
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	testStartDate := dates.NextWorkday(*calendar, testdatagen.DateInsidePeakRateCycle)
	testEndDate := dates.NextWorkday(*calendar, testStartDate.Add(time.Hour*24*2))

	// Two pickup times to check with ShipmentWithinBlackoutDates
	testPickupDateBetween := dates.NextWorkday(*calendar, testStartDate.Add(time.Hour*24))
	testPickupDateAfter := dates.NextWorkday(*calendar, testEndDate.Add(time.Hour*24*5))

	// Two shipments using these pickup dates to provide to ShipmentWithinBlackoutDates
	testShipmentBetween := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &testPickupDateBetween,
			ActualPickupDate:    &testPickupDateBetween,
			ActualDeliveryDate:  &testEndDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			BookDate:            &testdatagen.DateInsidePerformancePeriod,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	testShipmentAfter := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &testPickupDateAfter,
			ActualPickupDate:    &testPickupDateAfter,
			ActualDeliveryDate:  &testEndDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			BookDate:            &testdatagen.DateInsidePerformancePeriod,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *testShipmentBetween.TrafficDistributionList

	testdatagen.MakeBlackoutDate(suite.DB(), testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: testTSP1.ID,
			StartBlackoutDate:               testStartDate,
			EndBlackoutDate:                 testEndDate,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &sourceGBLOC,
			Market:                          &market,
		},
	})

	// One TSP with no blackout dates
	testTSP2 := testdatagen.MakeDefaultTSP(suite.DB())

	// Checks a date that falls within the blackout date range; returns true.
	test1, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, testShipmentBetween)

	if err != nil {
		t.Fatal(err)
	} else if !test1 {
		t.Errorf("Expected true, got false instead.")
	}

	// Checks a date that falls after the blackout date range; returns false.
	test2, err := queue.ShipmentWithinBlackoutDates(testTSP1.ID, testShipmentAfter)

	if err != nil {
		t.Fatal(err)
	} else if test2 {
		t.Errorf("Expected false, got true instead.")
	}

	// Checks a TSP with no blackout dates and returns false.
	test3, err := queue.ShipmentWithinBlackoutDates(testTSP2.ID, testShipmentAfter)

	if err != nil {
		t.Fatal(err)
	} else if test3 {
		t.Errorf("Expected false, got true instead.")
	}
}

func (suite *AwardQueueSuite) Test_FindAllUnassignedShipments() {
	t := suite.T()
	queue := NewAwardQueue(suite.DB(), suite.logger)
	_, err := queue.findAllUnassignedShipments()

	if err != nil {
		t.Error("Unable to find shipments: ", err)
	}
}

// Test that we can create a shipment that should be offered, and that
// it actually gets offered.
func (suite *AwardQueueSuite) Test_OfferSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.DB(), suite.logger)

	// Make a shipment
	calendar := dates.NewUSCalendar()
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := dates.NextWorkday(*calendar, testdatagen.DateInsidePeakRateCycle)
	deliveryDate := pickupDate

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			ActualDeliveryDate:  &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			Market:              &market,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	tdl := *shipment.TrafficDistributionList

	// Make a TSP to handle it
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	tspp, err := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)
	suite.Nil(err)

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(context.Background(), shipment)

	// See if shipment was offered
	if err != nil {
		t.Errorf("Shipment offer expected no errors, received: %v", err)
	} else if offer == nil {
		t.Error("ShipmentOffer was not found.")
	} else {
		if err := suite.DB().Find(&shipment, shipment.ID); err != nil {
			t.Errorf("Couldn't find shipment with ID: %v", shipment.ID)
		}
		if shipment.Status != models.ShipmentStatusAWARDED {
			t.Errorf("Shipment should have been awarded for shipment ID %v but instead had status %v", shipment.ID,
				shipment.Status)
		}
	}

	suite.Equal(tsp.ID, offer.TransportationServiceProviderID)
	suite.Equal(tspp.ID, offer.TransportationServiceProviderPerformanceID)
}

// Test that a shipment does NOT get offered because it is not in a TDL with
// any enabled TSPs.
func (suite *AwardQueueSuite) Test_FailOfferingSingleShipment() {
	t := suite.T()
	queue := NewAwardQueue(suite.DB(), suite.logger)

	// Make a shipment in a new TDL, which inherently has no TSPs
	market := "dHHG"
	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			ActualDeliveryDate:  &deliveryDate,
			SourceGBLOC:         &sourceGBLOC,
			DestinationGBLOC:    &destinationGBLOC,
			Market:              &market,
			BookDate:            &pickupDate,
			Status:              models.ShipmentStatusSUBMITTED,
		},
	})

	var scac = "NPEK"
	var supplierID = scac + "1234" // scac + payee code
	// Make a TSP in the same TDL, but that is NOT enrolled
	tdl := *shipment.TrafficDistributionList
	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID, //NPEK1234
			Enrolled:                 false,
		},
	})
	_, err := testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)
	suite.Nil(err)

	// Run the Award Queue
	offer, err := queue.attemptShipmentOffer(context.Background(), shipment)

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
	queue := NewAwardQueue(suite.DB(), suite.logger)

	const shipmentsToMake = 10

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make a few shipments in this TDL
	var shipments [shipmentsToMake]models.Shipment
	for i := 0; i < shipmentsToMake; i++ {
		shipments[i] = testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate: &pickupDate,
				ActualPickupDate:    &pickupDate,
				ActualDeliveryDate:  &deliveryDate,
				SourceGBLOC:         &sourceGBLOC,
				Market:              &market,
				Status:              models.ShipmentStatusSUBMITTED,
			},
		})
	}

	tdl := *shipments[0].TrafficDistributionList

	// Make a TSP in the same TDL to handle these shipments
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	// ... and give this TSP a performance record
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, swag.Int(1), mps+1, 0, .3, .3)

	// Run the Award Queue
	queue.assignShipments(context.Background())

	// Count the number of shipments offered to our TSP
	query := suite.DB().Where("transportation_service_provider_id = $1", tsp.ID)
	offers := []models.ShipmentOffer{}
	count, err := query.Count(&offers)

	if err != nil {
		t.Errorf("Error counting shipment offers: %v", err)
	}
	if count != shipmentsToMake {
		t.Errorf("Not all ShipmentOffers found. Expected %d found %d", shipmentsToMake, count)
	}

	for _, shipment := range shipments {
		if err := suite.DB().Find(&shipment, shipment.ID); err != nil {
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

	suite.DB().TruncateAll()

	queue := NewAwardQueue(suite.DB(), suite.logger)

	const shipmentsToMake = 17

	// Shipment details
	market := testdatagen.DefaultMarket
	sourceGBLOC := testdatagen.DefaultSrcGBLOC
	pickupDate := testdatagen.DateInsidePeakRateCycle
	deliveryDate := testdatagen.DateInsidePeakRateCycle.Add(time.Hour)

	// Make shipments in this TDL
	var shipments [shipmentsToMake]models.Shipment
	for i := 0; i < shipmentsToMake; i++ {
		shipments[i] = testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate: &pickupDate,
				ActualPickupDate:    &pickupDate,
				ActualDeliveryDate:  &deliveryDate,
				SourceGBLOC:         &sourceGBLOC,
				Market:              &market,
				Status:              models.ShipmentStatusSUBMITTED,
			},
		})
	}

	tdl := *shipments[0].TrafficDistributionList

	// Make TSPs in the same TDL to handle these shipments
	tsp1 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp2 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp3 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp4 := testdatagen.MakeDefaultTSP(suite.DB())
	tsp5 := testdatagen.MakeDefaultTSP(suite.DB())

	// TSPs should be orderd by offer_count first, then BVS.
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp1, tdl, swag.Int(1), mps+5, 0, .4, .4)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp2, tdl, swag.Int(1), mps+4, 0, .3, .3)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp3, tdl, swag.Int(2), mps+2, 0, .2, .2)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp4, tdl, swag.Int(3), mps+3, 0, .1, .1)
	testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp5, tdl, swag.Int(4), mps+1, 0, .6, .6)

	// Run the Award Queue
	queue.assignShipments(context.Background())

	// TODO: revert to [6, 5, 3, 2, 1] after the B&M pilot
	suite.verifyOfferCount(tsp1, 4)
	suite.verifyOfferCount(tsp2, 4)
	suite.verifyOfferCount(tsp3, 3)
	suite.verifyOfferCount(tsp4, 3)
	suite.verifyOfferCount(tsp5, 3)

	for _, shipment := range shipments {
		if err := suite.DB().Find(&shipment, shipment.ID); err != nil {
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
	queue := NewAwardQueue(suite.DB(), suite.logger)
	tspsToMake := 5

	tdl := testdatagen.MakeDefaultTDL(suite.DB())

	var lastTSPP models.TransportationServiceProviderPerformance
	for i := 0; i < tspsToMake; i++ {
		tsp := testdatagen.MakeDefaultTSP(suite.DB())
		score := float64(mps + i + 1)

		rate := unit.NewDiscountRateFromPercent(45.3)
		var err error
		lastTSPP, err = testdatagen.MakeTSPPerformanceDeprecated(suite.DB(), tsp, tdl, nil, score, 0, rate, rate)
		if err != nil {
			t.Errorf("Failed to MakeTSPPerformance: %v", err)
		}
	}

	err := queue.assignPerformanceBands(context.Background())

	if err != nil {
		t.Errorf("Failed to assign to performance bands: %v", err)
	}

	perfGroup := models.TSPPerformanceGroup{
		TrafficDistributionListID: lastTSPP.TrafficDistributionListID,
		PerformancePeriodStart:    lastTSPP.PerformancePeriodStart,
		PerformancePeriodEnd:      lastTSPP.PerformancePeriodEnd,
		RateCycleStart:            lastTSPP.RateCycleStart,
		RateCycleEnd:              lastTSPP.RateCycleEnd,
	}

	perfs, err := models.FetchTSPPerformancesForQualityBandAssignment(suite.DB(), perfGroup, mps)
	if err != nil {
		t.Errorf("Failed to fetch TSPPerformances: %v", err)
	}

	expectedBands := []int{1, 1, 2, 3, 4}

	for i, perf := range perfs {
		band := expectedBands[i]
		if perf.QualityBand == nil {
			t.Errorf("No quality band assigned for Performance #%v, got nil", perf.ID)
		} else if (*perf.QualityBand) != band {
			t.Errorf("Wrong quality band: expected %v, got %v", band, *perf.QualityBand)
		}
	}
}

// Test_AwardTSPsInDifferentRateCycles ensures that TSPs that service different
// rate cycles get awarded shipments appropriately
func (suite *AwardQueueSuite) Test_AwardTSPsInDifferentRateCycles() {
	t := suite.T()
	queue := NewAwardQueue(suite.DB(), suite.logger)

	sm := testdatagen.MakeDefaultServiceMember(suite.DB())
	twoMonths, _ := time.ParseDuration("2 months")
	twoMonthsLater := testdatagen.PerformancePeriodStart.Add(twoMonths)

	// Make Peak TSP and Shipment
	shipmentPeak := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate:    &testdatagen.DateInsidePeakRateCycle,
			RequestedPickupDate: &testdatagen.DateInsidePeakRateCycle,
			ActualDeliveryDate:  &twoMonthsLater,
			BookDate:            &testdatagen.PerformancePeriodStart,
			Status:              models.ShipmentStatusSUBMITTED,
			ServiceMemberID:     sm.ID,
			ServiceMember:       sm,
		},
	})

	tspPeak := testdatagen.MakeDefaultTSP(suite.DB())
	tspPerfPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TransportationServiceProviderID: tspPeak.ID,
		TrafficDistributionListID:       shipmentPeak.TrafficDistributionList.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err := suite.DB().ValidateAndSave(&tspPerfPeak)
	if err != nil {
		t.Error(err)
	}

	// Make Non-Peak TSP and Shipment
	shipmentNonPeak := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate:    &testdatagen.DateInsideNonPeakRateCycle,
			RequestedPickupDate: &testdatagen.DateInsideNonPeakRateCycle,
			ActualDeliveryDate:  &twoMonthsLater,
			BookDate:            &testdatagen.PerformancePeriodStart,
			Status:              models.ShipmentStatusSUBMITTED,
			ServiceMemberID:     sm.ID,
			ServiceMember:       sm,
		},
	})

	var scac = "NPEK"
	var supplierID = scac + "1234" // scac + payee code
	tspNonPeak := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID, // NPEK1234
			Enrolled:                 true,
		},
	})
	tspPerfNonPeak := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.NonPeakRateCycleStart,
		RateCycleEnd:                    testdatagen.NonPeakRateCycleEnd,
		TransportationServiceProviderID: tspNonPeak.ID,
		TrafficDistributionListID:       shipmentNonPeak.TrafficDistributionList.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  100,
		OfferCount:                      0,
	}
	_, err = suite.DB().ValidateAndSave(&tspPerfNonPeak)
	if err != nil {
		t.Error(err)
	}

	queue.assignShipments(context.Background())

	suite.verifyOfferCount(tspPeak, 1)
	suite.verifyOfferCount(tspNonPeak, 1)

	// Test that shipments were awarded appropriately.
	if err := suite.DB().Find(&shipmentPeak, shipmentPeak.ID); err != nil {
		t.Errorf("Couldn't find peak shipment with ID: %v", shipmentPeak.ID)
	}
	if shipmentPeak.Status != models.ShipmentStatusAWARDED {
		t.Errorf("Shipment should have been awarded for peak shipment ID %v but instead had status %v",
			shipmentPeak.ID, shipmentPeak.Status)
	}

	if err := suite.DB().Find(&shipmentNonPeak, shipmentNonPeak.ID); err != nil {
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

	query := suite.DB().Where("transportation_service_provider_id = $1", tsp.ID)
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

func (suite *AwardQueueSuite) Test_validateShipmentForAward() {
	t := suite.T()
	t.Helper()

	// A default shipment, which has all valid fields on it
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{},
	})
	err := validateShipmentForAward(shipment)
	suite.Nil(err)

	// A shipment with a nil TDL ID
	shipment = testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{},
	})
	shipment.TrafficDistributionListID = nil
	err = validateShipmentForAward(shipment)
	suite.NotNil(err)

	// A shipment with a nil BookDate
	shipment = testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{},
	})
	shipment.BookDate = nil
	err = validateShipmentForAward(shipment)
	suite.NotNil(err)

	// A shipment with a nil RequestedPickupDate
	shipment = testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{},
	})
	shipment.RequestedPickupDate = nil
	err = validateShipmentForAward(shipment)
	suite.NotNil(err)
}

func (suite *AwardQueueSuite) Test_waitForLock() {
	ctx := context.Background()
	ret := make(chan int)
	lockID := 1

	go func() {
		suite.DB().Transaction(func(tx *pop.Connection) error {
			suite.Nil(waitForLock(ctx, tx, lockID))
			time.Sleep(time.Second)
			ret <- 1
			return nil
		})
	}()

	go func() {
		suite.DB().Transaction(func(tx *pop.Connection) error {
			time.Sleep(time.Millisecond * 500)
			suite.Nil(waitForLock(ctx, tx, lockID))
			ret <- 2
			return nil
		})
	}()

	first := <-ret
	second := <-ret

	suite.Equal(1, first)
	suite.Equal(2, second)
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
	testingsuite.PopTestSuite
	logger *hnyzap.Logger
}

func (suite *AwardQueueSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestAwardQueueSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &AwardQueueSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       &hnyzap.Logger{Logger: logger},
	}
	suite.Run(t, hs)
}
