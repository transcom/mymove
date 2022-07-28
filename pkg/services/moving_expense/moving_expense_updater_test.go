package movingexpense

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite MovingExpenseSuite) TestUpdateMovingExpense() {

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.MovingExpense, hasDocumentUploads bool) *models.MovingExpense {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})

		expenseDocument := testdatagen.MakeDocument(appCtx.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		})

		if hasDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &expenseDocument.ID,
						Document:   expenseDocument,
						DeletedAt:  deletedAt,
					},
				})
			}
		}

		originalMovingExpense := models.MovingExpense{
			PPMShipmentID: ppmShipment.ID,
			Document:      expenseDocument,
			DocumentID:    expenseDocument.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalMovingExpense, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalMovingExpense)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalMovingExpense.ID)

		return &originalMovingExpense
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		notFoundMovingExpense := models.MovingExpense{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewMovingExpenseUpdater()

		updatedMovingExpense, err := updater.UpdateMovingExpense(suite.AppContextForTest(), notFoundMovingExpense, "")

		suite.Nil(updatedMovingExpense)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for MovingExpense", notFoundMovingExpense.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		appCtx := suite.AppContextForTest()

		originalMovingExpense := setupForTest(appCtx, nil, false)

		updater := NewMovingExpenseUpdater()

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *originalMovingExpense, "")

		suite.Nil(updatedMovingExpense)

		suite.Error(updateErr)
		suite.IsType(apperror.PreconditionFailedError{}, updateErr)

		suite.Equal(
			fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalMovingExpense.ID.String()),
			updateErr.Error(),
		)
	})

	suite.Run("Successfully updates", func() {
		appCtx := suite.AppContextForTest()

		originalMovingExpense := setupForTest(appCtx, nil, true)

		updater := NewMovingExpenseUpdater()
		contractedExpenseType := models.MovingExpenseReceiptTypeContractedExpense
		rejectedStatus := models.PPMDocumentStatusRejected

		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &contractedExpenseType,
			Description:       models.StringPointer("Dumpster rental"),
			PaidWithGTCC:      models.BoolPointer(true),
			MissingReceipt:    models.BoolPointer(true),
			Amount:            models.CentPointer(unit.Cents(67899)),
			Status:            &rejectedStatus,
			Reason:            models.StringPointer("Learn to recycle"),
			SITStartDate:      models.TimePointer(time.Now()),
			SITEndDate:        models.TimePointer(time.Now()),
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalMovingExpense.ID, updatedMovingExpense.ID)
		suite.Equal(originalMovingExpense.DocumentID, updatedMovingExpense.DocumentID)
		// filters out the deleted upload
		suite.Len(updatedMovingExpense.Document.UserUploads, 1)
		suite.Equal(*expectedMovingExpense.MovingExpenseType, *updatedMovingExpense.MovingExpenseType)
		suite.Equal(*expectedMovingExpense.Description, *updatedMovingExpense.Description)
		suite.Equal(*expectedMovingExpense.PaidWithGTCC, *updatedMovingExpense.PaidWithGTCC)
		suite.Equal(*expectedMovingExpense.Amount, *updatedMovingExpense.Amount)
		suite.Equal(*expectedMovingExpense.MissingReceipt, *updatedMovingExpense.MissingReceipt)
		suite.Equal(*expectedMovingExpense.Status, *updatedMovingExpense.Status)
		suite.Equal(*expectedMovingExpense.Reason, *updatedMovingExpense.Reason)
		// Only the storage type receipt should be able to set these fields, would we rather reject
		// the update outright than fail silently?
		suite.Nil(updatedMovingExpense.SITStartDate)
		suite.Nil(updatedMovingExpense.SITEndDate)

	})

	suite.Run("Successfully updates storage receipt type", func() {
		appCtx := suite.AppContextForTest()

		originalMovingExpense := setupForTest(appCtx, nil, true)

		updater := NewMovingExpenseUpdater()
		storageExpenseType := models.MovingExpenseReceiptTypeStorage
		storageStart := time.Now()
		storageEnd := storageStart.Add(7 * time.Hour * 24)

		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &storageExpenseType,
			Description:       models.StringPointer("Dolly Parton memorabilia"),
			PaidWithGTCC:      models.BoolPointer(true),
			MissingReceipt:    models.BoolPointer(true),
			Amount:            models.CentPointer(unit.Cents(67899)),
			SITStartDate:      &storageStart,
			SITEndDate:        &storageEnd,
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalMovingExpense.ID, updatedMovingExpense.ID)
		suite.Equal(originalMovingExpense.DocumentID, updatedMovingExpense.DocumentID)
		suite.Equal(*expectedMovingExpense.MovingExpenseType, *updatedMovingExpense.MovingExpenseType)
		suite.Equal(*expectedMovingExpense.Description, *updatedMovingExpense.Description)
		suite.Equal(*expectedMovingExpense.PaidWithGTCC, *updatedMovingExpense.PaidWithGTCC)
		suite.Equal(*expectedMovingExpense.Amount, *updatedMovingExpense.Amount)
		suite.Equal(*expectedMovingExpense.MissingReceipt, *updatedMovingExpense.MissingReceipt)
		suite.Equal(*expectedMovingExpense.SITStartDate, *updatedMovingExpense.SITStartDate)
		suite.Equal(*expectedMovingExpense.SITEndDate, *updatedMovingExpense.SITEndDate)
		suite.Nil(updatedMovingExpense.Status)
		suite.Nil(updatedMovingExpense.Reason)
	})

	suite.Run("Successfully clears storage dates if receipt type changes", func() {
		appCtx := suite.AppContextForTest()

		storageReceiptType := models.MovingExpenseReceiptTypeStorage
		originalMovingExpense := setupForTest(appCtx, &models.MovingExpense{
			MovingExpenseType: &storageReceiptType,
			SITStartDate:      models.TimePointer(time.Now()),
			SITEndDate:        models.TimePointer(time.Now()),
		}, true)

		updater := NewMovingExpenseUpdater()
		packingReceiptType := models.MovingExpenseReceiptTypePackingMaterials

		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &packingReceiptType,
			Description:       models.StringPointer("Foam"),
			PaidWithGTCC:      models.BoolPointer(true),
			MissingReceipt:    models.BoolPointer(true),
			Amount:            models.CentPointer(unit.Cents(67899)),
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalMovingExpense.ID, updatedMovingExpense.ID)
		suite.Equal(originalMovingExpense.DocumentID, updatedMovingExpense.DocumentID)
		suite.Equal(*expectedMovingExpense.MovingExpenseType, *updatedMovingExpense.MovingExpenseType)
		suite.Equal(*expectedMovingExpense.Description, *updatedMovingExpense.Description)
		suite.Equal(*expectedMovingExpense.PaidWithGTCC, *updatedMovingExpense.PaidWithGTCC)
		suite.Equal(*expectedMovingExpense.Amount, *updatedMovingExpense.Amount)
		suite.Equal(*expectedMovingExpense.MissingReceipt, *updatedMovingExpense.MissingReceipt)
		suite.Nil(updatedMovingExpense.SITStartDate)
		suite.Nil(updatedMovingExpense.SITEndDate)
		suite.Nil(updatedMovingExpense.Status)
		suite.Nil(updatedMovingExpense.Reason)
	})

	suite.Run("Successfully clears the reason when status is approved", func() {
		appCtx := suite.AppContextForTest()

		rejectedStatus := models.PPMDocumentStatusRejected
		originalMovingExpense := setupForTest(appCtx, &models.MovingExpense{
			Status: &rejectedStatus,
			Reason: models.StringPointer("Can't pump your own gas in New Jersey"),
		}, true)

		updater := NewMovingExpenseUpdater()
		oilExpenseType := models.MovingExpenseReceiptTypeOil

		approvedStatus := models.PPMDocumentStatusApproved
		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &oilExpenseType,
			Description:       models.StringPointer("Fuel"),
			PaidWithGTCC:      models.BoolPointer(false),
			MissingReceipt:    models.BoolPointer(false),
			Amount:            models.CentPointer(unit.Cents(67899)),
			Status:            &approvedStatus,
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalMovingExpense.ID, updatedMovingExpense.ID)
		suite.Equal(originalMovingExpense.DocumentID, updatedMovingExpense.DocumentID)
		suite.Equal(*expectedMovingExpense.MovingExpenseType, *updatedMovingExpense.MovingExpenseType)
		suite.Equal(*expectedMovingExpense.Description, *updatedMovingExpense.Description)
		suite.Equal(*expectedMovingExpense.PaidWithGTCC, *updatedMovingExpense.PaidWithGTCC)
		suite.Equal(*expectedMovingExpense.Amount, *updatedMovingExpense.Amount)
		suite.Equal(*expectedMovingExpense.MissingReceipt, *updatedMovingExpense.MissingReceipt)
		suite.Nil(updatedMovingExpense.SITStartDate)
		suite.Nil(updatedMovingExpense.SITEndDate)
		suite.Equal(*expectedMovingExpense.Status, *updatedMovingExpense.Status)
		suite.Nil(updatedMovingExpense.Reason)
	})

	suite.Run("Fails to update when files are missing", func() {
		appCtx := suite.AppContextForTest()

		originalMovingExpense := setupForTest(appCtx, nil, false)

		updater := NewMovingExpenseUpdater()
		oilExpenseType := models.MovingExpenseReceiptTypeOil

		approvedStatus := models.PPMDocumentStatusApproved
		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &oilExpenseType,
			Description:       models.StringPointer("Fuel"),
			PaidWithGTCC:      models.BoolPointer(false),
			MissingReceipt:    models.BoolPointer(false),
			Amount:            models.CentPointer(unit.Cents(67899)),
			Status:            &approvedStatus,
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updatedMovingExpense)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.ErrorContains(updateErr, "At least 1 receipt file is required")
	})
}
