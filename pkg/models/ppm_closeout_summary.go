package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

type PPMCloseoutSummary struct {
	ID                          uuid.UUID   `json:"id" db:"id"`
	PPMShipmentID               uuid.UUID   `json:"ppm_shipment_id" db:"ppm_shipment_id"`
	PPMShipment                 PPMShipment `belongs_to:"ppm_shipments" fk_id:"ppm_shipment_id"`
	MaxAdvance                  *unit.Cents `json:"max_advance" db:"max_advance"`
	GTCCPaidContractedExpense   *unit.Cents `json:"gtcc_paid_contracted_expense" db:"gtcc_paid_contracted_expense"`
	MemberPaidContractedExpense *unit.Cents `json:"member_paid_contracted_expense" db:"member_paid_contracted_expense"`
	GTCCPaidPackingMaterials    *unit.Cents `json:"gtcc_paid_packing_materials" db:"gtcc_paid_packing_materials"`
	MemberPaidPackingMaterials  *unit.Cents `json:"member_paid_packing_materials" db:"member_paid_packing_materials"`
	GTCCPaidWeighingFee         *unit.Cents `json:"gtcc_paid_weighing_fee" db:"gtcc_paid_weighing_fee"`
	MemberPaidWeighingFee       *unit.Cents `json:"member_paid_weighing_fee" db:"member_paid_weighing_fee"`
	GTCCPaidRentalEquipment     *unit.Cents `json:"gtcc_paid_rental_equipment" db:"gtcc_paid_rental_equipment"`
	MemberPaidRentalEquipment   *unit.Cents `json:"member_paid_rental_equipment" db:"member_paid_rental_equipment"`
	GTCCPaidTolls               *unit.Cents `json:"gtcc_paid_tolls" db:"gtcc_paid_tolls"`
	MemberPaidTolls             *unit.Cents `json:"member_paid_tolls" db:"member_paid_tolls"`
	GTCCPaidOil                 *unit.Cents `json:"gtcc_paid_oil" db:"gtcc_paid_oil"`
	MemberPaidOil               *unit.Cents `json:"member_paid_oil" db:"member_paid_oil"`
	GTCCPaidOther               *unit.Cents `json:"gtcc_paid_other" db:"gtcc_paid_other"`
	MemberPaidOther             *unit.Cents `json:"member_paid_other" db:"member_paid_other"`
	TotalGTCCPaidExpenses       *unit.Cents `json:"total_gtcc_paid_expenses" db:"total_gtcc_paid_expenses"`
	TotalMemberPaidExpenses     *unit.Cents `json:"total_member_paid_expenses" db:"total_member_paid_expenses"`
	RemainingIncentive          *unit.Cents `json:"remaining_incentive" db:"remaining_incentive"`
	GTCCPaidSIT                 *unit.Cents `json:"gtcc_paid_sit" db:"gtcc_paid_sit"`
	MemberPaidSIT               *unit.Cents `json:"member_paid_sit" db:"member_paid_sit"`
	GTCCPaidSmallPackage        *unit.Cents `json:"gtcc_paid_small_package" db:"gtcc_paid_small_package"`
	MemberPaidSmallPackage      *unit.Cents `json:"member_paid_small_package" db:"member_paid_small_package"`
	GTCCDisbursement            *unit.Cents `json:"gtcc_disbursement" db:"gtcc_disbursement"`
	MemberDisbursement          *unit.Cents `json:"member_disbursement" db:"member_disbursement"`
	CreatedAt                   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time   `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (g PPMCloseoutSummary) TableName() string {
	return "ppm_closeouts"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PPMCloseoutSummary) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PPMShipmentID, Name: "PPMShipmentID"},
	), nil
}

func FetchPPMCloseoutByPPMID(db *pop.Connection, ppmID uuid.UUID) (PPMCloseoutSummary, error) {
	var closeout PPMCloseoutSummary
	err := db.Q().Where("ppm_shipment_id = ?", ppmID).First(&closeout)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return PPMCloseoutSummary{}, ErrFetchNotFound
		}
		return PPMCloseoutSummary{}, err
	}

	return closeout, nil
}

// CalculatePPMCloseoutSummary calls the calculate_ppm_closeout stored procedure in the DB.
func CalculatePPMCloseoutSummary(db *pop.Connection, ppmID uuid.UUID, recalculateIfExists bool) error {
	// The procedure's second parameter indicates whether to skip recalculating if a closeout already exists, the inverse of the function's recalculateIfExists parameter.
	err := db.RawQuery("CALL calculate_ppm_closeout($1, $2)", ppmID, !recalculateIfExists).Exec()
	if err != nil {
		return fmt.Errorf("error executing calculate_ppm_closeout procedure: %w", err)
	}
	return nil
}

// GetPPMCloseoutSummary calls the stored procedure to calculate/update the summary and fetches it from the DB.
func GetPPMCloseoutSummary(db *pop.Connection, ppmID uuid.UUID, recalculateIfExists bool) (PPMCloseoutSummary, error) {
	if err := CalculatePPMCloseoutSummary(db, ppmID, recalculateIfExists); err != nil {
		return PPMCloseoutSummary{}, err
	}
	closeoutSummary, err := FetchPPMCloseoutByPPMID(db, ppmID)
	if err != nil {
		return PPMCloseoutSummary{}, err
	}
	return closeoutSummary, nil
}
