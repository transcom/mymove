package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// Form1299 is an application for shipment or storage of personal property
type Form1299 struct {
	ID                               uuid.UUID `json:"id" db:"id"`
	CreatedAt                        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                        time.Time `json:"updated_at" db:"updated_at"`
	DatePrepared                     time.Date `json:"date_prepared" db:"date_prepared"`
	ShipmentNumber                   string    `json:"shipment_number" db:"shipment_number"`
	NameOfPreparingOffice            string    `json:"name_of_preparing_office" db:"name_of_preparing_office"`
	DestOfficeName                   string    `json:"dest_office_name" db:"dest_office_name"`
	OriginOfficeAddressName          string    `json:"origin_office_address_name" db:"origin_office_address_name"`
	OriginOfficeAddress              string    `json:"origin_office_address" db:"origin_office_address"`
	ServiceMemberFirstName           string    `json:"service_member_first_name" db:"service_member_first_name"`
	ServiceMemberMiddleInitial       string    `json:"service_member_middle_initial" db:"service_member_middle_initial"`
	ServiceMemberLastName            string    `json:"service_member_last_name" db:"service_member_last_name"`
	ServiceMemberRank                string    `json:"service_member_rank" db:"service_member_rank"`
	ServiceMemberSsn                 string    `json:"service_member_ssn" db:"service_member_ssn"`
	ServiceMemberAgency              string    `json:"service_member_agency" db:"service_member_agency"`
	HhgTotalPounds                   int       `json:"hhg_total_pounds" db:"hhg_total_pounds"`
	HhgProgearPounds                 int       `json:"hhg_progear_pounds" db:"hhg_progear_pounds"`
	HhgValuableItemsCartons          int       `json:"hhg_valuable_items_cartons" db:"hhg_valuable_items_cartons"`
	MobileHomeSerialNumber           string    `json:"mobile_home_serial_number" db:"mobile_home_serial_number"`
	MobileHomeLengthFt               int       `json:"mobile_home_length_ft" db:"mobile_home_length_ft"`
	MobileHomeLengthInches           int       `json:"mobile_home_length_inches" db:"mobile_home_length_inches"`
	MobileHomeWidthFt                int       `json:"mobile_home_width_ft" db:"mobile_home_width_ft"`
	MobileHomeWidthInches            int       `json:"mobile_home_width_inches" db:"mobile_home_width_inches"`
	MobileHomeHeightFt               int       `json:"mobile_home_height_ft" db:"mobile_home_height_ft"`
	MobileHomeHeightInches           int       `json:"mobile_home_height_inches" db:"mobile_home_height_inches"`
	MobileHomeTypeExpando            string    `json:"mobile_home_type_expando" db:"mobile_home_type_expando"`
	MobileHomeServicesRequested      string    `json:"mobile_home_services_requested" db:"mobile_home_services_requested"`
	StationOrdersType                string    `json:"station_orders_type" db:"station_orders_type"`
	StationOrdersIssuedBy            string    `json:"station_orders_issued_by" db:"station_orders_issued_by"`
	StationOrdersNewAssignment       string    `json:"station_orders_new_assignment" db:"station_orders_new_assignment"`
	StationOrdersDate                time.Date `json:"station_orders_date" db:"station_orders_date"`
	StationOrdersNumber              string    `json:"station_orders_number" db:"station_orders_number"`
	StationOrdersParagraphNumber     string    `json:"station_orders_paragraph_number" db:"station_orders_paragraph_number"`
	StationOrdersInTransitTelephone  string    `json:"station_orders_in_transit_telephone" db:"station_orders_in_transit_telephone"`
	InTransitAddress                 string    `json:"in_transit_address" db:"in_transit_address"`
	PickupAddress                    string    `json:"pickup_address" db:"pickup_address"`
	PickupAddressMobileCourtName     string    `json:"pickup_address_mobile_court_name" db:"pickup_address_mobile_court_name"`
	PickupTelephone                  string    `json:"pickup_telephone" db:"pickup_telephone"`
	DestAddress                      string    `json:"dest_address" db:"dest_address"`
	DestAddressMobileCourtName       string    `json:"dest_address_mobile_court_name" db:"dest_address_mobile_court_name"`
	AgentToReceiveHhg                string    `json:"agent_to_receive_hhg" db:"agent_to_receive_hhg"`
	ExtraAddress                     string    `json:"extra_address" db:"extra_address"`
	PackScheduledDate                time.Date `json:"pack_scheduled_date" db:"pack_scheduled_date"`
	PickupScheduledDate              time.Date `json:"pickup_scheduled_date" db:"pickup_scheduled_date"`
	DeliveryScheduledDate            time.Date `json:"delivery_scheduled_date" db:"delivery_scheduled_date"`
	Remarks                          string    `json:"remarks" db:"remarks"`
	OtherMoveFrom                    string    `json:"other_move_from" db:"other_move_from"`
	OtherMoveTo                      string    `json:"other_move_to" db:"other_move_to"`
	OtherMoveNetPounds               int       `json:"other_move_net_pounds" db:"other_move_net_pounds"`
	OtherMoveProgearPounds           int       `json:"other_move_progear_pounds" db:"other_move_progear_pounds"`
	ServiceMemberSignature           string    `json:"service_member_signature" db:"service_member_signature"`
	DateSigned                       time.Date `json:"date_signed" db:"date_signed"`
	ContractorAddress                string    `json:"contractor_address" db:"contractor_address"`
	ContractorName                   string    `json:"contractor_name" db:"contractor_name"`
	NonavailabilityOfSignatureReason string    `json:"nonavailability_of_signature_reason" db:"nonavailability_of_signature_reason"`
	CertifiedBySignature             string    `json:"certified_by_signature" db:"certified_by_signature"`
	TitleOfCertifiedBySignature      string    `json:"title_of_certified_by_signature" db:"title_of_certified_by_signature"`
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

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Form1299) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: f.ShipmentNumber, Name: "ShipmentNumber"},
		&validators.StringIsPresent{Field: f.NameOfPreparingOffice, Name: "NameOfPreparingOffice"},
		&validators.StringIsPresent{Field: f.DestOfficeName, Name: "DestOfficeName"},
		&validators.StringIsPresent{Field: f.OriginOfficeAddressName, Name: "OriginOfficeAddressName"},
		&validators.StringIsPresent{Field: f.OriginOfficeAddress, Name: "OriginOfficeAddress"},
		&validators.StringIsPresent{Field: f.ServiceMemberFirstName, Name: "ServiceMemberFirstName"},
		&validators.StringIsPresent{Field: f.ServiceMemberMiddleInitial, Name: "ServiceMemberMiddleInitial"},
		&validators.StringIsPresent{Field: f.ServiceMemberLastName, Name: "ServiceMemberLastName"},
		&validators.StringIsPresent{Field: f.ServiceMemberRank, Name: "ServiceMemberRank"},
		&validators.StringIsPresent{Field: f.ServiceMemberSsn, Name: "ServiceMemberSsn"},
		&validators.StringIsPresent{Field: f.ServiceMemberAgency, Name: "ServiceMemberAgency"},
		&validators.IntIsPresent{Field: f.HhgTotalPounds, Name: "HhgTotalPounds"},
		&validators.IntIsPresent{Field: f.HhgProgearPounds, Name: "HhgProgearPounds"},
		&validators.IntIsPresent{Field: f.HhgValuableItemsCartons, Name: "HhgValuableItemsCartons"},
		&validators.StringIsPresent{Field: f.MobileHomeSerialNumber, Name: "MobileHomeSerialNumber"},
		&validators.IntIsPresent{Field: f.MobileHomeLengthFt, Name: "MobileHomeLengthFt"},
		&validators.IntIsPresent{Field: f.MobileHomeLengthInches, Name: "MobileHomeLengthInches"},
		&validators.IntIsPresent{Field: f.MobileHomeWidthFt, Name: "MobileHomeWidthFt"},
		&validators.IntIsPresent{Field: f.MobileHomeWidthInches, Name: "MobileHomeWidthInches"},
		&validators.IntIsPresent{Field: f.MobileHomeHeightFt, Name: "MobileHomeHeightFt"},
		&validators.IntIsPresent{Field: f.MobileHomeHeightInches, Name: "MobileHomeHeightInches"},
		&validators.StringIsPresent{Field: f.MobileHomeTypeExpando, Name: "MobileHomeTypeExpando"},
		&validators.StringIsPresent{Field: f.MobileHomeServicesRequested, Name: "MobileHomeServicesRequested"},
		&validators.StringIsPresent{Field: f.StationOrdersType, Name: "StationOrdersType"},
		&validators.StringIsPresent{Field: f.StationOrdersIssuedBy, Name: "StationOrdersIssuedBy"},
		&validators.StringIsPresent{Field: f.StationOrdersNewAssignment, Name: "StationOrdersNewAssignment"},
		&validators.StringIsPresent{Field: f.StationOrdersNumber, Name: "StationOrdersNumber"},
		&validators.StringIsPresent{Field: f.StationOrdersParagraphNumber, Name: "StationOrdersParagraphNumber"},
		&validators.StringIsPresent{Field: f.StationOrdersInTransitTelephone, Name: "StationOrdersInTransitTelephone"},
		&validators.StringIsPresent{Field: f.InTransitAddress, Name: "InTransitAddress"},
		&validators.StringIsPresent{Field: f.PickupAddress, Name: "PickupAddress"},
		&validators.StringIsPresent{Field: f.PickupAddressMobileCourtName, Name: "PickupAddressMobileCourtName"},
		&validators.StringIsPresent{Field: f.PickupTelephone, Name: "PickupTelephone"},
		&validators.StringIsPresent{Field: f.DestAddress, Name: "DestAddress"},
		&validators.StringIsPresent{Field: f.DestAddressMobileCourtName, Name: "DestAddressMobileCourtName"},
		&validators.StringIsPresent{Field: f.AgentToReceiveHhg, Name: "AgentToReceiveHhg"},
		&validators.StringIsPresent{Field: f.ExtraAddress, Name: "ExtraAddress"},
		&validators.StringIsPresent{Field: f.Remarks, Name: "Remarks"},
		&validators.StringIsPresent{Field: f.OtherMoveFrom, Name: "OtherMoveFrom"},
		&validators.StringIsPresent{Field: f.OtherMoveTo, Name: "OtherMoveTo"},
		&validators.IntIsPresent{Field: f.OtherMoveNetPounds, Name: "OtherMoveNetPounds"},
		&validators.IntIsPresent{Field: f.OtherMoveProgearPounds, Name: "OtherMoveProgearPounds"},
		&validators.StringIsPresent{Field: f.ServiceMemberSignature, Name: "ServiceMemberSignature"},
		&validators.StringIsPresent{Field: f.ContractorAddress, Name: "ContractorAddress"},
		&validators.StringIsPresent{Field: f.ContractorName, Name: "ContractorName"},
		&validators.StringIsPresent{Field: f.NonavailabilityOfSignatureReason, Name: "NonavailabilityOfSignatureReason"},
		&validators.StringIsPresent{Field: f.CertifiedBySignature, Name: "CertifiedBySignature"},
		&validators.StringIsPresent{Field: f.TitleOfCertifiedBySignature, Name: "TitleOfCertifiedBySignature"},
	), nil
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
