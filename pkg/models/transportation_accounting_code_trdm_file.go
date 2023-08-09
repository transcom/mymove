package models

// This file declares all of the column values that are received from TRDM in regard to the pipe delimited
// Transportation_Accounting_Code (TAC) .txt file

// This struct only applies to the received .txt file.
//
//nolint:revive
type TransportationAccountingCodeTrdmFileRecord struct {
	TAC_SYS_ID              string
	LOA_SYS_ID              string
	TRNSPRTN_ACNT_CD        string
	TAC_FY_TXT              string
	TAC_FN_BL_MOD_CD        string
	ORG_GRP_DFAS_CD         string
	TAC_MVT_DSG_ID          string
	TAC_TY_CD               string
	TAC_USE_CD              string
	TAC_MAJ_CLMT_ID         string
	TAC_BILL_ACT_TXT        string
	TAC_COST_CTR_NM         string
	BUIC                    string
	TAC_HIST_CD             string
	TAC_STAT_CD             string
	TRNSPRTN_ACNT_TX        string
	TRNSPRTN_ACNT_BGN_DT    string
	TRNSPRTN_ACNT_END_DT    string
	DD_ACTVTY_ADRS_ID       string
	TAC_BLLD_ADD_FRST_LN_TX string
	TAC_BLLD_ADD_SCND_LN_TX string
	TAC_BLLD_ADD_THRD_LN_TX string
	TAC_BLLD_ADD_FRTH_LN_TX string
	TAC_FNCT_POC_NM         string
}
