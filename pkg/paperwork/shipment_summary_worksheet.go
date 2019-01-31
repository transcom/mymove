package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"PreparationDate":                FormField(156, 22, 46, floatPtr(10), nil),
		"ServiceMemberName":              FormField(45, 42, 105, floatPtr(10), nil),
		"MaxSITStorageEntitlement":       FormField(153, 115, 49, floatPtr(10), nil),
		"WeightAllotmentSelf":            FormField(74, 104, 16, floatPtr(10), nil),
		"WeightAllotmentProgear":         FormField(74, 110, 16, floatPtr(10), nil),
		"WeightAllotmentProgearSpouse":   FormField(74, 115, 16, floatPtr(10), nil),
		"TotalWeightAllotment":           FormField(74, 120, 16, floatPtr(10), nil),
		"POVAuthorized":                  FormField(103, 115, 45, floatPtr(10), nil),
		"AuthorizedOrigin":               FormField(103, 104, 45, floatPtr(10), nil),
		"AuthorizedDestination":          FormField(153, 104, 45, floatPtr(10), nil),
		"OrdersIssueDate":                FormField(10, 85, 40, floatPtr(10), nil),
		"OrdersTypeAndOrdersNumber":      FormField(54, 85, 44, floatPtr(10), nil),
		"IssuingBranchOrAgency":          FormField(103, 85, 47, floatPtr(10), nil),
		"NewDutyAssignment":              FormField(154, 85, 60, floatPtr(10), nil),
		"TAC":                            FormField(53, 233, 45, floatPtr(10), nil),
		"Shipment1NumberAndType":         FormField(10, 135, 41, floatPtr(10), nil),
		"Shipment1PickUpDate":            FormField(55, 135, 46, floatPtr(10), nil),
		"Shipment1Weight":                FormField(104, 135, 41, floatPtr(10), nil),
		"Shipment1CurrentShipmentStatus": FormField(154, 135, 41, floatPtr(10), nil),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{},
}
