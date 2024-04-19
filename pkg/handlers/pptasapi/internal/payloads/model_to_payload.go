package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func MovesSince(appCtx appcontext.AppContext, moves models.Moves) *pptasmessages.GetMovesSinceResponse {
	searchMoves := make(pptasmessages.SearchMoves, len(moves))

	for i, move := range moves {
		customer := move.Orders.ServiceMember

		numShipments := 0
		for _, shipment := range move.MTOShipments {
			if shipment.Status != models.MTOShipmentStatusDraft {
				numShipments++
			}
		}

		var pickupDate, deliveryDate *strfmt.Date

		if numShipments > 0 && move.MTOShipments[0].ScheduledPickupDate != nil {
			pickupDate = handlers.FmtDatePtr(move.MTOShipments[0].ScheduledPickupDate)
		} else {
			pickupDate = nil
		}

		if numShipments > 0 && move.MTOShipments[0].ScheduledDeliveryDate != nil {
			deliveryDate = handlers.FmtDatePtr(move.MTOShipments[0].ScheduledDeliveryDate)
		} else {
			deliveryDate = nil
		}

		var originGBLOC pptasmessages.GBLOC
		if move.Status == models.MoveStatusNeedsServiceCounseling {
			originGBLOC = pptasmessages.GBLOC(*move.Orders.OriginDutyLocationGBLOC)
		} else if len(move.ShipmentGBLOC) > 0 {
			// There is a Pop bug that prevents us from using a has_one association for
			// Move.ShipmentGBLOC, so we have to treat move.ShipmentGBLOC as an array, even
			// though there can never be more than one GBLOC for a move.
			if move.ShipmentGBLOC[0].GBLOC != nil {
				originGBLOC = pptasmessages.GBLOC(*move.ShipmentGBLOC[0].GBLOC)
			}
		} else {
			// If the move's first shipment doesn't have a pickup address (like with an NTS-Release),
			// we need to fall back to the origin duty location GBLOC.  If that's not available for
			// some reason, then we should get the empty string (no GBLOC).
			originGBLOC = pptasmessages.GBLOC(*move.Orders.OriginDutyLocationGBLOC)
		}

		var destinationGBLOC pptasmessages.GBLOC
		var PostalCodeToGBLOC models.PostalCodeToGBLOC
		var err error
		if numShipments > 0 && move.MTOShipments[0].DestinationAddress != nil {
			PostalCodeToGBLOC, err = models.FetchGBLOCForPostalCode(appCtx.DB(), move.MTOShipments[0].DestinationAddress.PostalCode)
		} else {
			// If the move has no shipments or the shipment has no destination address fall back to the origin duty location GBLOC
			PostalCodeToGBLOC, err = models.FetchGBLOCForPostalCode(appCtx.DB(), move.Orders.NewDutyLocation.Address.PostalCode)
		}

		if err != nil {
			destinationGBLOC = *pptasmessages.NewGBLOC("")
		} else {
			destinationGBLOC = pptasmessages.GBLOC(PostalCodeToGBLOC.GBLOC)
		}

		searchMoves[i] = &pptasmessages.SearchMove{
			FirstName:                         customer.FirstName,
			LastName:                          customer.LastName,
			DodID:                             customer.Edipi,
			Branch:                            customer.Affiliation.String(),
			Status:                            pptasmessages.MoveStatus(move.Status),
			ID:                                *handlers.FmtUUID(move.ID),
			Locator:                           move.Locator,
			ShipmentsCount:                    int64(numShipments),
			OriginDutyLocationPostalCode:      move.Orders.OriginDutyLocation.Address.PostalCode,
			DestinationDutyLocationPostalCode: move.Orders.NewDutyLocation.Address.PostalCode,
			OrderType:                         string(move.Orders.OrdersType),
			RequestedPickupDate:               pickupDate,
			RequestedDeliveryDate:             deliveryDate,
			OriginGBLOC:                       originGBLOC,
			DestinationGBLOC:                  destinationGBLOC,
		}
	}

	payload := pptasmessages.GetMovesSinceResponse{
		MovesFound: searchMoves,
	}

	return &payload
}
