package ghcrateengine

import (
	"log"
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
		log.Fatal("Error creating uuid")
	}
	pricingData := DomesticServicePricingData{
		MoveDate:      time.Date(2019, time.February, 18, 0, 0, 0, 0, time.UTC),
		ServiceAreaID: serviceAreaID,
		Distance:      456,
		Weight:        25.80,
		IsPeakPeriod:  false,
		ContractCode:  "ABC",
	}

	actual := engine.CalculateBaseDomesticLinehaul(pricingData)
	weight := pricingData.Weight
	var rate unit.Millicents = 272700
	expected := rate.MultiplyFloat64(float64(weight))
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
}
