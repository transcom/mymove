package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// MovingExpenseReceiptType represents types of different moving expenses
type MovingExpenseReceiptType string

const (
	// MovingExpenseReceiptTypeContractedExpense captures enum value "CONTRACTED_EXPENSE"
	MovingExpenseReceiptTypeContractedExpense MovingExpenseReceiptType = "CONTRACTED_EXPENSE"
	// MovingExpenseReceiptTypeOil captures enum value "OIL"
	MovingExpenseReceiptTypeOil MovingExpenseReceiptType = "OIL"
	// MovingExpenseReceiptTypePackingMaterials captures enum value "PACKING_MATERIALS"
	MovingExpenseReceiptTypePackingMaterials MovingExpenseReceiptType = "PACKING_MATERIALS"
	// MovingExpenseReceiptTypeRentalEquipment captures enum value "RENTAL_EQUIPMENT"
	MovingExpenseReceiptTypeRentalEquipment MovingExpenseReceiptType = "RENTAL_EQUIPMENT"
	// MovingExpenseReceiptTypeStorage captures enum value "STORAGE"
	MovingExpenseReceiptTypeStorage MovingExpenseReceiptType = "STORAGE"
	// MovingExpenseReceiptTypeTolls captures enum value "TOLLS"
	MovingExpenseReceiptTypeTolls MovingExpenseReceiptType = "TOLLS"
	// MovingExpenseReceiptTypeWeighingFee captures enum value "WEIGHING_FEE"
	MovingExpenseReceiptTypeWeighingFee MovingExpenseReceiptType = "WEIGHING_FEE"
	// MovingExpenseReceiptTypeOther captures enum value "OTHER"
	MovingExpenseReceiptTypeOther MovingExpenseReceiptType = "OTHER"
)

var AllowedExpenseTypes = []string{
	string(MovingExpenseReceiptTypeContractedExpense),
	string(MovingExpenseReceiptTypeOil),
	string(MovingExpenseReceiptTypePackingMaterials),
	string(MovingExpenseReceiptTypeRentalEquipment),
	string(MovingExpenseReceiptTypeStorage),
	string(MovingExpenseReceiptTypeTolls),
	string(MovingExpenseReceiptTypeWeighingFee),
	string(MovingExpenseReceiptTypeOther),
}

type MovingExpense struct {
	ID                uuid.UUID                 `json:"id" db:"id"`
	PPMShipmentID     uuid.UUID                 `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment       PPMShipment               `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	DocumentID        uuid.UUID                 `json:"document_id" db:"document_id"`
	Document          Document                  `belongs_to:"documents" fk_id:"document_id"`
	CreatedAt         time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time                `json:"deleted_at" db:"deleted_at"`
	MovingExpenseType *MovingExpenseReceiptType `json:"moving_expense_type" db:"moving_expense_type"`
	Description       *string                   `json:"description" db:"description"`
	PaidWithGTCC      *bool                     `json:"paid_with_gtcc" db:"paid_with_gtcc"`
	Amount            *unit.Cents               `json:"amount" db:"amount"`
	MissingReceipt    *bool                     `json:"missing_receipt" db:"missing_receipt"`
	Status            *PPMDocumentStatus        `json:"status" db:"status"`
	Reason            *string                   `json:"reason" db:"reason"`
	SITStartDate      *time.Time                `json:"sit_start_date" db:"sit_start_date"`
	SITEndDate        *time.Time                `json:"sit_end_date" db:"sit_end_date"`
	WeightStored      *unit.Pound               `json:"weight_stored" db:"weight_stored"`
}

// TableName overrides the table name used by Pop.
func (m MovingExpense) TableName() string {
	return "moving_expenses"
}

type MovingExpenses []MovingExpense

func (e MovingExpenses) FilterDeleted() MovingExpenses {
	if len(e) == 0 {
		return e
	}

	nonDeletedExpenses := MovingExpenses{}
	for _, expense := range e {
		if expense.DeletedAt == nil {
			nonDeletedExpenses = append(nonDeletedExpenses, expense)
		}
	}

	return nonDeletedExpenses
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (m *MovingExpense) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "PPMShipmentID", Field: m.PPMShipmentID},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: m.DeletedAt},
		&OptionalStringInclusion{Name: "MovingExpenseType", Field: (*string)(m.MovingExpenseType), List: AllowedExpenseTypes},
		&StringIsNilOrNotBlank{Name: "Description", Field: m.Description},
		&OptionalStringInclusion{Name: "Status", Field: (*string)(m.Status), List: AllowedPPMDocumentStatuses},
		&StringIsNilOrNotBlank{Name: "Reason", Field: m.Reason},
		&OptionalTimeIsPresent{Name: "SITStartDate", Field: m.SITStartDate},
		&OptionalTimeIsPresent{Name: "SITEndDate", Field: m.SITEndDate},
	), nil
}
