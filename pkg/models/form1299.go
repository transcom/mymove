package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

// Form1299 is an application for shipment or storage of personal property
type Form1299 struct {
	ID                               uuid.UUID  `json:"id" db:"id"`
	CreatedAt                        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                        time.Time  `json:"updated_at" db:"updated_at"`
	DatePrepared                     *time.Time `json:"date_prepared" db:"date_prepared"`
	ShipmentNumber                   *string    `json:"shipment_number" db:"shipment_number"`
	NameOfPreparingOffice            *string    `json:"name_of_preparing_office" db:"name_of_preparing_office"`
	DestOfficeName                   *string    `json:"dest_office_name" db:"dest_office_name"`
	OriginOfficeAddressName          *string    `json:"origin_office_address_name" db:"origin_office_address_name"`
	OriginOfficeAddress              *string    `json:"origin_office_address" db:"origin_office_address"`
	ServiceMemberFirstName           *string    `json:"service_member_first_name" db:"service_member_first_name"`
	ServiceMemberMiddleInitial       *string    `json:"service_member_middle_initial" db:"service_member_middle_initial"`
	ServiceMemberLastName            *string    `json:"service_member_last_name" db:"service_member_last_name"`
	ServiceMemberRank                *string    `json:"service_member_rank" db:"service_member_rank"`
	ServiceMemberSsn                 *string    `json:"service_member_ssn" db:"service_member_ssn"`
	ServiceMemberAgency              *string    `json:"service_member_agency" db:"service_member_agency"`
	HhgTotalPounds                   *int64     `json:"hhg_total_pounds" db:"hhg_total_pounds"`
	HhgProgearPounds                 *int64     `json:"hhg_progear_pounds" db:"hhg_progear_pounds"`
	HhgValuableItemsCartons          *int64     `json:"hhg_valuable_items_cartons" db:"hhg_valuable_items_cartons"`
	MobileHomeSerialNumber           *string    `json:"mobile_home_serial_number" db:"mobile_home_serial_number"`
	MobileHomeLengthFt               *int64     `json:"mobile_home_length_ft" db:"mobile_home_length_ft"`
	MobileHomeLengthInches           *int64     `json:"mobile_home_length_inches" db:"mobile_home_length_inches"`
	MobileHomeWidthFt                *int64     `json:"mobile_home_width_ft" db:"mobile_home_width_ft"`
	MobileHomeWidthInches            *int64     `json:"mobile_home_width_inches" db:"mobile_home_width_inches"`
	MobileHomeHeightFt               *int64     `json:"mobile_home_height_ft" db:"mobile_home_height_ft"`
	MobileHomeHeightInches           *int64     `json:"mobile_home_height_inches" db:"mobile_home_height_inches"`
	MobileHomeTypeExpando            *string    `json:"mobile_home_type_expando" db:"mobile_home_type_expando"`
	MobileHomeServicesRequested      *string    `json:"mobile_home_services_requested" db:"mobile_home_services_requested"`
	StationOrdersType                *string    `json:"station_orders_type" db:"station_orders_type"`
	StationOrdersIssuedBy            *string    `json:"station_orders_issued_by" db:"station_orders_issued_by"`
	StationOrdersNewAssignment       *string    `json:"station_orders_new_assignment" db:"station_orders_new_assignment"`
	StationOrdersDate                *time.Time `json:"station_orders_date" db:"station_orders_date"`
	StationOrdersNumber              *string    `json:"station_orders_number" db:"station_orders_number"`
	StationOrdersParagraphNumber     *string    `json:"station_orders_paragraph_number" db:"station_orders_paragraph_number"`
	StationOrdersInTransitTelephone  *string    `json:"station_orders_in_transit_telephone" db:"station_orders_in_transit_telephone"`
	InTransitAddress                 *string    `json:"in_transit_address" db:"in_transit_address"`
	PickupAddress                    *string    `json:"pickup_address" db:"pickup_address"`
	PickupAddressMobileCourtName     *string    `json:"pickup_address_mobile_court_name" db:"pickup_address_mobile_court_name"`
	PickupTelephone                  *string    `json:"pickup_telephone" db:"pickup_telephone"`
	DestAddress                      *string    `json:"dest_address" db:"dest_address"`
	DestAddressMobileCourtName       *string    `json:"dest_address_mobile_court_name" db:"dest_address_mobile_court_name"`
	AgentToReceiveHhg                *string    `json:"agent_to_receive_hhg" db:"agent_to_receive_hhg"`
	ExtraAddress                     *string    `json:"extra_address" db:"extra_address"`
	PackScheduledDate                *time.Time `json:"pack_scheduled_date" db:"pack_scheduled_date"`
	PickupScheduledDate              *time.Time `json:"pickup_scheduled_date" db:"pickup_scheduled_date"`
	DeliveryScheduledDate            *time.Time `json:"delivery_scheduled_date" db:"delivery_scheduled_date"`
	Remarks                          *string    `json:"remarks" db:"remarks"`
	OtherMoveFrom                    *string    `json:"other_move_from" db:"other_move_from"`
	OtherMoveTo                      *string    `json:"other_move_to" db:"other_move_to"`
	OtherMoveNetPounds               *int64     `json:"other_move_net_pounds" db:"other_move_net_pounds"`
	OtherMoveProgearPounds           *int64     `json:"other_move_progear_pounds" db:"other_move_progear_pounds"`
	ServiceMemberSignature           *string    `json:"service_member_signature" db:"service_member_signature"`
	DateSigned                       *time.Time `json:"date_signed" db:"date_signed"`
	ContractorAddress                *string    `json:"contractor_address" db:"contractor_address"`
	ContractorName                   *string    `json:"contractor_name" db:"contractor_name"`
	NonavailabilityOfSignatureReason *string    `json:"nonavailability_of_signature_reason" db:"nonavailability_of_signature_reason"`
	CertifiedBySignature             *string    `json:"certified_by_signature" db:"certified_by_signature"`
	TitleOfCertifiedBySignature      *string    `json:"title_of_certified_by_signature" db:"title_of_certified_by_signature"`
}

// String is not required by pop and may be deleted
func (f Form1299) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Form1299s is not required by pop and may be deleted
type Form1299s []Form1299

// String is not required by pop and may be deleted
func (f Form1299s) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Form1299) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Form1299) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
