package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// Form1299 is an application for shipment or storage of personal property
type Form1299 struct {
	ID                                     uuid.UUID                           `db:"id"`
	CreatedAt                              time.Time                           `db:"created_at"`
	UpdatedAt                              time.Time                           `db:"updated_at"`
	DatePrepared                           *time.Time                          `db:"date_prepared"`
	ShipmentNumber                         *string                             `db:"shipment_number"`
	NameOfPreparingOffice                  *string                             `db:"name_of_preparing_office"`
	DestOfficeName                         *string                             `db:"dest_office_name"`
	OriginOfficeAddressName                *string                             `db:"origin_office_address_name"`
	OriginOfficeAddressID                  *uuid.UUID                          `db:"origin_office_address_id"`
	OriginOfficeAddress                    *Address                            `belongs_to:"address"`
	ServiceMemberFirstName                 *string                             `db:"service_member_first_name"`
	ServiceMemberMiddleInitial             *string                             `db:"service_member_middle_initial"`
	ServiceMemberLastName                  *string                             `db:"service_member_last_name"`
	ServiceMemberRank                      *internalmessages.ServiceMemberRank `db:"service_member_rank"`
	ServiceMemberSsn                       *string                             `db:"service_member_ssn"`
	ServiceMemberAgency                    *string                             `db:"service_member_agency"`
	HhgTotalPounds                         *int64                              `db:"hhg_total_pounds"`
	HhgProgearPounds                       *int64                              `db:"hhg_progear_pounds"`
	HhgValuableItemsCartons                *int64                              `db:"hhg_valuable_items_cartons"`
	MobileHomeSerialNumber                 *string                             `db:"mobile_home_serial_number"`
	MobileHomeLengthFt                     *int64                              `db:"mobile_home_length_ft"`
	MobileHomeLengthInches                 *int64                              `db:"mobile_home_length_inches"`
	MobileHomeWidthFt                      *int64                              `db:"mobile_home_width_ft"`
	MobileHomeWidthInches                  *int64                              `db:"mobile_home_width_inches"`
	MobileHomeHeightFt                     *int64                              `db:"mobile_home_height_ft"`
	MobileHomeHeightInches                 *int64                              `db:"mobile_home_height_inches"`
	MobileHomeTypeExpando                  *string                             `db:"mobile_home_type_expando"`
	MobileHomeContentsPackedRequested      bool                                `db:"mobile_home_contents_packed_requested"`
	MobileHomeBlockedRequested             bool                                `db:"mobile_home_blocked_requested"`
	MobileHomeUnblockedRequested           bool                                `db:"mobile_home_unblocked_requested"`
	MobileHomeStoredAtOriginRequested      bool                                `db:"mobile_home_stored_at_origin_requested"`
	MobileHomeStoredAtDestinationRequested bool                                `db:"mobile_home_stored_at_destination_requested"`
	StationOrdersType                      *string                             `db:"station_orders_type"`
	StationOrdersIssuedBy                  *string                             `db:"station_orders_issued_by"`
	StationOrdersNewAssignment             *string                             `db:"station_orders_new_assignment"`
	StationOrdersDate                      *time.Time                          `db:"station_orders_date"`
	StationOrdersNumber                    *string                             `db:"station_orders_number"`
	StationOrdersParagraphNumber           *string                             `db:"station_orders_paragraph_number"`
	StationOrdersInTransitTelephone        *string                             `db:"station_orders_in_transit_telephone"`
	InTransitAddressID                     *uuid.UUID                          `db:"in_transit_address_id"`
	InTransitAddress                       *Address                            `belongs_to:"address"`
	PickupAddressID                        *uuid.UUID                          `db:"pickup_address_id"`
	PickupAddress                          *Address                            `belongs_to:"address"`
	PickupTelephone                        *string                             `db:"pickup_telephone"`
	DestAddressID                          *uuid.UUID                          `db:"dest_address_id"`
	DestAddress                            *Address                            `belongs_to:"address"`
	AgentToReceiveHhg                      *string                             `db:"agent_to_receive_hhg"`
	ExtraAddressID                         *uuid.UUID                          `db:"extra_address_id"`
	ExtraAddress                           *Address                            `belongs_to:"address"`
	PackScheduledDate                      *time.Time                          `db:"pack_scheduled_date"`
	PickupScheduledDate                    *time.Time                          `db:"pickup_scheduled_date"`
	DeliveryScheduledDate                  *time.Time                          `db:"delivery_scheduled_date"`
	Remarks                                *string                             `db:"remarks"`
	OtherMove1From                         *string                             `db:"other_move_1_from"`
	OtherMove1To                           *string                             `db:"other_move_1_to"`
	OtherMove1NetPounds                    *int64                              `db:"other_move_1_net_pounds"`
	OtherMove1ProgearPounds                *int64                              `db:"other_move_1_progear_pounds"`
	OtherMove2From                         *string                             `db:"other_move_2_from"`
	OtherMove2To                           *string                             `db:"other_move_2_to"`
	OtherMove2NetPounds                    *int64                              `db:"other_move_2_net_pounds"`
	OtherMove2ProgearPounds                *int64                              `db:"other_move_2_progear_pounds"`
	ServiceMemberSignature                 *string                             `db:"service_member_signature"`
	DateSigned                             *time.Time                          `db:"date_signed"`
	ContractorAddressID                    *uuid.UUID                          `db:"contractor_address_id"`
	ContractorAddress                      *Address                            `belongs_to:"address"`
	ContractorName                         *string                             `db:"contractor_name"`
	NonavailabilityOfSignatureReason       *string                             `db:"nonavailability_of_signature_reason"`
	CertifiedBySignature                   *string                             `db:"certified_by_signature"`
	TitleOfCertifiedBySignature            *string                             `db:"title_of_certified_by_signature"`
}

// CreateForm1299WithAddresses takes a form1299 with Address structs and coordinates saving it all in a transaction
func CreateForm1299WithAddresses(dbConnection *pop.Connection, form1299 *Form1299) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {

		var transactionError error
		addressModels := []*Address{
			form1299.OriginOfficeAddress,
			form1299.InTransitAddress,
			form1299.PickupAddress,
			form1299.DestAddress,
			form1299.ExtraAddress,
			form1299.ContractorAddress,
		}

		for _, model := range addressModels {
			if model == nil {
				continue
			} else if verrs, err := dbConnection.ValidateAndCreate(model); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				transactionError = errors.New("Rollback The transaction")
				// Halt what we're doing if we get a database error
				if err != nil {
					responseError = err
					break
				}
			}
		}

		if transactionError == nil {
			form1299.OriginOfficeAddressID = GetAddressID(form1299.OriginOfficeAddress)
			form1299.InTransitAddressID = GetAddressID(form1299.InTransitAddress)
			form1299.PickupAddressID = GetAddressID(form1299.PickupAddress)
			form1299.DestAddressID = GetAddressID(form1299.DestAddress)
			form1299.ExtraAddressID = GetAddressID(form1299.ExtraAddress)
			form1299.ContractorAddressID = GetAddressID(form1299.ContractorAddress)

			if verrs, err := dbConnection.ValidateAndCreate(form1299); verrs.HasAny() || err != nil {
				transactionError = errors.New("Rollback The transaction")
				responseVErrors = verrs
				responseError = err
			}
		}

		return transactionError

	})

	return responseVErrors, responseError

}

// FetchAllForm1299s fetches all Form1299s and accompanying addresses
func FetchAllForm1299s(dbConnection *pop.Connection) (Form1299s, error) {
	var err error
	form1299s := []Form1299{}
	if err := dbConnection.Eager().All(&form1299s); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
	}
	return form1299s, err
}

// FetchForm1299ByID fetches a single Form1299 by ID and populated address fields
func FetchForm1299ByID(dbConnection *pop.Connection, id strfmt.UUID) (Form1299, error) {
	form1299 := Form1299{}
	err := dbConnection.Eager().Find(&form1299, id)
	if err != nil {
		zap.L().Error("DB Query", zap.Error(err))
	}
	return form1299, err
}

// Form1299s is not required by pop and may be deleted
type Form1299s []Form1299

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
