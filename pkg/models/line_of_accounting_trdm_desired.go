package models

import (
	"time"
)

// !IMPORTANT! This struct is not the file record nor the internal Line Of Accounting struct, see LineOfAccountingTrdmFileRecord model for this.
// This struct is what will be returned when the TRDM .txt file gets parsed.
// See LineOfAccountingTrdmFileRecord for the struct representing the .txt file.

type LineOfAccountingDesiredFromTRDM struct {
	LOA                                                  string    `json:"loa"`
	DepartmentID                                         string    `json:"department_id"`
	TransferDepartmentName                               string    `json:"transfer_department_name"`
	BasicAppropriationFundID                             string    `json:"basic_appropriation_fund_id"`
	TreasurySuffixText                                   string    `json:"treasury_suffix_text"`
	MajorClaimantName                                    string    `json:"major_claimant_name"`
	OperatingAgencyID                                    string    `json:"operating_agency_id"`
	AllotmentSerialNumberID                              string    `json:"allotment_serial_number_id"`
	ProgramElementID                                     string    `json:"program_element_id"`
	TaskBudgetSublineText                                string    `json:"task_budget_subline_text"`
	DefenseAgencyAllocationRecipientID                   string    `json:"defense_agency_allocation_recipient_id"`
	JobOrderName                                         string    `json:"job_order_name"`
	SubAllotmentRecipientId                              string    `json:"sub_allotment_recipient_id"`
	WorkCenterRecipientName                              string    `json:"work_center_recipient_name"`
	MajorReimbursementSourceID                           string    `json:"major_reimbursement_source_id"`
	DetailReimbursementSourceID                          string    `json:"detail_reimbursement_source_id"`
	CustomerName                                         string    `json:"customer_name"`
	ObjectClassID                                        string    `json:"object_class_id"`
	ServiceSourceID                                      string    `json:"service_source_id"`
	SpecialInterestID                                    string    `json:"special_interest_id"`
	BudgetAccountClassificationName                      string    `json:"budget_account_classification_name"`
	DocumentID                                           string    `json:"document_id"`
	ClassReferenceID                                     string    `json:"class_reference_id"`
	InstallationAccountingActivityID                     string    `json:"installation_accounting_activity_id"`
	LocalInstallationID                                  string    `json:"local_installation_id"`
	FMSTransactionID                                     string    `json:"fms_transaction_id"`
	DescriptionText                                      string    `json:"description_text"`
	BeginningDate                                        time.Time `json:"beginning_date"`
	EndDate                                              time.Time `json:"end_date"`
	FunctionalPersonName                                 string    `json:"functional_person_name"`
	StatusCode                                           string    `json:"status_code"`
	HistoryStatusCode                                    string    `json:"history_status_code"`
	HouseholdGoodsCode                                   string    `json:"household_goods_code"`
	OrganizationGroupDefenseFinanceAccountingServiceCode string    `json:"organization_group_defense_finance_accounting_service_code"`
	UnitIdentificationCode                               string    `json:"unit_identification_code"`
	TransactionID                                        string    `json:"transaction_id"`
	SubordinateAccountID                                 string    `json:"subordinate_account_id"`
	BusinessEventTypeCode                                string    `json:"business_event_type_code"`
	FundTypeFlagCode                                     string    `json:"fund_type_flag_code"`
	BudgetLineItemID                                     string    `json:"budget_line_item_id"`
	SecurityCooperationImplementingAgencyCode            string    `json:"security_cooperation_implementing_agency_code"`
	SecurityCooperationDesignatorID                      string    `json:"security_cooperation_designator_id"`
	SecurityCooperationLineItemID                        string    `json:"security_cooperation_line_item_id"`
	AgencyDisbursingCode                                 string    `json:"agency_disbursing_code"`
	AgencyAccountingCode                                 string    `json:"agency_accounting_code"`
	FundCenterID                                         string    `json:"fund_center_id"`
	CostCenterID                                         string    `json:"cost_center_id"`
	ProjectTaskID                                        string    `json:"project_task_id"`
	ActivityID                                           string    `json:"activity_id"`
	CostCode                                             string    `json:"cost_code"`
	WorkOrderID                                          string    `json:"work_order_id"`
	FunctionalAreaID                                     string    `json:"functional_area_id"`
	SecurityCooperationCustomerCode                      string    `json:"security_cooperation_customer_code"`
	EndingFiscalYear                                     int       `json:"ending_fiscal_year"`
	BeginningFiscalYear                                  int       `json:"beginning_fiscal_year"`
	BudgetRestrictionCode                                string    `json:"budget_restriction_code"`
	BudgetSubActivityCode                                string    `json:"budget_sub_activity_code"`
}
