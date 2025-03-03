package paymentrequest

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestBulkDownloadCreator struct {
	pdfGenerator *paperwork.Generator
}

func NewPaymentRequestBulkDownloadCreator(pdfGenerator *paperwork.Generator) services.PaymentRequestBulkDownloadCreator {
	return &paymentRequestBulkDownloadCreator{
		pdfGenerator,
	}
}

func (p *paymentRequestBulkDownloadCreator) CreatePaymentRequestBulkDownload(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (afero.File, string, error) {
	errMsgPrefix := "error creating Payment Request packet"
	dirName := uuid.Must(uuid.NewV4()).String()
	dirPath := p.pdfGenerator.GetWorkDir() + "/" + dirName

	paymentRequest := models.PaymentRequest{}
	err := appCtx.DB().Q().Eager(
		"MoveTaskOrder",
		"ProofOfServiceDocs",
		"ProofOfServiceDocs.PrimeUploads",
		"ProofOfServiceDocs.PrimeUploads.Upload",
	).Find(&paymentRequest, paymentRequestID)
	if err != nil || len(paymentRequest.ProofOfServiceDocs) < 1 {
		return nil, dirPath, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	var primeUploads models.Uploads
	for _, serviceDoc := range paymentRequest.ProofOfServiceDocs {
		for _, upload := range serviceDoc.PrimeUploads {
			primeUploads = append(primeUploads, upload.Upload)
		}
	}

	pdfs, err := p.pdfGenerator.ConvertUploadsToPDF(appCtx, primeUploads, false, dirName)
	if err != nil {
		return nil, dirPath, fmt.Errorf("%s error generating pdf", err)
	}

	pdfFile, err := p.pdfGenerator.MergePDFFiles(appCtx, pdfs, dirName)
	if err != nil {
		return nil, dirPath, fmt.Errorf("%s error generating merged pdf", err)
	}

	return pdfFile, dirPath, nil
}

// remove all of the packet files from the temp directory associated with creating the bulk payment request
func (p *paymentRequestBulkDownloadCreator) CleanupPaymentRequestBulkDir(dirPath string) error {
	// RemoveAll does not return an error if the directory doesn't exist it will just do nothing and return nil
	return p.pdfGenerator.FileSystem().RemoveAll(dirPath)
}
