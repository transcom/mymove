package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	FieldsLayout: map[string]FieldPos{
		"ServiceMemberName":            FormField(45, 42, 105, floatPtr(10), nil),
		"MaxSITStorageEntitlement":     FormField(153, 115, 49, floatPtr(10), nil),
		"WeightAllotmentSelf":          FormField(74, 104, 16, floatPtr(10), nil),
		"WeightAllotmentProgear":       FormField(74, 110, 16, floatPtr(10), nil),
		"WeightAllotmentProgearSpouse": FormField(74, 115, 16, floatPtr(10), nil),
		"TotalWeightAllotment":         FormField(74, 120, 16, floatPtr(10), nil),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{},
}
