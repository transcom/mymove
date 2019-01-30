package paperwork

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/mocks"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
	"testing"
)

type CreateFormSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *CreateFormSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestCreateFormSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &CreateFormSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func (suite *CreateFormSuite) GenerateGBLFormValues() models.GovBillOfLadingFormValues {
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	shipment := scenario.MakeHhgFromAwardedToAcceptedGBLReady(suite.DB(), tspUser)
	shipment.Move.Orders.TAC = models.StringPointer("NTA4")
	suite.MustSave(&shipment.Move.Orders)

	gbl, _ := models.FetchGovBillOfLadingFormValues(suite.DB(), shipment.ID)
	return gbl
}

func (suite *CreateFormSuite) TestCreateFormServiceSuccess() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	gbl := suite.GenerateGBLFormValues()
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, _ := afs.TempFile("", "ioutil-test")

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.GovBillOfLadingFormValues"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(f, nil)

	FormFiller.On("Output",
		f,
	).Return(nil)

	createFormService := CreateForm{FileStorer, FormFiller}
	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.NotNil(suite.T(), file)
	assert.Nil(suite.T(), err)
	FormFiller.AssertExpectations(suite.T())
}

func (suite *CreateFormSuite) TestCreateFormServiceFormFillerAppendPageFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	gbl := suite.GenerateGBLFormValues()

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.GovBillOfLadingFormValues"),
	).Return(errors.New("Error for FormFiller.AppendPage()")).Times(1)

	createFormService := CreateForm{FileStorer, FormFiller}
	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), file)
	FormFiller.AssertExpectations(suite.T())
}

func (suite *CreateFormSuite) TestCreateFormServiceFileStorerCreateFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	gbl := suite.GenerateGBLFormValues()

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.GovBillOfLadingFormValues"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(nil, errors.New("Error for FileStorer.Create()"))

	createFormService := CreateForm{FileStorer, FormFiller}
	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.Nil(suite.T(), file)
	assert.NotNil(suite.T(), err)
	// check error message
	//assert.Equal(suite.T(), err.msg, "Error creating a new temp file for some-form-type form.", "should be equal")
	FormFiller.AssertExpectations(suite.T())
}

func (suite *CreateFormSuite) TestCreateFormServiceFormFillerOutputFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	gbl := suite.GenerateGBLFormValues()
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, _ := afs.TempFile("", "ioutil-test")

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.GovBillOfLadingFormValues"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(f, nil)

	FormFiller.On("Output",
		f,
	).Return(errors.New("Error for FormFiller.Output()"))

	createFormService := CreateForm{FileStorer, FormFiller}
	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.Nil(suite.T(), file)
	assert.NotNil(suite.T(), err)

	errMsg := errors.Cause(err)
	fmt.Println(errMsg)
	fmt.Println(err)

	assert.Equal(suite.T(), "Failure exporting some-form-type form to file.: Error for FormFiller.Output()", errMsg.Error(), "should be equal")
	FormFiller.AssertExpectations(suite.T())
}
