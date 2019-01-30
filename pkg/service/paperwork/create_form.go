package paperwork

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/paperwork"
	"io"
)

// FileStorer is an interface for FileStorer
type FileStorer interface {
	Create(string) (afero.File, error)
}

// FormFiller is an interface for FormFiller
type FormFiller interface {
	AppendPage(io.ReadSeeker, map[string]paperwork.FieldPos, interface{}) error
	Output(io.Writer) error
}

// CreateForm is a service object to create a form with data
type CreateForm struct {
	FileStorer FileStorer
	FormFiller FormFiller
}

// CreateAssetByteReader creates a new byte reader based on the TemplateImagePath of the formLayout
func CreateAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, err
	}

	templateBuffer := bytes.NewReader(asset)
	return templateBuffer, nil
}

// Call creates a form with the given data
func (c CreateForm) Call(data interface{}, formLayout paperwork.FormLayout, fileName string, formType string) (afero.File, error) {
	// Read in bytes from Asset pkg
	templateBuffer, err := CreateAssetByteReader(formLayout.TemplateImagePath)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading template file")
	}

	// Populate form fields with data
	err = c.FormFiller.AppendPage(templateBuffer, formLayout.FieldsLayout, data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure writing %s data to form.", formType))
	}

	// Read the incoming data into a temporary afero.File for consumption
	file, err := c.FileStorer.Create(fileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error creating a new afero file for %s form.", formType))
	}

	err = c.FormFiller.Output(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure exporting %s form to file.", formType))
	}

	return file, nil
}
