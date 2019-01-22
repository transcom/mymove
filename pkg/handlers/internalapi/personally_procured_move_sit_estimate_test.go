package internalapi

import (
	"net/http/httptest"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithDcos() {
	t := suite.T()

	// Given: a TDL, TSP and TSP performance with SITRate for relevant location and date
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US68",
			DestinationRegion: "5",
			CodeOfService:     "D",
		},
	}) // Victoria, TX to Salina, KS
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	suite.MustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "748", Region: "6"})
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.MustSave(&destServiceArea)

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     models.IntPointer(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance)

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:   4,
		OriginZip:       "77901",
		DestinationZip:  "67401",
		WeightEstimate:  3000,
	}
	// And: ShowPPMSitEstimateHandler is queried
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(3060)
	if *sitCost.Estimate != expectedSitCost {
		t.Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}

func (suite *HandlerSuite) TestShowPPMSitEstimateHandler2cos() {
	t := suite.T()

	// Given: a TDL, TSP and TSP performance with SITRate for relevant location and date
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US68",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	}) // Victoria, TX to Salina, KS
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	suite.MustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "748", Region: "6"})
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.MustSave(&destServiceArea)

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     models.IntPointer(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance)

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:   4,
		OriginZip:       "77901",
		DestinationZip:  "67401",
		WeightEstimate:  3000,
	}
	// And: ShowPPMSitEstimateHandler is queried
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(3060)
	if *sitCost.Estimate != expectedSitCost {
		t.Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithError() {

	// Given: A PPM Estimate request with all relevant records except TSP performance
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "748", Region: "6"})
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.MustSave(&destServiceArea)

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains required auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:   4,
		OriginZip:       "77901",
		DestinationZip:  "67401",
		WeightEstimate:  3000,
	}
	// And: ShowPPMSitEstimateHandler is queried
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect bad request response
	suite.CheckResponseNotFound(showResponse)
}
