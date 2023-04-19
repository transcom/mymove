package paperwork

import (
	"bytes"
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/uploader"
)

// userUploadToPDFConverter is the concrete struct implementing the UserUploadToPDFConverter interface
type userUploadToPDFConverter struct {
	*uploader.UserUploader
}

// NewUserUploadToPDFConverter creates a new userUploadToPDFConverter struct with the service dependencies
func NewUserUploadToPDFConverter(userUploader *uploader.UserUploader) services.UserUploadToPDFConverter {
	return &userUploadToPDFConverter{
		userUploader,
	}
}

// ConvertUserUploadsToPDF converts user uploads to PDFs
func (u *userUploadToPDFConverter) ConvertUserUploadsToPDF(appCtx appcontext.AppContext, userUploads models.UserUploads) ([]*services.FileInfo, error) {
	convertedFiles := []*services.FileInfo{}

	// function to close all the PDFs we've opened if we encounter an error
	closePDFs := func() {
		for _, convertedFile := range convertedFiles {
			if err := convertedFile.PDFStream.Close(); err != nil {
				appCtx.Logger().Error("Failed to close PDF stream", zap.Error(err))
			}
		}
	}

	for _, userUpload := range userUploads {
		userUpload := userUpload

		download, downloadErr := u.UserUploader.Download(appCtx, &userUpload)

		// Normally you would want to set up a deferred function to close the download stream, but we're going to be
		// doing that later on because of how we're using the streams.

		errorMsgPrefix := fmt.Sprintf("failed to convert file %s (UserUpload ID: %d) to PDF", userUpload.Upload.Filename, userUpload.ID)

		if downloadErr != nil {
			// This should be nil, but just in case, we'll try closing it.
			if download != nil {
				if err := download.Close(); err != nil {
					appCtx.Logger().Error("Failed to close download stream", zap.Error(err))
				}
			}

			appCtx.Logger().Error(errorMsgPrefix, zap.Error(downloadErr))

			closePDFs()

			return nil, fmt.Errorf("%s: %w", errorMsgPrefix, downloadErr)
		}

		fileInfo := services.NewFileInfo(&userUpload, download)

		// No need to do anything to the file if it is already a PDF, so we'll add it to the running list and move on.
		// I had thought about running them through the conversion anyways to get them into a consistent format
		// ("PDF/A-1a"), but I'm getting an error for certain PDFs.
		// Details on error: https://dp3.atlassian.net/browse/MB-15340?focusedCommentId=25982
		if userUpload.Upload.ContentType == uploader.FileTypePDF {
			downloadContents, downloadReadErr := io.ReadAll(fileInfo.OriginalUploadStream)

			if downloadReadErr != nil {
				appCtx.Logger().Error("Failed to read download stream", zap.Error(downloadReadErr))

				closePDFs()

				return nil, fmt.Errorf("%s: failed to read download stream: %w", errorMsgPrefix, downloadReadErr)
			}

			fileInfo.PDFStream = io.NopCloser(bytes.NewReader(downloadContents))

			if err := fileInfo.OriginalUploadStream.Close(); err != nil {
				appCtx.Logger().Error("Failed to close download stream", zap.Error(err))
			}

			convertedFiles = append(convertedFiles, fileInfo)

			continue
		}

		conversionErr := convertFileToPDF(appCtx, fileInfo)

		// we'll close the downloaded file if we've converted it because we no longer need it, but if we didn't convert
		// it (and thus didn't get to this part), we don't close it because we're returning it and the caller will need
		// to close it.
		if err := fileInfo.OriginalUploadStream.Close(); err != nil {
			appCtx.Logger().Error("Failed to close download stream", zap.Error(err))
		}

		// Not setting up closing of outputPDF file since we're returning it. Caller will need to close it.

		if conversionErr != nil {
			// This should be nil, but just in case, we'll try closing it.
			if fileInfo.PDFStream != nil {
				if err := fileInfo.PDFStream.Close(); err != nil {
					appCtx.Logger().Error("Failed to close output PDF stream", zap.Error(err))
				}
			}

			closePDFs()

			return nil, conversionErr
		}

		convertedFiles = append(convertedFiles, fileInfo)
	}

	return convertedFiles, nil
}
