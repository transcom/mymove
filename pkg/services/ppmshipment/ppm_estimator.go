package ppmshipment

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

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
			checkSITRequiredFields(),
		},
		planner:              planner,
		paymentRequestHelper: paymentRequestHelper,
	}
}

// EstimateIncentiveWithDefaultChecks func that returns the estimate hard coded to 12K (because it'll be clear that the value is coming from teh service)
func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*unit.Cents, *unit.Cents, error) {
	return f.estimateIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func shouldSkipEstimatingIncentive(newPPMShipment *models.PPMShipment, oldPPMShipment *models.PPMShipment) bool {
	return oldPPMShipment.ExpectedDepartureDate.Equal(newPPMShipment.ExpectedDepartureDate) &&
		newPPMShipment.PickupPostalCode == oldPPMShipment.PickupPostalCode &&
		newPPMShipment.DestinationPostalCode == oldPPMShipment.DestinationPostalCode &&
		((newPPMShipment.EstimatedWeight == nil && oldPPMShipment.EstimatedWeight == nil) || (oldPPMShipment.EstimatedWeight != nil && newPPMShipment.EstimatedWeight.Int() == oldPPMShipment.EstimatedWeight.Int()))
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
		newPPMShipment.PickupPostalCode != oldPPMShipment.PickupPostalCode ||
		newPPMShipment.DestinationPostalCode != oldPPMShipment.DestinationPostalCode ||
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

	estimatedIncentive := oldPPMShipment.EstimatedIncentive
	if !skipCalculatingEstimatedIncentive {
		// Clear out advance and advance requested fields when the estimated incentive is reset.
		newPPMShipment.HasRequestedAdvance = nil
		newPPMShipment.AdvanceAmountRequested = nil

		estimatedIncentive, err = f.calculatePrice(appCtx, newPPMShipment)
		if err != nil {
			return nil, nil, err
		}
	}

	estimatedSITCost := oldPPMShipment.SITEstimatedCost
	if calculateSITEstimate {
		estimatedSITCost, err = f.calculateSITCost(appCtx, newPPMShipment)
		if err != nil {
			return nil, nil, err
		}
	}

	return estimatedIncentive, estimatedSITCost, nil
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
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return nil, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(mtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment)

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
			logger.Error("unable to calculate service item price", zap.Error(err))
			return nil, err
		}

		totalPrice = totalPrice.AddCents(centsValue)
	}

	return &totalPrice, nil
}

func (f estimatePPM) calculateSITCost(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*unit.Cents, error) {
	logger := appCtx.Logger()

	additionalDaysInSIT := additionalDaysInSIT(*ppmShipment.SITEstimatedEntryDate, *ppmShipment.SITEstimatedDepartureDate)

	serviceItemsToPrice := storageServiceItems(ppmShipment.ShipmentID, *ppmShipment.SITLocation, additionalDaysInSIT)

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
			price, err = priceFirstDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment)
		case services.DomesticOriginAdditionalDaysSITPricer, services.DomesticDestinationAdditionalDaysSITPricer:
			price, err = priceAdditionalDaySIT(appCtx, serviceItemPricer, serviceItem, ppmShipment, additionalDaysInSIT)
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

func priceFirstDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment) (*unit.Cents, error) {
	firstDayPricer, ok := pricer.(services.DomesticFirstDaySITPricer)
	if !ok {
		return nil, errors.New("ppm estimate pricer for SIT service item does not implement the first day pricer interface")
	}

	serviceAreaPostalCode := ppmShipment.PickupPostalCode
	if serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
		serviceAreaPostalCode = ppmShipment.DestinationPostalCode
	}

	serviceAreaLookup := serviceparamvaluelookups.ServiceAreaLookup{
		Address: models.Address{PostalCode: serviceAreaPostalCode},
	}
	serviceArea, err := serviceAreaLookup.ParamValue(appCtx, ghcrateengine.DefaultContractCode)
	if err != nil {
		return nil, err
	}

	price, pricingParams, err := firstDayPricer.Price(appCtx, ghcrateengine.DefaultContractCode, ppmShipment.ExpectedDepartureDate, *ppmShipment.SITEstimatedWeight, serviceArea, true)
	if err != nil {
		return nil, err
	}

	appCtx.Logger().Debug(fmt.Sprintf("Pricing params for first day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

	return &price, nil
}

func additionalDaysInSIT(sitEntryDate time.Time, sitDepartureDate time.Time) int {
	// The cost of first day SIT is already accounted for in DOFSIT or DDFSIT service items
	if sitEntryDate.Equal(sitDepartureDate) {
		return 0
	}

	difference := sitDepartureDate.Sub(sitEntryDate)
	return int(difference.Hours() / 24)
}

func priceAdditionalDaySIT(appCtx appcontext.AppContext, pricer services.ParamsPricer, serviceItem models.MTOServiceItem, ppmShipment *models.PPMShipment, additionalDaysInSIT int) (*unit.Cents, error) {
	additionalDaysPricer, ok := pricer.(services.DomesticAdditionalDaysSITPricer)
	if !ok {
		return nil, errors.New("ppm estimate pricer for SIT service item does not implement the additional days pricer interface")
	}

	serviceAreaPostalCode := ppmShipment.PickupPostalCode
	if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
		serviceAreaPostalCode = ppmShipment.DestinationPostalCode
	}
	serviceAreaLookup := serviceparamvaluelookups.ServiceAreaLookup{
		Address: models.Address{PostalCode: serviceAreaPostalCode},
	}

	serviceArea, err := serviceAreaLookup.ParamValue(appCtx, ghcrateengine.DefaultContractCode)
	if err != nil {
		return nil, err
	}

	price, pricingParams, err := additionalDaysPricer.Price(appCtx, ghcrateengine.DefaultContractCode, ppmShipment.ExpectedDepartureDate, *ppmShipment.SITEstimatedWeight, serviceArea, additionalDaysInSIT, true)
	if err != nil {
		return nil, err
	}

	appCtx.Logger().Debug(fmt.Sprintf("Pricing params for additional day SIT %+v", pricingParams), zap.String("shipmentId", ppmShipment.ShipmentID.String()))

	return &price, nil
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

func storageServiceItems(mtoShipmentID uuid.UUID, locationType models.SITLocationType, additionalDaysInSIT int) []models.MTOServiceItem {
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
