package models_test

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPPMCloseoutSummaryValidation() {
	suite.Run("Test Valid PPM Closeout Summary", func() {
		validPPMCloseoutSummary := models.PPMCloseoutSummary{
			ID:            uuid.Must(uuid.NewV4()),
			PPMShipmentID: uuid.Must(uuid.NewV4()),
			MaxAdvance:    models.CentPointer(100000),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPPMCloseoutSummary, expErrors, nil)
	})

	suite.Run("Test missing PPMShipmentID", func() {
		invalidPPMCloseoutSummary := models.PPMCloseoutSummary{
			ID:         uuid.Must(uuid.NewV4()),
			MaxAdvance: models.CentPointer(100000),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		expErrors := map[string][]string{
			"ppmshipment_id": {"PPMShipmentID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPPMCloseoutSummary, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchPPMCloseoutByPPMID_AllColumnsPopulated() {
	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(suite.DB(), nil)
	now := time.Now()

	closeout := models.PPMCloseoutSummary{
		ID:                          uuid.Must(uuid.NewV4()),
		PPMShipmentID:               ppmShipment.ID,
		MaxAdvance:                  models.CentPointer(10000),
		GTCCPaidContractedExpense:   models.CentPointer(20000),
		MemberPaidContractedExpense: models.CentPointer(15000),
		GTCCPaidPackingMaterials:    models.CentPointer(5000),
		MemberPaidPackingMaterials:  models.CentPointer(4000),
		GTCCPaidWeighingFee:         models.CentPointer(3000),
		MemberPaidWeighingFee:       models.CentPointer(2500),
		GTCCPaidRentalEquipment:     models.CentPointer(10000),
		MemberPaidRentalEquipment:   models.CentPointer(9000),
		GTCCPaidTolls:               models.CentPointer(1000),
		MemberPaidTolls:             models.CentPointer(800),
		GTCCPaidOil:                 models.CentPointer(600),
		MemberPaidOil:               models.CentPointer(500),
		GTCCPaidOther:               models.CentPointer(700),
		MemberPaidOther:             models.CentPointer(600),
		TotalGTCCPaidExpenses:       models.CentPointer(50000),
		TotalMemberPaidExpenses:     models.CentPointer(40000),
		RemainingIncentive:          models.CentPointer(10000),
		GTCCPaidSIT:                 models.CentPointer(3000),
		MemberPaidSIT:               models.CentPointer(2500),
		GTCCPaidSmallPackage:        models.CentPointer(1200),
		MemberPaidSmallPackage:      models.CentPointer(1100),
		GTCCDisbursement:            models.CentPointer(2000),
		MemberDisbursement:          models.CentPointer(1900),
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}

	verrs, err := suite.DB().ValidateAndCreate(&closeout)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "expected no validation errors")

	fetched, err := models.FetchPPMCloseoutByPPMID(suite.DB(), ppmShipment.ID)
	suite.NoError(err, "expected no error when fetching an existing record")
	suite.NotNil(fetched)
}

func (suite *ModelSuite) TestFetchPPMCloseoutByPPMID_NotFound() {
	notFoundId := uuid.Must(uuid.NewV4())

	fetched, err := models.FetchPPMCloseoutByPPMID(suite.DB(), notFoundId)
	suite.Error(err)
	suite.Equal(fetched, models.PPMCloseoutSummary{})
	suite.True(errors.Is(err, models.ErrFetchNotFound), "Expected FETCH_NOT_FOUND error")
}
