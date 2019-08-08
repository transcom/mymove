package migrate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

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
	require.Equal(t, len(res), len(files))
}

func TestListFilesForS3WithInvalidClient(t *testing.T) {
	//setup
	folder := "s3://home/files"

	sut := NewFileHelper()
	res, err := sut.ListFiles(folder, nil)
	require.Equal(t, 0, len(res))
	expectedErr := errors.New(fmt.Sprintf("No s3Client provided to list files at %s", folder))
	assert.Error(t, expectedErr, err)
}