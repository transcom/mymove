package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	FieldsLayout: map[string]FieldPos{
		"ServiceMemberName": FormField(45, 38, 105, floatPtr(10), nil),
	},
}
