package services

import (
	"bytes"
	"io"

	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	paperworkforms "github.com/transcom/mymove/pkg/paperwork"
)

// FormType defined as integer
type FormType int

// Form Types for CreateForm
const (
	GBL FormType = iota
	SSW FormType = iota
)

// String returns the string value of the Form Type
func (ft FormType) String() string {
	return [...]string{"GBL", "SSW"}[ft]
}

// FormTemplate are the struct fields defined to call CreateForm service object
type FormTemplate struct {
	Buffer       *bytes.Reader
	FieldsLayout map[string]paperworkforms.FieldPos
	FormType
	FileName string
	Data     interface{}
}

// FormCreator is the service object interface for CreateForm
type FormCreator interface {
	CreateForm(template FormTemplate) (afero.File, error)
}

// FileInfo is a struct that holds information about a file that we'll be converting to a PDF. It's helpful to have
// this as a single argument for both the work we need to do and for writing better error messages & logs.
type FileInfo struct {
	*models.UserUpload
	OriginalUploadStream io.ReadCloser
	PDFStream            io.ReadCloser
}

// NewFileInfo creates a new FileInfo struct.
func NewFileInfo(userUpload *models.UserUpload, stream io.ReadCloser) *FileInfo {
	return &FileInfo{
		UserUpload:           userUpload,
		OriginalUploadStream: stream,
	}
}

// UserUploadToPDFConverter converts user uploads to PDFs
//
//go:generate mockery --name UserUploadToPDFConverter
type UserUploadToPDFConverter interface {
	ConvertUserUploadsToPDF(appCtx appcontext.AppContext, userUploads models.UserUploads) ([]*FileInfo, error)
}

// PDFMerger merges PDFs
//
//go:generate mockery --name PDFMerger
type PDFMerger interface {
	MergePDFs(appCtx appcontext.AppContext, pdfsToMerge []io.ReadCloser) (io.ReadCloser, error)
}

// Prime move order upload to PDF generation for download
type MoveOrderUploadType int

const (
	MoveOrderUploadAll MoveOrderUploadType = iota
	MoveOrderUpload
	MoveOrderAmendmentUpload
)

//go:generate mockery --name PrimeDownloadMoveUploadPDFGenerator
type PrimeDownloadMoveUploadPDFGenerator interface {
	GenerateDownloadMoveUserUploadPDF(appCtx appcontext.AppContext, moveOrderUploadType MoveOrderUploadType, move models.Move, dirName string) (afero.File, error)
	CleanupFile(file afero.File) error
}
