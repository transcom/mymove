package invoice

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
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
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(db),
		logger:       logger,
	}
	suite.Run(t, hs)
}
