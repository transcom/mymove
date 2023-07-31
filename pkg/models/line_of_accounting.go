package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type LineOfAccounting struct {
	ID                                                   uuid.UUID `json:"id"`
	LOA                                                  string    `json:"loa"`
	CreatedAt                                            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                                            time.Time `json:"updated_at" db:"updated_at"`
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

// TableName overrides the table name used by Pop.
func (t LineOfAccounting) TableName() string {
	return "line_of_ccounting"
}

// This func will take the "Desired" values, compare it to the "Internal" values that are currently captured, and map to what is stored internally.
// So, for example, as of right now TAC codes are only utilizing the "TAC" value. But, the .txt file holds way more information. So, we created the desired values
// struct so that we can map only the desired values. This is the same for LOAs at the moment; however, LOAs are currently not being utilized
// via this struct at all yet. This is the groundwork for future feature updates.
// This func will take the "Desired" values, compare it to the "Internal" values that are currently captured,
// and map to what is stored internally.
func MapLineOfAccountingFileRecordToInternalStruct(loaFileRecord LineOfAccountingTrdmFileRecord) LineOfAccounting {
	return LineOfAccounting{
		LOA:                                loaFileRecord.LOA_SYS_ID,
		DepartmentID:                       loaFileRecord.LOA_DPT_ID,
		TransferDepartmentName:             loaFileRecord.LOA_TNSFR_DPT_NM,
		BasicAppropriationFundID:           loaFileRecord.LOA_BAF_ID,
		TreasurySuffixText:                 loaFileRecord.LOA_TRSY_SFX_TX,
		MajorClaimantName:                  loaFileRecord.LOA_MAJ_CLM_NM,
		OperatingAgencyID:                  loaFileRecord.LOA_OP_AGNCY_ID,
		AllotmentSerialNumberID:            loaFileRecord.LOA_ALLT_SN_ID,
		ProgramElementID:                   loaFileRecord.LOA_PGM_ELMNT_ID,
		TaskBudgetSublineText:              loaFileRecord.LOA_TSK_BDGT_SBLN_TX,
		DefenseAgencyAllocationRecipientID: loaFileRecord.LOA_DF_AGNCY_ALCTN_RCPNT_ID,
		JobOrderName:                       loaFileRecord.LOA_JB_ORD_NM,
		SubAllotmentRecipientId:            loaFileRecord.LOA_SBALTMT_RCPNT_ID,
		WorkCenterRecipientName:            loaFileRecord.LOA_WK_CNTR_RCPNT_NM,
		MajorReimbursementSourceID:         loaFileRecord.LOA_MAJ_RMBSMT_SRC_ID,
		DetailReimbursementSourceID:        loaFileRecord.LOA_DTL_RMBSMT_SRC_ID,
		CustomerName:                       loaFileRecord.LOA_CUST_NM,
		ObjectClassID:                      loaFileRecord.LOA_OBJ_CLS_ID,
		ServiceSourceID:                    loaFileRecord.LOA_SRV_SRC_ID,
		SpecialInterestID:                  loaFileRecord.LOA_SPCL_INTR_ID,
		BudgetAccountClassificationName:    loaFileRecord.LOA_BDGT_ACNT_CLS_NM,
		DocumentID:                         loaFileRecord.LOA_DOC_ID,
		ClassReferenceID:                   loaFileRecord.LOA_CLS_REF_ID,
		InstallationAccountingActivityID:   loaFileRecord.LOA_INSTL_ACNTG_ACT_ID,
		LocalInstallationID:                loaFileRecord.LOA_LCL_INSTL_ID,
		FMSTransactionID:                   loaFileRecord.LOA_FMS_TRNSACTN_ID,
		DescriptionText:                    loaFileRecord.LOA_DSC_TX,
		BeginningDate:                      loaFileRecord.LOA_BGN_DT,
		EndDate:                            loaFileRecord.LOA_END_DT,
		FunctionalPersonName:               loaFileRecord.LOA_FNCT_PRS_NM,
		StatusCode:                         loaFileRecord.LOA_STAT_CD,
		HistoryStatusCode:                  loaFileRecord.LOA_HIST_STAT_CD,
		HouseholdGoodsCode:                 loaFileRecord.LOA_HS_GDS_CD,
		OrganizationGroupDefenseFinanceAccountingServiceCode: loaFileRecord.ORG_GRP_DFAS_CD,
		UnitIdentificationCode:                               loaFileRecord.LOA_UIC,
		TransactionID:                                        loaFileRecord.LOA_TRNSN_ID,
		SubordinateAccountID:                                 loaFileRecord.LOA_SUB_ACNT_ID,
		BusinessEventTypeCode:                                loaFileRecord.LOA_BET_CD,
		FundTypeFlagCode:                                     loaFileRecord.LOA_FND_TY_FG_CD,
		BudgetLineItemID:                                     loaFileRecord.LOA_BGT_LN_ITM_ID,
		SecurityCooperationImplementingAgencyCode:            loaFileRecord.LOA_SCRTY_COOP_IMPL_AGNC_CD,
		SecurityCooperationDesignatorID:                      loaFileRecord.LOA_SCRTY_COOP_DSGNTR_CD,
		SecurityCooperationLineItemID:                        loaFileRecord.LOA_SCRTY_COOP_LN_ITM_ID,
		AgencyDisbursingCode:                                 loaFileRecord.LOA_AGNC_DSBR_CD,
		AgencyAccountingCode:                                 loaFileRecord.LOA_AGNC_ACNTNG_CD,
		FundCenterID:                                         loaFileRecord.LOA_FND_CNTR_ID,
		CostCenterID:                                         loaFileRecord.LOA_CST_CNTR_ID,
		ProjectTaskID:                                        loaFileRecord.LOA_PRJ_ID,
		ActivityID:                                           loaFileRecord.LOA_ACTVTY_ID,
		CostCode:                                             loaFileRecord.LOA_CST_CD,
		WorkOrderID:                                          loaFileRecord.LOA_WRK_ORD_ID,
		FunctionalAreaID:                                     loaFileRecord.LOA_FNCL_AR_ID,
		SecurityCooperationCustomerCode:                      loaFileRecord.LOA_SCRTY_COOP_CUST_CD,
		EndingFiscalYear:                                     loaFileRecord.LOA_END_FY_TX,
		BeginningFiscalYear:                                  loaFileRecord.LOA_BG_FY_TX,
		BudgetRestrictionCode:                                loaFileRecord.LOA_BGT_RSTR_CD,
		BudgetSubActivityCode:                                loaFileRecord.LOA_BGT_SUB_ACT_CD,
	}
}
