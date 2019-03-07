package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
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
		b[i] = alphanumericBytes[rand.Intn(len(alphanumericBytes))]
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
	existingTsp := models.TransportationServiceProvider{}
	if !isZeroUUID(assertions.TransportationServiceProvider.ID) {
		if err := db.Find(&existingTsp, assertions.TransportationServiceProvider.ID); err == nil {
			// Found existing TSP for this ID
			return existingTsp
		}
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
		name = stringPointer(defaultTspName)
	}

	pocGeneralName := assertions.TransportationServiceProvider.PocGeneralName
	if pocGeneralName == nil {
		pocGeneralName = stringPointer(defaultPocGeneralName)
	}

	pocGeneralEmail := assertions.TransportationServiceProvider.PocGeneralEmail
	if pocGeneralEmail == nil {
		pocGeneralEmail = stringPointer(defaultPocGeneralEmail)
	}

	pocGeneralPhone := assertions.TransportationServiceProvider.PocGeneralPhone
	if pocGeneralPhone == nil {
		pocGeneralPhone = stringPointer(defaultPocGeneralPhone)
	}

	pocClaimsName := assertions.TransportationServiceProvider.PocClaimsName
	if pocClaimsName == nil {
		pocClaimsName = stringPointer(defaultPocClaimsName)
	}

	pocClaimsEmail := assertions.TransportationServiceProvider.PocClaimsEmail
	if pocClaimsEmail == nil {
		pocClaimsEmail = stringPointer(defaultPocClaimsEmail)
	}

	pocClaimsPhone := assertions.TransportationServiceProvider.PocClaimsPhone
	if pocClaimsPhone == nil {
		pocClaimsPhone = stringPointer(defaultPocClaimsPhone)
	}

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: scac,
		SupplierID:               supplierID,
		Enrolled:                 assertions.TransportationServiceProvider.Enrolled,
		Name:                     name,
		PocGeneralName:           pocGeneralName,
		PocGeneralEmail:          pocGeneralEmail,
		PocGeneralPhone:          pocGeneralPhone,
		PocClaimsName:            pocGeneralName,
		PocClaimsEmail:           pocGeneralEmail,
		PocClaimsPhone:           pocGeneralPhone,
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
