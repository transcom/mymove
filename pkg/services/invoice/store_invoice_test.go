package invoice

import (
	"io/ioutil"
	"path/filepath"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func helperInvoice(suite *InvoiceServiceSuite) (*models.Invoice, *models.OfficeUser) {
	// Data setup
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	shipment := shipments[0]

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Invoice tied to a shipment of authed tsp user
	invoice := testdatagen.MakeInvoice(suite.DB(), testdatagen.Assertions{
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
	extantInvoice, err := models.FetchInvoice(suite.DB(), session, invoice.ID)
	suite.Nil(err)
	if suite.NoError(err) {
		suite.Equal(extantInvoice.ID, invoice.ID)
	}

	return &invoice, &officeUser
}

func helperExpectedEDIString(suite *InvoiceServiceSuite, name string) string {
	// TODO: Move this file to somewhere more central (or create one specific to this test).
	path := filepath.Join("..", "..", "edi", "invoice", "testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}

// TestStoreInvoice858C tests the store function EDI/Invoice to S3
func (suite *InvoiceServiceSuite) TestStoreInvoice858C() {
	invoice, officeUser := helperInvoice(suite)
	invoiceString := helperExpectedEDIString(suite, "expected_invoice.edi.golden")
	fs := suite.storer

	verrs, err := StoreInvoice858C{
		DB:     suite.DB(),
		Logger: suite.logger,
		Storer: &fs,
	}.Call(invoiceString, invoice, *officeUser.UserID)
	suite.Nil(err)
	suite.Empty(verrs.Error())

	// When: office user tries to access
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// Fetch invoice and verify upload is set
	invoice, err = models.FetchInvoice(suite.DB(), session, invoice.ID)
	suite.Nil(err)
	suite.NotNil(invoice)
	suite.NotNil(invoice.Upload)
	suite.NotNil(invoice.UploadID)
	// Check that StoragKey matches expected filepath name
	// {application-bucket}/app/invoice/{invoice_id}.edi
	suite.Regexp("^/app/invoice/([a-z0-9-])+\\.edi$", invoice.Upload.StorageKey)
}
