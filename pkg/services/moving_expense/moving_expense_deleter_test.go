package movingexpense

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MovingExpenseSuite) TestDeleteMovingExpense() {

	setupForTest := func(overrides *models.MovingExpense, hasDocumentUploads bool) *models.MovingExpense {
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

		verrs, err := suite.DB().ValidateAndCreate(&originalMovingExpense)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalMovingExpense.ID)

		return &originalMovingExpense
	}
	suite.Run("Returns an error if the original doesn't exist", func() {
		notFoundMovingExpenseID := uuid.Must(uuid.NewV4())
		ppmID := uuid.Must(uuid.NewV4())
		deleter := NewMovingExpenseDeleter()

		err := deleter.DeleteMovingExpense(suite.AppContextWithSessionForTest(&auth.Session{}), ppmID, notFoundMovingExpenseID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for MovingExpense", notFoundMovingExpenseID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Successfully deletes as a customer's moving expense", func() {
		originalMovingExpense := setupForTest(nil, true)
		deleter := NewMovingExpenseDeleter()

		suite.Nil(originalMovingExpense.DeletedAt)
		err := deleter.DeleteMovingExpense(suite.AppContextWithSessionForTest(&auth.Session{
			ServiceMemberID: originalMovingExpense.Document.ServiceMemberID,
		}), originalMovingExpense.PPMShipmentID, originalMovingExpense.ID)
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
