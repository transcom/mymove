package paperwork

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/hhrutter/pdfcpu/pkg/api"
	"github.com/hhrutter/pdfcpu/pkg/pdfcpu"

	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaperworkSuite) sha256ForPath(path string) (string, error) {
	file, err := os.Open(path)
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

func (suite *PaperworkSuite) TestPDFFromImages() {
	images := []inputFile{
		{Path: "testdata/orders1.jpg", ContentType: "image/jpeg"},
		{Path: "testdata/orders2.jpg", ContentType: "image/jpeg"},
	}
	generator, err := NewGenerator(suite.db, suite.logger, suite.uploader)
	suite.FatalNil(err)

	generatedPath, err := generator.pdfFromImages(images)
	suite.Nil(err, "failed to generate pdf")
	suite.NotEmpty(generatedPath, "got an empty path to the generated file")
	defer os.Remove(generatedPath)

	// verify that the images are in the pdf by extracting them and checking their checksums
	tmpdir, err := ioutil.TempDir("", "extracted_images")
	suite.Nil(err, "could not create temp dir")
	defer os.RemoveAll(tmpdir)

	pdfcpuConfig := pdfcpu.NewDefaultConfiguration()
	api.ExtractImages(generatedPath, tmpdir, []string{"-2"}, pdfcpuConfig)

	checksums := make([]string, 2)
	files, err := ioutil.ReadDir(tmpdir)
	suite.Nil(err)

	suite.Equal(2, len(files), "did not find 2 images")

	for _, file := range files {
		checksum, err := suite.sha256ForPath(path.Join(tmpdir, file.Name()))
		suite.FatalNil(err, "error calculating hash")
		if err != nil {
			suite.FailNow(err.Error())
		}
		checksums = append(checksums, checksum)
	}

	orders1Checksum, err := suite.sha256ForPath("testdata/orders1.jpg")
	suite.Nil(err, "error calculating hash")
	suite.Contains(checksums, orders1Checksum, "did not find hash for orders1.jpg")

	orders2Checksum, err := suite.sha256ForPath("testdata/orders2.jpg")
	suite.Nil(err, "error calculating hash")
	suite.Contains(checksums, orders2Checksum, "did not find hash for orders2.jpg")
}

func (suite *PaperworkSuite) TestGenerateUploadsPDF() {
	order := testdatagen.MakeDefaultOrder(suite.db)

	document := testdatagen.MakeDefaultDocument(suite.db)

	order.UploadedOrders = document
	order.UploadedOrdersID = document.ID
	suite.mustSave(&order)

	file, err := uploader.NewLocalFile("testdata/orders1.jpg")
	suite.FatalNil(err)
	_, _, err = suite.uploader.CreateUpload(&document.ID, document.ServiceMember.UserID, file)
	suite.FatalNil(err)

	file, err = uploader.NewLocalFile("testdata/orders1.pdf")
	suite.FatalNil(err)
	_, _, err = suite.uploader.CreateUpload(&document.ID, document.ServiceMember.UserID, file)
	suite.FatalNil(err)

	file, err = uploader.NewLocalFile("testdata/orders2.jpg")
	suite.Nil(err)
	_, _, err = suite.uploader.CreateUpload(&document.ID, document.ServiceMember.UserID, file)
	suite.Nil(err)

	err = suite.db.Load(&document, "Uploads")
	suite.FatalNil(err)
	suite.Equal(3, len(document.Uploads))

	generator, err := NewGenerator(suite.db, suite.logger, suite.uploader)
	suite.FatalNil(err)

	paths, err := generator.GenerateUploadsPDF(document.Uploads)
	suite.FatalNil(err)

	suite.Equal(3, len(paths), "wrong number of paths returned")
}
