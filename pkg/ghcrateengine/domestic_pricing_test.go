package ghcrateengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineSuite) Test_CalculateBaseDomesticLinehaul() {
	engine := NewGHCRateEngine(suite.DB(), suite.logger)

	serviceAreaID, err := uuid.FromString("9dda4dec-4dac-4aeb-b6ba-6736d689da8e")
	if err != nil {
		suite.logger.Fatal("Error creating uuid", zap.Error(err))
	}
	pricingData := DomesticServicePricingData{
		MoveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
		ServiceAreaID: serviceAreaID,
		Distance:      456,
		Weight:        25.80,
		IsPeakPeriod:  false,
		ContractCode:  "ABC",
	}

	actual, err := engine.CalculateBaseDomesticLinehaul(pricingData)
	if err != nil {
		suite.logger.Fatal("Could not calculate domestic linehaul", zap.Error(err))
	}
	weight := pricingData.Weight
	var rate unit.Millicents = 272700
	expected := rate.MultiplyFloat64(float64(weight))
	suite.Assertions.Equal(expected, actual)
}

func (suite *GHCRateEngineSuite) TestCalculateBaseDomesticPerWeightServiceCost() {
	engine := NewGHCRateEngine(suite.DB(), suite.logger)
	// test SIT day 1, SIT addt'l days, SIT50, pack, unpack, service area
	type testDataStruct struct {
			testName	  string
			isDomesticOther bool
			expectedRate  unit.Cents
			moveDate      time.Time
			serviceAreaID uuid.UUID
			distance      unit.Miles
			weight        unit.CWTFloat // record this here as 5.00 if actualWt less than minimum of 5.00 cwt (500lb)
			isPeakPeriod  bool
			contractCode  string
			serviceCode   string // may change to Service when model is available
	}
	testCases := [] testDataStruct {
		{
			testName: "test SIT P/D cost calculation",
			isDomesticOther: true,
			expectedRate:  23440,
			moveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()), // TODO replace with test UUID for Birmingham, AL (area 4)
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  false,
			contractCode:  "ABC",
			serviceCode:   "DSIT",
		},
		{
			testName: "test pack calculation",
			isDomesticOther: true,
			expectedRate:	7250,
			moveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()),
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  false,
			contractCode:  "ABC",
			serviceCode: "DPK",
		},
		{
			testName: "test unpack calculation",
			isDomesticOther: true,
			expectedRate:  597,
			moveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()),
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  false,
			contractCode:  "ABC",
			serviceCode: "DUPK",
		},
		{
			testName: "test origin service area calculation",
			isDomesticOther: false,
			expectedRate: 689,
			moveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()),
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  false,
			contractCode:  "ABC",
			serviceCode: "OSA",
		},
		{
			testName: "test destination service area calculation",
			isDomesticOther: false,
			expectedRate: 689,
			moveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()),
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  false,
			contractCode:  "ABC",
			serviceCode: "DSA",
		},
		{
			testName: "test peak destination service area calculation",
			isDomesticOther: false,
			expectedRate: 792,
			moveDate:      time.Date(2019, time.June, 18, 0, 0, 0, 0, time.UTC),
			serviceAreaID: uuid.Must(uuid.NewV4()),
			weight:        25.80,
			distance: 	   458,
			isPeakPeriod:  true,
			contractCode:  "ABC",
			serviceCode: "DSA",
		},
	}

	for _, test := range testCases {
		suite.T().Run(test.testName, func(t *testing.T) {
			pricingData := DomesticServicePricingData{
				MoveDate:      test.moveDate,
				ServiceAreaID: test.serviceAreaID,
				Weight:        test.weight,
				IsPeakPeriod:  test.isPeakPeriod,
				ContractCode:  test.contractCode,
				ServiceCode: test.serviceCode,
			}

			actual, err := engine.CalculateBaseDomesticPerWeightServiceCost(pricingData, test.isDomesticOther)
			if err != nil {
				suite.logger.Fatal(fmt.Sprintf("Could not calculate %s", pricingData.ServiceCode), zap.Error(err))
			}
			weight := pricingData.Weight
			expected := test.expectedRate.MultiplyFloat64(float64(weight))
			suite.Assertions.Equal(expected, actual)
		})
	}
}

func (suite *GHCRateEngineSuite) TestCalculateBaseDomesticShorthaulCost() {
	engine := NewGHCRateEngine(suite.DB(), suite.logger)

	serviceAreaID, err := uuid.FromString("9dda4dec-4dac-4aeb-b6ba-6736d689da8e")
	if err != nil {
		suite.logger.Fatal("Error creating uuid", zap.Error(err))
	}
	pricingData := DomesticServicePricingData{
		MoveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
		ServiceAreaID: serviceAreaID,
		Weight:        25.80,
		Distance: 	   458,
		IsPeakPeriod:  false,
		ContractCode:  "ABC",
		ServiceCode: "SH", // shorthaul service
	}

	actual, err := engine.CalculateBaseDomesticShorthaulCost(pricingData)
	if err != nil {
		suite.logger.Fatal(fmt.Sprintf("Could not calculate %s", pricingData.ServiceCode), zap.Error(err))
	}
	weight := pricingData.Weight
	rate := unit.Cents(689)
	expected := rate.MultiplyCWTFloat(weight).MultiplyMiles(458)
	suite.Assertions.Equal(expected, actual)
}

type GHCRateEngineSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *GHCRateEngineSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGHCRateEngineSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger, _ := zap.NewDevelopment()

	hs := &GHCRateEngineSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
