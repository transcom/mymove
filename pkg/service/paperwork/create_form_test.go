package paperwork

import (
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

func (suite *CreateFormSuite) TestCreateFormFileStorerCreateFail() {
	fileStorer := new(mocks.FileCreator)
	fileStorer.On("Create", "something.png").Return(nil, errors.New("File error")).Times(1)

	createFormService := CreateForm{FileStorer: fileStorer}

	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	shipment := scenario.MakeHhgFromAwardedToAcceptedGBLReady(suite.DB(), tspUser)
	shipment.Move.Orders.TAC = models.StringPointer("NTA4")
	suite.MustSave(&shipment.Move.Orders)

	gbl, _ := models.FetchGovBillOfLadingExtractor(suite.DB(), shipment.ID)

	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.Nil(suite.T(), file)
	assert.NotNil(suite.T(), err)
	//assert.Equal(suite.T(), err.msg, "Error creating a new temp file for some-form-type form.", "should be equal")

	fileStorer.AssertExpectations(suite.T())
}

func (suite *CreateFormSuite) TestCreateFormFileStorerCreateSuccess() {
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, _ := afs.TempFile("", "ioutil-test")

	fileStorer := new(mocks.FileCreator)
	fileStorer.On("Create", mock.AnythingOfType("string")).Return(f, nil).Times(2)

	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	shipment := scenario.MakeHhgFromAwardedToAcceptedGBLReady(suite.DB(), tspUser)
	shipment.Move.Orders.TAC = models.StringPointer("NTA4")
	suite.MustSave(&shipment.Move.Orders)

	gbl, _ := models.FetchGovBillOfLadingExtractor(suite.DB(), shipment.ID)

	createFormService := CreateForm{FileStorer: fileStorer}
	file, err := createFormService.Call(gbl, paperwork.Form1203Layout, "some-file-name", "some-form-type")

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), file)
	fileStorer.AssertExpectations(suite.T())
}
