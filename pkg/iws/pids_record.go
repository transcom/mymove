package iws

// PidsRecord contains the match reason code and optionally the matched EDIPI for a PIDS query to IWS: RBS
type PidsRecord struct {
	MtchRsnCd MatchReasonCode `xml:"MTCH_RSN_CD"`
	Edipi     uint64          `xml:"DOD_EDI_PN_ID,omitempty"`
}

// MatchReasonCode indicates the reason a DOD_EDI_PN_ID could or could not be returned.
// Reason codes that start with "P" are returned for Person Inquiries (SRC_DS_CD=P).
// Reason codes that start with "D" are returned for Dependent Inquiries (SRC_DS_CD = D).
type MatchReasonCode string

const (
	// MatchReasonCodeDAB means that more than one dependent matched the provided criteria, with sponsor identified by SPN_PN_ID and SPN_PN_ID_TYP_CD only
	MatchReasonCodeDAB MatchReasonCode = "DAB"
	// MatchReasonCodeDAC means that more than one dependent matched the provided criteria, with sponsor identified by SPN_PN_ID, SPN_PN_ID_TYP_CD and at least one additional criterion
	MatchReasonCodeDAC MatchReasonCode = "DAC"
	// MatchReasonCodeMultipleSponsors means that more than one SPN_PN_ID matched the provided criteria
	MatchReasonCodeMultipleSponsors MatchReasonCode = "DAS"
	// MatchReasonCodeDMB means that a dependent matched on at least one criterion, with sponsor identified by SPN_PN_ID and SPN_PN_ID_TYP_CD only
	MatchReasonCodeDMB MatchReasonCode = "DMB"
	// MatchReasonCodeDMC means that a dependent matched on at least one criterion, with sponsor identified by SPN_PN_ID, SPN_PN_ID_TYP_CD and at least one additional criterion
	MatchReasonCodeDMC MatchReasonCode = "DMC"
	// MatchReasonCodeDNB means that a sponsor was found using SPN_PN_ID and SPN_PN_ID_TYP_CD, but no dependents for this sponsor could be found that matched any of the provided criteria
	MatchReasonCodeDNB MatchReasonCode = "DNB"
	// MatchReasonCodeDNC means that a sponsor was found using SPN_PN_ID and SPN_PN_ID_TYP_CD and at least one additional criterion, but no dependents for this sponsor could be found that matched any of the provided criteria
	MatchReasonCodeDNC MatchReasonCode = "DNC"
	// MatchReasonCodeNoMatchingSponsor means that no sponsor matched the provided SPN_PN_ID and SPN_PN_ID_TYP_CD combination
	MatchReasonCodeNoMatchingSponsor MatchReasonCode = "DNS"
	// MatchReasonCodeMultiple means that more than one PN_ID matched the provided criteria
	MatchReasonCodeMultiple MatchReasonCode = "PAB"
	// MatchReasonCodeLimited means that the person matched on PN_ID and PN_ID_TYP_CD only
	MatchReasonCodeLimited MatchReasonCode = "PMB"
	// MatchReasonCodeFull means that the person matched on PN_ID, PN_ID_TYP_CD and at least one additional criterion
	MatchReasonCodeFull MatchReasonCode = "PMC"
	// MatchReasonCodeNone means that no person matched the provided PN_ID and PN_ID_TYP_CD combination
	MatchReasonCodeNone MatchReasonCode = "PNB"
)
