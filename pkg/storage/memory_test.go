package storage

import (
	"testing"
)

func TestMemoryPresignedURL(t *testing.T) {
	fsParams := MemoryParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	fs := NewMemory(fsParams)

	url, err := fs.PresignedURL("key/to/file/12345", "image/jpeg", "testimage.jpeg")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := "https://example.text/files/key/to/file/12345?contentType=image%2Fjpeg&filename=testimage.jpeg"
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}
