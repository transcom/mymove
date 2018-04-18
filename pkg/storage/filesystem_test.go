package storage

import (
	"testing"

	"go.uber.org/zap"
)

func TestPresignedURL(t *testing.T) {
	logger := zap.NewNop()
	fs := NewFilesystem("/home/username", "https://example.text/files", logger)

	url, err := fs.PresignedURL("key/to/file/12345", "image/jpeg")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := "https://example.text/files/key/to/file/12345?contentType=image%2Fjpeg"
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}
