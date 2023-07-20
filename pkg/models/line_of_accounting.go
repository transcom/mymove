package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type LineOfAccounting struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
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
}

// TODO Validate required fields?

// TableName overrides the table name used by Pop.
func (t LineOfAccounting) TableName() string {
	return "lines_of_accounting"
}
