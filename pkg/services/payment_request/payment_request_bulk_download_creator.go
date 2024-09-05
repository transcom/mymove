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

func (p *paymentRequestBulkDownloadCreator) CreatePaymentRequestBulkDownload(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (afero.File, error) {
	errMsgPrefix := "error creating Payment Request packet"

	paymentRequest := models.PaymentRequest{}
	err := appCtx.DB().Q().Eager(
		"MoveTaskOrder",
		"ProofOfServiceDocs",
		"ProofOfServiceDocs.PrimeUploads",
		"ProofOfServiceDocs.PrimeUploads.Upload",
	).Find(&paymentRequest, paymentRequestID)
	if err != nil || len(paymentRequest.ProofOfServiceDocs) < 1 {
		return nil, fmt.Errorf("%s: %w", errMsgPrefix, err)
	}

	var primeUploads models.Uploads
	for _, serviceDoc := range paymentRequest.ProofOfServiceDocs {
		for _, upload := range serviceDoc.PrimeUploads {
			primeUploads = append(primeUploads, upload.Upload)
		}
	}

	pdfs, err := p.pdfGenerator.ConvertUploadsToPDF(appCtx, primeUploads, false)
	if err != nil {
		return nil, fmt.Errorf("%s error generating pdf", err)
	}

	pdfFile, err := p.pdfGenerator.MergePDFFiles(appCtx, pdfs)
	if err != nil {
		return nil, fmt.Errorf("%s error generating merged pdf", err)
	}

	return pdfFile, nil
}
