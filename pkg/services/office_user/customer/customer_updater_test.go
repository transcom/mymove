package customer

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *CustomerServiceSuite) TestCustomerUpdater() {

	customerUpdater := NewCustomerUpdater()

	suite.Run("NewNotFoundError when customer if doesn't exist", func() {
		_, err := customerUpdater.UpdateCustomer(suite.AppContextForTest(), "", models.ServiceMember{})
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("PreconditionsError when etag is stale", func() {
		expectedCustomer := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)
		staleEtag := etag.GenerateEtag(expectedCustomer.UpdatedAt.Add(-1 * time.Minute))
		_, err := customerUpdater.UpdateCustomer(suite.AppContextForTest(), staleEtag, models.ServiceMember{ID: expectedCustomer.ID})
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Customer fields are updated", func() {
		defaultCustomer := factory.BuildServiceMember(suite.DB(), nil, nil)

		var backupContacts []models.BackupContact
		backupContact := models.BackupContact{
			Email:     "newbackup@mail.com",
			FirstName: "New",
			LastName:  "Backup Contact",
			Phone:     "445-345-1212",
		}
		backupContacts = append(backupContacts, backupContact)

		updatedCustomer := models.ServiceMember{
			ID:        defaultCustomer.ID,
			LastName:  models.StringPointer("Newlastname"),
			FirstName: models.StringPointer("Newfirstname"),
			Telephone: models.StringPointer("123-455-3399"),
			ResidentialAddress: &models.Address{
				StreetAddress1: "123 New Street",
				City:           "Newcity",
				State:          "MA",
				PostalCode:     "12345",
			},
			BackupContacts: backupContacts,
			CacValidated:   true,
		}

		expectedETag := etag.GenerateEtag(defaultCustomer.UpdatedAt)
		actualCustomer, err := customerUpdater.UpdateCustomer(suite.AppContextForTest(), expectedETag, updatedCustomer)

		suite.NoError(err)
		suite.Equal(updatedCustomer.ID, actualCustomer.ID)
		suite.Equal(updatedCustomer.LastName, actualCustomer.LastName)
		suite.Equal(updatedCustomer.FirstName, actualCustomer.FirstName)
		suite.Equal(updatedCustomer.Telephone, actualCustomer.Telephone)
		suite.Equal(updatedCustomer.ResidentialAddress.StreetAddress1, actualCustomer.ResidentialAddress.StreetAddress1)
		suite.Equal(updatedCustomer.ResidentialAddress.City, actualCustomer.ResidentialAddress.City)
		suite.Equal(updatedCustomer.ResidentialAddress.PostalCode, actualCustomer.ResidentialAddress.PostalCode)
		suite.Equal(updatedCustomer.ResidentialAddress.State, actualCustomer.ResidentialAddress.State)
		suite.Equal(updatedCustomer.BackupContacts[0].FirstName, actualCustomer.BackupContacts[0].FirstName)
		suite.Equal(updatedCustomer.BackupContacts[0].LastName, actualCustomer.BackupContacts[0].LastName)
		suite.Equal(updatedCustomer.BackupContacts[0].Phone, actualCustomer.BackupContacts[0].Phone)
		suite.Equal(updatedCustomer.BackupContacts[0].Email, actualCustomer.BackupContacts[0].Email)
		suite.Equal(updatedCustomer.CacValidated, actualCustomer.CacValidated)

		updatedCustomer.BackupContacts[0].FirstName = "Updated"
		updatedCustomer.BackupContacts[0].LastName = "Backup Contact"
		expectedETag = etag.GenerateEtag(actualCustomer.UpdatedAt)
		actualCustomer, err = customerUpdater.UpdateCustomer(suite.AppContextForTest(), expectedETag, updatedCustomer)
		suite.NoError(err)
		suite.Equal(actualCustomer.BackupContacts[0].FirstName, "Updated")
		suite.Equal(actualCustomer.BackupContacts[0].LastName, "Backup Contact")
	})

	suite.Run("Empty customer is updated", func() {
		defaultCustomer := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					BackupContacts: nil,
				},
			},
		}, nil)

		defaultCustomer.ResidentialAddressID = nil
		defaultCustomer.BackupMailingAddressID = nil
		err := suite.DB().Save(&defaultCustomer)
		suite.NoError(err)

		var backupContacts []models.BackupContact
		backupContact := models.BackupContact{
			Email:     "newbackup@mail.com",
			FirstName: "New",
			LastName:  "Backup Contact",
			Phone:     "445-345-1212",
		}
		backupContacts = append(backupContacts, backupContact)

		updatedCustomer := models.ServiceMember{
			ID:        defaultCustomer.ID,
			LastName:  models.StringPointer("Newlastname"),
			FirstName: models.StringPointer("Newfirstname"),
			Telephone: models.StringPointer("123-455-3399"),
			ResidentialAddress: &models.Address{
				StreetAddress1: "123 New Street",
				City:           "Newcity",
				State:          "MA",
				PostalCode:     "12345",
			},
			BackupMailingAddress: &models.Address{
				StreetAddress1: "456 Extra Street",
				City:           "Metropolis",
				State:          "TX",
				PostalCode:     "43320",
			},
			BackupContacts: backupContacts,
			CacValidated:   true,
		}

		expectedETag := etag.GenerateEtag(defaultCustomer.UpdatedAt)
		actualCustomer, err := customerUpdater.UpdateCustomer(suite.AppContextForTest(), expectedETag, updatedCustomer)
		suite.NoError(err)
		suite.Equal(updatedCustomer.ID, actualCustomer.ID)
		suite.Equal(updatedCustomer.LastName, actualCustomer.LastName)
		suite.Equal(updatedCustomer.FirstName, actualCustomer.FirstName)
		suite.Equal(updatedCustomer.Telephone, actualCustomer.Telephone)
		suite.Equal(updatedCustomer.ResidentialAddress.StreetAddress1, actualCustomer.ResidentialAddress.StreetAddress1)
		suite.Equal(updatedCustomer.ResidentialAddress.City, actualCustomer.ResidentialAddress.City)
		suite.Equal(updatedCustomer.ResidentialAddress.PostalCode, actualCustomer.ResidentialAddress.PostalCode)
		suite.Equal(updatedCustomer.ResidentialAddress.State, actualCustomer.ResidentialAddress.State)
		suite.Equal(updatedCustomer.BackupContacts[0].FirstName, actualCustomer.BackupContacts[0].FirstName)
		suite.Equal(updatedCustomer.BackupContacts[0].LastName, actualCustomer.BackupContacts[0].LastName)
		suite.Equal(updatedCustomer.BackupContacts[0].Phone, actualCustomer.BackupContacts[0].Phone)
		suite.Equal(updatedCustomer.BackupContacts[0].Email, actualCustomer.BackupContacts[0].Email)
	})
}
