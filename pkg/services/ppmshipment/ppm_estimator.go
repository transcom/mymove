package ppmshipment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
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
	if oldPPMShipment.Status != models.PPMShipmentStatusDraft && oldPPMShipment.EstimatedIncentive != nil && *newPPMShipment.EstimatedIncentive == 0 || oldPPMShipment.MaxIncentive == nil {
		return false
	} else {
		return oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate) &&
			newPPMShipment.PickupAddress.PostalCode == oldPPMShipment.PickupAddress.PostalCode &&
			newPPMShipment.DestinationAddress.PostalCode == oldPPMShipment.DestinationAddress.PostalCode &&
			((newPPMShipment.EstimatedWeight == nil && oldPPMShipment.EstimatedWeight == nil) || (oldPPMShipment.EstimatedWeight != nil && newPPMShipment.EstimatedWeight.Int() == oldPPMShipment.EstimatedWeight.Int()))
	}
}

func shouldSkipCalculatingFinalIncentive(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment, originalTotalWeight unit.Pound, newTotalWeight unit.Pound) bool {
	// If oldPPMShipment field value is nil we know that the value has been updated and we should return false - the adjusted net weight is accounted for in the
	// SumWeightTickets function and the change in weight is then checked with `newTotalWeight == originalTotalWeight`
	return (oldPPMShipment.ActualMoveDate != nil && newPPMShipment.ActualMoveDate.Equal(*oldPPMShipment.ActualMoveDate)) &&
		(oldPPMShipment.ActualPickupPostalCode != nil && *newPPMShipment.ActualPickupPostalCode == *oldPPMShipment.ActualPickupPostalCode) &&
		(oldPPMShipment.ActualDestinationPostalCode != nil && *newPPMShipment.ActualDestinationPostalCode == *oldPPMShipment.ActualDestinationPostalCode) &&
		newTotalWeight == originalTotalWeight
}

func shouldSetFinalIncentiveToNil(newPPMShipment *models.PPMShipment, newTotalWeight unit.Pound) bool {
	if newPPMShipment.ActualMoveDate == nil || newPPMShipment.ActualPickupPostalCode == nil || newPPMShipment.ActualDestinationPostalCode == nil || newTotalWeight <= 0 {
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

	calculateSITEstimate := shouldCalculateSITCost(newPPMShipment, &oldPPMShipment)

	// Clear out any previously calculated SIT estimated costs, if SIT is no longer expected
	if newPPMShipment.SITExpected != nil && !*newPPMShipment.SITExpected {
		newPPMShipment.SITEstimatedCost = nil
	}

	skipCalculatingEstimatedIncentive := shouldSkipEstimatingIncentive(newPPMShipment, &oldPPMShipment)

	if skipCalculatingEstimatedIncentive && !calculateSITEstimate {
		return oldPPMShipment.EstimatedIncentive, newPPMShipment.SITEstimatedCost, nil
	}

	contractDate := newPPMShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, nil, err
	}

	estimatedIncentive := oldPPMShipment.EstimatedIncentive
	if !skipCalculatingEstimatedIncentive {
		// Clear out advance and advance requested fields when the estimated incentive is reset.
		newPPMShipment.HasRequestedAdvance = nil
		newPPMShipment.AdvanceAmountRequested = nil

		estimatedIncentive, err = f.calculatePrice(appCtx, newPPMShipment, 0, contract, false)
		if err != nil {
			return nil, nil, err
		}
	}

	estimatedSITCost := oldPPMShipment.SITEstimatedCost
	if calculateSITEstimate {
		estimatedSITCost, err = CalculateSITCost(appCtx, newPPMShipment, contract)
		if err != nil {
			return nil, nil, err
		}
	}

	return estimatedIncentive, estimatedSITCost, nil
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
		"Orders.Entitlement",
	).Where("show = TRUE").Find(&move, newPPMShipment.Shipment.MoveTaskOrderID)
	if err != nil {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " error querying move")
	}
	orders := move.Orders
	if orders.Entitlement.DBAuthorizedWeight == nil {
		return nil, apperror.NewNotFoundError(newPPMShipment.ID, " DB authorized weight cannot be nil")
	}

	contractDate := newPPMShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	// since the max incentive is based off of the authorized weight entitlement and that value CAN change
	// we will calculate the max incentive each time it is called
	maxIncentive, err := f.calculatePrice(appCtx, newPPMShipment, unit.Pound(*orders.Entitlement.DBAuthorizedWeight), contract, true)
	if err != nil {
		return nil, err
	}

	return maxIncentive, nil
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
	originalTotalWeight, newTotalWeight := SumWeightTickets(oldPPMShipment, *newPPMShipment)

	if newPPMShipment.AllowableWeight != nil && *newPPMShipment.AllowableWeight < newTotalWeight {
		newTotalWeight = *newPPMShipment.AllowableWeight
	}

	isMissingInfo := shouldSetFinalIncentiveToNil(newPPMShipment, newTotalWeight)
	var skipCalculateFinalIncentive bool
	finalIncentive := oldPPMShipment.FinalIncentive

	if !isMissingInfo {
		skipCalculateFinalIncentive = shouldSkipCalculatingFinalIncentive(newPPMShipment, &oldPPMShipment, originalTotalWeight, newTotalWeight)
		if !skipCalculateFinalIncentive {
			contractDate := newPPMShipment.ExpectedDepartureDate
			if newPPMShipment.ActualMoveDate != nil {
				contractDate = *newPPMShipment.ActualMoveDate
			}
			contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
			if err != nil {
				return nil, err
			}

			finalIncentive, err = f.calculatePrice(appCtx, newPPMShipment, newTotalWeight, contract, false)
			if err != nil {
				return nil, err
			}
		}
	} else {
		finalIncentive = nil
	}

	return finalIncentive, nil
}

// SumWeightTickets return the total weight of all weightTickets associated with a PPMShipment, returns 0 if there is no valid weight
func SumWeightTickets(ppmShipment, newPPMShipment models.PPMShipment) (originalTotalWeight, newTotalWeight unit.Pound) {
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

	return originalTotalWeight, newTotalWeight
}

// calculatePrice returns an incentive value for the ppm shipment as if we were pricing the service items for
// an HHG shipment with the same values for a payment request.  In this case we're not persisting service items,
// MTOServiceItems or PaymentRequestServiceItems, to the database to avoid unnecessary work and get a quicker result.
// we use this when calculating the estimated, final, and max incentive values
func (f estimatePPM) calculatePrice(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, totalWeight unit.Pound, contract models.ReContract, isMaxIncentiveCheck bool) (*unit.Cents, error) {
	logger := appCtx.Logger()

	zeroTotal := false
	serviceItemsToPrice := BaseServiceItems(ppmShipment.ShipmentID)

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
		if ppmShipment.ActualPickupPostalCode != nil {
			pickupPostal = *ppmShipment.ActualPickupPostalCode
		} else if ppmShipment.PickupAddress.PostalCode != "" {
			pickupPostal = ppmShipment.PickupAddress.PostalCode
		}

		if ppmShipment.ActualDestinationPostalCode != nil {
			destPostal = *ppmShipment.ActualDestinationPostalCode
		} else if ppmShipment.DestinationAddress.PostalCode != "" {
			destPostal = ppmShipment.DestinationAddress.PostalCode
		}
	}

	if pickupPostal[0:3] == destPostal[0:3] {
		serviceItemsToPrice[0] = models.MTOServiceItem{ReService: models.ReService{Code: models.ReServiceCodeDSH}, MTOShipmentID: &ppmShipment.ShipmentID}
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
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code)

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
		logger.Debug(fmt.Sprintf("Service item price %s %d", serviceItem.ReService.Code, centsValue))
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

	serviceItemsToPrice := BaseServiceItems(ppmShipment.ShipmentID)

	// Replace linehaul pricer with shorthaul pricer if move is within the same Zip3
	var pickupPostal, destPostal string

	// Check different address values for a postal code
	if ppmShipment.ActualPickupPostalCode != nil {
		pickupPostal = *ppmShipment.ActualPickupPostalCode
	} else if ppmShipment.PickupAddress.PostalCode != "" {
		pickupPostal = ppmShipment.PickupAddress.PostalCode
	}

	// Same for destination
	if ppmShipment.ActualDestinationPostalCode != nil {
		destPostal = *ppmShipment.ActualDestinationPostalCode
	} else if ppmShipment.DestinationAddress.PostalCode != "" {
		destPostal = ppmShipment.DestinationAddress.PostalCode
	}

	if len(pickupPostal) >= 3 && len(destPostal) >= 3 && pickupPostal[:3] == destPostal[:3] {
		if pickupPostal[0:3] == destPostal[0:3] {
			serviceItemsToPrice[0] = models.MTOServiceItem{ReService: models.ReService{Code: models.ReServiceCodeDSH}, MTOShipmentID: &ppmShipment.ShipmentID}
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

	var totalWeightFromWeightTickets unit.Pound
	var blankPPM models.PPMShipment
	if ppmShipment.WeightTickets != nil {
		_, totalWeightFromWeightTickets = SumWeightTickets(blankPPM, *ppmShipment)
	} else {
		return emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, emptyPrice, apperror.NewPPMNoWeightTicketsError(ppmShipment.ID, " no weight tickets")
	}

	var mtoShipment models.MTOShipment
	if totalWeightFromWeightTickets > 0 {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		mtoShipment = MapPPMShipmentFinalFields(*ppmShipment, totalWeightFromWeightTickets)
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
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code)

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

func CalculateSITCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, contract models.ReContract) (*unit.Cents, error) {
	logger := appCtx.Logger()

	additionalDaysInSIT := additionalDaysInSIT(*ppmShipment.SITEstimatedEntryDate, *ppmShipment.SITEstimatedDepartureDate)

	serviceItemsToPrice := StorageServiceItems(ppmShipment.ShipmentID, *ppmShipment.SITLocation, additionalDaysInSIT)

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
}

func CalculateSITCostBreakdown(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, contract models.ReContract) (*models.PPMSITEstimatedCostInfo, error) {
	logger := appCtx.Logger()

	ppmSITEstimatedCostInfoData := &models.PPMSITEstimatedCostInfo{}

	additionalDaysInSIT := additionalDaysInSIT(*ppmShipment.SITEstimatedEntryDate, *ppmShipment.SITEstimatedDepartureDate)

	serviceItemsToPrice := StorageServiceItems(ppmShipment.ShipmentID, *ppmShipment.SITLocation, additionalDaysInSIT)

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
		case services.DomesticOriginAdditionalDaysSITPricer, services.DomesticDestinationAdditionalDaysSITPricer:
			price, ppmSITEstimatedCostInfoData, err = calculateAdditionalDaySITCostBreakdown(appCtx, serviceItemPricer, serviceItem, ppmShipment, contract, additionalDaysInSIT, ppmSITEstimatedCostInfoData, logger)
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

func priceFirstDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, contract models.ReContract) (*unit.Cents, services.PricingDisplayParams, error) {
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

func additionalDaysInSIT(sitEntryDate time.Time, sitDepartureDate time.Time) int {
	// The cost of first day SIT is already accounted for in DOFSIT or DDFSIT service items
	if sitEntryDate.Equal(sitDepartureDate) {
		return 0
	}

	difference := sitDepartureDate.Sub(sitEntryDate)
	return int(difference.Hours() / 24)
}

func priceAdditionalDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, additionalDaysInSIT int, contract models.ReContract) (*unit.Cents, services.PricingDisplayParams, error) {
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

// mapPPMShipmentEstimatedFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func MapPPMShipmentEstimatedFields(appCtx appcontext.AppContext, ppmShipment models.PPMShipment) (models.MTOShipment, error) {

	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: ppmShipment.PickupAddress.PostalCode}
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

	ppmShipment.Shipment.ActualPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.RequestedPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: *ppmShipment.ActualPickupPostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: *ppmShipment.ActualDestinationPostalCode}
	ppmShipment.Shipment.PrimeActualWeight = &totalWeight

	return ppmShipment.Shipment
}

// baseServiceItems returns a list of the MTOServiceItems that makeup the price of the estimated incentive.  These
// are the same non-accesorial service items that get auto-created and approved when the TOO approves an HHG shipment.
func BaseServiceItems(mtoShipmentID uuid.UUID) []models.MTOServiceItem {
	return []models.MTOServiceItem{
		{ReService: models.ReService{Code: models.ReServiceCodeDLH}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeFSC}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDOP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDDP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDPK}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDUPK}, MTOShipmentID: &mtoShipmentID},
	}
}

func StorageServiceItems(mtoShipmentID uuid.UUID, locationType models.SITLocationType, additionalDaysInSIT int) []models.MTOServiceItem {
	if locationType == models.SITLocationTypeOrigin {
		if additionalDaysInSIT > 0 {
			return []models.MTOServiceItem{
				{ReService: models.ReService{Code: models.ReServiceCodeDOFSIT}, MTOShipmentID: &mtoShipmentID},
				{ReService: models.ReService{Code: models.ReServiceCodeDOASIT}, MTOShipmentID: &mtoShipmentID},
			}
		}
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeDOFSIT}, MTOShipmentID: &mtoShipmentID}}
	}

	if additionalDaysInSIT > 0 {
		return []models.MTOServiceItem{
			{ReService: models.ReService{Code: models.ReServiceCodeDDFSIT}, MTOShipmentID: &mtoShipmentID},
			{ReService: models.ReService{Code: models.ReServiceCodeDDASIT}, MTOShipmentID: &mtoShipmentID},
		}
	}

	return []models.MTOServiceItem{
		{ReService: models.ReService{Code: models.ReServiceCodeDDFSIT}, MTOShipmentID: &mtoShipmentID}}
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
