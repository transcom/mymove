package invoice

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type InvoiceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *InvoiceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
