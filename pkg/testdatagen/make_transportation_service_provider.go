package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/random"
)

const alphanumericBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const defaultTspName = "Truss Transport LLC"
const defaultPocGeneralName = "Joey Joe-Joe Schabadoo"
const defaultPocGeneralEmail = "joey.j@example.com"
const defaultPocGeneralPhone = "(555) 101-0101"
const defaultPocClaimsName = "Art Vandelay"
const defaultPocClaimsEmail = "vandelay.ind@example.com"
const defaultPocClaimsPhone = "(555) 321-4321"
const defaultPayeeCode = "2708"

// RandomSCAC generates a random 4 figure string from allowed alphanumeric bytes to represent the SCAC.
func RandomSCAC() string {
	b := make([]byte, 4)
	for i := range b {
		randInt, err := random.GetRandomInt(len(alphanumericBytes))
		if err != nil {
			log.Panicf("failed to create random SCAC %v", err)
			return ""
		}
		b[i] = alphanumericBytes[randInt]
	}
	return string(b)
}

// DefaultSupplierID generates a default SupplierID for a given SCAC
func DefaultSupplierID(scac string) *string {
	var supplierID = scac + defaultPayeeCode
	return &supplierID
}

// MakeTSP makes a single transportation service provider record.
func MakeTSP(db *pop.Connection, assertions Assertions) models.TransportationServiceProvider {

	// Check to see if TSP has already been created
	tspID := assertions.TransportationServiceProvider.ID
	existingTsp := models.TransportationServiceProvider{}
	if !isZeroUUID(tspID) {
		if err := db.Find(&existingTsp, assertions.TransportationServiceProvider.ID); err == nil {
			// Found existing TSP for this ID
			return existingTsp
		}
	} else {
		tspID = uuid.Must(uuid.NewV4())
	}

	scac := assertions.TransportationServiceProvider.StandardCarrierAlphaCode
	if scac == "" {
		scac = "ABBV" //Valid SCAC for Syncada sandbox environment
	}

	supplierID := assertions.TransportationServiceProvider.SupplierID
	if supplierID == nil || *supplierID == "" {
		supplierID = DefaultSupplierID(scac)
	}

	name := assertions.TransportationServiceProvider.Name
	if name == nil {
		name = models.StringPointer(defaultTspName)
	}

	pocGeneralName := assertions.TransportationServiceProvider.PocGeneralName
	if pocGeneralName == nil {
		pocGeneralName = models.StringPointer(defaultPocGeneralName)
	}

	pocGeneralEmail := assertions.TransportationServiceProvider.PocGeneralEmail
	if pocGeneralEmail == nil {
		pocGeneralEmail = models.StringPointer(defaultPocGeneralEmail)
	}

	pocGeneralPhone := assertions.TransportationServiceProvider.PocGeneralPhone
	if pocGeneralPhone == nil {
		pocGeneralPhone = models.StringPointer(defaultPocGeneralPhone)
	}

	pocClaimsName := assertions.TransportationServiceProvider.PocClaimsName
	if pocClaimsName == nil {
		pocClaimsName = models.StringPointer(defaultPocClaimsName)
	}

	pocClaimsEmail := assertions.TransportationServiceProvider.PocClaimsEmail
	if pocClaimsEmail == nil {
		pocClaimsEmail = models.StringPointer(defaultPocClaimsEmail)
	}

	pocClaimsPhone := assertions.TransportationServiceProvider.PocClaimsPhone
	if pocClaimsPhone == nil {
		pocClaimsPhone = models.StringPointer(defaultPocClaimsPhone)
	}

	tsp := models.TransportationServiceProvider{
		ID:                       tspID,
		StandardCarrierAlphaCode: scac,
		SupplierID:               supplierID,
		Enrolled:                 assertions.TransportationServiceProvider.Enrolled,
		Name:                     name,
		PocGeneralName:           pocGeneralName,
		PocGeneralEmail:          pocGeneralEmail,
		PocGeneralPhone:          pocGeneralPhone,
		PocClaimsName:            pocClaimsName,
		PocClaimsEmail:           pocClaimsEmail,
		PocClaimsPhone:           pocClaimsPhone,
	}

	verrs, err := db.ValidateAndCreate(&tsp)
	if verrs.HasAny() {
		err = fmt.Errorf("TSP validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tsp
}

// MakeDefaultTSP makes a TSP with default values
func MakeDefaultTSP(db *pop.Connection) models.TransportationServiceProvider {
	return MakeTSP(db, Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			Enrolled: true,
		},
	})
}

// MakeTSPs creates numTSP number of TSP records
// numTSP specifies how many TSPs to create
func MakeTSPs(db *pop.Connection, numTSP int) {

	for i := 0; i < numTSP; i++ {
		MakeTSP(db, Assertions{
			TransportationServiceProvider: models.TransportationServiceProvider{
				StandardCarrierAlphaCode: RandomSCAC(),
			},
		})
	}
}
