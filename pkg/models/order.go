package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// UploadedOrdersDocumentName is the name of an uploaded orders document
const UploadedOrdersDocumentName = "uploaded_orders"

const SupplyAndServicesCostEstimate string = "Prices for services under this task order will be in accordance with rates provided in GHC Attachment 2 - Pricing Rate Table. It is the responsibility of the contractor to provide the estimated weight quantity to apply to services on this task order, when applicable (See Attachment 1 - Performance Work Statement)."
const MethodOfPayment string = "Payment will be made using the Third-Party Payment System (TPPS) Automated Payment System"
const NAICS string = "488510 - FREIGHT TRANSPORTATION ARRANGEMENT"
const InstructionsBeforeContractNumber string = "Packaging, packing, and shipping instructions as identified in the Conformed Copy of"
const InstructionsAfterContractNumber string = "Attachment 1 Performance Work Statement"

// OrderStatus represents the state of an order record in the UX manual orders flow
type OrderStatus string

const (
	// OrderStatusDRAFT captures enum value "DRAFT"
	OrderStatusDRAFT OrderStatus = "DRAFT"
	// OrderStatusSUBMITTED captures enum value "SUBMITTED"
	OrderStatusSUBMITTED OrderStatus = "SUBMITTED"
	// OrderStatusAPPROVED captures enum value "APPROVED"
	OrderStatusAPPROVED OrderStatus = "APPROVED"
	// OrderStatusCANCELED captures enum value "CANCELED"
	OrderStatusCANCELED OrderStatus = "CANCELED"
)

// DepartmentIndicator represents an order's department indicator
type DepartmentIndicator string

// String is a string representation of a DepartmentIndicator
func (d DepartmentIndicator) String() string {
	return string(d)
}

const (
	// DepartmentIndicatorARMY captures enum value "ARMY"
	DepartmentIndicatorARMY DepartmentIndicator = "ARMY"
	// DepartmentIndicatorARMYCORPSOFENGINEERS captures enum value "ARMY_CORPS_OF_ENGINEERS"
	DepartmentIndicatorARMYCORPSOFENGINEERS DepartmentIndicator = "ARMY_CORPS_OF_ENGINEERS"
	// DepartmentIndicatorCOASTGUARD captures enum value "COAST_GUARD"
	DepartmentIndicatorCOASTGUARD DepartmentIndicator = "COAST_GUARD"
	// DepartmentIndicatorNAVYANDMARINES captures enum value "NAVY_AND_MARINES"
	DepartmentIndicatorNAVYANDMARINES DepartmentIndicator = "NAVY_AND_MARINES"
	// DepartmentIndicatorAIRANDSPACEFORCE captures enum value "AIR_AND_SPACE_FORCE"
	DepartmentIndicatorAIRANDSPACEFORCE DepartmentIndicator = "AIR_AND_SPACE_FORCE"
	// DepartmentIndicatorOFFICEOFSECRETARYOFDEFENSE captures enum value "OFFICE_OF_SECRETARY_OF_DEFENSE"
	DepartmentIndicatorOFFICEOFSECRETARYOFDEFENSE DepartmentIndicator = "OFFICE_OF_SECRETARY_OF_DEFENSE"
)

// Order is a set of orders received by a service member
type Order struct {
	ID                             uuid.UUID                          `json:"id" db:"id"`
	CreatedAt                      time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time                          `json:"updated_at" db:"updated_at"`
	ServiceMemberID                uuid.UUID                          `json:"service_member_id" db:"service_member_id"`
	ServiceMember                  ServiceMember                      `json:"service_member" belongs_to:"service_members" fk_id:"service_member_id"`
	IssueDate                      time.Time                          `json:"issue_date" db:"issue_date"`
	ReportByDate                   time.Time                          `json:"report_by_date" db:"report_by_date"`
	OrdersType                     internalmessages.OrdersType        `json:"orders_type" db:"orders_type"`
	OrdersTypeDetail               *internalmessages.OrdersTypeDetail `json:"orders_type_detail" db:"orders_type_detail"`
	HasDependents                  bool                               `json:"has_dependents" db:"has_dependents"`
	SpouseHasProGear               bool                               `json:"spouse_has_pro_gear" db:"spouse_has_pro_gear"`
	OriginDutyLocation             *DutyLocation                      `json:"origin_duty_location" belongs_to:"duty_locations" fk_id:"origin_duty_location_id"`
	OriginDutyLocationID           *uuid.UUID                         `json:"origin_duty_location_id" db:"origin_duty_location_id"`
	NewDutyLocationID              uuid.UUID                          `json:"new_duty_location_id" db:"new_duty_location_id"`
	NewDutyLocation                DutyLocation                       `json:"new_duty_location" belongs_to:"duty_locations" fk_id:"new_duty_location_id"`
	DestinationGBLOC               *string                            `json:"destination_duty_location_gbloc" db:"destination_gbloc"`
	UploadedOrders                 Document                           `belongs_to:"documents" fk_id:"uploaded_orders_id"`
	UploadedOrdersID               uuid.UUID                          `json:"uploaded_orders_id" db:"uploaded_orders_id"`
	OrdersNumber                   *string                            `json:"orders_number" db:"orders_number"`
	Moves                          Moves                              `has_many:"moves" fk_id:"orders_id" order_by:"created_at desc"`
	Status                         OrderStatus                        `json:"status" db:"status"`
	TAC                            *string                            `json:"tac" db:"tac"`
	SAC                            *string                            `json:"sac" db:"sac"`
	NtsTAC                         *string                            `json:"nts_tac" db:"nts_tac"`
	NtsSAC                         *string                            `json:"nts_sac" db:"nts_sac"`
	DepartmentIndicator            *string                            `json:"department_indicator" db:"department_indicator"`
	Grade                          *internalmessages.OrderPayGrade    `json:"grade" db:"grade"`
	Entitlement                    *Entitlement                       `belongs_to:"entitlements" fk_id:"entitlement_id"`
	EntitlementID                  *uuid.UUID                         `json:"entitlement_id" db:"entitlement_id"`
	UploadedAmendedOrders          *Document                          `belongs_to:"documents" fk_id:"uploaded_amended_orders_id"`
	UploadedAmendedOrdersID        *uuid.UUID                         `json:"uploaded_amended_orders_id" db:"uploaded_amended_orders_id"`
	AmendedOrdersAcknowledgedAt    *time.Time                         `json:"amended_orders_acknowledged_at" db:"amended_orders_acknowledged_at"`
	OriginDutyLocationGBLOC        *string                            `json:"origin_duty_location_gbloc" db:"gbloc"`
	SupplyAndServicesCostEstimate  string                             `json:"supply_and_services_cost_estimate" db:"supply_and_services_cost_estimate"`
	PackingAndShippingInstructions string                             `json:"packing_and_shipping_instructions" db:"packing_and_shipping_instructions"`
	MethodOfPayment                string                             `json:"method_of_payment" db:"method_of_payment"`
	NAICS                          string                             `json:"naics" db:"naics"`
	ProvidesServicesCounseling     *bool                              `belongs_to:"duty_locations" fk_id:"origin_duty_location_id"`
	RankID                         *uuid.UUID                         `db:"rank_id"`
	Rank                           *Rank                              `belongs_to:"ranks" fk_id:"rank_id"`
}

// TableName overrides the table name used by Pop.
func (o Order) TableName() string {
	return "orders"
}

type Orders []Order

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *Order) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&OrdersTypeIsPresent{Field: o.OrdersType, Name: "OrdersType"},
		&validators.TimeIsPresent{Field: o.IssueDate, Name: "IssueDate"},
		&validators.TimeIsPresent{Field: o.ReportByDate, Name: "ReportByDate"},
		&validators.UUIDIsPresent{Field: o.ServiceMemberID, Name: "ServiceMemberID"},
		&validators.UUIDIsPresent{Field: o.NewDutyLocationID, Name: "NewDutyLocationID"},
		&StringIsNilOrNotBlank{Field: o.DestinationGBLOC, Name: "DestinationDutyLocationGBLOC"},
		&validators.StringIsPresent{Field: string(o.Status), Name: "Status"},
		&StringIsNilOrNotBlank{Field: o.TAC, Name: "TransportationAccountingCode"},
		&StringIsNilOrNotBlank{Field: o.SAC, Name: "SAC"},
		&StringIsNilOrNotBlank{Field: o.NtsTAC, Name: "NtsTAC"},
		&StringIsNilOrNotBlank{Field: o.NtsSAC, Name: "NtsSAC"},
		&StringIsNilOrNotBlank{Field: o.DepartmentIndicator, Name: "DepartmentIndicator"},
		&CannotBeTrueIfFalse{Field1: o.SpouseHasProGear, Name1: "SpouseHasProGear", Field2: o.HasDependents, Name2: "HasDependents"},
		&OptionalUUIDIsPresent{Field: o.EntitlementID, Name: "EntitlementID"},
		&OptionalUUIDIsPresent{Field: o.OriginDutyLocationID, Name: "OriginDutyLocationID"},
		&OptionalRegexMatch{Name: "TransportationAccountingCode", Field: o.TAC, Expr: `\A([A-Za-z0-9]){4}\z`, Message: "TAC must be exactly 4 alphanumeric characters."},
		// block * and " in TAC & SAC fields
		&OptionalRegexMatch{
			Name:    "TransportationAccountingCode",
			Field:   o.TAC,
			Expr:    `\A[^*"]*\z`,
			Message: `TAC cannot contain * or " characters.`,
		},
		&OptionalRegexMatch{
			Name:    "SAC",
			Field:   o.SAC,
			Expr:    `\A[^*"]*\z`,
			Message: `SAC cannot contain * or " characters.`,
		},
		&OptionalRegexMatch{
			Name:    "NtsTAC",
			Field:   o.NtsTAC,
			Expr:    `\A[^*"]*\z`,
			Message: `NTS TAC cannot contain * or " characters.`,
		},
		&OptionalRegexMatch{
			Name:    "NtsSAC",
			Field:   o.NtsSAC,
			Expr:    `\A[^*"]*\z`,
			Message: `NTS SAC cannot contain * or " characters.`,
		},
		&validators.UUIDIsPresent{Field: o.UploadedOrdersID, Name: "UploadedOrdersID"},
		&OptionalUUIDIsPresent{Field: o.UploadedAmendedOrdersID, Name: "UploadedAmendedOrdersID"},
		&StringIsNilOrNotBlank{Field: o.OriginDutyLocationGBLOC, Name: "OriginDutyLocationGBLOC"},
		&validators.StringIsPresent{Field: o.SupplyAndServicesCostEstimate, Name: "SupplyAndServicesCostEstimate"},
		&validators.StringIsPresent{Field: o.PackingAndShippingInstructions, Name: "PackingAndShippingInstructions"},
		&validators.StringIsPresent{Field: o.MethodOfPayment, Name: "MethodOfPayment"},
		&validators.StringIsPresent{Field: o.NAICS, Name: "NAICS"},
	), nil
}

// SaveOrder saves an order
func SaveOrder(db *pop.Connection, order *Order) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionErr := db.Transaction(func(_ *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if verrs, err := db.ValidateAndSave(order); verrs.HasAny() || err != nil {
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

// State Machine
// Avoid calling Order.Status = ... ever. Use these methods to change the state.

// Submit submits the Order
func (o *Order) Submit() error {
	if o.Status != OrderStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Submit")
	}

	o.Status = OrderStatusSUBMITTED
	return nil
}

// Cancel cancels the Order
func (o *Order) Cancel() error {
	if o.Status == OrderStatusCANCELED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	o.Status = OrderStatusCANCELED
	return nil
}

var ordersTypeToAllotmentAppParamName = map[internalmessages.OrdersType]string{
	internalmessages.OrdersTypeSTUDENTTRAVEL: "studentTravelHhgAllowance",
}

// Helper func to enforce strict unmarshal of application param values into a given  interface
func strictUnmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	// Fail on unknown fields
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// FetchOrderForUser returns orders only if it is allowed for the given user to access those orders.
func FetchOrderForUser(db *pop.Connection, session *auth.Session, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().EagerPreload("ServiceMember.User",
		"OriginDutyLocation.Address",
		"OriginDutyLocation.TransportationOffice",
		"NewDutyLocation.Address",
		"NewDutyLocation.TransportationOffice",
		"UploadedOrders",
		"UploadedAmendedOrders",
		"Moves.SignedCertifications",
		"Moves.CloseoutOffice.Address",
		"Entitlement",
		"OriginDutyLocation",
		"OriginDutyLocation.ProvidesServicesCounseling").
		Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}

	// Conduct allotment lookup
	if order.Entitlement != nil && order.Grade != nil {
		switch order.OrdersType {
		case internalmessages.OrdersTypeSTUDENTTRAVEL:
			// We currently store orders allotment overrides as an application parameter
			// as it is a current one-off use case introduced by E-06189
			var jsonData json.RawMessage
			if paramName, ok := ordersTypeToAllotmentAppParamName[order.OrdersType]; ok {
				err := db.RawQuery(`
				SELECT parameter_json
				FROM application_parameters
				WHERE parameter_name = $1
				`, paramName).First(&jsonData)

				if err != nil {
					return Order{}, fmt.Errorf("failed to fetch weight allotment for orders type %s: %w", order.OrdersType, err)
				}

				// Convert the JSON data to the WeightAllotment struct
				err = strictUnmarshal(jsonData, &order.Entitlement.WeightAllotted)
				if err != nil {
					return Order{}, fmt.Errorf("failed to parse weight allotment JSON for orders type %s: %w", order.OrdersType, err)
				}
			}
		default:
			var hhgAllowance HHGAllowance
			err = db.RawQuery(`
				SELECT hhg_allowances.*
				FROM hhg_allowances
				INNER JOIN pay_grades ON hhg_allowances.pay_grade_id = pay_grades.id
				WHERE pay_grades.grade = $1
				LIMIT 1
			`, order.Grade).First(&hhgAllowance)
			if err != nil {
				return Order{}, fmt.Errorf("failed to parse weight allotment JSON for orders type %s: %w", order.OrdersType, err)
			}
			order.Entitlement.WeightAllotted = &WeightAllotment{
				TotalWeightSelf:               hhgAllowance.TotalWeightSelf,
				TotalWeightSelfPlusDependents: hhgAllowance.TotalWeightSelfPlusDependents,
				ProGearWeight:                 hhgAllowance.ProGearWeight,
				ProGearWeightSpouse:           hhgAllowance.ProGearWeightSpouse,
			}

		}
	}

	// TODO: Handle case where more than one user is authorized to modify orders
	if session.IsMilApp() && order.ServiceMember.ID != session.ServiceMemberID {
		return Order{}, ErrFetchForbidden
	}

	// Eager loading of nested has_many associations is broken
	var userUploads UserUploads
	err = db.Q().
		Scope(utilities.ExcludeDeletedScope()).EagerPreload("Upload").
		Where("document_id = ?", order.UploadedOrders.ID).
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

	return order, nil
}

// Fetch order containing only base amendment information
func FetchOrderAmendmentsInfo(db *pop.Connection, session *auth.Session, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().EagerPreload("ServiceMember.User",
		"UploadedAmendedOrders").
		Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}

	if session != nil && session.IsMilApp() && order.ServiceMember.ID != session.ServiceMemberID {
		return Order{}, ErrFetchForbidden
	}

	if order.UploadedAmendedOrders != nil {
		var amendedUserUploads UserUploads
		err = db.Q().
			Where("document_id = ?", order.UploadedAmendedOrdersID).
			All(&amendedUserUploads)
		if err != nil {
			return Order{}, err
		}
		order.UploadedAmendedOrders.UserUploads = amendedUserUploads
	}

	return order, nil
}

// FetchOrder returns orders without REGARDLESS OF USER.
// DO NOT USE IF YOU NEED USER AUTH
func FetchOrder(db *pop.Connection, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}

	return order, nil
}

// FetchOrderForPDFConversion returns orders and any attached uploads
func FetchOrderForPDFConversion(db *pop.Connection, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Eager("UploadedOrders.UserUploads.Upload").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}
	return order, nil
}

// CreateNewMove creates a move associated with these Orders
func (o *Order) CreateNewMove(db *pop.Connection, moveOptions MoveOptions) (*Move, *validate.Errors, error) {
	return createNewMove(db, *o, moveOptions)
}

/*
 * GetOriginPostalCode returns the GBLOC for the postal code of the the origin duty location of the order.
 */
func (o Order) GetOriginPostalCode(db *pop.Connection) (string, error) {
	// Since this requires looking up the order in the DB, the order must have an ID. This means, the order has to have been created first.
	if uuid.UUID.IsNil(o.ID) {
		return "", errors.WithMessage(ErrInvalidOrderID, "You must create the order in the DB before getting the origin GBLOC.")
	}

	err := db.Load(&o, "OriginDutyLocation.Address")
	if err != nil {
		if err.Error() == RecordNotFoundErrorString {
			return "", errors.WithMessage(err, "No Origin Duty Location was found for the order ID "+o.ID.String())
		}
		return "", err
	}

	return o.OriginDutyLocation.Address.PostalCode, nil
}

/*
 * GetOriginGBLOC returns the GBLOC for the postal code of the the origin duty location of the order.
 */
func (o Order) GetOriginGBLOC(db *pop.Connection) (string, error) {
	// Since this requires looking up the order in the DB, the order must have an ID. This means, the order has to have been created first.
	if uuid.UUID.IsNil(o.ID) {
		return "", errors.WithMessage(ErrInvalidOrderID, "You must created the order in the DB before getting the destination GBLOC.")
	}

	originPostalCode, err := o.GetOriginPostalCode(db)
	if err != nil {
		return "", err
	}

	var originGBLOC PostalCodeToGBLOC
	originGBLOC, err = FetchGBLOCForPostalCode(db, originPostalCode)
	if err != nil {
		return "", err
	}

	return originGBLOC.GBLOC, nil
}

// IsComplete checks if orders have all fields necessary to approve a move
func (o *Order) IsComplete() bool {

	if o.OrdersTypeDetail == nil {
		return false
	}
	if o.OrdersNumber == nil {
		return false
	}
	if o.DepartmentIndicator == nil {
		return false
	}
	// HHG TAC
	if o.TAC == nil {
		return false
	}

	return true
}

// FetchAllShipmentsExcludingRejected returns all the shipments associated with an order excluding rejected shipments
func (o Order) FetchAllShipmentsExcludingRejected(db *pop.Connection) (map[uuid.UUID]MTOShipments, error) {
	// Since this requires looking up the order in the DB, the order must have an ID. This means, the order has to have been created first.
	if uuid.UUID.IsNil(o.ID) {
		return nil, errors.WithMessage(ErrInvalidOrderID, "You must created the order in the DB before fetching associated shipments.")
	}

	var err error

	err = db.Load(&o, "Moves")
	if err != nil {
		if err.Error() == RecordNotFoundErrorString {
			return nil, errors.WithMessage(err, "No Moves were found for the order ID "+o.ID.String())
		}
		return nil, errors.WithMessage(err, "Could not load moves for order ID "+o.ID.String())
	}

	// shipmentsMap is a map of key, value pairs where the key is the move id and the value is a list of associated MTOShipments
	shipmentsMap := make(map[uuid.UUID]MTOShipments)

	for _, m := range o.Moves {
		var shipments MTOShipments
		err = db.Load(&m, "MTOShipments")
		if err != nil {
			return nil, errors.WithMessage(err, "Could not load shipments for move "+m.ID.String())
		}

		for _, s := range m.MTOShipments {
			err = db.Load(&s, "Status", "DeletedAt", "CreatedAt", "DestinationAddress")
			if err != nil {
				return nil, errors.WithMessage(err, "Could not load shipment with ID of "+s.ID.String()+" for move ID "+m.ID.String())
			}

			if s.Status != MTOShipmentStatusRejected && s.Status != MTOShipmentStatusCanceled && s.DeletedAt == nil {
				shipments = append(shipments, s)
			}
		}

		sort.Slice(shipments, func(i, j int) bool {
			return shipments[i].CreatedAt.Before(shipments[j].CreatedAt)
		})

		shipmentsMap[m.ID] = shipments
	}

	return shipmentsMap, nil
}

/*
* GetDestinationAddressForAssociatedMoves returns the destination Address for the first shipments from each of
* the moves that are associated with an order. If there are no shipments returned, it will return the
* Address of the new duty station addresses.
 */
func (o Order) GetDestinationAddressForAssociatedMoves(db *pop.Connection) (*Address, error) {
	if uuid.UUID.IsNil(o.ID) {
		return nil, errors.WithMessage(ErrInvalidOrderID, "You must create the order in the DB before getting the destination Address.")
	}

	err := db.Load(&o, "Moves", "NewDutyLocation.Address")
	if err != nil {
		if err.Error() == RecordNotFoundErrorString {
			return nil, errors.WithMessage(err, "No Moves were found for the order ID "+o.ID.String())
		}
		return nil, err
	}

	// addrMap is a map of key, value pairs where the key is the move id and the value is the destination address
	var destinationAddress Address
	for i, m := range o.Moves {
		err = db.Load(&o.Moves[i], "MTOShipments")
		if err != nil {
			if err.Error() == RecordNotFoundErrorString {
				return nil, errors.WithMessage(err, "Could not find shipments for move "+m.ID.String())
			}
			return nil, err
		}

		var shipments MTOShipments
		for j, s := range o.Moves[i].MTOShipments {
			err = db.Load(&o.Moves[i].MTOShipments[j], "CreatedAt", "Status", "DeletedAt", "DestinationAddress", "ShipmentType", "PPMShipment", "PPMShipment.Status", "PPMShipment.DestinationAddress")
			if err != nil {
				if err.Error() == RecordNotFoundErrorString {
					return nil, errors.WithMessage(err, "Could not load shipment with ID of "+s.ID.String()+" for move ID "+m.ID.String())
				}
				return nil, err
			}

			if o.Moves[i].MTOShipments[j].Status != MTOShipmentStatusRejected &&
				o.Moves[i].MTOShipments[j].Status != MTOShipmentStatusCanceled &&
				o.Moves[i].MTOShipments[j].ShipmentType != MTOShipmentTypeHHGIntoNTS &&
				o.Moves[i].MTOShipments[j].DeletedAt == nil {
				shipments = append(shipments, o.Moves[i].MTOShipments[j])
			}
		}

		// If we have valid shipments, use the first one's destination address
		if len(shipments) > 0 {
			sort.Slice(shipments, func(i, j int) bool {
				return shipments[i].CreatedAt.Before(shipments[j].CreatedAt)
			})

			addressResult, err := shipments[0].GetDestinationAddress(db)
			if err != nil {
				return nil, err
			}

			if addressResult != nil {
				destinationAddress = *addressResult
			} else {
				return nil, errors.WithMessage(ErrMissingDestinationAddress, "No destination address was able to be found for the order ID "+o.ID.String())
			}
		} else {
			// No valid shipments, use new duty location
			destinationAddress = o.NewDutyLocation.Address
		}
	}

	return &destinationAddress, nil
}

// IsCompleteForGBL checks if orders have all fields necessary to generate a GBL
func (o *Order) IsCompleteForGBL() bool {

	if o.DepartmentIndicator == nil {
		return false
	}
	// HHG TAC
	if o.TAC == nil {
		return false
	}
	return true
}

func (o *Order) CanSendEmailWithOrdersType() bool {
	if o.OrdersType == internalmessages.OrdersTypeBLUEBARK || o.OrdersType == internalmessages.OrdersTypeSAFETY {
		return false
	}
	return true
}
