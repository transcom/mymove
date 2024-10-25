package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReServiceCode is the code of service
type ReServiceCode string

func (r ReServiceCode) String() string {
	return string(r)
}

const (
	// ReServiceCodeCS Counseling
	ReServiceCodeCS ReServiceCode = "CS"
	// ReServiceCodeDBHF Domestic haul away boat factor
	ReServiceCodeDBHF ReServiceCode = "DBHF"
	// ReServiceCodeDBTF Domestic tow away boat factor
	ReServiceCodeDBTF ReServiceCode = "DBTF"
	// ReServiceCodeDCRT Domestic crating
	ReServiceCodeDCRT ReServiceCode = "DCRT"
	// ReServiceCodeDCRTSA Domestic crating - standalone
	ReServiceCodeDCRTSA ReServiceCode = "DCRTSA"
	// ReServiceCodeDDASIT Domestic destination add'l SIT
	ReServiceCodeDDASIT ReServiceCode = "DDASIT"
	// ReServiceCodeDDDSIT Domestic destination SIT delivery
	ReServiceCodeDDDSIT ReServiceCode = "DDDSIT"
	// ReServiceCodeDDSFSC Domestic destination SIT FSC
	ReServiceCodeDDSFSC ReServiceCode = "DDSFSC"
	// ReServiceCodeDDFSIT Domestic destination 1st day SIT
	ReServiceCodeDDFSIT ReServiceCode = "DDFSIT"
	// ReServiceCodeDDP Domestic destination price
	ReServiceCodeDDP ReServiceCode = "DDP"
	// ReServiceCodeDDSHUT Domestic destination shuttle service
	ReServiceCodeDDSHUT ReServiceCode = "DDSHUT"
	// ReServiceCodeDLH Domestic linehaul
	ReServiceCodeDLH ReServiceCode = "DLH"
	// ReServiceCodeDMHF Domestic mobile home factor
	ReServiceCodeDMHF ReServiceCode = "DMHF"
	// ReServiceCodeDNPK Domestic NTS packing
	ReServiceCodeDNPK ReServiceCode = "DNPK"
	// ReServiceCodeDOASIT Domestic origin add'l SIT
	ReServiceCodeDOASIT ReServiceCode = "DOASIT"
	// ReServiceCodeDOFSIT Domestic origin 1st day SIT
	ReServiceCodeDOFSIT ReServiceCode = "DOFSIT"
	// ReServiceCodeDOP Domestic origin price
	ReServiceCodeDOP ReServiceCode = "DOP"
	// ReServiceCodeDOPSIT Domestic origin SIT pickup
	ReServiceCodeDOPSIT ReServiceCode = "DOPSIT"
	// ReServiceCodeDOSFSC Domestic origin SIT FSC
	ReServiceCodeDOSFSC ReServiceCode = "DOSFSC"
	// ReServiceCodeDOSHUT Domestic origin shuttle service
	ReServiceCodeDOSHUT ReServiceCode = "DOSHUT"
	// ReServiceCodeDPK Domestic packing
	ReServiceCodeDPK ReServiceCode = "DPK"
	// ReServiceCodeDSH Domestic shorthaul
	ReServiceCodeDSH ReServiceCode = "DSH"
	// ReServiceCodeDUCRT Domestic uncrating
	ReServiceCodeDUCRT ReServiceCode = "DUCRT"
	// ReServiceCodeDUPK Domestic unpacking
	ReServiceCodeDUPK ReServiceCode = "DUPK"
	// ReServiceCodeFSC Fuel Surcharge
	ReServiceCodeFSC ReServiceCode = "FSC"
	// ReServiceCodeIBHF International haul away boat factor
	ReServiceCodeIBHF ReServiceCode = "IBHF"
	// ReServiceCodeIBTF International tow away boat factor
	ReServiceCodeIBTF ReServiceCode = "IBTF"
	// ReServiceCodeICOLH International C->O shipping & LH
	ReServiceCodeICOLH ReServiceCode = "ICOLH"
	// ReServiceCodeICOUB International C->O UB
	ReServiceCodeICOUB ReServiceCode = "ICOUB"
	// ReServiceCodeICRT International crating
	ReServiceCodeICRT ReServiceCode = "ICRT"
	// ReServiceCodeICRTSA International crating - standalone
	ReServiceCodeICRTSA ReServiceCode = "ICRTSA"
	// ReServiceCodeIDASIT International destination add'l day SIT
	ReServiceCodeIDASIT ReServiceCode = "IDASIT"
	// ReServiceCodeIDDSIT International destination SIT delivery
	ReServiceCodeIDDSIT ReServiceCode = "IDDSIT"
	// ReServiceCodeIDFSIT International destination 1st day SIT
	ReServiceCodeIDFSIT ReServiceCode = "IDFSIT"
	// ReServiceCodeIDSHUT International destination shuttle service
	ReServiceCodeIDSHUT ReServiceCode = "IDSHUT"
	// ReServiceCodeIHPK International HHG pack
	ReServiceCodeIHPK ReServiceCode = "IHPK"
	// ReServiceCodeIHUPK International HHG unpack
	ReServiceCodeIHUPK ReServiceCode = "IHUPK"
	// ReServiceCodeINPK International NTS packing
	ReServiceCodeINPK ReServiceCode = "INPK"
	// ReServiceCodeIOASIT International origin add'l day SIT
	ReServiceCodeIOASIT ReServiceCode = "IOASIT"
	// ReServiceCodeIOCLH International O->C shipping & LH
	ReServiceCodeIOCLH ReServiceCode = "IOCLH"
	// ReServiceCodeIOCUB International O->C UB
	ReServiceCodeIOCUB ReServiceCode = "IOCUB"
	// ReServiceCodeIOFSIT International origin 1st day SIT
	ReServiceCodeIOFSIT ReServiceCode = "IOFSIT"
	// ReServiceCodeIOOLH International O->O shipping & LH
	ReServiceCodeIOOLH ReServiceCode = "IOOLH"
	// ReServiceCodeIOOUB International O->O UB
	ReServiceCodeIOOUB ReServiceCode = "IOOUB"
	// ReServiceCodeIOPSIT International origin SIT pickup
	ReServiceCodeIOPSIT ReServiceCode = "IOPSIT"
	// ReServiceCodeIOSHUT International origin shuttle service
	ReServiceCodeIOSHUT ReServiceCode = "IOSHUT"
	// ReServiceCodeIUBPK International UB pack
	ReServiceCodeIUBPK ReServiceCode = "IUBPK"
	// ReServiceCodeIUBUPK International UB unpack
	ReServiceCodeIUBUPK ReServiceCode = "IUBUPK"
	// ReServiceCodeIUCRT International uncrating
	ReServiceCodeIUCRT ReServiceCode = "IUCRT"
	// ReServiceCodeMS Move management
	ReServiceCodeMS ReServiceCode = "MS"
	// ReServiceCodeNSTH Nonstandard HHG
	ReServiceCodeNSTH ReServiceCode = "NSTH"
	// ReServiceCodeNSTUB Nonstandard UB
	ReServiceCodeNSTUB ReServiceCode = "NSTUB"
	// ReServiceCodeUBP International UB
	ReServiceCodeUBP ReServiceCode = "UBP"
	// ReServiceCodeISLH Shipping & Linehaul
	ReServiceCodeISLH ReServiceCode = "ISLH"
	// ReServiceCodePOEFSC International POE Fuel Surcharge
	ReServiceCodePOEFSC ReServiceCode = "POEFSC"
	// ReServiceCodePODFSC International POD Fuel Surcharge
	ReServiceCodePODFSC ReServiceCode = "PODFSC"
)

type ServiceLocationType string

const (
	// ServiceLocationO Origin
	ServiceLocationO ServiceLocationType = "O"
	// ServiceLocationD Destination
	ServiceLocationD ServiceLocationType = "D"
	// ServiceLocationB Both
	ServiceLocationB ServiceLocationType = "B"
)

// ReService model struct
type ReService struct {
	ID              uuid.UUID           `json:"id" db:"id"`
	Code            ReServiceCode       `json:"code" db:"code"`
	Priority        int                 `db:"priority"`
	Name            string              `json:"name" db:"name"`
	CreatedAt       time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
	ServiceLocation ServiceLocationType `json:"service_location" db:"service_location"`
}

// Hold groupings of SIT for the shipment
type SITServiceItemGroupings []SITServiceItemGrouping

// Holds the relevant SIT ReServiceCodes for Domestic Origin and Destination SIT
// service items, and provides a top-level summary due to our Service Item architecture
type SITServiceItemGrouping struct {
	Summary      SITSummary
	ServiceItems []MTOServiceItem
}

// Holds the summary of "Sub-Groupings" of SIT.
// For example, this will list the overall summary for an array of DOFSIT, DOPSIT, DOASIT, etc.,
// and the same for destination
type SITSummary struct {
	FirstDaySITServiceItemID uuid.UUID // TODO: Refactor this out and instead base payments off the entire grouping rather than just DOFSIT/DOASIT
	Location                 string
	DaysInSIT                int
	SITEntryDate             time.Time
	SITDepartureDate         *time.Time
	SITAuthorizedEndDate     time.Time
	SITCustomerContacted     *time.Time
	SITRequestedDelivery     *time.Time
}

// Definition of valid Domestic Origin SIT ReServiceCodes
var ValidDomesticOriginSITReServiceCodes = []ReServiceCode{
	ReServiceCodeDOASIT,
	ReServiceCodeDOFSIT,
	ReServiceCodeDOPSIT,
	ReServiceCodeDOSFSC,
	ReServiceCodeDOSHUT,
}

// Definition of valid Domestic Destination SIT ReServiceCodes
var ValidDomesticDestinationSITReServiceCodes = []ReServiceCode{
	ReServiceCodeDDASIT,
	ReServiceCodeDDDSIT,
	ReServiceCodeDDSFSC,
	ReServiceCodeDDFSIT,
	ReServiceCodeDDSHUT,
}

// Definition of valid International Origin SIT ReServiceCodes
var ValidInternationalOriginSITReServiceCodes = []ReServiceCode{
	ReServiceCodeIOASIT,
	ReServiceCodeIOFSIT,
	ReServiceCodeIOPSIT,
	ReServiceCodeIOSHUT,
}

// Definition of valid International Destination SIT ReServiceCodes
var ValidInternationalDestinationSITReServiceCodes = []ReServiceCode{
	ReServiceCodeIDASIT,
	ReServiceCodeIDDSIT,
	ReServiceCodeIDFSIT,
	ReServiceCodeIDSHUT,
}

// TableName overrides the table name used by Pop.
func (r ReService) TableName() string {
	return "re_services"
}

type ReServices []ReService

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReService) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(r.Code), Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}
