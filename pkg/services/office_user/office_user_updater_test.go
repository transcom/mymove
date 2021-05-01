package officeuser

import (
	"database/sql"
	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *OfficeUserServiceSuite) TestUpdateOfficeUser() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	updater := NewOfficeUserUpdater(queryBuilder)
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			TransportationOffice: models.TransportationOffice{
				Name: "Random Office",
			},
		},
	})
	transportationOffice := testdatagen.MakeDefaultTransportationOffice(suite.DB())

	// Happy path
	suite.T().Run("If the user is updated successfully it should be returned", func(t *testing.T) {
		firstName := "Lea"
		payload := &adminmessages.OfficeUserUpdatePayload{
			FirstName:              &firstName,
			TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
		}

		updatedOfficeUser, verrs, err := updater.UpdateOfficeUser(officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, officeUser.LastName)
	})

	// Bad office user ID
	suite.T().Run("If we are provided an office user that doesn't exist, the create should fail", func(t *testing.T) {
		payload := &adminmessages.OfficeUserUpdatePayload{}

		_, _, err := updater.UpdateOfficeUser(uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.T().Run("If we are provided a transportation office that doesn't exist, the create should fail", func(t *testing.T) {
		payload := &adminmessages.OfficeUserUpdatePayload{
			TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
		}

		_, _, err := updater.UpdateOfficeUser(officeUser.ID, payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})
}
