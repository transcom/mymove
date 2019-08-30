package invoice

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type InvoiceServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
	storer storage.FileStorer
}

func (suite *InvoiceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestInvoiceSuite(t *testing.T) {

	hs := &InvoiceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("invoice_service")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
		storer:       storageTest.NewFakeS3Storage(true),
	}
	suite.Run(t, hs)
}
