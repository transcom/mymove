package paperwork

import (
	"os"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkSuite) setupOrdersDocumentForMoveOrderDownloadGenerator() (*MoveOrderDownloadPdfGenerator, models.Order) {
	order := factory.BuildOrder(suite.DB(), nil, nil)

	document := factory.BuildDocument(suite.DB(), nil, nil)

	generator, err := NewMoveOrderPdfGeneratorForTesting(suite.AppContextForTest(), suite.userUploader.Uploader().Storer)
	suite.FatalNil(err)

	file, err := suite.openLocalFile("testdata/orders1.jpg", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("testdata/orders1.pdf", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("testdata/orders2.jpg", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	err = suite.DB().Load(&document, "UserUploads.Upload")
	suite.FatalNil(err)
	suite.Equal(3, len(document.UserUploads))

	order.UploadedOrders = document
	order.UploadedOrdersID = document.ID
	//order.UploadedAmendedOrdersID = &document.ID
	suite.MustSave(&order)

	return generator, order
}

func (suite *PaperworkSuite) TestGenerate() {

	generator, order := suite.setupOrdersDocumentForMoveOrderDownloadGenerator()

	locator := "AAAA"

	customMoveWithOnlyOrders := models.Move{
		Locator: locator,
		Orders:  order,
	}

	pdfFileWithOnlyOrders, err := generator.GeneratePdf(DownloadAll, customMoveWithOnlyOrders)
	suite.FatalNil(err)

	pdfGenerator, err := NewGenerator(suite.userUploader.Uploader())
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
	pdfFileWithAll, err := generator.GeneratePdf(DownloadAll, customMoveWithOrdersAndAmendments)
	suite.FatalNil(err)
	// Verify generated files have (3 x 2) pages for both orders and amendments. see setup data for upload count
	fileInfoAll, err := pdfGenerator.GetPdfFileInfo(pdfFileWithAll.Name())
	suite.FatalNil(err)
	suite.Equal(6, fileInfoAll.PageCount)

	pdfFileWithOnlyAmendments, err := generator.GeneratePdf(DownloadOnlyAmendments, customMoveWithOrdersAndAmendments)
	suite.FatalNil(err)
	// Verify only amendments are generated
	fileInfoOnlyAmendments, err := pdfGenerator.GetPdfFileInfo(pdfFileWithOnlyAmendments.Name())
	suite.FatalNil(err)
	suite.Equal(3, fileInfoOnlyAmendments.PageCount)
}

func (suite *PaperworkSuite) TestCleanupForMoveOrderDownloadGenerator() {
	generator, order := suite.setupOrdersDocument()

	uploads, err := models.UploadsFromUserUploads(suite.DB(), order.UploadedOrders.UserUploads)
	suite.FatalNil(err)
	_, err = generator.CreateMergedPDFUpload(suite.AppContextForTest(), uploads)
	suite.FatalNil(err)

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return value in the file is used for test database teardown
	//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	generator.Cleanup(suite.AppContextForTest())

	fs := suite.userUploader.FileSystem()
	exists, existsErr := fs.DirExists(generator.workDir)
	suite.Nil(existsErr)

	if exists {
		suite.Failf("expected %s to not be a directory, but it was", generator.workDir)

		var paths []string
		walkErr := fs.Walk(generator.workDir, func(path string, info os.FileInfo, err error) error {
			if path != generator.workDir { // Walk starts off with the directory passed to it
				paths = append(paths, path)
			}
			return nil
		})
		suite.Nil(walkErr)
		if len(paths) > 0 {
			suite.Failf("did not clean up", "expected %s to be empty, but it contained %v", generator.workDir, paths)
		}
	}
}
