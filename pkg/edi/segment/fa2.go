package edisegment

import (
	"fmt"
)

type FA2DetailCode string

func (d FA2DetailCode) String() string {
	return string(d)
}

const (
	// FA2DetailCodeTA is Transportation Account Code (TAC)
	FA2DetailCodeTA FA2DetailCode = "TA"
	// FA2DetailCodeZZ is Mutually Defined
	FA2DetailCodeZZ FA2DetailCode = "ZZ"
	// FA2DetailCodeA1 is Department Indicator
	FA2DetailCodeA1 FA2DetailCode = "A1"
	// FA2DetailCodeA2 is Transfer from Department
	FA2DetailCodeA2 FA2DetailCode = "A2"
	// FA2DetailCodeA3 is Fiscal Year Indicator
	FA2DetailCodeA3 FA2DetailCode = "A3"
	// FA2DetailCodeA4 is Basic Symbol Number
	FA2DetailCodeA4 FA2DetailCode = "A4"
	// FA2DetailCodeA5 is Sub-class
	FA2DetailCodeA5 FA2DetailCode = "A5"
	// FA2DetailCodeA6 is Sub-Account Symbol
	FA2DetailCodeA6 FA2DetailCode = "A6"
	// FA2DetailCodeB1 is Budget Activity Number
	FA2DetailCodeB1 FA2DetailCode = "B1"
	// FA2DetailCodeB2 is Budget Sub-activity Number
	FA2DetailCodeB2 FA2DetailCode = "B2"
	// FA2DetailCodeB3 is Budget Program Activity
	FA2DetailCodeB3 FA2DetailCode = "B3"
	// FA2DetailCodeC1 is Program Element
	FA2DetailCodeC1 FA2DetailCode = "C1"
	// FA2DetailCodeC2 is Project Task or Budget Subline
	FA2DetailCodeC2 FA2DetailCode = "C2"
	// FA2DetailCodeD1 is Defense Agency Allocation Recipient
	FA2DetailCodeD1 FA2DetailCode = "D1"
	// FA2DetailCodeD4 is Component Sub-allocation Recipient
	FA2DetailCodeD4 FA2DetailCode = "D4"
	// FA2DetailCodeD6 is Sub-allotment Recipient
	FA2DetailCodeD6 FA2DetailCode = "D6"
	// FA2DetailCodeD7 is Work Center Recipient
	FA2DetailCodeD7 FA2DetailCode = "D7"
	// FA2DetailCodeE1 is Major Reimbursement Source Code
	FA2DetailCodeE1 FA2DetailCode = "E1"
	// FA2DetailCodeE2 is Detail Reimbursement Source Code
	FA2DetailCodeE2 FA2DetailCode = "E2"
	// FA2DetailCodeE3 is Customer Indicator
	FA2DetailCodeE3 FA2DetailCode = "E3"
	// FA2DetailCodeF1 is Object Class
	FA2DetailCodeF1 FA2DetailCode = "F1"
	// FA2DetailCodeF3 is Government or Public Sector Identifier
	FA2DetailCodeF3 FA2DetailCode = "F3"
	// FA2DetailCodeG2 is Special Interest Code or Special Program Cost Code
	FA2DetailCodeG2 FA2DetailCode = "G2"
	// FA2DetailCodeI1 is Abbreviated Department of Defense (DoD) Budget and Accounting Classification Code (BACC)
	FA2DetailCodeI1 FA2DetailCode = "I1"
	// FA2DetailCodeJ1 is Document or Record Reference Number
	FA2DetailCodeJ1 FA2DetailCode = "J1"
	// FA2DetailCodeK6 is Accounting Classification Reference Code
	FA2DetailCodeK6 FA2DetailCode = "K6"
	// FA2DetailCodeL1 is Accounting Installation Number
	FA2DetailCodeL1 FA2DetailCode = "L1"
	// FA2DetailCodeM1 is Local Installation Data
	FA2DetailCodeM1 FA2DetailCode = "M1"
	// FA2DetailCodeN1 is Transaction Type
	FA2DetailCodeN1 FA2DetailCode = "N1"
	// FA2DetailCodeP5 is Security Cooperation Case Line Item Identifier
	FA2DetailCodeP5 FA2DetailCode = "P5"
)

// FA2 represents the FA2 EDI segment
type FA2 struct {
	BreakdownStructureDetailCode FA2DetailCode `validate:"oneof=TA ZZ A1 A2 A3 A4 A5 A6 B1 B2 B3 C1 C2 D1 D4 D6 D7 E1 E2 E3 F1 F3 G2 I1 J1 K6 L1 M1 N1 P5"`
	FinancialInformationCode     string        `validate:"min=1,max=80"`
}

// StringArray converts FA2 to an array of strings
func (s *FA2) StringArray() []string {
	return []string{"FA2", s.BreakdownStructureDetailCode.String(), s.FinancialInformationCode}
}

// Parse parses an X12 string that's split into an array into the FA2 struct
func (s *FA2) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("fA2: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.BreakdownStructureDetailCode = FA2DetailCode(elements[0])
	s.FinancialInformationCode = elements[1]
	return nil
}
