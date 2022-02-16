package paperwork

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkSuite) sha256ForPath(path string, fs *afero.Afero) (string, error) {
	var file afero.File
	var err error
	if fs != nil {
		file, err = fs.Open(filepath.Clean(path))
	} else {
		file, err = os.Open(filepath.Clean(path))
	}
	if err != nil {
		suite.NoError(err)
	}
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to close a local server connection to ensure a unit test server is not left running indefinitely
	//RA: Given the functions causing the lint errors are used to close a local server connection for testing purposes, it is not deemed a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
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

	generator, err := NewGenerator(suite.userUploader.Uploader())
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
	suite.MustSave(&order)

	return generator, order
}

func (suite *PaperworkSuite) TestPDFFromImages() {
	generator, newGeneratorErr := NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(newGeneratorErr)

	images := []inputFile{
		{Path: "testdata/orders1.jpg", ContentType: "image/jpeg"},
		{Path: "testdata/orders2.jpg", ContentType: "image/jpeg"},
	}
	for _, image := range images {
		_, err := suite.openLocalFile(image.Path, generator.fs)
		suite.FatalNil(err)
	}

	generatedPath, err := generator.PDFFromImages(suite.AppContextForTest(), images)
	suite.FatalNil(err, "failed to generate pdf")
	aferoFile, err := generator.fs.Open(generatedPath)
	suite.FatalNil(err, "afero failed to open pdf")

	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
	suite.FatalNil(err)

	// verify that the images are in the pdf by extracting them and checking their checksums
	file, err := afero.ReadAll(aferoFile)
	suite.FatalNil(err)
	tmpDir, err := ioutil.TempDir("", "images")
	suite.FatalNil(err)
	f, err := ioutil.TempFile(tmpDir, "")
	suite.FatalNil(err)
	err = ioutil.WriteFile(f.Name(), file, os.ModePerm)
	suite.FatalNil(err)
	err = api.ExtractImages(f, tmpDir, []string{"-2"}, generator.pdfConfig)
	suite.FatalNil(err)
	err = os.Remove(f.Name())
	suite.FatalNil(err)

	checksums := make([]string, 2)
	files, err := ioutil.ReadDir(tmpDir)
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
	generator, err := NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)

	images := []inputFile{
		// The below image isn't getting extracted by pdfcpu for some reason.
		// We're adding it because gofpdf can't process 16-bit PNG images, so we
		// just care that PDFFromImages doesn't error
		{Path: "testdata/16bitpng.png", ContentType: "image/png"},
	}
	_, err = suite.openLocalFile(images[0].Path, generator.fs)
	suite.FatalNil(err)

	generatedPath, err := generator.PDFFromImages(suite.AppContextForTest(), images)
	suite.FatalNil(err, "failed to generate pdf")
	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
}

func (suite *PaperworkSuite) TestPDFFromImagesRotation() {
	generator, err := NewGenerator(suite.userUploader.Uploader())
	suite.FatalNil(err)

	images := []inputFile{
		// The below image is best viewed in landscape, but will rotate in
		// PDFFromImages. Since we can't analyze the final contents, we'll
		// just ensure it doesn't error.
		{Path: "testdata/example_landscape.png", ContentType: "image/png"},
		{Path: "testdata/example_landscape.jpg", ContentType: "image/jpeg"},
	}
	_, err = suite.openLocalFile(images[0].Path, generator.fs)
	suite.FatalNil(err)
	_, err = suite.openLocalFile(images[1].Path, generator.fs)
	suite.FatalNil(err)

	generatedPath, err := generator.PDFFromImages(suite.AppContextForTest(), images)
	suite.FatalNil(err, "failed to generate pdf")
	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
}

func (suite *PaperworkSuite) TestGenerateUploadsPDF() {
	generator, order := suite.setupOrdersDocument()

	uploads, err := models.UploadsFromUserUploads(suite.DB(), order.UploadedOrders.UserUploads)
	suite.FatalNil(err)
	paths, err := generator.ConvertUploadsToPDF(suite.AppContextForTest(), uploads)
	suite.FatalNil(err)

	suite.Equal(3, len(paths), "wrong number of paths returned")
}

func (suite *PaperworkSuite) TestCreateMergedPDF() {
	generator, order := suite.setupOrdersDocument()

	uploads, err := models.UploadsFromUserUploads(suite.DB(), order.UploadedOrders.UserUploads)
	suite.FatalNil(err)
	file, err := generator.CreateMergedPDFUpload(suite.AppContextForTest(), uploads)
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
