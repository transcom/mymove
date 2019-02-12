package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"PreparationDate":                 FormField(155.5, 22, 46, floatPtr(10), nil),
		"ServiceMemberName":               FormField(45, 42, 105, floatPtr(10), nil),
		"MaxSITStorageEntitlement":        FormField(153.75, 115, 49, floatPtr(10), nil),
		"WeightAllotment":                 FormField(74, 104, 16, floatPtr(10), nil),
		"WeightAllotmentProgear":          FormField(74, 110, 16, floatPtr(10), nil),
		"WeightAllotmentProgearSpouse":    FormField(74, 115, 16, floatPtr(10), nil),
		"TotalWeightAllotment":            FormField(74, 120, 16, floatPtr(10), nil),
		"POVAuthorized":                   FormField(103.25, 115, 45, floatPtr(10), nil),
		"AuthorizedOrigin":                FormField(103.25, 103, 45, floatPtr(10), nil),
		"AuthorizedDestination":           FormField(153.5, 103, 45, floatPtr(10), nil),
		"OrdersIssueDate":                 FormField(9.5, 85, 40, floatPtr(10), nil),
		"OrdersTypeAndOrdersNumber":       FormField(54, 85, 44, floatPtr(10), nil),
		"IssuingBranchOrAgency":           FormField(102.5, 85, 47, floatPtr(10), nil),
		"NewDutyAssignment":               FormField(153, 85, 60, floatPtr(10), nil),
		"TAC":                             FormField(53, 233, 45, floatPtr(10), nil),
		"ShipmentNumberAndTypes":          FormField(9.5, 135, 41, floatPtr(10), nil),
		"ShipmentPickUpDates":             FormField(54, 135, 46, floatPtr(10), nil),
		"ShipmentWeights":                 FormField(103, 135, 41, floatPtr(10), nil),
		"ShipmentCurrentShipmentStatuses": FormField(153.5, 135, 41, floatPtr(10), nil),
		"GCC100":                          FormField(40, 194.5, 22, floatPtr(10), nil),
		"TotalWeightAllotmentRepeat":      FormField(74, 194.5, 16, floatPtr(10), nil),
		"GCC95":                           FormField(40, 201, 22, floatPtr(10), nil),
		"GCCMaxAdvance":                   FormField(40, 213, 22, floatPtr(10), nil),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{},
}
