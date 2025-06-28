package storage

import (
	"io"
	"strings"
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

func TestMemoryReturnsSuccessful(t *testing.T) {
	fsParams := MemoryParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	memory := NewMemory(fsParams)
	if memory == nil {
		t.Fatal("could not create new memory")
	}

	storeValue := strings.NewReader("anyValue")
	_, err := memory.Store("anyKey", storeValue, "", nil)
	if err != nil {
		t.Fatalf("could not store in memory: %s", err)
	}

	retReader, err := memory.Fetch("anyKey")
	if err != nil {
		t.Fatalf("could not fetch from memory: %s", err)
	}

	err = memory.Delete("anyKey")
	if err != nil {
		t.Fatalf("could not delete on memory: %s", err)
	}

	retValue, err := io.ReadAll(retReader)
	if strings.Compare(string(retValue[:]), "anyValue") != 0 {
		t.Fatalf("could not fetch from memory: %s", err)
	}

	fileSystem := memory.FileSystem()
	if fileSystem == nil {
		t.Fatal("could not retrieve filesystem from memory")
	}

	tempFileSystem := memory.TempFileSystem()
	if tempFileSystem == nil {
		t.Fatal("could not retrieve filesystem from memory")
	}
}

func TestMemoryTags(t *testing.T) {
	fsParams := MemoryParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
	}
	fs := NewMemory(fsParams)

	tags, err := fs.Tags("anyKey")
	if err != nil {
		t.Fatalf("could not get tags: %s", err)
	}

	if tag, exists := tags["aGuardDutyMalwareScanStatus"]; exists && strings.Compare(tag, "NO_THREATS_FOUND") != 0 {
		t.Fatal("tag 'GuardDutyMalwareScanStatus' should return NO_THREATS_FOUND")
	}
}
