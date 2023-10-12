package sitentrydateupdate

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UpdateSitEntryDateServiceSuite) TestUpdateSitEntryDate() {

	updater := NewSitEntryDateUpdater()

	setupSitEntryDateUpdateModel := func() models.SITEntryDateUpdate {
		serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		sitEntryDateUpdateModel := models.SITEntryDateUpdate{
			ID:           serviceItem.ID,
			SITEntryDate: serviceItem.SITEntryDate,
		}
		return sitEntryDateUpdateModel
	}

	// Test not found error
	suite.Run("Not Found Error", func() {
		serviceItem := setupSitEntryDateUpdateModel()
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundServiceItem := serviceItem
		notFoundServiceItem.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &notFoundServiceItem)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test successful update
	suite.Run("Successful update of service item ", func() {
		serviceItem := setupSitEntryDateUpdateModel()
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		newServiceItem := serviceItem
		newServiceItem.SITEntryDate = &sitEntryDate

		updatedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &newServiceItem)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(serviceItem.ID, updatedServiceItem.ID)
		suite.Equal(newServiceItem.SITEntryDate.Local(), updatedServiceItem.SITEntryDate.Local())
	})

}
