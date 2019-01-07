package invoice

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
	"go.uber.org/zap"
)

func (suite *InvoiceServiceSuite) openLocalFile(path string) (afero.File, error) {
	var fs = afero.NewMemMapFs()

	file, err := os.Open(path)
	if err != nil {
		suite.logger.Fatal("Error opening local file", zap.Error(err))
	}

	suite.NotNil(fs)
	outputFile, err := fs.Create(path)
	if err != nil {
		suite.logger.Fatal("Error creating afero file", zap.Error(err))
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		suite.logger.Fatal("Error copying to afero file", zap.Error(err))
	}

	return outputFile, nil
}

// fixture creates a File for testing. Caller responsible to close file
// when done using it.
func (suite *InvoiceServiceSuite) fixture(name string) afero.File {
	fixtureDir := "testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatalf("failed to get current directory: %s", err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)
	file, err := suite.openLocalFile(fixturePath)

	if err != nil {
		suite.T().Fatalf("failed to create a fixture file: %s", err)
	}
	// Caller should call close on file when finished
	return file
}

func (suite *InvoiceServiceSuite) helperCreateUpload(storer *storage.FileStorer) *models.Upload {
	document := testdatagen.MakeDefaultDocument(suite.db)
	userID := document.ServiceMember.UserID
	up := uploader.NewUploader(suite.db, suite.logger, *storer)

	// Create file to use for upload
	file := suite.fixture("test.pdf")
	if file == nil {
		suite.T().Fatal("test.pdf is missing")
	}
	_, err := file.Stat()
	if err != nil {
		suite.T().Fatalf("file.Stat() err: %s", err.Error())
	}

	// Create Upload and save it
	upload, verrs, err := up.CreateUploadNoDocument(userID, &file)
	suite.Nil(err, "CreateUploadNoDocument() failed to create upload")
	suite.Empty(verrs.Error(), "CreateUploadNoDocument() verrs returned error")
	suite.NotNil(upload, "CreateUploadNoDocument() failed to create upload structure")
	if upload == nil {
		suite.T().Fatalf("failed to create a upload object: %s", err)
	}
	// Call Close on file after CreateUpload is complete
	file.Close()
	return upload
}

func (suite *InvoiceServiceSuite) helperCreateInvoice() *models.Invoice {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)
	shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)

	createInvoice := CreateInvoice{
		suite.db,
		clock.NewMock(),
	}
	var invoice models.Invoice
	verrs, err := createInvoice.Call(officeUser, &invoice, shipmentLineItem.Shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)
	updateInvoicesSubmitted := UpdateInvoiceSubmitted{
		DB: suite.db,
	}
	shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

	verrs, err = updateInvoicesSubmitted.Call(&invoice, shipmentLineItems)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(models.InvoiceStatusSUBMITTED, invoice.Status)
	suite.Equal(invoice.ID, *shipmentLineItems[0].InvoiceID)

	return &invoice
}

func (suite *InvoiceServiceSuite) helperCreateFileStorer() *storage.FileStorer {
	var storer storage.FileStorer
	fakeS3 := storageTest.NewFakeS3Storage(true)
	storer = fakeS3
	return &storer
}

func (suite *InvoiceServiceSuite) helperFetchInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	err := suite.db.Eager().Find(&invoice, invoiceID)
	if err != nil {
		fmt.Print(err.Error())
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return nil, errors.New("Record not found")
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	return &invoice, nil
}

// TestUpdateInvoiceUploadCall Test the Service UpdateInvoiceUpload{}.Call() function
func (suite *InvoiceServiceSuite) TestUpdateInvoiceUploadCall() {
	storer := suite.helperCreateFileStorer()
	invoice := suite.helperCreateInvoice()
	suite.NotNil(invoice)
	upload := suite.helperCreateUpload(storer)
	suite.NotNil(upload)

	up := uploader.NewUploader(suite.db, suite.logger, *storer)

	// Add first upload to invoice
	verrs, err := UpdateInvoiceUpload{DB: suite.db, Uploader: up}.Call(invoice, upload)
	suite.Nil(err)
	suite.Empty(verrs.Error())
	suite.Equal(upload.ID, *invoice.UploadID)

	// Add second upload to invoice -- will force delete of previous upload
	upload = suite.helperCreateUpload(storer)
	suite.NotNil(upload)
	verrs, err = UpdateInvoiceUpload{DB: suite.db, Uploader: up}.Call(invoice, upload)
	suite.Nil(err)
	suite.Empty(verrs.Error())
	suite.Equal(upload.ID, *invoice.UploadID)

	// Fetch Invoice from database and compare Upload IDs
	fetchInvoice, err := suite.helperFetchInvoice(invoice.ID)
	suite.Nil(err)
	suite.NotNil(fetchInvoice)
	suite.NotNil(fetchInvoice.UploadID)
	suite.NotNil(fetchInvoice.Upload)
	suite.Equal(upload.ID, *(fetchInvoice).UploadID)
}
