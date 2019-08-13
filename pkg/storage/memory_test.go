package storage

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"go.uber.org/zap"
)

func TestMemoryPresignedURL(t *testing.T) {
	// create in memory file system
	fsParams := MemoryParams{
		root:    "/home/username",
		webRoot: "https://example.text/files",
		logger:  zap.NewNop(),
	}
	fs := NewMemory(fsParams)

	// open the file
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fixturePath := path.Join(cwd, "fixtures", "test.pdf")
	file, err := os.Open(fixturePath)
	if err != nil {
		t.Fatal("Error opening fixture file", zap.Error(err))
	}

	// save file to inmemory file system
	checksum := "1B2M2Y8AsgTpgAmY7PhCfg=="
	_, err = fs.Store("123456", file, checksum)
	if err != nil {
		t.Fatal("Error storing file", zap.Error(err))
	}

	// get files base64 encoding
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("file seek error %s", err)
	}
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal("Error opening fixture file", zap.Error(err))
	}
	encoded := base64.StdEncoding.EncodeToString(content)

	// assert url matches base 64 encoding of file
	url, err := fs.PresignedURL("123456", "image/pdf")
	if err != nil {
		t.Fatalf("could not get presigned url: %s", err)
	}

	expected := fmt.Sprintf("data:%s;base64, %s", "image/pdf", encoded)
	if url != expected {
		t.Errorf("wrong presigned url: expected %s, got %s", expected, url)
	}
}
