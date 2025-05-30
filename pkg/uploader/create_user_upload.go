package uploader

import (
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	// weightticketparser "github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/utils"
)

const weightEstimatePages = 11

// Func taken from weight_ticket_parser.go Cannot import function from package weight_ticket_parser because it creates a cyclic import
func IsWeightEstimatorFile(appCtx appcontext.AppContext, file io.ReadCloser) (bool, error) {
	const WeightEstimatorSpreadsheetName = "CUBE SHEET-ITO-TMO-ONLY"

	excelFile, err := excelize.OpenReader(file)
	log.Println("\033[1;33m[WARNING] Something unusual happened!\033[0m")
	log.Println("excelFile.SheetCount")
	log.Println(excelFile.SheetCount)

	if err != nil {
		return false, errors.Wrap(err, "Opening excel file")
	}

	defer func() {
		// Close the spreadsheet
		if closeErr := excelFile.Close(); err != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
		}
	}()

	// Check for a spreadhsheet with the same name the Weight Estimator template uses, if we find it we assume its a Weight Estimator spreadsheet
	_, err = excelFile.GetRows(WeightEstimatorSpreadsheetName)

	if err != nil {
		return false, nil
	}

	return true, nil
}

// CreateUserUploadForDocumentWrapper wrapper/helper function to create a user upload
func CreateUserUploadForDocumentWrapper(
	appCtx appcontext.AppContext, userID uuid.UUID,
	storer storage.FileStorer, file io.ReadCloser,
	filename string,
	fileSizeLimit ByteSize,
	allowedFileTypes AllowedFileTypes,
	docID *uuid.UUID,
	uploadType models.UploadType,
) (*models.UserUpload, string, *validate.Errors, error) {

	// log.Println("\033[1;33m[WARNING] Something unusual happened!\033[0m")
	// log.Println("In: create_user_upload: CreateUserUPloadForDocumentWrapper : line 60ish looking for file size")
	// log.Println(file)
	var url string
	var userUploader *UserUploader
	var newUserUpload *models.UserUpload
	var verrs *validate.Errors
	isWeightEstimatorFile := false
	// buf := make([]byte, 35000) // create a []byte of length 35
	// // now buf[:n] holds the bytes you read

	// bytes, err := file.Read(buf)
	// file.Close()

	// log.Println("\033[1;33m[WARNING] Something unusual happened!\033[0m")
	// log.Println("create_user_upload.go: in CreateUserUploadForDocumentWrapper, bytes of file size")
	// log.Println(bytes)

	if uploadType == models.UploadTypeUSER {
		var err error
		userUploader, err = NewUserUploader(storer, fileSizeLimit)
		if err != nil {
			appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
		}
	} else if uploadType == models.UploadTypeOFFICE {
		var err error
		userUploader, err = NewOfficeUploader(storer, fileSizeLimit)
		if err != nil {
			appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
		}
	} else {
		appCtx.Logger().Fatal("could not instantiate uploader")
		return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: "could not instantiate uploader"}
	}

	timestampPattern := regexp.MustCompile(`-(\d{14})$`)

	timestamp := ""
	filenameWithoutTimestamp := ""
	if matches := timestampPattern.FindStringSubmatch(filename); len(matches) > 1 {
		timestamp = matches[1]
		filenameWithoutTimestamp = strings.TrimSuffix(filename, "-"+timestamp)
	} else {
		filenameWithoutTimestamp = filename
	}

	extension := filepath.Ext(filenameWithoutTimestamp)
	extensionLower := strings.ToLower(extension)
	if extensionLower == ".xlsx" {
		var err error
		var createErr error

		isWeightEstimatorFile, err = IsWeightEstimatorFile(appCtx, file)
		log.Println("bool val and error")
		log.Println(isWeightEstimatorFile)
		log.Println(err)
		if err != nil {
			appCtx.Logger().Fatal("Weight Estimator File determination failed", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
		}
		if !isWeightEstimatorFile {
			//return error from in here to line 430 of uploads.go in both file
			return nil, "", &validate.Errors{}, ErrWrongXlsxFormat{message: "Wrong Xlsx Format"}
		}
		// file is io.ReadCloser
		if seeker, ok := file.(io.Seeker); ok {
			// rewind back to byte 0
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				// handle error
				appCtx.Logger().Error("File seeking created an error", zap.Error(err))
			}
		} else {
			// not seekable
			appCtx.Logger().Error("The file isn't seekable")
		}
		// _, err = file.Data.Seek(0, io.SeekStart)
		// if isWeightEstimatorFile {
		pageValues, err := weightticketparser.WeightTicketComputer.ParseWeightEstimatorExcelFile(appCtx, file)

		if err != nil {
			appCtx.Logger().Error("Failed to parse Weight Estimator Xlsx File", zap.Error(err))
		}

		pdfFileName := utils.AppendTimeStampToFilename(filename)

		aFile, pdfInfo, err := weightticketparser.WeightTicketGenerator.FillWeightEstimatorPDFForm(*pageValues, pdfFileName)

		if err != nil || pdfInfo.PageCount != weightEstimatePages {
			cleanupErr := weightticketparser.WeightTicketGenerator.CleanupFile(aFile)

			if cleanupErr != nil {
				appCtx.Logger().Warn("failed to cleanup weight ticket file", zap.Error(cleanupErr), zap.String("verrs", verrs.Error()))
			}
			// return an error here
			appCtx.Logger().Error("Failed to transfer data to Weight Estimator PDF Form", zap.Error(err))
		}

		newUserUpload, verrs, createErr = userUploader.CreateUserUploadForDocument(appCtx, &docID, userID, aFile, AllowedTypesPPMDocuments)

		if verrs.HasAny() || createErr != nil {
			appCtx.Logger().Error("failed to create new user upload".zap.Error(createErr), zap.String("verrs", verrs.Error()))
			cleanupErr := CleanupFile(aFile)

			if cleanupErr != nil {
				appCtx.Logger().Warn("failed to clean up weight ticket file", zap.Error(cleanupErr), zap.String("verrs", verrs.Error()))
				return nil, "", &validate.Errors{}, ErrFile{message: cleanupErr.Error()}

			}
			// return an error here
			return nil, "", &validate.Errors{}, ErrFile{message: createErr.Error()}
		}

		newUserUpload, verrs, err := userUploader.UpdateUserXlsxUploadFilename(appCtx, newUserUpload, filename)

		if verrs.HasAny() || err != nil {
			appCtx.Logger().Error("failed to rename filename", zap.Error(createErr), zap.String("verrs", verrs.Error()))

		}
		err = weightticketparser.WeightTicketGenerator.CleanupFile(aFile)
		if err != nil {
			appCtx.Logger().Warn("failed to cean up weight ticket file", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
		}

		url, err = userUploader.PresignedURL(appCtx, newUserUpload)
		if err != nil {
			return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
		}
		return newUserUpload, url, &validate.Errors{}, err
		// }
	} else {
		aFile, err := userUploader.PrepareFileForUpload(appCtx, file, filename)
		log.Println("Prepare File for Upload error")
		log.Println(err)
		if err != nil {
			appCtx.Logger().Fatal("could not prepare file for uploader", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
		}
		log.Println("\033[1;33m[WARNING] Something unusual happened!\033[0m")
		log.Println("create_user_upload : after PrepareFile For Upload before CreateUserUploadForDocument")
		log.Println("afile size")
		info, err2 := aFile.Stat()
		if err2 != nil {
			appCtx.Logger().Error("Could not get file info", zap.Error(err2))
		}
		log.Println(info.Size())

		newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(appCtx, docID, userID, File{File: aFile}, allowedFileTypes)
		log.Println("\033[1;33m[WARNING] Something unusual happened!\033[0m")
		log.Println("In create_user_upload, after CreateUserUploadForDocument line 124 file size of aFile")
		info1, err1 := aFile.Stat()
		if err1 != nil {
			appCtx.Logger().Error("Could not get file info", zap.Error(err1))
		}
		log.Println(info1.Size())
		if verrs.HasAny() || err != nil {
			return nil, "", verrs, err
		}

		url, err := userUploader.PresignedURL(appCtx, newUserUpload)
		if err != nil {
			appCtx.Logger().Error("failed to get presigned url", zap.Error(err))
			return nil, "", &validate.Errors{}, err
		}
		return newUserUpload, url, &validate.Errors{}, err
	}
}
