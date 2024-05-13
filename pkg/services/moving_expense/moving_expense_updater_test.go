package movingexpense

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MovingExpenseSuite) TestUpdateMovingExpense() {

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.MovingExpense, hasDocumentUploads bool) *models.MovingExpense {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)

		expenseDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		if hasDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    expenseDocument,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{
							DeletedAt: deletedAt,
						},
					},
				}, nil)
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

		updater := NewCustomerMovingExpenseUpdater()

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		updatedMovingExpense, err := updater.UpdateMovingExpense(appCtx, notFoundMovingExpense, "")

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
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMovingExpense := setupForTest(appCtx, nil, false)

		updater := NewCustomerMovingExpenseUpdater()

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *originalMovingExpense, "")

		suite.Nil(updatedMovingExpense)

		suite.Error(updateErr)
		suite.IsType(apperror.PreconditionFailedError{}, updateErr)

		suite.Equal(
			fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalMovingExpense.ID.String()),
			updateErr.Error(),
		)
	})

	suite.Run("Returns not found if user is unauthorized", func() {
		setupAppCtx := suite.AppContextWithSessionForTest(&auth.Session{})
		originalMovingExpense := setupForTest(setupAppCtx, nil, false)

		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: unauthorizedUser.ID,
		})

		updater := NewCustomerMovingExpenseUpdater()

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *originalMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updatedMovingExpense)

		suite.Error(updateErr)
		suite.IsType(apperror.NotFoundError{}, updateErr)
	})

	suite.Run("Successfully updates as a MilMove customer", func() {
		// It's obnoxious, but: we can't use the setupForTest function here,
		// since we need to get the service member ID for the AppContext.
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		// Code ported from `setupForTest`

		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)

		expenseDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		for i := 0; i < 2; i++ {
			var deletedAt *time.Time
			if i == 1 {
				deletedAt = models.TimePointer(time.Now())
			}
			factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    expenseDocument,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{
						DeletedAt: deletedAt,
					},
				},
			}, nil)
		}

		originalMovingExpense := models.MovingExpense{
			PPMShipmentID: ppmShipment.ID,
			Document:      expenseDocument,
			DocumentID:    expenseDocument.ID,
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalMovingExpense)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalMovingExpense.ID)

		// Actual test starts here

		updater := NewCustomerMovingExpenseUpdater()
		contractedExpenseType := models.MovingExpenseReceiptTypeContractedExpense

		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &contractedExpenseType,
			Description:       models.StringPointer("Dumpster rental"),
			PaidWithGTCC:      models.BoolPointer(true),
			MissingReceipt:    models.BoolPointer(true),
			Amount:            models.CentPointer(unit.Cents(67899)),
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
		suite.Equal(*expectedMovingExpense.Amount, *updatedMovingExpense.SubmittedAmount)
		// Only the storage type receipt should be able to set these fields, would we rather reject
		// the update outright than fail silently?
		suite.Nil(updatedMovingExpense.SITStartDate)
		suite.Nil(updatedMovingExpense.SITEndDate)
	})

	suite.Run("Successfully updates as an office user", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMovingExpense := setupForTest(appCtx, nil, true)

		updater := NewOfficeMovingExpenseUpdater()
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
		// Only the storage type receipt should be able to set these fields
		suite.Nil(updatedMovingExpense.SITStartDate)
		suite.Nil(updatedMovingExpense.SITEndDate)
		// Office user updates should not update SubmittedAmount
		suite.Nil(updatedMovingExpense.SubmittedAmount)
	})

	suite.Run("Successfully updates storage receipt type", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMovingExpense := setupForTest(appCtx, nil, true)

		updater := NewCustomerMovingExpenseUpdater()
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
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		storageReceiptType := models.MovingExpenseReceiptTypeStorage
		originalMovingExpense := setupForTest(appCtx, &models.MovingExpense{
			MovingExpenseType: &storageReceiptType,
			SITStartDate:      models.TimePointer(time.Now()),
			SITEndDate:        models.TimePointer(time.Now()),
		}, true)

		updater := NewCustomerMovingExpenseUpdater()
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
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		rejectedStatus := models.PPMDocumentStatusRejected
		originalMovingExpense := setupForTest(appCtx, &models.MovingExpense{
			Status: &rejectedStatus,
			Reason: models.StringPointer("Can't pump your own gas in New Jersey"),
		}, true)

		updater := NewOfficeMovingExpenseUpdater()
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
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMovingExpense := setupForTest(appCtx, nil, false)

		updater := NewCustomerMovingExpenseUpdater()
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

	suite.Run("Fails to update when a reason isn't provided for non-approved status", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMovingExpense := setupForTest(appCtx, nil, true)

		updater := NewOfficeMovingExpenseUpdater()
		oilExpenseType := models.MovingExpenseReceiptTypeOil

		excludedStatus := models.PPMDocumentStatusExcluded
		expectedMovingExpense := &models.MovingExpense{
			ID:                originalMovingExpense.ID,
			MovingExpenseType: &oilExpenseType,
			Description:       models.StringPointer("Fuel"),
			PaidWithGTCC:      models.BoolPointer(false),
			MissingReceipt:    models.BoolPointer(false),
			Amount:            models.CentPointer(unit.Cents(67899)),
			Status:            &excludedStatus,
		}

		updatedMovingExpense, updateErr := updater.UpdateMovingExpense(appCtx, *expectedMovingExpense, etag.GenerateEtag(originalMovingExpense.UpdatedAt))

		suite.Nil(updatedMovingExpense)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.ErrorContains(updateErr, "reason is mandatory if the status is Excluded or Rejected")
	})
}
