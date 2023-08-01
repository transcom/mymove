package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type LineOfAccounting struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	LoaSysID               *int       `json:"loa_sys_id" db:"loa_sys_id"`
	LoaDptID               *string    `json:"loa_dpt_id" db:"loa_dpt_id"`
	LoaTnsfrDptNm          *string    `json:"loa_tnsfr_dpt_nm" db:"loa_tnsfr_dpt_nm"`
	LoaBafID               *string    `json:"loa_baf_id" db:"loa_baf_id"`
	LoaTrsySfxTx           *string    `json:"loa_trsy_sfx_tx" db:"loa_trsy_sfx_tx"`
	LoaMajClmNm            *string    `json:"loa_maj_clm_nm" db:"loa_maj_clm_nm"`
	LoaOpAgncyID           *string    `json:"loa_op_agncy_id" db:"loa_op_agncy_id"`
	LoaAlltSnID            *string    `json:"loa_allt_sn_id" db:"loa_allt_sn_id"`
	LoaPgmElmntID          *string    `json:"loa_pgm_elmnt_id" db:"loa_pgm_elmnt_id"`
	LoaTskBdgtSblnTx       *string    `json:"loa_tsk_bdgt_sbln_tx" db:"loa_tsk_bdgt_sbln_tx"`
	LoaDfAgncyAlctnRcpntID *string    `json:"loa_df_agncy_alctn_rcpnt_id" db:"loa_df_agncy_alctn_rcpnt_id"`
	LoaJbOrdNm             *string    `json:"loa_jb_ord_nm" db:"loa_jb_ord_nm"`
	LoaSbaltmtRcpntID      *string    `json:"loa_sbaltmt_rcpnt_id" db:"loa_sbaltmt_rcpnt_id"`
	LoaWkCntrRcpntNm       *string    `json:"loa_wk_cntr_rcpnt_nm" db:"loa_wk_cntr_rcpnt_nm"`
	LoaMajRmbsmtSrcID      *string    `json:"loa_maj_rmbsmt_src_id" db:"loa_maj_rmbsmt_src_id"`
	LoaDtlRmbsmtSrcID      *string    `json:"loa_dtl_rmbsmt_src_id" db:"loa_dtl_rmbsmt_src_id"`
	LoaCustNm              *string    `json:"loa_cust_nm" db:"loa_cust_nm"`
	LoaObjClsID            *string    `json:"loa_obj_cls_id" db:"loa_obj_cls_id"`
	LoaSrvSrcID            *string    `json:"loa_srv_src_id" db:"loa_srv_src_id"`
	LoaSpclIntrID          *string    `json:"loa_spcl_intr_id" db:"loa_spcl_intr_id"`
	LoaBdgtAcntClsNm       *string    `json:"loa_bdgt_acnt_cls_nm" db:"loa_bdgt_acnt_cls_nm"`
	LoaDocID               *string    `json:"loa_doc_id" db:"loa_doc_id"`
	LoaClsRefID            *string    `json:"loa_cls_ref_id" db:"loa_cls_ref_id"`
	LoaInstlAcntgActID     *string    `json:"loa_instl_acntg_act_id" db:"loa_instl_acntg_act_id"`
	LoaLclInstlID          *string    `json:"loa_lcl_instl_id" db:"loa_lcl_instl_id"`
	LoaFmsTrnsactnID       *string    `json:"loa_fms_trnsactn_id" db:"loa_fms_trnsactn_id"`
	LoaDscTx               *string    `json:"loa_dsc_tx" db:"loa_dsc_tx"`
	LoaBgnDt               *time.Time `json:"loa_bgn_dt" db:"loa_bgn_dt"`
	LoaEndDt               *time.Time `json:"loa_end_dt" db:"loa_end_dt"`
	LoaFnctPrsNm           *string    `json:"loa_fnct_prs_nm" db:"loa_fnct_prs_nm"`
	LoaStatCd              *string    `json:"loa_stat_cd" db:"loa_stat_cd"`
	LoaHistStatCd          *string    `json:"loa_hist_stat_cd" db:"loa_hist_stat_cd"`
	LoaHsGdsCd             *string    `json:"loa_hs_gds_cd" db:"loa_hs_gds_cd"`
	OrgGrpDfasCd           *string    `json:"org_grp_dfas_cd" db:"org_grp_dfas_cd"`
	LoaUic                 *string    `json:"loa_uic" db:"loa_uic"`
	LoaTrnsnID             *string    `json:"loa_trnsn_id" db:"loa_trnsn_id"`
	LoaSubAcntID           *string    `json:"loa_sub_acnt_id" db:"loa_sub_acnt_id"`
	LoaBetCd               *string    `json:"loa_bet_cd" db:"loa_bet_cd"`
	LoaFndTyFgCd           *string    `json:"loa_fnd_ty_fg_cd" db:"loa_fnd_ty_fg_cd"`
	LoaBgtLnItmID          *string    `json:"loa_bgt_ln_itm_id" db:"loa_bgt_ln_itm_id"`
	LoaScrtyCoopImplAgncCd *string    `json:"loa_scrty_coop_impl_agnc_cd" db:"loa_scrty_coop_impl_agnc_cd"`
	LoaScrtyCoopDsgntrCd   *string    `json:"loa_scrty_coop_dsgntr_cd" db:"loa_scrty_coop_dsgntr_cd"`
	LoaScrtyCoopLnItmID    *string    `json:"loa_scrty_coop_ln_itm_id" db:"loa_scrty_coop_ln_itm_id"`
	LoaAgncDsbrCd          *string    `json:"loa_agnc_dsbr_cd" db:"loa_agnc_dsbr_cd"`
	LoaAgncAcntngCd        *string    `json:"loa_agnc_acntng_cd" db:"loa_agnc_acntng_cd"`
	LoaFndCntrID           *string    `json:"loa_fnd_cntr_id" db:"loa_fnd_cntr_id"`
	LoaCstCntrID           *string    `json:"loa_cst_cntr_id" db:"loa_cst_cntr_id"`
	LoaPrjID               *string    `json:"loa_prj_id" db:"loa_prj_id"`
	LoaActvtyID            *string    `json:"loa_actvty_id" db:"loa_actvty_id"`
	LoaCstCd               *string    `json:"loa_cst_cd" db:"loa_cst_cd"`
	LoaWrkOrdID            *string    `json:"loa_wrk_ord_id" db:"loa_wrk_ord_id"`
	LoaFnclArID            *string    `json:"loa_fncl_ar_id" db:"loa_fncl_ar_id"`
	LoaScrtyCoopCustCd     *string    `json:"loa_scrty_coop_cust_cd" db:"loa_scrty_coop_cust_cd"`
	LoaEndFyTx             *int       `json:"loa_end_fy_tx" db:"loa_end_fy_tx"`
	LoaBgFyTx              *int       `json:"loa_bg_fy_tx" db:"loa_bg_fy_tx"`
	LoaBgtRstrCd           *string    `json:"loa_bgt_rstr_cd" db:"loa_bgt_rstr_cd"`
	LoaBgtSubActCd         *string    `json:"loa_bgt_sub_act_cd" db:"loa_bgt_sub_act_cd"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (l LineOfAccounting) TableName() string {
	return "lines_of_accounting"
}

// This func will take the "Desired" values, compare it to the "Internal" values that are currently captured, and map to what is stored internally.
// So, for example referencing TACs, as of right now TAC codes are only utilizing the "TAC" value. But, the .txt file holds way more information. So, we created the desired values
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
