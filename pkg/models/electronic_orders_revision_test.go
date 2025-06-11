package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestElectronicOrdersRevisionValidate() {
	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerAirForce,
		OrdersNumber: "8675309",
		ID:           uuid.Must(uuid.NewV4()),
	}

	rev := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       models.ElectronicOrdersAffiliationAirForce,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	verrs, err := rev.Validate(nil)

	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestElectronicOrdersRevisionValidations() {
	empty := ""
	revision := &models.ElectronicOrdersRevision{
		SeqNum:                -1,
		MiddleName:            &empty,
		NameSuffix:            &empty,
		Title:                 &empty,
		LosingUIC:             &empty,
		LosingUnitName:        &empty,
		LosingUnitCity:        &empty,
		LosingUnitLocality:    &empty,
		LosingUnitCountry:     &empty,
		LosingUnitPostalCode:  &empty,
		GainingUIC:            &empty,
		GainingUnitName:       &empty,
		GainingUnitCity:       &empty,
		GainingUnitLocality:   &empty,
		GainingUnitCountry:    &empty,
		GainingUnitPostalCode: &empty,
		HhgTAC:                &empty,
		HhgSDN:                &empty,
		HhgLOA:                &empty,
		NtsTAC:                &empty,
		NtsSDN:                &empty,
		NtsLOA:                &empty,
		PovShipmentTAC:        &empty,
		PovShipmentSDN:        &empty,
		PovShipmentLOA:        &empty,
		PovStorageTAC:         &empty,
		PovStorageSDN:         &empty,
		PovStorageLOA:         &empty,
		UbTAC:                 &empty,
		UbSDN:                 &empty,
		UbLOA:                 &empty,
	}

	var expErrors = map[string][]string{
		"electronic_order_id":      {"ElectronicOrderID can not be blank."},
		"seq_num":                  {"-1 is not greater than -1."},
		"given_name":               {"GivenName can not be blank."},
		"middle_name":              {"MiddleName can not be blank."},
		"family_name":              {"FamilyName can not be blank."},
		"name_suffix":              {"NameSuffix can not be blank."},
		"paygrade":                 {"Paygrade is not in the list [aviation-cadet, cadet, civilian, E-1, E-2, E-3, E-4, E-5, E-6, E-7, E-8, E-9, midshipman, O-1, O-2, O-3, O-4, O-5, O-6, O-7, O-8, O-9, O-10, W-1, W-2, W-3, W-4, W-5]."},
		"affiliation":              {"Affiliation is not in the list [air-force, army, civilian-agency, coast-guard, marine-corps, navy]."},
		"title":                    {"Title can not be blank."},
		"date_issued":              {"DateIssued can not be blank."},
		"status":                   {"Status is not in the list [authorized, rfo, canceled]."},
		"tour_type":                {"TourType is not in the list [accompanied, unaccompanied, unaccompanied-dependents-restricted]."},
		"orders_type":              {"OrdersType is not in the list [accession, between-duty-stations, brac, cot, emergency-evac, ipcot, low-cost-travel, operational, oteip, rotational, separation, special-purpose, training, unit-move]."},
		"losing_uic":               {"LosingUIC can not be blank."},
		"losing_unit_name":         {"LosingUnitName can not be blank."},
		"losing_unit_city":         {"LosingUnitCity can not be blank."},
		"losing_unit_locality":     {"LosingUnitLocality can not be blank."},
		"losing_unit_postal_code":  {"LosingUnitPostalCode can not be blank."},
		"losing_unit_country":      {"LosingUnitCountry can not be blank."},
		"gaining_uic":              {"GainingUIC can not be blank."},
		"gaining_unit_name":        {"GainingUnitName can not be blank."},
		"gaining_unit_city":        {"GainingUnitCity can not be blank."},
		"gaining_unit_locality":    {"GainingUnitLocality can not be blank."},
		"gaining_unit_postal_code": {"GainingUnitPostalCode can not be blank."},
		"gaining_unit_country":     {"GainingUnitCountry can not be blank."},
		"hhg_tac":                  {"HhgTAC can not be blank."},
		"hhg_sdn":                  {"HhgSDN can not be blank."},
		"hhg_loa":                  {"HhgLOA can not be blank."},
		"nts_tac":                  {"NtsTAC can not be blank."},
		"nts_sdn":                  {"NtsSDN can not be blank."},
		"nts_loa":                  {"NtsLOA can not be blank."},
		"pov_shipment_tac":         {"PovShipmentTAC can not be blank."},
		"pov_shipment_sdn":         {"PovShipmentSDN can not be blank."},
		"pov_shipment_loa":         {"PovShipmentLOA can not be blank."},
		"pov_storage_tac":          {"PovStorageTAC can not be blank."},
		"pov_storage_sdn":          {"PovStorageSDN can not be blank."},
		"pov_storage_loa":          {"PovStorageLOA can not be blank."},
		"ub_tac":                   {"UbTAC can not be blank."},
		"ub_loa":                   {"UbLOA can not be blank."},
		"ub_sdn":                   {"UbSDN can not be blank."},
	}

	suite.verifyValidationErrors(revision, expErrors, nil)
}

func (suite *ModelSuite) TestCreateElectronicOrdersRevision() {
	order := &models.ElectronicOrder{
		ID:           uuid.Must(uuid.NewV4()),
		Edipi:        "1234567890",
		Issuer:       models.IssuerAirForce,
		OrdersNumber: "8675309",
	}
	verrs, err := order.Validate(nil)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	rev := &models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   *order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       models.ElectronicOrdersAffiliationAirForce,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	verrs, err = rev.Validate(nil)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestCreateElectronicOrdersRevision_Amendment() {
	order := &models.ElectronicOrder{
		ID:           uuid.Must(uuid.NewV4()),
		Edipi:        "1234567890",
		Issuer:       models.IssuerAirForce,
		OrdersNumber: "8675309",
	}

	rev0 := &models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   *order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       models.ElectronicOrdersAffiliationAirForce,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	rev1 := &models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   *order,
		SeqNum:            1,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       models.ElectronicOrdersAffiliationAirForce,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	verrs, err := suite.DB().ValidateAndCreate(order)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	verrs, err = suite.DB().ValidateAndCreate(rev0)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	verrs, err = suite.DB().ValidateAndCreate(rev1)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	retrievedOrder, err := models.FetchElectronicOrderByID(suite.DB(), order.ID)
	suite.NoError(err)
	suite.Len(retrievedOrder.Revisions, 2)
	expectedIDs := []uuid.UUID{rev0.ID, rev1.ID}
	suite.Contains(expectedIDs, retrievedOrder.Revisions[0].ID)
	suite.Contains(expectedIDs, retrievedOrder.Revisions[1].ID)
	suite.NotEqual(retrievedOrder.Revisions[0].ID, retrievedOrder.Revisions[1].ID)
}
