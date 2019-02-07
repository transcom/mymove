package storage

import (
	"testing"

	"go.uber.org/zap"
)

func TestMemoryPresignedURL(t *testing.T) {
	logger := zap.NewNop()
	fsParams := MemoryParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
		logger:  logger,
	}
	fs := NewMemory(fsParams)

	url, err := fs.PresignedURL("key/to/file/12345", "image/jpeg")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := "https://example.text/files/key/to/file/12345?contentType=image%2Fjpeg"
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}
