package models

import "time"

// !IMPORTANT! This struct is not the file record, see TransportationAccountingCodeTrdmFileRecord model for this.
// This struct is what will be returned when the TRDM .txt file gets parsed.
// See TransportationAccountingCodeTrdmFileRecord for the struct representing the .txt file.
// The below columns are the chosen fields necessary from the parsed file.
// TRNSPRTN_ACNT_CD, TAC_BLLD_ADD_FRST_LN_TX, TAC_BLLD_ADD_SCND_LN_TX, TAC_BLLD_ADD_THRD_LN_TX,
// TAC_BLLD_ADD_FRTH_LN_TX, TRNSPRTN_ACNT_TX, TRNSPRTN_ACNT_BGN_DT, TRNSPRTN_ACNT_END_DT
// TAC_FY_TXT

type TransportationAccountingCodeDesiredFromTRDM struct {
	TAC/*Third in line, values[2]*/ string                      `json:"tac"`
	BillingAddressFirstLine/*20th in line, values[19]*/ string  `json:"billing_address_first_line"`
	BillingAddressSecondLine/*21st in line, values[20]*/ string `json:"billing_address_second_line"`
	BillingAddressThirdLine/*22nd in line, values[21]*/ string  `json:"billing_address_third_line"`
	BillingAddressFourthLine/*23rd in line, values[22]*/ string `json:"billing_address_fourth_line"`
	Transaction/*16th in line, values[15]*/ string              `json:"transaction"`
	EffectiveDate/*17th in line, values[16]*/ time.Time         `json:"effective_date"`
	ExpirationDate/*18th in line, values[17]*/ time.Time        `json:"expiration_date"`
	FiscalYear/*4th in line, values[3]*/ string                 `json:"fiscal_year"`
}
