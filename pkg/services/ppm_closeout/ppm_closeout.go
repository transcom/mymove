package ppmcloseout

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
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
	planner              route.Planner
	paymentRequestHelper paymentrequesthelper.Helper
	estimator            services.PPMEstimator
}

type serviceItemPrices struct {
	ddp                       *unit.Cents
	dop                       *unit.Cents
	packPrice                 *unit.Cents
	unpackPrice               *unit.Cents
	storageReimbursementCosts *unit.Cents
	haulPrice                 *unit.Cents
	haulFSC                   *unit.Cents
	haulType                  models.HaulType
}

func NewPPMCloseoutFetcher(planner route.Planner, paymentRequestHelper paymentrequesthelper.Helper, estimator services.PPMEstimator) services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{
		planner,
		paymentRequestHelper,
		estimator,
	}
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
			return nil, apperror.NewNotFoundError(ppmShipmentID, "unable to find PPMShipment")
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
	ppmCloseoutObj.DOP = nil
	ppmCloseoutObj.DDP = nil
	ppmCloseoutObj.PackPrice = nil
	ppmCloseoutObj.UnpackPrice = nil
	ppmCloseoutObj.SITReimbursement = nil

	// Check if PPM shipment is in "NEEDS_CLOSEOUT" or "CLOSEOUT_COMPLETE" status, if not, it's not ready for closeout
	if ppmShipment.Status != models.PPMShipmentStatusNeedsCloseout && ppmShipment.Status != models.PPMShipmentStatusCloseoutComplete {
		return nil, apperror.NewPPMNotReadyForCloseoutError(ppmShipmentID, "")
	}

	return &ppmShipment, err
}

func (p *ppmCloseoutFetcher) GetActualWeight(ppmShipment *models.PPMShipment) (unit.Pound, error) {
	var totalWeight unit.Pound
	if len(ppmShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			if weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil && (weightTicket.Status == nil || *weightTicket.Status != models.PPMDocumentStatusRejected) {
				totalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
			}
		}
	} else {
		return unit.Pound(0), apperror.NewPPMNoWeightTicketsError(ppmShipment.ID, "")
	}
	return totalWeight, nil
}

func (p *ppmCloseoutFetcher) GetExpenseStoragePrice(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (unit.Cents, error) {
	var expenseItems []models.MovingExpense
	var storageExpensePrice unit.Cents
	err := appCtx.DB().Where("ppm_shipment_id = ?", ppmShipmentID).All(&expenseItems)
	if err != nil {
		return unit.Cents(0), err
	}

	for _, movingExpense := range expenseItems {
		if movingExpense.MovingExpenseType != nil && *movingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage && *movingExpense.Status == models.PPMDocumentStatusApproved {
			storageExpensePrice += *movingExpense.Amount
		}
	}
	return storageExpensePrice, err
}

func (p *ppmCloseoutFetcher) GetEntitlement(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Entitlement, error) {
	var moveModel models.Move
	err := appCtx.DB().EagerPreload(
		"OrdersID",
	).Find(&moveModel, moveID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "unable to find Move")
		}
	}

	var order models.Order
	orderID := &moveModel.OrdersID
	errOrder := appCtx.DB().EagerPreload(
		"EntitlementID",
	).Find(&order, orderID)

	if errOrder != nil {
		switch errOrder {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*orderID, "while looking for Order")
		default:
			return nil, apperror.NewQueryError("Order", errOrder, "unable to find Order")
		}
	}

	var entitlement models.Entitlement
	entitlementID := order.EntitlementID
	errEntitlement := appCtx.DB().EagerPreload(
		"DBAuthorizedWeight",
	).Find(&entitlement, entitlementID)

	if errEntitlement != nil {
		switch errEntitlement {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*entitlementID, "while looking for Entitlement")
		default:
			return nil, apperror.NewQueryError("Entitlement", errEntitlement, "unable to find Entitlement")
		}
	}
	return &entitlement, nil
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

func (p *ppmCloseoutFetcher) getServiceItemPrices(appCtx appcontext.AppContext, ppmShipment models.PPMShipment) (serviceItemPrices, error) {
	// Get all DLH, FSC, DOP, DDP, DPK, and DUPK service items for the shipment
	var serviceItemsToPrice []models.MTOServiceItem
	var returnPriceObj serviceItemPrices
	logger := appCtx.Logger()

	err := appCtx.DB().Where("mto_shipment_id = ?", ppmShipment.ShipmentID).All(&serviceItemsToPrice)
	if err != nil {
		return serviceItemPrices{}, err
	}

	serviceItemsToPrice = ppmshipment.BaseServiceItems(ppmShipment.ShipmentID)

	// Change DLH to DSH if move within same Zip3
	actualPickupPostal := *ppmShipment.ActualPickupPostalCode
	actualDestPostal := *ppmShipment.ActualDestinationPostalCode
	if actualPickupPostal[0:3] == actualDestPostal[0:3] {
		serviceItemsToPrice[0] = models.MTOServiceItem{ReService: models.ReService{Code: models.ReServiceCodeDSH}, MTOShipmentID: &ppmShipment.ShipmentID}
	}
	contractDate := ppmShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return serviceItemPrices{}, err
	}

	paramsForServiceItems, paramErr := p.paymentRequestHelper.FetchServiceParamsForServiceItems(appCtx, serviceItemsToPrice)
	if paramErr != nil {
		return serviceItemPrices{}, paramErr
	}
	var totalPrice, packPrice, unpackPrice, destinationPrice, originPrice, haulPrice, haulFSC unit.Cents
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
		ppmToMtoShipment = ppmshipment.MapPPMShipmentFinalFields(ppmShipment, totalWeight)
	} else {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		ppmToMtoShipment = ppmshipment.MapPPMShipmentEstimatedFields(ppmShipment)
	}

	sitCosts, err := p.GetExpenseStoragePrice(appCtx, ppmShipment.ID)
	if err != nil {
		logger.Error("Error calculating SIT Reimbursement Costs", zap.Error(err))
		return serviceItemPrices{}, err
	}

	validCodes := map[models.ReServiceCode]string{
		models.ReServiceCodeDPK:  "DPK",
		models.ReServiceCodeDUPK: "DUPK",
		models.ReServiceCodeDOP:  "DOP",
		models.ReServiceCodeDDP:  "DDP",
		models.ReServiceCodeDSH:  "DSH",
		models.ReServiceCodeDLH:  "DLH",
		models.ReServiceCodeFSC:  "FSC",
	}

	// If service item is of a type we need for a specific calculation, get its price
	for _, serviceItem := range serviceItemsToPrice {
		_, isValidCode := validCodes[serviceItem.ReService.Code]
		if !isValidCode {
			continue
		} // Next iteration of loop if we don't need this service type

		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return serviceItemPrices{}, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(ppmToMtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(p.planner, serviceItemLookups, serviceItem, ppmToMtoShipment, contract.Code)

		// The distance value gets saved to the mto shipment model to reduce repeated api calls.
		var shipmentWithDistance models.MTOShipment
		err = appCtx.DB().Find(&shipmentWithDistance, ppmShipment.Shipment.ID)
		if err != nil {
			logger.Error("could not find shipment in the database")
			return serviceItemPrices{}, err
		}
		serviceItem.MTOShipment = shipmentWithDistance
		// set this to avoid potential eTag errors because the MTOShipment.Distance field was likely updated
		ppmShipment.Shipment = shipmentWithDistance

		var paramValues models.PaymentServiceItemParams
		for _, param := range paramsForServiceCode(serviceItem.ReService.Code, paramsForServiceItems) {
			paramKey := param.ServiceItemParamKey
			// This is where the lookup() method of each service item param is actually evaluated
			paramValue, serviceParamErr := keyData.ServiceParamValue(appCtx, paramKey.Key) // Fails with "DistanceZip" param?
			if serviceParamErr != nil {
				logger.Error("could not calculate param value lookup", zap.Error(serviceParamErr))
				return serviceItemPrices{}, serviceParamErr
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
			return serviceItemPrices{}, fmt.Errorf("no params were found for service item %s", serviceItem.ReService.Code)
		}

		// Middle var here can give you info on payment params like FSC multiplier, price rate/factor, etc. if needed.
		centsValue, _, err := pricer.PriceUsingParams(appCtx, paramValues)
		if err != nil {
			return serviceItemPrices{}, err
		}

		totalPrice = totalPrice.AddCents(centsValue)

		switch serviceItem.ReService.Code {
		case models.ReServiceCodeDPK:
			packPrice += centsValue
		case models.ReServiceCodeDUPK:
			unpackPrice += centsValue
		case models.ReServiceCodeDOP:
			originPrice += centsValue
		case models.ReServiceCodeDDP:
			destinationPrice += centsValue
		case models.ReServiceCodeDSH, models.ReServiceCodeDLH:
			haulPrice += centsValue
			_, linehaulOk := pricer.(services.DomesticLinehaulPricer)
			if linehaulOk {
				returnPriceObj.haulType = models.HaulType(models.LINEHAUL)
			} else {
				_, shorthaulOk := pricer.(services.DomesticShorthaulPricer)
				if shorthaulOk {
					returnPriceObj.haulType = models.HaulType(models.SHORTHAUL)
				} else { // Fallback in case pricer comparison fails
					if ppmToMtoShipment.DestinationAddress.PostalCode[0:3] == ppmToMtoShipment.PickupAddress.PostalCode[0:3] {
						returnPriceObj.haulType = models.HaulType(models.SHORTHAUL)
					} else {
						returnPriceObj.haulType = models.HaulType(models.LINEHAUL)
					}
				}
			}
		case models.ReServiceCodeFSC:
			haulFSC += centsValue
		}
	}
	returnPriceObj.ddp = &destinationPrice
	returnPriceObj.dop = &originPrice
	returnPriceObj.packPrice = &packPrice
	returnPriceObj.unpackPrice = &unpackPrice
	returnPriceObj.storageReimbursementCosts = &sitCosts
	returnPriceObj.haulPrice = &haulPrice
	returnPriceObj.haulFSC = &haulFSC

	return returnPriceObj, nil
}
