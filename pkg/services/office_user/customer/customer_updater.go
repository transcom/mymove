package customer

import (
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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
		return nil, err
	}

	existingETag := etag.GenerateEtag(existingCustomer.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(customer.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		if customer.ResidentialAddress != nil && strings.TrimSpace(customer.ResidentialAddress.PostalCode) != "" && strings.TrimSpace(customer.ResidentialAddress.City) != "" {
			usprc, err := models.FindByZipCodeAndCity(appCtx.DB(), customer.ResidentialAddress.PostalCode, strings.ToUpper(customer.ResidentialAddress.City))
			if err != nil {
				return err
			}

			customer.ResidentialAddress.UsPostRegionCity = usprc
			customer.ResidentialAddress.UsPostRegionCityID = &usprc.ID
		}

		if customer.BackupMailingAddress != nil && strings.TrimSpace(customer.BackupMailingAddress.PostalCode) != "" && strings.TrimSpace(customer.BackupMailingAddress.City) != "" {
			usprc, err := models.FindByZipCodeAndCity(appCtx.DB(), customer.BackupMailingAddress.PostalCode, strings.ToUpper(customer.BackupMailingAddress.City))
			if err != nil {
				return err
			}

			customer.BackupMailingAddress.UsPostRegionCity = usprc
			customer.BackupMailingAddress.UsPostRegionCityID = &usprc.ID
		}

		if residentialAddress := customer.ResidentialAddress; residentialAddress != nil {
			if existingCustomer.ResidentialAddress != nil {
				existingCustomer.ResidentialAddress.StreetAddress1 = residentialAddress.StreetAddress1
				existingCustomer.ResidentialAddress.City = residentialAddress.City
				existingCustomer.ResidentialAddress.State = residentialAddress.State
				existingCustomer.ResidentialAddress.PostalCode = residentialAddress.PostalCode
				if residentialAddress.StreetAddress2 != nil {
					existingCustomer.ResidentialAddress.StreetAddress2 = residentialAddress.StreetAddress2
				}
				if residentialAddress.StreetAddress3 != nil {
					existingCustomer.ResidentialAddress.StreetAddress3 = residentialAddress.StreetAddress3
				}
				if residentialAddress.UsPostRegionCityID != nil {
					existingCustomer.ResidentialAddress.UsPostRegionCityID = residentialAddress.UsPostRegionCityID
				}
			} else {
				newResidentialAddress := models.Address{
					StreetAddress1: residentialAddress.StreetAddress1,
					City:           residentialAddress.City,
					State:          residentialAddress.State,
					PostalCode:     residentialAddress.PostalCode,
				}

				isOconus, err := models.IsAddressOconus(txnAppCtx.DB(), newResidentialAddress)
				if err != nil {
					return err
				}
				newResidentialAddress.IsOconus = &isOconus

				if residentialAddress.StreetAddress2 != nil {
					newResidentialAddress.StreetAddress2 = residentialAddress.StreetAddress2
				}
				if residentialAddress.StreetAddress3 != nil {
					newResidentialAddress.StreetAddress3 = residentialAddress.StreetAddress3
				}
				if residentialAddress.UsPostRegionCityID != nil {
					newResidentialAddress.UsPostRegionCityID = residentialAddress.UsPostRegionCityID
				}

				existingCustomer.ResidentialAddress = &newResidentialAddress
			}

			verrs, dbErr := txnAppCtx.DB().ValidateAndSave(existingCustomer.ResidentialAddress)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(customer.ID, dbErr, verrs, "")
			}
			if dbErr != nil {
				return dbErr
			}

			existingCustomer.ResidentialAddressID = &existingCustomer.ResidentialAddress.ID
		}

		if backupAddress := customer.BackupMailingAddress; backupAddress != nil {
			if existingCustomer.BackupMailingAddress != nil {
				existingCustomer.BackupMailingAddress.StreetAddress1 = backupAddress.StreetAddress1
				existingCustomer.BackupMailingAddress.City = backupAddress.City
				existingCustomer.BackupMailingAddress.State = backupAddress.State
				existingCustomer.BackupMailingAddress.PostalCode = backupAddress.PostalCode
				if backupAddress.StreetAddress2 != nil {
					existingCustomer.BackupMailingAddress.StreetAddress2 = backupAddress.StreetAddress2
				}
				if backupAddress.StreetAddress3 != nil {
					existingCustomer.BackupMailingAddress.StreetAddress3 = backupAddress.StreetAddress3
				}
				if backupAddress.UsPostRegionCityID != nil {
					existingCustomer.BackupMailingAddress.UsPostRegionCityID = backupAddress.UsPostRegionCityID
				}
			} else {
				newBackupAddress := models.Address{
					StreetAddress1: backupAddress.StreetAddress1,
					City:           backupAddress.City,
					State:          backupAddress.State,
					PostalCode:     backupAddress.PostalCode,
				}

				isOconus, err := models.IsAddressOconus(txnAppCtx.DB(), newBackupAddress)
				if err != nil {
					return err
				}
				newBackupAddress.IsOconus = &isOconus

				if backupAddress.StreetAddress2 != nil {
					newBackupAddress.StreetAddress2 = backupAddress.StreetAddress2
				}
				if backupAddress.StreetAddress3 != nil {
					newBackupAddress.StreetAddress3 = backupAddress.StreetAddress3
				}
				if backupAddress.UsPostRegionCityID != nil {
					newBackupAddress.UsPostRegionCityID = backupAddress.UsPostRegionCityID
				}

				existingCustomer.BackupMailingAddress = &newBackupAddress
			}

			verrs, dbErr := txnAppCtx.DB().ValidateAndSave(existingCustomer.BackupMailingAddress)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(customer.ID, dbErr, verrs, "")
			}
			if dbErr != nil {
				return dbErr
			}

			existingCustomer.BackupMailingAddressID = &existingCustomer.BackupMailingAddress.ID
		}

		if backupContacts := customer.BackupContacts; len(backupContacts) > 0 {
			// added this check to prevent crashes when the customer doesn't finish creating their profile
			if len(existingCustomer.BackupContacts) > 0 {
				existingCustomer.BackupContacts[0].FirstName = backupContacts[0].FirstName
				existingCustomer.BackupContacts[0].LastName = backupContacts[0].LastName
				existingCustomer.BackupContacts[0].Email = backupContacts[0].Email
				existingCustomer.BackupContacts[0].Phone = backupContacts[0].Phone
			} else {
				backupContact, verrs, dbErr := existingCustomer.CreateBackupContact(
					txnAppCtx.DB(),
					backupContacts[0].FirstName,
					backupContacts[0].LastName,
					backupContacts[0].Email,
					backupContacts[0].Phone,
					models.BackupContactPermissionNONE,
				)
				existingCustomer.BackupContacts = append(existingCustomer.BackupContacts, backupContact)

				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(customer.ID, dbErr, verrs, "")
				}
				if dbErr != nil {
					return dbErr
				}
			}

			verrs, dbErr := txnAppCtx.DB().ValidateAndSave(existingCustomer.BackupContacts)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(customer.ID, dbErr, verrs, "")
			}
			if dbErr != nil {
				return dbErr
			}
		}

		if customer.PreferredName != nil {
			existingCustomer.PreferredName = customer.PreferredName
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

		if customer.SecondaryTelephone != nil {
			existingCustomer.SecondaryTelephone = customer.SecondaryTelephone
		}

		if customer.PhoneIsPreferred != nil {
			existingCustomer.PhoneIsPreferred = customer.PhoneIsPreferred
		}

		if customer.EmailIsPreferred != nil {
			existingCustomer.EmailIsPreferred = customer.EmailIsPreferred
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

		if customer.CacValidated != existingCustomer.CacValidated {
			existingCustomer.CacValidated = customer.CacValidated
		}

		// optimistic locking handled before transaction block
		verrs, updateErr := txnAppCtx.DB().ValidateAndUpdate(existingCustomer)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(customer.ID, err, verrs, "")
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
