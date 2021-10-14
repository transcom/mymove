package storage

import (
	"testing"
)

func TestFilesystemPresignedURL(t *testing.T) {
	fsParams := FilesystemParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	fs := NewFilesystem(fsParams)

	url, err := fs.PresignedURL("key/to/file/12345", "image/jpeg")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := "https://example.text/files/key/to/file/12345?contentType=image%2Fjpeg"
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}
