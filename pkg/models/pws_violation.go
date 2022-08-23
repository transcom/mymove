package models

type PWSViolationCategory string

const (
	PWSViolationCategoryPreMoveServices      PWSViolationCategory = "Pre-Move Services"
	PWSViolationCategoryPhysicalMoveServices PWSViolationCategory = "Physical Move Services"
	PWSViolationCategoryLiability            PWSViolationCategory = "Liability"
)

// type PWSViolationSubCategory string

// const (
// 	PWSViolationSubCategoryCustomerSupport    PWSViolationSubCategory = "Customer Support"
// 	PWSViolationSubCategoryCounseling         PWSViolationSubCategory = "Counseling"
// 	PWSViolationSubCategoryWeightEstimate     PWSViolationSubCategory = "Weight Estimate"
// 	PWSViolationSubCategoryAdditionalServices PWSViolationSubCategory = "Additional Services"
// 	PWSViolationSubCategoryInventory          PWSViolationSubCategory = "Inventory & Documentation"
// 	PWSViolationSubCategoryPackingUnpacking   PWSViolationSubCategory = "Packing/Unpacking"
// 	PWSViolationSubCategoryShipmentSchedule   PWSViolationSubCategory = "Shipment Schedule"
// 	PWSViolationSubCategoryShipmentWeights    PWSViolationSubCategory = "Shipment Weights"
// 	PWSViolationSubCategoryStorage            PWSViolationSubCategory = "Storage"
// 	PWSViolationSubCategoryWorkforce          PWSViolationSubCategory = "Workforce/Sub-Contractor Management"
// 	PWSViolationSubCategoryLossAndDamage      PWSViolationSubCategory = "Loss & Damage"
// 	PWSViolationSubCategoryInconvenience      PWSViolationSubCategory = "Inconvenience & Hardship Claims"
// )

type PWSViolation struct {
	ID                   int                  `json:"id" db:"id"`
	ParagraphNumber      string               `db:"paragraph_number"`
	Title                string               `db:"title"`
	Category             PWSViolationCategory `db:"category"`
	SubCategory          string               `db:"sub_category"`
	RequirementSummary   string               `db:"requirement_summary"`
	RequirementStatement string               `db:"requirement_statement"`
	IsKpi                bool                 `db:"is_kpi"`
	AdditionalDataElem   string               `db:"additional_data_elem"`
}

type PWSViolations []PWSViolation

func (p PWSViolations) TableName() string {
	return "pws_violations"
}
