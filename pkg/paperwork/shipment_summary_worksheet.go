package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"PreparationDate":                 FormField(155.5, 22, 46, floatPtr(10), nil, nil),
		"ServiceMemberName":               FormField(10, 42, 105, floatPtr(10), nil, nil),
		"WeightAllotment":                 FormField(74, 92, 16, floatPtr(10), nil, nil),
		"WeightAllotmentProgear":          FormField(74, 98, 16, floatPtr(10), nil, nil),
		"WeightAllotmentProgearSpouse":    FormField(74, 103, 16, floatPtr(10), nil, nil),
		"TotalWeightAllotment":            FormField(74, 108, 16, floatPtr(10), nil, nil),
		"POVAuthorized":                   FormField(102.25, 104, 45, floatPtr(10), nil, nil),
		"AuthorizedOrigin":                FormField(102.25, 91, 45, floatPtr(10), nil, nil),
		"MaxSITStorageEntitlement":        FormField(153.5, 104, 49, floatPtr(10), nil, nil),
		"AuthorizedDestination":           FormField(153.5, 91, 45, floatPtr(10), nil, nil),
		"OrdersIssueDate":                 FormField(9.5, 73, 40, floatPtr(10), nil, nil),
		"OrdersTypeAndOrdersNumber":       FormField(54, 73, 44, floatPtr(10), nil, nil),
		"IssuingBranchOrAgency":           FormField(102.5, 73, 47, floatPtr(10), nil, nil),
		"NewDutyAssignment":               FormField(153, 73, 60, floatPtr(10), nil, nil),
		"TAC":                             FormField(10, 233, 45, floatPtr(10), nil, nil),
		"ShipmentNumberAndTypes":          FormField(9.5, 122.5, 41, floatPtr(10), nil, nil),
		"ShipmentPickUpDates":             FormField(54, 122.5, 46, floatPtr(10), nil, nil),
		"ShipmentWeights":                 FormField(103, 122.5, 41, floatPtr(10), nil, nil),
		"ShipmentCurrentShipmentStatuses": FormField(153.5, 122.5, 41, floatPtr(10), nil, nil),
		"GCC100":                          FormField(40, 182.5, 22, floatPtr(10), nil, nil),
		"TotalWeightAllotmentRepeat":      FormField(74, 182.5, 16, floatPtr(10), nil, nil),
		"GCC95":                           FormField(40, 188.5, 22, floatPtr(10), nil, nil),
		"GCCMaxAdvance":                   FormField(40, 201, 22, floatPtr(10), nil, nil),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{

	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{
		"PreparationDate":             FormField(155.5, 22, 46, floatPtr(10), nil, nil),
		"ContractedExpenseMemberPaid": FormField(156.5, 49, 20, floatPtr(10), nil, nil),
		"RentalEquipmentMemberPaid":   FormField(156.5, 55.5, 20, floatPtr(10), nil, nil),
		"PackingMaterialsMemberPaid":  FormField(156.5, 61.5, 20, floatPtr(10), nil, nil),
		"WeighingFeesMemberPaid":      FormField(156.5, 68, 20, floatPtr(10), nil, nil),
		"GasMemberPaid":               FormField(156.5, 74, 20, floatPtr(10), nil, nil),
		"TollsMemberPaid":             FormField(156.5, 80, 20, floatPtr(10), nil, nil),
		"OilMemberPaid":               FormField(156.5, 86.5, 20, floatPtr(10), nil, nil),
		"OtherMemberPaid":             FormField(156.5, 93, 20, floatPtr(10), nil, nil),
		"TotalMemberPaid":             FormField(156.5, 99.5, 20, floatPtr(10), nil, nil),
		"ContractedExpenseGTCCPaid":   FormField(181.5, 49, 20, floatPtr(10), nil, nil),
		"RentalEquipmentGTCCPaid":     FormField(181.5, 55.5, 20, floatPtr(10), nil, nil),
		"PackingMaterialsGTCCPaid":    FormField(181.5, 61.5, 20, floatPtr(10), nil, nil),
		"WeighingFeesGTCCPaid":        FormField(181.5, 68, 20, floatPtr(10), nil, nil),
		"GasGTCCPaid":                 FormField(181.5, 74, 20, floatPtr(10), nil, nil),
		"TollsGTCCPaid":               FormField(181.5, 80, 20, floatPtr(10), nil, nil),
		"OilGTCCPaid":                 FormField(181.5, 86.5, 20, floatPtr(10), nil, nil),
		"OtherGTCCPaid":               FormField(181.5, 93, 20, floatPtr(10), nil, nil),
		"TotalGTCCPaid":               FormField(181.5, 99.5, 20, floatPtr(10), nil, nil),
	},
}
