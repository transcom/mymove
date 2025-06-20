package ppmshipment

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/featureflag"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// estimatePPM Struct
type estimatePPM struct {
	checks               []ppmShipmentValidator
	planner              route.Planner
	paymentRequestHelper paymentrequesthelper.Helper
}

// NewEstimatePPM returns the estimatePPM (pass in checkRequiredFields() and checkEstimatedWeight)
func NewEstimatePPM(planner route.Planner, paymentRequestHelper paymentrequesthelper.Helper) services.PPMEstimator {
	return &estimatePPM{
		checks: []ppmShipmentValidator{
			checkRequiredFields(),
			checkEstimatedWeight(),
			checkSITRequiredFields(),
		},
		planner:              planner,
		paymentRequestHelper: paymentRequestHelper,
	}
}

func (f *estimatePPM) CalculatePPMSITEstimatedCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*unit.Cents, error) {
	if ppmShipment == nil {
		return nil, nil
	}

	oldPPMShipment, err := FindPPMShipment(appCtx, ppmShipment.ID)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment, err := mergePPMShipment(*ppmShipment, oldPPMShipment)
	if err != nil {
		return nil, err
	}

	err = validatePPMShipment(appCtx, *updatedPPMShipment, oldPPMShipment, &oldPPMShipment.Shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	contractDate := ppmShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	estimatedSITCost, err := CalculateSITCost(appCtx, updatedPPMShipment, contract)
	if err != nil {
		return nil, err
	}

	return estimatedSITCost, nil
}

func (f *estimatePPM) CalculatePPMSITEstimatedCostBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMSITEstimatedCostInfo, error) {

	if ppmShipment == nil {
		return nil, nil
	}

	oldPPMShipment, err := FindPPMShipment(appCtx, ppmShipment.ID)
	if err != nil {
		return nil, err
	}

	updatedPPMShipment, err := mergePPMShipment(*ppmShipment, oldPPMShipment)
	if err != nil {
		return nil, err
	}

	err = validatePPMShipment(appCtx, *updatedPPMShipment, oldPPMShipment, &oldPPMShipment.Shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	// Use actual departure date if possible
	contractDate := ppmShipment.ExpectedDepartureDate
	if ppmShipment.ActualMoveDate != nil {
		contractDate = *ppmShipment.ActualMoveDate
	}
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	ppmSITEstimatedCostInfoData, err := CalculateSITCostBreakdown(appCtx, updatedPPMShipment, contract)
	if err != nil {
		return nil, err
	}

	return ppmSITEstimatedCostInfoData, nil
}

func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, *unit.Cents, error) {
	return f.estimateIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) MaxIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, error) {
	return f.maxIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) FinalIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, error) {
	return f.finalIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) PriceBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, error) {
	return f.priceBreakdown(appCtx, ppmShipment)
}

func shouldSkipEstimatingIncentive(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment) bool {
	// check if GCC multipliers have changed or do not match
	if newPPMShipment.GCCMultiplierID != nil && oldPPMShipment.GCCMultiplierID != nil && *newPPMShipment.GCCMultiplierID != *oldPPMShipment.GCCMultiplierID {
		return false
	}
	if oldPPMShipment.Status != models.PPMShipmentStatusDraft && oldPPMShipment.EstimatedIncentive != nil && *newPPMShipment.EstimatedIncentive == 0 || oldPPMShipment.MaxIncentive == nil {
		return false
	} else {
		return oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate) &&
			newPPMShipment.PickupAddress.PostalCode == oldPPMShipment.PickupAddress.PostalCode &&
			newPPMShipment.DestinationAddress.PostalCode == oldPPMShipment.DestinationAddress.PostalCode &&
			((newPPMShipment.EstimatedWeight == nil && oldPPMShipment.EstimatedWeight == nil) || (oldPPMShipment.EstimatedWeight != nil && newPPMShipment.EstimatedWeight.Int() == oldPPMShipment.EstimatedWeight.Int()))
	}
}

func shouldSkipMaxIncentive(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment) bool {
	// check if GCC multipliers have changed or do not match
	if newPPMShipment.GCCMultiplierID != nil && oldPPMShipment.GCCMultiplierID != nil && *newPPMShipment.GCCMultiplierID != *oldPPMShipment.GCCMultiplierID {
		return false
	}

	// handle mismatches including nil and uuid.Nil
	newMultiplier := uuid.Nil
	if newPPMShipment.GCCMultiplierID != nil {
		newMultiplier = *newPPMShipment.GCCMultiplierID
	}

	oldMultiplier := uuid.Nil
	if oldPPMShipment.GCCMultiplierID != nil {
		oldMultiplier = *oldPPMShipment.GCCMultiplierID
	}

	if newMultiplier != oldMultiplier {
		return false
	}

	// if the max incentive is nil or 0, we want to update it
	if oldPPMShipment.MaxIncentive == nil || *oldPPMShipment.MaxIncentive == 0 {
		return false
	}

	// if the actual move date is being updated/added we want to re-run the max incentive
	if (oldPPMShipment.ActualMoveDate == nil && newPPMShipment.ActualMoveDate != nil) ||
		(oldPPMShipment.ActualMoveDate != nil && newPPMShipment.ActualMoveDate != nil && !newPPMShipment.ActualMoveDate.Equal(*oldPPMShipment.ActualMoveDate)) {
		return false
	} else {
		// if the departure date has changed, we want to recalculate
		return oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate)
	}
}

func shouldSkipCalculatingFinalIncentive(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment, originalTotalWeight unit.Pound, newTotalWeight unit.Pound) bool {
	// check if GCC multipliers have changed or do not match
	if newPPMShipment.GCCMultiplierID != nil && oldPPMShipment.GCCMultiplierID != nil && *newPPMShipment.GCCMultiplierID != *oldPPMShipment.GCCMultiplierID {
		return false
	}

	// handle mismatches including nil and uuid.Nil
	newMultiplier := uuid.Nil
	if newPPMShipment.GCCMultiplierID != nil {
		newMultiplier = *newPPMShipment.GCCMultiplierID
	}

	oldMultiplier := uuid.Nil
	if oldPPMShipment.GCCMultiplierID != nil {
		oldMultiplier = *oldPPMShipment.GCCMultiplierID
	}

	if newMultiplier != oldMultiplier {
		return false
	}
	// If oldPPMShipment field value is nil we know that the value has been updated and we should return false - the adjusted net weight is accounted for in the
	// SumWeights function and the change in weight is then checked with `newTotalWeight == originalTotalWeight`
	return (oldPPMShipment.ActualMoveDate != nil && newPPMShipment.ActualMoveDate.Equal(*oldPPMShipment.ActualMoveDate)) &&
		(oldPPMShipment.PickupAddress != nil && oldPPMShipment.DestinationAddress != nil && newPPMShipment.PickupAddress != nil && newPPMShipment.DestinationAddress != nil) &&
		(oldPPMShipment.PickupAddress.PostalCode != "" && newPPMShipment.PickupAddress.PostalCode == oldPPMShipment.PickupAddress.PostalCode) &&
		(oldPPMShipment.DestinationAddress.PostalCode != "" && newPPMShipment.DestinationAddress.PostalCode == oldPPMShipment.DestinationAddress.PostalCode) &&
		newTotalWeight == originalTotalWeight
}

func shouldSetFinalIncentiveToNil(newPPMShipment *models.PPMShipment, newTotalWeight unit.Pound) bool {
	if newPPMShipment.ActualMoveDate == nil ||
		newPPMShipment.PickupAddress == nil ||
		newPPMShipment.DestinationAddress == nil ||
		newPPMShipment.PickupAddress.PostalCode == "" ||
		newPPMShipment.DestinationAddress.PostalCode == "" ||
		newTotalWeight <= 0 {
		return true
	}

	return false
}

func shouldCalculateSITCost(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment) bool {
	// storage has not been selected yet or storage is not needed
	if newPPMShipment.SITExpected == nil || !*newPPMShipment.SITExpected {
		return false
	}

	// the service member could request storage but it can't be calculated until the services counselor provides info
	if newPPMShipment.SITLocation == nil {
		return false
	}

	// storage inputs weren't previously saved so we're calculating the cost for the first time
	if oldPPMShipment.SITLocation == nil ||
		oldPPMShipment.SITEstimatedWeight == nil ||
		oldPPMShipment.SITEstimatedEntryDate == nil ||
		oldPPMShipment.SITEstimatedDepartureDate == nil {
		return true
	}

	// the shipment had a previous storage cost but some of the inputs have changed, including some of the shipment
	// locations or departure date
	return *newPPMShipment.SITLocation != *oldPPMShipment.SITLocation ||
		*newPPMShipment.SITEstimatedWeight != *oldPPMShipment.SITEstimatedWeight ||
		*newPPMShipment.SITEstimatedEntryDate != *oldPPMShipment.SITEstimatedEntryDate ||
		*newPPMShipment.SITEstimatedDepartureDate != *oldPPMShipment.SITEstimatedDepartureDate ||
		newPPMShipment.PickupAddress.PostalCode != oldPPMShipment.PickupAddress.PostalCode ||
		newPPMShipment.DestinationAddress.PostalCode != oldPPMShipment.DestinationAddress.PostalCode ||
		newPPMShipment.ExpectedDepartureDate != oldPPMShipment.ExpectedDepartureDate
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*unit.Cents, *unit.Cents, error) {
	if newPPMShipment.Status != models.PPMShipmentStatusDraft && newPPMShipment.Status != models.PPMShipmentStatusSubmitted {
		return oldPPMShipment.EstimatedIncentive, oldPPMShipment.SITEstimatedCost, nil
	}
	// Check that all the required fields we need are present.
	err := validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	// If a field does not pass validation return nil as error handling is happening in the validator
	if err != nil {
		switch err.(type) {
		case apperror.InvalidInputError:
			return nil, nil, nil
		default:
			return nil, nil, err
		}
	}

	contractDate := newPPMShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, nil, err
	}

	calculateSITEstimate := shouldCalculateSITCost(newPPMShipment, &oldPPMShipment)

	// Clear out any previously calculated SIT estimated costs, if SIT is no longer expected
	if newPPMShipment.SITExpected != nil && !*newPPMShipment.SITExpected {
		newPPMShipment.SITEstimatedCost = nil
	}

	skipCalculatingEstimatedIncentive := shouldSkipEstimatingIncentive(newPPMShipment, &oldPPMShipment)

	if skipCalculatingEstimatedIncentive && !calculateSITEstimate {
		return oldPPMShipment.EstimatedIncentive, newPPMShipment.SITEstimatedCost, nil
	}

	estimatedIncentive := oldPPMShipment.EstimatedIncentive
	estimatedSITCost := oldPPMShipment.SITEstimatedCost

	// if the PPM is international, we will use a db func
	if newPPMShipment.Shipment.MarketCode != models.MarketCodeInternational {
		if !skipCalculatingEstimatedIncentive {
			// Clear out advance and advance requested fields when the estimated incentive is reset.
			newPPMShipment.HasRequestedAdvance = nil
			newPPMShipment.AdvanceAmountRequested = nil

			estimatedIncentive, err = f.calculatePrice(appCtx, newPPMShipment, 0, contract, false)
			if err != nil {
				return nil, nil, err
			}
		}

		if calculateSITEstimate {
			estimatedSITCost, err = CalculateSITCost(appCtx, newPPMShipment, contract)
			if err != nil {
				return nil, nil, err
			}
		}

		return estimatedIncentive, estimatedSITCost, nil

	} else {
		pickupAddress := newPPMShipment.PickupAddress
		destinationAddress := newPPMShipment.DestinationAddress

		if !skipCalculatingEstimatedIncentive {
			// Clear out advance and advance requested fields when the estimated incentive is reset.
			newPPMShipment.HasRequestedAdvance = nil
			newPPMShipment.AdvanceAmountRequested = nil

			estimatedIncentive, err = f.CalculateOCONUSIncentive(appCtx, newPPMShipment.ID, *pickupAddress, *destinationAddress, contractDate, newPPMShipment.EstimatedWeight.Int(), true, false, false)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to calculate estimated PPM incentive: %w", err)
			}
		}

		if calculateSITEstimate {
			var sitAddress models.Address
			isOrigin := *newPPMShipment.SITLocation == models.SITLocationTypeOrigin
			if isOrigin {
				sitAddress = *newPPMShipment.PickupAddress
			} else if !isOrigin {
				sitAddress = *newPPMShipment.DestinationAddress
			} else {
				return estimatedIncentive, estimatedSITCost, nil
			}
			daysInSIT := additionalDaysInSIT(*newPPMShipment.SITEstimatedEntryDate, *newPPMShipment.SITEstimatedDepartureDate)
			estimatedSITCost, err = f.CalculateOCONUSSITCosts(appCtx, newPPMShipment.ID, sitAddress.ID, isOrigin, contractDate, newPPMShipment.EstimatedWeight.Int(), daysInSIT)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to calculate estimated PPM incentive: %w", err)
			}
		}

		return estimatedIncentive, estimatedSITCost, nil
	}
}

func (f *estimatePPM) maxIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*unit.Cents, error) {
	// Check that all the required fields we need are present.
	err := validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	// If a field does not pass validation return nil as error handling is happening in the validator
	if err != nil {
		switch err.(type) {
		case apperror.InvalidInputError:
			return nil, nil
		default:
			return nil, err
		}
	}

	// we have access to the MoveTaskOrderID in the ppmShipment object so we can use that to get the customer's maximum weight entitlement
	var move models.Move
	err = appCtx.DB().Q().Eager(
		"Orders.Entitlement", "Orders.OriginDutyLocation.Address", "Orders.NewDutyLocation.Address",
	).Where("show = TRUE").Find(&move, newPPMShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " error querying move")
	}
	orders := move.Orders
	if orders.Entitlement.DBAuthorizedWeight == nil {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " DB authorized weight cannot be nil")
	}

	// max incentive should use the total of all of these weights
	dbAuthorizedWeight := 0
	proGearWeight := 0
	proGearWeightSpouse := 0
	gunSafe := 0

	// this has already been nil checked above
	dbAuthorizedWeight = *orders.Entitlement.DBAuthorizedWeight
	proGearWeight = orders.Entitlement.ProGearWeight
	// if dependents are authorized we can add the spouse pro gear
	if orders.Entitlement.DependentsAuthorized != nil || *orders.Entitlement.DependentsAuthorized {
		proGearWeightSpouse = orders.Entitlement.ProGearWeightSpouse
	}
	isGunSafeFeatureOn := false
	featureFlagName := "gun_safe"
	config := cli.GetFliptFetcherConfig(viper.GetViper())
	flagFetcher, err := featureflag.NewFeatureFlagFetcher(config)
	if err != nil {
		appCtx.Logger().Error("Error initializing FeatureFlagFetcher", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	}

	flag, err := flagFetcher.GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagName, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	} else {
		isGunSafeFeatureOn = flag.Match
	}
	// if the gun safe FF is on, we add the gun safe weight to the calculation here
	if isGunSafeFeatureOn {
		gunSafe = orders.Entitlement.GunSafeWeight
	}

	totalWeight := dbAuthorizedWeight + proGearWeight + proGearWeightSpouse + gunSafe

	contractDate := newPPMShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	maxIncentive := oldPPMShipment.MaxIncentive

	if newPPMShipment.Shipment.MarketCode != models.MarketCodeInternational {

		skipCalculatingMaxIncentive := shouldSkipMaxIncentive(newPPMShipment, &oldPPMShipment)

		if !skipCalculatingMaxIncentive {
			// since the max incentive is based off of the authorized weight entitlement and that value CAN change
			// we will calculate the max incentive each time it is called
			maxIncentive, err = f.calculatePrice(appCtx, newPPMShipment, unit.Pound(totalWeight), contract, true)
			if err != nil {
				return nil, err
			}
		}

		return maxIncentive, nil
	} else {
		pickupAddress := orders.OriginDutyLocation.Address
		destinationAddress := orders.NewDutyLocation.Address

		maxIncentive, err := f.CalculateOCONUSIncentive(appCtx, newPPMShipment.ID, pickupAddress, destinationAddress, contractDate, totalWeight, false, false, true)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate estimated PPM incentive: %w", err)
		}

		return maxIncentive, nil
	}
}

func (f *estimatePPM) finalIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*unit.Cents, error) {
	if newPPMShipment.Status != models.PPMShipmentStatusWaitingOnCustomer && newPPMShipment.Status != models.PPMShipmentStatusNeedsCloseout {
		return oldPPMShipment.FinalIncentive, nil
	}

	err := validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	if err != nil {
		switch err.(type) {
		case apperror.InvalidInputError:
			return nil, err
		default:
			return nil, err
		}
	}
	originalTotalWeight, newTotalWeight := SumWeights(oldPPMShipment, *newPPMShipment)

	if newPPMShipment.AllowableWeight != nil && *newPPMShipment.AllowableWeight < newTotalWeight {
		newTotalWeight = *newPPMShipment.AllowableWeight
	}

	// we have access to the MoveTaskOrderID in the ppmShipment object so we can use that to get the orders and allotment
	// This allows us to ensure the final incentive calculates with the actual weight, allowable weight, or total weight,
	// whichever is lowest.
	var move models.Move
	var entitlement models.WeightAllotment
	err = appCtx.DB().Q().Eager(
		"Orders.Entitlement",
	).Where("show = TRUE").Find(&move, newPPMShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " error querying move")
	}
	if move.Orders.Grade != nil {
		waf := entitlements.NewWeightAllotmentFetcher()
		entitlement, err = waf.GetWeightAllotment(appCtx, string(*move.Orders.Grade), move.Orders.OrdersType)
	} else {
		return nil, apperror.NewNotFoundError(move.ID, " orders.grade nil when getting weight allotment")
	}
	if err != nil {
		return nil, err
	}

	allotment := entitlement.TotalWeightSelfPlusDependents

	if newTotalWeight > unit.Pound(allotment) {
		newTotalWeight = unit.Pound(allotment)
	}

	contractDate := newPPMShipment.ExpectedDepartureDate
	if newPPMShipment.ActualMoveDate != nil {
		contractDate = *newPPMShipment.ActualMoveDate
	}
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	if newPPMShipment.Shipment.MarketCode != models.MarketCodeInternational {
		isMissingInfo := shouldSetFinalIncentiveToNil(newPPMShipment, newTotalWeight)
		var skipCalculateFinalIncentive bool
		finalIncentive := oldPPMShipment.FinalIncentive
		if !isMissingInfo {
			skipCalculateFinalIncentive = shouldSkipCalculatingFinalIncentive(newPPMShipment, &oldPPMShipment, originalTotalWeight, newTotalWeight)
			if !skipCalculateFinalIncentive {

				finalIncentive, err := f.calculatePrice(appCtx, newPPMShipment, newTotalWeight, contract, false)
				if err != nil {
					return nil, err
				}
				return capFinalIncentive(finalIncentive, newPPMShipment)
			}
		} else {
			finalIncentive = nil

			return finalIncentive, nil
		}
		return capFinalIncentive(finalIncentive, newPPMShipment)
	} else {
		pickupAddress := newPPMShipment.PickupAddress
		destinationAddress := newPPMShipment.DestinationAddress

		// we can't calculate actual incentive without the weight
		if newTotalWeight != 0 {
			finalIncentive, err := f.CalculateOCONUSIncentive(appCtx, newPPMShipment.ID, *pickupAddress, *destinationAddress, contractDate, newTotalWeight.Int(), false, true, false)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate estimated PPM incentive: %w", err)
			}
			return capFinalIncentive(finalIncentive, newPPMShipment)
		} else {
			return nil, nil
		}
	}
}

func capFinalIncentive(finalIncentive *unit.Cents, newPPMShipment *models.PPMShipment) (*unit.Cents, error) {
	if finalIncentive != nil && newPPMShipment.MaxIncentive != nil {
		if *finalIncentive > *newPPMShipment.MaxIncentive {
			finalIncentive = newPPMShipment.MaxIncentive
		}
		return finalIncentive, nil
	} else {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " MaxIncentive missing and/or finalIncentive nil when comparing")
	}

}

// SumWeights return the total weight of all weightTickets associated with a PPMShipment, returns 0 if there is no valid weight
func SumWeights(ppmShipment, newPPMShipment models.PPMShipment) (originalTotalWeight, newTotalWeight unit.Pound) {
	// small package PPMs will not have weight tickets, so we need to instead use moving expenses
	if newPPMShipment.PPMType != models.PPMTypeSmallPackage {
		if len(ppmShipment.WeightTickets) >= 1 {
			for _, weightTicket := range ppmShipment.WeightTickets {
				if weightTicket.Status != nil && *weightTicket.Status == models.PPMDocumentStatusRejected {
					originalTotalWeight += 0
				} else if weightTicket.AdjustedNetWeight != nil {
					originalTotalWeight += *weightTicket.AdjustedNetWeight
				} else if weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil {
					originalTotalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
				}
			}
		}

		if len(newPPMShipment.WeightTickets) >= 1 {
			for _, weightTicket := range newPPMShipment.WeightTickets {
				if weightTicket.Status != nil && *weightTicket.Status == models.PPMDocumentStatusRejected {
					newTotalWeight += 0
				} else if weightTicket.AdjustedNetWeight != nil {
					newTotalWeight += *weightTicket.AdjustedNetWeight
				} else if weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil {
					newTotalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
				}
			}
		}
	} else {
		for _, expense := range ppmShipment.MovingExpenses {
			if expense.Status != nil && *expense.Status == models.PPMDocumentStatusRejected {
				originalTotalWeight += 0
			} else if expense.WeightShipped != nil {
				originalTotalWeight += *expense.WeightShipped
			}
		}

		for _, expense := range newPPMShipment.MovingExpenses {
			if expense.Status != nil && *expense.Status == models.PPMDocumentStatusRejected {
				newTotalWeight += 0
			} else if expense.WeightShipped != nil {
				newTotalWeight += *expense.WeightShipped
			}
		}
	}

	return originalTotalWeight, newTotalWeight
}

// calculatePrice returns an incentive value for the ppm shipment as if we were pricing the service items for
// an HHG shipment with the same values for a payment request.  In this case we're not persisting service items,
// MTOServiceItems or PaymentRequestServiceItems, to the database to avoid unnecessary work and get a quicker result.
// we use this when calculating the estimated, final, and max incentive values
func (f estimatePPM) calculatePrice(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, totalWeight unit.Pound, contract models.ReContract, isMaxIncentiveCheck bool) (*unit.Cents, error) {
	logger := appCtx.Logger()

	zeroTotal := false
	serviceItemsToPrice := BaseServiceItems(*ppmShipment)

	var move models.Move
	err := appCtx.DB().Q().Eager(
		"Orders.Entitlement",
		"Orders.OriginDutyLocation.Address",
		"Orders.NewDutyLocation.Address",
	).Where("show = TRUE").Find(&move, ppmShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return nil, apperror.NewNotFoundError(ppmShipment.ID, " error querying move")
	}
	orders := move.Orders
	if orders.Entitlement.DBAuthorizedWeight == nil {
		return nil, apperror.NewNotFoundError(ppmShipment.ID, " DB authorized weight cannot be nil")
	}

	// Replace linehaul pricer with shorthaul pricer if move is within the same Zip3
	var pickupPostal, destPostal string

	gccMultiplier := ppmShipment.GCCMultiplier

	// if we are getting the max incentive, we want to use the addresses on the orders, else use what's on the shipment
	if isMaxIncentiveCheck {
		if orders.OriginDutyLocation.Address.PostalCode != "" {
			pickupPostal = orders.OriginDutyLocation.Address.PostalCode
		} else {
			return nil, apperror.NewNotFoundError(ppmShipment.ID, " No postal code for origin duty location on orders when comparing postal codes")
		}

		if orders.NewDutyLocation.Address.PostalCode != "" {
			destPostal = orders.NewDutyLocation.Address.PostalCode
		} else {
			return nil, apperror.NewNotFoundError(ppmShipment.ID, " No postal code for destination duty location on orders when comparing postal codes")
		}
	} else {
		if ppmShipment.PickupAddress != nil && ppmShipment.PickupAddress.PostalCode != "" {
			pickupPostal = ppmShipment.PickupAddress.PostalCode
		} else {
			return nil, apperror.NewNotFoundError(ppmShipment.ID, " no pickup address or zip on PPM - unable to calculate incentive")
		}

		if ppmShipment.DestinationAddress != nil && ppmShipment.DestinationAddress.PostalCode != "" {
			destPostal = ppmShipment.DestinationAddress.PostalCode
		} else {
			return nil, apperror.NewNotFoundError(ppmShipment.ID, " no destination address or zip on PPM - unable to calculate incentive")
		}
	}

	// if the ZIPs are the same, we need to replace the DLH service item with DSH
	if len(pickupPostal) >= 3 && len(destPostal) >= 3 && pickupPostal[:3] == destPostal[:3] {
		if pickupPostal[0:3] == destPostal[0:3] {
			for i, serviceItem := range serviceItemsToPrice {
				if serviceItem.ReService.Code == models.ReServiceCodeDLH {
					serviceItemsToPrice[i] = models.MTOServiceItem{ReService: models.ReService{Code: models.ReServiceCodeDSH}, MTOShipmentID: &ppmShipment.ShipmentID}
					break
				}
			}
		}
	}

	// Get a list of all the pricing params needed to calculate the price for each service item
	paramsForServiceItems, err := f.paymentRequestHelper.FetchServiceParamsForServiceItems(appCtx, serviceItemsToPrice)
	if err != nil {
		logger.Error("fetching PPM estimate ServiceParams failed", zap.Error(err))
		return nil, err
	}

	var mtoShipment models.MTOShipment
	if totalWeight > 0 && !isMaxIncentiveCheck {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		mtoShipment = MapPPMShipmentFinalFields(*ppmShipment, totalWeight)
	} else if totalWeight > 0 && isMaxIncentiveCheck {
		mtoShipment, err = MapPPMShipmentMaxIncentiveFields(appCtx, *ppmShipment, totalWeight)
		if err != nil {
			logger.Error("unable to map PPM fields for max incentive", zap.Error(err))
			return nil, err
		}
	} else {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		mtoShipment, err = MapPPMShipmentEstimatedFields(appCtx, *ppmShipment)
		if err != nil {
			logger.Error("unable to map PPM estimated fields", zap.Error(err))
			return nil, err
		}
	}

	totalPrice := unit.Cents(0)
	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return nil, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(appCtx, mtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code, contract.ID)

		// The distance value gets saved to the mto shipment model to reduce repeated api calls.
		var shipmentWithDistance models.MTOShipment
		err = appCtx.DB().Find(&shipmentWithDistance, mtoShipment.ID)
		if err != nil {
			logger.Error("could not find shipment in the database")
			return nil, err
		}
		serviceItem.MTOShipment = shipmentWithDistance
		// set this to avoid potential eTag errors because the MTOShipment.Distance field was likely updated
		ppmShipment.Shipment = shipmentWithDistance

		var paramValues models.PaymentServiceItemParams
		for _, param := range paramsForServiceCode(serviceItem.ReService.Code, paramsForServiceItems) {
			paramKey := param.ServiceItemParamKey
			// This is where the lookup() method of each service item param is actually evaluated
			paramValue, valueErr := keyData.ServiceParamValue(appCtx, paramKey.Key)
			if valueErr != nil {
				logger.Error("could not calculate param value lookup", zap.Error(valueErr))
				return nil, valueErr
			}

			// Gather all the param values for the service item to pass to the pricer's Price() method
			paymentServiceItemParam := models.PaymentServiceItemParam{
				// Some pricers like Fuel Surcharge try to requery the shipment through the service item, this is a
				// workaround to avoid a not found error because our PPM shipment has no service items saved in the db.
				// I think the FSC service item should really be relying on one of the zip distance params.
				PaymentServiceItem: models.PaymentServiceItem{
					MTOServiceItem: serviceItem,
				},
				ServiceItemParamKey: paramKey,
				Value:               paramValue,
			}
			paramValues = append(paramValues, paymentServiceItemParam)
		}

		if len(paramValues) == 0 {
			return nil, fmt.Errorf("no params were found for service item %s", serviceItem.ReService.Code)
		}

		centsValue, paymentParams, err := pricer.PriceUsingParams(appCtx, paramValues)
		// only apply the multiplier if centsValue is positive
		if gccMultiplier != nil && gccMultiplier.Multiplier > 0 && centsValue > 0 {
			oldCentsValue := centsValue
			multiplier := gccMultiplier.Multiplier
			multipliedPrice := float64(centsValue) * multiplier
			centsValue = unit.Cents(int(math.Round(multipliedPrice)))
			logger.Debug(fmt.Sprintf("Applying GCC multiplier: %f to service item price %s, original price: %d, new price: %d", multiplier, serviceItem.ReService.Code, oldCentsValue, centsValue))
		} else {
			logger.Debug(fmt.Sprintf("Service item price %s %d, no GCC multiplier applied (negative price or no multiplier)",
				serviceItem.ReService.Code, centsValue))
		}
		logger.Debug(fmt.Sprintf("Payment service item params %+v", paymentParams))

		if err != nil {
			if appCtx.Session().IsServiceMember() && ppmShipment.Shipment.Distance != nil && *ppmShipment.Shipment.Distance == unit.Miles(0) {
				zeroTotal = true
			} else {
				logger.Error("unable to calculate service item price", zap.Error(err))
				return nil, err
			}
		}

		totalPrice = totalPrice.AddCents(centsValue)
	}

	if zeroTotal {
		totalPrice = unit.Cents(0)
		return &totalPrice, nil
	}

	return &totalPrice, nil
}

// returns the final price breakdown of a ppm into linehaul, fuel, packing, unpacking, destination, and origin costs
func (f estimatePPM) priceBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, unit.Cents, error) {
	logger := appCtx.Logger()

	var emptyPrice unit.Cents
	var linehaul unit.Cents
	var fuel unit.Cents
	var origin unit.Cents
	var dest unit.Cents
	var packing unit.Cents
	var unpacking unit.Cents
	var storage unit.Cents

	serviceItemsToPrice := BaseServiceItems(*ppmShipment)

	// Replace linehaul pricer with shorthaul pricer if move is within the same Zip3
	var pickupPostal, destPostal string
	gccMultiplier := ppmShipment.GCCMultiplier

	// Check different address values for a postal code
	if ppmShipment.PickupAddress != nil && ppmShipment.PickupAddress.PostalCode != "" {
		pickupPostal = ppmShipment.PickupAddress.PostalCode
	} else {
		return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, apperror.NewNotFoundError(ppmShipment.ID, " no pickup address or zip on PPM - unable to calculate incentive")
	}

	// Same for destination
	if ppmShipment.DestinationAddress != nil && ppmShipment.DestinationAddress.PostalCode != "" {
		destPostal = ppmShipment.DestinationAddress.PostalCode
	} else {
		return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, apperror.NewNotFoundError(ppmShipment.ID, " no destination address or zip on PPM - unable to calculate incentive")
	}

	// if the ZIPs are the same, we need to replace the DLH service item with DSH
	if len(pickupPostal) >= 3 && len(destPostal) >= 3 && pickupPostal[:3] == destPostal[:3] {
		if pickupPostal[0:3] == destPostal[0:3] {
			for i, serviceItem := range serviceItemsToPrice {
				if serviceItem.ReService.Code == models.ReServiceCodeDLH {
					serviceItemsToPrice[i] = models.MTOServiceItem{ReService: models.ReService{Code: models.ReServiceCodeDSH}, MTOShipmentID: &ppmShipment.ShipmentID}
					break
				}
			}
		}
	}

	// Get a list of all the pricing params needed to calculate the price for each service item
	paramsForServiceItems, err := f.paymentRequestHelper.FetchServiceParamsForServiceItems(appCtx, serviceItemsToPrice)
	if err != nil {
		logger.Error("fetching PPM estimate ServiceParams failed", zap.Error(err))
		return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
	}

	contractDate := ppmShipment.ExpectedDepartureDate
	if ppmShipment.ActualMoveDate != nil {
		contractDate = *ppmShipment.ActualMoveDate
	}
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
	}

	var totalWeightFromWeightTicketsOrExpenses unit.Pound
	var blankPPM models.PPMShipment
	if ppmShipment.PPMType != models.PPMTypeSmallPackage {
		// for incentive-based/actual expense PPMs, weight tickets are required
		if ppmShipment.WeightTickets == nil {
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice,
				apperror.NewPPMNoWeightTicketsError(ppmShipment.ID, " no weight tickets")
		}
		_, totalWeightFromWeightTicketsOrExpenses = SumWeights(blankPPM, *ppmShipment)
	} else {
		// for small package PPM-SPRs, moving expenses are used
		if ppmShipment.MovingExpenses == nil {
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice,
				apperror.NewPPMNoMovingExpensesError(ppmShipment.ID, " no moving expenses")
		}
		_, totalWeightFromWeightTicketsOrExpenses = SumWeights(blankPPM, *ppmShipment)
	}

	var mtoShipment models.MTOShipment
	if totalWeightFromWeightTicketsOrExpenses > 0 {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		mtoShipment = MapPPMShipmentFinalFields(*ppmShipment, totalWeightFromWeightTicketsOrExpenses)
	} else {
		mtoShipment, err = MapPPMShipmentEstimatedFields(appCtx, *ppmShipment)
		if err != nil {
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
		}
	}

	doSITCalculation := *ppmShipment.SITExpected
	if doSITCalculation {
		estimatedSITCost, err := CalculateSITCost(appCtx, ppmShipment, contract)
		if err != nil {
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
		}

		if *estimatedSITCost > unit.Cents(0) {
			storage = *estimatedSITCost
		}
	}

	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(appCtx, mtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code, contract.ID)

		// The distance value gets saved to the mto shipment model to reduce repeated api calls.
		var shipmentWithDistance models.MTOShipment
		err = appCtx.DB().Find(&shipmentWithDistance, ppmShipment.ShipmentID)
		if err != nil {
			logger.Error("could not find shipment in the database")
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
		}
		serviceItem.MTOShipment = shipmentWithDistance
		// set this to avoid potential eTag errors because the MTOShipment.Distance field was likely updated
		ppmShipment.Shipment = shipmentWithDistance

		var paramValues models.PaymentServiceItemParams
		for _, param := range paramsForServiceCode(serviceItem.ReService.Code, paramsForServiceItems) {
			paramKey := param.ServiceItemParamKey
			// This is where the lookup() method of each service item param is actually evaluated
			paramValue, valueErr := keyData.ServiceParamValue(appCtx, paramKey.Key)
			if valueErr != nil {
				logger.Error("could not calculate param value lookup", zap.Error(valueErr))
				return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
			}

			// Gather all the param values for the service item to pass to the pricer's Price() method
			paymentServiceItemParam := models.PaymentServiceItemParam{
				// Some pricers like Fuel Surcharge try to requery the shipment through the service item, this is a
				// workaround to avoid a not found error because our PPM shipment has no service items saved in the db.
				// I think the FSC service item should really be relying on one of the zip distance params.
				PaymentServiceItem: models.PaymentServiceItem{
					MTOServiceItem: serviceItem,
				},
				ServiceItemParamKey: paramKey,
				Value:               paramValue,
			}
			paramValues = append(paramValues, paymentServiceItemParam)
		}

		if len(paramValues) == 0 {
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, fmt.Errorf("no params were found for service item %s", serviceItem.ReService.Code)
		}

		centsValue, _, err := pricer.PriceUsingParams(appCtx, paramValues)
		// only apply the multiplier if centsValue is positive
		if gccMultiplier != nil && gccMultiplier.Multiplier > 0 && centsValue > 0 {
			multiplier := gccMultiplier.Multiplier
			multipliedPrice := float64(centsValue) * multiplier
			centsValue = unit.Cents(int(math.Round(multipliedPrice)))
		}

		if err != nil {
			logger.Error("unable to calculate service item price", zap.Error(err))
			return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, err
		}

		switch serviceItem.ReService.Code {
		case models.ReServiceCodeDSH:
		case models.ReServiceCodeDLH:
			linehaul = centsValue
		case models.ReServiceCodeFSC:
			fuel = centsValue
		case models.ReServiceCodeDOP:
			origin = centsValue
		case models.ReServiceCodeDDP:
			dest = centsValue
		case models.ReServiceCodeDPK:
			packing = centsValue
		case models.ReServiceCodeDUPK:
			unpacking = centsValue
		}
	}

	return linehaul, fuel, origin, dest, packing, unpacking, storage, nil
}

// function for calculating incentives for OCONUS PPM shipments
// this uses a db function that takes in values needed to come up with the estimated/actual/max incentives
// this simulates the reimbursement for an iHHG move with ISLH, IHPK, IHUPK, and CONUS portion of FSC
func (f *estimatePPM) CalculateOCONUSIncentive(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID, pickupAddress models.Address, destinationAddress models.Address, moveDate time.Time, weight int, isEstimated bool, isActual bool, isMax bool) (*unit.Cents, error) {
	var mileage int
	ppmPort, err := models.FetchPortLocationByCode(appCtx.DB(), "4E1") // Tacoma, WA port
	if err != nil {
		return nil, fmt.Errorf("failed to fetch port location: %w", err)
	}

	// check if addresses are OCONUS or CONUS -> this determines how we check mileage to/from the authorized port
	isPickupOconus := pickupAddress.IsOconus != nil && *pickupAddress.IsOconus
	isDestinationOconus := destinationAddress.IsOconus != nil && *destinationAddress.IsOconus

	switch {
	case isPickupOconus && isDestinationOconus:
		// OCONUS -> OCONUS, we only reimburse for the CONUS mileage of the PPM
		mileage = 0
	case isPickupOconus && !isDestinationOconus:
		// OCONUS -> CONUS (port ZIP -> address ZIP)
		mileage, err = f.planner.ZipTransitDistance(appCtx, ppmPort.UsPostRegionCity.UsprZipID, destinationAddress.PostalCode)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate OCONUS to CONUS mileage: %w", err)
		}
	case !isPickupOconus && isDestinationOconus:
		// CONUS -> OCONUS (address ZIP -> port ZIP)
		mileage, err = f.planner.ZipTransitDistance(appCtx, pickupAddress.PostalCode, ppmPort.UsPostRegionCity.UsprZipID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate CONUS to OCONUS mileage: %w", err)
		}
	default:
		// covering down on CONUS -> CONUS moves - they should not appear here
		return nil, fmt.Errorf("invalid pickup and destination configuration: pickup isOconus=%v, destination isOconus=%v", isPickupOconus, isDestinationOconus)
	}

	incentive, err := models.CalculatePPMIncentive(appCtx.DB(), ppmShipmentID, pickupAddress.ID, destinationAddress.ID, moveDate, mileage, weight, isEstimated, isActual, isMax)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate PPM incentive: %w", err)
	}

	return (*unit.Cents)(&incentive.TotalIncentive), nil
}

func (f *estimatePPM) CalculateOCONUSSITCosts(appCtx appcontext.AppContext, ppmID uuid.UUID, addressID uuid.UUID, isOrigin bool, moveDate time.Time, weight int, sitDays int) (*unit.Cents, error) {
	if sitDays <= 0 {
		return nil, fmt.Errorf("SIT days must be greater than zero")
	}

	if weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than zero")
	}

	sitCosts, err := models.CalculatePPMSITCost(appCtx.DB(), ppmID, addressID, isOrigin, moveDate, weight, sitDays)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SIT costs: %w", err)
	}

	return (*unit.Cents)(&sitCosts.TotalSITCost), nil
}

func CalculateSITCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, contract models.ReContract) (*unit.Cents, error) {
	additionalDaysInSIT := additionalDaysInSIT(*ppmShipment.SITEstimatedEntryDate, *ppmShipment.SITEstimatedDepartureDate)

	if ppmShipment.Shipment.MarketCode != models.MarketCodeInternational {
		logger := appCtx.Logger()

		serviceItemsToPrice := StorageServiceItems(*ppmShipment, *ppmShipment.SITLocation, additionalDaysInSIT)

		totalPrice := unit.Cents(0)
		for _, serviceItem := range serviceItemsToPrice {
			pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
			if err != nil {
				logger.Error("unable to find pricer for service item", zap.Error(err))
				return nil, err
			}

			var price *unit.Cents
			switch serviceItemPricer := pricer.(type) {
			case services.DomesticOriginFirstDaySITPricer, services.DomesticDestinationFirstDaySITPricer:
				price, _, err = priceFirstDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract)
			case services.DomesticOriginAdditionalDaysSITPricer, services.DomesticDestinationAdditionalDaysSITPricer:
				price, _, err = priceAdditionalDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, additionalDaysInSIT, contract)
			default:
				return nil, fmt.Errorf("unknown SIT pricer type found for service item code %s", serviceItem.ReService.Code)
			}

			if err != nil {
				return nil, err
			}

			logger.Debug(fmt.Sprintf("Price of service item %s %d", serviceItem.ReService.Code, *price))
			totalPrice += *price
		}

		return &totalPrice, nil
	} else {
		var sitAddress models.Address
		isOrigin := *ppmShipment.SITLocation == models.SITLocationTypeOrigin
		if isOrigin {
			sitAddress = *ppmShipment.PickupAddress
		} else {
			sitAddress = *ppmShipment.DestinationAddress
		}

		contractDate := ppmShipment.ExpectedDepartureDate
		totalSITCost, err := models.CalculatePPMSITCost(appCtx.DB(), ppmShipment.ID, sitAddress.ID, isOrigin, contractDate, ppmShipment.SITEstimatedWeight.Int(), additionalDaysInSIT)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate PPM SIT incentive: %w", err)
		}
		return (*unit.Cents)(&totalSITCost.TotalSITCost), nil
	}
}

func CalculateSITCostBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, contract models.ReContract) (*models.PPMSITEstimatedCostInfo, error) {
	logger := appCtx.Logger()

	ppmSITEstimatedCostInfoData := &models.PPMSITEstimatedCostInfo{}

	additionalDaysInSIT := additionalDaysInSIT(*ppmShipment.SITEstimatedEntryDate, *ppmShipment.SITEstimatedDepartureDate)

	serviceItemsToPrice := StorageServiceItems(*ppmShipment, *ppmShipment.SITLocation, additionalDaysInSIT)

	totalPrice := unit.Cents(0)
	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return nil, err
		}

		var price *unit.Cents
		switch serviceItemPricer := pricer.(type) {
		case services.DomesticOriginFirstDaySITPricer, services.DomesticDestinationFirstDaySITPricer:
			price, ppmSITEstimatedCostInfoData, err = calculateFirstDaySITCostBreakdown(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract, ppmSITEstimatedCostInfoData, logger)
		case services.IntlOriginFirstDaySITPricer, services.IntlDestinationFirstDaySITPricer:
			price, ppmSITEstimatedCostInfoData, err = calculateIntlFirstDaySITCostBreakdown(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract, ppmSITEstimatedCostInfoData, logger)
		case services.DomesticOriginAdditionalDaysSITPricer, services.DomesticDestinationAdditionalDaysSITPricer:
			price, ppmSITEstimatedCostInfoData, err = calculateAdditionalDaySITCostBreakdown(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract, additionalDaysInSIT, ppmSITEstimatedCostInfoData, logger)
		case services.IntlOriginAdditionalDaySITPricer, services.IntlDestinationAdditionalDaySITPricer:
			price, ppmSITEstimatedCostInfoData, err = calculateIntlAdditionalDaySITCostBreakdown(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract, additionalDaysInSIT, ppmSITEstimatedCostInfoData, logger)
		default:
			return nil, fmt.Errorf("unknown SIT pricer type found for service item code %s", serviceItem.ReService.Code)
		}

		if err != nil {
			return nil, err
		}

		logger.Debug(fmt.Sprintf("Price of service item %s %d", serviceItem.ReService.Code, price))
		totalPrice += *price
	}

	ppmSITEstimatedCostInfoData.EstimatedSITCost = &totalPrice
	return ppmSITEstimatedCostInfoData, nil
}

func calculateFirstDaySITCostBreakdown(appCtx appcontext.AppContext, serviceItemPricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract, ppmSITEstimatedCostInfoData *models.PPMSITEstimatedCostInfo, logger *zap.Logger) (*unit.Cents, *models.PPMSITEstimatedCostInfo, error) {
	price, priceParams, err := priceFirstDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract)
	if err != nil {
		return nil, nil, err
	}
	ppmSITEstimatedCostInfoData.PriceFirstDaySIT = price
	for _, param := range priceParams {
		switch param.Key {
		case models.ServiceItemParamNameServiceAreaOrigin:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ServiceAreaOrigin = param.Value
		case models.ServiceItemParamNameServiceAreaDest:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ServiceAreaDestination = param.Value
		case models.ServiceItemParamNameIsPeak:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.IsPeak = param.Value
		case models.ServiceItemParamNameContractYearName:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ContractYearName = param.Value
		case models.ServiceItemParamNamePriceRateOrFactor:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.PriceRateOrFactor = param.Value
		case models.ServiceItemParamNameEscalationCompounded:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.EscalationCompounded = param.Value
		default:
			logger.Debug(fmt.Sprintf("Unexpected ServiceItemParam in PPM First Day SIT: %s, %s", param.Key, param.Value))
		}
	}
	return price, ppmSITEstimatedCostInfoData, nil
}

func calculateIntlFirstDaySITCostBreakdown(appCtx appcontext.AppContext, serviceItemPricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract, ppmSITEstimatedCostInfoData *models.PPMSITEstimatedCostInfo, logger *zap.Logger) (*unit.Cents, *models.PPMSITEstimatedCostInfo, error) {
	price, priceParams, err := priceFirstDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract)
	if err != nil {
		return nil, nil, err
	}
	ppmSITEstimatedCostInfoData.PriceFirstDaySIT = price
	for _, param := range priceParams {
		switch param.Key {
		case models.ServiceItemParamNameServiceAreaOrigin:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ServiceAreaOrigin = param.Value
		case models.ServiceItemParamNameServiceAreaDest:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ServiceAreaDestination = param.Value
		case models.ServiceItemParamNameIsPeak:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.IsPeak = param.Value
		case models.ServiceItemParamNameContractYearName:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.ContractYearName = param.Value
		case models.ServiceItemParamNamePriceRateOrFactor:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.PriceRateOrFactor = param.Value
		case models.ServiceItemParamNameEscalationCompounded:
			ppmSITEstimatedCostInfoData.ParamsFirstDaySIT.EscalationCompounded = param.Value
		default:
			logger.Debug(fmt.Sprintf("Unexpected ServiceItemParam in PPM First Day SIT: %s, %s", param.Key, param.Value))
		}
	}
	return price, ppmSITEstimatedCostInfoData, nil
}

func calculateAdditionalDaySITCostBreakdown(appCtx appcontext.AppContext, serviceItemPricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract, additionalDaysInSIT int, ppmSITEstimatedCostInfoData *models.PPMSITEstimatedCostInfo, logger *zap.Logger) (*unit.Cents, *models.PPMSITEstimatedCostInfo, error) {
	price, priceParams, err := priceAdditionalDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, additionalDaysInSIT, contract)
	if err != nil {
		return nil, nil, err
	}
	ppmSITEstimatedCostInfoData.PriceAdditionalDaySIT = price
	for _, param := range priceParams {
		switch param.Key {
		case models.ServiceItemParamNameServiceAreaOrigin:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ServiceAreaOrigin = param.Value
		case models.ServiceItemParamNameServiceAreaDest:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ServiceAreaDestination = param.Value
		case models.ServiceItemParamNameIsPeak:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.IsPeak = param.Value
		case models.ServiceItemParamNameContractYearName:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ContractYearName = param.Value
		case models.ServiceItemParamNamePriceRateOrFactor:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.PriceRateOrFactor = param.Value
		case models.ServiceItemParamNameEscalationCompounded:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.EscalationCompounded = param.Value
		case models.ServiceItemParamNameNumberDaysSIT:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.NumberDaysSIT = param.Value
		default:
			logger.Debug(fmt.Sprintf("Unexpected ServiceItemParam in PPM Additional Day SIT: %s, %s", param.Key, param.Value))
		}
	}
	return price, ppmSITEstimatedCostInfoData, nil
}

func calculateIntlAdditionalDaySITCostBreakdown(appCtx appcontext.AppContext, serviceItemPricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract, additionalDaysInSIT int, ppmSITEstimatedCostInfoData *models.PPMSITEstimatedCostInfo, logger *zap.Logger) (*unit.Cents, *models.PPMSITEstimatedCostInfo, error) {
	price, priceParams, err := priceAdditionalDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, additionalDaysInSIT, contract)
	if err != nil {
		return nil, nil, err
	}
	ppmSITEstimatedCostInfoData.PriceAdditionalDaySIT = price
	for _, param := range priceParams {
		switch param.Key {
		case models.ServiceItemParamNameServiceAreaOrigin:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ServiceAreaOrigin = param.Value
		case models.ServiceItemParamNameServiceAreaDest:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ServiceAreaDestination = param.Value
		case models.ServiceItemParamNameIsPeak:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.IsPeak = param.Value
		case models.ServiceItemParamNameContractYearName:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.ContractYearName = param.Value
		case models.ServiceItemParamNamePriceRateOrFactor:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.PriceRateOrFactor = param.Value
		case models.ServiceItemParamNameEscalationCompounded:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.EscalationCompounded = param.Value
		case models.ServiceItemParamNameNumberDaysSIT:
			ppmSITEstimatedCostInfoData.ParamsAdditionalDaySIT.NumberDaysSIT = param.Value
		default:
			logger.Debug(fmt.Sprintf("Unexpected ServiceItemParam in PPM Additional Day SIT: %s, %s", param.Key, param.Value))
		}
	}
	return price, ppmSITEstimatedCostInfoData, nil
}

func priceFirstDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract) (*unit.Cents, services.PricingDisplayParams, error) {
	if serviceItem.ReService.Code == models.ReServiceCodeIOFSIT || serviceItem.ReService.Code == models.ReServiceCodeIDFSIT {
		var addressID uuid.UUID
		if serviceItem.ReService.Code == models.ReServiceCodeIOFSIT {
			addressID = *ppmShipment.PickupAddressID
		} else {
			addressID = *ppmShipment.DestinationAddressID
		}
		reServiceID, _ := models.FetchReServiceByCode(appCtx.DB(), serviceItem.ReService.Code)
		intlOtherPrice, _ := models.FetchReIntlOtherPrice(appCtx.DB(), addressID, reServiceID.ID, contract.ID, &ppmShipment.ExpectedDepartureDate)
		firstDayPricer, ok := pricer.(services.IntlOriginFirstDaySITPricer)
		if !ok {
			return nil, nil, errors.New("ppm estimate pricer for SIT service item does not implement the first day pricer interface")
		}
		if ppmShipment.ActualMoveDate != nil {
			price, pricingParams, err := firstDayPricer.Price(appCtx, contract.Code, *ppmShipment.ActualMoveDate, *ppmShipment.SITEstimatedWeight, intlOtherPrice.PerUnitCents.Int())
			if err != nil {
				return nil, nil, err
			}

			appCtx.Logger().Debug(fmt.Sprintf("Pricing params for first day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

			return &price, pricingParams, nil
		}

		price, pricingParams, err := firstDayPricer.Price(appCtx, contract.Code, ppmShipment.ExpectedDepartureDate, *ppmShipment.SITEstimatedWeight, intlOtherPrice.PerUnitCents.Int())
		if err != nil {
			return nil, nil, err
		}

		appCtx.Logger().Debug(fmt.Sprintf("Pricing params for first day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

		return &price, pricingParams, nil
	} else {
		firstDayPricer, ok := pricer.(services.DomesticFirstDaySITPricer)
		if !ok {
			return nil, nil, errors.New("ppm estimate pricer for SIT service item does not implement the first day pricer interface")
		}

		// Need to declare if origin or destination for the serviceAreaLookup, otherwise we already have it
		serviceAreaPostalCode := ppmShipment.PickupAddress.PostalCode
		serviceAreaKey := models.ServiceItemParamNameServiceAreaOrigin
		if serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
			serviceAreaPostalCode = ppmShipment.DestinationAddress.PostalCode
			serviceAreaKey = models.ServiceItemParamNameServiceAreaDest
		}

		serviceAreaLookup := serviceparamvaluelookups.ServiceAreaLookup{
			Address: models.Address{PostalCode: serviceAreaPostalCode},
		}
		serviceArea, err := serviceAreaLookup.ParamValue(appCtx, contract.Code)
		if err != nil {
			return nil, nil, err
		}

		serviceAreaParam := services.PricingDisplayParam{
			Key:   serviceAreaKey,
			Value: serviceArea,
		}

		// Since this function may be ran before closeout, we need to account for if there's no actual move date yet.
		if ppmShipment.ActualMoveDate != nil {
			price, pricingParams, err := firstDayPricer.Price(appCtx, contract.Code, *ppmShipment.ActualMoveDate, *ppmShipment.SITEstimatedWeight, serviceArea, true)
			if err != nil {
				return nil, nil, err
			}

			pricingParams = append(pricingParams, serviceAreaParam)

			appCtx.Logger().Debug(fmt.Sprintf("Pricing params for first day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

			return &price, pricingParams, nil
		}
		price, pricingParams, err := firstDayPricer.Price(appCtx, contract.Code, ppmShipment.ExpectedDepartureDate, *ppmShipment.SITEstimatedWeight, serviceArea, true)
		if err != nil {
			return nil, nil, err
		}

		pricingParams = append(pricingParams, serviceAreaParam)

		appCtx.Logger().Debug(fmt.Sprintf("Pricing params for first day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

		return &price, pricingParams, nil
	}
}

func additionalDaysInSIT(sitEntryDate time.Time, sitDepartureDate time.Time) int {
	// The cost of first day SIT is already accounted for in DOFSIT or DDFSIT service items
	if sitEntryDate.Equal(sitDepartureDate) {
		return 0
	}

	difference := sitDepartureDate.Sub(sitEntryDate)
	return int(difference.Hours() / 24)
}

func priceAdditionalDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, additionalDaysInSIT int, contract models.ReContract) (*unit.Cents, services.PricingDisplayParams, error) {
	// international shipment logic
	if serviceItem.ReService.Code == models.ReServiceCodeIOASIT || serviceItem.ReService.Code == models.ReServiceCodeIDASIT {
		// address we need for the per_unit_cents is dependent on if it's origin/destination SIT
		var addressID uuid.UUID
		if serviceItem.ReService.Code == models.ReServiceCodeIOASIT {
			addressID = *ppmShipment.PickupAddressID
		} else {
			addressID = *ppmShipment.DestinationAddressID
		}

		var moveDate time.Time
		if ppmShipment.ActualMoveDate != nil {
			moveDate = *ppmShipment.ActualMoveDate
		} else {
			moveDate = ppmShipment.ExpectedDepartureDate
		}

		reServiceID, _ := models.FetchReServiceByCode(appCtx.DB(), serviceItem.ReService.Code)
		intlOtherPrice, _ := models.FetchReIntlOtherPrice(appCtx.DB(), addressID, reServiceID.ID, contract.ID, &moveDate)

		sitDaysParam := services.PricingDisplayParam{
			Key:   models.ServiceItemParamNameNumberDaysSIT,
			Value: strconv.Itoa(additionalDaysInSIT),
		}

		additionalDayPricer, ok := pricer.(services.IntlOriginAdditionalDaySITPricer)
		if !ok {
			return nil, nil, errors.New("ppm estimate pricer for SIT service item does not implement the first day pricer interface")
		}

		price, pricingParams, err := additionalDayPricer.Price(appCtx, contract.Code, moveDate, additionalDaysInSIT, *ppmShipment.SITEstimatedWeight, intlOtherPrice.PerUnitCents.Int())
		if err != nil {
			return nil, nil, err
		}

		pricingParams = append(pricingParams, sitDaysParam)

		appCtx.Logger().Debug(fmt.Sprintf("Pricing params for additional day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

		return &price, pricingParams, nil
	} else {
		// domestic PPMs
		additionalDaysPricer, ok := pricer.(services.DomesticAdditionalDaysSITPricer)
		if !ok {
			return nil, nil, errors.New("ppm estimate pricer for SIT service item does not implement the additional days pricer interface")
		}

		// Need to declare if origin or destination for the serviceAreaLookup, otherwise we already have it
		serviceAreaPostalCode := ppmShipment.PickupAddress.PostalCode
		serviceAreaKey := models.ServiceItemParamNameServiceAreaOrigin
		if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
			serviceAreaPostalCode = ppmShipment.DestinationAddress.PostalCode
			serviceAreaKey = models.ServiceItemParamNameServiceAreaDest
		}
		serviceAreaLookup := serviceparamvaluelookups.ServiceAreaLookup{
			Address: models.Address{PostalCode: serviceAreaPostalCode},
		}

		serviceArea, err := serviceAreaLookup.ParamValue(appCtx, contract.Code)
		if err != nil {
			return nil, nil, err
		}

		serviceAreaParam := services.PricingDisplayParam{
			Key:   serviceAreaKey,
			Value: serviceArea,
		}

		sitDaysParam := services.PricingDisplayParam{
			Key:   models.ServiceItemParamNameNumberDaysSIT,
			Value: strconv.Itoa(additionalDaysInSIT),
		}

		// Since this function may be ran before closeout, we need to account for if there's no actual move date yet.
		if ppmShipment.ActualMoveDate != nil {
			price, pricingParams, err := additionalDaysPricer.Price(appCtx, contract.Code, *ppmShipment.ActualMoveDate, *ppmShipment.SITEstimatedWeight, serviceArea, additionalDaysInSIT, true)
			if err != nil {
				return nil, nil, err
			}

			pricingParams = append(pricingParams, serviceAreaParam, sitDaysParam)

			appCtx.Logger().Debug(fmt.Sprintf("Pricing params for additional day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

			return &price, pricingParams, nil
		}
		price, pricingParams, err := additionalDaysPricer.Price(appCtx, contract.Code, ppmShipment.ExpectedDepartureDate, *ppmShipment.SITEstimatedWeight, serviceArea, additionalDaysInSIT, true)
		if err != nil {
			return nil, nil, err
		}

		pricingParams = append(pricingParams, serviceAreaParam, sitDaysParam)

		appCtx.Logger().Debug(fmt.Sprintf("Pricing params for additional day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

		return &price, pricingParams, nil
	}
}

// mapPPMShipmentEstimatedFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func MapPPMShipmentEstimatedFields(appCtx appcontext.AppContext, ppmShipment models.PPMShipment) (models.MTOShipment, error) {

	ppmShipment.Shipment.PPMShipment = &ppmShipment
	ppmShipment.Shipment.ShipmentType = models.MTOShipmentTypePPM
	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = ppmShipment.PickupAddress
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: ppmShipment.PickupAddress.PostalCode}
	ppmShipment.Shipment.DestinationAddress = ppmShipment.DestinationAddress
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: ppmShipment.DestinationAddress.PostalCode}
	ppmShipment.Shipment.PrimeActualWeight = ppmShipment.EstimatedWeight

	return ppmShipment.Shipment, nil
}

// MapPPMShipmentMaxIncentiveFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func MapPPMShipmentMaxIncentiveFields(appCtx appcontext.AppContext, ppmShipment models.PPMShipment, totalWeight unit.Pound) (models.MTOShipment, error) {
	var move models.Move
	err := appCtx.DB().Q().Eager(
		"Orders.Entitlement",
		"Orders.OriginDutyLocation.Address",
		"Orders.NewDutyLocation.Address",
	).Where("show = TRUE").Find(&move, ppmShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return models.MTOShipment{}, apperror.NewNotFoundError(ppmShipment.ID, " error querying move")
	}
	orders := move.Orders
	if orders.Entitlement.DBAuthorizedWeight == nil {
		return models.MTOShipment{}, apperror.NewNotFoundError(ppmShipment.ID, " DB authorized weight cannot be nil")
	}

	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: orders.OriginDutyLocation.Address.PostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: orders.NewDutyLocation.Address.PostalCode}
	ppmShipment.Shipment.PrimeActualWeight = &totalWeight

	return ppmShipment.Shipment, nil
}

// mapPPMShipmentFinalFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func MapPPMShipmentFinalFields(ppmShipment models.PPMShipment, totalWeight unit.Pound) models.MTOShipment {

	ppmShipment.Shipment.PPMShipment = &ppmShipment
	ppmShipment.Shipment.ShipmentType = models.MTOShipmentTypePPM
	ppmShipment.Shipment.ActualPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.RequestedPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.PickupAddress = ppmShipment.PickupAddress
	ppmShipment.Shipment.DestinationAddress = ppmShipment.DestinationAddress
	ppmShipment.Shipment.PrimeActualWeight = &totalWeight

	return ppmShipment.Shipment
}

// baseServiceItems returns a list of the MTOServiceItems that makeup the price of the estimated incentive.  These
// are the same non-accesorial service items that get auto-created and approved when the TOO approves an HHG shipment.
func BaseServiceItems(ppmShipment models.PPMShipment) []models.MTOServiceItem {
	mtoShipmentID := ppmShipment.ShipmentID
	isInternationalShipment := ppmShipment.Shipment.MarketCode == models.MarketCodeInternational

	if isInternationalShipment {
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeFSC}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeIHPK}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeIHUPK}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeISLH}, MTOShipmentID: &mtoShipmentID},
		}
	} else {
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeDLH}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeFSC}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeDOP}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeDDP}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeDPK}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeDUPK}, MTOShipmentID: &mtoShipmentID},
		}
	}
}

func StorageServiceItems(ppmShipment models.PPMShipment, locationType models.SITLocationType, additionalDaysInSIT int) []models.MTOServiceItem {
	mtoShipmentID := ppmShipment.ShipmentID
	isInternationalShipment := ppmShipment.Shipment.MarketCode == models.MarketCodeInternational

	// domestic shipments
	if locationType == models.SITLocationTypeOrigin && !isInternationalShipment {
		if additionalDaysInSIT > 0 {
			return []models.MTOServiceItem{
				{ReService: models.ReService{Code: models.ReServiceCodeDOFSIT}, MTOShipmentID: &mtoShipmentID},
				{ReService: models.ReService{Code: models.ReServiceCodeDOASIT}, MTOShipmentID: &mtoShipmentID},
			}
		}
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeDOFSIT}, MTOShipmentID: &mtoShipmentID}}
	}

	if locationType == models.SITLocationTypeDestination && !isInternationalShipment {
		if additionalDaysInSIT > 0 {
			return []models.MTOServiceItem{
				{ReService: models.ReService{Code: models.ReServiceCodeDDFSIT}, MTOShipmentID: &mtoShipmentID},
				{ReService: models.ReService{Code: models.ReServiceCodeDDASIT}, MTOShipmentID: &mtoShipmentID},
			}
		}
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeDDFSIT}, MTOShipmentID: &mtoShipmentID}}
	}

	// international shipments
	if locationType == models.SITLocationTypeOrigin && isInternationalShipment {
		if additionalDaysInSIT > 0 {
			return []models.MTOServiceItem{
				{ReService: models.ReService{Code: models.ReServiceCodeIOFSIT}, MTOShipmentID: &mtoShipmentID},
				{ReService: models.ReService{Code: models.ReServiceCodeIOASIT}, MTOShipmentID: &mtoShipmentID},
			}
		}
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeIOFSIT}, MTOShipmentID: &mtoShipmentID}}
	}

	if locationType == models.SITLocationTypeDestination && isInternationalShipment {
		if additionalDaysInSIT > 0 {
			return []models.MTOServiceItem{
				{ReService: models.ReService{Code: models.ReServiceCodeIDFSIT}, MTOShipmentID: &mtoShipmentID},
				{ReService: models.ReService{Code: models.ReServiceCodeIDASIT}, MTOShipmentID: &mtoShipmentID},
			}
		}
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeIDFSIT}, MTOShipmentID: &mtoShipmentID}}
	}

	return nil
}

// paramsForServiceCode filters the list of all service params for service items, to only those matching the service
// item code.  This allows us to make one initial query to the database instead of six just to filter by service item.
func paramsForServiceCode(code models.ReServiceCode, serviceParams models.ServiceParams) models.ServiceParams {
	var serviceItemParams models.ServiceParams
	for _, serviceParam := range serviceParams {
		if serviceParam.Service.Code == code {
			serviceItemParams = append(serviceItemParams, serviceParam)
		}
	}
	return serviceItemParams
}
