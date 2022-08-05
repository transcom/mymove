package paperwork

// ShipmentSummaryPage1Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage1Layout = FormLayout{
	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page1.png",

	// For now only lists a single shipment. Will need to update to accommodate multiple shipments
	FieldsLayout: map[string]FieldPos{
		"CUIBanner":                       FormField(0, 1.5, 216, floatPtr(10), nil, stringPtr("CM")),
		"PreparationDate":                 FormField(155.5, 23, 46, floatPtr(10), nil, nil),
		"ServiceMemberName":               FormField(10, 43, 90, floatPtr(10), nil, nil),
		"DODId":                           FormField(153.5, 43, 60, floatPtr(10), nil, nil),
		"ServiceBranch":                   FormField(10, 54, 40, floatPtr(10), nil, nil),
		"RankGrade":                       FormField(54, 54, 44, floatPtr(10), nil, nil),
		"PreferredEmail":                  FormField(102.5, 54, 47, floatPtr(10), nil, nil),
		"PreferredPhoneNumber":            FormField(153.5, 54, 60, floatPtr(10), nil, nil),
		"WeightAllotment":                 FormField(73.5, 96, 16, floatPtr(10), nil, stringPtr("RM")),
		"WeightAllotmentProgear":          FormField(73.5, 103, 16, floatPtr(10), nil, stringPtr("RM")),
		"WeightAllotmentProgearSpouse":    FormField(73.5, 110, 16, floatPtr(10), nil, stringPtr("RM")),
		"TotalWeightAllotment":            FormField(73.5, 116.5, 16, floatPtr(10), nil, stringPtr("RM")),
		"POVAuthorized":                   FormField(102.25, 91, 45, floatPtr(10), nil, nil),
		"AuthorizedOrigin":                FormField(153.5, 91, 49, floatPtr(10), nil, nil),
		"MaxSITStorageEntitlement":        FormField(102.25, 104, 45, floatPtr(10), nil, nil),
		"AuthorizedDestination":           FormField(153.5, 104, 49, floatPtr(10), nil, nil),
		"MileageTotal":                    FormField(153.5, 116.5, 49, floatPtr(10), nil, nil),
		"OrdersIssueDate":                 FormField(9.5, 73, 40, floatPtr(10), nil, nil),
		"OrdersTypeAndOrdersNumber":       FormField(54, 73, 44, floatPtr(10), nil, nil),
		"IssuingBranchOrAgency":           FormField(102.5, 73, 47, floatPtr(10), nil, nil),
		"NewDutyAssignment":               FormField(153, 73, 60, floatPtr(10), nil, nil),
		"ShipmentNumberAndTypes":          FormField(9.5, 141, 41, floatPtr(10), nil, nil),
		"ShipmentPickUpDates":             FormField(54, 141, 46, floatPtr(10), nil, nil),
		"ShipmentWeights":                 FormField(103, 141, 41, floatPtr(10), nil, nil),
		"ShipmentCurrentShipmentStatuses": FormField(153.5, 141, 41, floatPtr(10), nil, nil),
		"SITNumberAndTypes":               FormField(9.5, 180, 41, floatPtr(10), nil, nil),
		"SITEntryDates":                   FormField(54, 180, 46, floatPtr(10), nil, nil),
		"SITEndDates":                     FormField(103, 180, 41, floatPtr(10), nil, nil),
		"SITDaysInStorage":                FormField(153.5, 180, 41, floatPtr(10), nil, nil),
		"MaxObligationGCC100":             FormField(39, 225.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"TotalWeightAllotmentRepeat":      FormField(73, 225.5, 16, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationGCC95":              FormField(39, 233, 22, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationSIT":                FormField(39, 240, 22, floatPtr(10), nil, stringPtr("RM")),
		"MaxObligationGCCMaxAdvance":      FormField(39, 247, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationGCC100":          FormField(133, 225.5, 22, floatPtr(10), nil, stringPtr("RM")),
		"PPMRemainingEntitlement":         FormField(167, 225.5, 16, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationGCC95":           FormField(133, 233, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationSIT":             FormField(133, 240, 22, floatPtr(10), nil, stringPtr("RM")),
		"ActualObligationAdvance":         FormField(133, 247, 22, floatPtr(10), nil, stringPtr("RM")),
	},
}

// ShipmentSummaryPage2Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage2Layout = FormLayout{

	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page2.png",

	FieldsLayout: map[string]FieldPos{
		"CUIBanner":                   FormField(0, 2, 216, floatPtr(10), nil, stringPtr("CM")),
		"PreparationDate":             FormField(155.5, 23, 46, floatPtr(10), nil, nil),
		"SAC":                         FormField(10, 43, 45, floatPtr(10), nil, nil),
		"TAC":                         FormField(10, 54, 45, floatPtr(10), nil, nil),
		"ContractedExpenseMemberPaid": FormField(55, 98, 20, floatPtr(10), nil, stringPtr("RM")),
		"RentalEquipmentMemberPaid":   FormField(55, 104, 20, floatPtr(10), nil, stringPtr("RM")),
		"PackingMaterialsMemberPaid":  FormField(55, 110.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"WeighingFeesMemberPaid":      FormField(55, 117, 20, floatPtr(10), nil, stringPtr("RM")),
		"GasMemberPaid":               FormField(55, 123, 20, floatPtr(10), nil, stringPtr("RM")),
		"TollsMemberPaid":             FormField(55, 129, 20, floatPtr(10), nil, stringPtr("RM")),
		"OilMemberPaid":               FormField(55, 135.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"OtherMemberPaid":             FormField(55, 142, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalMemberPaid":             FormField(55, 148, 20, floatPtr(10), nil, stringPtr("RM")),
		"ContractedExpenseGTCCPaid":   FormField(82.5, 98, 20, floatPtr(10), nil, stringPtr("RM")),
		"RentalEquipmentGTCCPaid":     FormField(82.5, 104, 20, floatPtr(10), nil, stringPtr("RM")),
		"PackingMaterialsGTCCPaid":    FormField(82.5, 110.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"WeighingFeesGTCCPaid":        FormField(82.5, 117, 20, floatPtr(10), nil, stringPtr("RM")),
		"GasGTCCPaid":                 FormField(82.5, 123, 20, floatPtr(10), nil, stringPtr("RM")),
		"TollsGTCCPaid":               FormField(82.5, 129, 20, floatPtr(10), nil, stringPtr("RM")),
		"OilGTCCPaid":                 FormField(82.5, 135.5, 20, floatPtr(10), nil, stringPtr("RM")),
		"OtherGTCCPaid":               FormField(82.5, 142, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalGTCCPaid":               FormField(82.5, 148, 20, floatPtr(10), nil, stringPtr("RM")),
		"TotalMemberPaidRepeated":     FormField(169.5, 98, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalGTCCPaidRepeated":       FormField(169.5, 104.5, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalMemberPaidSIT":          FormField(169.5, 135.5, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalGTCCPaidSIT":            FormField(169.5, 142, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalPaidNonSIT":             FormField(169.5, 111, 30, floatPtr(10), nil, stringPtr("RM")),
		"TotalPaidSIT":                FormField(169.5, 148.5, 30, floatPtr(10), nil, stringPtr("RM")),
	},
}

// ShipmentSummaryPage3Layout specifies the layout and template of a
// Shipment Summary Worksheet
var ShipmentSummaryPage3Layout = FormLayout{

	TemplateImagePath: "pkg/paperwork/formtemplates/shipment_summary_worksheet_page3.png",

	FieldsLayout: map[string]FieldPos{
		"CUIBanner":              FormField(0, 2, 216, floatPtr(10), nil, stringPtr("CM")),
		"PreparationDate":        FormField(155.5, 23, 46, floatPtr(10), nil, nil),
		"Descriptions":           FormField(10, 64, 135, floatPtr(10), nil, nil),
		"AmountsPaid":            FormField(155.5, 64, 46, floatPtr(10), nil, stringPtr("RM")),
		"ServiceMemberSignature": FormField(10, 125, 200, floatPtr(10), nil, nil),
		"SignatureDate":          FormField(10, 136, 200, floatPtr(10), nil, nil),
	},
}
