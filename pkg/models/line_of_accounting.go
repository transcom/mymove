package models

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	LineOfAccountingHouseholdGoodsCodeCivilian string = "HC"
	LineOfAccountingHouseholdGoodsCodeEnlisted string = "HE"
	LineOfAccountingHouseholdGoodsCodeOfficer  string = "HO"
	LineOfAccountingHouseholdGoodsCodeOther    string = "HT"
	LineOfAccountingHouseholdGoodsCodeDual     string = "HD"
	LineOfAccountingHouseholdGoodsCodeNTS      string = "HS"
)

type LineOfAccounting struct {
	ID                     uuid.UUID  `json:"id" db:"id" pipe:"-"`
	LoaSysID               *string    `json:"loa_sys_id" db:"loa_sys_id" pipe:"LOA_SYS_ID"`
	LoaDptID               *string    `json:"loa_dpt_id" db:"loa_dpt_id" pipe:"LOA_DPT_ID"`
	LoaTnsfrDptNm          *string    `json:"loa_tnsfr_dpt_nm" db:"loa_tnsfr_dpt_nm" pipe:"LOA_TNSFR_DPT_NM"`
	LoaBafID               *string    `json:"loa_baf_id" db:"loa_baf_id" pipe:"LOA_BAF_ID"`
	LoaTrsySfxTx           *string    `json:"loa_trsy_sfx_tx" db:"loa_trsy_sfx_tx" pipe:"LOA_TRSY_SFX_TX"`
	LoaMajClmNm            *string    `json:"loa_maj_clm_nm" db:"loa_maj_clm_nm" pipe:"LOA_MAJ_CLM_NM"`
	LoaOpAgncyID           *string    `json:"loa_op_agncy_id" db:"loa_op_agncy_id" pipe:"LOA_OP_AGNCY_ID"`
	LoaAlltSnID            *string    `json:"loa_allt_sn_id" db:"loa_allt_sn_id" pipe:"LOA_ALLT_SN_ID"`
	LoaPgmElmntID          *string    `json:"loa_pgm_elmnt_id" db:"loa_pgm_elmnt_id" pipe:"LOA_PGM_ELMNT_ID"`
	LoaTskBdgtSblnTx       *string    `json:"loa_tsk_bdgt_sbln_tx" db:"loa_tsk_bdgt_sbln_tx" pipe:"LOA_TSK_BDGT_SBLN_TX"`
	LoaDfAgncyAlctnRcpntID *string    `json:"loa_df_agncy_alctn_rcpnt_id" db:"loa_df_agncy_alctn_rcpnt_id" pipe:"LOA_DF_AGNCY_ALCTN_RCPNT_ID"`
	LoaJbOrdNm             *string    `json:"loa_jb_ord_nm" db:"loa_jb_ord_nm" pipe:"LOA_JB_ORD_NM"`
	LoaSbaltmtRcpntID      *string    `json:"loa_sbaltmt_rcpnt_id" db:"loa_sbaltmt_rcpnt_id" pipe:"LOA_SBALTMT_RCPNT_ID"`
	LoaWkCntrRcpntNm       *string    `json:"loa_wk_cntr_rcpnt_nm" db:"loa_wk_cntr_rcpnt_nm" pipe:"LOA_WK_CNTR_RCPNT_NM"`
	LoaMajRmbsmtSrcID      *string    `json:"loa_maj_rmbsmt_src_id" db:"loa_maj_rmbsmt_src_id" pipe:"LOA_MAJ_RMBSMT_SRC_ID"`
	LoaDtlRmbsmtSrcID      *string    `json:"loa_dtl_rmbsmt_src_id" db:"loa_dtl_rmbsmt_src_id" pipe:"LOA_DTL_RMBSMT_SRC_ID"`
	LoaCustNm              *string    `json:"loa_cust_nm" db:"loa_cust_nm" pipe:"LOA_CUST_NM"`
	LoaObjClsID            *string    `json:"loa_obj_cls_id" db:"loa_obj_cls_id" pipe:"LOA_OBJ_CLS_ID"`
	LoaSrvSrcID            *string    `json:"loa_srv_src_id" db:"loa_srv_src_id" pipe:"LOA_SRV_SRC_ID"`
	LoaSpclIntrID          *string    `json:"loa_spcl_intr_id" db:"loa_spcl_intr_id" pipe:"LOA_SPCL_INTR_ID"`
	LoaBdgtAcntClsNm       *string    `json:"loa_bdgt_acnt_cls_nm" db:"loa_bdgt_acnt_cls_nm" pipe:"LOA_BDGT_ACNT_CLS_NM"`
	LoaDocID               *string    `json:"loa_doc_id" db:"loa_doc_id" pipe:"LOA_DOC_ID"`
	LoaClsRefID            *string    `json:"loa_cls_ref_id" db:"loa_cls_ref_id" pipe:"LOA_CLS_REF_ID"`
	LoaInstlAcntgActID     *string    `json:"loa_instl_acntg_act_id" db:"loa_instl_acntg_act_id" pipe:"LOA_INSTL_ACNTG_ACT_ID"`
	LoaLclInstlID          *string    `json:"loa_lcl_instl_id" db:"loa_lcl_instl_id" pipe:"LOA_LCL_INSTL_ID"`
	LoaFmsTrnsactnID       *string    `json:"loa_fms_trnsactn_id" db:"loa_fms_trnsactn_id" pipe:"LOA_FMS_TRNSACTN_ID"`
	LoaDscTx               *string    `json:"loa_dsc_tx" db:"loa_dsc_tx" pipe:"LOA_DSC_TX"`
	LoaBgnDt               *time.Time `json:"loa_bgn_dt" db:"loa_bgn_dt" pipe:"LOA_BGN_DT"`
	LoaEndDt               *time.Time `json:"loa_end_dt" db:"loa_end_dt" pipe:"LOA_END_DT"`
	LoaFnctPrsNm           *string    `json:"loa_fnct_prs_nm" db:"loa_fnct_prs_nm" pipe:"LOA_FNCT_PRS_NM"`
	LoaStatCd              *string    `json:"loa_stat_cd" db:"loa_stat_cd" pipe:"LOA_STAT_CD"`
	LoaHistStatCd          *string    `json:"loa_hist_stat_cd" db:"loa_hist_stat_cd" pipe:"LOA_HIST_STAT_CD"`
	LoaHsGdsCd             *string    `json:"loa_hs_gds_cd" db:"loa_hs_gds_cd" pipe:"LOA_HS_GDS_CD"`
	OrgGrpDfasCd           *string    `json:"org_grp_dfas_cd" db:"org_grp_dfas_cd" pipe:"ORG_GRP_DFAS_CD"`
	LoaUic                 *string    `json:"loa_uic" db:"loa_uic" pipe:"LOA_UIC"`
	LoaTrnsnID             *string    `json:"loa_trnsn_id" db:"loa_trnsn_id" pipe:"LOA_TRNSN_ID"`
	LoaSubAcntID           *string    `json:"loa_sub_acnt_id" db:"loa_sub_acnt_id" pipe:"LOA_SUB_ACNT_ID"`
	LoaBetCd               *string    `json:"loa_bet_cd" db:"loa_bet_cd" pipe:"LOA_BET_CD"`
	LoaFndTyFgCd           *string    `json:"loa_fnd_ty_fg_cd" db:"loa_fnd_ty_fg_cd" pipe:"LOA_FND_TY_FG_CD"`
	LoaBgtLnItmID          *string    `json:"loa_bgt_ln_itm_id" db:"loa_bgt_ln_itm_id" pipe:"LOA_BGT_LN_ITM_ID"`
	LoaScrtyCoopImplAgncCd *string    `json:"loa_scrty_coop_impl_agnc_cd" db:"loa_scrty_coop_impl_agnc_cd" pipe:"LOA_SCRTY_COOP_IMPL_AGNC_CD"`
	LoaScrtyCoopDsgntrCd   *string    `json:"loa_scrty_coop_dsgntr_cd" db:"loa_scrty_coop_dsgntr_cd" pipe:"LOA_SCRTY_COOP_DSGNTR_CD"`
	LoaScrtyCoopLnItmID    *string    `json:"loa_scrty_coop_ln_itm_id" db:"loa_scrty_coop_ln_itm_id" pipe:"LOA_SCRTY_COOP_LN_ITM_ID"`
	LoaAgncDsbrCd          *string    `json:"loa_agnc_dsbr_cd" db:"loa_agnc_dsbr_cd" pipe:"LOA_AGNC_DSBR_CD"`
	LoaAgncAcntngCd        *string    `json:"loa_agnc_acntng_cd" db:"loa_agnc_acntng_cd" pipe:"LOA_AGNC_ACNTNG_CD"`
	LoaFndCntrID           *string    `json:"loa_fnd_cntr_id" db:"loa_fnd_cntr_id" pipe:"LOA_FND_CNTR_ID"`
	LoaCstCntrID           *string    `json:"loa_cst_cntr_id" db:"loa_cst_cntr_id" pipe:"LOA_CST_CNTR_ID"`
	LoaPrjID               *string    `json:"loa_prj_id" db:"loa_prj_id" pipe:"LOA_PRJ_ID"`
	LoaActvtyID            *string    `json:"loa_actvty_id" db:"loa_actvty_id" pipe:"LOA_ACTVTY_ID"`
	LoaCstCd               *string    `json:"loa_cst_cd" db:"loa_cst_cd" pipe:"LOA_CST_CD"`
	LoaWrkOrdID            *string    `json:"loa_wrk_ord_id" db:"loa_wrk_ord_id" pipe:"LOA_WRK_ORD_ID"`
	LoaFnclArID            *string    `json:"loa_fncl_ar_id" db:"loa_fncl_ar_id" pipe:"LOA_FNCL_AR_ID"`
	LoaScrtyCoopCustCd     *string    `json:"loa_scrty_coop_cust_cd" db:"loa_scrty_coop_cust_cd" pipe:"LOA_SCRTY_COOP_CUST_CD"`
	LoaEndFyTx             *int       `json:"loa_end_fy_tx" db:"loa_end_fy_tx" pipe:"LOA_END_FY_TX"`
	LoaBgFyTx              *int       `json:"loa_bg_fy_tx" db:"loa_bg_fy_tx" pipe:"LOA_BG_FY_TX"`
	LoaBgtRstrCd           *string    `json:"loa_bgt_rstr_cd" db:"loa_bgt_rstr_cd" pipe:"LOA_BGT_RSTR_CD"`
	LoaBgtSubActCd         *string    `json:"loa_bgt_sub_act_cd" db:"loa_bgt_sub_act_cd" pipe:"LOA_BGT_SUB_ACT_CD"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at" pipe:"-"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at" pipe:"-"`
}

// TableName overrides the table name used by Pop.
func (l LineOfAccounting) TableName() string {
	return "lines_of_accounting"
}
