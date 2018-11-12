package models

/*
	This file contains the public interface to the User/people related parts of the DB model. It should not leak
    any of the implementation details, if possible
*/
import (
	"github.com/gobuffalo/uuid"
	"time"
)

// ServiceMemberRank represents a service member's rank
type ServiceMemberRank string

const (
	// ServiceMemberRankE1 captures enum value "E_1"
	ServiceMemberRankE1 ServiceMemberRank = "E_1"
	// ServiceMemberRankE2 captures enum value "E_2"
	ServiceMemberRankE2 ServiceMemberRank = "E_2"
	// ServiceMemberRankE3 captures enum value "E_3"
	ServiceMemberRankE3 ServiceMemberRank = "E_3"
	// ServiceMemberRankE4 captures enum value "E_4"
	ServiceMemberRankE4 ServiceMemberRank = "E_4"
	// ServiceMemberRankE5 captures enum value "E_5"
	ServiceMemberRankE5 ServiceMemberRank = "E_5"
	// ServiceMemberRankE6 captures enum value "E_6"
	ServiceMemberRankE6 ServiceMemberRank = "E_6"
	// ServiceMemberRankE7 captures enum value "E_7"
	ServiceMemberRankE7 ServiceMemberRank = "E_7"
	// ServiceMemberRankE8 captures enum value "E_8"
	ServiceMemberRankE8 ServiceMemberRank = "E_8"
	// ServiceMemberRankE9 captures enum value "E_9"
	ServiceMemberRankE9 ServiceMemberRank = "E_9"
	// ServiceMemberRankO1W1ACADEMYGRADUATE captures enum value "O_1_W_1_ACADEMY_GRADUATE"
	ServiceMemberRankO1W1ACADEMYGRADUATE ServiceMemberRank = "O_1_W_1_ACADEMY_GRADUATE"
	// ServiceMemberRankO2W2 captures enum value "O_2_W_2"
	ServiceMemberRankO2W2 ServiceMemberRank = "O_2_W_2"
	// ServiceMemberRankO3W3 captures enum value "O_3_W_3"
	ServiceMemberRankO3W3 ServiceMemberRank = "O_3_W_3"
	// ServiceMemberRankO4W4 captures enum value "O_4_W_4"
	ServiceMemberRankO4W4 ServiceMemberRank = "O_4_W_4"
	// ServiceMemberRankO5W5 captures enum value "O_5_W_5"
	ServiceMemberRankO5W5 ServiceMemberRank = "O_5_W_5"
	// ServiceMemberRankO6 captures enum value "O_6"
	ServiceMemberRankO6 ServiceMemberRank = "O_6"
	// ServiceMemberRankO7 captures enum value "O_7"
	ServiceMemberRankO7 ServiceMemberRank = "O_7"
	// ServiceMemberRankO8 captures enum value "O_8"
	ServiceMemberRankO8 ServiceMemberRank = "O_8"
	// ServiceMemberRankO9 captures enum value "O_9"
	ServiceMemberRankO9 ServiceMemberRank = "O_9"
	// ServiceMemberRankO10 captures enum value "O_10"
	ServiceMemberRankO10 ServiceMemberRank = "O_10"
	// ServiceMemberRankAVIATIONCADET captures enum value "AVIATION_CADET"
	ServiceMemberRankAVIATIONCADET ServiceMemberRank = "AVIATION_CADET"
	// ServiceMemberRankCIVILIANEMPLOYEE captures enum value "CIVILIAN_EMPLOYEE"
	ServiceMemberRankCIVILIANEMPLOYEE ServiceMemberRank = "CIVILIAN_EMPLOYEE"
	// ServiceMemberRankACADEMYCADETMIDSHIPMAN captures enum value "ACADEMY_CADET_MIDSHIPMAN"
	ServiceMemberRankACADEMYCADETMIDSHIPMAN ServiceMemberRank = "ACADEMY_CADET_MIDSHIPMAN"
)

// ServiceMemberAffiliation represents a service member's branch
type ServiceMemberAffiliation string

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
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                     uuid.UUID                 `json:"id" db:"id"`
	CreatedAt              time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                 `json:"updated_at" db:"updated_at"`
	UserID                 uuid.UUID                 `json:"user_id" db:"user_id"`
	User                   User                      `belongs_to:"user"`
	Edipi                  *string                   `json:"edipi" db:"edipi"`
	Affiliation            *ServiceMemberAffiliation `json:"affiliation" db:"affiliation"`
	Rank                   *ServiceMemberRank        `json:"rank" db:"rank"`
	FirstName              *string                   `json:"first_name" db:"first_name"`
	MiddleName             *string                   `json:"middle_name" db:"middle_name"`
	LastName               *string                   `json:"last_name" db:"last_name"`
	Suffix                 *string                   `json:"suffix" db:"suffix"`
	Telephone              *string                   `json:"telephone" db:"telephone"`
	SecondaryTelephone     *string                   `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail          *string                   `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred       *bool                     `json:"phone_is_preferred" db:"phone_is_preferred"`
	TextMessageIsPreferred *bool                     `json:"text_message_is_preferred" db:"text_message_is_preferred"`
	EmailIsPreferred       *bool                     `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID   *uuid.UUID                `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress     *Address                  `belongs_to:"address"`
	BackupMailingAddressID *uuid.UUID                `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress   *Address                  `belongs_to:"address"`
	SocialSecurityNumberID *uuid.UUID                `json:"social_security_number_id" db:"social_security_number_id"`
	SocialSecurityNumber   *SocialSecurityNumber     `belongs_to:"address"`
	Orders                 Orders                    `has_many:"orders" order_by:"created_at desc"`
	BackupContacts         BackupContacts            `has_many:"backup_contacts"`
	DutyStationID          *uuid.UUID                `json:"duty_station_id" db:"duty_station_id"`
	DutyStation            DutyStation               `belongs_to:"duty_stations"`
}

// ServiceMembers is not required by pop and may be deleted
type ServiceMembers []ServiceMember

// ServiceMemberDB defines the functions needed from the DB to access models.ServiceMembers
type ServiceMemberDB interface {
	Save(serviceMember *ServiceMember) (ValidationErrors, error)
	Fetch(id uuid.UUID, loadAssociations bool) (*ServiceMember, error)
	IsTspManagingShipment(tspUserID uuid.UUID, serviceMemberID uuid.UUID) (bool, error)
}
