package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
)

// Form1299 is an application for shipment or storage of personal property
type Form1299 struct {
	ID                                     uuid.UUID                   `json:"id" db:"id"`
	CreatedAt                              time.Time                   `json:"created_at" db:"created_at"`
	UpdatedAt                              time.Time                   `json:"updated_at" db:"updated_at"`
	DatePrepared                           *time.Time                  `json:"date_prepared" db:"date_prepared"`
	ShipmentNumber                         *string                     `json:"shipment_number" db:"shipment_number"`
	NameOfPreparingOffice                  *string                     `json:"name_of_preparing_office" db:"name_of_preparing_office"`
	DestOfficeName                         *string                     `json:"dest_office_name" db:"dest_office_name"`
	OriginOfficeAddressName                *string                     `json:"origin_office_address_name" db:"origin_office_address_name"`
	OriginOfficeAddressID                  *uuid.UUID                  `json:"origin_office_address_id" db:"origin_office_address_id"`
	OriginOfficeAddress                    *Address                    `db:"-"`
	ServiceMemberFirstName                 *string                     `json:"service_member_first_name" db:"service_member_first_name"`
	ServiceMemberMiddleInitial             *string                     `json:"service_member_middle_initial" db:"service_member_middle_initial"`
	ServiceMemberLastName                  *string                     `json:"service_member_last_name" db:"service_member_last_name"`
	ServiceMemberRank                      *messages.ServiceMemberRank `json:"service_member_rank" db:"service_member_rank"`
	ServiceMemberSsn                       *string                     `json:"service_member_ssn" db:"service_member_ssn"`
	ServiceMemberAgency                    *string                     `json:"service_member_agency" db:"service_member_agency"`
	HhgTotalPounds                         *int64                      `json:"hhg_total_pounds" db:"hhg_total_pounds"`
	HhgProgearPounds                       *int64                      `json:"hhg_progear_pounds" db:"hhg_progear_pounds"`
	HhgValuableItemsCartons                *int64                      `json:"hhg_valuable_items_cartons" db:"hhg_valuable_items_cartons"`
	MobileHomeSerialNumber                 *string                     `json:"mobile_home_serial_number" db:"mobile_home_serial_number"`
	MobileHomeLengthFt                     *int64                      `json:"mobile_home_length_ft" db:"mobile_home_length_ft"`
	MobileHomeLengthInches                 *int64                      `json:"mobile_home_length_inches" db:"mobile_home_length_inches"`
	MobileHomeWidthFt                      *int64                      `json:"mobile_home_width_ft" db:"mobile_home_width_ft"`
	MobileHomeWidthInches                  *int64                      `json:"mobile_home_width_inches" db:"mobile_home_width_inches"`
	MobileHomeHeightFt                     *int64                      `json:"mobile_home_height_ft" db:"mobile_home_height_ft"`
	MobileHomeHeightInches                 *int64                      `json:"mobile_home_height_inches" db:"mobile_home_height_inches"`
	MobileHomeTypeExpando                  *string                     `json:"mobile_home_type_expando" db:"mobile_home_type_expando"`
	MobileHomeContentsPackedRequested      bool                        `json:"mobile_home_contents_packed_requested" db:"mobile_home_contents_packed_requested"`
	MobileHomeBlockedRequested             bool                        `json:"mobile_home_blocked_requested" db:"mobile_home_blocked_requested"`
	MobileHomeUnblockedRequested           bool                        `json:"mobile_home_unblocked_requested" db:"mobile_home_unblocked_requested"`
	MobileHomeStoredAtOriginRequested      bool                        `json:"mobile_home_stored_at_origin_requested" db:"mobile_home_stored_at_origin_requested"`
	MobileHomeStoredAtDestinationRequested bool                        `json:"mobile_home_stored_at_destination_requested" db:"mobile_home_stored_at_destination_requested"`
	StationOrdersType                      *string                     `json:"station_orders_type" db:"station_orders_type"`
	StationOrdersIssuedBy                  *string                     `json:"station_orders_issued_by" db:"station_orders_issued_by"`
	StationOrdersNewAssignment             *string                     `json:"station_orders_new_assignment" db:"station_orders_new_assignment"`
	StationOrdersDate                      *time.Time                  `json:"station_orders_date" db:"station_orders_date"`
	StationOrdersNumber                    *string                     `json:"station_orders_number" db:"station_orders_number"`
	StationOrdersParagraphNumber           *string                     `json:"station_orders_paragraph_number" db:"station_orders_paragraph_number"`
	StationOrdersInTransitTelephone        *string                     `json:"station_orders_in_transit_telephone" db:"station_orders_in_transit_telephone"`
	InTransitAddressID                     *uuid.UUID                  `json:"in_transit_address_id" db:"in_transit_address_id"`
	InTransitAddress                       *Address                    `db:"-"`
	PickupAddressID                        *uuid.UUID                  `json:"pickup_address_id" db:"pickup_address_id"`
	PickupAddress                          *Address                    `db:"-"`
	PickupAddressMobileCourtName           *string                     `json:"pickup_address_mobile_court_name" db:"pickup_address_mobile_court_name"`
	PickupTelephone                        *string                     `json:"pickup_telephone" db:"pickup_telephone"`
	DestAddressID                          *uuid.UUID                  `json:"dest_address_id" db:"dest_address_id"`
	DestAddress                            *Address                    `db:"-"`
	DestAddressMobileCourtName             *string                     `json:"dest_address_mobile_court_name" db:"dest_address_mobile_court_name"`
	AgentToReceiveHhg                      *string                     `json:"agent_to_receive_hhg" db:"agent_to_receive_hhg"`
	ExtraAddressID                         *uuid.UUID                  `json:"extra_address_id" db:"extra_address_id"`
	ExtraAddress                           *Address                    `db:"-"`
	PackScheduledDate                      *time.Time                  `json:"pack_scheduled_date" db:"pack_scheduled_date"`
	PickupScheduledDate                    *time.Time                  `json:"pickup_scheduled_date" db:"pickup_scheduled_date"`
	DeliveryScheduledDate                  *time.Time                  `json:"delivery_scheduled_date" db:"delivery_scheduled_date"`
	Remarks                                *string                     `json:"remarks" db:"remarks"`
	OtherMoveFrom                          *string                     `json:"other_move_from" db:"other_move_from"`
	OtherMoveTo                            *string                     `json:"other_move_to" db:"other_move_to"`
	OtherMoveNetPounds                     *int64                      `json:"other_move_net_pounds" db:"other_move_net_pounds"`
	OtherMoveProgearPounds                 *int64                      `json:"other_move_progear_pounds" db:"other_move_progear_pounds"`
	ServiceMemberSignature                 *string                     `json:"service_member_signature" db:"service_member_signature"`
	DateSigned                             *time.Time                  `json:"date_signed" db:"date_signed"`
	ContractorAddressID                    *uuid.UUID                  `json:"contractor_address_id" db:"contractor_address_id"`
	ContractorAddress                      *Address                    `db:"-"`
	ContractorName                         *string                     `json:"contractor_name" db:"contractor_name"`
	NonavailabilityOfSignatureReason       *string                     `json:"nonavailability_of_signature_reason" db:"nonavailability_of_signature_reason"`
	CertifiedBySignature                   *string                     `json:"certified_by_signature" db:"certified_by_signature"`
	TitleOfCertifiedBySignature            *string                     `json:"title_of_certified_by_signature" db:"title_of_certified_by_signature"`
}

// CreateForm1299WithAddresses takes a form1299 with Address structs and coordinates saving it all in a transaction
func CreateForm1299WithAddresses(dbConnection *pop.Connection, form1299 *Form1299) (*validate.Errors, []error) {
	transactionVErrors := validate.NewErrors()
	transactionErrors := []error{}

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {

		saveAndPopulateErrors := func(dbStruct interface{}) bool {
			success := false
			if verrs, err := dbConnection.ValidateAndCreate(dbStruct); verrs.HasAny() || err != nil {
				transactionVErrors.Append(verrs)
				if err != nil {
					transactionErrors = append(transactionErrors, err)
				}
			} else {
				success = true
			}
			return success
		}

		if form1299.OriginOfficeAddress != nil && saveAndPopulateErrors(form1299.OriginOfficeAddress) {
			form1299.OriginOfficeAddressID = &form1299.OriginOfficeAddress.ID
		}
		if form1299.InTransitAddress != nil && saveAndPopulateErrors(form1299.InTransitAddress) {
			form1299.InTransitAddressID = &form1299.InTransitAddress.ID
		}
		if form1299.PickupAddress != nil && saveAndPopulateErrors(form1299.PickupAddress) {
			form1299.PickupAddressID = &form1299.PickupAddress.ID
		}
		if form1299.DestAddress != nil && saveAndPopulateErrors(form1299.DestAddress) {
			form1299.DestAddressID = &form1299.DestAddress.ID
		}
		if form1299.ExtraAddress != nil && saveAndPopulateErrors(form1299.ExtraAddress) {
			form1299.ExtraAddressID = &form1299.ExtraAddress.ID
		}
		if form1299.ContractorAddress != nil && saveAndPopulateErrors(form1299.ContractorAddress) {
			form1299.ContractorAddressID = &form1299.ContractorAddress.ID
		}

		saveAndPopulateErrors(form1299)

		var transactionError error
		if transactionVErrors.HasAny() || len(transactionErrors) > 0 {
			transactionError = errors.New("Rollback The transaction")
		}
		return transactionError

	})

	return transactionVErrors, transactionErrors

}

// FetchAllForm1299s fetches all Form1299s and accompanying addresses
func FetchAllForm1299s(dbConnection *pop.Connection) (Form1299s, error) {
	var err error
	form1299s := []Form1299{}
	if err := dbConnection.All(&form1299s); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
	} else {
		for i, form1299 := range form1299s {
			form1299.populateAddresses(dbConnection)
			form1299s[i] = form1299
		}
	}
	return form1299s, err
}

// FetchForm1299ByID fetches a single Form1299 by ID and populated address fields
func FetchForm1299ByID(dbConnection *pop.Connection, id strfmt.UUID) (Form1299, error) {
	form1299 := Form1299{}
	err := dbConnection.Find(&form1299, id)
	if err != nil {
		zap.L().Error("DB Query", zap.Error(err))
	} else {
		form1299.populateAddresses(dbConnection)
	}
	return form1299, err
}

// Populates address fields for form1299 structs if ID is present
func (f *Form1299) populateAddresses(dbConnection *pop.Connection) {
	if f.OriginOfficeAddressID != nil {
		f.OriginOfficeAddress = FetchAddressByID(dbConnection, f.OriginOfficeAddressID)
	}

	if f.InTransitAddressID != nil {
		f.InTransitAddress = FetchAddressByID(dbConnection, f.InTransitAddressID)
	}

	if f.PickupAddressID != nil {
		f.PickupAddress = FetchAddressByID(dbConnection, f.PickupAddressID)
	}

	if f.DestAddressID != nil {
		f.DestAddress = FetchAddressByID(dbConnection, f.DestAddressID)
	}

	if f.ExtraAddressID != nil {
		f.ExtraAddress = FetchAddressByID(dbConnection, f.ExtraAddressID)
	}

	if f.ContractorAddressID != nil {
		f.ContractorAddress = FetchAddressByID(dbConnection, f.ContractorAddressID)
	}
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
