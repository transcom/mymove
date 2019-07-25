package paperwork

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/hhrutter/pdfcpu/pkg/pdfcpu/validate"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/hhrutter/pdfcpu/pkg/api"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaperworkSuite) sha256ForPath(path string, fs *afero.Afero) (string, error) {
	var file afero.File
	var err error
	if fs != nil {
		file, err = fs.Open(path)
	} else {
		file, err = os.Open(path)
	}
	if err != nil {
		suite.NoError(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		suite.NoError(err)
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
	generator, newGeneratorErr := NewGenerator(suite.DB(), suite.logger, suite.uploader)
	suite.FatalNil(newGeneratorErr)

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
	aferoFile, err := generator.fs.Open(generatedPath)
	suite.FatalNil(err, "afero failed to open pdf")

	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
	suite.FatalNil(err)

	// verify that the images are in the pdf by extracting them and checking their checksums
	file, err := afero.ReadAll(aferoFile)
	suite.FatalNil(err)
	tmpDir, err := ioutil.TempDir("", "images")
	f, err := ioutil.TempFile(tmpDir, "")
	err = ioutil.WriteFile(f.Name(), file, os.ModePerm)
	suite.FatalNil(err)
	err = api.ExtractImages(f, tmpDir, []string{"-2"}, generator.pdfConfig)
	suite.FatalNil(err)
	err = os.Remove(f.Name())
	suite.FatalNil(err)

	checksums := make([]string, 2)
	files, err := ioutil.ReadDir(tmpDir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
	suite.FatalNil(err)

	suite.Equal(2, len(files), "did not find 2 images")

	for _, file := range files {
		checksum, sha256ForPathErr := suite.sha256ForPath(path.Join(tmpDir, file.Name()), nil)
		suite.FatalNil(sha256ForPathErr, "error calculating hash")
		if sha256ForPathErr != nil {
			suite.FailNow(sha256ForPathErr.Error())
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

	uploads := order.UploadedOrders.Uploads
	file, err := generator.CreateMergedPDFUpload(uploads)
	suite.FatalNil(err)

	// Read merged file and verify page count
	ctx, err := api.ReadContext(file, generator.pdfConfig)
	suite.FatalNil(err)

	err = validate.XRefTable(ctx.XRefTable)
	suite.FatalNil(err)

	suite.Equal(3, ctx.PageCount)
}

func (suite *PaperworkSuite) TestCleanup() {
	generator, order := suite.setupOrdersDocument()

	uploads := order.UploadedOrders.Uploads
	_, err := generator.CreateMergedPDFUpload(uploads)
	suite.FatalNil(err)

	generator.Cleanup()

	fs := suite.uploader.Storer.FileSystem()
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
		suite.Failf("did not clean up", "expected %s to be empty, but it contained %v", generator.workDir, paths)
	}
}
