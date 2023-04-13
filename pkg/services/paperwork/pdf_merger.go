package paperwork

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// pdfMerger is the concrete struct implementing the PDFMerger interface
type pdfMerger struct{}

// NewPDFMerger creates a new pdfMerger struct
func NewPDFMerger() services.PDFMerger {
	return &pdfMerger{}
}

// MergePDFs merges PDFs into a single PDF
func (p *pdfMerger) MergePDFs(appCtx appcontext.AppContext, pdfsToMerge []io.ReadCloser) (io.ReadCloser, error) {
	// The endpoint we'll be using accepts multipart/form-data, so we set up a multipart writer to prepare our request.
	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	// We'll use this function to close any remaining PDF streams if we encounter an error.
	closeRemainingPDFs := func(i int) {
		for j := i + 1; j < len(pdfsToMerge); j++ {
			if err := pdfsToMerge[j].Close(); err != nil {
				appCtx.Logger().Error(fmt.Sprintf("failed to close PDF stream %d", j), zap.Error(err))
			}
		}
	}

	for i, pdf := range pdfsToMerge {
		pdf := pdf

		defer func() {
			if err := pdf.Close(); err != nil {
				appCtx.Logger().Error(fmt.Sprintf("failed to close PDF stream %d", i), zap.Error(err))
			}
		}()

		// It's important that we use a different filename (second arg) for each file. Name clashes mean that only one
		// of the files with that name actually gets converted. I think Gotenberg overwrites the "duplicate" file
		// because the size of the stream we send does increase if there are two files with the same name, but I'm not
		// sure. Either way, we don't want to skip any accidentally, so we'll just use a unique name for each file.
		part, formFileErr := writer.CreateFormFile("files", fmt.Sprintf("file-%d.pdf", i))

		if formFileErr != nil {
			errMsg := fmt.Sprintf("failed to create form file for PDF stream %d", i)

			appCtx.Logger().Error(errMsg, zap.Error(formFileErr))

			closeRemainingPDFs(i)

			return nil, fmt.Errorf("%s: %w", errMsg, formFileErr)
		}

		if _, err := io.Copy(part, pdf); err != nil {
			errMsg := fmt.Sprintf("failed to copy PDF stream %d to request", i)

			appCtx.Logger().Error(errMsg, zap.Error(err))

			closeRemainingPDFs(i)

			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}
	}

	// Note that this endpoint has a different field name for setting the format than the other gotenberg endpoint.
	if err := writer.WriteField("pdfFormat", uploader.AccessiblePDFFormat); err != nil {
		return nil, err
	}

	// We need to close the writer so that the trailer is written, otherwise our request will fail.
	if err := writer.Close(); err != nil {
		errMsg := "failed to close multipart writer"

		appCtx.Logger().Error(errMsg, zap.Error(err))

		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// endpoint docs: https://gotenberg.dev/docs/modules/pdf-engines#merge
	url := formGotenbergURL("forms/pdfengines/merge")

	req, requestErr := http.NewRequest("POST", url, buf)

	if requestErr != nil {
		errMsg := "failed to create request to merge PDFs"

		appCtx.Logger().Error(errMsg, zap.Error(requestErr))

		return nil, fmt.Errorf("%s: %w", errMsg, requestErr)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, clientErr := http.DefaultClient.Do(req)

	errMsgPrefix := "failed to merge PDFs"

	if clientErr != nil {
		appCtx.Logger().Error(errMsgPrefix, zap.Error(clientErr))

		return nil, fmt.Errorf("%s: %w", errMsgPrefix, clientErr)
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(res.Body)

		var body string
		if readErr != nil {
			body = "failed to read body"
		} else {
			body = string(bodyBytes)
		}

		appCtx.Logger().Error(
			"Did not get a 200 status code when merging PDF files",
			zap.Int("status code", res.StatusCode),
			zap.String("status", res.Status),
			zap.Any("body", body),
		)

		return nil, fmt.Errorf("%s: bad status | code: %d | status: %s", errMsgPrefix, res.StatusCode, res.Status)
	}

	return res.Body, nil
}
