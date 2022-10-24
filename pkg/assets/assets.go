package assets

import (
	"embed"
)

//go:embed notifications paperwork
var embeddedFiles embed.FS

// Asset loads and returns the asset for the given path (relative to the "pkg/assets" directory). It returns an error if
// the asset could not be found or could not be loaded.
func Asset(path string) ([]byte, error) {
	file, err := embeddedFiles.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// MustAsset is like Asset but panics when Asset would return an error. It simplifies safe initialization of global
// variables.
func MustAsset(path string) []byte {
	a, err := Asset(path)
	if err != nil {
		panic("asset: Asset(" + path + "): " + err.Error())
	}

	return a
}
