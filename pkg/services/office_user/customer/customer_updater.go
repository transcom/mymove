package customer

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type customerUpdater struct {
	fetchCustomer
}

// NewCustomerUpdater creates a new struct with the service dependencies
func NewCustomerUpdater() services.CustomerUpdater {
	return &customerUpdater{fetchCustomer{}}
}

// UpdateCustomer updates the Customer model
func (s *customerUpdater) UpdateCustomer(appCtx appcontext.AppContext, eTag string, customer models.ServiceMember) (*models.ServiceMember, error) {
	existingCustomer, err := s.fetchCustomer.FetchCustomer(appCtx, customer.ID)
	if err != nil {
		return nil, services.NewNotFoundError(customer.ID, "while looking for customer")
	}

	existingETag := etag.GenerateEtag(existingCustomer.UpdatedAt)
	if existingETag != eTag {
		return nil, services.NewPreconditionFailedError(customer.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if residentialAddress := customer.ResidentialAddress; residentialAddress != nil {
			existingCustomer.ResidentialAddress.StreetAddress1 = residentialAddress.StreetAddress1
			existingCustomer.ResidentialAddress.City = residentialAddress.City
			existingCustomer.ResidentialAddress.State = residentialAddress.State
			existingCustomer.ResidentialAddress.PostalCode = residentialAddress.PostalCode
			if residentialAddress.StreetAddress2 != nil {
				existingCustomer.ResidentialAddress.StreetAddress2 = residentialAddress.StreetAddress2
			}

			verrs, dbErr := txnAppCtx.DB().ValidateAndSave(existingCustomer.ResidentialAddress)
			if verrs != nil && verrs.HasAny() {
				return services.NewInvalidInputError(customer.ID, dbErr, verrs, "")
			}
			if dbErr != nil {
				return dbErr
			}
		}

		if backupContacts := customer.BackupContacts; len(backupContacts) > 0 {
			existingCustomer.BackupContacts[0].Name = backupContacts[0].Name
			existingCustomer.BackupContacts[0].Email = backupContacts[0].Email
			existingCustomer.BackupContacts[0].Phone = backupContacts[0].Phone

			verrs, dbErr := txnAppCtx.DB().ValidateAndSave(existingCustomer.BackupContacts)
			if verrs != nil && verrs.HasAny() {
				return services.NewInvalidInputError(customer.ID, dbErr, verrs, "")
			}
			if dbErr != nil {
				return dbErr
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

		if customer.Suffix != nil {
			if len(*customer.Suffix) == 0 {
				existingCustomer.Suffix = nil
			} else {
				existingCustomer.Suffix = customer.Suffix
			}
		}

		if customer.MiddleName != nil {
			if len(*customer.MiddleName) == 0 {
				existingCustomer.MiddleName = nil
			} else {
				existingCustomer.MiddleName = customer.MiddleName
			}
		}

		// optimistic locking handled before transaction block
		verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(existingCustomer)

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
