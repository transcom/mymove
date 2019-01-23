package iws

// SSN is a test SSN value
const SSN = "666839559"
const edipi = 1234567890

// TestingPersonLookup is a mock of RBS that returns dummy data
type TestingPersonLookup struct{}

// GetPersonUsingEDIPI returns a static dummy RBS result
func (r TestingPersonLookup) GetPersonUsingEDIPI(edipi uint64) (*Person, []Personnel, error) {
	return getTestPerson(), []Personnel{getTestPersonnel()}, nil
}

// GetPersonUsingSSN returns a static dummy RBS result
func (r TestingPersonLookup) GetPersonUsingSSN(params GetPersonUsingSSNParams) (MatchReasonCode, uint64, *Person, []Personnel, error) {
	return MatchReasonCodeFull, edipi, getTestPerson(), []Personnel{getTestPersonnel()}, nil
}

// GetPersonUsingWorkEmail returns a static dummy RBS result
func (r TestingPersonLookup) GetPersonUsingWorkEmail(workEmail string) (uint64, *Person, []Personnel, error) {
	return edipi, getTestPerson(), []Personnel{getTestPersonnel()}, nil
}

func getTestPerson() *Person {
	person := Person{
		ID:         SSN,
		TypeCode:   PersonTypeCodeSSN,
		LastName:   "McTestface",
		FirstName:  "Testy",
		MiddleName: "Test",
		CdncyName:  "",
		BirthDate:  "19900101",
	}

	return &person
}

func getTestPersonnel() Personnel {
	return Personnel{
		PnlCatCd:  PersonnelCategoryCodeActiveDuty,
		OrgCd:     OrgCodeAirForceActive,
		Email:     "testy.mctestface@example.com",
		RankCd:    "MSGT",
		PgCd:      PayGradeCode07,
		PayPlanCd: PayPlanCodeCG,
		SvcCd:     ServiceCodeAirForce,
	}
}
