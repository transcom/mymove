package paperwork

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/hhrutter/pdfcpu/pkg/api"
	"github.com/hhrutter/pdfcpu/pkg/pdfcpu"
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
	workDir  string
}

// NewGenerator creates a new Generator.
func NewGenerator(db *pop.Connection, logger *zap.Logger, uploader *uploader.Uploader) (*Generator, error) {
	directory, err := ioutil.TempDir("", "generator")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Generator{
		db:       db,
		logger:   logger,
		uploader: uploader,
		workDir:  directory,
	}, nil
}

type inputFile struct {
	Path        string
	ContentType string
}

func (g *Generator) newTempFile() (*os.File, error) {
	outputFile, err := ioutil.TempFile(g.workDir, "temp")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return outputFile, nil
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

	outputFile, err := g.newTempFile()
	if err != nil {
		return "", err
	}

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

// GenerateAdvancePaperwork generates the advance paperwork for a move.
// Outputs to a tempfile
func (g *Generator) GenerateAdvancePaperwork(moveID uuid.UUID, build string) (string, error) {
	move, err := models.FetchMoveForAdvancePaperwork(g.db, moveID)
	if err != nil {
		return "", err
	}

	summary := NewShipmentSummary(&move)
	outfile, err := g.newTempFile()
	if err != nil {
		return "", err
	}
	if err := summary.DrawForm(outfile); err != nil {
		return "", err
	}
	outfile.Close()

	generatedPath := outfile.Name()
	ordersPaths, err := g.GenerateOrderPDF(move.OrdersID)
	if err != nil {
		return "", err
	}

	mergedFile, err := g.newTempFile()
	if err != nil {
		return "", err
	}

	var inputFiles []string
	g.logger.Debug("adding orders and shipment summary to packet", zap.Any("inputFiles", inputFiles))
	inputFiles = append(ordersPaths, generatedPath)

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance.MethodOfReceipt == models.MethodOfReceiptOTHERDD {
			g.logger.Debug("adding direct deposit form to packet", zap.Any("inputFiles", inputFiles))
			ddFormPath := filepath.Join(build, "/downloads/direct_deposit_form.pdf")
			inputFiles = append(inputFiles, ddFormPath)
			break
		}
	}

	config := pdfcpu.NewDefaultConfiguration()
	if err = api.Merge(inputFiles, mergedFile.Name(), config); err != nil {
		return "", err
	}

	return mergedFile.Name(), nil

}

// MergeImagesToPDF creates a PDF containing the images at the specified paths.
//
// The content type of the image is inferred from its extension. If this proves to
// be insufficient, storage.DetectContentType and contentTypeToImageType above can
// be used.
func (g *Generator) MergeImagesToPDF(paths []string) (string, error) {
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
