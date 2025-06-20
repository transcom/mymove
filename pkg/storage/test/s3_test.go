package test

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// Tests all functions of FakeS3Storage
func TestFakeS3ReturnsSuccessful(t *testing.T) {
	fakeS3 := NewFakeS3Storage(true)
	if fakeS3 == nil {
		t.Fatal("could not create new fakeS3")
	}

	storeValue := strings.NewReader("anyValue")
	_, err := fakeS3.Store("anyKey", storeValue, "", nil)
	if err != nil {
		t.Fatalf("could not store in fakeS3: %s", err)
	}

	retReader, err := fakeS3.Fetch("anyKey")
	if err != nil {
		t.Fatalf("could not fetch from fakeS3: %s", err)
	}

	err = fakeS3.Delete("anyKey")
	if err != nil {
		t.Fatalf("could not delete on fakeS3: %s", err)
	}

	retValue, err := io.ReadAll(retReader)
	if strings.Compare(string(retValue[:]), "anyValue") != 0 {
		t.Fatalf("could not fetch from fakeS3: %s", err)
	}

	fileSystem := fakeS3.FileSystem()
	if fileSystem == nil {
		t.Fatal("could not retrieve filesystem from fakeS3")
	}

	tempFileSystem := fakeS3.TempFileSystem()
	if tempFileSystem == nil {
		t.Fatal("could not retrieve filesystem from fakeS3")
	}

	tags, err := fakeS3.Tags("anyKey")
	if err != nil {
		t.Fatalf("could not fetch from fakeS3: %s", err)
	}
	if len(tags) != 2 {
		t.Fatal("return tags must have GuardDutyMalwareScanStatus key assigned for fakeS3")
	}

	presignedUrl, err := fakeS3.PresignedURL("anyKey", "anyContentType", "anyFileName")
	if err != nil {
		t.Fatal("could not retrieve presignedUrl from fakeS3")
	}

	if strings.Compare(presignedUrl, "https://example.com/dir/anyKey?response-content-disposition=attachment%3B+filename%3D%22anyFileName%22&response-content-type=anyContentType&signed=test") != 0 {
		t.Fatalf("could not retrieve proper presignedUrl from fakeS3 %s", presignedUrl)
	}
}

// Test for willSucceed false
func TestFakeS3WillNotSucceed(t *testing.T) {
	fakeS3 := NewFakeS3Storage(false)
	if fakeS3 == nil {
		t.Fatalf("could not create new fakeS3")
	}

	storeValue := strings.NewReader("anyValue")
	_, err := fakeS3.Store("anyKey", storeValue, "", nil)
	if err == nil || errors.Is(err, errors.New("failed to push")) {
		t.Fatalf("should not be able to store when willSucceed false: %s", err)
	}

	_, err = fakeS3.Fetch("anyKey")
	if err == nil || errors.Is(err, errors.New("failed to fetch file")) {
		t.Fatalf("should not find file on Fetch for willSucceed false: %s", err)
	}
}

// Tests empty tag returns empty tags on FakeS3Storage
func TestFakeS3ReturnsEmptyTags(t *testing.T) {
	fakeS3 := NewFakeS3Storage(true)
	if fakeS3 == nil {
		t.Fatal("could not create new fakeS3")
	}

	fakeS3.EmptyTags = true

	tags, err := fakeS3.Tags("anyKey")
	if err != nil {
		t.Fatalf("could not fetch from fakeS3: %s", err)
	}
	if len(tags) != 0 {
		t.Fatal("return tags must be empty for FakeS3 when EmptyTags set to true")
	}
}
