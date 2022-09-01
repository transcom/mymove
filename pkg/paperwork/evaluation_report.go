package paperwork

import (
	"github.com/transcom/mymove/pkg/models"
)

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
)

// EvaluationReportPage1Layout specifies the layout and template of the first page of an evaluation report
var EvaluationReportPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/qae_spike_template.png",

	FieldsLayout: map[string]FieldPos{
		"CUIBanner":        FormField(0, 1.5, 216, floatPtr(10), nil, stringPtr("CM")),
		"ReportID":         FormField(80, 40, 216, floatPtr(10), nil, stringPtr("LM")),
		"DateOfInspection": FormField(80, 60, 216, floatPtr(10), nil, stringPtr("LM")),
		"EvaluationType":   FormField(80, 80, 216, floatPtr(10), nil, stringPtr("LM")),
		"Remarks":          FormField(80, 105, 216, floatPtr(10), nil, stringPtr("LM")),
	},
}

// EvaluationReportPage1Values is an object representing an evaluation report
type EvaluationReportPage1Values struct {
	CUIBanner        string
	ReportID         string
	DateOfInspection string
	EvaluationType   string
	Remarks          string
}

func FormatValuesEvaluationReportPage1(data models.EvaluationReport) EvaluationReportPage1Values {
	page1 := EvaluationReportPage1Values{
		CUIBanner:      controlledUnclassifiedInformationText,
		ReportID:       data.ID.String(),
		EvaluationType: string(data.Type),
	}
	if data.Remarks != nil {
		page1.Remarks = *data.Remarks
	}
	if data.InspectionDate != nil {
		page1.DateOfInspection = data.InspectionDate.Format("2006-01-02")
	}
	return page1
}
