package iws

// TidsRecord contains information related to a TOKEN query
type TidsRecord struct {
	// TidsMtchRsnCd indicates the reason code a TOKEN could or could not be returned.
	TidsMtchRsnCd TidsMatchReasonCode `xml:"TIDS_MTCH_RSN_CD"`
	Edipi         uint64              `xml:"DOD_EDI_PN_ID"`
	// OrgID is the identifier of the organization.
	OrgID       OrgID                `xml:"ORG_ID"`
	OrgAscCatCd OrgAssocCategoryCode `xml:"ORG_ASC_CAT_CD"`
}

// TidsMatchReasonCode is the reason code a TOKEN could or could not be returned.
type TidsMatchReasonCode string

const (
	// TidsMatchReasonCodeValid indicates that the Token is valid and current.
	TidsMatchReasonCodeValid TidsMatchReasonCode = "M"
	// TidsMatchReasonCodeValidMultiple indicates that the Token(s) is valid for multiple individuals (data for most recently issued token is returned; token may be a partial match).
	TidsMatchReasonCodeValidMultiple TidsMatchReasonCode = "A"
	// TidsMatchReasonCodeLimited indicates that the Token is valid and current for person & ORG_ID only.
	TidsMatchReasonCodeLimited TidsMatchReasonCode = "P"
	// TidsMatchReasonCodeInvalid indicates that the Token is invalid or expired.
	TidsMatchReasonCodeInvalid TidsMatchReasonCode = "N"
	// TidsMatchReasonCodeLost that the Token is Lost or Stolen.
	TidsMatchReasonCodeLost TidsMatchReasonCode = "L"
	// TidsMatchReasonCodeTerminated indicates that the Token has been Terminated.
	TidsMatchReasonCodeTerminated TidsMatchReasonCode = "T"
	// TidsMatchReasonCodeExpired indicates that the Token has Expired.
	TidsMatchReasonCodeExpired TidsMatchReasonCode = "E"
)

// OrgID is the 4 digit identifier of the organization.
type OrgID string

const (
	// OrgIDNOAA identifies the National Oceanic and Atmospheric Administration (NOAA)
	OrgIDNOAA OrgID = "1330"
	// OrgIDNavy identifies the United States Navy
	OrgIDNavy OrgID = "1700"
	// OrgIDMarineCorps identifies the United States Marine Corps
	OrgIDMarineCorps OrgID = "1727"
	// OrgIDArmy identifies the United States Army
	OrgIDArmy OrgID = "2100"
	// OrgIDAirForce identifies the United States Air Force
	OrgIDAirForce OrgID = "5700"
	// OrgIDCoastGuard identifies the United States Coast Guard
	OrgIDCoastGuard OrgID = "6950"
	// OrgIDPublicHealth identifies the United States Public Health Service
	OrgIDPublicHealth OrgID = "7520"
	// OrgIDDeptOfDefense identifies the Department of Defense
	OrgIDDeptOfDefense OrgID = "9700"
)

// OrgAssocCategoryCode represents the category of the organization association.
type OrgAssocCategoryCode string

const (
	// OrgAssocCategoryCodeCivil means Civilian Employee
	OrgAssocCategoryCodeCivil OrgAssocCategoryCode = "01"
	// OrgAssocCategoryCodeAppointee means Political Appointee/SES
	OrgAssocCategoryCodeAppointee OrgAssocCategoryCode = "02"
	// OrgAssocCategoryCodeUniformedService means Uniformed Service Member
	OrgAssocCategoryCodeUniformedService OrgAssocCategoryCode = "03"
	// OrgAssocCategoryCodeContractor means Contractor
	OrgAssocCategoryCodeContractor OrgAssocCategoryCode = "04"
	// OrgAssocCategoryCodeAffiliate means Affiliate (e.g., Foreign Military or National / Federal or Non-Federal Agency)
	OrgAssocCategoryCodeAffiliate OrgAssocCategoryCode = "05"
	// OrgAssocCategoryCodeBeneficiary means Beneficiary (e.g., Retiree, Family Member)
	OrgAssocCategoryCodeBeneficiary OrgAssocCategoryCode = "06"
)
