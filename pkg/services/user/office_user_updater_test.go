package user

import (
	"testing"

	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UserServiceSuite) TestUpdateOfficeUser() {
	transportationOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{})

	newUUID, _ := uuid.NewV4()

	userInfo := models.OfficeUser{
		ID:                     newUUID,
		LastName:               "Spaceman",
		FirstName:              "Leo",
		Email:                  "spaceman@leo.org",
		TransportationOfficeID: transportationOffice.ID,
		Telephone:              "312-111-1111",
		TransportationOffice:   transportationOffice,
	}

	// Happy path
	suite.T().Run("If the user is updated successfully it should be returned", func(t *testing.T) {
		fakeUpdateOne := func(interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(model interface{}) error {
			return nil
		}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewOfficeUserUpdater(builder)
		_, verrs, err := updater.UpdateOfficeUser(&userInfo)
		suite.NoError(err)
		suite.Nil(verrs)
	})

	// Bad transportation office ID
	suite.T().Run("If we are provided a transportation office that doesn't exist, the create should fail", func(t *testing.T) {
		fakeUpdateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(model interface{}) error {
			return models.ErrFetchNotFound
		}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}

		updater := NewOfficeUserUpdater(builder)
		_, _, err := updater.UpdateOfficeUser(&userInfo)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())

	})

}
