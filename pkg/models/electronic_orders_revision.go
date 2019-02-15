package models

import (
	"context"
	"encoding/json"
	"time"

	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// ElectronicOrdersRevision represents a complete amendment of one set of electronic orders
type ElectronicOrdersRevision struct {
	ID                    uuid.UUID                  `json:"id" db:"id"`
	CreatedAt             time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time                  `json:"updated_at" db:"updated_at"`
	ElectronicOrderID     uuid.UUID                  `json:"electronic_order_id" db:"electronic_order_id"`
	ElectronicOrder       ElectronicOrder            `belongs_to:"electronic_order"`
	SeqNum                int                        `json:"seq_num" db:"seq_num"`
	GivenName             string                     `json:"given_name" db:"given_name"`
	MiddleName            *string                    `json:"middle_name" db:"middle_name"`
	FamilyName            string                     `json:"family_name" db:"family_name"`
	NameSuffix            *string                    `json:"name_suffix" db:"name_suffix"`
	Affiliation           ordersmessages.Affiliation `json:"affiliation" db:"affiliation"`
	Paygrade              ordersmessages.Rank        `json:"paygrade" db:"paygrade"`
	Title                 *string                    `json:"title" db:"title"`
	Status                ordersmessages.Status      `json:"status" db:"status"`
	DateIssued            time.Time                  `json:"date_issued" db:"date_issued"`
	NoCostMove            bool                       `json:"no_cost_move" db:"no_cost_move"`
	TdyEnRoute            bool                       `json:"tdy_en_route" db:"tdy_en_route"`
	TourType              ordersmessages.TourType    `json:"tour_type" db:"tour_type"`
	OrdersType            ordersmessages.OrdersType  `json:"orders_type" db:"orders_type"`
	HasDependents         bool                       `json:"has_dependents" db:"has_dependents"`
	LosingUIC             *string                    `json:"losing_uic" db:"losing_uic"`
	LosingUnitName        *string                    `json:"losing_unit_name" db:"losing_unit_name"`
	LosingUnitCity        *string                    `json:"losing_unit_city" db:"losing_unit_city"`
	LosingUnitLocality    *string                    `json:"losing_unit_locality" db:"losing_unit_locality"`
	LosingUnitCountry     *string                    `json:"losing_unit_country" db:"losing_unit_country"`
	LosingUnitPostalCode  *string                    `json:"losing_unit_postal_code" db:"losing_unit_postal_code"`
	GainingUIC            *string                    `json:"gaining_uic" db:"gaining_uic"`
	GainingUnitName       *string                    `json:"gaining_unit_name" db:"gaining_unit_name"`
	GainingUnitCity       *string                    `json:"gaining_unit_city" db:"gaining_unit_city"`
	GainingUnitLocality   *string                    `json:"gaining_unit_locality" db:"gaining_unit_locality"`
	GainingUnitCountry    *string                    `json:"gaining_unit_country" db:"gaining_unit_country"`
	GainingUnitPostalCode *string                    `json:"gaining_unit_postal_code" db:"gaining_unit_postal_code"`
	ReportNoEarlierThan   *time.Time                 `json:"report_no_earlier_than" db:"report_no_earlier_than"`
	ReportNoLaterThan     *time.Time                 `json:"report_no_later_than" db:"report_no_later_than"`
	HhgTAC                *string                    `json:"hhg_tac" db:"hhg_tac"`
	HhgSDN                *string                    `json:"hhg_sdn" db:"hhg_sdn"`
	HhgLOA                *string                    `json:"hhg_loa" db:"hhg_loa"`
	NtsTAC                *string                    `json:"nts_tac" db:"nts_tac"`
	NtsSDN                *string                    `json:"nts_sdn" db:"nts_sdn"`
	NtsLOA                *string                    `json:"nts_loa" db:"nts_loa"`
	PovShipmentTAC        *string                    `json:"pov_shipment_tac" db:"pov_shipment_tac"`
	PovShipmentSDN        *string                    `json:"pov_shipment_sdn" db:"pov_shipment_sdn"`
	PovShipmentLOA        *string                    `json:"pov_shipment_loa" db:"pov_shipment_loa"`
	PovStorageTAC         *string                    `json:"pov_storage_tac" db:"pov_storage_tac"`
	PovStorageSDN         *string                    `json:"pov_storage_sdn" db:"pov_storage_sdn"`
	PovStorageLOA         *string                    `json:"pov_storage_loa" db:"pov_storage_loa"`
	UbTAC                 *string                    `json:"ub_tac" db:"ub_tac"`
	UbSDN                 *string                    `json:"ub_sdn" db:"ub_sdn"`
	UbLOA                 *string                    `json:"ub_loa" db:"ub_loa"`
	Comments              *string                    `json:"comments" db:"comments"`
}

// String is not required by pop and may be deleted
func (e ElectronicOrdersRevision) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// ElectronicOrdersRevisions is not required by pop and may be deleted
type ElectronicOrdersRevisions []ElectronicOrdersRevision

// String is not required by pop and may be deleted
func (e ElectronicOrdersRevisions) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *ElectronicOrdersRevision) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: e.ElectronicOrderID, Name: "ElectronicOrderID"},
		&validators.IntIsGreaterThan{Field: e.SeqNum, Name: "SeqNum", Compared: -1},
		&validators.StringIsPresent{Field: e.GivenName, Name: "GivenName"},
		&StringIsNilOrNotBlank{Field: e.MiddleName, Name: "MiddleName"},
		&validators.StringIsPresent{Field: e.FamilyName, Name: "FamilyName"},
		&StringIsNilOrNotBlank{Field: e.NameSuffix, Name: "NameSuffix"},
		&validators.StringIsPresent{Field: string(e.Affiliation), Name: "Affiliation"},
		&validators.StringIsPresent{Field: string(e.Paygrade), Name: "Paygrade"},
		&StringIsNilOrNotBlank{Field: e.Title, Name: "Title"},
		&validators.StringIsPresent{Field: string(e.Status), Name: "Status"},
		&validators.TimeIsPresent{Field: e.DateIssued, Name: "DateIssued"},
		&validators.StringIsPresent{Field: string(e.TourType), Name: "TourType"},
		&validators.StringIsPresent{Field: string(e.OrdersType), Name: "OrdersType"},
		&StringIsNilOrNotBlank{Field: e.LosingUIC, Name: "LosingUIC"},
		&StringIsNilOrNotBlank{Field: e.LosingUnitName, Name: "LosingUnitName"},
		&StringIsNilOrNotBlank{Field: e.LosingUnitCity, Name: "LosingUnitCity"},
		&StringIsNilOrNotBlank{Field: e.LosingUnitLocality, Name: "LosingUnitLocality"},
		&StringIsNilOrNotBlank{Field: e.LosingUnitCountry, Name: "LosingUnitCountry"},
		&StringIsNilOrNotBlank{Field: e.LosingUnitPostalCode, Name: "LosingUnitPostalCode"},
		&StringIsNilOrNotBlank{Field: e.GainingUIC, Name: "GainingUIC"},
		&StringIsNilOrNotBlank{Field: e.GainingUnitName, Name: "GainingUnitName"},
		&StringIsNilOrNotBlank{Field: e.GainingUnitCity, Name: "GainingUnitCity"},
		&StringIsNilOrNotBlank{Field: e.GainingUnitLocality, Name: "GainingUnitLocality"},
		&StringIsNilOrNotBlank{Field: e.GainingUnitCountry, Name: "GainingUnitCountry"},
		&StringIsNilOrNotBlank{Field: e.GainingUnitPostalCode, Name: "GainingUnitPostalCode"},
		&StringIsNilOrNotBlank{Field: e.HhgTAC, Name: "HhgTAC"},
		&StringIsNilOrNotBlank{Field: e.HhgSDN, Name: "HhgSDN"},
		&StringIsNilOrNotBlank{Field: e.HhgLOA, Name: "HhgLOA"},
		&StringIsNilOrNotBlank{Field: e.NtsTAC, Name: "NtsTAC"},
		&StringIsNilOrNotBlank{Field: e.NtsSDN, Name: "NtsSDN"},
		&StringIsNilOrNotBlank{Field: e.NtsLOA, Name: "NtsLOA"},
		&StringIsNilOrNotBlank{Field: e.PovShipmentTAC, Name: "PovShipmentTAC"},
		&StringIsNilOrNotBlank{Field: e.PovShipmentSDN, Name: "PovShipmentSDN"},
		&StringIsNilOrNotBlank{Field: e.PovShipmentLOA, Name: "PovShipmentLOA"},
		&StringIsNilOrNotBlank{Field: e.PovStorageTAC, Name: "PovStorageTAC"},
		&StringIsNilOrNotBlank{Field: e.PovStorageSDN, Name: "PovStorageSDN"},
		&StringIsNilOrNotBlank{Field: e.PovStorageLOA, Name: "PovStorageLOA"},
		&StringIsNilOrNotBlank{Field: e.UbTAC, Name: "UbTAC"},
		&StringIsNilOrNotBlank{Field: e.UbSDN, Name: "UbSDN"},
		&StringIsNilOrNotBlank{Field: e.UbLOA, Name: "UbLOA"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *ElectronicOrdersRevision) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *ElectronicOrdersRevision) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// CreateElectronicOrdersRevision inserts a revision into the database
func CreateElectronicOrdersRevision(ctx context.Context, dbConnection *pop.Connection, revision *ElectronicOrdersRevision) (*validate.Errors, error) {
	ctx, span := beeline.StartSpan(ctx, "CreateElectronicOrdersRevision")
	defer span.Send()

	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")
		if verrs, err := dbConnection.ValidateAndCreate(revision); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError
}
