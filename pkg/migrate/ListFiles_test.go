package migrate

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var s3folder = "s3://home/files"

func TestListFilesForLocalFilesOnly(t *testing.T) {

	//setup
	folder := "file://home/files"
	files := []string{"home/files/migration1", "home/files/migration2", "home/files/migration3"}

	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("home/files", 0755)
	for _, file := range files {
		_ = afero.WriteFile(appFS, file, []byte(file), 0644)
	}

	//sut = subject under test
	sut := NewFileHelper()
	sut.SetFileSystem(appFS)
	res, _ := sut.ListFiles(folder, nil)
	assert.Equal(t, len(res), len(files))
}

func TestListFilesForS3WithInvalidClient(t *testing.T) {
	//setup
	folder := "s3://home/files"

	// sut - subject under test
	sut := NewFileHelper()
	res, err := sut.ListFiles(folder, nil)
	require.Equal(t, 0, len(res))
	expectedErr := fmt.Errorf("No s3Client provided to list files at %s", folder)
	assert.Error(t, expectedErr, err)
}

// mock interface
type mockS3Client struct {
}

// mock function
func (m *mockS3Client) ListObjectsV2(_ context.Context, _ *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {

	// mock response/functionality
	var files []s3types.Object
	files = append(files, s3types.Object{Key: &s3folder})
	return &s3.ListObjectsV2Output{Contents: files}, nil
}

func TestListFilesForS3(t *testing.T) {
	//setup
	mockSvc := &mockS3Client{}

	sut := NewFileHelper()
	res, err := sut.ListFiles(s3folder, mockSvc)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 1)
}

func TestListFilesForBadPrefix(t *testing.T) {
	sut := NewFileHelper()
	res, err := sut.ListFiles("/home/files", nil)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 0)
}
