package move

import (
	"strings"

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

func (s moveSearcher) SearchMoves(appCtx appcontext.AppContext, locator *string, dodID *string, customerName *string) (models.Moves, error) {
	if locator == nil && dodID == nil && customerName == nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "move locator, DOD ID, or customer name must be provided")
		return models.Moves{}, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}
	if locator != nil && dodID != nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "search by multiple keys is not supported")
		return models.Moves{}, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}

	query := appCtx.DB().EagerPreload(
		"MTOShipments",
		"Orders.ServiceMember",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
	).
		Join("orders", "orders.id = moves.orders_id").
		Join("service_members", "service_members.id = orders.service_member_id").
		Where("show = TRUE")

	if customerName != nil && len(*customerName) > 0 {
		query = query.Where("? % translate(lower(first_name || ' ' || last_name),'àáâãäåèéêëìíîïòóôõöùúûüñ', 'aaaaaaeeeeiiiiooooouuuun')", *customerName).Order("similarity(translate(lower(first_name || ' ' || last_name),'àáâãäåèéêëìíîïòóôõöùúûüñ', 'aaaaaaeeeeiiiiooooouuuun'), ?) desc", *customerName)
	}

	if locator != nil {
		searchLocator := strings.ToUpper(*locator)
		query = query.Where("locator = ?", searchLocator)
	}

	if dodID != nil {
		query = query.Where("service_members.edipi = ?", *dodID)
	}

	var moves models.Moves
	err := query.All(&moves)

	if err != nil {
		return models.Moves{}, apperror.NewQueryError("Move", err, "")
	}
	return moves, nil
}
