package invoice

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/facebookgo/clock"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/unit"

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
	document := testdatagen.MakeDefaultDocument(suite.DB())
	userID := document.ServiceMember.UserID
	up := uploader.NewUploader(suite.DB(), suite.logger, *storer)

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

func (suite *InvoiceServiceSuite) helperShipment() models.Shipment {
	var weight unit.Pound
	weight = 2000
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			NetWeight: &weight,
		},
	})
	err := shipment.AssignGBLNumber(suite.DB())
	//TODO: mustSave(db, shipment)
	if err != nil {
		log.Fatalf("could not assign GBLNumber: %v", err)
	}

	// Create an accepted shipment offer and the associated TSP.
	scac := "ABBV"
	supplierID := scac + "2708" //scac + payee code -- ABBV2708

	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID,
		},
	})

	tspp := testdatagen.MakeTSPPerformance(suite.DB(), testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			TransportationServiceProvider:   tsp,
			TransportationServiceProviderID: tsp.ID,
		},
	})

	shipmentOffer := testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			Shipment:                                   shipment,
			Accepted:                                   swag.Bool(true),
			TransportationServiceProvider:              tsp,
			TransportationServiceProviderID:            tsp.ID,
			TransportationServiceProviderPerformance:   tspp,
			TransportationServiceProviderPerformanceID: tspp.ID,
		},
	})
	shipment.ShipmentOffers = models.ShipmentOffers{shipmentOffer}

	// Create some shipment line items.
	var lineItems []models.ShipmentLineItem
	codes := []string{"LHS", "135A", "135B", "105A", "16A", "105C", "125B", "105B", "130B", "46A"}
	amountCents := unit.Cents(12325)

	for _, code := range codes {
		appliedRate := unit.Millicents(2537234)
		var measurementUnit1 models.Tariff400ngItemMeasurementUnit
		var location models.ShipmentLineItemLocation

		switch code {
		case "LHS":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
			appliedRate = 0
		case "16A":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
		case "105B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitCUBICFOOT

		case "130B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitEACH

		case "125B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE

		default:
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitWEIGHT
		}

		// default location created in testdatagen shipmentLineItem is DESTINATION
		if code == "135A" || code == "105A" {
			location = models.ShipmentLineItemLocationORIGIN
		}
		if code == "135B" {
			location = models.ShipmentLineItemLocationDESTINATION
		}
		if code == "LHS" || code == "46A" {
			location = models.ShipmentLineItemLocationNEITHER
		}

		item := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
			Tariff400ngItem: models.Tariff400ngItem{
				Code:             code,
				MeasurementUnit1: measurementUnit1,
			},
		})
		lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment:          shipment,
				Tariff400ngItemID: item.ID,
				Tariff400ngItem:   item,
				Quantity1:         unit.BaseQuantityFromInt(2000),
				AppliedRate:       &appliedRate,
				AmountCents:       &amountCents,
				Location:          location,
			},
		})

		lineItems = append(lineItems, lineItem)
	}
	shipment.ShipmentLineItems = lineItems

	return shipment
}

func (suite *InvoiceServiceSuite) helperShipmentInvoice(shipment models.Shipment) *models.Invoice {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	var invoiceModel models.Invoice
	verrs, err := CreateInvoice{DB: suite.DB(), Clock: clock.NewMock()}.Call(officeUser, &invoiceModel, shipment)
	if err != nil {
		log.Fatalf("error when creating invoice: %v", err)
	}
	if verrs.HasAny() {
		log.Fatalf("validation errors when creating invoice: %s", verrs.String())
	}

	return &invoiceModel
}

func (suite *InvoiceServiceSuite) helperCreateFileStorer() *storage.FileStorer {
	var storer storage.FileStorer
	fakeS3 := storageTest.NewFakeS3Storage(true)
	storer = fakeS3
	return &storer
}

func (suite *InvoiceServiceSuite) helperFetchInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	err := suite.DB().Eager().Find(&invoice, invoiceID)
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
	shipment := suite.helperShipment()
	invoice := suite.helperShipmentInvoice(shipment)
	suite.NotNil(invoice)
	upload := suite.helperCreateUpload(storer)
	suite.NotNil(upload)

	up := uploader.NewUploader(suite.DB(), suite.logger, *storer)

	// Add upload to invoice
	verrs, err := UpdateInvoiceUpload{DB: suite.DB(), Uploader: up}.Call(invoice, upload)
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

	// Delete upload
	upload = suite.helperCreateUpload(storer)
	suite.NotNil(upload)
	err = UpdateInvoiceUpload{DB: suite.DB(), Uploader: up}.DeleteUpload(invoice)
	suite.Nil(err)
	suite.Empty(verrs.Error())
	suite.Nil(invoice.UploadID)
	suite.Nil(invoice.Upload)
}
