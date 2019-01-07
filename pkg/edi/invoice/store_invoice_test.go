package ediinvoice_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"go.uber.org/zap"
)

type StoreInvoiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
	storer storage.FileStorer
}

// SetupTest
func (suite *StoreInvoiceSuite) SetupTest() {
	suite.db.TruncateAll()
}

// TestStoreInvoiceSuite
func TestStoreInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	fakeS3 := storageTest.NewFakeS3Storage(true)

	hs := &StoreInvoiceSuite{db: db, logger: logger, storer: fakeS3}
	suite.Run(t, hs)
}

// helperFileStorer is a simple setter for storage private field
func (suite *StoreInvoiceSuite) helperFileStorer() storage.FileStorer {
	return suite.storer
}

func helperInvoice(suite *StoreInvoiceSuite) (*models.Invoice, *models.OfficeUser) {
	// Data setup
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	// Invoice tied to a shipment of authed tsp user
	invoice := testdatagen.MakeInvoice(suite.db, testdatagen.Assertions{
		Invoice: models.Invoice{
			ShipmentID: shipment.ID,
			Approver:   officeUser,
			ApproverID: officeUser.ID,
		},
	})

	// When: office user tries to access
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// Then: invoice is returned
	extantInvoice, err := models.FetchInvoice(suite.db, session, invoice.ID)
	suite.Nil(err)
	if suite.NoError(err) {
		suite.Equal(extantInvoice.ID, invoice.ID)
	}

	return &invoice, &officeUser
}

func helperExpectedEDIString(suite *StoreInvoiceSuite, name string) string {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}

// TestStoreInvoice858C tests the store function EDI/Invoice to S3
func (suite *StoreInvoiceSuite) TestStoreInvoice858C() {
	invoice, officeUser := helperInvoice(suite)
	invoiceString := helperExpectedEDIString(suite, "expected_invoice.edi.golden")
	fs := suite.helperFileStorer()
	verrs, err := ediinvoice.StoreInvoice858C(invoiceString, invoice, &fs, suite.logger, *officeUser.UserID, suite.db)
	suite.Nil(err)
	suite.Empty(verrs.Error())
}
