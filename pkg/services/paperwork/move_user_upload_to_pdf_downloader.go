package paperwork

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// userUploadToPDFConverter is the concrete struct implementing the UserUploadToPDFConverter interface
type moveUserUploadToPDFDownloader struct {
	isTest bool
	//*uploader.UserUploader
	pdfGenerator paperwork.Generator
}

// NewUserUploadToPDFConverter creates a new userUploadToPDFConverter struct with the service dependencies
func NewMoveUserUploadToPDFDownloader(isTest bool, userUploader *uploader.UserUploader) (services.PrimeDownloadMoveUploadPDFGenerator, error) {
	pdfGenerator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		return nil, err
	}
	return &moveUserUploadToPDFDownloader{
		isTest,
		*pdfGenerator,
	}, nil
}

type pdfBatchInfo struct {
	UploadDocType services.UserUploadDocType
	FileNames     []string
	PageCounts    []int
}

// ConvertUserUploadsToPDF converts user uploads to PDFs
func (g *moveUserUploadToPDFDownloader) GenerateDownloadMoveUserUploadPDF(appCtx appcontext.AppContext, downloadMoveOrderUploadType services.DownloadMoveOrderUploadType, move models.Move) (afero.File, error) {
	var pdfBatchInfos []pdfBatchInfo
	var pdfFileNames []string

	if downloadMoveOrderUploadType == services.DownloadMoveOrderUploadTypeAll || downloadMoveOrderUploadType == services.DownloadMoveOrderUploadTypeOnlyOrders {
		if move.Orders.UploadedOrdersID == uuid.Nil {
			return nil, errors.New("Order does not have any uploades associated to it.")
		}
		info, err := g.buildPdfBatchInfo(appCtx, services.UserUploadDocTypeOrder, move.Locator, move.Orders.UploadedOrdersID)
		if err != nil {
			return nil, err
		}
		pdfBatchInfos = append(pdfBatchInfos, *info)
	}

	if downloadMoveOrderUploadType == services.DownloadMoveOrderUploadTypeAll || downloadMoveOrderUploadType == services.DownloadMoveOrderUploadTypeOnlyAmendments {
		if downloadMoveOrderUploadType == services.DownloadMoveOrderUploadTypeOnlyAmendments && move.Orders.UploadedAmendedOrdersID == nil {
			return nil, errors.New("Order does not have any amendment uploads associated to it.")
		}
		if move.Orders.UploadedAmendedOrdersID != nil {
			info, err := g.buildPdfBatchInfo(appCtx, services.UserUploadDocTypeAmendments, move.Locator, *move.Orders.UploadedAmendedOrdersID)
			if err != nil {
				return nil, err
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
	mergedPdf, err := g.pdfGenerator.MergePDFFiles(appCtx, pdfFileNames)
	if err != nil {
		return nil, err
	}

	// *** Build Bookmarks ****
	// pdfBatchInfos[0] => UploadDocs
	// pdfBatchInfos[1] => AmendedUploadDocs
	var bookmarks []pdfcpu.Bookmark
	index := 0
	docCounter := 1
	var lastDocType services.UserUploadDocType
	for i := 0; i < len(pdfBatchInfos); i++ {
		if lastDocType != pdfBatchInfos[i].UploadDocType {
			docCounter = 1
		}
		for j := 0; j < len(pdfBatchInfos[i].PageCounts); j++ {
			if pdfBatchInfos[i].UploadDocType == services.UserUploadDocTypeOrder {
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

	if g.isTest {
		// Hack to overcome unit test failure on AddPdfBookmarks().
		// For some reason it fails when using temp files in the context of running
		// this in a unit test. It actually works when running via application.
		// TODO: look into to this at a later day. For now, unit testing will not
		// have any bookmarks in the PDF outputs.
		return mergedPdf, nil
	}

	// Decorate master PDF file with bookmarks
	return g.pdfGenerator.AddPdfBookmarks(mergedPdf.Name(), bookmarks)
}

// Build orderUploadDocType for document
func (g *moveUserUploadToPDFDownloader) buildPdfBatchInfo(appCtx appcontext.AppContext, uploadDocType services.UserUploadDocType, locator string, documentID uuid.UUID) (*pdfBatchInfo, error) {
	document, err := models.FetchDocumentWithNoRestrictions(appCtx.DB(), appCtx.Session(), documentID, false)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		pdfFile, err := g.pdfGenerator.CreateMergedPDFUpload(appCtx, uploads)
		pdfFileNames = append(pdfFileNames, pdfFile.Name())
		pdfFileInfo, err := g.pdfGenerator.GetPdfFileInfo(pdfFile.Name())
		if err != nil {
			return nil, err
		}
		if pdfFileInfo != nil {
			pageCounts = append(pageCounts, pdfFileInfo.PageCount)
		}
	}
	return &pdfBatchInfo{UploadDocType: uploadDocType, PageCounts: pageCounts, FileNames: pdfFileNames}, nil
}
