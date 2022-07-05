package fuelprice

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FuelPriceServiceSuite struct {
	*testingsuite.PopTestSuite
	storer storage.FileStorer
}

func TestFuelPriceSuite(t *testing.T) {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	ts := &FuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		storer:       fakeS3,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
