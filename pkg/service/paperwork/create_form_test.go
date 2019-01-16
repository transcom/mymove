package paperwork

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/paperwork"
	"testing"
)

/*
Mocking approaches in Golang
1. stub out functions - basically rewrite fake/mock stubbed out functions

*/

type FileStorerMock struct {
	mock.Mock
}

type FormData struct {
	attributeOne string
	attributeTwo string
}

func (f *FileStorerMock) Create(name string) (afero.File, error) {
	fmt.Println("Mocked create called")
	args := f.Called(name)
	return args.Get(0).(afero.File), args.Error(1)
}

func TestCreateForm(t *testing.T) {
	fileStorer := new(FileStorerMock)
	fileStorer.On("Create").Return(errors.New("File error"))

	createFormService := CreateForm{FileStorer: fileStorer}
	createFormService.Call(FormData{attributeOne: "attribute1", attributeTwo: "attributeTwo"}, paperwork.Form1203Layout, "some-file-name", "some-form-type")
	fileStorer.AssertExpectations(t)
}
