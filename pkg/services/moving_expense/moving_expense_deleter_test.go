package movingexpense

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MovingExpenseSuite) TestDeleteMovingExpense() {

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
		notFoundMovingExpenseID := uuid.Must(uuid.NewV4())
		deleter := NewMovingExpenseDeleter()

		err := deleter.DeleteMovingExpense(suite.AppContextForTest(), notFoundMovingExpenseID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for MovingExpense", notFoundMovingExpenseID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Successfully deletes as a customer's moving expense", func() {
		appCtx := suite.AppContextForTest()

		originalMovingExpense := setupForTest(appCtx, nil, true)

		deleter := NewMovingExpenseDeleter()

		suite.Nil(originalMovingExpense.DeletedAt)
		err := deleter.DeleteMovingExpense(appCtx, originalMovingExpense.ID)
		suite.NoError(err)

		var movingExpenseInDB models.MovingExpense
		err = suite.DB().Find(&movingExpenseInDB, originalMovingExpense.ID)
		suite.NoError(err)
		suite.NotNil(movingExpenseInDB.DeletedAt)

		// Should not delete associated PPM shipment
		var dbPPMShipment models.PPMShipment
		suite.NotNil(originalMovingExpense.PPMShipmentID)
		err = suite.DB().Find(&dbPPMShipment, originalMovingExpense.PPMShipmentID)
		suite.NoError(err)
		suite.Nil(dbPPMShipment.DeletedAt)

		// Should delete associated document
		var dbDocument models.Document
		suite.NotNil(originalMovingExpense.DocumentID)
		err = suite.DB().Find(&dbDocument, originalMovingExpense.DocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)
	})
}
