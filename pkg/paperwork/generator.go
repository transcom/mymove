package paperwork

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/trussworks/pdfcpu/pkg/api"
	"github.com/trussworks/pdfcpu/pkg/pdfcpu"
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
	db        *pop.Connection
	fs        *afero.Afero
	logger    Logger
	uploader  *uploader.Uploader
	pdfConfig *pdfcpu.Configuration
	workDir   string
}

// Converts an image of any type to a PNG with 8-bit color depth
func convertTo8BitPNG(in io.Reader, out io.Writer) error {
	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	b := img.Bounds()
	imgSet := image.NewRGBA(b)
	// Converts each pixel to a 32-bit RGBA pixel
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			newPixel := color.RGBAModel.Convert(img.At(x, y))
			imgSet.Set(x, y, newPixel)
		}
	}

	err = png.Encode(out, imgSet)
	if err != nil {
		return err
	}

	return nil
}

// NewGenerator creates a new Generator.
func NewGenerator(db *pop.Connection, logger Logger, uploader *uploader.Uploader) (*Generator, error) {
	afs := uploader.Storer.FileSystem()

	pdfConfig := pdfcpu.NewInMemoryConfiguration()
	pdfConfig.FileSystem = afs.Fs

	directory, err := afs.TempDir("", "generator")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Generator{
		db:        db,
		fs:        afs,
		logger:    logger,
		uploader:  uploader,
		pdfConfig: pdfConfig,
		workDir:   directory,
	}, nil
}

type inputFile struct {
	Path        string
	ContentType string
}

func (g *Generator) newTempFile() (afero.File, error) {
	outputFile, err := g.fs.TempFile(g.workDir, "temp")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return outputFile, nil
}

// CreateMergedPDFUpload converts Uploads to PDF and merges them into a single PDF
func (g *Generator) CreateMergedPDFUpload(uploads models.Uploads) (afero.File, error) {
	pdfs, err := g.ConvertUploadsToPDF(uploads)
	if err != nil {
		return nil, errors.Wrap(err, "Error while converting uploads")
	}

	mergedPdf, err := g.MergePDFFiles(pdfs)
	if err != nil {
		return nil, errors.Wrap(err, "Error while merging PDFs")
	}

	return mergedPdf, err
}

// ConvertUploadsToPDF turns a slice of Uploads into a slice of paths to converted PDF files
func (g *Generator) ConvertUploadsToPDF(uploads models.Uploads) ([]string, error) {
	// tempfile paths to be returned
	pdfs := make([]string, 0)

	// path for each image once downloaded
	images := make([]inputFile, 0)

	for _, upload := range uploads {
		if upload.ContentType == "application/pdf" {
			if len(images) > 0 {
				// We want to retain page order and will generate a PDF for images
				// that have already been encountered before handling this PDF.
				pdf, err := g.PDFFromImages(images)
				if err != nil {
					return nil, errors.Wrap(err, "Converting images")
				}
				pdfs = append(pdfs, pdf)
				images = make([]inputFile, 0)
			}
		}

		download, err := g.uploader.Download(&upload)
		if err != nil {
			return nil, errors.Wrap(err, "Downloading file from upload")
		}
		defer download.Close()

		outputFile, err := g.newTempFile()
		if err != nil {
			return nil, errors.Wrap(err, "Creating temp file")
		}

		_, err = io.Copy(outputFile, download)
		if err != nil {
			return nil, errors.Wrap(err, "Copying to afero file")
		}

		path := outputFile.Name()

		if upload.ContentType == "application/pdf" {
			pdfs = append(pdfs, path)
		} else {
			images = append(images, inputFile{Path: path, ContentType: upload.ContentType})
		}
	}

	// Merge all remaining images in urls into a new PDF
	if len(images) > 0 {
		pdf, err := g.PDFFromImages(images)
		if err != nil {
			return nil, errors.Wrap(err, "Converting remaining images to pdf")
		}
		pdfs = append(pdfs, pdf)
	}

	for _, f := range pdfs {
		err := api.Validate(f, g.pdfConfig)
		if err != nil {
			return nil, errors.Wrap(err, "Validating pdfs")
		}
	}

	return pdfs, nil
}

// convert between image MIME types and the values expected by gofpdf
var contentTypeToImageType = map[string]string{
	"image/jpeg": "JPG",
	"image/png":  "PNG",
}

// PDFFromImages returns the path to tempfile PDF containing all images included
// in urls.
//
// The files at those paths will be tempfiles that will need to be cleaned
// up by the caller.
func (g *Generator) PDFFromImages(images []inputFile) (string, error) {
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
		file, _ := g.fs.Open(image.Path)
		if image.ContentType == "image/png" {
			// gofpdf isn't able to process 16-bit PNGs, so to be safe we convert all PNGs to an 8-bit color depth
			newFile, err := g.newTempFile()
			if err != nil {
				return "", errors.Wrap(err, "Creating temp file for png conversion")
			}
			err = convertTo8BitPNG(file, newFile)
			if err != nil {
				return "", errors.Wrap(err, "Converting to 8-bit png")
			}
			defer file.Close()
			file = newFile
			file.Seek(0, io.SeekStart)
		}
		// Need to register the image using an afero reader, else it uses default filesystem
		pdf.RegisterImageReader(image.Path, contentTypeToImageType[image.ContentType], file)
		opt.ImageType = contentTypeToImageType[image.ContentType]
		pdf.ImageOptions(image.Path, horizontalMargin, topMargin, bodyWidth, 0, false, opt, 0, "")
	}

	if err = pdf.OutputAndClose(outputFile); err != nil {
		return "", errors.Wrap(err, "could not write PDF to outputfile")
	}
	return outputFile.Name(), nil
}

// MergePDFFiles Merges a slice of paths to PDF files into a single PDF
func (g *Generator) MergePDFFiles(paths []string) (afero.File, error) {
	mergedFile, err := g.newTempFile()
	if err != nil {
		return mergedFile, err
	}

	if err = api.Merge(paths, mergedFile.Name(), g.pdfConfig); err != nil {
		return mergedFile, err
	}

	// Reload the file from memstore
	mergedFile, err = g.fs.Open(mergedFile.Name())
	if err != nil {
		return mergedFile, err
	}

	return mergedFile, nil
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

	return g.PDFFromImages(images)
}
