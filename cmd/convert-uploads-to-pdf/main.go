package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// This sets up the command so that cobra knows what we want it to run and gives some docs for CLI usage.
// Docs: https://cobra.dev/
var root = &cobra.Command{
	Use:   "convert-uploads-to-pdf",
	Short: "Converts a PPM shipment's uploads to a PDF",
	Long: "Converts a PPM shipment's uploads to PDF format, merges them all into a single PDF, and saves the merged " +
		"PDF. This is part of the payment packet the customer needs to turn into the finance office.",
	RunE: runPPMShipmentDocumentUploadToPDFConverter,
	Args: cobra.NoArgs,
}

const PPMShipmentIDFlag string = "ppm-shipment-id"

// runPPMShipmentDocumentUploadToPDFConverter Sets up command, validates flags, and converts a PPM shipment's document
// uploads into a PDF. Note that this command is only meant to run locally, not in a deployed environment. It would
// require updates in order to run in a deployed environment.
func runPPMShipmentDocumentUploadToPDFConverter(cmd *cobra.Command, _ []string) error {
	v, viperErr := initViper(cmd)

	if viperErr != nil {
		log.Fatalf("Failed to initialize viper due to %v", viperErr)
	}

	logger, loggerErr := setUpLogger(v)

	if loggerErr != nil {
		log.Fatalf("Failed to initialize logger due to %v", loggerErr)
	}

	if err := checkConfig(v, logger); err != nil {
		log.Fatalf("Issue with config: %v", err)

		return err
	}

	awsSession := initializeAwsSession(v, logger)

	db := initializeDB(v, logger, awsSession)

	// We shouldn't get here if we can't connect, but checking just in case.
	if db == nil {
		logger.Fatal("Could not init database")
	}

	storer := storage.InitStorage(v, awsSession, logger)

	userUploader, uploaderErr := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if uploaderErr != nil {
		logger.Fatal("could not instantiate user uploader", zap.Error(uploaderErr))
	}

	appCtx := appcontext.NewAppContext(db, logger, nil)

	ppmShipmentID, err := uuid.FromString(v.GetString(PPMShipmentIDFlag))

	if err != nil {
		logger.Fatal("Could not parse PPM shipment ID", zap.Error(err))
	}

	if err := convertPPMShipmentDocumentUploadsToPDF(appCtx, userUploader, ppmShipmentID); err != nil {
		logger.Fatal("Could not convert uploads to PDF", zap.Error(err))
	}

	defer cleanUp(appCtx)

	return nil
}

// initUploadConverterFlags initializes the flags needed for the command, including the flags specific to this command
// as well as flags for the things we need, e.g. the database.
func initUploadConverterFlags() {
	flags := root.Flags()

	// This command's config
	flags.StringP(PPMShipmentIDFlag, "p", "", "The PPM shipment ID for the shipment to convert uploads for.")

	// Environment
	cli.InitEnvironmentFlags(flags)

	// DB Config
	cli.InitDatabaseFlags(flags)

	// Logging Levels
	cli.InitLoggingFlags(flags)

	// Storage
	cli.InitStorageFlags(flags)

	// sort flags for help output
	flags.SortFlags = true
}

// initViper initializes the viper config object and sets up the environment variables.
// https://github.com/spf13/viper#what-is-viper
func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()

	// https://github.com/spf13/viper#working-with-flags
	errParseFlags := cmd.ParseFlags(nil)

	if errParseFlags != nil {
		return nil, fmt.Errorf("Could not parse args: %w", errParseFlags)
	}

	errBindPFlags := v.BindPFlags(cmd.Flags())

	if errBindPFlags != nil {
		return nil, fmt.Errorf("Could not bind flags: %w", errBindPFlags)
	}

	// https://github.com/spf13/viper#working-with-environment-variables
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	return v, nil
}

// setUpLogger sets up the logger for the command.
func setUpLogger(v *viper.Viper) (*zap.Logger, error) {
	if err := cli.CheckLogging(v); err != nil {
		return nil, err
	}

	logger, _, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)

	if err != nil {
		return nil, err
	}

	// Replace the global logger with the one we just created so that we use the same logger throughout the app when we
	// don't have access to the app context (which is where we normally get the logger from).
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// checkConfig checks input flags to ensure they are valid.
func checkConfig(v *viper.Viper, logger *zap.Logger) error {
	if err := cli.CheckEnvironment(v); err != nil {
		return err
	}

	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckStorage(v); err != nil {
		return err
	}

	ppmShipmentIDString := v.GetString(PPMShipmentIDFlag)
	if ppmShipmentIDString == "" {
		return errors.New("must provide ppm-shipment-id")
	} else if ppmShipmentID, err := uuid.FromString(ppmShipmentIDString); ppmShipmentID.IsNil() || err != nil {
		return fmt.Errorf("ppm-shipment-id is not a valid UUID: %w", err)
	}

	return nil
}

// initializeAwsSession initializes the AWS session, if needed. This was mostly copied from cmd/milmove/serve.go.
func initializeAwsSession(v *viper.Viper, logger *zap.Logger) *awssession.Session {
	if !v.GetBool(cli.DbIamFlag) && !(v.GetString(cli.StorageBackendFlag) == "s3") {
		return nil
	}

	c := &aws.Config{
		Region: aws.String(v.GetString(cli.AWSRegionFlag)),
	}

	s, errorSession := awssession.NewSession(c)

	if errorSession != nil {
		logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
	}

	return s
}

// initializeDB initializes the database connection. This was mostly copied from cmd/milmove/serve.go.
func initializeDB(v *viper.Viper, logger *zap.Logger, awsSession *awssession.Session) *pop.Connection {
	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	var dbCreds *credentials.Credentials

	if v.GetBool(cli.DbIamFlag) {
		if awsSession != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly, because the
			// session conflates the environment, shared, and container metadata config within NewSession. With
			// stscreds, we use the Secure Token Service, to assume the given role (w/RDS DB connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)

			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))

			dbCreds = stscreds.NewCredentials(awsSession, dbIamRole)

			stsService := sts.New(awsSession)

			callerIdentity, callerIdentityErr := stsService.GetCallerIdentity(&sts.GetCallerIdentityInput{})

			if callerIdentityErr != nil {
				logger.Error(errors.Wrap(callerIdentityErr, "error getting aws sts caller identity").Error())
			} else {
				logger.Info(fmt.Sprintf("STS Caller Identity - Account: %s, ARN: %s, UserId: %s", *callerIdentity.Account, *callerIdentity.Arn, *callerIdentity.UserId))
			}
		}
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)

	if err != nil {
		logger.Fatal("Invalid DB Configuration", zap.Error(err))
	}

	err = cli.PingPopConnection(dbConnection, logger)

	if err != nil {
		logger.Fatal("Can't connect to the DB", zap.Error(err))
	}

	return dbConnection
}

// convertPPMShipmentDocumentUploadsToPDF converts a PPM shipment's document uploads to a PDFs, merges the PDFs into a
// single one, and then saves the merged PDF. This merged PDF can be used for the payment packet the customer needs to
// submit to the finance office.
func convertPPMShipmentDocumentUploadsToPDF(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, ppmShipmentID uuid.UUID) error {
	ppmShipment, err := ppmshipment.FindPPMShipment(appCtx, ppmShipmentID)

	if err != nil {
		return err
	}

	userUploads := gatherPPMShipmentUploads(appCtx, ppmShipment)

	pdfsToMerge, conversionErr := convertUserUploadsToPDFs(appCtx, userUploader, userUploads)

	if conversionErr != nil {
		return conversionErr
	}

	if len(pdfsToMerge) == 0 {
		return nil
	}

	// mergePDFs will close each of the individual PDF streams (each item in pdfsToMerge) once it is done.
	mergedPDF, mergeErr := mergePDFs(appCtx, pdfsToMerge)

	if mergedPDF != nil {
		defer func() {
			if err := mergedPDF.Close(); err != nil {
				appCtx.Logger().Error("Failed to close merged PDF stream", zap.Error(err))
			}
		}()
	}

	if mergeErr != nil {
		return mergeErr
	}

	if err := saveMergedPDF(appCtx, userUploader, ppmShipment, mergedPDF); err != nil {
		return err
	}

	return nil
}

// gatherPPMShipmentUploads gathers the uploads for a PPM shipment. This is mainly a helper function to help keep the
// primary function clean.
func gatherPPMShipmentUploads(_ appcontext.AppContext, ppmShipment *models.PPMShipment) models.UserUploads {
	var userUploads models.UserUploads

	for _, weightTicket := range ppmShipment.WeightTickets {
		weightTicket := weightTicket

		userUploads = append(userUploads, weightTicket.EmptyDocument.UserUploads...)
		userUploads = append(userUploads, weightTicket.FullDocument.UserUploads...)
		userUploads = append(userUploads, weightTicket.ProofOfTrailerOwnershipDocument.UserUploads...)
	}

	for _, progearWeightTicket := range ppmShipment.ProgearWeightTickets {
		progearWeightTicket := progearWeightTicket

		userUploads = append(userUploads, progearWeightTicket.Document.UserUploads...)
	}

	for _, movingExpense := range ppmShipment.MovingExpenses {
		movingExpense := movingExpense

		userUploads = append(userUploads, movingExpense.Document.UserUploads...)
	}

	return userUploads
}

// convertUserUploadsToPDFs converts a PPM shipment's document uploads to a PDFs. This goes through the process of
// downloading the actual uploads because models.UserUpload and related models.Upload contain metadata about the upload
// and its associations, but not the actual file itself. After downloading, then we convert the file to a PDF if it
// isn't already one, and then return the list of PDF streams once we finish with all the uploads.
func convertUserUploadsToPDFs(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, userUploads models.UserUploads) ([]io.ReadCloser, error) {
	pdfsToMerge := []io.ReadCloser{}

	for _, userUpload := range userUploads {
		userUpload := userUpload

		download, downloadErr := userUploader.Download(appCtx, &userUpload)

		// Normally you would want to set up a deferred function to close the download stream, but we're going to be
		// doing that later on because of how we're using the streams.

		if downloadErr != nil {
			return nil, downloadErr
		}

		// No need to do anything to the file if it is already a PDF, so we'll add it to the running list and move on.
		// I had thought about running them through the conversion anyways to get them into a consistent format
		// ("PDF/A-1a"), but I'm getting an error for certain PDFs.
		// Details on error: https://dp3.atlassian.net/browse/MB-15340?focusedCommentId=25982
		if userUpload.Upload.ContentType == uploader.FileTypePDF {
			pdfsToMerge = append(pdfsToMerge, download)

			continue
		}

		fileName := filepath.Base(userUpload.Upload.Filename)

		outputPDF, conversionErr := convertFileToPDF(appCtx, download, fileName)

		// we'll close the downloaded file if we've converted it because we're not returning it, but if we didn't
		// convert it (and thus didn't get to this part), we don't close it because we're returning it and the caller
		// will need to close it.
		if err := download.Close(); err != nil {
			appCtx.Logger().Error("Failed to close download stream", zap.Error(err))
		}

		// Not setting up closing of outputPDF file since we're returning it. Caller will need to close it.

		if conversionErr != nil {
			return nil, conversionErr
		}

		pdfsToMerge = append(pdfsToMerge, outputPDF)
	}

	return pdfsToMerge, nil
}

// convertFileToPDF converts a single file to a PDF stream. This is one of the functions that actually interacts with
// Gotenberg.
func convertFileToPDF(appCtx appcontext.AppContext, fileToConvert io.ReadCloser, fileName string) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)

	// The endpoint we'll be using accepts multipart/form-data, so we set that up here.
	writer := multipart.NewWriter(buf)

	part, formFileErr := writer.CreateFormFile("files", fileName)

	if formFileErr != nil {
		return nil, formFileErr
	}

	if _, err := io.Copy(part, fileToConvert); err != nil {
		return nil, err
	}

	// Note that this endpoint has a different field name for setting the format than the other endpoint.
	if err := writer.WriteField("nativePdfFormat", uploader.AccessiblePDFFormat); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// endpoint docs: https://gotenberg.dev/docs/modules/libreoffice#route
	req, requestErr := http.NewRequest("POST", "http://localhost:2000/forms/libreoffice/convert", buf)
	if requestErr != nil {
		return nil, requestErr
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, clientErr := http.DefaultClient.Do(req)

	if clientErr != nil {
		appCtx.Logger().Error("Failed to convert file to PDF", zap.Error(clientErr))

		return nil, clientErr
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

		return nil, fmt.Errorf("bad status | code: %d | status: %s", res.StatusCode, res.Status)
	}

	// If all is good, we'll just return the whole body, which should be the PDF stream.
	return res.Body, nil
}

// mergePDFs merges a list of PDFs into a single PDF. This is one of the functions that actually interacts with
// Gotenberg.
func mergePDFs(appCtx appcontext.AppContext, pdfsToMerge []io.ReadCloser) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)

	// The endpoint we'll be using accepts multipart/form-data, so we set that up here.
	writer := multipart.NewWriter(buf)

	for i, pdf := range pdfsToMerge {
		pdf := pdf

		defer func() {
			if err := pdf.Close(); err != nil {
				appCtx.Logger().Error("Failed to close PDF stream", zap.Error(err))
			}
		}()

		// It's important that we use a different filename (second arg) for each file. Name clashes mean that only one
		// of the files with that name actually gets converted. Tbh, I'm not sure 100% if it's that we don't even send
		// the other file, or if Gotenberg overwrites it. My guess is the later because the size of the stream we send
		// does increase if there are two files with the same name, but I'm not sure. Either way, we don't want to
		// skip any accidentally, so we'll just use a different name for each file.
		part, formFileErr := writer.CreateFormFile("files", fmt.Sprintf("file-%d.pdf", i))

		if formFileErr != nil {
			return nil, formFileErr
		}

		if _, err := io.Copy(part, pdf); err != nil {
			return nil, err
		}
	}

	// Note that this endpoint has a different field name for setting the format than the other endpoint.
	if err := writer.WriteField("pdfFormat", uploader.AccessiblePDFFormat); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// endpoint docs: https://gotenberg.dev/docs/modules/pdf-engines#merge
	req, requestErr := http.NewRequest("POST", "http://localhost:2000/forms/pdfengines/merge", buf)

	if requestErr != nil {
		return nil, requestErr
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, clientErr := http.DefaultClient.Do(req)

	if clientErr != nil {
		appCtx.Logger().Error("Failed to convert file to PDF", zap.Error(clientErr))

		return nil, clientErr
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

		return nil, fmt.Errorf("bad status | code: %d | status: %s", res.StatusCode, res.Status)
	}

	return res.Body, nil
}

// saveMergedPDF uploads the merged PDF to storage and saves the relevant DB info.
func saveMergedPDF(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, ppmShipment *models.PPMShipment, mergedPDF io.ReadCloser) error {
	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// We'll need to decide if we want to load this like this, or if we'll have a custom loader for PPMShipment that
		// includes the ServiceMember already.
		if err := txnAppCtx.DB().Load(&ppmShipment.Shipment,
			"MoveTaskOrder.Orders.ServiceMember",
		); err != nil {
			return fmt.Errorf("failed to load move task order: %w", err)
		}

		// This is here to ease re-running for now. Need to think about what we'd want to happen after the first one is
		// created. Do we delete these at any point? Would we have reason to re-generate them if we don't delete them?
		// Would a SC be able to trigger a regeneration if they change something?
		if ppmShipment.PaymentPacketID == nil {
			document := models.Document{
				ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
			}

			verrs, err := txnAppCtx.DB().ValidateAndCreate(&document)

			if verrs.HasAny() || err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}

			ppmShipment.PaymentPacketID = &document.ID
			ppmShipment.PaymentPacket = &document

			// Hacky saving of PPM Shipment. For the real implementation of this, if it's even done this way, we should
			// use the ppm shipment updater service object.
			verrs, err = txnAppCtx.DB().ValidateAndUpdate(ppmShipment)

			if verrs.HasAny() || err != nil {
				return fmt.Errorf("failed to update ppm shipment: %w", err)
			}
		}

		fileToUpload, prepErr := userUploader.PrepareFileForUpload(txnAppCtx, mergedPDF, "payment_packet.pdf")

		if prepErr != nil {
			txnAppCtx.Logger().Error("Failed to prepare file for upload", zap.Error(prepErr))

			return prepErr
		}

		newUpload, uploadVerrs, uploadErr := userUploader.CreateUserUploadForDocument(
			txnAppCtx,
			ppmShipment.PaymentPacketID,
			// Do we have a system user ID? This is meant to be the uploader ID, but it'll be the system, not a specific
			// user.
			ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
			uploader.File{File: fileToUpload},
			uploader.AllowedTypesPPMDocuments,
		)

		if uploadVerrs.HasAny() || uploadErr != nil {
			return fmt.Errorf("failed to upload file: %w", uploadErr)
		}

		ppmShipment.PaymentPacket.UserUploads = append(ppmShipment.PaymentPacket.UserUploads, *newUpload)

		// The download is just for testing purposes. We don't need to do this in the real implementation. It's so that
		// we can quickly see the final file that was created.
		download, downloadErr := userUploader.Download(txnAppCtx, newUpload)

		if downloadErr != nil {
			return fmt.Errorf("failed to download file: %w", downloadErr)
		}

		if err := writeToDisk(txnAppCtx, download, fmt.Sprintf("payment-packet-%s.pdf", time.Now().String())); err != nil {
			return fmt.Errorf("failed to write file to disk: %w", err)
		}

		return nil
	})

	if txnErr != nil {
		return txnErr
	}

	return nil
}

// writeToDisk writes a file to disk. Helper function for testing.
func writeToDisk(_ appcontext.AppContext, fileToSave io.ReadCloser, fileName string) error {
	out, createErr := os.Create(filepath.Join("tmp", fileName))

	if createErr != nil {
		return createErr
	}

	defer out.Close()

	if _, err := io.Copy(out, fileToSave); err != nil {
		return err
	}

	return nil
}

// cleanUp cleans up after the command is finished, ensuring we close the DB connection.
func cleanUp(appCtx appcontext.AppContext) {
	if appCtx.Logger() != nil {
		if r := recover(); r != nil {
			appCtx.Logger().Error(" panic", zap.Any("recover", r))
		}

		if appCtx.DB() != nil {
			appCtx.Logger().Info("closing database connections")

			if err := appCtx.DB().Close(); err != nil {
				appCtx.Logger().Error("error closing database connections", zap.Error(err))
			}

		}
	}
}

func main() {
	// We need to initialize the flags before we execute the command, otherwise it won't know about the flags.
	initUploadConverterFlags()

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
