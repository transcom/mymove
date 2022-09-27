package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakePWSViolation(db *pop.Connection, assertions Assertions) models.PWSViolation {

	violation := models.PWSViolation{
		ID:                   uuid.Must(uuid.NewV4()),
		DisplayOrder:         1,
		ParagraphNumber:      "1.2.3",
		Title:                "Title",
		Category:             "Category",
		SubCategory:          "Customer Support",
		RequirementSummary:   "RequirementSummary",
		RequirementStatement: "RequirementStatement",
		IsKpi:                false,
		AdditionalDataElem:   "",
	}

	mergeModels(&violation, assertions.Violation)
	mustCreate(db, &violation, assertions.Stub)

	return violation
}
