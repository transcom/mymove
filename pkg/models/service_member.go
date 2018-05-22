package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                     uuid.UUID                           `json:"id" db:"id"`
	CreatedAt              time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                           `json:"updated_at" db:"updated_at"`
	UserID                 uuid.UUID                           `json:"user_id" db:"user_id"`
	User                   User                                `belongs_to:"user"`
	Edipi                  *string                             `json:"edipi" db:"edipi"`
	Affiliation            *internalmessages.Affiliation       `json:"affiliation" db:"affiliation"`
	Rank                   *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	FirstName              *string                             `json:"first_name" db:"first_name"`
	MiddleName             *string                             `json:"middle_name" db:"middle_name"`
	LastName               *string                             `json:"last_name" db:"last_name"`
	Suffix                 *string                             `json:"suffix" db:"suffix"`
	Telephone              *string                             `json:"telephone" db:"telephone"`
	SecondaryTelephone     *string                             `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail          *string                             `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred       *bool                               `json:"phone_is_preferred" db:"phone_is_preferred"`
	TextMessageIsPreferred *bool                               `json:"text_message_is_preferred" db:"text_message_is_preferred"`
	EmailIsPreferred       *bool                               `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID   *uuid.UUID                          `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress     *Address                            `belongs_to:"address"`
	BackupMailingAddressID *uuid.UUID                          `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress   *Address                            `belongs_to:"address"`
	SocialSecurityNumberID *uuid.UUID                          `json:"social_security_number_id" db:"social_security_number_id"`
	SocialSecurityNumber   *SocialSecurityNumber               `belongs_to:"address"`
	Orders                 Orders                              `has_many:"orders" order_by:"created_at desc"`
	BackupContacts         *BackupContacts                     `has_many:"backup_contacts"`
	DutyStationID          *uuid.UUID                          `json:"duty_station_id" db:"duty_station_id"`
	DutyStation            *DutyStation                        `belongs_to:"duty_stations"`
}

// ServiceMembers is not required by pop and may be deleted
type ServiceMembers []ServiceMember

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

// FetchServiceMember returns a service member only if it is allowed for the given user to access that service member.
// This method is thereby a useful way of performing access control checks.
func FetchServiceMember(db *pop.Connection, user User, reqApp string, id uuid.UUID) (ServiceMember, error) {
	var serviceMember ServiceMember
	err := db.Q().Eager().Find(&serviceMember, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return ServiceMember{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return ServiceMember{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify serviceMember
	if reqApp == app.MyApp && serviceMember.UserID != user.ID {
		return ServiceMember{}, ErrFetchForbidden
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

	if serviceMember.DutyStationID == nil {
		serviceMember.DutyStation = nil
	} else {
		// Need to do this because Pop's nested eager loading seems to be broken
		db.Q().Eager().Find(&serviceMember.DutyStation.Address, serviceMember.DutyStation.AddressID)
	}

	return serviceMember, nil
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
func (s ServiceMember) CreateBackupContact(db *pop.Connection, name string, email string, phone *string, permission internalmessages.BackupContactPermission) (BackupContact, *validate.Errors, error) {
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
	newDutyStation DutyStation) (Order, *validate.Errors, error) {

	var newOrders Order
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")
		uploadedOrders := Document{
			ServiceMemberID: s.ID,
			ServiceMember:   s,
			Name:            UploadedOrdersDocumentName,
		}
		verrs, err := db.ValidateAndCreate(&uploadedOrders)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		newOrders = Order{
			ServiceMemberID:  s.ID,
			ServiceMember:    s,
			IssueDate:        issueDate,
			ReportByDate:     reportByDate,
			OrdersType:       ordersType,
			HasDependents:    hasDependents,
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
			UploadedOrders:   uploadedOrders,
			UploadedOrdersID: uploadedOrders.ID,
			Status:           OrderStatusDRAFT,
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
	if s.ResidentialAddress == nil {
		return false
	}
	if s.BackupMailingAddress == nil {
		return false
	}
	if s.SocialSecurityNumberID == nil {
		return false
	}
	if s.DutyStationID == nil {
		return false
	}
	if s.BackupContacts == nil {
		return false
	}
	// All required fields have a set value
	return true
}

// FetchLatestOrder gets the latest order for a service member
func (s ServiceMember) FetchLatestOrder(db *pop.Connection) (Order, error) {
	var order Order
	query := db.Where("service_member_id = $1", s.ID).Order("created_at desc")
	err := query.Eager("ServiceMember.User", "NewDutyStation.Address", "UploadedOrders.Uploads").First(&order)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		return Order{}, err
	}
	return order, nil
}
