package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

const civilianBaseUBAllowanceTestConstant = 350
const dependents12AndOverUBAllowanceTestConstant = 350
const depedentsUnder12UBAllowanceTestConstant = 175
const maxWholeFamilyCivilianUBAllowanceTestConstant = 2000
const studentTravelMaxAllowance = 350

func (suite *ModelSuite) TestAuthorizedWeightWhenExistsInDB() {
	aw := 3000
	entitlement := models.Entitlement{DBAuthorizedWeight: &aw}
	err := suite.DB().Create(&entitlement)
	suite.NoError(err)

	suite.Equal(entitlement.DBAuthorizedWeight, entitlement.AuthorizedWeight())
}

func (suite *ModelSuite) TestProGearAndProGearSpouseWeight() {
	suite.Run("no validation errors for ProGearWeight and ProGearSpouseWeight", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       2000,
			ProGearWeightSpouse: 500,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.False(verrs.HasAny(), "Should not have validation errors")
	})

	suite.Run("validation errors for ProGearWeight and ProGearSpouseWeight over max value", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       2001,
			ProGearWeightSpouse: 501,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("pro_gear_weight"))
		suite.NotNil(verrs.Get("pro_gear_weight_spouse"))
	})

	suite.Run("validation errors for ProGearWeight and ProGearSpouseWeight under min value", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       -1,
			ProGearWeightSpouse: -1,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("pro_gear_weight"))
		suite.NotNil(verrs.Get("pro_gear_weight_spouse"))
	})
}

func (suite *ModelSuite) TestOconusFields() {
	suite.Run("no validation errors for valid DependentsUnderTwelve, DependentsTwelveAndOver, and UBAllowance", func() {
		entitlement := models.Entitlement{
			ProGearWeight:           2000,
			ProGearWeightSpouse:     500,
			DependentsUnderTwelve:   models.IntPointer(1),
			DependentsTwelveAndOver: models.IntPointer(2),
			UBAllowance:             models.IntPointer(100),
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.False(verrs.HasAny())
	})

	suite.Run("validation errors for DependentsUnderTwelve and DependentsTwelveAndOver less than 0", func() {
		entitlement := models.Entitlement{
			DependentsTwelveAndOver: models.IntPointer(-1),
			DependentsUnderTwelve:   models.IntPointer(-1),
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("dependents_under_twelve"))
		suite.NotNil(verrs.Get("dependents_twelve_and_over"))
	})
}

func (suite *ModelSuite) TestTotalDependentsCalculation() {
	suite.Run("calculates total dependents correctly when both fields are set", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: models.IntPointer(3),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(5, *fetchedEntitlement.TotalDependents) // sum of 2 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsUnderTwelve is nil", func() {
		entitlement := models.Entitlement{
			DependentsTwelveAndOver: models.IntPointer(3),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(3, *fetchedEntitlement.TotalDependents) // sum of 0 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsTwelveAndOver is nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve: models.IntPointer(2),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(2, *fetchedEntitlement.TotalDependents) // sum of 2 + 0
	})
	suite.Run("sets total dependents to nil when both fields are nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   nil,
			DependentsTwelveAndOver: nil,
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.Nil(fetchedEntitlement.TotalDependents) // NOT 0, NOT A SUM, nil + nil is NULL
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
	civilianTDYUBAllowance := 0

	suite.Run("UB allowance is zero when origin and new duty location are both CONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})

	suite.Run("UB allowance is zero for orders type OrdersTypeLOCALMOVE", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})

	originDutyLocationIsOconus = true
	orderType = internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	dependentsAuthorized = false
	suite.Run("UB allowance is zero for nonexistent combination of branch, grade, orders_type, dependents_authorized, and accompanied_tour", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
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
	civilianTDYUBAllowance := 0

	dependentsUnderTwelve := 0
	dependentsTwelveAndOver := 0
	suite.Run("UB allowance is calculated for Civilian Employee pay grade with no dependents", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(civilianBaseUBAllowanceTestConstant, ubAllowance)
	})

	dependentsUnderTwelve = 2
	dependentsTwelveAndOver = 1
	suite.Run("UB allowance is calculated for Civilian Employee pay grade", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.Equal(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
	})

	dependentsTwelveAndOver = 4
	// this combination of depdendents would tally up to 2100 pounds normally
	// however, we limit the max ub allowance for a family to 2000
	suite.Run("UB allowance is set to 2000 for the max weight for a family", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		civilianPlusDependentsTotalBaggageAllowance := civilianBaseUBAllowanceTestConstant + (dependentsUnderTwelve * depedentsUnder12UBAllowanceTestConstant) + (dependentsTwelveAndOver * dependents12AndOverUBAllowanceTestConstant)
		suite.Assertions.NotEqual(civilianPlusDependentsTotalBaggageAllowance, ubAllowance)
		suite.Assertions.Equal(maxWholeFamilyCivilianUBAllowanceTestConstant, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetUBWeightAllowanceCivilianTDY() {
	appCtx := suite.AppContextForTest()
	branch := models.AffiliationCOASTGUARD
	originDutyLocationIsOconus := true
	newDutyLocationIsOconus := false
	grade := models.ServiceMemberGradeCIVILIANEMPLOYEE
	orderType := internalmessages.OrdersTypeTEMPORARYDUTY
	dependentsAuthorized := true
	isAccompaniedTour := true
	var civilianTDYUBAllowance int
	dependentsUnderTwelve := 0
	dependentsTwelveAndOver := 0

	suite.Run("UB allowance is 0 if no value is specified by the customer or office user for Civilian Employee with TDY orders type", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(0, ubAllowance)
	})

	civilianTDYUBAllowance = 350
	suite.Run("UB allowance is specified by the customer or office user for Civilian Employee with TDY orders type", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(350, ubAllowance)
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
	civilianTDYUBAllowance := 0

	suite.Run("Air Force gets a UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationSPACEFORCE
	suite.Run("Space Force gets the same UB allowance as Air Force", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationNAVY
	grade = models.ServiceMemberGradeE9
	suite.Run("Pay grade E9 gets a UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	grade = models.ServiceMemberGradeE9SPECIALSENIORENLISTED
	suite.Run("Pay grade E9 Special Senior Enlisted and pay grade E9 get the same UB allowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeSTUDENTTRAVEL
	suite.Run("orders type of STUDENT TRAVEL gets studentTravelMaxAllowance", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(studentTravelMaxAllowance, ubAllowance)
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
	civilianTDYUBAllowance := 0

	suite.Run("UB allowance is calculated when origin duty location is OCONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	originDutyLocationIsOconus = false
	newDutyLocationIsOconus = true
	suite.Run("UB allowance is calculated when new duty location is OCONUS", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	suite.Run("OCONUS, Marines, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	suite.Run("OCONUS, Marines, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	branch = models.AffiliationAIRFORCE
	suite.Run("OCONUS, Air Force, E1, PCS, dependents are authorized, is accompanied = 2000 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(2000, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeTEMPORARYDUTY
	dependentsAuthorized = false
	isAccompaniedTour = false
	suite.Run("OCONUS, Air Force, E1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 400 lbs", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(400, ubAllowance)
	})

	grade = models.ServiceMemberGradeW2
	retirementOrderType := internalmessages.OrdersTypeRETIREMENT
	orderType = internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	suite.Run("Orders type of Retirement returns same entitlement value as the PCS orders type in the database", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &retirementOrderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})

	orderType = internalmessages.OrdersTypeTEMPORARYDUTY
	suite.Run("OCONUS, Air Force, W1, Temporary Duty, dependents are NOT authorized, is NOT accompanied = 600", func() {
		ubAllowance, err := models.GetUBWeightAllowance(appCtx, &originDutyLocationIsOconus, &newDutyLocationIsOconus, &branch, &grade, &orderType, &dependentsAuthorized, &isAccompaniedTour, &dependentsUnderTwelve, &dependentsTwelveAndOver, &civilianTDYUBAllowance)
		suite.NoError(err)
		suite.Assertions.Equal(600, ubAllowance)
	})
}

func (suite *ModelSuite) TestGetMaxGunSafeAllowance() {
	suite.Run("returns the correct max gun safe allowance when parameter exists", func() {
		appCtx := suite.AppContextForTest()

		parameterName := "maxGunSafeAllowance"
		parameterValue := "500"

		param := models.ApplicationParameters{
			ParameterName:  &parameterName,
			ParameterValue: &parameterValue,
		}
		suite.MustSave(&param)

		val, err := models.GetMaxGunSafeAllowance(appCtx)
		suite.NoError(err)
		suite.Equal(500, val)
	})

	suite.Run("returns an error when parameter does not exist", func() {
		appCtx := suite.AppContextForTest()

		val, err := models.GetMaxGunSafeAllowance(appCtx)
		suite.Error(err)
		suite.Contains(err.Error(), "error fetching max gun safe allowance")
		suite.Equal(0, val)
	})
}
