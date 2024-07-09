package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ServiceMemberAffiliation represents a service member's branch
type ServiceMemberAffiliation string

// String is a string representation of a ServiceMemberAffiliation
func (s ServiceMemberAffiliation) String() string {
	return string(s)
}

const (
	// AffiliationARMY captures enum value "ARMY"
	AffiliationARMY ServiceMemberAffiliation = "ARMY"
	// AffiliationNAVY captures enum value "NAVY"
	AffiliationNAVY ServiceMemberAffiliation = "NAVY"
	// AffiliationMARINES captures enum value "MARINES"
	AffiliationMARINES ServiceMemberAffiliation = "MARINES"
	// AffiliationAIRFORCE captures enum value "AIR_FORCE"
	AffiliationAIRFORCE ServiceMemberAffiliation = "AIR_FORCE"
	// AffiliationCOASTGUARD captures enum value "COAST_GUARD"
	AffiliationCOASTGUARD ServiceMemberAffiliation = "COAST_GUARD"
	// AffiliationCOASTGUARD captures enum value "SPACE_FORCE"
	AffiliationSPACEFORCE ServiceMemberAffiliation = "SPACE_FORCE"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                     uuid.UUID                 `json:"id" db:"id"`
	CreatedAt              time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                 `json:"updated_at" db:"updated_at"`
	UserID                 uuid.UUID                 `json:"user_id" db:"user_id"`
	User                   User                      `belongs_to:"user" fk_id:"user_id"`
	Edipi                  *string                   `json:"edipi" db:"edipi"`
	Emplid                 *string                   `json:"emplid" db:"emplid"`
	Affiliation            *ServiceMemberAffiliation `json:"affiliation" db:"affiliation"`
	FirstName              *string                   `json:"first_name" db:"first_name"`
	MiddleName             *string                   `json:"middle_name" db:"middle_name"`
	LastName               *string                   `json:"last_name" db:"last_name"`
	Suffix                 *string                   `json:"suffix" db:"suffix"`
	Telephone              *string                   `json:"telephone" db:"telephone"`
	SecondaryTelephone     *string                   `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail          *string                   `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred       *bool                     `json:"phone_is_preferred" db:"phone_is_preferred"`
	EmailIsPreferred       *bool                     `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID   *uuid.UUID                `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress     *Address                  `belongs_to:"address" fk_id:"residential_address_id"`
	BackupMailingAddressID *uuid.UUID                `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress   *Address                  `belongs_to:"address" fk_id:"backup_mailing_address_id"`
	Orders                 Orders                    `has_many:"orders" fk_id:"service_member_id" order_by:"created_at desc" `
	BackupContacts         BackupContacts            `has_many:"backup_contacts" fk_id:"service_member_id"`
	CacValidated           bool                      `json:"cac_validated" db:"cac_validated"`
}

// TableName overrides the table name used by Pop.
func (s ServiceMember) TableName() string {
	return "service_members"
}

type ServiceMembers []ServiceMember

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *ServiceMember) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.UserID, Name: "UserID"},
	), nil
}

// FetchServiceMemberForUser returns a service member only if it is allowed for the given user to access that service member.
// This method is thereby a useful way of performing access control checks.
func FetchServiceMemberForUser(db *pop.Connection, session *auth.Session, id uuid.UUID) (ServiceMember, error) {

	var serviceMember ServiceMember
	err := db.Q().Eager("User",
		"BackupMailingAddress",
		"BackupContacts",
		"Orders.NewDutyLocation.TransportationOffice",
		"Orders.OriginDutyLocation",
		"Orders.UploadedOrders.UserUploads.Upload",
		"Orders.Moves",
		"ResidentialAddress").Find(&serviceMember, id)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ServiceMember{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return ServiceMember{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify serviceMember
	if session.IsMilApp() && serviceMember.ID != session.ServiceMemberID {
		return ServiceMember{}, ErrFetchForbidden
	}

	// TODO: Remove this when Pop's eager loader stops populating blank structs into these fields
	if serviceMember.ResidentialAddressID == nil {
		serviceMember.ResidentialAddress = nil
	}
	if serviceMember.BackupMailingAddressID == nil {
		serviceMember.BackupMailingAddress = nil
	}

	return serviceMember, nil
}

// FetchServiceMember returns a service member by id REGARDLESS OF USER.
// Does not fetch nested models.
// DO NOT USE IF YOU NEED USER AUTH
func FetchServiceMember(db *pop.Connection, id uuid.UUID) (ServiceMember, error) {
	var serviceMember ServiceMember
	err := db.Q().Find(&serviceMember, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ServiceMember{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return ServiceMember{}, err
	}

	return serviceMember, nil
}

// SaveServiceMember takes a serviceMember with Address structs and coordinates saving it all in a transaction
func SaveServiceMember(appCtx appcontext.AppContext, serviceMember *ServiceMember) (*validate.Errors, error) {

	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		transactionError := errors.New("Rollback The transaction")

		if serviceMember.ResidentialAddress != nil {
			county, err := FindCountyByZipCode(appCtx.DB(), serviceMember.ResidentialAddress.PostalCode)
			if err != nil {
				responseError = err
				return err
			}
			serviceMember.ResidentialAddress.County = county
			if verrs, err := txnAppCtx.DB().ValidateAndSave(serviceMember.ResidentialAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.ResidentialAddressID = &serviceMember.ResidentialAddress.ID
		}

		if serviceMember.BackupMailingAddress != nil {
			county, err := FindCountyByZipCode(appCtx.DB(), serviceMember.BackupMailingAddress.PostalCode)
			if err != nil {
				responseError = err
				return err
			}
			serviceMember.BackupMailingAddress.County = county
			if verrs, err := txnAppCtx.DB().ValidateAndSave(serviceMember.BackupMailingAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			serviceMember.BackupMailingAddressID = &serviceMember.BackupMailingAddress.ID
		}

		if verrs, err := txnAppCtx.DB().ValidateAndSave(serviceMember); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	if transactionErr != nil {
		return responseVErrors, responseError
	}

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
func (s ServiceMember) CreateOrder(appCtx appcontext.AppContext,
	issueDate time.Time,
	reportByDate time.Time,
	ordersType internalmessages.OrdersType,
	hasDependents bool,
	spouseHasProGear bool,
	newDutyLocation DutyLocation,
	ordersNumber *string,
	tac *string,
	sac *string,
	departmentIndicator *string,
	originDutyLocation *DutyLocation,
	grade *internalmessages.OrderPayGrade,
	entitlement *Entitlement,
	originDutyLocationGBLOC *string,
	packingAndShippingInstructions string) (Order, *validate.Errors, error) {

	var newOrders Order
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		transactionError := errors.New("Rollback The transaction")
		uploadedOrders := Document{
			ServiceMemberID: s.ID,
			ServiceMember:   s,
		}
		verrs, err := txnAppCtx.DB().ValidateAndCreate(&uploadedOrders)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		newOrders = Order{
			ServiceMemberID:                s.ID,
			ServiceMember:                  s,
			IssueDate:                      issueDate,
			ReportByDate:                   reportByDate,
			OrdersType:                     ordersType,
			HasDependents:                  hasDependents,
			SpouseHasProGear:               spouseHasProGear,
			NewDutyLocationID:              newDutyLocation.ID,
			NewDutyLocation:                newDutyLocation,
			UploadedOrders:                 uploadedOrders,
			UploadedOrdersID:               uploadedOrders.ID,
			Status:                         OrderStatusDRAFT,
			OrdersNumber:                   ordersNumber,
			TAC:                            tac,
			SAC:                            sac,
			DepartmentIndicator:            departmentIndicator,
			Grade:                          grade,
			OriginDutyLocation:             originDutyLocation,
			Entitlement:                    entitlement,
			OriginDutyLocationGBLOC:        originDutyLocationGBLOC,
			SupplyAndServicesCostEstimate:  SupplyAndServicesCostEstimate,
			MethodOfPayment:                MethodOfPayment,
			NAICS:                          NAICS,
			PackingAndShippingInstructions: packingAndShippingInstructions,
		}

		verrs, err = txnAppCtx.DB().ValidateAndCreate(&newOrders)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	if transactionErr != nil {
		return newOrders, responseVErrors, responseError
	}

	return newOrders, responseVErrors, responseError
}

// UpdateServiceMemberDoDID is called if Safety Move order is created to clear out the DoDID
func UpdateServiceMemberDoDID(db *pop.Connection, serviceMember *ServiceMember, dodid *string) error {

	serviceMember.Edipi = dodid

	verrs, err := db.ValidateAndUpdate(serviceMember)
	if verrs.HasAny() {
		return verrs
	} else if err != nil {
		err = errors.Wrap(err, "Unable to update service member edipi")
		return err
	}

	return nil
}

// UpdateServiceMemberEMPLID is called if Safety Move order is created to clear out the EMPLID
func UpdateServiceMemberEMPLID(db *pop.Connection, serviceMember *ServiceMember, emplid *string) error {

	serviceMember.Emplid = emplid

	verrs, err := db.ValidateAndUpdate(serviceMember)
	if verrs.HasAny() {
		return verrs
	} else if err != nil {
		err = errors.Wrap(err, "Unable to update service member emplid")
		return err
	}

	return nil
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
	if s.PhoneIsPreferred == nil && s.EmailIsPreferred == nil {
		return false
	}
	if s.ResidentialAddressID == nil {
		return false
	}
	if s.BackupMailingAddressID == nil {
		return false
	}
	if len(s.BackupContacts) == 0 {
		return false
	}
	// All required fields have a set value
	return true
}

// FetchLatestOrder gets the latest order for a service member
func FetchLatestOrder(session *auth.Session, db *pop.Connection) (Order, error) {
	var order Order
	query := db.Where("orders.service_member_id = $1", session.ServiceMemberID).Order("created_at desc")
	err := query.EagerPreload("ServiceMember.User",
		"OriginDutyLocation.Address",
		"OriginDutyLocation.TransportationOffice",
		"NewDutyLocation.Address",
		"UploadedOrders",
		"UploadedAmendedOrders",
		"Moves.SignedCertifications",
		"Entitlement").
		First(&order)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		return Order{}, err
	}

	// Eager loading of nested has_many associations is broken
	var userUploads UserUploads
	err = db.Q().
		Scope(utilities.ExcludeDeletedScope()).EagerPreload("Upload").
		Where("document_id = ?", order.UploadedOrdersID).
		All(&userUploads)
	if err != nil {
		return Order{}, err
	}

	order.UploadedOrders.UserUploads = userUploads

	// Eager loading of nested has_many associations is broken
	if order.UploadedAmendedOrders != nil {
		var amendedUserUploads UserUploads
		err = db.Q().
			Scope(utilities.ExcludeDeletedScope()).EagerPreload("Upload").
			Where("document_id = ?", order.UploadedAmendedOrdersID).
			All(&amendedUserUploads)
		if err != nil {
			return Order{}, err
		}
		order.UploadedAmendedOrders.UserUploads = amendedUserUploads
	}

	// User must be logged in service member
	if session.IsMilApp() && order.ServiceMember.ID != session.ServiceMemberID {
		return Order{}, ErrFetchForbidden
	}
	return order, nil
}

// ReverseNameLineFormat returns the service member's name as a string in Last, First, M format.
func (s *ServiceMember) ReverseNameLineFormat() string {
	names := []string{}
	if s.LastName != nil && len(*s.LastName) > 0 {
		names = append(names, *s.LastName)
	}
	if s.FirstName != nil && len(*s.FirstName) > 0 {
		names = append(names, *s.FirstName)
	}
	if s.MiddleName != nil && len(*s.MiddleName) > 0 {
		names = append(names, *s.MiddleName)
	}
	return strings.Join(names, ", ")
}
