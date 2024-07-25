package adminuser

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *RequestedOfficeUsersServiceSuite) TestUpdateRequestedOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	updater := NewRequestedOfficeUserUpdater(queryBuilder)
	setupTestData := func() models.OfficeUser {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		return officeUser
	}

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())

		firstName := "Jimmy"
		lastName := "Jim"
		status := "APPROVED"
		payload := &adminmessages.RequestedOfficeUserUpdate{
			FirstName:              &firstName,
			LastName:               &lastName,
			TransportationOfficeID: handlers.FmtUUID(transportationOffice.ID),
			Status:                 status,
		}
		updatedOfficeUser, verrs, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.NotEqual(updatedOfficeUser.TransportationOfficeID.String(), officeUser.TransportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, lastName)
		suite.Equal(updatedOfficeUser.Active, true)

	})

	// Bad office user ID
	suite.Run("If we are provided an office user that doesn't exist, the create should fail", func() {
		payload := &adminmessages.RequestedOfficeUserUpdate{}

		_, _, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		officeUser := setupTestData()
		badID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
		payload := &adminmessages.RequestedOfficeUserUpdate{
			TransportationOfficeID: handlers.FmtUUID(badID),
		}

		_, _, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

}
