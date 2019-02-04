package services

import (
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/services/paperwork/forms"
)

// FormCreator is the service object interface for CreateForm
type FormCreator interface {
	CreateForm(template forms.FormTemplate) (afero.File, error)
}
