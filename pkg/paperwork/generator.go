package paperwork

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jung-kurt/gofpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// Default values for PDF generation
const (
	PdfOrientation string  = "P"
	PdfUnit        string  = "pt"
	PdfPageWidth   float64 = 612.0
	PdfPageHeight  float64 = 792.0
	PdfPageSize    string  = "Letter"
	PdfFontDir     string  = ""
)

// Generator encapsulates the prerequisites for PDF generation.
type Generator struct {
	fs        *afero.Afero
	uploader  *uploader.Uploader
	pdfConfig *model.Configuration
	workDir   string
	pdfLib    PDFLibrary
}

type pdfCPUWrapper struct {
	*model.Configuration
}

// Merge merges files
func (pcw pdfCPUWrapper) Merge(files []io.ReadSeeker, w io.Writer) error {
	var rscs []io.ReadSeeker
	rscs = append(rscs, files...)
	return api.MergeRaw(rscs, w, false, pcw.Configuration) // Todo: False refers to a divider page. Find out what this does
}

// Validate validates the api configuration
func (pcw pdfCPUWrapper) Validate(rs io.ReadSeeker) error {
	return api.Validate(rs, pcw.Configuration)
}

// PDFLibrary is the PDF library interface
type PDFLibrary interface {
	Merge(rsc []io.ReadSeeker, w io.Writer) error
	Validate(rs io.ReadSeeker) error
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
func NewGenerator(uploader *uploader.Uploader) (*Generator, error) {
	// Use in memory filesystem for generation. Purpose is to not write
	// to hard disk due to restrictions in AWS storage. May need better long term solution.
	afs := storage.NewMemory(storage.NewMemoryParams("", "")).FileSystem()

	// Disable ConfiDir for AWS deployment purposes.
	// PDFCPU will attempt to create temp dir using os.create(hard disk).This will prevent it.
	api.DisableConfigDir()
	pdfConfig := model.NewDefaultConfiguration()
	pdfCPU := pdfCPUWrapper{Configuration: pdfConfig}

	directory, err := afs.TempDir("", "generator")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Generator{
		fs:        afs,
		uploader:  uploader,
		pdfConfig: pdfConfig,
		workDir:   directory,
		pdfLib:    pdfCPU,
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

func (g *Generator) newTempFileWithName(fileName string) (afero.File, error) {
	name := "temp"

	if fileName != "" {
		// by adding a * before the extension we tell TempFile to put its random number before the extension instead of after it
		extensionIndex := strings.LastIndex(fileName, ".")
		name = fileName[:extensionIndex] + strings.Replace(fileName[extensionIndex:], ".", "*.", 1)
	}

	outputFile, err := g.fs.TempFile(g.workDir, name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return outputFile, nil
}

// Cleanup removes filesystem working dir
func (g *Generator) Cleanup(_ appcontext.AppContext) error {
	return g.fs.RemoveAll(g.workDir)
}

// Get PDF Configuration (For Testing)
func (g *Generator) FileSystem() *afero.Afero {
	return g.fs
}

// Add bookmarks into a single PDF
func (g *Generator) AddPdfBookmarks(inputFile afero.File, bookmarks []pdfcpu.Bookmark) (afero.File, error) {

	buf := new(bytes.Buffer)
	replace := true
	err := api.AddBookmarks(inputFile, buf, bookmarks, replace, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error pdfcpu.api.AddBookmarks")
	}

	tempFile, err := g.newTempFile()
	if err != nil {
		return nil, err
	}

	// copy byte[] to temp file
	_, err = io.Copy(tempFile, buf)
	if err != nil {
		return nil, errors.Wrap(err, "error io.Copy on byte[] to temp")
	}

	// Reload the file from memstore
	pdfWithBookmarks, err := g.fs.Open(tempFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "error g.fs.Open on reload from memstore")
	}

	return pdfWithBookmarks, nil
}

// Get PDF Configuration (For Testing)
func (g *Generator) PdfConfiguration() *model.Configuration {
	return g.pdfConfig
}

// Get file information of a single PDF
func (g *Generator) GetPdfFileInfo(fileName string) (*pdfcpu.PDFInfo, error) {
	file, err := g.fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return api.PDFInfo(file, fileName, nil, g.pdfConfig)
}

func (g *Generator) GetPdfFileInfoForReadSeeker(rs io.ReadSeeker) (*pdfcpu.PDFInfo, error) {
	return api.PDFInfo(rs, "", nil, g.pdfConfig)
}

// Get file information of a single PDF
func (g *Generator) GetPdfFileInfoByContents(file afero.File) (*pdfcpu.PDFInfo, error) {
	return api.PDFInfo(file, file.Name(), nil, g.pdfConfig)
}

// CreateMergedPDFUpload converts Uploads to PDF and merges them into a single PDF
func (g *Generator) CreateMergedPDFUpload(appCtx appcontext.AppContext, uploads models.Uploads) (afero.File, error) {
	pdfs, err := g.ConvertUploadsToPDF(appCtx, uploads)
	if err != nil {
		return nil, errors.Wrap(err, "Error while converting uploads")
	}

	mergedPdf, err := g.MergePDFFiles(appCtx, pdfs)
	if err != nil {
		return nil, errors.Wrap(err, "Error while merging PDFs")
	}

	return mergedPdf, err
}

// ConvertUploadsToPDF turns a slice of Uploads into a slice of paths to converted PDF files
func (g *Generator) ConvertUploadsToPDF(appCtx appcontext.AppContext, uploads models.Uploads) ([]string, error) {
	// tempfile paths to be returned
	pdfs := make([]string, 0)

	// path for each image once downloaded
	images := make([]inputFile, 0)

	for _, upload := range uploads {
		copyOfUpload := upload // Make copy to avoid implicit memory aliasing of items from a range statement.
		if copyOfUpload.ContentType == uploader.FileTypePDF {
			if len(images) > 0 {
				// We want to retain page order and will generate a PDF for images
				// that have already been encountered before handling this PDF.
				pdf, err := g.PDFFromImages(appCtx, images)
				if err != nil {
					return nil, errors.Wrap(err, "Converting images")
				}
				pdfs = append(pdfs, pdf)
				images = make([]inputFile, 0)
			}
		}

		download, err := g.uploader.Download(appCtx, &copyOfUpload)
		if err != nil {
			return nil, errors.Wrap(err, "Downloading file from upload")
		}

		defer func() {
			if downloadErr := download.Close(); downloadErr != nil {
				appCtx.Logger().Debug("Failed to close file", zap.Error(downloadErr))
			}
		}()

		outputFile, err := g.newTempFile()

		if err != nil {
			return nil, errors.Wrap(err, "Creating temp file")
		}

		_, err = io.Copy(outputFile, download)
		if err != nil {
			return nil, errors.Wrap(err, "Copying to afero file")
		}

		path := outputFile.Name()

		if copyOfUpload.ContentType == uploader.FileTypePDF {
			pdfs = append(pdfs, path)
		} else {
			images = append(images, inputFile{Path: path, ContentType: copyOfUpload.ContentType})
		}
	}

	// Merge all remaining images in urls into a new PDF
	if len(images) > 0 {
		pdf, err := g.PDFFromImages(appCtx, images)
		if err != nil {
			return nil, errors.Wrap(err, "Converting remaining images to pdf")
		}
		pdfs = append(pdfs, pdf)
	}

	for _, fn := range pdfs {
		f, err := g.fs.Open(fn)
		if err != nil {
			return nil, errors.Wrap(err, "Validating pdfs")
		}
		err = g.pdfLib.Validate(f)
		if err != nil {
			return nil, errors.Wrap(err, "Validating pdfs")
		}
	}

	return pdfs, nil
}

func (g *Generator) ConvertUploadToPDF(appCtx appcontext.AppContext, upload models.Upload) (string, error) {

	download, err := g.uploader.Download(appCtx, &upload)
	if err != nil {
		return "nil", errors.Wrap(err, "Downloading file from upload")
	}

	defer func() {
		if downloadErr := download.Close(); downloadErr != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(downloadErr))
		}
	}()

	outputFile, err := g.newTempFile()

	if err != nil {
		return "nil", errors.Wrap(err, "Creating temp file")
	}

	_, err = io.Copy(outputFile, download)
	if err != nil {
		return "nil", errors.Wrap(err, "Copying to afero file")
	}

	path := outputFile.Name()

	if upload.ContentType == uploader.FileTypePDF {
		return path, nil
	}

	images := make([]inputFile, 0)
	images = append(images, inputFile{Path: path, ContentType: upload.ContentType})
	return g.PDFFromImages(appCtx, images)
}

// convert between image MIME types and the values expected by gofpdf
var contentTypeToImageType = map[string]string{
	uploader.FileTypeJPEG: "JPG",
	uploader.FileTypePNG:  "PNG",
}

// ReduceUnusedSpace reduces unused space
func ReduceUnusedSpace(_ appcontext.AppContext, file afero.File, g *Generator, contentType string) (imgFile afero.File, width float64, height float64, err error) {
	// Figure out if the image should be rotated by calculating height and width of image.
	pic, _, decodeErr := image.Decode(file)
	if decodeErr != nil {
		return nil, 0.0, 0.0, errors.Wrapf(decodeErr, "file %s was not decodable", file.Name())
	}
	rect := pic.Bounds()
	w := float64(rect.Max.X - rect.Min.X)
	h := float64(rect.Max.Y - rect.Min.Y)

	// If the image is landscape, then turn it to portrait orientation
	if w > h {
		newFile, newTemplateFileErr := g.newTempFile()
		if newTemplateFileErr != nil {
			return nil, 0.0, 0.0, errors.Wrap(newTemplateFileErr, "Creating temp file for image rotation")
		}

		// Rotate and save new file
		newPic := imaging.Rotate90(pic)
		if contentType == uploader.FileTypePNG {
			err := png.Encode(newFile, newPic)
			if err != nil {
				return nil, 0.0, 0.0, errors.Wrap(err, "Encountered an error rotating and encoding the png")
			}
		} else {
			err := jpeg.Encode(newFile, newPic, nil)
			if err != nil {
				return nil, 0.0, 0.0, errors.Wrap(err, "Encountered an error rotating and encoding the jpg")
			}
		}

		// The original width is now the height and vice versa.
		w, h = h, w

		// Use newFile instead of oldFile
		file = newFile

		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			return nil, 0.0, 0.0, errors.Wrap(fileCloseErr, "Encountered an error closing the file")
		}

		return newFile, w, h, nil
	}
	return file, w, h, nil
}

// PDFFromImages returns the path to tempfile PDF containing all images included
// in urls.
//
// Images will be rotated to have as little white space as possible.
//
// The files at those paths will be tempfiles that will need to be cleaned
// up by the caller.
func (g *Generator) PDFFromImages(appCtx appcontext.AppContext, images []inputFile) (string, error) {
	// These constants are based on A4 page size, which we currently default to.
	horizontalMargin := 0.0
	topMargin := 0.0
	bodyWidth := PdfPageWidth - (horizontalMargin * 2)
	bodyHeight := PdfPageHeight - (topMargin * 2)
	wToHRatio := bodyWidth / bodyHeight

	pdf := gofpdf.New(PdfOrientation, PdfUnit, PdfPageSize, PdfFontDir)
	pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	if len(images) == 0 {
		return "", errors.New("No images provided")
	}

	appCtx.Logger().Debug("generating PDF from image files", zap.Any("images", images))

	outputFile, err := g.newTempFile()
	if err != nil {
		return "", err
	}

	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
		}
	}()

	var opt gofpdf.ImageOptions
	for _, img := range images {
		pdf.AddPage()
		file, openErr := g.fs.Open(img.Path)
		if openErr != nil {
			return "", errors.Wrap(openErr, "Opening image file")
		}

		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
			}
		}()

		if img.ContentType == uploader.FileTypePNG {
			appCtx.Logger().Debug("Converting png to 8-bit")
			// gofpdf isn't able to process 16-bit PNGs, so to be safe we convert all PNGs to an 8-bit color depth
			newFile, newTemplateFileErr := g.newTempFile()
			if newTemplateFileErr != nil {
				return "", errors.Wrap(newTemplateFileErr, "Creating temp file for png conversion")
			}

			defer func() {
				if closeErr := newFile.Close(); closeErr != nil {
					appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
				}
			}()

			convertTo8BitPNGErr := convertTo8BitPNG(file, newFile)
			if convertTo8BitPNGErr != nil {
				return "", errors.Wrap(convertTo8BitPNGErr, "Converting to 8-bit png")
			}
			file = newFile
			_, fileSeekErr := file.Seek(0, io.SeekStart)
			if fileSeekErr != nil {
				return "", errors.Wrapf(fileSeekErr, "file.Seek offset: 0 whence: %d", io.SeekStart)
			}
		}

		optimizedFile, w, h, rotateErr := ReduceUnusedSpace(appCtx, file, g, img.ContentType)
		if rotateErr != nil {
			return "", errors.Wrapf(rotateErr, "Rotating image if in landscape orientation")
		}

		widthInPdf := bodyWidth
		heightInPdf := 0.0

		// Scale using the imageOptions below
		// BodyWidth should be set to 0 when the image height the proportion of the page
		// is taller than wide as compared to an A4 page.
		//
		// The opposite is true and defaulted for when the image is wider than it is tall,
		// in comparison to an A4 page.
		if float64(w/h) < wToHRatio {
			widthInPdf = 0
			heightInPdf = bodyHeight
		}

		// Rotation may have closed the file, so reopen the file before we use it.
		optimizedFile, err = g.fs.Open(optimizedFile.Name())
		if err != nil {
			return "", err
		}

		// Seek to the beginning of the file so when we register the image, it doesn't start
		// at the end of the file.
		_, fileSeekErr := optimizedFile.Seek(0, io.SeekStart)
		if fileSeekErr != nil {
			return "", errors.Wrapf(fileSeekErr, "file.Seek offset: 0 whence: %d", io.SeekStart)
		}
		// Need to register the image using an afero reader, else it uses default filesystem
		pdf.RegisterImageReader(img.Path, contentTypeToImageType[img.ContentType], optimizedFile)
		opt.ImageType = contentTypeToImageType[img.ContentType]

		pdf.ImageOptions(img.Path, horizontalMargin, topMargin, widthInPdf, heightInPdf, false, opt, 0, "")
		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			return "", errors.Wrapf(err, "error closing file: %s", optimizedFile.Name())
		}
	}

	if err = pdf.OutputAndClose(outputFile); err != nil {
		return "", errors.Wrap(err, "could not write PDF to outputfile")
	}
	return outputFile.Name(), nil
}

// MergePDFFiles Merges a slice of paths to PDF files into a single PDF
func (g *Generator) MergePDFFiles(_ appcontext.AppContext, paths []string) (afero.File, error) {
	var err error
	mergedFile, err := g.newTempFile()
	if err != nil {
		return mergedFile, err
	}

	var files []io.ReadSeeker
	for _, p := range paths {
		f, fileOpenErr := g.fs.Open(p)
		if fileOpenErr != nil {
			return mergedFile, fileOpenErr
		}
		files = append(files, f)
	}
	if err = g.pdfLib.Merge(files, mergedFile); err != nil {
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
func (g *Generator) MergeImagesToPDF(appCtx appcontext.AppContext, paths []string) (string, error) {
	// path and type for each image
	images := make([]inputFile, 0)

	for _, path := range paths {
		extension := filepath.Ext(path)[1:]
		images = append(images, inputFile{
			Path:        path,
			ContentType: strings.ToUpper(extension),
		})
	}

	return g.PDFFromImages(appCtx, images)
}

func (g *Generator) FillPDFForm(jsonData []byte, templateReader io.ReadSeeker, fileName string) (SSWWorksheet afero.File, err error) {
	var conf = g.pdfConfig
	// Change type to reader
	readJSON := strings.NewReader(string(jsonData))
	buf := new(bytes.Buffer)
	// Fills form using the template reader with json reader, outputs to byte, to be saved to afero file.
	formerr := api.FillForm(templateReader, readJSON, buf, conf)
	if formerr != nil {
		return nil, err
	}

	tempFile, err := g.newTempFileWithName(fileName) // Will use g.newTempFileWithName for proper memory usage, saves the new temp file with the fileName
	if err != nil {
		return nil, err
	}

	// copy byte[] to temp file
	_, err = io.Copy(tempFile, buf)
	if err != nil {
		return nil, errors.Wrap(err, "error io.Copy on byte[] to temp")
	}

	// Reload the file from memstore
	outputFile, err := g.FileSystem().Open(tempFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "error g.fs.Open on reload from memstore")
	}
	return outputFile, nil
}

// MergePDFFiles Merges a slice of paths to PDF files into a single PDF
func (g *Generator) MergePDFFilesByContents(_ appcontext.AppContext, fileReaders []io.ReadSeeker) (afero.File, error) {
	var err error

	// Create a merged file
	mergedFile, err := g.newTempFile()
	if err != nil {
		return nil, err
	}
	defer mergedFile.Close() // Close merged file after finishing

	// Merge files
	if err = g.pdfLib.Merge(fileReaders, mergedFile); err != nil {
		return nil, err
	}

	// Reload the merged file
	mergedFile, err = g.fs.Open(mergedFile.Name())
	if err != nil {
		return nil, err
	}

	return mergedFile, nil
}

func (g *Generator) AddWatermarks(inputFile afero.File, m map[int][]*model.Watermark) (afero.File, error) {
	buf := new(bytes.Buffer)
	err := api.AddWatermarksSliceMap(inputFile, buf, m, g.pdfConfig)
	if err != nil {
		return nil, err
	}

	tempFile, err := g.newTempFile()
	if err != nil {
		return nil, err
	}

	// copy byte[] to temp file
	_, err = io.Copy(tempFile, buf)
	if err != nil {
		return nil, errors.Wrap(err, "error io.Copy on byte[] to temp")
	}

	// Reload the file from memstore
	pdfWithWatermarks, err := g.fs.Open(tempFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "error g.fs.Open on reload from memstore")
	}

	return pdfWithWatermarks, nil
}

func (g *Generator) CreateTextWatermark(text, desc string, onTop, update bool, u types.DisplayUnit) (*model.Watermark, error) {
	return api.TextWatermark(text, desc, onTop, update, u)
}
