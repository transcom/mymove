package iws

// PersonLookup is the interface used to look up a service member in DEERS
type PersonLookup interface {
	GetPersonUsingEDIPI(edipi uint64) (*Person, []Personnel, error)
	GetPersonUsingSSN(params GetPersonUsingSSNParams) (MatchReasonCode, uint64, *Person, []Personnel, error)
	GetPersonUsingWorkEmail(workEmail string) (uint64, *Person, []Personnel, error)
}
