package edisegment

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// AffiliationToAgency is a map from our affiliation to the FA1 segment's AgencyQualifierCode field
var AffiliationToAgency = map[models.ServiceMemberAffiliation]string{
	models.AffiliationARMY:     "DZ",
	models.AffiliationNAVY:     "DN",
	models.AffiliationMARINES:  "DX",
	models.AffiliationAIRFORCE: "DY",
}

// FA1 represents the FA1 EDI segment
type FA1 struct {
	AgencyQualifierCode string
}

// StringArray converts FA1 to an array of strings
func (s *FA1) StringArray() []string {
	return []string{"FA1", s.AgencyQualifierCode}
}

// Parse parses an X12 string that's split into an array into the FA1 struct
func (s *FA1) Parse(elements []string) error {
	expectedNumElements := 1
	if len(elements) != expectedNumElements {
		return fmt.Errorf("FA1: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.AgencyQualifierCode = elements[0]
	return nil
}
