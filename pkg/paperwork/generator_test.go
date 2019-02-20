package paperwork

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/spf13/afero"
	"github.com/trussworks/pdfcpu/pkg/api"
	"github.com/trussworks/pdfcpu/pkg/pdfcpu"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaperworkSuite) sha256ForPath(path string, fs *afero.Afero) (string, error) {
	file, err := fs.Open(path)
	if err != nil {
		suite.Nil(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		suite.Nil(err)
	}

	byteArray := hash.Sum(nil)
	return fmt.Sprintf("%x", byteArray), nil
}

func (suite *PaperworkSuite) setupOrdersDocument() (*Generator, models.Order) {
	order := testdatagen.MakeDefaultOrder(suite.DB())

	document := testdatagen.MakeDefaultDocument(suite.DB())

	generator, err := NewGenerator(suite.DB(), suite.logger, suite.uploader)
	suite.FatalNil(err)

	file, err := suite.openLocalFile("testdata/orders1.jpg", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.uploader.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, file, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("testdata/orders1.pdf", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.uploader.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, file, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	file, err = suite.openLocalFile("testdata/orders2.jpg", generator.fs)
	suite.FatalNil(err)

	_, _, err = suite.uploader.CreateUploadForDocument(&document.ID, document.ServiceMember.UserID, file, uploader.AllowedTypesAny)
	suite.FatalNil(err)

	err = suite.DB().Load(&document, "Uploads")
	suite.FatalNil(err)
	suite.Equal(3, len(document.Uploads))

	order.UploadedOrders = document
	order.UploadedOrdersID = document.ID
	suite.MustSave(&order)

	return generator, order
}

func (suite *PaperworkSuite) TestPDFFromImages() {
	generator, err := NewGenerator(suite.DB(), suite.logger, suite.uploader)
	suite.FatalNil(err)

	images := []inputFile{
		{Path: "testdata/orders1.jpg", ContentType: "image/jpeg"},
		{Path: "testdata/orders2.jpg", ContentType: "image/jpeg"},
	}
	for _, image := range images {
		_, err := suite.openLocalFile(image.Path, generator.fs)
		suite.FatalNil(err)
	}

	generatedPath, err := generator.PDFFromImages(images)
	suite.FatalNil(err, "failed to generate pdf")
	suite.NotEmpty(generatedPath, "got an empty path to the generated file")

	// verify that the images are in the pdf by extracting them and checking their checksums
	tmpdir, err := afero.TempDir(generator.fs, "", "extracted_images")
	suite.FatalNil(err, "could not create temp dir")

	api.ExtractImages(generatedPath, tmpdir, []string{"-2"}, generator.pdfConfig)

	checksums := make([]string, 2)
	files, err := afero.ReadDir(generator.fs, tmpdir)
	suite.FatalNil(err)

	suite.Equal(2, len(files), "did not find 2 images")

	for _, file := range files {
		checksum, err := suite.sha256ForPath(path.Join(tmpdir, file.Name()), generator.fs)
		suite.FatalNil(err, "error calculating hash")
		if err != nil {
			suite.FailNow(err.Error())
		}
		checksums = append(checksums, checksum)
	}

	orders1Checksum, err := suite.sha256ForPath("testdata/orders1.jpg", generator.fs)
	suite.Nil(err, "error calculating hash")
	suite.Contains(checksums, orders1Checksum, "did not find hash for orders1.jpg")

	orders2Checksum, err := suite.sha256ForPath("testdata/orders2.jpg", generator.fs)
	suite.Nil(err, "error calculating hash")
	suite.Contains(checksums, orders2Checksum, "did not find hash for orders2.jpg")
}

func (suite *PaperworkSuite) TestPDFFromImages16BitPNG() {
	generator, err := NewGenerator(suite.DB(), suite.logger, suite.uploader)
	suite.FatalNil(err)

	images := []inputFile{
		// The below image isn't getting extracted by pdfcpu for some reason.
		// We're adding it because gofpdf can't process 16-bit PNG images, so we
		// just care that PDFFromImages doesn't error
		{Path: "testdata/16bitpng.png", ContentType: "image/png"},
	}
	_, err = suite.openLocalFile(images[0].Path, generator.fs)
	suite.FatalNil(err)

	generatedPath, err := generator.PDFFromImages(images)
	suite.FatalNil(err, "failed to generate pdf")
	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
}

func (suite *PaperworkSuite) TestGenerateUploadsPDF() {
	generator, order := suite.setupOrdersDocument()

	paths, err := generator.ConvertUploadsToPDF(order.UploadedOrders.Uploads)
	suite.FatalNil(err)

	suite.Equal(3, len(paths), "wrong number of paths returned")
}

func (suite *PaperworkSuite) TestCreateMergedPDF() {
	generator, order := suite.setupOrdersDocument()

	file, err := generator.CreateMergedPDFUpload(order.UploadedOrders.Uploads)
	suite.FatalNil(err)

	// Read merged file and verify page count
	ctx, err := api.Read(file.Name(), generator.pdfConfig)
	suite.FatalNil(err)

	err = pdfcpu.ValidateXRefTable(ctx.XRefTable)
	suite.FatalNil(err)

	suite.Equal(3, ctx.PageCount)
}
