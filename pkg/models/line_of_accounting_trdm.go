package models

import (
	"time"
)

// This file declares all of the column values that are received from TRDM in regard to the pipe delimited
// Line_Of_Accounting (LOA) .txt file

// This struct only applies to the received .txt file.
// 57 values pulled from TRDM
type LineOfAccountingTrdmFileRecord struct {
	LOA_SYS_ID                  string // Yes, this is a string not an int to stay in line with the matrix. The LineOfAccounting struct uses an int
	LOA_DPT_ID                  string
	LOA_TNSFR_DPT_NM            string
	LOA_BAF_ID                  string
	LOA_TRSY_SFX_TX             string
	LOA_MAJ_CLM_NM              string
	LOA_OP_AGNCY_ID             string
	LOA_ALLT_SN_ID              string
	LOA_PGM_ELMNT_ID            string
	LOA_TSK_BDGT_SBLN_TX        string
	LOA_DF_AGNCY_ALCTN_RCPNT_ID string
	LOA_JB_ORD_NM               string
	LOA_SBALTMT_RCPNT_ID        string
	LOA_WK_CNTR_RCPNT_NM        string
	LOA_MAJ_RMBSMT_SRC_ID       string
	LOA_DTL_RMBSMT_SRC_ID       string
	LOA_CUST_NM                 string
	LOA_OBJ_CLS_ID              string
	LOA_SRV_SRC_ID              string
	LOA_SPCL_INTR_ID            string
	LOA_BDGT_ACNT_CLS_NM        string
	LOA_DOC_ID                  string
	LOA_CLS_REF_ID              string
	LOA_INSTL_ACNTG_ACT_ID      string
	LOA_LCL_INSTL_ID            string
	LOA_FMS_TRNSACTN_ID         string
	LOA_DSC_TX                  string
	LOA_BGN_DT                  time.Time
	LOA_END_DT                  time.Time
	LOA_FNCT_PRS_NM             string
	LOA_STAT_CD                 string
	LOA_HIST_STAT_CD            string
	LOA_HS_GDS_CD               string
	ORG_GRP_DFAS_CD             string
	LOA_UIC                     string
	LOA_TRNSN_ID                string
	LOA_SUB_ACNT_ID             string
	LOA_BET_CD                  string
	LOA_FND_TY_FG_CD            string
	LOA_BGT_LN_ITM_ID           string
	LOA_SCRTY_COOP_IMPL_AGNC_CD string
	LOA_SCRTY_COOP_DSGNTR_CD    string
	LOA_SCRTY_COOP_LN_ITM_ID    string
	LOA_AGNC_DSBR_CD            string
	LOA_AGNC_ACNTNG_CD          string
	LOA_FND_CNTR_ID             string
	LOA_CST_CNTR_ID             string
	LOA_PRJ_ID                  string
	LOA_ACTVTY_ID               string
	LOA_CST_CD                  string
	LOA_WRK_ORD_ID              string
	LOA_FNCL_AR_ID              string
	LOA_SCRTY_COOP_CUST_CD      string
	LOA_END_FY_TX               int
	LOA_BG_FY_TX                int
	LOA_BGT_RSTR_CD             string
	LOA_BGT_SUB_ACT_CD          string
}
