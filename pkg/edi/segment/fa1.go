package edisegment

import (
	"fmt"
	"strings"

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

// String converts FA1 to its X12 single line string representation
func (s *FA1) String(delimiter string) string {
	return strings.Join([]string{"FA1", s.AgencyQualifierCode}, delimiter) + "\n"
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
