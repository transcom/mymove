package paperwork

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/assets"
	paperworkforms "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

// Storer is an interface for FileStorer
type Storer interface {
	Create(string) (afero.File, error)
}

// Filler is an interface for FormFiller
type Filler interface {
	AppendPage(io.ReadSeeker, map[string]paperworkforms.FieldPos, interface{}) error
	Output(io.Writer) error
}

// CreateForm is a service object to create a form with data
type createForm struct {
	FileStorer Storer
	FormFiller Filler
}

// CreateAssetByteReader creates a new byte reader based on the TemplateImagePath of the formLayout
func CreateAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating asset from path. Check image path.")
	}

	templateBuffer := bytes.NewReader(asset)
	return templateBuffer, nil
}

// MakeFormTemplate creates form template with all needed parameters from handler
func MakeFormTemplate(data interface{}, fileName string, formLayout paperworkforms.FormLayout, formType services.FormType) (services.FormTemplate, error) {
	// Read in bytes from Asset pkg
	templateBuffer, err := CreateAssetByteReader(formLayout.TemplateImagePath)
	if err != nil {
		return services.FormTemplate{}, errors.Wrap(err, "Error reading template file and creating form template")
	}
	return services.FormTemplate{Buffer: templateBuffer, FieldsLayout: formLayout.FieldsLayout, FormType: formType, FileName: fileName, Data: data}, nil
}

// NewCreateForm creates a new struct with service dependencies
func NewCreateForm(FileStorer Storer, FormFiller Filler) services.FormCreator {
	return &createForm{FileStorer: FileStorer, FormFiller: FormFiller}
}

// Call creates a form with the given data
func (c createForm) CreateForm(template services.FormTemplate) (afero.File, error) {
	// Populate form fields with data
	err := c.FormFiller.AppendPage(template.Buffer, template.FieldsLayout, template.Data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure writing %s data to form.", template.FormType.String()))
	}

	// Read the incoming data into a temporary afero.File for consumption
	file, err := c.FileStorer.Create(template.FileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error creating a new afero file for %s form.", template.FormType.String()))
	}

	// Export file from form filler
	err = c.FormFiller.Output(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failure exporting %s form to file.", template.FormType.String()))
	}

	return file, nil
}
