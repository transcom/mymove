package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// RunScenarioThree creates some duty stations.
func RunScenarioThree(db *pop.Connection) {
	// These are some random Marine bases found on the internet, not for production use.
	MakeDutyStation(db, "Air Station Yuma", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Yuma", State: "Arizona", PostalCode: "85364"})
	MakeDutyStation(db, "Air Station Miramar", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "San Diego", State: "California", PostalCode: "92145"})
	MakeDutyStation(db, "Air Station Cherry Point", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Cherry Point", State: "North Carolina", PostalCode: "28533"})
	MakeDutyStation(db, "Air Station New River", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Jacksonville", State: "North Carolina", PostalCode: "28540"})
	MakeDutyStation(db, "Air Station Beaufort", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Beaufort", State: "South Carolina", PostalCode: "29904"})
	MakeDutyStation(db, "Air Ground Combat Center Twentynine Palms", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Twentynine Palms", State: "California", PostalCode: "92278"})
	MakeDutyStation(db, "Base Camp Pendleton", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Oceanside", State: "California", PostalCode: "92058"})
	MakeDutyStation(db, "Recruit Depot San Diego", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "San Diego", State: "California", PostalCode: "92140"})
	MakeDutyStation(db, "Base Hawaii", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Honolulu", State: "Hawaii", PostalCode: "96734"})
	MakeDutyStation(db, "Base Camp Lejeune", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Lejeune", State: "North Carolina", PostalCode: "28547"})
	MakeDutyStation(db, "Recruit Depot Parris Island", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Parris Island", State: "South Carolina", PostalCode: "29905"})
	MakeDutyStation(db, "Base Quantico", internalmessages.MilitaryBranchMARINES,
		models.Address{StreetAddress1: "duty station", City: "Quantico", State: "Virginia", PostalCode: "22134"})
}
