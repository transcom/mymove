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

	err := cli.CheckTPPSFlags(v)
	if err != nil {
		return err
	}

	err = cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

// initProcessTPPSFlags initializes TPPS processing flags
func initProcessTPPSFlags(flag *pflag.FlagSet) {

	// TPPS Config
	cli.InitTPPSFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

const (
	// AVStatusCLEAN string NO_THREATS_FOUND
	AVStatusCLEAN string = "NO_THREATS_FOUND"

	// AVStatusUNKNOWN string UNKNOWN
	// Placeholder for error when scanning, actual scan results from GuardDuty are NO_THREATS_FOUND or INFECTED
	AVStatusUNKNOWN string = "UNKNOWN"

	// Default value for parameter store environment variable
	tppsSFTPFileFormatNoCustomDate string = "MILMOVE-enYYYYMMDD.csv"

	//LEGACY clam av scan status
	LegacyAVStatusCLEAN string = "CLEAN"
)

type S3API interface {
	GetObjectTagging(ctx context.Context, input *s3.GetObjectTaggingInput, optFns ...func(*s3.Options)) (*s3.GetObjectTaggingOutput, error)
	GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

var s3Client S3API

func processTPPS(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	if flags.Lookup(cli.DbEnvFlag) == nil {
		flag := pflag.CommandLine
		cli.InitDatabaseFlags(flag)
	}
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
		logger.Info(fmt.Sprintf("Duration of processTPPS task: %v", elapsedTime))
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

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil, nil)

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

	customFilePathToProcess := v.GetString(cli.ProcessTPPSCustomDateFile)
	logger.Info(fmt.Sprintf("customFilePathToProcess: %s", customFilePathToProcess))

	timezone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		logger.Error("Error loading timezone for process-tpps ECS task", zap.Error(err))
	}

	tppsFilename := ""
	if customFilePathToProcess == tppsSFTPFileFormatNoCustomDate || customFilePathToProcess == "" {
		// Process the previous day's payment file
		logger.Info("No custom filepath provided to process, processing payment file for yesterday's date.")
		yesterday := time.Now().In(timezone).AddDate(0, 0, -1)
		previousDay := yesterday.Format("20060102")
		tppsFilename = fmt.Sprintf("MILMOVE-en%s.csv", previousDay)
		previousDayFormatted := yesterday.Format("January 02, 2006")
		logger.Info(fmt.Sprintf("Starting processing of TPPS data for %s: %s", previousDayFormatted, tppsFilename))
	} else {
		// Process the custom date specified by the ProcessTPPSCustomDateFile AWS parameter store value
		logger.Info("Custom filepath provided to process")
		tppsFilename = customFilePathToProcess
		logger.Info(fmt.Sprintf("Starting transfer of TPPS data file: %s", tppsFilename))
	}

	s3Region := v.GetString(cli.AWSS3RegionFlag)
	if s3Client == nil {
		cfg, errCfg := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(s3Region),
		)
		if errCfg != nil {
			logger.Error("error loading AWS config", zap.Error(errCfg))
		}
		s3Client = s3.NewFromConfig(cfg)
	}

	tppsS3Bucket := v.GetString(cli.TPPSS3Bucket)
	tppsS3Folder := v.GetString(cli.TPPSS3Folder)
	s3Key := tppsS3Folder + tppsFilename

	avStatus, s3ObjectTags, err := getS3ObjectTags(s3Client, tppsS3Bucket, s3Key)
	if err != nil {
		logger.Error("Failed to get S3 object tags", zap.Error(err))
		return fmt.Errorf("failed to get S3 object tags: %w", err)
	}

	if avStatus == AVStatusCLEAN || avStatus == LegacyAVStatusCLEAN {
		logger.Info(fmt.Sprintf("GuardDutyMalwareScanStatus is NO_THREATS_FOUND for TPPS file: %s", tppsFilename))

		// get the S3 object, download file to /tmp dir for processing if clean
		localFilePath, err := downloadS3File(logger, s3Client, tppsS3Bucket, s3Key)
		if err != nil {
			logger.Error("Error with getting the S3 object data via GetObject", zap.Error(err))
		}

		err = tppsInvoiceProcessor.ProcessFile(appCtx, localFilePath, "")

		if err != nil {
			logger.Error("Error processing TPPS Paid Invoice Report", zap.Error(err))
		} else {
			logger.Info("Successfully processed TPPS Paid Invoice Report")
		}
	} else {
		logger.Warn("Skipping unclean file",
			zap.String("bucket", tppsS3Bucket),
			zap.String("key", s3Key),
			zap.Any("tags", s3ObjectTags))
		logger.Info("avStatus is not CLEAN, not attempting file download")
		return nil
	}

	return nil
}

func getS3ObjectTags(s3Client S3API, bucket, key string) (string, map[string]string, error) {
	tagResp, err := s3Client.GetObjectTagging(context.Background(),
		&s3.GetObjectTaggingInput{
			Bucket: &bucket,
			Key:    &key,
		})
	if err != nil {
		return AVStatusUNKNOWN, nil, err
	}

	tags := make(map[string]string)
	avStatus := AVStatusUNKNOWN

	for _, tag := range tagResp.TagSet {
		tags[*tag.Key] = *tag.Value
		if *tag.Key == "GuardDutyMalwareScanStatus" {
			avStatus = *tag.Value
		}
	}

	return avStatus, tags, nil
}

func downloadS3File(logger *zap.Logger, s3Client S3API, bucket, key string) (string, error) {
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
		return "", err
	}
	defer response.Body.Close()

	// create a temp file in /tmp directory to store the CSV from the S3 bucket
	// the /tmp directory will only exist for the duration of the task, so no cleanup is required
	tempDir := os.TempDir()
	if !isDirMutable(tempDir) {
		return "", fmt.Errorf("tmp directory (%s) is not mutable, cannot write /tmp file for TPPS processing", tempDir)
	}

	localFilePath := filepath.Join(tempDir, filepath.Base(filepath.Clean(key)))
	absoluteLocalFilePath, err := filepath.Abs(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	if !strings.HasPrefix(absoluteLocalFilePath, tempDir) {
		return "", fmt.Errorf("path traversal detected, rejecting file: %s", absoluteLocalFilePath)
	}

	file, err := os.Create(absoluteLocalFilePath)
	if err != nil {
		logger.Error("Failed to create tmp file", zap.Error(err))
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logger.Error("Failed to write S3 object to tmp file", zap.Error(err))
		return "", err
	}

	if _, err := os.Stat(absoluteLocalFilePath); err != nil {
		logger.Error("File does not exist or is inaccessible", zap.Error(err))
		return "", err
	}

	logger.Info(fmt.Sprintf("Successfully wrote S3 file contents to local file: %s", absoluteLocalFilePath))
	logFileContents(logger, absoluteLocalFilePath)

	return absoluteLocalFilePath, nil
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

	const maxPreviewSize = 5000
	utf8Content := convertToUTF8(content)

	cleanedContent := cleanLogOutput(utf8Content)

	preview := cleanedContent
	if len(cleanedContent) > maxPreviewSize {
		preview = cleanedContent[:maxPreviewSize] + "..."
	}

	logger.Info("File contents preview:",
		zap.String("filePath", filePath),
		zap.Int64("fileSize", stat.Size()),
		zap.String("content-preview", preview),
	)
}

func cleanLogOutput(input string) string {
	cleaned := strings.ReplaceAll(input, "\t", ", ")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	return cleaned
}
