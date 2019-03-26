package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"PreparationDate":                 FormField(155.5, 22, 46, floatPtr(10), nil, nil),
		"ServiceMemberName":               FormField(10, 42, 105, floatPtr(10), nil, nil),
		"DODId":                           FormField(10, 54, 40, floatPtr(10), nil, nil),
		"ServiceBranch":                   FormField(54, 54, 44, floatPtr(10), nil, nil),
		"RankGrade":                       FormField(102.5, 54, 47, floatPtr(10), nil, nil),
		"PreferredEmail":                  FormField(153.5, 54, 60, floatPtr(10), nil, nil),
		"PreferredPhoneNumber":            FormField(153.5, 42, 60, floatPtr(10), nil, nil),
		"WeightAllotment":                 FormField(74, 92, 16, floatPtr(10), nil, stringPtr("RM")),
		"WeightAllotmentProgear":          FormField(74, 98, 16, floatPtr(10), nil, stringPtr("RM")),
		"WeightAllotmentProgearSpouse":    FormField(74, 103, 16, floatPtr(10), nil, stringPtr("RM")),
		"TotalWeightAllotment":            FormField(74, 108, 16, floatPtr(10), nil, stringPtr("RM")),
		"POVAuthorized":                   FormField(102.25, 104, 45, floatPtr(10), nil, nil),
		"AuthorizedOrigin":                FormField(102.25, 91, 45, floatPtr(10), nil, nil),
		"MaxSITStorageEntitlement":        FormField(153.5, 104, 49, floatPtr(10), nil, nil),
		"AuthorizedDestination":           FormField(153.5, 91, 60, floatPtr(10), nil, nil),
		"OrdersIssueDate":                 FormField(9.5, 73, 40, floatPtr(10), nil, nil),
		"OrdersTypeAndOrdersNumber":       FormField(54, 73, 44, floatPtr(10), nil, nil),
		"IssuingBranchOrAgency":           FormField(102.5, 73, 47, floatPtr(10), nil, nil),
		"NewDutyAssignment":               FormField(153, 73, 60, floatPtr(10), nil, nil),
		"TAC":                             FormField(10, 233, 45, floatPtr(10), nil, nil),
		"SAC":                             FormField(10, 222, 45, floatPtr(10), nil, nil),
		"ShipmentNumberAndTypes":          FormField(9.5, 122.5, 41, floatPtr(10), nil, nil),
		"ShipmentPickUpDates":             FormField(54, 122.5, 46, floatPtr(10), nil, nil),
		"ShipmentWeights":                 FormField(103, 122.5, 41, floatPtr(10), nil, nil),
		"ShipmentCurrentShipmentStatuses": FormField(153.5, 122.5, 41, floatPtr(10), nil, nil),
		"MaxObligationGCC100":             FormField(40, 182.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"TotalWeightAllotmentRepeat":      FormField(74, 182.5, 16, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationGCC95":              FormField(40, 188.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationSIT":                FormField(40, 194.75, 22, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationGCCMaxAdvance":      FormField(40, 201, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationGCC100":          FormField(133, 182.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"PPMRemainingEntitlement":         FormField(167, 182.5, 16, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationGCC95":           FormField(133, 188.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationSIT":             FormField(133, 194.75, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationAdvance":         FormField(133, 201, 22, floatPtr(10), nil, stringPtr("RM")),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{

	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{
		"PreparationDate":             FormField(155.5, 22, 46, floatPtr(10), nil, nil),
		"ContractedExpenseMemberPaid": FormField(156.5, 49, 20, floatPtr(10), nil, stringPtr("RM")),
		"RentalEquipmentMemberPaid":   FormField(156.5, 55.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"PackingMaterialsMemberPaid":  FormField(156.5, 61.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"WeighingFeesMemberPaid":      FormField(156.5, 68, 20, floatPtr(10), nil, stringPtr("RM")),
		"GasMemberPaid":               FormField(156.5, 74, 20, floatPtr(10), nil, stringPtr("RM")),
		"TollsMemberPaid":             FormField(156.5, 80, 20, floatPtr(10), nil, stringPtr("RM")),
		"OilMemberPaid":               FormField(156.5, 86.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"OtherMemberPaid":             FormField(156.5, 93, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalMemberPaid":             FormField(156.5, 99.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"ContractedExpenseGTCCPaid":   FormField(181.5, 49, 20, floatPtr(10), nil, stringPtr("RM")),
		"RentalEquipmentGTCCPaid":     FormField(181.5, 55.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"PackingMaterialsGTCCPaid":    FormField(181.5, 61.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"WeighingFeesGTCCPaid":        FormField(181.5, 68, 20, floatPtr(10), nil, stringPtr("RM")),
		"GasGTCCPaid":                 FormField(181.5, 74, 20, floatPtr(10), nil, stringPtr("RM")),
		"TollsGTCCPaid":               FormField(181.5, 80, 20, floatPtr(10), nil, stringPtr("RM")),
		"OilGTCCPaid":                 FormField(181.5, 86.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"OtherGTCCPaid":               FormField(181.5, 93, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalGTCCPaid":               FormField(181.5, 99.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalMemberPaidRepeated":     FormField(74, 42, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalGTCCPaidRepeated":       FormField(74, 53, 30, floatPtr(10), nil, stringPtr("RM")),
	},
}
