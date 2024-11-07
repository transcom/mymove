package models_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const civilianBaseUBAllowanceTestConstant = 350
const dependents12AndOverUBAllowanceTestConstant = 350
const depedentsUnder12UBAllowanceTestConstant = 175
const maxWholeFamilyCivilianUBAllowanceTestConstant = 2000

type EntitlementsModelSuite struct {
	*testingsuite.PopTestSuite
}

func TestEntitlementsModelSuite(t *testing.T) {
	ts := &EntitlementsModelSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *EntitlementsModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberGradeE1

	suite.Run("E1 with dependents", func() {
		E1FullLoad := models.GetWeightAllotment(E1)
		suite.Assertions.Equal(8000, E1FullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("E1 without dependents", func() {
		E1Solo := models.GetWeightAllotment(E1)
		suite.Assertions.Equal(5000, E1Solo.TotalWeightSelf)
	})

	suite.Run("E1 Pro Gear", func() {
		E1ProGear := models.GetWeightAllotment(E1)
		suite.Assertions.Equal(2000, E1ProGear.ProGearWeight)
	})

	suite.Run("E1 Pro Gear Spouse", func() {
		E1ProGearSpouse := models.GetWeightAllotment(E1)
		suite.Assertions.Equal(500, E1ProGearSpouse.ProGearWeightSpouse)
	})
}

func (suite *EntitlementsModelSuite) TestGetUBWeightAllowanceIsZero() {
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

func (suite *EntitlementsModelSuite) TestGetUBWeightAllowanceCivilians() {
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

	dependentsUnderTwelve = 2
	dependentsTwelveAndOver = 1
	suite.Run("UB allowance is calculated for Civilian Employee pay grade", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
	})

	dependentsTwelveAndOver = 4
	// this combination of depdendents would tally up to 2100 pounds normally
	// however, we limit the max ub allowance for a family to 2000
	suite.Run("UB allowance is set to 2000 for the max weight for a family", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.NotEqual(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
		suite.Assertions.Equal(maxWholeFamilyCivilianUBAllowanceTestConstant, ubAllowance)
	})
}

func (suite *EntitlementsModelSuite) TestGetUBWeightAllowanceEdgeCases() {
	appCtx := suite.AppContextForTest()
	branch := models.AffiliationAIRFORCE
	originDutyLocationIsOconus := true
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeO4
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

func (suite *EntitlementsModelSuite) TestGetUBWeightAllowanceWithValidValues() {
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

	suite.Run("OCONUS, MARINES, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	suite.Run("OCONUS, MARINES, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
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

	isAccompaniedTour = false
	suite.Run("OCONUS, Air Force, E1, PCS, dependents are authorized, is NOT accompanied = 500 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(500, ubAllowance)
	})

	dependentsAuthorized = false
	suite.Run("OCONUS, Air Force, E1, PCS, dependents are NOTauthorized, is NOT accompanied = 500 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(500, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeTEMPORARYDUTY
	suite.Run("OCONUS, Air Force, E1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 400 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(400, ubAllowance)
	})

	grade = models.ServiceMemberGradeW1
	orderType = internalmessages.OrdersTypeRETIREMENT
	suite.Run("OCONUS, Air Force, W1, Retirement, dependents are NOT authorized, is NOT accompanied = 600 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})

	suite.Run("OCONUS, Air Force, W1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 600", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})

	branch = models.AffiliationSPACEFORCE
	suite.Run("Space Force And Air Force Get Same Entitlements: OCONUS, Space Force, W1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 600", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})
}
