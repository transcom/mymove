package paperwork

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

// moveUserUploadToPDFDownloader is the concrete struct implementing the services.PrimeDownloadMoveUploadPDFGenerator interface
type moveUserUploadToPDFDownloader struct {
	pdfGenerator paperwork.Generator
}

// NewMoveUserUploadToPDFDownloader creates a new userUploadToPDFDownloader struct with the service dependencies
func NewMoveUserUploadToPDFDownloader(pdfGenerator *paperwork.Generator) (services.PrimeDownloadMoveUploadPDFGenerator, error) {
	return &moveUserUploadToPDFDownloader{
		*pdfGenerator,
	}, nil
}

type pdfBatchInfo struct {
	UploadDocType services.MoveOrderUploadType
	FileNames     []string
	PageCounts    []int
}

// MoveUserUploadToPDFDownloader converts user uploads to PDFs to download
func (g *moveUserUploadToPDFDownloader) GenerateDownloadMoveUserUploadPDF(appCtx appcontext.AppContext, downloadMoveOrderUploadType services.MoveOrderUploadType, move models.Move, dirName string) (mergedPdf afero.File, returnErr error) {
	var pdfBatchInfos []pdfBatchInfo
	var pdfFileNames []string
	var err error

	if downloadMoveOrderUploadType == services.MoveOrderUploadAll || downloadMoveOrderUploadType == services.MoveOrderUpload {
		if move.Orders.UploadedOrdersID == uuid.Nil {
			return nil, apperror.NewUnprocessableEntityError(fmt.Sprintf("order does not have any uploads associated to it, move.Orders.ID: %s", move.Orders.ID))
		}
		info, err := g.buildPdfBatchInfo(appCtx, services.MoveOrderUpload, move.Orders.UploadedOrdersID, dirName)
		if err != nil {
			return nil, errors.Wrap(err, "error building PDF batch information for order docs")
		}
		pdfBatchInfos = append(pdfBatchInfos, *info)
	}

	if downloadMoveOrderUploadType == services.MoveOrderUploadAll || downloadMoveOrderUploadType == services.MoveOrderAmendmentUpload {
		if downloadMoveOrderUploadType == services.MoveOrderAmendmentUpload && move.Orders.UploadedAmendedOrdersID == nil {
			return nil, apperror.NewUnprocessableEntityError(fmt.Sprintf("order does not have any amendment uploads associated to it, move.Orders.ID: %s", move.Orders.ID))
		}
		if move.Orders.UploadedAmendedOrdersID != nil {
			info, err := g.buildPdfBatchInfo(appCtx, services.MoveOrderAmendmentUpload, *move.Orders.UploadedAmendedOrdersID, dirName)
			if err != nil {
				return nil, errors.Wrap(err, "error building PDF batch information for amendment docs")
			}
			pdfBatchInfos = append(pdfBatchInfos, *info)
		}
	}

	// Merge all pdfFileNames from pdfBatchInfos into one array for PDF merge
	for i := 0; i < len(pdfBatchInfos); i++ {
		for j := 0; j < len(pdfBatchInfos[i].FileNames); j++ {
			pdfFileNames = append(pdfFileNames, pdfBatchInfos[i].FileNames[j])
		}
	}

	// Take all of generated PDFs and merge into a single PDF.
	mergedPdf, err = g.pdfGenerator.MergePDFFiles(appCtx, pdfFileNames, dirName)
	if err != nil {
		return nil, errors.Wrap(err, "error merging PDF files into one")
	}

	defer func() {
		// if a panic occurred we set an error message that we can use to check for a recover in the calling method
		if r := recover(); r != nil {
			appCtx.Logger().Error("Panic creating move order download", zap.Error(err))
			returnErr = fmt.Errorf("panic creating move order download")
		}
	}()

	return mergedPdf, nil
}

func (g *moveUserUploadToPDFDownloader) CleanupFile(file afero.File) error {
	if file != nil {
		fs := g.pdfGenerator.FileSystem()
		exists, err := afero.Exists(fs, file.Name())

		if err != nil {
			return err
		}

		if exists {
			err := fs.Remove(file.Name())

			if err != nil {
				if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
					// File does not exist treat it as non-error:
					return nil
				}

				// Return the error if it's not a "file not found" error
				return err
			}
		}
	}

	return nil
}

// Build orderUploadDocType for document
func (g *moveUserUploadToPDFDownloader) buildPdfBatchInfo(appCtx appcontext.AppContext, uploadDocType services.MoveOrderUploadType, documentID uuid.UUID, dirName string) (*pdfBatchInfo, error) {
	document, err := models.FetchDocumentWithNoRestrictions(appCtx.DB(), appCtx.Session(), documentID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error fetching document domain by id: %s", documentID))
	}

	var pdfFileNames []string
	var pageCounts []int
	// Document has one or more uploads. Create PDF file for each.
	for _, uu := range document.UserUploads {
		// Build temp array for current userUpload
		var currentUserUpload []models.UserUpload
		currentUserUpload = append(currentUserUpload, uu)

		uploads, err := models.UploadsFromUserUploads(appCtx.DB(), currentUserUpload)
		if err != nil {
			return nil, errors.Wrap(err, "error retrieving user uploads")
		}

		pdfFile, err := g.pdfGenerator.CreateMergedPDFUpload(appCtx, uploads, dirName)
		if err != nil {
			return nil, errors.Wrap(err, "error generating a merged PDF file")
		}
		pdfFileNames = append(pdfFileNames, pdfFile.Name())
		pdfFileInfo, err := g.pdfGenerator.GetPdfFileInfo(pdfFile.Name())
		if err != nil {
			return nil, errors.Wrap(err, "error getting fileInfo from generated PDF file")
		}
		if pdfFileInfo != nil {
			pageCounts = append(pageCounts, pdfFileInfo.PageCount)
		}
	}
	return &pdfBatchInfo{UploadDocType: uploadDocType, PageCounts: pageCounts, FileNames: pdfFileNames}, nil
}
