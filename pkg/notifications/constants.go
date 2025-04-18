package notifications

import "github.com/transcom/mymove/pkg/models"

const OneSourceTransportationOfficeLink = "https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL"
const MyMoveLink = "https://my.move.mil/"
const WashingtonHQServicesLink = "https://www.esd.whs.mil"
const SmartVoucherLink = "https://smartvoucher.dfas.mil/"

const DTODFailureErrorMessage = "We are unable to calculate your distance. It may be that you have entered an invalid ZIP Code. Please check your ZIP Code to ensure it was entered correctly and is not a PO Box."
const DTODDownErrorMessage = "We are having an issue with the system we use to calculate mileage (DTOD) and cannot proceed."

var affiliationDisplayValue = map[models.ServiceMemberAffiliation]string{
	models.AffiliationARMY:       "Army",
	models.AffiliationNAVY:       "Marine Corps, Navy, and Coast Guard",
	models.AffiliationMARINES:    "Marine Corps, Navy, and Coast Guard",
	models.AffiliationAIRFORCE:   "Air Force and Space Force",
	models.AffiliationSPACEFORCE: "Air Force and Space Force",
	models.AffiliationCOASTGUARD: "Marine Corps, Navy, and Coast Guard",
}

func GetAffiliationDisplayValues() map[models.ServiceMemberAffiliation]string {
	return affiliationDisplayValue
}
