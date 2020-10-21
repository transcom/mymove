package invoice

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"
)

type GHCJobRunnerSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *GHCJobRunnerSuite) SetupTest() {
	errTruncateAll := suite.DB().TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

func TestGHCJobRunnerSuite(t *testing.T) {
	ts := &GHCJobRunnerSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("ghcjobrunner")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCJobRunnerSuite) TestGHCJobRunner() {
	jobRunner := NewGHCJobRunner(suite.DB())

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusPending,
		},
	})

	suite.T().Run("check for reviewed payment requests", func(t *testing.T) {
		result, err := jobRunner.ApprovedPaymentRequestFetcher()
		suite.NoError(err)
		suite.Equal(1, len(result))
	})

}
