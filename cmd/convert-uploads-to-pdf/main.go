package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
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

var root = &cobra.Command{
	Use:   "convert-uploads-to-pdf",
	Short: "Converts a PPM shipment's uploads to a PDF",
	Long:  "Converts a PPM shipment's uploads to a PDF to include in AOA and payment packets.",
	RunE:  runPPMShipmentDocumentUploadToPDFConverter,
	Args:  cobra.NoArgs,
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

	db, dbErr := cli.InitDatabase(v, nil, logger)

	if dbErr != nil {
		logger.Fatal("Could not init database", zap.Error(dbErr))
	}

	storer := storage.InitStorage(v, nil, logger)

	userUploader, uploaderErr := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if uploaderErr != nil {
		logger.Fatal("could not instantiate user uploader", zap.Error(uploaderErr))
	}

	appCtx := appcontext.NewAppContext(db, logger, nil)

	ppmShipmentID, err := uuid.FromString(v.GetString(PPMShipmentIDFlag))

	if err != nil {
		logger.Fatal("Could not parse PPM shipment ID", zap.Error(err))
	}

	if err := convertPPMShipmentDocumentUploadsToPDF(appCtx, ppmShipmentID, userUploader); err != nil {
		logger.Fatal("Could not convert uploads to PDF", zap.Error(err))
	}

	defer cleanUp(appCtx)

	return nil
}

// initUploadConverterFlags initializes the flags needed for the command, including the flags specific to this command as well as the
// database and logging flags.
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
func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()

	errParseFlags := cmd.ParseFlags(nil)

	if errParseFlags != nil {
		return nil, fmt.Errorf("Could not parse args: %w", errParseFlags)
	}

	errBindPFlags := v.BindPFlags(cmd.Flags())

	if errBindPFlags != nil {
		return nil, fmt.Errorf("Could not bind flags: %w", errBindPFlags)
	}

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

// convertPPMShipmentDocumentUploadsToPDF converts a PPM shipment's document uploads to a PDF to include in AOA and
// payment packets.
func convertPPMShipmentDocumentUploadsToPDF(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, userUploader *uploader.UserUploader) error {
	ppmShipment, err := ppmshipment.FindPPMShipment(appCtx, ppmShipmentID)

	if err != nil {
		return err
	}

	userUploads := gatherPPMShipmentUploads(appCtx, ppmShipment)

	for _, userUpload := range userUploads {
		userUpload := userUpload

		// TODO: We'll want to handle things that are already PDFs, but for now, we'll assume everything needs to be
		//  converted.

		download, downloadErr := userUploader.Download(appCtx, &userUpload)

		if downloadErr != nil {
			return downloadErr
		}

		defer func() {
			if err := download.Close(); err != nil {
				appCtx.Logger().Error("Failed to close userUpload download stream", zap.Error(err))
			}
		}()

		fileName := filepath.Base(userUpload.Upload.Filename)
		outputName := fmt.Sprintf("%s-%s.pdf", strings.TrimSuffix(fileName, filepath.Ext(fileName)), userUpload.ID.String())
		outputPath := filepath.Join("tmp", outputName)
		if conversionErr := convertFileToPDF(appCtx, download, fileName, outputPath); conversionErr != nil {
			return conversionErr
		}
	}

	return nil
}

// gatherPPMShipmentUploads gathers the uploads for a PPM shipment.
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

// convertFileToPDF converts a file to a PDF.
func convertFileToPDF(_ appcontext.AppContext, fileToConvert io.ReadCloser, fileName string, outputPath string) error {
	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	part, formFileErr := writer.CreateFormFile("file", fileName)
	if formFileErr != nil {
		return formFileErr
	}

	if _, err := io.Copy(part, fileToConvert); err != nil {
		return err
	}

	if err := writer.WriteField("nativePdfFormat", "PDF/A-1a"); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, requestErr := http.NewRequest("POST", "http://localhost:2000/forms/libreoffice/convert", buf)
	if requestErr != nil {
		return requestErr
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, clientErr := http.DefaultClient.Do(req)
	if clientErr != nil {
		return clientErr
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	defer res.Body.Close()

	out, createErr := os.Create(outputPath)
	if createErr != nil {
		return createErr
	}

	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
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

func execute() {
	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func main() {
	initUploadConverterFlags()

	execute()
}
