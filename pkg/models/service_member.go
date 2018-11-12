package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

/*
   This file contains the implementaion details of the ServiceMember model. It should be pushed down into the users_impl package
*/

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceMember) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.UserID, Name: "UserID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SaveServiceMember takes a serviceMember with Address structs and coordinates saving it all in a transaction
func SaveServiceMember(dbConnection *pop.Connection, serviceMember *ServiceMember) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if serviceMember.ResidentialAddress != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.ResidentialAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.ResidentialAddressID = &serviceMember.ResidentialAddress.ID
		}

		if serviceMember.BackupMailingAddress != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.BackupMailingAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.BackupMailingAddressID = &serviceMember.BackupMailingAddress.ID
		}

		if serviceMember.SocialSecurityNumber != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.SocialSecurityNumber); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.SocialSecurityNumberID = &serviceMember.SocialSecurityNumber.ID
		}

		if verrs, err := dbConnection.ValidateAndSave(serviceMember); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError

}

// CreateBackupContact creates a backup contact model tied to the service member
func (s ServiceMember) CreateBackupContact(db *pop.Connection, name string, email string, phone *string, permission BackupContactPermission) (BackupContact, *validate.Errors, error) {
	newContact := BackupContact{
		ServiceMemberID: s.ID,
		ServiceMember:   s,
		Name:            name,
		Email:           email,
		Phone:           phone,
		Permission:      permission,
	}

	verrs, err := db.ValidateAndCreate(&newContact)
	if err != nil || verrs.HasAny() {
		newContact = BackupContact{}
	}
	return newContact, verrs, err
}

// CreateOrder creates an order model tied to the service member
func (s ServiceMember) CreateOrder(db *pop.Connection,
	issueDate time.Time,
	reportByDate time.Time,
	ordersType internalmessages.OrdersType,
	hasDependents bool,
	spouseHasProGear bool,
	newDutyStation DutyStation,
	ordersNumber *string,
	paragraphNumber *string,
	ordersIssuingAgency *string,
	tac *string,
	sac *string,
	departmentIndicator *string) (Order, *validate.Errors, error) {

	var newOrders Order
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")
		uploadedOrders := Document{
			ServiceMemberID: s.ID,
			ServiceMember:   s,
		}
		verrs, err := db.ValidateAndCreate(&uploadedOrders)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		newOrders = Order{
			ServiceMemberID:     s.ID,
			ServiceMember:       s,
			IssueDate:           issueDate,
			ReportByDate:        reportByDate,
			OrdersType:          ordersType,
			HasDependents:       hasDependents,
			SpouseHasProGear:    spouseHasProGear,
			NewDutyStationID:    newDutyStation.ID,
			NewDutyStation:      newDutyStation,
			UploadedOrders:      uploadedOrders,
			UploadedOrdersID:    uploadedOrders.ID,
			Status:              OrderStatusDRAFT,
			OrdersNumber:        ordersNumber,
			ParagraphNumber:     paragraphNumber,
			OrdersIssuingAgency: ordersIssuingAgency,
			TAC:                 tac,
			SAC:                 sac,
			DepartmentIndicator: departmentIndicator,
		}

		verrs, err = db.ValidateAndCreate(&newOrders)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	return newOrders, responseVErrors, responseError
}

// IsProfileComplete checks if the profile has been completely filled out
func (s *ServiceMember) IsProfileComplete() bool {

	// The following fields are required to be set for a profile to be complete
	if s.Edipi == nil {
		return false
	}
	if s.Affiliation == nil {
		return false
	}
	if s.Rank == nil {
		return false
	}
	if s.FirstName == nil {
		return false
	}
	if s.LastName == nil {
		return false
	}
	if s.Telephone == nil {
		return false
	}
	if s.PersonalEmail == nil {
		return false
	}
	if s.PhoneIsPreferred == nil && s.TextMessageIsPreferred == nil && s.EmailIsPreferred == nil {
		return false
	}
	if s.ResidentialAddressID == nil {
		return false
	}
	if s.BackupMailingAddressID == nil {
		return false
	}
	if s.SocialSecurityNumberID == nil {
		return false
	}
	if s.DutyStationID == nil {
		return false
	}
	if len(s.BackupContacts) == 0 {
		return false
	}
	// All required fields have a set value
	return true
}

// FetchLatestOrder gets the latest order for a service member
func (s ServiceMember) FetchLatestOrder(db *pop.Connection) (Order, error) {
	var order Order
	query := db.Where("service_member_id = $1", s.ID).Order("created_at desc")
	err := query.Eager("ServiceMember.User",
		"NewDutyStation.Address",
		"UploadedOrders.Uploads",
		"Moves.PersonallyProcuredMoves",
		"Moves.SignedCertifications").First(&order)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		return Order{}, err
	}
	return order, nil
}

// ReverseNameLineFormat returns the service member's name as a string in Last, First, M format.
func (s *ServiceMember) ReverseNameLineFormat() string {
	var names []string
	if s.FirstName != nil && len(*s.FirstName) > 0 {
		names = append(names, *s.FirstName)
	}
	if s.LastName != nil && len(*s.LastName) > 0 {
		names = append(names, *s.LastName)
	}
	if s.MiddleName != nil && len(*s.MiddleName) > 0 {
		names = append(names, *s.MiddleName)
	}
	return strings.Join(names, ", ")
}

type popServiceMemberDB struct {
	db *pop.Connection
}

// NewServiceMemberDB is the DI provider to create a pop based ServiceMemberDB
func NewServiceMemberDB(db *pop.Connection) ServiceMemberDB {
	return &popServiceMemberDB{db}
}

func (pdb *popServiceMemberDB) Save(serviceMember *ServiceMember) (ValidationErrors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	pdb.db.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if serviceMember.ResidentialAddress != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.ResidentialAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.ResidentialAddressID = &serviceMember.ResidentialAddress.ID
		}

		if serviceMember.BackupMailingAddress != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.BackupMailingAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.BackupMailingAddressID = &serviceMember.BackupMailingAddress.ID
		}

		if serviceMember.SocialSecurityNumber != nil {
			if verrs, err := dbConnection.ValidateAndSave(serviceMember.SocialSecurityNumber); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.SocialSecurityNumberID = &serviceMember.SocialSecurityNumber.ID
		}

		if verrs, err := dbConnection.ValidateAndSave(serviceMember); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})
	return responseVErrors, responseError
}

func (pdb *popServiceMemberDB) Fetch(id uuid.UUID, loadAssociations bool) (*ServiceMember, error) {
	var serviceMember ServiceMember
	q := pdb.db.Q()
	if loadAssociations {
		q = q.Eager()
	}
	err := q.Find(&serviceMember, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	// TODO: Remove this when Pop's eager loader stops populating blank structs into these fields
	if serviceMember.ResidentialAddressID == nil {
		serviceMember.ResidentialAddress = nil
	}
	if serviceMember.BackupMailingAddressID == nil {
		serviceMember.BackupMailingAddress = nil
	}
	if serviceMember.SocialSecurityNumberID == nil {
		serviceMember.SocialSecurityNumber = nil
	}
	return &serviceMember, nil
}

func (pdb *popServiceMemberDB) IsTspManagingShipment(tspUserID uuid.UUID, serviceMemberID uuid.UUID) (bool, error) {
	// A TspUser is only allowed to interact with a service member if they are associated with one of their shipments.
	query := `
			SELECT tsp_users.id FROM tsp_users, shipment_offers, shipments
			WHERE
				tsp_users.transportation_service_provider_id = shipment_offers.transportation_service_provider_id
				AND shipment_offers.shipment_id = shipments.id
				AND shipment_offers.accepted IS NOT FALSE
				AND tsp_users.id = $1
				AND shipments.service_member_id = $2
		`

	count, err := pdb.db.RawQuery(query, tspUserID, serviceMemberID).Count(TspUser{})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
