package officeuser

import (
	"database/sql"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OfficeUserServiceSuite) TestUpdateOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	updater := NewOfficeUserUpdater(queryBuilder)
	setupTestData := func() models.OfficeUser {
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				TransportationOffice: models.TransportationOffice{
					Name: "Random Office",
				},
			},
		})
		return officeUser
	}

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		officeUser := setupTestData()
		transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())

		firstName := "Lea"
		payload := &adminmessages.OfficeUserUpdatePayload{
			FirstName:              &firstName,
			TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
		}

		updatedOfficeUser, verrs, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, officeUser.LastName)
	})

	// Bad office user ID
	suite.Run("If we are provided an office user that doesn't exist, the create should fail", func() {
		payload := &adminmessages.OfficeUserUpdatePayload{}

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		officeUser := setupTestData()
		payload := &adminmessages.OfficeUserUpdatePayload{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})
}
