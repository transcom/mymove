package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func TestLogFileContents_FailedToOpenFile(t *testing.T) {
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
