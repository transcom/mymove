package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from the command line with go run ./cmd/milmove-tasks process-tpps
func checkProcessTPPSConfig(v *viper.Viper, logger *zap.Logger) error {

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

// initProcessTPPSFlags initializes TPPS processing flags
func initProcessTPPSFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func processTPPS(cmd *cobra.Command, args []string) error {
	flag := pflag.CommandLine
	flags := cmd.Flags()
	cli.InitDatabaseFlags(flag)

	err := cmd.ParseFlags(args)
	if err != nil {
		return fmt.Errorf("could not parse args: %w", err)
	}
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return fmt.Errorf("could not bind flags: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, _, err := logging.Config(
		logging.WithEnvironment(dbEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		logger.Fatal("Failed to initialized Zap logging for process-tpps")
	}

	zap.ReplaceGlobals(logger)

	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		logger.Info(fmt.Sprintf("Duration of processTPPS task:: %v", elapsedTime))
	}()

	err = checkProcessTPPSConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	tppsInvoiceProcessor := invoice.NewTPPSPaidInvoiceReportProcessor()

	// Process TPPS paid invoice report
	// The daily run of the task will process the previous day's payment file (matching the TPPS lambda schedule of working with the previous day's file).
	// Example for running the task February 3, 2025 - we process February 2's payment file: MILMOVE-en20250202.csv

	// Should we need to process a filename from a specific day instead of the daily scheduled run:
	// 1. Find the ProcessTPPSCustomDateFile in the AWS parameter store
	// 2. Verify that it has default value of "MILMOVE-enYYYYMMDD.csv"
	// 3. Fill in the YYYYMMDD with the desired date value of the file needing processed
	// 4. Manually run the process-tpps task
	// 5. *IMPORTANT*: Set the ProcessTPPSCustomDateFile value back to default value of "MILMOVE-enYYYYMMDD.csv" in the environment that it was modified in

	s3BucketTPPSPaidInvoiceReport := v.GetString(cli.ProcessTPPSInvoiceReportPickupDirectory)
	logger.Info(fmt.Sprintf("s3BucketTPPSPaidInvoiceReport: %s\n", s3BucketTPPSPaidInvoiceReport))

	tppsS3Bucket := v.GetString(cli.TPPSS3Bucket)
	logger.Info(fmt.Sprintf("tppsS3Bucket: %s\n", tppsS3Bucket))
	tppsS3Folder := v.GetString(cli.TPPSS3Folder)
	logger.Info(fmt.Sprintf("tppsS3Folder: %s\n", tppsS3Folder))

	customFilePathToProcess := v.GetString(cli.ProcessTPPSCustomDateFile)
	logger.Info(fmt.Sprintf("customFilePathToProcess: %s\n", customFilePathToProcess))

	tppsFilename := ""

	timezone, err := time.LoadLocation("UTC")
	if err != nil {
		logger.Error("Error loading timezone for process-tpps ECS task", zap.Error(err))
	}

	logger.Info(tppsFilename)
	const tppsSFTPFileFormatNoCustomDate = "MILMOVE-enYYYYMMDD.csv"
	if customFilePathToProcess == tppsSFTPFileFormatNoCustomDate || customFilePathToProcess == "" {
		// Process the previous day's payment file
		logger.Info("No custom filepath provided to process, processing payment file for yesterday's date.")
		yesterday := time.Now().In(timezone).AddDate(0, 0, -1)
		previousDay := yesterday.Format("20060102")
		tppsFilename = fmt.Sprintf("MILMOVE-en%s.csv", previousDay)
		previousDayFormatted := yesterday.Format("January 02, 2006")
		logger.Info(fmt.Sprintf("Starting processing of TPPS data for %s: %s\n", previousDayFormatted, tppsFilename))
	} else {
		// Process the custom date specified by the ProcessTPPSCustomDateFile AWS parameter store value
		logger.Info("Custom filepath provided to process")
		tppsFilename = customFilePathToProcess
		logger.Info(fmt.Sprintf("Starting transfer of TPPS data file: %s\n", tppsFilename))
	}

	pathTPPSPaidInvoiceReport := s3BucketTPPSPaidInvoiceReport + "/" + tppsFilename
	// temporarily adding logging here to see that s3 path was found
	logger.Info(fmt.Sprintf("Entire TPPS filepath pathTPPSPaidInvoiceReport: %s", pathTPPSPaidInvoiceReport))

	var s3Client *s3.Client
	s3Region := v.GetString(cli.AWSS3RegionFlag)
	cfg, errCfg := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(s3Region),
	)
	if errCfg != nil {
		logger.Info("error loading rds aws config", zap.Error(errCfg))
	}
	s3Client = s3.NewFromConfig(cfg)

	logger.Info("Created S3 client")

	logger.Info("Getting S3 object tags to check av-status")

	s3Bucket := tppsS3Bucket
	s3Key := tppsS3Folder + tppsFilename
	logger.Info(fmt.Sprintf("s3Bucket: %s\n", s3Bucket))
	logger.Info(fmt.Sprintf("s3Key: %s\n", s3Key))

	awsBucket := aws.String("app-tpps-transfer-exp-us-gov-west-1")
	bucket := *awsBucket
	awskey := aws.String("connector-files/MILMOVE-en20250116.csv")
	key := *awskey
	avStatus, s3ObjectTags, err := getS3ObjectTags(logger, s3Client, bucket, key)
	if err != nil {
		logger.Info("Failed to get S3 object tags")
	}
	logger.Info(fmt.Sprintf("avStatus from calling getS3ObjectTags: %s\n", avStatus))

	if avStatus == "INFECTED" {
		logger.Warn("Skipping infected file",
			zap.String("bucket", bucket),
			zap.String("key", key),
			zap.Any("tags", s3ObjectTags))
		logger.Info("avStatus is INFECTED, not attempting file download")
		return nil
	}

	if avStatus == "CLEAN" {
		logger.Info("avStatus is clean, attempting file download")

		// get the S3 object, check the ClamAV results, download file to /tmp dir for processing if clean
		localFilePath, scanResult, err := downloadS3File(logger, s3Client, bucket, key)
		if err != nil {
			logger.Error("Error with getting the S3 object data via GetObject", zap.Error(err))
		}

		logger.Info(fmt.Sprintf("localFilePath from calling downloadS3File: %s\n", localFilePath))
		logger.Info(fmt.Sprintf("scanResult from calling downloadS3File: %s\n", scanResult))

		logger.Info("Scan result was clean")

		err = tppsInvoiceProcessor.ProcessFile(appCtx, localFilePath, "")

		if err != nil {
			logger.Error("Error reading TPPS Paid Invoice Report application advice responses", zap.Error(err))
		} else {
			logger.Info("Successfully processed TPPS Paid Invoice Report application advice responses")
		}
	}

	return nil
}

func getS3ObjectTags(logger *zap.Logger, s3Client *s3.Client, bucket, key string) (string, map[string]string, error) {
	tagResp, err := s3Client.GetObjectTagging(context.Background(),
		&s3.GetObjectTaggingInput{
			Bucket: &bucket,
			Key:    &key,
		})
	if err != nil {
		return "unknown", nil, err
	}

	tags := make(map[string]string)
	avStatus := "unknown"

	for _, tag := range tagResp.TagSet {
		tags[*tag.Key] = *tag.Value
		if *tag.Key == "av-status" {
			avStatus = *tag.Value
		}
	}

	return avStatus, tags, nil
}

func downloadS3File(logger *zap.Logger, s3Client *s3.Client, bucket, key string) (string, string, error) {
	response, err := s3Client.GetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

	if err != nil {
		logger.Error("Failed to get S3 object",
			zap.String("bucket", bucket),
			zap.String("key", key),
			zap.Error(err))
		return "", "", err
	}
	defer response.Body.Close()

	// create a temp file in /tmp directory to store the CSV from the S3 bucket
	// the /tmp directory will only exist for the duration of the task, so no cleanup is required
	tempDir := os.TempDir()
	if !isDirMutable(tempDir) {
		return "", "", fmt.Errorf("tmp directory (%s) is not mutable, cannot configure default pdfcpu generator settings", tempDir)
	}

	localFilePath := filepath.Join(tempDir, filepath.Base(key))
	logger.Info(fmt.Sprintf("localFilePath: %s\n", localFilePath))

	file, err := os.Create(localFilePath)
	if err != nil {
		logger.Error("Failed to create temporary file", zap.Error(err))
		return "", "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logger.Error("Failed to write S3 object to file", zap.Error(err))
		return "", "", err
	}

	content, err := os.ReadFile(localFilePath)
	if err != nil {
		logger.Error("Failed to read file contents for logging", zap.Error(err))
		return "", "", err
	}

	maxPreviewSize := 5000
	preview := string(content)
	if len(content) > maxPreviewSize {
		preview = string(content[:maxPreviewSize]) + "..."
	}

	logger.Info("File contents preview before closing:",
		zap.String("filePath", localFilePath),
		zap.String("content", preview),
	)

	// Final success message
	logger.Info("Successfully wrote to tmp file",
		zap.String("filePath", localFilePath),
	)

	logFileContents(logger, localFilePath)

	return localFilePath, "", nil
}

// convert to UTF-8 encoding
func convertToUTF8(data []byte) string {
	if len(data) >= 2 {
		if data[0] == 0xFF && data[1] == 0xFE { // UTF-16 LE
			decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()
			utf8Bytes, _, _ := transform.Bytes(decoder, data)
			return string(utf8Bytes)
		} else if data[0] == 0xFE && data[1] == 0xFF { // UTF-16 BE
			decoder := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewDecoder()
			utf8Bytes, _, _ := transform.Bytes(decoder, data)
			return string(utf8Bytes)
		}
	}
	return string(data)
}

// Identifies if a filepath directory is mutable
// This is needed in to write contents of S3 stream to
// local file so that we can open it with os.Open() in the parser
func isDirMutable(path string) bool {
	testFile := filepath.Join(path, "tmp")
	file, err := os.Create(testFile)
	if err != nil {
		log.Printf("isDirMutable: failed for %s: %v\n", path, err)
		return false
	}
	file.Close()
	os.Remove(testFile) // Cleanup the test file, it is mutable here
	return true
}

func logFileContents(logger *zap.Logger, filePath string) {
	stat, err := os.Stat(filePath)
	if err != nil {
		logger.Error("File does not exist or cannot be accessed", zap.String("filePath", filePath), zap.Error(err))
		return
	}
	if stat.Size() == 0 {
		logger.Warn("File is empty", zap.String("filePath", filePath))
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open file for logging", zap.String("filePath", filePath), zap.Error(err))
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file contents", zap.String("filePath", filePath), zap.Error(err))
		return
	}

	const maxPreviewSize = 5000 // Adjust this if needed
	// preview := string(content)
	// if len(content) > maxPreviewSize {
	// 	preview = preview[:maxPreviewSize] + "..." // Indicate truncation
	// }
	utf8Content := convertToUTF8(content)

	preview := utf8Content
	if len(utf8Content) > maxPreviewSize {
		preview = utf8Content[:maxPreviewSize] + "..."
	}

	// Log file preview
	logger.Info("File contents preview:",
		zap.String("filePath", filePath),
		zap.Int64("fileSize", stat.Size()), // Log the full file size
		zap.String("content-preview", preview),
	)
}
