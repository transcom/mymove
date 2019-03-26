package models

import (
	"context"
	"encoding/json"
	"time"

	beeline "github.com/honeycombio/beeline-go"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ElectronicOrdersAffiliation is a service member's military branch of service, or identifies a civilian
type ElectronicOrdersAffiliation string

const (
	// ElectronicOrdersAffiliationAirForce captures enum value "air-force"
	ElectronicOrdersAffiliationAirForce ElectronicOrdersAffiliation = "air-force"
	// ElectronicOrdersAffiliationArmy captures enum value "army"
	ElectronicOrdersAffiliationArmy ElectronicOrdersAffiliation = "army"
	// ElectronicOrdersAffiliationCivilianAgency captures enum value "civilian-agency"
	ElectronicOrdersAffiliationCivilianAgency ElectronicOrdersAffiliation = "civilian-agency"
	// ElectronicOrdersAffiliationCoastGuard captures enum value "coast-guard"
	ElectronicOrdersAffiliationCoastGuard ElectronicOrdersAffiliation = "coast-guard"
	// ElectronicOrdersAffiliationMarineCorps captures enum value "marine-corps"
	ElectronicOrdersAffiliationMarineCorps ElectronicOrdersAffiliation = "marine-corps"
	// ElectronicOrdersAffiliationNavy captures enum value "navy"
	ElectronicOrdersAffiliationNavy ElectronicOrdersAffiliation = "navy"
)

// Paygrade is the "rank" of a member. Some of these paygrades will have identical entitlements.
type Paygrade string

const (
	// PaygradeAviationCadet captures enum value "aviation-cadet"
	PaygradeAviationCadet Paygrade = "aviation-cadet"
	// PaygradeCadet captures enum value "cadet"
	PaygradeCadet Paygrade = "cadet"
	// PaygradeCivilian captures enum value "civilian"
	PaygradeCivilian Paygrade = "civilian"
	// PaygradeE1 captures enum value "e-1"
	PaygradeE1 Paygrade = "e-1"
	// PaygradeE2 captures enum value "e-2"
	PaygradeE2 Paygrade = "e-2"
	// PaygradeE3 captures enum value "e-3"
	PaygradeE3 Paygrade = "e-3"
	// PaygradeE4 captures enum value "e-4"
	PaygradeE4 Paygrade = "e-4"
	// PaygradeE5 captures enum value "e-5"
	PaygradeE5 Paygrade = "e-5"
	// PaygradeE6 captures enum value "e-6"
	PaygradeE6 Paygrade = "e-6"
	// PaygradeE7 captures enum value "e-7"
	PaygradeE7 Paygrade = "e-7"
	// PaygradeE8 captures enum value "e-8"
	PaygradeE8 Paygrade = "e-8"
	// PaygradeE9 captures enum value "e-9"
	PaygradeE9 Paygrade = "e-9"
	// PaygradeMidshipman captures enum value "midshipman"
	PaygradeMidshipman Paygrade = "midshipman"
	// PaygradeO1 captures enum value "o-1"
	PaygradeO1 Paygrade = "o-1"
	// PaygradeO2 captures enum value "o-2"
	PaygradeO2 Paygrade = "o-2"
	// PaygradeO3 captures enum value "o-3"
	PaygradeO3 Paygrade = "o-3"
	// PaygradeO4 captures enum value "o-4"
	PaygradeO4 Paygrade = "o-4"
	// PaygradeO5 captures enum value "o-5"
	PaygradeO5 Paygrade = "o-5"
	// PaygradeO6 captures enum value "o-6"
	PaygradeO6 Paygrade = "o-6"
	// PaygradeO7 captures enum value "o-7"
	PaygradeO7 Paygrade = "o-7"
	// PaygradeO8 captures enum value "o-8"
	PaygradeO8 Paygrade = "o-8"
	// PaygradeO9 captures enum value "o-9"
	PaygradeO9 Paygrade = "o-9"
	// PaygradeO10 captures enum value "o-10"
	PaygradeO10 Paygrade = "o-10"
	// PaygradeW1 captures enum value "w-1"
	PaygradeW1 Paygrade = "w-1"
	// PaygradeW2 captures enum value "w-2"
	PaygradeW2 Paygrade = "w-2"
	// PaygradeW3 captures enum value "w-3"
	PaygradeW3 Paygrade = "w-3"
	// PaygradeW4 captures enum value "w-4"
	PaygradeW4 Paygrade = "w-4"
	// PaygradeW5 captures enum value "w-5"
	PaygradeW5 Paygrade = "w-5"
)

// ElectronicOrdersStatus indicates whether these Orders are authorized, RFO (Request For Orders), or canceled. An RFO is not sufficient to authorize moving expenses; only authorized Orders can do that.
type ElectronicOrdersStatus string

const (
	// ElectronicOrdersStatusAuthorized captures enum value "authorized"
	ElectronicOrdersStatusAuthorized ElectronicOrdersStatus = "authorized"
	// ElectronicOrdersStatusRfo captures enum value "rfo"
	ElectronicOrdersStatusRfo ElectronicOrdersStatus = "rfo"
	// ElectronicOrdersStatusCanceled captures enum value "canceled"
	ElectronicOrdersStatusCanceled ElectronicOrdersStatus = "canceled"
)

// TourType indicates whether the travel is Accompanied or Unaccompanied; i.e., are dependents authorized to accompany the service member on the move. For certain OCONUS destinations, the tour type affects the member's entitlement. Otherwise, it doesn't matter.
type TourType string

const (
	// TourTypeAccompanied captures enum value "accompanied"
	TourTypeAccompanied TourType = "accompanied"
	// TourTypeUnaccompanied captures enum value "unaccompanied"
	TourTypeUnaccompanied TourType = "unaccompanied"
	// TourTypeUnaccompaniedDependentsRestricted captures enum value "unaccompanied-dependents-restricted"
	TourTypeUnaccompaniedDependentsRestricted TourType = "unaccompanied-dependents-restricted"
)

// ElectronicOrdersType is the type of travel or move for a set of Orders
type ElectronicOrdersType string

const (
	// ElectronicOrdersTypeAccession captures enum value "accession"
	ElectronicOrdersTypeAccession ElectronicOrdersType = "accession"
	// ElectronicOrdersTypeBetweenDutyStations captures enum value "between-duty-stations"
	ElectronicOrdersTypeBetweenDutyStations ElectronicOrdersType = "between-duty-stations"
	// ElectronicOrdersTypeBrac captures enum value "brac"
	ElectronicOrdersTypeBrac ElectronicOrdersType = "brac"
	// ElectronicOrdersTypeCot captures enum value "cot"
	ElectronicOrdersTypeCot ElectronicOrdersType = "cot"
	// ElectronicOrdersTypeEmergencyEvac captures enum value "emergency-evac"
	ElectronicOrdersTypeEmergencyEvac ElectronicOrdersType = "emergency-evac"
	// ElectronicOrdersTypeIpcot captures enum value "ipcot"
	ElectronicOrdersTypeIpcot ElectronicOrdersType = "ipcot"
	// ElectronicOrdersTypeLowCostTravel captures enum value "low-cost-travel"
	ElectronicOrdersTypeLowCostTravel ElectronicOrdersType = "low-cost-travel"
	// ElectronicOrdersTypeOperational captures enum value "operational"
	ElectronicOrdersTypeOperational ElectronicOrdersType = "operational"
	// ElectronicOrdersTypeOteip captures enum value "oteip"
	ElectronicOrdersTypeOteip ElectronicOrdersType = "oteip"
	// ElectronicOrdersTypeRotational captures enum value "rotational"
	ElectronicOrdersTypeRotational ElectronicOrdersType = "rotational"
	// ElectronicOrdersTypeSeparation captures enum value "separation"
	ElectronicOrdersTypeSeparation ElectronicOrdersType = "separation"
	// ElectronicOrdersTypeSpecialPurpose captures enum value "special-purpose"
	ElectronicOrdersTypeSpecialPurpose ElectronicOrdersType = "special-purpose"
	// ElectronicOrdersTypeTraining captures enum value "training"
	ElectronicOrdersTypeTraining ElectronicOrdersType = "training"
	// ElectronicOrdersTypeUnitMove captures enum value "unit-move"
	ElectronicOrdersTypeUnitMove ElectronicOrdersType = "unit-move"
)

// ElectronicOrdersRevision represents a complete amendment of one set of electronic orders
type ElectronicOrdersRevision struct {
	ID                    uuid.UUID                   `json:"id" db:"id"`
	CreatedAt             time.Time                   `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time                   `json:"updated_at" db:"updated_at"`
	ElectronicOrderID     uuid.UUID                   `json:"electronic_order_id" db:"electronic_order_id"`
	ElectronicOrder       ElectronicOrder             `belongs_to:"electronic_order"`
	SeqNum                int                         `json:"seq_num" db:"seq_num"`
	GivenName             string                      `json:"given_name" db:"given_name"`
	MiddleName            *string                     `json:"middle_name" db:"middle_name"`
	FamilyName            string                      `json:"family_name" db:"family_name"`
	NameSuffix            *string                     `json:"name_suffix" db:"name_suffix"`
	Affiliation           ElectronicOrdersAffiliation `json:"affiliation" db:"affiliation"`
	Paygrade              Paygrade                    `json:"paygrade" db:"paygrade"`
	Title                 *string                     `json:"title" db:"title"`
	Status                ElectronicOrdersStatus      `json:"status" db:"status"`
	DateIssued            time.Time                   `json:"date_issued" db:"date_issued"`
	NoCostMove            bool                        `json:"no_cost_move" db:"no_cost_move"`
	TdyEnRoute            bool                        `json:"tdy_en_route" db:"tdy_en_route"`
	TourType              TourType                    `json:"tour_type" db:"tour_type"`
	OrdersType            ElectronicOrdersType        `json:"orders_type" db:"orders_type"`
	HasDependents         bool                        `json:"has_dependents" db:"has_dependents"`
	LosingUIC             *string                     `json:"losing_uic" db:"losing_uic"`
	LosingUnitName        *string                     `json:"losing_unit_name" db:"losing_unit_name"`
	LosingUnitCity        *string                     `json:"losing_unit_city" db:"losing_unit_city"`
	LosingUnitLocality    *string                     `json:"losing_unit_locality" db:"losing_unit_locality"`
	LosingUnitCountry     *string                     `json:"losing_unit_country" db:"losing_unit_country"`
	LosingUnitPostalCode  *string                     `json:"losing_unit_postal_code" db:"losing_unit_postal_code"`
	GainingUIC            *string                     `json:"gaining_uic" db:"gaining_uic"`
	GainingUnitName       *string                     `json:"gaining_unit_name" db:"gaining_unit_name"`
	GainingUnitCity       *string                     `json:"gaining_unit_city" db:"gaining_unit_city"`
	GainingUnitLocality   *string                     `json:"gaining_unit_locality" db:"gaining_unit_locality"`
	GainingUnitCountry    *string                     `json:"gaining_unit_country" db:"gaining_unit_country"`
	GainingUnitPostalCode *string                     `json:"gaining_unit_postal_code" db:"gaining_unit_postal_code"`
	ReportNoEarlierThan   *time.Time                  `json:"report_no_earlier_than" db:"report_no_earlier_than"`
	ReportNoLaterThan     *time.Time                  `json:"report_no_later_than" db:"report_no_later_than"`
	HhgTAC                *string                     `json:"hhg_tac" db:"hhg_tac"`
	HhgSDN                *string                     `json:"hhg_sdn" db:"hhg_sdn"`
	HhgLOA                *string                     `json:"hhg_loa" db:"hhg_loa"`
	NtsTAC                *string                     `json:"nts_tac" db:"nts_tac"`
	NtsSDN                *string                     `json:"nts_sdn" db:"nts_sdn"`
	NtsLOA                *string                     `json:"nts_loa" db:"nts_loa"`
	PovShipmentTAC        *string                     `json:"pov_shipment_tac" db:"pov_shipment_tac"`
	PovShipmentSDN        *string                     `json:"pov_shipment_sdn" db:"pov_shipment_sdn"`
	PovShipmentLOA        *string                     `json:"pov_shipment_loa" db:"pov_shipment_loa"`
	PovStorageTAC         *string                     `json:"pov_storage_tac" db:"pov_storage_tac"`
	PovStorageSDN         *string                     `json:"pov_storage_sdn" db:"pov_storage_sdn"`
	PovStorageLOA         *string                     `json:"pov_storage_loa" db:"pov_storage_loa"`
	UbTAC                 *string                     `json:"ub_tac" db:"ub_tac"`
	UbSDN                 *string                     `json:"ub_sdn" db:"ub_sdn"`
	UbLOA                 *string                     `json:"ub_loa" db:"ub_loa"`
	Comments              *string                     `json:"comments" db:"comments"`
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
		&validators.StringInclusion{Field: string(e.Affiliation), Name: "Affiliation", List: []string{
			string(ElectronicOrdersAffiliationAirForce),
			string(ElectronicOrdersAffiliationArmy),
			string(ElectronicOrdersAffiliationCivilianAgency),
			string(ElectronicOrdersAffiliationCoastGuard),
			string(ElectronicOrdersAffiliationMarineCorps),
			string(ElectronicOrdersAffiliationNavy),
		}},
		&validators.StringInclusion{Field: string(e.Paygrade), Name: "Paygrade", List: []string{
			string(PaygradeAviationCadet),
			string(PaygradeCadet),
			string(PaygradeCivilian),
			string(PaygradeE1),
			string(PaygradeE2),
			string(PaygradeE3),
			string(PaygradeE4),
			string(PaygradeE5),
			string(PaygradeE6),
			string(PaygradeE7),
			string(PaygradeE8),
			string(PaygradeE9),
			string(PaygradeMidshipman),
			string(PaygradeO1),
			string(PaygradeO2),
			string(PaygradeO3),
			string(PaygradeO4),
			string(PaygradeO5),
			string(PaygradeO6),
			string(PaygradeO7),
			string(PaygradeO8),
			string(PaygradeO9),
			string(PaygradeO10),
			string(PaygradeW1),
			string(PaygradeW2),
			string(PaygradeW3),
			string(PaygradeW4),
			string(PaygradeW5),
		}},
		&StringIsNilOrNotBlank{Field: e.Title, Name: "Title"},
		&validators.StringInclusion{Field: string(e.Status), Name: "Status", List: []string{
			string(ElectronicOrdersStatusAuthorized),
			string(ElectronicOrdersStatusRfo),
			string(ElectronicOrdersStatusCanceled),
		}},
		&validators.TimeIsPresent{Field: e.DateIssued, Name: "DateIssued"},
		&validators.StringInclusion{Field: string(e.TourType), Name: "TourType", List: []string{
			string(TourTypeAccompanied),
			string(TourTypeUnaccompanied),
			string(TourTypeUnaccompaniedDependentsRestricted),
		}},
		&validators.StringInclusion{Field: string(e.OrdersType), Name: "OrdersType", List: []string{
			string(ElectronicOrdersTypeAccession),
			string(ElectronicOrdersTypeBetweenDutyStations),
			string(ElectronicOrdersTypeBrac),
			string(ElectronicOrdersTypeCot),
			string(ElectronicOrdersTypeEmergencyEvac),
			string(ElectronicOrdersTypeIpcot),
			string(ElectronicOrdersTypeLowCostTravel),
			string(ElectronicOrdersTypeOperational),
			string(ElectronicOrdersTypeOteip),
			string(ElectronicOrdersTypeRotational),
			string(ElectronicOrdersTypeSeparation),
			string(ElectronicOrdersTypeSpecialPurpose),
			string(ElectronicOrdersTypeTraining),
			string(ElectronicOrdersTypeUnitMove),
		}},
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
	verrs, responseError := dbConnection.ValidateAndCreate(revision)
	if verrs.HasAny() {
		responseVErrors.Append(verrs)
	}

	return responseVErrors, responseError
}
