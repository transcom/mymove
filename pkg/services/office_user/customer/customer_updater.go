package customer

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type customerUpdater struct {
	db *pop.Connection
	fetchCustomer
}

// NewCustomerUpdater creates a new struct with the service dependencies
func NewCustomerUpdater(db *pop.Connection) services.CustomerUpdater {
	return &customerUpdater{db, fetchCustomer{db}}
}

// UpdateCustomer updates the Customer model
func (s *customerUpdater) UpdateCustomer(eTag string, customer models.ServiceMember) (*models.ServiceMember, error) {
	existingCustomer, err := s.fetchCustomer.FetchCustomer(customer.ID)
	if err != nil {
		return nil, services.NewNotFoundError(customer.ID, "while looking for customer")
	}

	existingETag := etag.GenerateEtag(existingCustomer.UpdatedAt)
	if existingETag != eTag {
		return nil, services.NewPreconditionFailedError(customer.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := s.db.Transaction(func(tx *pop.Connection) error {
		if residentialAddress := customer.ResidentialAddress; residentialAddress != nil {
			existingCustomer.ResidentialAddress.StreetAddress1 = residentialAddress.StreetAddress1
			existingCustomer.ResidentialAddress.City = residentialAddress.City
			existingCustomer.ResidentialAddress.State = residentialAddress.State
			existingCustomer.ResidentialAddress.PostalCode = residentialAddress.PostalCode
			if residentialAddress.StreetAddress2 != nil {
				existingCustomer.ResidentialAddress.StreetAddress2 = residentialAddress.StreetAddress2
			}

			err = tx.Save(existingCustomer.ResidentialAddress)
			if err != nil {
				return err
			}
		}

		if backupContacts := customer.BackupContacts; backupContacts != nil {
			existingCustomer.BackupContacts[0].Name = backupContacts[0].Name
			existingCustomer.BackupContacts[0].Email = backupContacts[0].Email
			existingCustomer.BackupContacts[0].Phone = backupContacts[0].Phone

			err = tx.Save(existingCustomer.BackupContacts)
			if err != nil {
				return err
			}
		}

		if customer.FirstName != nil {
			existingCustomer.FirstName = customer.FirstName
		}

		if customer.LastName != nil {
			existingCustomer.LastName = customer.LastName
		}

		if customer.PersonalEmail != nil {
			existingCustomer.PersonalEmail = customer.PersonalEmail
		}

		if customer.Telephone != nil {
			existingCustomer.Telephone = customer.Telephone
		}

		// optimistic locking handled before transaction block
		verrs, updateErr := tx.ValidateAndUpdate(existingCustomer)

		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidInputError(customer.ID, err, verrs, "")
		}

		if updateErr != nil {
			return updateErr
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return existingCustomer, err
}
