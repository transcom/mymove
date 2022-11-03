package models

import "github.com/gofrs/uuid"

type PWSViolationCategory string

const (
	PWSViolationCategoryPreMoveServices      PWSViolationCategory = "Pre-Move Services"
	PWSViolationCategoryPhysicalMoveServices PWSViolationCategory = "Physical Move Services"
	PWSViolationCategoryLiability            PWSViolationCategory = "Liability"
)

type PWSViolation struct {
	ID                   uuid.UUID            `json:"id" db:"id"`
	DisplayOrder         int                  `json:"display_order" db:"display_order"`
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

func (p PWSViolation) TableName() string {
	return "pws_violations"
}
