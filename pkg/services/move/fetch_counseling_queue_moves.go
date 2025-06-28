package move

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
)

type counselingQueueFetcher struct {
}

func NewCounselingQueueFetcher() services.CounselingQueueFetcher {
	return &counselingQueueFetcher{}
}

func (o *counselingQueueFetcher) FetchCounselingQueue(appCtx appcontext.AppContext, counselingQueueParams services.CounselingQueueParams) ([]models.Move, int64, error) {
	movesWithCount, err := getCounselingQueueDbFunc(counselingQueueParams, appCtx)
	if err != nil {
		appCtx.Logger().
			Error("error in method getCounselingQueueDbFunc. Failed to fetch list of moves for the counseling queue", zap.Error(err))
		return models.Moves{}, 0, err
	}

	count := int64(0)
	if len(movesWithCount) > 0 {
		count = movesWithCount[0].TotalCount
	}

	moves, err := movesWithCountToMoves(movesWithCount)
	if err != nil {
		return moves, 0, err
	}

	return moves, count, nil
}

func getCounselingQueueDbFunc(counselingQueueParams services.CounselingQueueParams, appCtx appcontext.AppContext) ([]MoveWithCount, error) {
	var movesWithCount []MoveWithCount

	var officeUserGbloc string
	if counselingQueueParams.ViewAsGBLOC != nil {
		officeUserGbloc = *counselingQueueParams.ViewAsGBLOC
	} else {
		var gblocErr error
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		officeUserGbloc, gblocErr = gblocFetcher.FetchGblocForOfficeUser(appCtx, appCtx.Session().OfficeUserID)
		if gblocErr != nil {
			return movesWithCount, gblocErr
		}
	}

	err := appCtx.DB().
		RawQuery(
			`SELECT * FROM get_counseling_queue($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			officeUserGbloc,
			counselingQueueParams.CustomerName,
			counselingQueueParams.Edipi,
			counselingQueueParams.Emplid,
			pq.Array(counselingQueueParams.Status),
			counselingQueueParams.Locator,
			counselingQueueParams.RequestedMoveDate,
			counselingQueueParams.SubmittedAt,
			counselingQueueParams.Branch,
			counselingQueueParams.OriginDutyLocationName,
			counselingQueueParams.CounselingOffice,
			counselingQueueParams.SCAssignedUser,
			counselingQueueParams.HasSafetyPrivilege,
			counselingQueueParams.Page,
			counselingQueueParams.PerPage,
			counselingQueueParams.Sort,
			counselingQueueParams.Order,
		).
		All(&movesWithCount)

	if err != nil {
		appCtx.Logger().
			Error("error fetching moves for counseling queue", zap.Error(err))
		return movesWithCount, err
	}

	return movesWithCount, nil
}

type MoveWithCount struct {
	models.Move
	OrdersRaw                     json.RawMessage              `json:"orders" db:"orders"`
	Orders                        *models.Order                `json:"-"`
	MTOShipmentsRaw               json.RawMessage              `json:"mto_shipments" db:"mto_shipments"`
	MTOShipments                  *models.MTOShipments         `json:"-"`
	CounselingOfficeRaw           json.RawMessage              `json:"counseling_transportation_office" db:"counseling_transportation_office"`
	CounselingOffice              *models.TransportationOffice `json:"-"`
	SCAssignedUserRaw             json.RawMessage              `json:"sc_assigned" db:"sc_assigned"`
	SCCounselingAssignedUser      *models.OfficeUser           `json:"-"`
	TotalCount                    int64                        `json:"total_count" db:"total_count"`
	EarliestRequestedPickupDate   *time.Time                   `json:"mtos_earliest_requested_pickup_date" db:"mtos_earliest_requested_pickup_date"`
	EarliestRequestedDeliveryDate *time.Time                   `json:"mtos_earliest_requested_delivery_date" db:"mtos_earliest_requested_delivery_date"`
	EarliestExpectedDepartureDate *time.Time                   `json:"ppms_earliest_expected_departure_date" db:"ppms_earliest_expected_departure_date"`
	PPMShipmentsRaw               json.RawMessage              `json:"ppm_shipments" db:"ppm_shipments"`
	PPMShipments                  *models.PPMShipments         `json:"-"`
}

func movesWithCountToMoves(movesWithCount []MoveWithCount) ([]models.Move, error) {
	var moves models.Moves

	// we have to manually loop through each move and populate the nested objects that the queue uses/needs
	for i := range movesWithCount {
		// populating Move.Orders struct
		var order models.Order
		if err := json.Unmarshal(movesWithCount[i].OrdersRaw, &order); err != nil {
			return moves, fmt.Errorf("error unmarshaling orders JSON: %w", err)
		}
		movesWithCount[i].OrdersRaw = nil
		movesWithCount[i].Orders = &order

		// populating Move.MTOShipments array
		var shipments models.MTOShipments
		if err := json.Unmarshal(movesWithCount[i].MTOShipmentsRaw, &shipments); err != nil {
			return moves, fmt.Errorf("error unmarshaling shipments JSON: %w", err)
		}
		movesWithCount[i].MTOShipmentsRaw = nil
		movesWithCount[i].MTOShipments = &shipments

		// populating Move.PPMShipments array
		var ppmShipments models.PPMShipments
		if err := json.Unmarshal(movesWithCount[i].PPMShipmentsRaw, &ppmShipments); err != nil {
			return moves, fmt.Errorf("error unmarshaling ppmShipments JSON: %w", err)
		}
		movesWithCount[i].PPMShipmentsRaw = nil
		movesWithCount[i].PPMShipments = &ppmShipments

		for i, shipment := range shipments {
			for _, ppmShipment := range ppmShipments {
				if ppmShipment.ShipmentID == shipment.ID {
					shipments[i].PPMShipment = &ppmShipment
				}
			}
		}

		// populating Moves.CounselingOffice struct
		var counselingTransportationOffice models.TransportationOffice
		if err := json.Unmarshal(movesWithCount[i].CounselingOfficeRaw, &counselingTransportationOffice); err != nil {
			return moves, fmt.Errorf("error unmarshaling counseling_transportation_office JSON: %w", err)
		}
		movesWithCount[i].CounselingOfficeRaw = nil
		movesWithCount[i].CounselingOffice = &counselingTransportationOffice

		var scAssigned models.OfficeUser
		if err := json.Unmarshal(movesWithCount[i].SCAssignedUserRaw, &scAssigned); err != nil {
			return moves, fmt.Errorf("error unmarshaling sc_assigned JSON: %w", err)
		}
		movesWithCount[i].SCAssignedUserRaw = nil
		movesWithCount[i].SCCounselingAssignedUser = &scAssigned
	}

	// the handler consumes a Move object and NOT the MoveWithCount struct used in this func
	// so we have to copy our custom struct into the Move struct
	for _, moveWithCount := range movesWithCount {
		var move models.Move
		if err := copier.Copy(&move, &moveWithCount); err != nil {
			return moves, fmt.Errorf("error copying movesWithCount into Moves struct: %w", err)
		}
		moves = append(moves, move)
	}

	return moves, nil
}
