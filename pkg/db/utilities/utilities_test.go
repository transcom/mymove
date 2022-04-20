package utilities_test

import (
	"context"
	"testing"

	"github.com/gobuffalo/pop/v5"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type UtilitiesSuite struct {
	testingsuite.PopTestSuite
}

func TestUtilitiesSuite(t *testing.T) {
	hs := &UtilitiesSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *UtilitiesSuite) TestSoftDestroy_NotModel() {
	arbitaryFetcher := &mocks.AdminUserFetcher{}

	err := utilities.SoftDestroy(suite.DB(), &arbitaryFetcher)

	suite.Equal("can only soft delete type model", err.Error())
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithoutAssociations() {
	//model without deleted_at with no associations
	user := testdatagen.MakeDefaultUser(suite.DB())

	err := utilities.SoftDestroy(suite.DB(), &user)

	suite.Equal("this model does not have deleted_at field", err.Error())
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithoutAssociations() {
	//model with deleted_at with no associations
	expenseDocumentModel := testdatagen.MakeMovingExpenseDocument(suite.DB(), testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MovingExpenseType:    models.MovingExpenseTypeCONTRACTEDEXPENSE,
			PaymentMethod:        "GTCC",
			RequestedAmountCents: unit.Cents(10000),
		},
	})

	suite.MustSave(&expenseDocumentModel)
	suite.Nil(expenseDocumentModel.DeletedAt)

	err := utilities.SoftDestroy(suite.DB(), &expenseDocumentModel)
	suite.NoError(err)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithoutDeletedAtWithAssociations() {
	// model without deleted_at with associations
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	suite.MustSave(&serviceMember)

	err := utilities.SoftDestroy(suite.DB(), &serviceMember)
	suite.Equal("this model does not have deleted_at field", err.Error())
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithHasOneAssociations() {
	// model with deleted_at with "has one" associations
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusPAYMENTREQUESTED,
		},
	})
	move := ppm.Move
	moveDoc := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
				Status:                   models.MoveDocumentStatusOK,
			},
		})
	suite.MustSave(&moveDoc)
	suite.Nil(moveDoc.DeletedAt)

	vehicleNickname := "My Car"
	emptyWeight := unit.Pound(1000)
	fullWeight := unit.Pound(2500)
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDoc.ID,
		MoveDocument:             moveDoc,
		EmptyWeight:              &emptyWeight,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	suite.MustSave(&weightTicketSetDocument)
	suite.Nil(weightTicketSetDocument.DeletedAt)

	err := utilities.SoftDestroy(suite.DB(), &moveDoc)

	suite.NoError(err)
	suite.NotNil(moveDoc.DeletedAt)
	suite.NotNil(moveDoc.WeightTicketSetDocument.DeletedAt)
}

func (suite *UtilitiesSuite) TestSoftDestroy_ModelWithDeletedAtWithHasManyAssociations() {
	// model with deleted_at with "has many" associations
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
		Document: models.Document{
			ServiceMemberID: serviceMember.ID,
			ServiceMember:   serviceMember,
		},
	})
	suite.MustSave(&document)
	suite.Nil(document.DeletedAt)

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}
	suite.MustSave(&upload)
	userUpload1 := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
		UploadID:   upload.ID,
		Upload:     upload,
	}
	suite.MustSave(&userUpload1)
	upload2 := models.Upload{
		Filename:    "test2.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}
	suite.MustSave(&upload2)
	userUpload2 := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
		UploadID:   upload2.ID,
		Upload:     upload2,
	}
	suite.MustSave(&userUpload2)
	suite.Nil(upload.DeletedAt)
	suite.Nil(upload2.DeletedAt)
	suite.Nil(userUpload1.DeletedAt)
	suite.Nil(userUpload2.DeletedAt)

	err := utilities.SoftDestroy(suite.DB(), &document)

	suite.NoError(err)
	suite.NotNil(document.DeletedAt)
	suite.NotNil(document.UserUploads[0].DeletedAt)
	suite.NotNil(document.UserUploads[1].DeletedAt)
}

func (suite *UtilitiesSuite) TestExcludeDeletedScope() {
	suite.Run("successfully adds scope with no model args", func() {
		query := suite.DB().Q().Scope(utilities.ExcludeDeletedScope())

		sqlString, _ := query.ToSQL(pop.NewModel(models.MTOShipment{}, context.Background()))
		suite.Contains(sqlString, "WHERE deleted_at IS NULL")
	})

	suite.Run("successfully adds scope with default reflection table name", func() {
		query := suite.DB().Q().Scope(utilities.ExcludeDeletedScope(models.UserUpload{}))

		sqlString, _ := query.ToSQL(pop.NewModel(models.UserUpload{}, context.Background()))
		suite.Contains(sqlString, "WHERE user_uploads.deleted_at IS NULL")
	})

	suite.Run("successfully adds scope with overridden table name", func() {
		query := suite.DB().Q().Scope(utilities.ExcludeDeletedScope(models.Reimbursement{}))

		sqlString, _ := query.ToSQL(pop.NewModel(models.Reimbursement{}, context.Background()))
		suite.Contains(sqlString, "WHERE archived_reimbursements.deleted_at IS NULL")
	})

	suite.Run("successfully adds scope when there are multiple models", func() {
		query := suite.DB().Q().Scope(utilities.ExcludeDeletedScope(models.MTOShipment{}, models.UserUpload{}))

		sqlString, _ := query.ToSQL(pop.NewModel(models.MTOShipment{}, context.Background()))
		suite.Contains(sqlString, "WHERE mto_shipments.deleted_at IS NULL AND user_uploads.deleted_at IS NULL")
	})
}
