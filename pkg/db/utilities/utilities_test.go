package utilities_test

import (
	"context"
	"testing"

	"github.com/gobuffalo/pop/v6"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type UtilitiesSuite struct {
	*testingsuite.PopTestSuite
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
	paidWithGTCC := false
	amount := unit.Cents(10000)
	contractExpense := models.MovingExpenseReceiptTypeContractedExpense
	expenseModel := testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
		MovingExpense: models.MovingExpense{
			MovingExpenseType: &contractExpense,
			PaidWithGTCC:      &paidWithGTCC,
			Amount:            &amount,
		},
	})

	suite.MustSave(&expenseModel)
	suite.Nil(expenseModel.DeletedAt)

	err := utilities.SoftDestroy(suite.DB(), &expenseModel)
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
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ShipmentType: models.MTOShipmentTypePPM,
		},
		Stub: true,
	})
	ppmShipment := testdatagen.MakePPMShipment(suite.DB(),
		testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
	suite.MustSave(&ppmShipment)
	suite.Nil(ppmShipment.DeletedAt)

	err := utilities.SoftDestroy(suite.DB(), &ppmShipment)

	suite.NoError(err)
	suite.NotNil(ppmShipment.DeletedAt)
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
