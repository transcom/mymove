package paperwork

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
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
func (g *moveUserUploadToPDFDownloader) GenerateDownloadMoveUserUploadPDF(appCtx appcontext.AppContext, downloadMoveOrderUploadType services.MoveOrderUploadType, move models.Move, addBookmarks bool, dirName string) (mergedPdf afero.File, returnErr error) {
	var pdfBatchInfos []pdfBatchInfo
	var pdfFileNames []string
	var err error

	if downloadMoveOrderUploadType == services.MoveOrderUploadAll || downloadMoveOrderUploadType == services.MoveOrderUpload {
		if move.Orders.UploadedOrdersID == uuid.Nil {
			return nil, apperror.NewUnprocessableEntityError(fmt.Sprintf("order does not have any uploads associated to it, move.Orders.ID: %s", move.Orders.ID))
		}
		info, err := g.buildPdfBatchInfo(appCtx, services.MoveOrderUpload, move.Orders.UploadedOrdersID, dirName)
		if err != nil {
			return nil, errors.Wrap(err, "error building PDF batch information for bookmark generation for order docs")
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
				return nil, errors.Wrap(err, "error building PDF batch information for bookmark generation for amendment docs")
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

	if !addBookmarks {
		return mergedPdf, nil
	}

	// *** Build Bookmarks ****
	// pdfBatchInfos[0] => UploadDocs
	// pdfBatchInfos[1] => AmendedUploadDocs
	var bookmarks []pdfcpu.Bookmark
	index := 0
	docCounter := 1
	var lastDocType services.MoveOrderUploadType
	for i := 0; i < len(pdfBatchInfos); i++ {
		if lastDocType != pdfBatchInfos[i].UploadDocType {
			docCounter = 1
		}
		for j := 0; j < len(pdfBatchInfos[i].PageCounts); j++ {
			if pdfBatchInfos[i].UploadDocType == services.MoveOrderUpload {
				if index == 0 {
					bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: 1, PageThru: pdfBatchInfos[i].PageCounts[j], Title: fmt.Sprintf("Customer Order for MTO %s Doc #%s", move.Locator, strconv.Itoa(docCounter))})
				} else {
					bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: bookmarks[index-1].PageThru + 1, PageThru: bookmarks[index-1].PageThru + pdfBatchInfos[i].PageCounts[j], Title: fmt.Sprintf("Customer Order for MTO %s Doc #%s", move.Locator, strconv.Itoa(docCounter))})
				}
			} else {
				if index == 0 {
					bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: 1, PageThru: pdfBatchInfos[i].PageCounts[j], Title: fmt.Sprintf("Amendment #%s to Customer Order for MTO %s", strconv.Itoa(docCounter), move.Locator)})
				} else {
					bookmarks = append(bookmarks, pdfcpu.Bookmark{PageFrom: bookmarks[index-1].PageThru + 1, PageThru: bookmarks[index-1].PageThru + pdfBatchInfos[i].PageCounts[j], Title: fmt.Sprintf("Amendment #%s to Customer Order for MTO %s", strconv.Itoa(docCounter), move.Locator)})
				}
			}
			lastDocType = pdfBatchInfos[i].UploadDocType
			index++
			docCounter++
		}
	}

	defer func() {
		// if a panic occurred we set an error message that we can use to check for a recover in the calling method
		if r := recover(); r != nil {
			appCtx.Logger().Error("Panic creating move order download", zap.Error(err))
			returnErr = fmt.Errorf("panic creating move order download")
		}
	}()

	// Decorate master PDF file with bookmarks
	return g.pdfGenerator.AddPdfBookmarks(mergedPdf, bookmarks, dirName)
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
	// For each PDF gather metadata as pdfBatchInfo type used for Bookmarking.
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
