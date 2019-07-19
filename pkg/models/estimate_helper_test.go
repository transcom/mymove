package models_test

import (
	"time"

	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) setupPPMDiscountFetchRates() {
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "395",
		BasepointCity: "Saucier",
		State:         "MS",
		ServiceArea:   "428",
		RateArea:      "US48",
		Region:        "11",
	}
	suite.MustSave(&originZip3)
	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Gulfport, MS",
		ServiceArea:        "428",
		LinehaulFactor:     57,
		ServiceChargeCents: 350,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}
	suite.MustSave(&originServiceArea)
	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          "336",
		BasepointCity: "Tampa",
		State:         "FL",
		ServiceArea:   "197",
		RateArea:      "US4964400",
		Region:        "13",
	}
	suite.MustSave(&destinationZip3)
	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Tampa, FL",
		ServiceArea:        "197",
		LinehaulFactor:     69,
		ServiceChargeCents: 663,
		ServicesSchedule:   1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(5550),
		SIT185BRateCents:   unit.Cents(222),
		SITPDSchedule:      1,
	}
	suite.MustSave(&destinationServiceArea)
	newBaseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 10000,
		WeightLbsLower:     1000,
		WeightLbsUpper:     4000,
		RateCents:          20000,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&newBaseLinehaul)
}

func (suite *ModelSuite) TestPPMDiscountFetch() {
	suite.setupPPMDiscountFetchRates()
	logger, _ := zap.NewDevelopment()
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US48",
			DestinationRegion: "13",
			CodeOfService:     "2",
		},
	})
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50.0),
	}
	suite.MustSave(&tspPerformance)
	originZip := "39574"
	destinationZip := "33633"

	//
	// Successful move date
	//
	allowBookDate := false
	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(suite.DB(),
		logger,
		originZip,
		destinationZip,
		testdatagen.RateEngineDate,
		time.Time{},
		allowBookDate,
	)
	suite.Nil(err)
	suite.Equal(0.505, lhDiscount.Float64(), "Discount rate is not 0.505")
	suite.Equal(0.5, sitDiscount.Float64(), "SIT rate is not 0.5")

	//
	// Successful move date
	//
	allowBookDate = true
	lhDiscount, sitDiscount, err = models.PPMDiscountFetch(suite.DB(),
		logger,
		originZip,
		destinationZip,
		testdatagen.RateEngineDate,
		time.Time{},
		allowBookDate,
	)
	suite.Nil(err)
	suite.Equal(0.505, lhDiscount.Float64(), "Discount rate is not 0.505")
	suite.Equal(0.5, sitDiscount.Float64(), "SIT rate is not 0.5")

	//
	// Failed move date and not using book date
	//
	allowBookDate = false
	lhDiscount, sitDiscount, err = models.PPMDiscountFetch(suite.DB(),
		logger,
		originZip,
		destinationZip,
		testdatagen.RateEngineDate.AddDate(2, 0, 0),
		time.Time{},
		allowBookDate,
	)
	// Expect to get FETCH_NOT_FOUND
	suite.NotNil(err)
	suite.Equal("FETCH_NOT_FOUND", err.Error(), "Expect FETCH_NOT_FOUND for move date")
	suite.Equal(0.0, lhDiscount.Float64(), "Discount rate is 0.0")
	suite.Equal(0.0, sitDiscount.Float64(), "SIT rate is not 0.0")

	//
	// Failed move date and Successful book date
	//
	allowBookDate = true
	lhDiscount, sitDiscount, err = models.PPMDiscountFetch(suite.DB(),
		logger,
		originZip,
		destinationZip,
		testdatagen.RateEngineDate.AddDate(2, 0, 0),
		testdatagen.RateEngineDate,
		allowBookDate,
	)
	suite.Nil(err)
	suite.Equal(0.505, lhDiscount.Float64(), "Discount rate is not 0.505")
	suite.Equal(0.5, sitDiscount.Float64(), "SIT rate is not 0.5")
}
