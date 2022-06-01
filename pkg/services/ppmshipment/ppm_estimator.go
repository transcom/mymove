package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"
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
		},
		planner:              planner,
		paymentRequestHelper: paymentRequestHelper,
	}
}

// EstimateIncentiveWithDefaultChecks func that returns the estimate hard coded to 12K (because it'll be clear that the value is coming from teh service)
func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, error) {
	return f.estimateIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*unit.Cents, error) {
	// Check that the PPMShipment has an ID
	var err error

	if newPPMShipment.Status != models.PPMShipmentStatusDraft && newPPMShipment.Status != models.PPMShipmentStatusSubmitted {
		return nil, nil
	}
	// Check that all the required fields we need are present.
	err = validatePPMShipment(appCtx, *newPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
	// If a field does not pass validation return nil as error handling is happening in the validator
	if err != nil {
		switch err.(type) {
		case apperror.InvalidInputError:
			return nil, nil
		default:
			return nil, err
		}
	}

	if oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate) && newPPMShipment.PickupPostalCode == oldPPMShipment.PickupPostalCode && newPPMShipment.DestinationPostalCode == oldPPMShipment.DestinationPostalCode && oldPPMShipment.EstimatedWeight != nil && *newPPMShipment.EstimatedWeight == *oldPPMShipment.EstimatedWeight {
		return oldPPMShipment.EstimatedIncentive, nil
	}
	// Clear out advance and advance requested fields when the estimated incentive is reset.
	newPPMShipment.AdvanceRequested = nil
	newPPMShipment.HasRequestedAdvance = nil
	newPPMShipment.Advance = nil
	newPPMShipment.AdvanceAmountRequested = nil

	estimatedIncentive, err := f.calculatePrice(appCtx, newPPMShipment)
	if err != nil {
		return nil, err
	}

	return estimatedIncentive, nil
}

// calculatePrice returns an estimated incentive value for the ppm shipment as if we were pricing the service items for
// an HHG shipment with the same values for a payment request.  In this case we're not persisting service items,
// MTOServiceItems or PaymentRequestServiceItems, to the database to avoid unnecessary work and get a quicker result.
func (f estimatePPM) calculatePrice(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*unit.Cents, error) {
	logger := appCtx.Logger()

	serviceItemsToPrice := estimateServiceItems(ppmShipment.ShipmentID)

	// Get a list of all the pricing params needed to calculate the price for each service item
	paramsForServiceItems, err := f.paymentRequestHelper.FetchServiceParamsForServiceItems(appCtx, serviceItemsToPrice)
	if err != nil {
		logger.Error("fetching PPM estimate ServiceParams failed", zap.Error(err))
		return nil, err
	}

	// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
	mtoShipment := mapPPMShipmentFields(*ppmShipment)

	totalPrice := unit.Cents(0)
	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("not able to find pricer for service item", zap.Error(err))
			return nil, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(mtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment)

		// The distance value gets saved to the mto shipment model to reduce repeated api calls. We'll need to update
		// our in memory copy for pricers, like FSC, that try using that saved value directly.
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
				logger.Error("could not calculate param value lookup", zap.Error(err))
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

		centsValue, _, err := pricer.PriceUsingParams(appCtx, paramValues)
		if err != nil {
			logger.Error("unable to calculate service item price", zap.Error(err))
			return nil, err
		}

		totalPrice = totalPrice.AddCents(centsValue)
	}

	return &totalPrice, nil
}

// mapPPMShipmentFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func mapPPMShipmentFields(ppmShipment models.PPMShipment) models.MTOShipment {

	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: ppmShipment.PickupPostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: ppmShipment.DestinationPostalCode}
	ppmShipment.Shipment.PrimeActualWeight = ppmShipment.EstimatedWeight

	return ppmShipment.Shipment
}

// estimateServiceItems returns a list of the MTOServiceItems that makeup the price of the estimated incentive.  These
// are the same non-accesorial service items that get auto-created and approved when the TOO approves an HHG shipment.
func estimateServiceItems(mtoShipmentID uuid.UUID) []models.MTOServiceItem {
	return []models.MTOServiceItem{
		{ReService: models.ReService{Code: models.ReServiceCodeDLH}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeFSC}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDOP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDDP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDPK}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDUPK}, MTOShipmentID: &mtoShipmentID},
	}
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
