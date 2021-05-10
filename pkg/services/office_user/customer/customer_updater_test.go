package customer

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerServiceSuite) TestCustomerUpdater() {
	expectedCustomer := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{})

	customerUpdater := NewCustomerUpdater((suite.DB()))

	suite.T().Run("NewNotFoundError when customer if doesn't exist", func(t *testing.T) {
		_, err := customerUpdater.UpdateCustomer("", models.ServiceMember{})
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("PreconditionsError when etag is stale", func(t *testing.T) {
		staleEtag := etag.GenerateEtag(expectedCustomer.UpdatedAt.Add(-1 * time.Minute))
		_, err := customerUpdater.UpdateCustomer(staleEtag, models.ServiceMember{ID: expectedCustomer.ID})
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Customer fields are updated", func(t *testing.T) {
		defaultCustomer := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{})

		var backupContacts []models.BackupContact
		backupContact := models.BackupContact{
			Email: "newbackup@mail.com",
			Name:  "New Backup Contact",
			Phone: swag.String("445-345-1212"),
		}
		backupContacts = append(backupContacts, backupContact)

		updatedCustomer := models.ServiceMember{
			ID:        defaultCustomer.ID,
			LastName:  swag.String("Newlastname"),
			FirstName: swag.String("Newfirstname"),
			Telephone: swag.String("123-455-3399"),
			ResidentialAddress: &models.Address{
				StreetAddress1: "123 New Street",
				City:           "Newcity",
				State:          "MA",
				PostalCode:     "12345",
			},
			BackupContacts: backupContacts,
		}

		expectedETag := etag.GenerateEtag(defaultCustomer.UpdatedAt)
		actualCustomer, err := customerUpdater.UpdateCustomer(expectedETag, updatedCustomer)

		suite.NoError(err)
		suite.Equal(updatedCustomer.ID, actualCustomer.ID)
		suite.Equal(updatedCustomer.LastName, actualCustomer.LastName)
		suite.Equal(updatedCustomer.FirstName, actualCustomer.FirstName)
		suite.Equal(updatedCustomer.Telephone, actualCustomer.Telephone)
		suite.Equal(updatedCustomer.ResidentialAddress.StreetAddress1, actualCustomer.ResidentialAddress.StreetAddress1)
		suite.Equal(updatedCustomer.ResidentialAddress.City, actualCustomer.ResidentialAddress.City)
		suite.Equal(updatedCustomer.ResidentialAddress.PostalCode, actualCustomer.ResidentialAddress.PostalCode)
		suite.Equal(updatedCustomer.ResidentialAddress.State, actualCustomer.ResidentialAddress.State)
		suite.Equal(updatedCustomer.BackupContacts[0].Name, actualCustomer.BackupContacts[0].Name)
		suite.Equal(updatedCustomer.BackupContacts[0].Phone, actualCustomer.BackupContacts[0].Phone)
		suite.Equal(updatedCustomer.BackupContacts[0].Email, actualCustomer.BackupContacts[0].Email)
	})
}
