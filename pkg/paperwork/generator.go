package paperwork

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// Default values for PDF generation
const (
	PdfOrientation string  = "P"
	PdfUnit        string  = "mm"
	PdfPageWidth   float64 = 210.0
	PdfPageSize    string  = "A4"
	PdfFontDir     string  = ""
)

// Generator encapsulates the prerequisites for PDF generation.
type Generator struct {
	db       *pop.Connection
	logger   *zap.Logger
	uploader *uploader.Uploader
}

// NewGenerator creates a new Generator.
func NewGenerator(db *pop.Connection, logger *zap.Logger, uploader *uploader.Uploader) *Generator {
	return &Generator{
		db:       db,
		logger:   logger,
		uploader: uploader,
	}
}

// GenerateAdvancePaperwork generates the advance paperwork for a move.
func (g *Generator) GenerateAdvancePaperwork(moveID uuid.UUID) error {
	move, err := models.FetchMoveForAdvancePaperwork(g.db, moveID)
	if err != nil {
		return err
	}
	fmt.Println(move)
	return nil
}

type inputFile struct {
	Path        string
	ContentType string
}

// GenerateOrderPDF returns a slice of paths to PDF files that represent all files
// uploaded for Orders.
func (g *Generator) GenerateOrderPDF(orderID uuid.UUID) ([]string, error) {
	order, err := models.FetchOrderForPDFConversion(g.db, orderID)
	if err != nil {
		return nil, err
	}

	// tempfile paths to be returned
	pdfs := make([]string, 0)

	// path for each image once downloaded
	images := make([]inputFile, 0)

	for _, upload := range order.UploadedOrders.Uploads {
		if upload.ContentType == "application/pdf" {
			if len(images) > 0 {
				// We want to retain page order and will generate a PDF for images
				// that have already been encountered before handling this PDF.
				pdf, err := g.pdfFromImages(images)
				if err != nil {
					return nil, err
				}
				pdfs = append(pdfs, pdf)
				images = make([]inputFile, 0)
			}
		}

		path, err := g.uploader.Download(&upload)
		if err != nil {
			return nil, err
		}
		if upload.ContentType == "application/pdf" {
			pdfs = append(pdfs, path)
		} else {
			images = append(images, inputFile{Path: path, ContentType: upload.ContentType})
		}
	}

	// Merge all images in urls into a new PDF
	pdf, err := g.pdfFromImages(images)
	if err != nil {
		return nil, err
	}
	pdfs = append(pdfs, pdf)

	return pdfs, nil
}

// convert between image MIME types and the values expected by gofpdf
var contentTypeToImageType = map[string]string{
	"image/jpeg": "JPG",
	"image/png":  "PNG",
}

// pdfFromImageURLs returns the path to tempfile PDF containing all images included
// in urls.
//
// The files at those paths will be tempfiles that will need to be cleaned
// up by the caller.
func (g *Generator) pdfFromImages(images []inputFile) (string, error) {
	horizontalMargin := 0.0
	topMargin := 0.0
	bodyWidth := PdfPageWidth - (horizontalMargin * 2)

	pdf := gofpdf.New(PdfOrientation, PdfUnit, PdfPageSize, PdfFontDir)
	pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	if len(images) == 0 {
		return "", errors.New("No images provided")
	}

	g.logger.Debug("generating PDF from image files", zap.Any("images", images))

	// TODO create a temp dir for use by this generator that can be easily cleaned up
	outputFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		return "", errors.WithStack(err)
	}
	// defer os.Remove(outputFile.Name())

	var opt gofpdf.ImageOptions
	for _, image := range images {
		pdf.AddPage()
		opt.ImageType = contentTypeToImageType[image.ContentType]
		pdf.ImageOptions(image.Path, horizontalMargin, topMargin, bodyWidth, 0, false, opt, 0, "")
	}

	if err = pdf.OutputAndClose(outputFile); err != nil {
		return "", errors.Wrap(err, "could not write PDF to outputfile")
	}
	return outputFile.Name(), nil
}

// MergeLocalFiles creates a PDF containing the images at the specified paths.
//
// The content type of the image is inferred from its extension. If this proves to
// be insufficient, storage.DetectContentType and contentTypeToImageType above can
// be used.
func (g *Generator) MergeLocalFiles(paths []string) (string, error) {
	// path and type for each image
	images := make([]inputFile, 0)

	for _, path := range paths {
		extension := filepath.Ext(path)[1:]
		images = append(images, inputFile{
			Path:        path,
			ContentType: strings.ToUpper(extension),
		})
	}

	return g.pdfFromImages(images)
}
