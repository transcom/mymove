package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// TransportationAccountingCode model struct that represents transportation accounting codes
// TODO: Update this model and internal use to reflect incoming TransportationAccountingCode model updates.
// Don't forget to update the MakeDefaultTransportationAccountingCode function inside of the testdatagen package.
type TransportationAccountingCode struct {
	ID                 uuid.UUID         `json:"id" db:"id" pipe:"-"` // Internal
	TAC                string            `json:"tac" db:"tac" pipe:"TRNSPRTN_ACNT_CD"`
	LoaID              *uuid.UUID        `json:"loa_id" db:"loa_id" pipe:"-"`                       // Internal
	LineOfAccounting   *LineOfAccounting `belongs_to:"lines_of_accounting" fk_id:"loa_id" pipe:"-"` // Internal
	TacSysID           *string           `json:"tac_sys_id" db:"tac_sys_id" pipe:"TAC_SYS_ID"`
	LoaSysID           *string           `json:"loa_sys_id" db:"loa_sys_id" pipe:"LOA_SYS_ID"`
	TacFyTxt           *string           `json:"tac_fy_txt" db:"tac_fy_txt" pipe:"TAC_FY_TXT"`
	TacFnBlModCd       *string           `json:"tac_fn_bl_mod_cd" db:"tac_fn_bl_mod_cd" pipe:"TAC_FN_BL_MOD_CD"`
	OrgGrpDfasCd       *string           `json:"org_grp_dfas_cd" db:"org_grp_dfas_cd" pipe:"ORG_GRP_DFAS_CD"`
	TacMvtDsgID        *string           `json:"tac_mvt_dsg_id" db:"tac_mvt_dsg_id" pipe:"TAC_MVT_DSG_ID"`
	TacTyCd            *string           `json:"tac_ty_cd" db:"tac_ty_cd" pipe:"TAC_TY_CD"`
	TacUseCd           *string           `json:"tac_use_cd" db:"tac_use_cd" pipe:"TAC_USE_CD"`
	TacMajClmtID       *string           `json:"tac_maj_clmt_id" db:"tac_maj_clmt_id" pipe:"TAC_MAJ_CLMT_ID"`
	TacBillActTxt      *string           `json:"tac_bill_act_txt" db:"tac_bill_act_txt" pipe:"TAC_BILL_ACT_TXT"`
	TacCostCtrNm       *string           `json:"tac_cost_ctr_nm" db:"tac_cost_ctr_nm" pipe:"TAC_COST_CTR_NM"`
	Buic               *string           `json:"buic" db:"buic" pipe:"BUIC"`
	TacHistCd          *string           `json:"tac_hist_cd" db:"tac_hist_cd" pipe:"TAC_HIST_CD"`
	TacStatCd          *string           `json:"tac_stat_cd" db:"tac_stat_cd" pipe:"TAC_STAT_CD"`
	TrnsprtnAcntTx     *string           `json:"trnsprtn_acnt_tx" db:"trnsprtn_acnt_tx" pipe:"TRNSPRTN_ACNT_TX"`
	TrnsprtnAcntBgnDt  *time.Time        `json:"trnsprtn_acnt_bgn_dt" db:"trnsprtn_acnt_bgn_dt" pipe:"TRNSPRTN_ACNT_BGN_DT"`
	TrnsprtnAcntEndDt  *time.Time        `json:"trnsprtn_acnt_end_dt" db:"trnsprtn_acnt_end_dt" pipe:"TRNSPRTN_ACNT_END_DT"`
	DdActvtyAdrsID     *string           `json:"dd_actvty_adrs_id" db:"dd_actvty_adrs_id" pipe:"DD_ACTVTY_ADRS_ID"`
	TacBlldAddFrstLnTx *string           `json:"tac_blld_add_frst_ln_tx" db:"tac_blld_add_frst_ln_tx" pipe:"TAC_BLLD_ADD_FRST_LN_TX"`
	TacBlldAddScndLnTx *string           `json:"tac_blld_add_scnd_ln_tx" db:"tac_blld_add_scnd_ln_tx" pipe:"TAC_BLLD_ADD_SCND_LN_TX"`
	TacBlldAddThrdLnTx *string           `json:"tac_blld_add_thrd_ln_tx" db:"tac_blld_add_thrd_ln_tx" pipe:"TAC_BLLD_ADD_THRD_LN_TX"`
	TacBlldAddFrthLnTx *string           `json:"tac_blld_add_frth_ln_tx" db:"tac_blld_add_frth_ln_tx" pipe:"TAC_BLLD_ADD_FRTH_LN_TX"`
	TacFnctPocNm       *string           `json:"tac_fnct_poc_nm" db:"tac_fnct_poc_nm" pipe:"TAC_FNCT_POC_NM"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at" pipe:"-"` // Internal
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at" pipe:"-"` // Internal
}

func (t *TransportationAccountingCode) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.TAC, Name: "TAC"},
	), nil
}

// TableName overrides the table name used by Pop.
func (t TransportationAccountingCode) TableName() string {
	return "transportation_accounting_codes"
}
