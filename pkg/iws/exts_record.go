package iws

// ExtsRecord contains information used to identify an individual within a system outside of DMDC
type ExtsRecord struct {
	// SydPnID is the identifier used to identify an individual within a system outside of DMDC.
	SydPnID string `xml:"SYS_PN_ID"`
	// SysPnXRef is the identifier that is used to represent cross-reference between a person's System Person Identifiers. If this field is present, this is the newest System Person Identifier to be applied as the SYS_PN_ID.
	SysPnXRef string `xml:"SYS_PN_XR_ID"`
	// SysPnIDTypeCode represents a specific kind of person identifier.
	SysPnIDTypCd SysPnIDTypeCode `xml:"SYS_PN_ID_TYP_CD"`
	Edipi        uint64          `xml:"DOD_EDI_PN_ID"`
}

// SysPnIDTypeCode represents a specific kind of person identifier.
type SysPnIDTypeCode string

const (
	// SysPnIDTypeCodeVAIC indicates a Veterans Administration Integration Control Number
	SysPnIDTypeCodeVAIC SysPnIDTypeCode = "IC"
	// SysPnIDTypeCodeInterim indicates an Interim Person Identifier
	SysPnIDTypeCodeInterim SysPnIDTypeCode = "IP"
)
