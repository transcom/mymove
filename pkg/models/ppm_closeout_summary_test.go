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

func (suite *ModelSuite) TestCalculateAndGetPPMCloseoutSummary() {
	suite.Run("CalculatePPMCloseoutSummary creates a closeout record", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)

		// Should not exist yet
		_, err := models.FetchPPMCloseoutByPPMID(suite.DB(), ppmShipment.ID)
		suite.Error(err)

		// Should succeed and create a record
		err = models.CalculatePPMCloseoutSummary(suite.DB(), ppmShipment.ID, false)
		suite.NoError(err)

		// Should succeed and find a record
		closeout, err := models.FetchPPMCloseoutByPPMID(suite.DB(), ppmShipment.ID)
		suite.NoError(err)
		suite.Equal(ppmShipment.ID, closeout.PPMShipmentID)
	})

	suite.Run("GetPPMCloseoutSummary returns the closeout record and calls calculate", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)

		// Should succeed and return the record
		closeout, err := models.GetPPMCloseoutSummary(suite.DB(), ppmShipment.ID, false)
		suite.NoError(err)
		suite.Equal(ppmShipment.ID, closeout.PPMShipmentID)
	})

	suite.Run("GetPPMCloseoutSummary & CalculatePPMCloseoutSummary don't recalculate if already calculated and recalculateIfExists is false", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		closeout, err := models.GetPPMCloseoutSummary(suite.DB(), ppmShipment.ID, false)
		suite.NoError(err)
		closeoutUpdatedAt := closeout.UpdatedAt

		closeout, err = models.GetPPMCloseoutSummary(suite.DB(), ppmShipment.ID, false)
		suite.NoError(err)
		suite.Equal(ppmShipment.ID, closeout.PPMShipmentID)
		suite.Equal(closeoutUpdatedAt, closeout.UpdatedAt, "Expected no change in updated_at when not recalculating")
	})

	suite.Run("GetPPMCloseoutSummary & CalculatePPMCloseoutSummary do recalculate if already calculated and recalculateIfExists is true", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		closeout, err := models.GetPPMCloseoutSummary(suite.DB(), ppmShipment.ID, false)
		suite.NoError(err)
		initialCloseoutMaxAdvance := closeout.MaxAdvance

		ppmShipment.EstimatedIncentive = initialCloseoutMaxAdvance
		err = suite.DB().Update(&ppmShipment)
		suite.NoError(err)

		closeout, err = models.GetPPMCloseoutSummary(suite.DB(), ppmShipment.ID, true)
		suite.NoError(err)
		suite.Equal(ppmShipment.ID, closeout.PPMShipmentID)
		suite.Equal(initialCloseoutMaxAdvance.Float64()*0.6, closeout.MaxAdvance.Float64(), "Expected max advance to change when recalculating")
	})

	suite.Run("GetPPMCloseoutSummary returns error for nonexistent PPM", func() {
		fakePPMID := uuid.Must(uuid.NewV4())
		_, err := models.GetPPMCloseoutSummary(suite.DB(), fakePPMID, false)
		suite.Error(err)
	})
}
