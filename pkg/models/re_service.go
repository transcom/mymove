package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReServiceCode is the code of service
type ReServiceCode string

const (
	// ReServiceCodeCS Counseling Services
	ReServiceCodeCS ReServiceCode = "CS"
	// ReServiceCodeDBHF Dom. Haul Away Boat Factor
	ReServiceCodeDBHF ReServiceCode = "DBHF"
	// ReServiceCodeDBTF Dom. Tow Away Boat Factor
	ReServiceCodeDBTF ReServiceCode = "DBTF"
	// ReServiceCodeDCRT Dom. Crating
	ReServiceCodeDCRT ReServiceCode = "DCRT"
	// ReServiceCodeDCRTSA Dom. Crating - Standalone
	ReServiceCodeDCRTSA ReServiceCode = "DCRTSA"
	// ReServiceCodeDDASIT "Dom. Destination Add'l SIT"
	ReServiceCodeDDASIT ReServiceCode = "DDASIT"
	// ReServiceCodeDDDSIT Dom. Destination SIT Delivery
	ReServiceCodeDDDSIT ReServiceCode = "DDDSIT"
	// ReServiceCodeDDFSIT Dom. Destination 1st Day SIT
	ReServiceCodeDDFSIT ReServiceCode = "DDFSIT"
	// ReServiceCodeDDP Dom. Destination Price
	ReServiceCodeDDP ReServiceCode = "DDP"
	// ReServiceCodeDDSHUT Dom. Destination Shuttle Service
	ReServiceCodeDDSHUT ReServiceCode = "DDSHUT"
	// ReServiceCodeDLH Dom. Linehaul
	ReServiceCodeDLH ReServiceCode = "DLH"
	// ReServiceCodeDMHF Dom. Mobile Home Factor
	ReServiceCodeDMHF ReServiceCode = "DMHF"
	// ReServiceCodeDNPKF Dom. NTS Packing Factor
	ReServiceCodeDNPKF ReServiceCode = "DNPKF"
	// ReServiceCodeDOASIT "Dom. Origin Add'l SIT"
	ReServiceCodeDOASIT ReServiceCode = "DOASIT"
	// ReServiceCodeDOFSIT Dom. Origin 1st Day SIT
	ReServiceCodeDOFSIT ReServiceCode = "DOFSIT"
	// ReServiceCodeDOP Dom. Origin Price
	ReServiceCodeDOP ReServiceCode = "DOP"
	// ReServiceCodeDOPSIT Dom. Origin SIT Pickup
	ReServiceCodeDOPSIT ReServiceCode = "DOPSIT"
	// ReServiceCodeDOSHUT Dom. Origin Shuttle Service
	ReServiceCodeDOSHUT ReServiceCode = "DOSHUT"
	// ReServiceCodeDPK Dom. Packing
	ReServiceCodeDPK ReServiceCode = "DPK"
	// ReServiceCodeDSH Dom. Shorthaul
	ReServiceCodeDSH ReServiceCode = "DSH"
	// ReServiceCodeDUCRT Dom. Uncrating
	ReServiceCodeDUCRT ReServiceCode = "DUCRT"
	// ReServiceCodeDUPK Dom. Unpacking
	ReServiceCodeDUPK ReServiceCode = "DUPK"
	// ReServiceCodeFSC Fuel Surcharge
	ReServiceCodeFSC ReServiceCode = "FSC"
	// ReServiceCodeIBHF Int’l. Haul Away Boat Factor
	ReServiceCodeIBHF ReServiceCode = "IBHF"
	// ReServiceCodeIBTF Int’l. Tow Away Boat Factor
	ReServiceCodeIBTF ReServiceCode = "IBTF"
	// ReServiceCodeICOLH "Int'l. C->O Shipping & LH"
	ReServiceCodeICOLH ReServiceCode = "ICOLH"
	// ReServiceCodeICOUB "Int'l. C->O UB"
	ReServiceCodeICOUB ReServiceCode = "ICOUB"
	// ReServiceCodeICRT "Int'l. Crating"
	ReServiceCodeICRT ReServiceCode = "ICRT"
	// ReServiceCodeICRTSA "Int'l. Crating - Standalone"
	ReServiceCodeICRTSA ReServiceCode = "ICRTSA"
	// ReServiceCodeIDASIT "Int'l. Destination Add'l Day SIT"
	ReServiceCodeIDASIT ReServiceCode = "IDASIT"
	// ReServiceCodeIDDSIT "Int'l. Destination SIT Delivery"
	ReServiceCodeIDDSIT ReServiceCode = "IDDSIT"
	// ReServiceCodeIDFSIT "Int'l. Destination 1st Day SIT"
	ReServiceCodeIDFSIT ReServiceCode = "IDFSIT"
	// ReServiceCodeIDSHUT "Int'l. Destination Shuttle Service"
	ReServiceCodeIDSHUT ReServiceCode = "IDSHUT"
	// ReServiceCodeIHPK "Int'l. HHG Pack"
	ReServiceCodeIHPK ReServiceCode = "IHPK"
	// ReServiceCodeIHUPK "Int'l. HHG Unpack"
	ReServiceCodeIHUPK ReServiceCode = "IHUPK"
	// ReServiceCodeINPKF Int’l. NTS Packing Factor
	ReServiceCodeINPKF ReServiceCode = "INPKF"
	// ReServiceCodeIOASIT "Int'l. Origin Add'l Day SIT"
	ReServiceCodeIOASIT ReServiceCode = "IOASIT"
	// ReServiceCodeIOCLH "Int'l. O->C Shipping & LH"
	ReServiceCodeIOCLH ReServiceCode = "IOCLH"
	// ReServiceCodeIOCUB "Int'l. O->C UB"
	ReServiceCodeIOCUB ReServiceCode = "IOCUB"
	// ReServiceCodeIOFSIT "Int'l. Origin 1st Day SIT"
	ReServiceCodeIOFSIT ReServiceCode = "IOFSIT"
	// ReServiceCodeIOOLH "Int'l. O->O Shipping & LH"
	ReServiceCodeIOOLH ReServiceCode = "IOOLH"
	// ReServiceCodeIOOUB "Int'l. O->O UB"
	ReServiceCodeIOOUB ReServiceCode = "IOOUB"
	// ReServiceCodeIOPSIT "Int'l. Origin SIT Pickup"
	ReServiceCodeIOPSIT ReServiceCode = "IOPSIT"
	// ReServiceCodeIOSHUT "Int'l. Origin Shuttle Service"
	ReServiceCodeIOSHUT ReServiceCode = "IOSHUT"
	// ReServiceCodeIUBPK "Int'l. UB Pack"
	ReServiceCodeIUBPK ReServiceCode = "IUBPK"
	// ReServiceCodeIUBUPK "Int'l. UB Unpack"
	ReServiceCodeIUBUPK ReServiceCode = "IUBUPK"
	// ReServiceCodeIUCRT "Int'l. Uncrating"
	ReServiceCodeIUCRT ReServiceCode = "IUCRT"
	// ReServiceCodeMS Shipment Mgmt. Services
	ReServiceCodeMS ReServiceCode = "MS"
	// ReServiceCodeNSTH NonStd. HHG
	ReServiceCodeNSTH ReServiceCode = "NSTH"
	// ReServiceCodeNSTUB NonStd. UB
	ReServiceCodeNSTUB ReServiceCode = "NSTUB"
)

// ReService model struct
type ReService struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	Code      ReServiceCode `json:"code" db:"code"`
	Priority  int           `db:"priority"`
	Name      string        `json:"name" db:"name"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

// ReServices is not required by pop and may be deleted
type ReServices []ReService

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReService) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(r.Code), Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}
