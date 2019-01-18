package paperwork

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/paperwork"
)

// FileCreator is an interface for FileStorer
type FileCreator interface {
	Create(string) (afero.File, error)
}

// CreateForm is a service object to create a form with data
type CreateForm struct {
	FileStorer FileCreator
}

//type CreateForm struct {
//	FileStorer *afero.Afero
//}

// Call creates a form with the given data
// TODO should it be *CreateForm
func (c CreateForm) Call(data interface{}, formLayout paperwork.FormLayout, fileName string, formType string) (afero.File, error) {
	// Read in bytes from Asset pkg
	asset, err := assets.Asset(formLayout.TemplateImagePath)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading template file")
	}

	templateBuffer := bytes.NewReader(asset)
	formFiller := paperwork.NewFormFiller()

	// Populate form fields with data
	err = formFiller.AppendPage(templateBuffer, formLayout.FieldsLayout, data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure writing %s data to form.", formType))
	}

	// Read the incoming data into a temporary afero.File for consumption
	file, err := c.FileStorer.Create(fileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error creating a new afero file for %s form.", formType))
	}

	err = formFiller.Output(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure exporting %s form to file.", formType))
	}

	return file, nil
}
