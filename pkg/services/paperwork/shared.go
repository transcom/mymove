package paperwork

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// These are helper constants to make it easier to refer to the expected environment variables.
const (
	// GotenbergProtocol is the environment variable that contains the protocol to use when communicating with
	// Gotenberg.
	GotenbergProtocol string = "GOTENBERG_PROTOCOL"
	// GotenbergHost is the environment variable that contains the host to use when communicating with Gotenberg.
	GotenbergHost string = "GOTENBERG_HOST"
	// GotenbergPort is the environment variable that contains the port to use when communicating with Gotenberg.
	GotenbergPort string = "GOTENBERG_PORT"
)

// formGotenbergURL is a convenience function to build a URL for a Gotenberg endpoint.
func formGotenbergURL(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s:%s/%s",
		os.Getenv(GotenbergProtocol),
		os.Getenv(GotenbergHost),
		os.Getenv(GotenbergPort),
		endpoint,
	)
}

// convertFileToPDF converts a single file to a PDF stream. If the conversions is successful, it will set the PDF stream
// onto on input FileInfo struct.
// This is one of the functions that actually interacts with Gotenberg.
func convertFileToPDF(appCtx appcontext.AppContext, fileInfo *services.FileInfo) error {
	// The endpoint we'll be using accepts multipart/form-data, so we set up a multipart writer to prepare our request.
	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	// Using "files" as the field name because it follows the documentation examples.
	part, formFileErr := writer.CreateFormFile("files", fileInfo.UserUpload.Upload.Filename)

	// We'll use this prefix for all of our error messages to make it easier to identify which file failed.
	errorMsgPrefix := fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", fileInfo.UserUpload.Upload.Filename, fileInfo.UserUpload.ID)

	if formFileErr != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(formFileErr))

		return fmt.Errorf("%s: %w", errorMsgPrefix, formFileErr)
	}

	if _, err := io.Copy(part, fileInfo.OriginalUploadStream); err != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(err))

		return fmt.Errorf("%s: %w", errorMsgPrefix, err)
	}

	// Note that this endpoint has a different field name for setting the format than the other gotenberg endpoint.
	if err := writer.WriteField("nativePdfFormat", uploader.AccessiblePDFFormat); err != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(err))

		return fmt.Errorf("%s: %w", errorMsgPrefix, err)
	}

	// We need to close the writer so that the trailer is written, otherwise our request will fail.
	if err := writer.Close(); err != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(err))

		return fmt.Errorf("%s: %w", errorMsgPrefix, err)
	}

	// endpoint docs: https://gotenberg.dev/docs/modules/libreoffice#route
	url := formGotenbergURL("forms/libreoffice/convert")

	req, requestErr := http.NewRequest("POST", url, buf)

	if requestErr != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(requestErr))

		return fmt.Errorf("%s: %w", errorMsgPrefix, requestErr)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, clientErr := http.DefaultClient.Do(req)

	if clientErr != nil {
		appCtx.Logger().Error(errorMsgPrefix, zap.Error(clientErr))

		return fmt.Errorf("%s: %w", errorMsgPrefix, clientErr)
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
			"Did not get a 200 status code when converting to a PDF",
			zap.Int("status code", res.StatusCode),
			zap.String("status", res.Status),
			zap.Any("body", body),
		)

		return fmt.Errorf("%s: bad status | code: %d | status: %s", errorMsgPrefix, res.StatusCode, res.Status)
	}

	// If all is good, we'll just set the whole body since it should be the PDF stream.
	fileInfo.PDFStream = res.Body

	return nil
}
