package migrate

import (
	"testing"

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
