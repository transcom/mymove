package storage

import (
	"io"
	"strings"
	"testing"
)

func TestFilesystemPresignedURL(t *testing.T) {
	fsParams := FilesystemParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	fs := NewFilesystem(fsParams)

	url, err := fs.PresignedURL("key/to/file/12345", "image/jpeg", "testimage.jpeg")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := "https://example.text/files/key/to/file/12345?contentType=image%2Fjpeg&filename=testimage.jpeg"
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}

func TestFilesystemReturnsSuccessful(t *testing.T) {
	fsParams := FilesystemParams{
		root:    "./",
		webRoot: "https://example.text/files",
	}
	filesystem := NewFilesystem(fsParams)
	if filesystem == nil {
		t.Fatal("could not create new filesystem")
	}

	storeValue := strings.NewReader("anyValue")
	_, err := filesystem.Store("anyKey", storeValue, "", nil)
	if err != nil {
		t.Fatalf("could not store in filesystem: %s", err)
	}

	retReader, err := filesystem.Fetch("anyKey")
	if err != nil {
		t.Fatalf("could not fetch from filesystem: %s", err)
	}

	err = filesystem.Delete("anyKey")
	if err != nil {
		t.Fatalf("could not delete on filesystem: %s", err)
	}

	retValue, err := io.ReadAll(retReader)
	if strings.Compare(string(retValue[:]), "anyValue") != 0 {
		t.Fatalf("could not fetch from filesystem: %s", err)
	}

	fileSystem := filesystem.FileSystem()
	if fileSystem == nil {
		t.Fatal("could not retrieve filesystem from filesystem")
	}

	tempFileSystem := filesystem.TempFileSystem()
	if tempFileSystem == nil {
		t.Fatal("could not retrieve filesystem from filesystem")
	}
}

func TestFilesystemTags(t *testing.T) {
	fsParams := FilesystemParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	fs := NewFilesystem(fsParams)

	tags, err := fs.Tags("anyKey")
	if err != nil {
		t.Fatalf("could not get tags: %s", err)
	}

	if tag, exists := tags["GuardDutyMalwareScanStatus"]; exists && strings.Compare(tag, "NO_THREATS_FOUND") != 0 {
		t.Fatal("tag 'GuardDutyMalwareScanStatus' should return NO_THREATS_FOUND")
	}
}
