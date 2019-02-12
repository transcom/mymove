package paperwork

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/mocks"
	"github.com/transcom/mymove/pkg/models"
	paperworkforms "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
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

	createForm := NewCreateForm(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(gbl, "some-file-name", paperworkforms.Form1203Layout, services.GBL)
	file, err := createForm.CreateForm(template)

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

	createForm := NewCreateForm(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(gbl, "some-file-name", paperworkforms.Form1203Layout, services.GBL)
	file, err := createForm.CreateForm(template)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), file)
	serviceErrMsg := errors.Cause(err)
	assert.Equal(suite.T(), "Error for FormFiller.AppendPage()", serviceErrMsg.Error(), "should be equal")
	assert.Equal(suite.T(), "Failure writing GBL data to form.: Error for FormFiller.AppendPage()", err.Error(), "should be equal")
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

	createForm := NewCreateForm(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(gbl, "some-file-name", paperworkforms.Form1203Layout, services.GBL)
	file, err := createForm.CreateForm(template)

	assert.Nil(suite.T(), file)
	assert.NotNil(suite.T(), err)
	serviceErrMsg := errors.Cause(err)
	assert.Equal(suite.T(), "Error for FileStorer.Create()", serviceErrMsg.Error(), "should be equal")
	assert.Equal(suite.T(), "Error creating a new afero file for GBL form.: Error for FileStorer.Create()", err.Error(), "should be equal")
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

	createForm := NewCreateForm(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(gbl, "some-file-name", paperworkforms.Form1203Layout, services.GBL)
	file, err := createForm.CreateForm(template)

	assert.Nil(suite.T(), file)
	assert.NotNil(suite.T(), err)
	serviceErrMsg := errors.Cause(err)
	assert.Equal(suite.T(), "Error for FormFiller.Output()", serviceErrMsg.Error(), "should be equal")
	assert.Equal(suite.T(), "Failure exporting GBL form to file.: Error for FormFiller.Output()", err.Error(), "should be equal")
	FormFiller.AssertExpectations(suite.T())
}

func (suite *CreateFormSuite) TestCreateFormServiceCreateAssetByteReaderFailure() {
	badAssetPath := "pkg/paperwork/formtemplates/someUndefinedTemplatePath.png"
	templateBuffer, err := CreateAssetByteReader(badAssetPath)
	assert.Nil(suite.T(), templateBuffer)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Error creating asset from path. Check image path.: Asset pkg/paperwork/formtemplates/someUndefinedTemplatePath.png not found", err.Error(), "should be equal")
}
