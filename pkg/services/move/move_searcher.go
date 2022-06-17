package move

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveSearcher struct {
}

func NewMoveSearcher() services.MoveSearcher {
	return &moveSearcher{}
}

func (s moveSearcher) SearchMoves(appCtx appcontext.AppContext, params *services.SearchMovesParams) (models.Moves, int, error) {
	appCtx.Logger().Warn("🐢 SearchMoves", zap.Any("params", params))
	if params.Locator == nil && params.DodID == nil && params.CustomerName == nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "move locator, DOD ID, or customer name must be provided")
		return models.Moves{}, 0, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}
	if params.Locator != nil && params.DodID != nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "search by multiple keys is not supported")
		return models.Moves{}, 0, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}
	//query := appCtx.DB().Select("moves.locator", "service_members.first_name", "service_members.last_name").

	query := appCtx.DB().EagerPreload(
		"MTOShipments",
		"Orders.ServiceMember",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
	).
		Join("orders", "orders.id = moves.orders_id").
		Join("service_members", "service_members.id = orders.service_member_id").
		Join("duty_locations as origin_duty_locations", "origin_duty_locations.id = orders.origin_duty_location_id").
		Join("addresses as origin_addresses", "origin_addresses.id = origin_duty_locations.address_id").
		Join("duty_locations as new_duty_locations", "new_duty_locations.id = orders.new_duty_location_id").
		Join("addresses as new_addresses", "new_addresses.id = new_duty_locations.address_id").
		Join("mto_shipments", "mto_shipments.move_id = moves.id").
		GroupBy("moves.id", "service_members.id", "origin_addresses.id", "new_addresses.id").
		Where("show = TRUE")

	if params.CustomerName != nil && len(*params.CustomerName) > 0 {
		query = query.Where("f_unaccent(lower(?)) % searchable_full_name(first_name, last_name)", *params.CustomerName)

		if params.Sort == nil || params.Order == nil {
			query = query.Order("similarity(searchable_full_name(first_name, last_name), f_unaccent(lower(?))) desc", *params.CustomerName)
		}
	}

	if params.Locator != nil {
		searchLocator := strings.ToUpper(*params.Locator)
		query = query.Where("locator = ?", searchLocator)
	}

	if params.DodID != nil {
		query = query.Where("service_members.edipi = ?", *params.DodID)
	}

	if params.OriginPostalCode != nil {
		query = query.Where("origin_addresses.postal_code = ?", *params.OriginPostalCode)
	}
	if params.DestinationPostalCode != nil {
		query = query.Where("new_addresses.postal_code = ?", *params.DestinationPostalCode)
	}
	if params.Status != nil && len(params.Status) > 0 {
		query = query.Where("moves.status in (?)", params.Status)
	}
	if params.Branch != nil {
		query = query.Where("service_members.affiliation = ?", params.Branch)
	}
	if params.ShipmentsCount != nil {
		query = query.Having("COUNT(mto_shipments.id) = ?", *params.ShipmentsCount)
	}
	if params.Sort != nil && params.Order != nil {
		appCtx.Logger().Warn("Ordering!!!", zap.String("col", qualifySortColumn(*params.Sort)), zap.String("ord", *params.Order))
		query = query.Order(fmt.Sprintf("%s %s", qualifySortColumn(*params.Sort), *params.Order))
	}

	var moves models.Moves
	err := query.Paginate(int(params.Page), int(params.PerPage)).All(&moves)

	if err != nil {
		return models.Moves{}, 0, apperror.NewQueryError("Move", err, "")
	}
	return moves, query.Paginator.TotalEntriesSize, nil
}

// TODO rename
// TODO and have this modify the query
func qualifySortColumn(sort string) string {
	if sort == "customerName" {
		return "service_members.last_name"
	}
	if sort == "status" {
		return "moves.status"
	}
	if sort == "originPostalCode" {
		return "origin_addresses.postal_code"
	}
	if sort == "destinationPostalCode" {
		return "new_addresses.postal_code"
	}
	if sort == "branch" {
		return "service_members.affiliation"
	}
	if sort == "shipmentsCount" {
		return "COUNT(mto_shipments.id)"
	}

	return "moves.locator" // TODO what do we do in the default case?
}
