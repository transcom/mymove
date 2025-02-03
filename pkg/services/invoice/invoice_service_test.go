package invoice

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type InvoiceServiceSuite struct {
	*testingsuite.PopTestSuite
	storer storage.FileStorer
}

func TestInvoiceSuite(t *testing.T) {

	ts := &InvoiceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("invoice_service"),
			testingsuite.WithPerTestTransaction()),
		storer: storageTest.NewFakeS3Storage(true),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
