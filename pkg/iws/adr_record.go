package iws

// AdrRecord is the container tag for the majority of the response data from DMDC's Identity Web Services: Real-time Broker Service REST API
type AdrRecord struct {
	// The identifier that is used to represent the person within a Department of Defense Electronic Data Interchange. Externally the EDI-PI is referred to as the DoD ID, or the DoD ID Number.&#13; XML Tag - dodEdiPersonId
	Edipi *uint64 `xml:"DOD_EDI_PN_ID,omitempty"` // <xsd:element minOccurs="0" name="DOD_EDI_PN_ID" type="tns:DOD_EDI_PN_ID"/>
	// The identifier that is used to represent cross-reference between a person's Department of Defense Electronic Data Interchange identifiers. If the code is invalidated, this is the new DoD EDI PN ID to use instead of the current one. This ID will be zero unless the INVL_DEPI_NTFCN_CD is Y.
	EdipiXRef *uint64 `xml:"DOD_EDI_PN_XR_ID,omitempty"` // <xsd:element minOccurs="0" name="DOD_EDI_PN_XR_ID" type="tns:DOD_EDI_PN_XR_ID"/>
	// The date the customer ended their association with ADR/ADW. - Date format is YYYYMMDD
	CstrAscEndDt string `xml:"CSTR_ASC_END_DT,omitempty"` // <xsd:element minOccurs="0" name="CSTR_ASC_END_DT" type="tns:CSTR_ASC_END_DT"/>
	// The code that represents the reason that the customer's association with ADR/ADW ended or is expected to end (see PN_LOSS_RSN_CD).
	CstrAscErsnCd *CustomerAssocEndReasonCode `xml:"CSTR_ASC_ERSN_CD,omitempty"` // <xsd:element minOccurs="0" name="CSTR_ASC_ERSN_CD" type="tns:CSTR_ASC_ERSN_CD"/>
	PidsRecord    *PidsRecord                 `xml:"PIDSRecord,omitempty"`       // <xsd:element minOccurs="0" name="PIDSRecord" type="tns:PIDSRecord"/>
	TidsRecord    *TidsRecord                 `xml:"TIDSRecord,omitempty"`       // <xsd:element minOccurs="0" name="TIDSRecord" type="tns:TIDSRecord"/>
	ExtsRecord    *ExtsRecord                 `xml:"EXTSRecord,omitempty"`       // <xsd:element minOccurs="0" name="EXTSRecord" type="tns:EXTSRecord"/>
	OldEdipis     []uint64                    `xml:"identifierHistory>OLD_DOD_EDI_PN_ID,omitempty"`
	WorkEmail     *WkEmaRecord                `xml:"WKEMARecord,omitempty"`
	Person        *Person                     `xml:"person,omitempty"`
	Personnel     []Personnel                 `xml:"personnel,omitempty"`
}

// CustomerAssocEndReasonCode represents the reason that the customer's association with ADR/ADW ended or is expected to end (see PN_LOSS_RSN_CD).
type CustomerAssocEndReasonCode string

const (
	// CustomerAssocEndReasonCodeNotInPopulation means that the provided ID is associated with a person who is not in the population supported by this integration
	CustomerAssocEndReasonCodeNotInPopulation CustomerAssocEndReasonCode = "N"
	// CustomerAssocEndReasonCodeSeparated means that the provided ID is associated with a person who has separated
	CustomerAssocEndReasonCodeSeparated CustomerAssocEndReasonCode = "S"
	// CustomerAssocEndReasonCodeIneligible means that the provided ID is associated with a person who is not eligible
	CustomerAssocEndReasonCodeIneligible CustomerAssocEndReasonCode = "W"
	// CustomerAssocEndReasonCodeNoLongerMatches means that the search no longer matches customer criteria.
	CustomerAssocEndReasonCodeNoLongerMatches CustomerAssocEndReasonCode = "Y"
)
