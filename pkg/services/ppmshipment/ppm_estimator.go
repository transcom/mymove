package ppmshipment

import (
	"github.com/gofrs/uuid"

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
	checks  []ppmShipmentValidator
	planner route.Planner
}

// NewEstimatePPM returns the estimatePPM (pass in checkRequiredFields() and checkEstimatedWeight)
func NewEstimatePPM(planner route.Planner) services.PPMEstimator {
	return &estimatePPM{
		checks: []ppmShipmentValidator{
			checkRequiredFields(),
			checkEstimatedWeight(),
		},
		planner: planner,
	}
}

// EstimateIncentiveWithDefaultChecks func that returns the estimate hard coded to 12K (because it'll be clear that the value is coming from teh service)
func (f *estimatePPM) EstimateIncentiveWithDefaultChecks(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment) (*int32, error) {
	return f.estimateIncentive(appCtx, oldPPMShipment, newPPMShipment, f.checks...)
}

func (f *estimatePPM) estimateIncentive(appCtx appcontext.AppContext, oldPPMShipment models.PPMShipment, newPPMShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*int32, error) {
	// Check that the PPMShipment has an ID
	var err error

	if newPPMShipment.Status != models.PPMShipmentStatusDraft {
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
	newPPMShipment.Advance = nil

	estimatedIncentive, err := f.calculatePrice(appCtx, *newPPMShipment)
	if err != nil {
		return nil, err
	}

	incentiveNumeric := int32(estimatedIncentive.Int())
	return &incentiveNumeric, nil
}

func (f estimatePPM) calculatePrice(appCtx appcontext.AppContext, ppmShipment models.PPMShipment) (*unit.Cents, error) {
	serviceItemsToPrice := f.estimateServiceItems(ppmShipment.ShipmentID)

	paymentHelper := paymentrequesthelper.RequestPaymentHelper{}
	// Get a unique list of all pricing params needed to calculate the estimate across all service items
	paramsForServiceItems, err := paymentHelper.FetchDistinctSystemServiceParamList(appCtx, serviceItemsToPrice)
	if err != nil {
		return nil, err
	}

	// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
	mtoShipment := mapPPMShipmentFields(ppmShipment)

	totalPrice := unit.Cents(0)
	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			return nil, err
		}

		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(mtoShipment, serviceItem)

		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(f.planner, serviceItemLookups, serviceItem, mtoShipment)

		// the distance value gets saved to the mto shipment model to reduce repeated api calls, we'll need to update
		// our in memory copy for pricers like FSC that try using that value directly.
		var shipmentWithDistance models.MTOShipment
		err = appCtx.DB().Find(&shipmentWithDistance, mtoShipment.ID)
		if err != nil {
			return nil, err
		}
		serviceItem.MTOShipment = shipmentWithDistance

		var paramValues models.PaymentServiceItemParams
		for _, param := range paramsForServiceItems {
			paramValue, valueErr := keyData.ServiceParamValue(appCtx, param.Key)
			if valueErr != nil {
				return nil, valueErr
			}

			paymentServiceItemParam := models.PaymentServiceItemParam{
				PaymentServiceItem: models.PaymentServiceItem{
					MTOServiceItem: serviceItem,
				},
				ServiceItemParamKey: param,
				Value:               paramValue,
			}
			paramValues = append(paramValues, paymentServiceItemParam)
		}

		centsValue, _, err := pricer.PriceUsingParams(appCtx, paramValues)

		if err != nil {
			return nil, err
		}

		totalPrice = totalPrice.AddCents(centsValue)
	}

	return &totalPrice, nil
}

func mapPPMShipmentFields(ppmShipment models.PPMShipment) models.MTOShipment {

	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: ppmShipment.PickupPostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: ppmShipment.DestinationPostalCode}
	ppmShipment.Shipment.PrimeActualWeight = ppmShipment.EstimatedWeight

	return ppmShipment.Shipment
}

func (f estimatePPM) estimateServiceItems(mtoShipmentID uuid.UUID) []models.MTOServiceItem {
	return []models.MTOServiceItem{
		{ReService: models.ReService{Code: models.ReServiceCodeDLH}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeFSC}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDOP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDDP}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDPK}, MTOShipmentID: &mtoShipmentID},
		{ReService: models.ReService{Code: models.ReServiceCodeDUPK}, MTOShipmentID: &mtoShipmentID},
	}
}
