package paperwork

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkServiceSuite) TestPrimeDownloadMoveUploadPDFGenerator() {
	service, order := suite.setupOrdersDocument()

	locator := "AAAA"

	customMoveWithOnlyOrders := models.Move{
		Locator: locator,
		Orders:  order,
	}

	pdfFileWithOnlyOrders, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.DownloadMoveOrderUploadTypeAll, customMoveWithOnlyOrders)
	suite.FatalNil(err)

	pdfGenerator, err := paperwork.NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)

	// Verify generated files have 3 pages. see setup data for upload count
	fileInfo, err := pdfGenerator.GetPdfFileInfo(pdfFileWithOnlyOrders.Name())
	suite.FatalNil(err)
	// Verify the 3 uploads were seperated out and images were mot grouped together
	suite.Equal(3, fileInfo.PageCount)

	// Point amendments doc to UploadedOrdersID.
	order.UploadedAmendedOrdersID = &order.UploadedOrdersID
	customMoveWithOrdersAndAmendments := models.Move{
		Locator: "AAAA",
		Orders:  order,
	}
	pdfFileWithAll, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.DownloadMoveOrderUploadTypeAll, customMoveWithOrdersAndAmendments)
	suite.FatalNil(err)
	// Verify generated files have (3 x 2) pages for both orders and amendments. see setup data for upload count
	fileInfoAll, err := pdfGenerator.GetPdfFileInfo(pdfFileWithAll.Name())
	suite.FatalNil(err)
	suite.Equal(6, fileInfoAll.PageCount)

	pdfFileWithOnlyAmendments, err := service.GenerateDownloadMoveUserUploadPDF(suite.AppContextForTest(), services.DownloadMoveOrderUploadTypeOnlyAmendments, customMoveWithOrdersAndAmendments)
	suite.FatalNil(err)
	// Verify only amendments are generated
	fileInfoOnlyAmendments, err := pdfGenerator.GetPdfFileInfo(pdfFileWithOnlyAmendments.Name())
	suite.FatalNil(err)
	suite.Equal(3, fileInfoOnlyAmendments.PageCount)
	suite.AfterTest()
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

	service, err := NewMoveUserUploadToPDFDownloader(true, suite.userUploader)
	if err != nil {
		suite.FatalNil(err)
	}
	return service, order
}
