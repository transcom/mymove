package ppmcloseout

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmCloseoutFetcher struct {
	planner route.Planner
}

func NewPPMCloseoutFetcher(planner route.Planner) services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{planner: planner}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var ppmShipment models.PPMShipment
	var mtoShipment models.MTOShipment

	errPPM := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ShipmentID",
			"ExpectedDepartureDate",
			"ActualMoveDate",
			"EstimatedWeight",
			"HasProGear",
			"ProGearWeight",
			"SpouseProGearWeight",
			"FinalIncentive",
		).
		Find(&ppmShipment, ppmShipmentID)

	if errPPM != nil {
		switch errPPM {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", errPPM, "unable to find PPMShipment")
		}
	}

	mtoShipmentID := &ppmShipment.ShipmentID
	errMTO := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ScheduledPickupDate",
			"ActualPickupDate",
			"Distance",
			"PrimeActualWeight",
		).
		Find(&mtoShipment, mtoShipmentID)

	if errMTO != nil {
		switch errMTO {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*mtoShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", errMTO, "unable to find MTOShipment")
		}
	}

	// Get all DLH, FSC, DOP, DDP, DPK, and DUPK service items for the shipment
	var serviceItemsToPrice []models.MTOServiceItem
	var err error
	logger := appCtx.Logger()
	idString := ppmShipment.ShipmentID.String()
	fmt.Print(idString)
	err = appCtx.DB().Where("mto_shipment_id = ?", ppmShipment.ShipmentID).All(&serviceItemsToPrice)
	if err != nil {
		return nil, err
	}
	serviceItemsToPrice = ppmshipment.BaseServiceItems(ppmShipment.ShipmentID)
	logger.Debug(fmt.Sprintf("serviceItemsToPrice %+v", serviceItemsToPrice))
	// itemPricer := ghcrateengine.NewServiceItemPricer()
	contractDate := ppmShipment.ExpectedDepartureDate
	contract, contractErr := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if contractErr != nil {
		return nil, contractErr
	}

	paramsForServiceItems, paramErr := paymentrequesthelper.Helper.FetchServiceParamsForServiceItems(&paymentrequesthelper.RequestPaymentHelper{}, appCtx, serviceItemsToPrice)
	if paramErr != nil {
		return nil, paramErr
	}
	var totalPrice, packPrice, unpackPrice, destinationPrice, originPrice unit.Cents
	var totalWeight unit.Pound
	var ppmToMtoShipment models.MTOShipment

	if len(ppmShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			if weightTicket.Status != nil && *weightTicket.Status == models.PPMDocumentStatusRejected {
				totalWeight += 0
			} else if weightTicket.AdjustedNetWeight != nil {
				totalWeight += *weightTicket.AdjustedNetWeight
			} else if weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil {
				totalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
			}
		}
	}
	if totalWeight > 0 {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		ppmToMtoShipment = mapPPMShipmentFinalFields(ppmShipment, *ppmShipment.EstimatedWeight)
	} else {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		ppmToMtoShipment = mapPPMShipmentEstimatedFields(ppmShipment)
	}

	for _, serviceItem := range serviceItemsToPrice {
		pricer, pricerErr := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if pricerErr != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return nil, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(ppmToMtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(p.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code)

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
			paramValue, valueErr := keyData.ServiceParamValue(appCtx, paramKey.Key) // Fails with "DistanceZip" param?
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

		centsValue, paymentParams, priceErr := pricer.PriceUsingParams(appCtx, paramValues)
		if priceErr != nil {
			return nil, priceErr
		}
		logger.Debug(fmt.Sprintf("Service item price %s %d", serviceItem.ReService.Code, centsValue))
		logger.Debug(fmt.Sprintf("Payment service item params %+v", paymentParams))

		if err != nil {
			logger.Error("unable to calculate service item price", zap.Error(err))
			return nil, err
		}

		totalPrice = totalPrice.AddCents(centsValue)

		switch serviceItem.ReService.Code {
		case "DPK":
			packPrice += centsValue
		case "DUPK":
			unpackPrice += centsValue
		case "DOP":
			originPrice += centsValue
		case "DDP":
			destinationPrice += centsValue
			// TODO: Others can put cases for FSC DLH, etc. here, and the pricer *should* handle it.
		}
	}

	factor := float32(*mtoShipment.PrimeActualWeight) / float32(ghcrateengine.GetMinDomesticWeight())

	ppmCloseoutObj.ID = &ppmShipmentID
	ppmCloseoutObj.PlannedMoveDate = mtoShipment.ScheduledPickupDate
	ppmCloseoutObj.ActualMoveDate = mtoShipment.ActualPickupDate
	ppmCloseoutObj.Miles = (*int)(mtoShipment.Distance)
	ppmCloseoutObj.EstimatedWeight = ppmShipment.EstimatedWeight
	ppmCloseoutObj.ActualWeight = mtoShipment.PrimeActualWeight
	ppmCloseoutObj.ProGearWeightCustomer = ppmShipment.ProGearWeight
	ppmCloseoutObj.ProGearWeightSpouse = ppmShipment.SpouseProGearWeight
	ppmCloseoutObj.GrossIncentive = ppmShipment.FinalIncentive
	ppmCloseoutObj.GCC = nil
	ppmCloseoutObj.AOA = nil
	ppmCloseoutObj.RemainingReimbursementOwed = nil
	ppmCloseoutObj.HaulPrice = nil
	ppmCloseoutObj.HaulFSC = nil
	ppmCloseoutObj.DOP = &originPrice
	ppmCloseoutObj.DDP = &destinationPrice
	ppmCloseoutObj.Factor = &factor
	ppmCloseoutObj.PackPrice = &packPrice
	ppmCloseoutObj.UnpackPrice = &unpackPrice
	ppmCloseoutObj.SITReimbursement = ppmShipment.SITEstimatedCost

	if err != nil {
		return nil, apperror.NewQueryError("SITReimbursement", err, "error calculating SIT costs.")
	}

	return &ppmCloseoutObj, nil
}

func paramsForServiceCode(code models.ReServiceCode, serviceParams models.ServiceParams) models.ServiceParams {
	var serviceItemParams models.ServiceParams
	for _, serviceParam := range serviceParams {
		if serviceParam.Service.Code == code {
			serviceItemParams = append(serviceItemParams, serviceParam)
		}
	}
	return serviceItemParams
}

func mapPPMShipmentFinalFields(ppmShipment models.PPMShipment, totalWeight unit.Pound) models.MTOShipment {

	ppmShipment.Shipment.ActualPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.RequestedPickupDate = ppmShipment.ActualMoveDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: *ppmShipment.ActualPickupPostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: *ppmShipment.ActualDestinationPostalCode}
	ppmShipment.Shipment.PrimeActualWeight = &totalWeight

	return ppmShipment.Shipment
}

// mapPPMShipmentEstimatedFields remaps our PPMShipment specific information into the fields where the service param lookups
// expect to find them on the MTOShipment model.  This is only in-memory and shouldn't get saved to the database.
func mapPPMShipmentEstimatedFields(ppmShipment models.PPMShipment) models.MTOShipment {

	ppmShipment.Shipment.ActualPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.RequestedPickupDate = &ppmShipment.ExpectedDepartureDate
	ppmShipment.Shipment.PickupAddress = &models.Address{PostalCode: ppmShipment.PickupPostalCode}
	ppmShipment.Shipment.DestinationAddress = &models.Address{PostalCode: ppmShipment.DestinationPostalCode}
	ppmShipment.Shipment.PrimeActualWeight = ppmShipment.EstimatedWeight

	return ppmShipment.Shipment
}