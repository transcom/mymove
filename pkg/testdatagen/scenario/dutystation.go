package scenario

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// RunDutyStationScenario3 creates some duty stations.
func RunDutyStationScenario3(db *pop.Connection) {
	// These are some random Marine bases found on the internet, not for production use.
	testdatagen.MakeDutyStation(db, "Air Station Yuma", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Yuma", State: "Arizona", PostalCode: "85364"})
	testdatagen.MakeDutyStation(db, "Air Station Miramar", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "San Diego", State: "California", PostalCode: "92145"})
	testdatagen.MakeDutyStation(db, "Air Station Cherry Point", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Cherry Point", State: "North Carolina", PostalCode: "28533"})
	testdatagen.MakeDutyStation(db, "Air Station New River", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Jacksonville", State: "North Carolina", PostalCode: "28540"})
	testdatagen.MakeDutyStation(db, "Air Station Beaufort", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Beaufort", State: "South Carolina", PostalCode: "29904"})
	testdatagen.MakeDutyStation(db, "Air Ground Combat Center Twentynine Palms", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Twentynine Palms", State: "California", PostalCode: "92278"})
	testdatagen.MakeDutyStation(db, "Base Camp Pendleton", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Oceanside", State: "California", PostalCode: "92058"})
	testdatagen.MakeDutyStation(db, "Recruit Depot San Diego", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "San Diego", State: "California", PostalCode: "92140"})
	testdatagen.MakeDutyStation(db, "Base Hawaii", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Honolulu", State: "Hawaii", PostalCode: "96734"})
	testdatagen.MakeDutyStation(db, "Base Camp Lejeune", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Lejeune", State: "North Carolina", PostalCode: "28547"})
	testdatagen.MakeDutyStation(db, "Recruit Depot Parris Island", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Parris Island", State: "South Carolina", PostalCode: "29905"})
	testdatagen.MakeDutyStation(db, "Base Quantico", internalmessages.AffiliationMARINES,
		models.Address{StreetAddress1: "duty station", City: "Quantico", State: "Virginia", PostalCode: "22134"})
}
