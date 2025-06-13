package paperwork

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkServiceSuite) TestPrimeDownloadMoveUploadPDFGenerator() {
	service, order := suite.setupOrdersDocument()

	pdfGenerator, err := paperwork.NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)

	locator := "AAAA"

	customMoveWithOnlyOrders := models.Move{
		Locator: locator,
		Orders:  order,
	}

	pdfFileTest1, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderUploadAll, customMoveWithOnlyOrders, "")
	suite.FatalNil(err)
	// Verify generated files have 3 pages. see setup data for upload count
	fileInfo, err := suite.pdfFileInfo(pdfGenerator, pdfFileTest1)
	suite.FatalNil(err)
	suite.Equal(3, fileInfo.PageCount)

	// Point amendments doc to UploadedOrdersID.
	order.UploadedAmendedOrdersID = &order.UploadedOrdersID
	customMoveWithOrdersAndAmendments := models.Move{
		Locator: locator,
		Orders:  order,
	}
	pdfFileTest2, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderUploadAll, customMoveWithOrdersAndAmendments, "")
	suite.FatalNil(err)
	// Verify generated files have (3 x 2) pages for both orders and amendments. see setup data for upload count
	fileInfoAll, err := suite.pdfFileInfo(pdfGenerator, pdfFileTest2)
	suite.FatalNil(err)
	suite.Equal(6, fileInfoAll.PageCount)

	pdfFileTest3, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderUpload, customMoveWithOrdersAndAmendments, "")
	suite.FatalNil(err)
	// Verify generated files have (3 x 1) pages for order. see setup data for upload count
	fileInfoAll1, err := suite.pdfFileInfo(pdfGenerator, pdfFileTest3)
	suite.FatalNil(err)
	suite.Equal(3, fileInfoAll1.PageCount)

	pdfFileTest4, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderAmendmentUpload, customMoveWithOrdersAndAmendments, "")
	suite.FatalNil(err)
	// Verify only amendments are generated
	fileInfoOnlyAmendments, err := suite.pdfFileInfo(pdfGenerator, pdfFileTest4)
	suite.FatalNil(err)
	suite.Equal(3, fileInfoOnlyAmendments.PageCount)
	// cleanup the created files
	err = service.CleanupFile(pdfFileTest1)
	suite.NoError(err)
	err = service.CleanupFile(pdfFileTest2)
	suite.NoError(err)
	err = service.CleanupFile(pdfFileTest3)
	suite.NoError(err)
	err = service.CleanupFile(pdfFileTest4)
	suite.NoError(err)
	suite.AfterTest()
}

func (suite *PaperworkServiceSuite) TestPrimeDownloadMoveUploadPDFGeneratorUnprocessableEntityError() {
	pdfGenerator, err := paperwork.NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)
	service, _ := NewMoveUserUploadToPDFDownloader(pdfGenerator)

	locator := "AAAA"

	testOrder1 := models.Move{
		Locator: locator,
		Orders:  models.Order{},
	}

	outputputTest1, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderUpload, testOrder1, "")
	suite.FatalNil(outputputTest1)
	suite.Assertions.IsType(apperror.UnprocessableEntityError{}, err)
	err = service.CleanupFile(outputputTest1)
	suite.NoError(err)
	testOrder2 := models.Move{
		Locator: locator,
		Orders:  models.Order{},
	}
	testOrder3, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.MoveOrderAmendmentUpload, testOrder2, "")
	suite.FatalNil(testOrder3)
	suite.Assertions.IsType(apperror.UnprocessableEntityError{}, err)
	err = service.CleanupFile(testOrder3)
	suite.NoError(err)
}

func (suite *PaperworkServiceSuite) pdfFileInfo(generator *paperwork.Generator, file afero.File) (*pdfcpu.PDFInfo, error) {
	return api.PDFInfo(file, file.Name(), nil, false, generator.PdfConfiguration())
}

func (suite *PaperworkServiceSuite) setupOrdersDocument() (services.PrimeDownloadMoveUploadPDFGenerator, models.Order) {
	order := factory.BuildOrder(suite.DB(), nil, nil)

	document := factory.BuildDocument(suite.DB(), nil, nil)

	file, err := suite.openLocalFile("../../paperwork/testdata/orders1.jpg", suite.userUploader.FileSystem())
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("../../paperwork/testdata/orders1.pdf", suite.userUploader.FileSystem())
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("../../paperwork/testdata/orders2.jpg", suite.userUploader.FileSystem())
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	err = suite.DB().Load(&document, "UserUploads.Upload")
	suite.FatalNil(err)
	suite.Equal(3, len(document.UserUploads))

	order.UploadedOrders = document
	order.UploadedOrdersID = document.ID
	suite.MustSave(&order)

	pdfGenerator, err := paperwork.NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)
	service, err := NewMoveUserUploadToPDFDownloader(pdfGenerator)
	if err != nil {
		suite.FatalNil(err)
	}
	return service, order
}
