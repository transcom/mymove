package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildUBAllowance creates a UB allowance
// Does not create other models
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUBAllowance(db *pop.Connection, customs []Customization, traits []Trait) models.UBAllowances {
	customs = setupCustomizations(customs, traits)

	// Find UBAllowances Customization and extract the custom UBAllowances
	var cUBAllowance models.UBAllowances
	if result := findValidCustomization(customs, UBAllowance); result != nil {
		cUBAllowance = result.Model.(models.UBAllowances)
		if result.LinkOnly {
			return cUBAllowance
		}
	}
	ubAllowanceValue := 2000
	branch := models.AffiliationAIRFORCE
	hasDependents := true
	accompaniedTour := true
	ubAllowance := models.UBAllowances{
		BranchOfService: (*string)(&branch),
		OrderPayGrade:   (*string)(models.ServiceMemberGradeE1.Pointer()),
		OrdersType:      (*string)(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION.Pointer()),
		HasDependents:   &hasDependents,
		AccompaniedTour: &accompaniedTour,
		UBAllowance:     &ubAllowanceValue,
	}

	// Overwrite default values with those from custom UB allowance
	testdatagen.MergeModels(&ubAllowance, cUBAllowance)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &ubAllowance)
	}

	return ubAllowance
}

// FetchOrBuildUBAllowance tries fetching a UBAllowance using a provided customization, then falls back to creating a default UBAllowance
func FetchOrBuildUBAllowance(db *pop.Connection, customs []Customization, traits []Trait) models.UBAllowances {
	if db == nil {
		return BuildUBAllowance(db, customs, traits)
	}

	customs = setupCustomizations(customs, traits)

	var cUBAllowance models.UBAllowances
	if result := findValidCustomization(customs, UBAllowance); result != nil {
		cUBAllowance = result.Model.(models.UBAllowances)
		if result.LinkOnly {
			return cUBAllowance
		}
	}

	var ubAllowance models.UBAllowances
	if !ubAllowance.ID.IsNil() {
		err := db.Where("ID = $1", ubAllowance.ID).First(&ubAllowance)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return ubAllowance
		}
	}

	if !ubAllowance.ID.IsNil() {
		err := db.Where("ID = $1", ubAllowance.ID).First(&ubAllowance)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return ubAllowance
		}
	}

	ubAllowanceValue := 2000
	branch := models.AffiliationAIRFORCE
	hasDependents := true
	accompaniedTour := true
	// search for the default UBAllowances if one is not provided
	defaultUBAllowance := models.UBAllowances{

		BranchOfService: (*string)(&branch),
		OrderPayGrade:   (*string)(models.ServiceMemberGradeE1.Pointer()),
		OrdersType:      (*string)(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION.Pointer()),
		HasDependents:   &hasDependents,
		AccompaniedTour: &accompaniedTour,
		UBAllowance:     &ubAllowanceValue,
	}
	err := db.Where("branch = ? AND grade = ? AND orders_type = ? AND dependents_authorized = ? AND accompanied_tour = ?", defaultUBAllowance.BranchOfService, defaultUBAllowance.OrderPayGrade, defaultUBAllowance.OrdersType, defaultUBAllowance.HasDependents, defaultUBAllowance.AccompaniedTour).First(&ubAllowance)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return ubAllowance
	}

	return BuildUBAllowance(db, customs, traits)
}
