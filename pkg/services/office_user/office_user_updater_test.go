package officeuser

import (
	"database/sql"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *OfficeUserServiceSuite) TestUpdateOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	updater := NewOfficeUserUpdater(queryBuilder)
	setupTestData := func() models.OfficeUser {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		return officeUser
	}

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())

		firstName := "Lea"
		payload := &adminmessages.OfficeUserUpdate{
			FirstName:              &firstName,
			TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
		}

		updatedOfficeUser, verrs, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.NotEqual(updatedOfficeUser.TransportationOfficeID.String(), officeUser.TransportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, officeUser.LastName)
	})

	// Bad office user ID
	suite.Run("If we are provided an office user that doesn't exist, the create should fail", func() {
		payload := &adminmessages.OfficeUserUpdate{}

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		officeUser := setupTestData()
		payload := &adminmessages.OfficeUserUpdate{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})
}
