package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

const civilianBaseUBAllowanceTestConstant = 350
const dependents12AndOverUBAllowanceTestConstant = 350
const depedentsUnder12UBAllowanceTestConstant = 175
const maxWholeFamilyCivilianUBAllowanceTestConstant = 2000

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberGradeE1
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION

	suite.Run("E1 with dependents", func() {
		E1FullLoad := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(8000, E1FullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("E1 without dependents", func() {
		E1Solo := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(5000, E1Solo.TotalWeightSelf)
	})

	suite.Run("E1 Pro Gear", func() {
		E1ProGear := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(2000, E1ProGear.ProGearWeight)
	})

	suite.Run("E1 Pro Gear Spouse", func() {
		E1ProGearSpouse := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(500, E1ProGearSpouse.ProGearWeightSpouse)
	})
}

func (suite *ModelSuite) TestGetUBWeightAllowanceIsZero() {
	appCtx := suite.AppContextForTest()
	branch := models.AffiliationMARINES
	originDutyLocationIsOconus := false
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeE1
	orderType := internalmessages.OrdersTypeLOCALMOVE
	dependentsAuthorized := true
	isAccompaniedTour := true
	dependentsUnderTwelve := 2
	dependentsTwelveAndOver := 1

	suite.Run("UB allowance is zero when origin and new duty location are both CONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})

	suite.Run("UB allowance is zero for orders type OrdersTypeLOCALMOVE", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})

	originDutyLocationIsOconus = true
	orderType = internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	dependentsAuthorized = false
	suite.Run("UB allowance is zero for nonexistent combination of branch, grade, orders_type, dependents_authorized, and accompanied_tour", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetUBWeightAllowanceCivilians() {
	appCtx := suite.AppContextForTest()
	branch := models.AffiliationCOASTGUARD
	originDutyLocationIsOconus := true
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeCIVILIANEMPLOYEE
	orderType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	dependentsAuthorized := true
	isAccompaniedTour := true

	dependentsUnderTwelve := 0
	dependentsTwelveAndOver := 0
	suite.Run("UB allowance is calculated for Civilian Employee pay grade with no dependents", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(civilianBaseUBAllowanceTestConstant, ubAllowance)
	})

	dependentsUnderTwelve = 0
	dependentsTwelveAndOver = 2
	suite.Run("UB allowance is calculated for Civilian Employee pay grade when dependentsUnderTwelve is 0 and dependentsTwelveAndOver is > 0", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(1050, civilianPlusDependentsTotalBaggageAllowance)
		suite.Assertions.Equal(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
	})

	dependentsUnderTwelve = 3
	dependentsTwelveAndOver = 0
	suite.Run("UB allowance is calculated for Civilian Employee pay grade when dependentsUnderTwelve is > 0 and dependentsTwelveAndOver is 0", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(875, civilianPlusDependentsTotalBaggageAllowance)
		suite.Assertions.Equal(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
	})

	dependentsUnderTwelve = 2
	dependentsTwelveAndOver = 1
	suite.Run("UB allowance is calculated for Civilian Employee pay grade", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(1050, civilianPlusDependentsTotalBaggageAllowance)
		suite.Assertions.Equal(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
	})

	dependentsUnderTwelve = 3
	dependentsTwelveAndOver = 4
	// this combination of depdendents would tally up to 2275 pounds normally
	// however, we limit the max ub allowance for a family to 2000
	suite.Run("UB allowance is set to 2000 for the max weight for a family", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(2275, civilianPlusDependentsTotalBaggageAllowance)
		suite.Assertions.NotEqual(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
		suite.Assertions.Equal(maxWholeFamilyCivilianUBAllowanceTestConstant, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetUBWeightAllowanceEdgeCases() {

	appCtx := suite.AppContextForTest()
	branch := models.AffiliationAIRFORCE
	originDutyLocationIsOconus := true
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeE1
	orderType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	dependentsAuthorized := true
	isAccompaniedTour := true
	dependentsUnderTwelve := 0
	dependentsTwelveAndOver := 0

	suite.Run("Air Force gets a UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationSPACEFORCE
	suite.Run("Space Force gets the same UB allowance as Air Force", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationNAVY
	grade = models.ServiceMemberGradeE9
	suite.Run("Pay grade E9 gets a UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	grade = models.ServiceMemberGradeE9SPECIALSENIORENLISTED
	suite.Run("Pay grade E9 Special Senior Enlisted and pay grade E9 get the same UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetUBWeightAllowanceWithValidValues() {
	appCtx := suite.AppContextForTest()
	branch := models.AffiliationMARINES
	originDutyLocationIsOconus := true
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeE1
	orderType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	dependentsAuthorized := true
	isAccompaniedTour := true
	dependentsUnderTwelve := 2
	dependentsTwelveAndOver := 4

	suite.Run("UB allowance is calculated when origin duty location is OCONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	originDutyLocationIsOconus = false
	newDutyLocationIsOconus = true
	suite.Run("UB allowance is calculated when new duty location is OCONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	suite.Run("OCONUS, Marines, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	suite.Run("OCONUS, Marines, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationAIRFORCE
	suite.Run("OCONUS, Air Force, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeTEMPORARYDUTY
	dependentsAuthorized = false
	isAccompaniedTour = false
	suite.Run("OCONUS, Air Force, E1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 400 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(400, ubAllowance)
	})

	grade = models.ServiceMemberGradeW2
	retirementOrderType := internalmessages.OrdersTypeRETIREMENT
	orderType = internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	suite.Run("Orders type of Retirement returns same entitlement value as the PCS orders type in the database", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &retirementOrderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeTEMPORARYDUTY
	suite.Run("OCONUS, Air Force, W1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 600", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetEntitlementByOrdersTypeWithValidValues() {
	E1 := models.ServiceMemberGradeE1
	ordersType := internalmessages.OrdersTypeSTUDENTTRAVEL

	suite.Run("Student Travel with dependents", func() {
		STFullLoad := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(350, STFullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("Student Travel without dependents", func() {
		STSolo := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(350, STSolo.TotalWeightSelf)
	})

	suite.Run("Student Travel Pro Gear", func() {
		STProGear := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(0, STProGear.ProGearWeight)
	})

	suite.Run("Student Travel Pro Gear Spouse", func() {
		STProGearSpouse := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(0, STProGearSpouse.ProGearWeightSpouse)
	})
}
