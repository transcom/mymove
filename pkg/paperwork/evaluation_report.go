package paperwork

import (
	"github.com/transcom/mymove/pkg/models"
)

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
)

// EvaluationReortPage1Layout specifies the layout and template of a
var EvaluationReportPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/qae_spike_template.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"CUIBanner":        FormField(0, 1.5, 216, floatPtr(10), nil, stringPtr("CM")),
		"ReportID":         FormField(80, 40, 216, floatPtr(10), nil, stringPtr("LM")),
		"DateOfInspection": FormField(80, 60, 216, floatPtr(10), nil, stringPtr("LM")),
		"EvaluationType":   FormField(80, 80, 216, floatPtr(10), nil, stringPtr("LM")),
		"Remarks":          FormField(80, 105, 216, floatPtr(10), nil, stringPtr("LM")),
	},
}

// EvaluationReportPage1Values is an object representing a Shipment Summary Worksheet
type EvaluationReportPage1Values struct {
	CUIBanner        string
	ReportID         string
	DateOfInspection string
	EvaluationType   string
	Remarks          string
}

func FormatValuesEvaluationReportPage1(data models.EvaluationReport) EvaluationReportPage1Values {
	page1 := EvaluationReportPage1Values{
		CUIBanner:        controlledUnclassifiedInformationText,
		ReportID:         data.ID.String(),
		DateOfInspection: data.InspectionDate.Format("2006-01-02"),
		EvaluationType:   string(data.Type),
	}
	if data.Remarks != nil {
		page1.Remarks = *data.Remarks
	}
	return page1
}
