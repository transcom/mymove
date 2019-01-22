package iws

// Person contains the PII returned by an EDI query or use to search in a PIDS query with DMDC's Identity Web Services: Real-time Broker Service REST API
type Person struct {
	ID         string         `xml:"PN_ID"`
	TypeCode   PersonTypeCode `xml:"PN_ID_TYP_CD"`
	LastName   string         `xml:"PN_LST_NM,omitempty"`
	FirstName  string         `xml:"PN_1ST_NM,omitempty"`
	MiddleName string         `xml:"PN_MID_NM,omitempty"`
	CdncyName  string         `xml:"PN_CDNCY_NM,omitempty"`
	BirthDate  string         `xml:"PN_BRTH_DT,omitempty"`
}

// PersonTypeCode is the code that represents a specific kind of person identifier.
type PersonTypeCode string

const (
	// PersonTypeCodeDODBenefitNum indicates a DOD Benefit Number
	PersonTypeCodeDODBenefitNum PersonTypeCode = "B"
	// PersonTypeCodePlaceholder indicates a special 9-digit code created for individuals (i.e., babies) who do not have or have not provided an SSN when the record is added to DEERS (dependents only)
	PersonTypeCodePlaceholder PersonTypeCode = "D"
	// PersonTypeCodeEDIPI indicates a 10-digit Electronic Data Interchange Identifier; i.e., a DoD ID Number
	PersonTypeCodeEDIPI PersonTypeCode = "E"
	// PersonTypeCodeForeign indicates a special 9-digit code created for foreign military and nationals
	PersonTypeCodeForeign PersonTypeCode = "F"
	// PersonTypeCodeTaxID indicates a tax identification number
	PersonTypeCodeTaxID PersonTypeCode = "I"
	// PersonTypeCodePatient indicates a Patient Identifier
	PersonTypeCodePatient PersonTypeCode = "M"
	// PersonTypeCodeInvalid indicates an invalid SSN. The PN_ID was submitted as an SSN, but does not conform to the valid SSN structure. Obsolete value, no longer applied.
	PersonTypeCodeInvalid PersonTypeCode = "N"
	// PersonTypeCodePreSSNMilitary indicates a special 9-digit code created for U.S. military personnel from Service Numbers before the switch to Social Security Numbers
	PersonTypeCodePreSSNMilitary PersonTypeCode = "P"
	// PersonTypeCodeShyContractor indicates a special 9-digit code created for a DoD contractor who refused to give his or her SSN to RAPIDS; the associated PN_ID will begin with 99
	PersonTypeCodeShyContractor PersonTypeCode = "R"
	// PersonTypeCodeSSN indicates a 9-digit Social Security Number
	PersonTypeCodeSSN PersonTypeCode = "S"
	// PersonTypeCodeTest indicates a Test (858 series) identifier
	PersonTypeCodeTest PersonTypeCode = "T"
	// PersonTypeCodeNotAPersonID indicates "Not a Person Identifier" (Used only in DoD Bar Codes)
	PersonTypeCodeNotAPersonID PersonTypeCode = "X"
)
