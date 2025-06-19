package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
)

// Helper struct to hold the outputted table from the PPM Closeout DB func
type PPMCloseoutQueueItem struct {
	ID                               *uuid.UUID         `json:"id" db:"id"`
	Show                             *bool              `json:"show" db:"show"`
	Locator                          *string            `json:"locator" db:"locator"`
	FullOrPartialPPM                 *string            `json:"full_or_partial_ppm" db:"full_or_partial_ppm"`
	OrdersID                         *uuid.UUID         `json:"orders_id" db:"orders_id"`
	LockedBy                         *uuid.UUID         `json:"locked_by" db:"locked_by"`
	LockExpiresAt                    *time.Time         `json:"lock_expires_at" db:"lock_expires_at"`
	SCCloseoutAssignedID             *uuid.UUID         `json:"sc_closeout_assigned_id" db:"sc_closeout_assigned_id"`
	CounselingTransportationOfficeID *uuid.UUID         `json:"counseling_transportation_office_id" db:"counseling_transportation_office_id"`
	Orders                           json.RawMessage    `json:"orders" db:"orders"`
	PpmShipments                     json.RawMessage    `json:"ppm_shipments" db:"ppm_shipments"`
	CounselingTransportationOffice   json.RawMessage    `json:"counseling_transportation_office" db:"counseling_transportation_office"`
	PpmCloseoutLocation              json.RawMessage    `json:"ppm_closeout_location" db:"ppm_closeout_location"`
	ScAssigned                       json.RawMessage    `json:"sc_assigned" db:"sc_assigned"`
	MoveStatus                       *models.MoveStatus `json:"status" db:"status"`
	MtoShipments                     json.RawMessage    `json:"mto_shipments" db:"mto_shipments"`
	TotalCount                       int                `json:"total_count" db:"total_count"`
}

func (f orderFetcher) ListPPMCloseoutOrders(
	appCtx appcontext.AppContext,
	officeUserID uuid.UUID,
	params *services.ListOrderParams,
) ([]models.Move, int, error) {
	if params.SubmittedAt != nil {
		// If this is present, most likely somebody is passing this parameter
		// thinking it is the closeout date. SubmittedAt param is the
		// move's submitted at filter. params.CloseoutInitiated is the ppm shipment
		// submitted at filter
		err := errors.New("submitted at parameter should not be used for PPM closeout queue. Please use closeout initiated instead")
		appCtx.Logger().Error("Incorrect parameter used for PPM closeout queue", zap.Error(err))
		return nil, 0, err
	}
	var ppmCloseoutQueueItems []PPMCloseoutQueueItem

	var officeUserGbloc string
	if params.ViewAsGBLOC != nil {
		officeUserGbloc = *params.ViewAsGBLOC
	} else {
		var err error
		officeUserGbloc, err = officeuser.NewOfficeUserGblocFetcher().
			FetchGblocForOfficeUser(appCtx, officeUserID)
		if err != nil {
			return nil, 0, err
		}
	}

	hasSafetyPrivilege := false
	if privs, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID); err == nil {
		hasSafetyPrivilege = privs.HasPrivilege(roles.PrivilegeTypeSafety)
	} else {
		appCtx.Logger().Error("Error retrieving user privileges", zap.Error(err))
		return nil, 0, err
	}

	// Leave as nil pointers to rely on proc default
	var page *int
	if params.Page != nil {
		paramPagePtr := int(*params.Page)
		page = &paramPagePtr
	}
	var perPage *int
	if params.PerPage != nil {
		paramPerPagePtr := int(*params.PerPage)
		perPage = &paramPerPagePtr
	}

	const q = `
        SELECT *
          FROM get_ppm_closeout_queue(
            $1,  $2,  $3,  $4,  $5,  $6,  $7,  $8,  $9,  $10,
           $11, $12, $13, $14, $15, $16, $17, $18, $19
        )`

	err := appCtx.DB().
		RawQuery(q,
			officeUserGbloc,
			params.CustomerName,
			params.Edipi,
			params.Emplid,
			pq.Array(params.Status),
			params.Locator,
			params.CloseoutInitiated,
			params.Branch,
			params.PPMType,
			strings.Join(params.OriginDutyLocation, " "),
			params.CounselingOffice,
			params.DestinationDutyLocation,
			params.CloseoutLocation,
			params.AssignedTo,
			hasSafetyPrivilege,
			page,
			perPage,
			params.Sort,
			params.Order,
		).
		All(&ppmCloseoutQueueItems)
	if err != nil {
		return nil, 0, err
	}

	moves, err := mapPPMCloseoutQueueItemsToMoves(ppmCloseoutQueueItems)
	if err != nil {
		return nil, 0, err
	}

	var totalCount int
	if len(ppmCloseoutQueueItems) > 0 {
		totalCount = ppmCloseoutQueueItems[0].TotalCount
	}

	return moves, totalCount, nil
}

func mapPPMCloseoutQueueItemsToMoves(queueItems []PPMCloseoutQueueItem) ([]models.Move, error) {
	var moves []models.Move

	for _, queueItem := range queueItems {
		var move models.Move

		if queueItem.ID != nil {
			move.ID = *queueItem.ID
		}
		if queueItem.Show != nil {
			move.Show = queueItem.Show
		}
		if queueItem.Locator != nil {
			move.Locator = *queueItem.Locator
		}
		if queueItem.MoveStatus != nil {
			move.Status = *queueItem.MoveStatus
		}

		move.PPMType = queueItem.FullOrPartialPPM

		if queueItem.OrdersID != nil {
			move.OrdersID = *queueItem.OrdersID
		}
		if queueItem.LockedBy != nil {
			move.LockedByOfficeUserID = queueItem.LockedBy
		}
		if queueItem.LockExpiresAt != nil {
			move.LockExpiresAt = queueItem.LockExpiresAt
		}

		if queueItem.SCCloseoutAssignedID != nil {
			move.SCCounselingAssignedID = queueItem.SCCloseoutAssignedID
		}
		if queueItem.CounselingTransportationOfficeID != nil {
			move.CounselingOfficeID = queueItem.CounselingTransportationOfficeID
		}

		var order models.Order
		if err := json.Unmarshal(queueItem.Orders, &order); err != nil {
			return nil, fmt.Errorf("unmarshal Orders JSON: %w", err)
		}
		move.OrdersID = order.ID
		move.Orders = order

		var counselOffice *models.TransportationOffice
		if err := json.Unmarshal(queueItem.CounselingTransportationOffice, &counselOffice); err != nil {
			return nil, fmt.Errorf("unmarshal CounselingTransportationOffice JSON: %w", err)
		}
		if counselOffice != nil {
			move.CounselingOfficeID = &counselOffice.ID
			move.CounselingOffice = counselOffice
		}

		var closeOffice *models.TransportationOffice
		if err := json.Unmarshal(queueItem.PpmCloseoutLocation, &closeOffice); err != nil {
			return nil, fmt.Errorf("unmarshal PpmCloseoutLocation JSON: %w", err)
		}
		if closeOffice != nil {
			move.CloseoutOfficeID = &closeOffice.ID
			move.CloseoutOffice = closeOffice
		}

		var scUser *models.OfficeUser
		if err := json.Unmarshal(queueItem.ScAssigned, &scUser); err != nil {
			return nil, fmt.Errorf("unmarshal ScAssigned JSON: %w", err)
		}
		if scUser != nil {
			move.SCCloseoutAssignedID = &scUser.ID
			move.SCCloseoutAssignedUser = scUser
		}

		// Account for the json agg of multiple mts
		var mtoShipments models.MTOShipments
		if err := json.Unmarshal(queueItem.MtoShipments, &mtoShipments); err != nil {
			return nil, fmt.Errorf("unmarshal MtoShipments JSON: %w", err)
		}
		// Tie the PPM before appending
		var ppmShipments models.PPMShipments
		if err := json.Unmarshal(queueItem.PpmShipments, &ppmShipments); err != nil {
			return nil, fmt.Errorf("unmarshal PpmShipments JSON: %w", err)
		}
		for i, shipment := range mtoShipments {
			for _, ppmShipment := range ppmShipments {
				if ppmShipment.ShipmentID == shipment.ID {
					mtoShipments[i].PPMShipment = &ppmShipment
				}
			}
		}

		move.MTOShipments = append(move.MTOShipments, mtoShipments...)

		moves = append(moves, move)
	}

	return moves, nil
}
