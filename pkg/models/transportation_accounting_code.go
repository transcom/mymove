package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// TransportationAccountingCode model struct that represents transportation accounting codes
type TransportationAccountingCode struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	TAC                string            `json:"tac" db:"tac"`
	LoaID              *uuid.UUID        `json:"loa_id" db:"loa_id"`
	LineOfAccounting   *LineOfAccounting `belongs_to:"lines_of_accounting" fk_id:"loa_id"`
	TacSysID           *int              `json:"tac_sys_id" db:"tac_sys_id"`
	LoaSysID           *int              `json:"loa_sys_id" db:"loa_sys_id"`
	TacFyTxt           *int              `json:"tac_fy_txt" db:"tac_fy_txt"`
	TacFnBlModCd       *string           `json:"tac_fn_bl_mod_cd" db:"tac_fn_bl_mod_cd"`
	OrgGrpDfasCd       *string           `json:"org_grp_dfas_cd" db:"org_grp_dfas_cd"`
	TacMvtDsgID        *string           `json:"tac_mvt_dsg_id" db:"tac_mvt_dsg_id"`
	TacTyCd            *string           `json:"tac_ty_cd" db:"tac_ty_cd"`
	TacUseCd           *string           `json:"tac_use_cd" db:"tac_use_cd"`
	TacMajClmtID       *string           `json:"tac_maj_clmt_id" db:"tac_maj_clmt_id"`
	TacBillActTxt      *string           `json:"tac_bill_act_txt" db:"tac_bill_act_txt"`
	TacCostCtrNm       *string           `json:"tac_cost_ctr_nm" db:"tac_cost_ctr_nm"`
	Buic               *string           `json:"buic" db:"buic"`
	TacHistCd          *string           `json:"tac_hist_cd" db:"tac_hist_cd"`
	TacStatCd          *string           `json:"tac_stat_cd" db:"tac_stat_cd"`
	TrnsprtnAcntTx     *string           `json:"trnsprtn_acnt_tx" db:"trnsprtn_acnt_tx"`
	TrnsprtnAcntBgnDt  *time.Time        `json:"trnsprtn_acnt_bgn_dt" db:"trnsprtn_acnt_bgn_dt"`
	TrnsprtnAcntEndDt  *time.Time        `json:"trnsprtn_acnt_end_dt" db:"trnsprtn_acnt_end_dt"`
	DdActvtyAdrsID     *string           `json:"dd_actvty_adrs_id" db:"dd_actvty_adrs_id"`
	TacBlldAddFrstLnTx *string           `json:"tac_blld_add_frst_ln_tx" db:"tac_blld_add_frst_ln_tx"`
	TacBlldAddScndLnTx *string           `json:"tac_blld_add_scnd_ln_tx" db:"tac_blld_add_scnd_ln_tx"`
	TacBlldAddThrdLnTx *string           `json:"tac_blld_add_thrd_ln_tx" db:"tac_blld_add_thrd_ln_tx"`
	TacBlldAddFrthLnTx *string           `json:"tac_blld_add_frth_ln_tx" db:"tac_blld_add_frth_ln_tx"`
	TacFnctPocNm       *string           `json:"tac_fnct_poc_nm" db:"tac_fnct_poc_nm"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`
}

// TODO Validate required fields?

// TableName overrides the table name used by Pop.
func (t TransportationAccountingCode) TableName() string {
	return "transportation_accounting_codes"
}
