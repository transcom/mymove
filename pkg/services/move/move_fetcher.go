package move

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveFetcher struct {
}

// NewMoveFetcher creates a new moveFetcher service
func NewMoveFetcher() services.MoveFetcher {
	return &moveFetcher{}
}

// FetchMove retrieves a Move if it is visible for a given locator
func (f moveFetcher) FetchMove(appCtx appcontext.AppContext, locator string, searchParams *services.MoveFetcherParams) (*models.Move, error) {
	move := &models.Move{}
	query := appCtx.DB().
		EagerPreload("CloseoutOffice.Address", "Contractor", "ShipmentGBLOC", "LockedByOfficeUser", "LockedByOfficeUser.TransportationOffice", "AdditionalDocuments",
			"AdditionalDocuments.UserUploads").
		LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
		LeftJoin("office_users", "office_users.id = moves.locked_by").
		Where("locator = $1", locator)

	if searchParams == nil || !searchParams.IncludeHidden {
		query.Where("show = TRUE")
	}

	err := query.First(move)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.Move{}, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.Move{}, apperror.NewQueryError("Move", err, "")
		}
	}

	if move.AdditionalDocumentsID != nil {
		var additionalDocumentUploads models.UserUploads
		err = appCtx.DB().Q().
			Scope(utilities.ExcludeDeletedScope()).EagerPreload("Upload").
			Where("document_id = ?", move.AdditionalDocumentsID).
			All(&additionalDocumentUploads)
		if err != nil {
			return move, err
		}
		move.AdditionalDocuments.UserUploads = additionalDocumentUploads
	}

	return move, nil
}

// Fetches moves for Navy servicemembers with approved shipments. Ignores gbloc rules
func (f moveFetcher) FetchMovesForPPTASReports(appCtx appcontext.AppContext, params *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	var moves models.Moves

	query := appCtx.DB().EagerPreload(
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MTOShipments.Reweigh",
		"MTOShipments.PPMShipment",
		"Orders.ServiceMember",
		"Orders.ServiceMember.ResidentialAddress",
		"Orders.ServiceMember.BackupContacts",
		"Orders.Entitlement",
		"Orders.Entitlement.WeightAllotted",
		"Orders.NewDutyLocation.Address.Country",
		"Orders.NewDutyLocation.TransportationOffice.Gbloc",
		"Orders.OriginDutyLocation.Address.Country",
		"Orders.TAC",
	).
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("entitlements", "entitlements.id = orders.entitlement_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("mto_shipments", "mto_shipments.move_id = moves.id").
		LeftJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
		LeftJoin("addresses", "addresses.id in (mto_shipments.pickup_address_id, mto_shipments.destination_address_id)").
		Where("mto_shipments.status = 'APPROVED'").
		Where("service_members.affiliation = ?", models.AffiliationNAVY).
		GroupBy("moves.id")

	if params.Since != nil {
		query.Where("mto_shipments.updated_at >= ?", params.Since)
	}

	err := query.All(&moves)

	if err != nil {
		return nil, err
	}

	if len(moves) < 1 {
		return nil, nil
	}

	return moves, nil
}

type moveFetcherBulkAssignment struct {
}

// NewMoveFetcherBulkAssignment creates a new moveFetcherBulkAssignment service
func NewMoveFetcherBulkAssignment() services.MoveFetcherBulkAssignment {
	return &moveFetcherBulkAssignment{}
}

func (f moveFetcherBulkAssignment) FetchMovesForBulkAssignmentCounseling(appCtx appcontext.AppContext, gbloc string, officeId uuid.UUID) ([]models.MoveWithEarliestDate, error) {
	var moves []models.MoveWithEarliestDate

	err := appCtx.DB().
		RawQuery(`SELECT
					moves.id,
					MIN(LEAST(
						COALESCE(mto_shipments.requested_pickup_date, '9999-12-31'),
						COALESCE(mto_shipments.requested_delivery_date, '9999-12-31'),
						COALESCE(ppm_shipments.expected_departure_date, '9999-12-31')
					)) AS earliest_date
				FROM moves
				INNER JOIN orders ON orders.id = moves.orders_id
				INNER JOIN mto_shipments ON mto_shipments.move_id = moves.id
				LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
				WHERE
					mto_shipments.deleted_at IS NULL
					AND moves.status = 'NEEDS SERVICE COUNSELING'
					AND orders.gbloc = $1
					AND moves.show = $2
					AND moves.sc_assigned_id IS NULL
					AND moves.counseling_transportation_office_id = $3
					AND (ppm_shipments.status IS NULL OR ppm_shipments.status NOT IN ($4, $5, $6))
					AND (orders.orders_type NOT IN ($7, $8, $9))
				GROUP BY moves.id
				ORDER BY earliest_date ASC`,
			gbloc,
			models.BoolPointer(true),
			officeId,
			models.PPMShipmentStatusWaitingOnCustomer,
			models.PPMShipmentStatusNeedsCloseout,
			models.PPMShipmentStatusCloseoutComplete,
			internalmessages.OrdersTypeBLUEBARK,
			internalmessages.OrdersTypeWOUNDEDWARRIOR,
			internalmessages.OrdersTypeSAFETY).
		All(&moves)

	if err != nil {
		return nil, fmt.Errorf("error fetching moves for office: %s with error %w", officeId, err)
	}

	if len(moves) < 1 {
		return nil, nil
	}

	return moves, nil
}

func (f moveFetcherBulkAssignment) FetchMovesForBulkAssignmentCloseout(appCtx appcontext.AppContext, gbloc string, officeId uuid.UUID) ([]models.MoveWithEarliestDate, error) {
	var moves []models.MoveWithEarliestDate

	query := `SELECT
					moves.id,
					ppm_shipments.submitted_at AS earliest_date
				FROM moves
				INNER JOIN orders ON orders.id = moves.orders_id
				INNER JOIN service_members ON service_members.id = orders.service_member_id
				INNER JOIN mto_shipments ON mto_shipments.move_id = moves.id
				INNER JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
				WHERE
					(moves.status IN ('APPROVED', 'SERVICE COUNSELING COMPLETED'))
					AND moves.show = $1
					AND moves.sc_assigned_id IS NULL`

	switch gbloc {
	case "NAVY":
		query += ` AND (service_members.affiliation in ('` + string(models.AffiliationNAVY) + `'))`
	case "TVCB":
		query += ` AND (service_members.affiliation in ('` + string(models.AffiliationMARINES) + `'))`
	case "USCG":
		query += ` AND (service_members.affiliation in ('` + string(models.AffiliationCOASTGUARD) + `'))`
	default:
		query += ` AND moves.closeout_office_id = '` + officeId.String() + `'
				   AND (service_members.affiliation NOT IN ('` +
			string(models.AffiliationNAVY) + `', '` +
			string(models.AffiliationMARINES) + `', '` +
			string(models.AffiliationCOASTGUARD) + `'))`
	}

	query += ` AND (ppm_shipments.status IN ($2))
					AND (orders.orders_type NOT IN ($3, $4, $5))
				GROUP BY moves.id, ppm_shipments.submitted_at
				ORDER BY earliest_date ASC`

	err := appCtx.DB().RawQuery(query,
		models.BoolPointer(true),
		models.PPMShipmentStatusNeedsCloseout,
		internalmessages.OrdersTypeBLUEBARK,
		internalmessages.OrdersTypeWOUNDEDWARRIOR,
		internalmessages.OrdersTypeSAFETY).All(&moves)

	if err != nil {
		return nil, fmt.Errorf("error fetching moves for office: %s with error %w", officeId, err)
	}

	if len(moves) < 1 {
		return nil, nil
	}

	return moves, nil
}

func (f moveFetcherBulkAssignment) FetchMovesForBulkAssignmentTaskOrder(appCtx appcontext.AppContext, gbloc string, officeId uuid.UUID) ([]models.MoveWithEarliestDate, error) {
	var moves []models.MoveWithEarliestDate

	sqlQuery := `
		SELECT
			moves.id,
			MIN(LEAST(
				COALESCE(mto_shipments.requested_pickup_date, '9999-12-31'),
				COALESCE(mto_shipments.requested_delivery_date, '9999-12-31'),
				COALESCE(ppm_shipments.expected_departure_date, '9999-12-31')
			)) AS earliest_date
		FROM moves
		INNER JOIN orders ON orders.id = moves.orders_id
		INNER JOIN service_members ON orders.service_member_id = service_members.id
		INNER JOIN mto_shipments ON mto_shipments.move_id = moves.id
		INNER JOIN duty_locations as origin_dl ON orders.origin_duty_location_id = origin_dl.id
		LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
		LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
		WHERE
			mto_shipments.deleted_at IS NULL
			AND (moves.status IN ('APPROVALS REQUESTED', 'SUBMITTED', 'SERVICE COUNSELING COMPLETED'))
			AND moves.show = $1
			AND moves.too_assigned_id IS NULL
			AND (orders.orders_type NOT IN ($2, $3, $4))
			AND (moves.ppm_type IS NULL OR (moves.ppm_type = 'PARTIAL' or (moves.ppm_type = 'FULL' and origin_dl.provides_services_counseling = 'false'))) `
	if gbloc == "USMC" {
		sqlQuery += `
			AND service_members.affiliation ILIKE 'MARINES' `
	} else {
		sqlQuery += `
		AND service_members.affiliation != 'MARINES'
		AND ((mto_shipments.shipment_type != 'HHG_OUTOF_NTS' AND move_to_gbloc.gbloc = '` + gbloc + `')
		OR (mto_shipments.shipment_type = 'HHG_OUTOF_NTS' AND orders.gbloc = '` + gbloc + `')) `
	}
	sqlQuery += `
		GROUP BY moves.id
		ORDER BY earliest_date ASC`

	err := appCtx.DB().RawQuery(sqlQuery,
		models.BoolPointer(true),
		internalmessages.OrdersTypeBLUEBARK,
		internalmessages.OrdersTypeWOUNDEDWARRIOR,
		internalmessages.OrdersTypeSAFETY).
		All(&moves)

	if err != nil {
		return nil, fmt.Errorf("error fetching moves for office: %s with error %w", officeId, err)
	}

	if len(moves) < 1 {
		return nil, nil
	}

	return moves, nil
}

func (f moveFetcherBulkAssignment) FetchMovesForBulkAssignmentPaymentRequest(appCtx appcontext.AppContext, gbloc string, officeId uuid.UUID) ([]models.MoveWithEarliestDate, error) {
	var moves []models.MoveWithEarliestDate

	sqlQuery := `
		SELECT
			moves.id,
			min(payment_requests.requested_at) AS earliest_date
		FROM moves
		INNER JOIN orders ON orders.id = moves.orders_id
		INNER JOIN service_members ON orders.service_member_id = service_members.id
		INNER JOIN payment_requests on payment_requests.move_id = moves.id
		LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
		WHERE moves.show = $1
		AND (orders.orders_type NOT IN ($2, $3, $4))
		AND moves.tio_assigned_id IS NULL
		AND payment_requests.status = 'PENDING' `
	if gbloc == "USMC" {
		sqlQuery += `
			AND service_members.affiliation ILIKE 'MARINES' `
	} else {
		sqlQuery += `
		AND service_members.affiliation != 'MARINES'
		AND move_to_gbloc.gbloc = '` + gbloc + `' `
	}
	sqlQuery += `
		GROUP BY moves.id
        ORDER BY earliest_date ASC`

	err := appCtx.DB().RawQuery(sqlQuery,
		models.BoolPointer(true),
		internalmessages.OrdersTypeBLUEBARK,
		internalmessages.OrdersTypeWOUNDEDWARRIOR,
		internalmessages.OrdersTypeSAFETY).
		All(&moves)

	if err != nil {
		return nil, fmt.Errorf("error fetching moves for office: %s with error %w", officeId, err)
	}

	if len(moves) < 1 {
		return nil, nil
	}

	return moves, nil
}
