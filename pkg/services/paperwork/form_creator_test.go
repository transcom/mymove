//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package paperwork

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/paperwork/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	paperworkforms "github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaperworkServiceSuite) GenerateSSWFormPage1Values() models.ShipmentSummaryWorksheetPage1Values {
	moveID, _ := uuid.NewV4()
	serviceMemberID, _ := uuid.NewV4()
	//advanceID, _ := uuid.NewV4()
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	fortGordon := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	rank := models.ServiceMemberRankE9
	moveType := models.SelectedMoveTypeHHGPPM

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:               moveID,
			SelectedMoveType: &moveType,
		},
		Order: models.Order{
			OrdersType:       ordersType,
			NewDutyStationID: fortGordon.ID,
		},
		ServiceMember: models.ServiceMember{
			ID:            serviceMemberID,
			DutyStationID: &yuma.ID,
			Rank:          &rank,
		},
	})

	netWeight := unit.Pound(10000)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:              move.ID,
			NetWeight:           &netWeight,
			HasRequestedAdvance: true,
		},
	})

	movedocuments := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm.Move.ID,
			Move:                     ppm.Move,
			PersonallyProcuredMoveID: &ppm.ID,
			Status:                   models.MoveDocumentStatusOK,
			MoveDocumentType:         "EXPENSE",
		},
		Document: models.Document{
			ServiceMemberID: serviceMemberID,
			ServiceMember:   move.Orders.ServiceMember,
		},
	}
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), movedocuments)

	session := auth.Session{
		UserID:          move.Orders.ServiceMember.UserID,
		ServiceMemberID: serviceMemberID,
		ApplicationName: auth.MilApp,
	}
	ppm.Move.Submit()
	ppm.Move.Approve()
	// This is the same PPM model as ppm, but this is the one that will be saved by SaveMoveDependencies
	ppm.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(suite.DB(), &ppm.Move)
	certificationType := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   moveID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	ssd, _ := models.FetchDataShipmentSummaryWorksheetFormData(suite.DB(), &session, moveID)
	page1Data, _, _, _ := models.FormatValuesShipmentSummaryWorksheet(ssd)
	return page1Data
}

func (suite *PaperworkServiceSuite) TestCreateFormServiceSuccess() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	ssd := suite.GenerateSSWFormPage1Values()
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, _ := afs.TempFile("", "ioutil-test")

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.ShipmentSummaryWorksheetPage1Values"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(f, nil)

	FormFiller.On("Output",
		f,
	).Return(nil)

	formCreator := NewFormCreator(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(ssd, "some-file-name", paperworkforms.ShipmentSummaryPage1Layout, services.SSW)
	file, err := formCreator.CreateForm(template)

	suite.NotNil(file)
	suite.NoError(err)
	FormFiller.AssertExpectations(suite.T())
}

func (suite *PaperworkServiceSuite) TestCreateFormServiceFormFillerAppendPageFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	ssd := suite.GenerateSSWFormPage1Values()

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.ShipmentSummaryWorksheetPage1Values"),
	).Return(errors.New("Error for FormFiller.AppendPage()")).Times(1)

	formCreator := NewFormCreator(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(ssd, "some-file-name", paperworkforms.ShipmentSummaryPage1Layout, services.SSW)
	file, err := formCreator.CreateForm(template)

	suite.NotNil(err)
	suite.Nil(file)
	serviceErrMsg := errors.Cause(err)
	suite.Equal("Error for FormFiller.AppendPage()", serviceErrMsg.Error())
	suite.Equal("Failure writing SSW data to form.: Error for FormFiller.AppendPage()", err.Error())
	FormFiller.AssertExpectations(suite.T())
}

func (suite *PaperworkServiceSuite) TestCreateFormServiceFileStorerCreateFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	ssd := suite.GenerateSSWFormPage1Values()

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.ShipmentSummaryWorksheetPage1Values"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(nil, errors.New("Error for FileStorer.Create()"))

	formCreator := NewFormCreator(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(ssd, "some-file-name", paperworkforms.ShipmentSummaryPage1Layout, services.SSW)
	file, err := formCreator.CreateForm(template)

	suite.Nil(file)
	suite.NotNil(err)
	serviceErrMsg := errors.Cause(err)
	suite.Equal("Error for FileStorer.Create()", serviceErrMsg.Error())
	suite.Equal("Error creating a new afero file for SSW form.: Error for FileStorer.Create()", err.Error())
	FormFiller.AssertExpectations(suite.T())
}

func (suite *PaperworkServiceSuite) TestCreateFormServiceFormFillerOutputFailure() {
	FileStorer := &mocks.FileStorer{}
	FormFiller := &mocks.FormFiller{}

	ssd := suite.GenerateSSWFormPage1Values()
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	f, _ := afs.TempFile("", "ioutil-test")

	FormFiller.On("AppendPage",
		mock.AnythingOfType("*bytes.Reader"),
		mock.AnythingOfType("map[string]paperwork.FieldPos"),
		mock.AnythingOfType("models.ShipmentSummaryWorksheetPage1Values"),
	).Return(nil).Times(1)

	FileStorer.On("Create",
		mock.AnythingOfType("string"),
	).Return(f, nil)

	FormFiller.On("Output",
		f,
	).Return(errors.New("Error for FormFiller.Output()"))

	formCreator := NewFormCreator(FileStorer, FormFiller)
	template, _ := MakeFormTemplate(ssd, "some-file-name", paperworkforms.ShipmentSummaryPage1Layout, services.SSW)
	file, err := formCreator.CreateForm(template)

	suite.Nil(file)
	suite.NotNil(err)
	serviceErrMsg := errors.Cause(err)
	suite.Equal("Error for FormFiller.Output()", serviceErrMsg.Error())
	suite.Equal("Failure exporting SSW form to file.: Error for FormFiller.Output()", err.Error())
	FormFiller.AssertExpectations(suite.T())
}

func (suite *PaperworkServiceSuite) TestCreateFormServiceCreateAssetByteReaderFailure() {
	badAssetPath := "pkg/paperwork/formtemplates/someUndefinedTemplatePath.png"
	templateBuffer, err := createAssetByteReader(badAssetPath)
	suite.Nil(templateBuffer)
	suite.NotNil(err)
	suite.Equal("Error creating asset from path. Check image path.: Asset pkg/paperwork/formtemplates/someUndefinedTemplatePath.png not found", err.Error())
}
