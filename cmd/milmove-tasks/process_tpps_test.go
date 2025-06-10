package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
)

type MockTPPSPaidInvoiceReportProcessor struct {
	mock.Mock
}

func (m *MockTPPSPaidInvoiceReportProcessor) ProcessFile(appCtx appcontext.AppContext, syncadaPath string, text string) error {
	args := m.Called(appCtx, syncadaPath, text)
	return args.Error(0)
}

type MockS3Client struct {
	mock.Mock
}

var globalFlagSet = func() *pflag.FlagSet {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.InitDatabaseFlags(fs)
	return fs
}()

func setupTestCommand() *cobra.Command {
	mockCmd := &cobra.Command{}
	mockCmd.Flags().AddFlagSet(globalFlagSet)
	mockCmd.Flags().String(cli.ProcessTPPSCustomDateFile, "", "Custom TPPS file date")
	mockCmd.Flags().String(cli.TPPSS3Bucket, "", "S3 bucket")
	mockCmd.Flags().String(cli.TPPSS3Folder, "", "S3 folder")
	return mockCmd
}

func (m *MockS3Client) GetObjectTagging(ctx context.Context, input *s3.GetObjectTaggingInput, opts ...func(*s3.Options)) (*s3.GetObjectTaggingOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*s3.GetObjectTaggingOutput), args.Error(1)
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func runProcessTPPSWithMockS3(cmd *cobra.Command, args []string, mockS3 S3API) error {
	originalS3Client := s3Client
	defer func() { s3Client = originalS3Client }()
	s3Client = mockS3
	return processTPPS(cmd, args)
}

func TestMain(m *testing.M) {
	// make sure global flag set is fresh before running tests
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	os.Exit(m.Run())
}

func TestInitProcessTPPSFlags(t *testing.T) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	initProcessTPPSFlags(flagSet)

	dbFlag := flagSet.Lookup(cli.DbEnvFlag)
	assert.NotNil(t, dbFlag, "Expected DbEnvFlag to be initialized")

	logFlag := flagSet.Lookup(cli.LoggingLevelFlag)
	assert.NotNil(t, logFlag, "Expected LoggingLevelFlag to be initialized")

	assert.False(t, flagSet.SortFlags, "Expected flag sorting to be disabled")
}

func TestProcessTPPSSuccess(t *testing.T) {
	mockCmd := setupTestCommand()

	args := []string{
		"--process_tpps_custom_date_file=MILMOVE-en20250210.csv",
		"--tpps_s3_bucket=test-bucket",
		"--tpps_s3_folder=test-folder",
	}

	err := mockCmd.ParseFlags(args)
	assert.NoError(t, err)

	mockS3 := new(MockS3Client)
	mockS3.On("GetObjectTagging", mock.Anything, mock.Anything).
		Return(&s3.GetObjectTaggingOutput{
			TagSet: []types.Tag{
				{Key: aws.String("GuardDutyMalwareScanStatus"), Value: aws.String(AVStatusCLEAN)},
			},
		}, nil).Once()

	mockS3.On("GetObject", mock.Anything, mock.Anything).
		Return(&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader("test-data"))}, nil).Once()

	err = runProcessTPPSWithMockS3(mockCmd, args, mockS3)
	assert.NoError(t, err)
	mockS3.AssertExpectations(t)
}

func TestProcessTPPSS3Failure(t *testing.T) {
	mockCmd := setupTestCommand()

	args := []string{
		"--tpps_s3_bucket=test-bucket",
		"--tpps_s3_folder=test-folder",
		"--process_tpps_custom_date_file=MILMOVE-en20250212.csv",
	}

	err := mockCmd.ParseFlags(args)
	assert.NoError(t, err)

	mockS3 := new(MockS3Client)
	mockS3.On("GetObjectTagging", mock.Anything, mock.Anything).
		Return(&s3.GetObjectTaggingOutput{}, fmt.Errorf("S3 error")).Once()

	err = runProcessTPPSWithMockS3(mockCmd, args, mockS3)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get S3 object tags")
	mockS3.AssertExpectations(t)
}

func TestConvertToUTF8(t *testing.T) {
	utf8Data := []byte("Invoice")
	assert.Equal(t, "Invoice", convertToUTF8(utf8Data))

	utf16LEData := []byte{0xFF, 0xFE, 'I', 0, 'n', 0, 'v', 0, 'o', 0, 'i', 0, 'c', 0, 'e', 0}
	assert.Equal(t, "Invoice", convertToUTF8(utf16LEData))

	utf16BEData := []byte{0xFE, 0xFF, 0, 'I', 0, 'n', 0, 'v', 0, 'o', 0, 'i', 0, 'c', 0, 'e'}
	assert.Equal(t, "Invoice", convertToUTF8(utf16BEData))

	emptyData := []byte{}
	assert.Equal(t, "", convertToUTF8(emptyData))
}

func TestIsDirMutable(t *testing.T) {
	// using the OS temp dir, should be mutable
	assert.True(t, isDirMutable("/tmp"))

	// non-writable paths should not be mutable
	assert.False(t, isDirMutable("/root"))
}

func captureLogs(fn func(logger *zap.Logger)) string {
	var logs strings.Builder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&logs),
		zapcore.DebugLevel,
	)
	logger := zap.New(core)

	fn(logger)
	return logs.String()
}

func TestLogFileContentsFailedToOpenFile(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "write-only-file.txt")
	// 0000 = no permissions
	err := os.WriteFile(tempFile, []byte("test"), 0000)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, tempFile)
	})

	assert.Contains(t, logOutput, "Failed to open file for logging")
}

func TestLogFileContentsFailedToReadFileContents(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "unopenable-dir")
	err := os.Mkdir(tempDir, 0755)
	assert.NoError(t, err)
	defer os.Remove(tempDir)

	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, tempDir)
	})

	assert.Contains(t, logOutput, "Failed to read file contents")
}

func TestLogFileContentsFileDoesNotExistOrCantBeAccessed(t *testing.T) {
	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, "nonexistent-file.txt")
	})

	assert.Contains(t, logOutput, "File does not exist or cannot be accessed")
}

func TestLogFileContentsEmptyFile(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "empty-file.txt")
	err := os.WriteFile(tempFile, []byte(""), 0600)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, tempFile)
	})

	assert.Contains(t, logOutput, "File is empty")
}

func TestLogFileContentsShortFilePreview(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test-file.txt")
	content := "Test test test short file"
	err := os.WriteFile(tempFile, []byte(content), 0600)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, tempFile)
	})

	fmt.Println("Captured log output:", logOutput)
	rawContent, _ := os.ReadFile(tempFile)
	fmt.Println("Actual file content:", string(rawContent))

	assert.Contains(t, logOutput, "File contents preview:")
	assert.Contains(t, logOutput, content)
}

func TestLogFileContentsLongFilePreview(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "large-file.txt")
	// larger than maxPreviewSize of 5000 bytes
	longContent := strings.Repeat("M", 6000)
	err := os.WriteFile(tempFile, []byte(longContent), 0600)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	logOutput := captureLogs(func(logger *zap.Logger) {
		logFileContents(logger, tempFile)
	})

	assert.Contains(t, logOutput, "File contents preview:")
	assert.Contains(t, logOutput, "MMMMM")
	assert.Contains(t, logOutput, "...")
}
