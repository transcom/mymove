package move

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"
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

// QueryOption defines the type for the functional arguments passed to SearchMoves
type QueryOption func(*pop.Query)

// SearchMoves returns a list of results for a QAE/CSR move search query
func (s moveSearcher) SearchMoves(appCtx appcontext.AppContext, params *services.SearchMovesParams) (models.Moves, int, error) {
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

	// The SQL % operator filters out strings that are below this similarity threshold
	// We have to set it here because other areas of the code that do a trigram search
	// (eg Duty Location search) may set a different threshold.
	// If the threshold is too high, we may filter out too many results and make searching harder.
	// If it's too low, the query will get slower/more memory intensive.
	err := appCtx.DB().RawQuery("SET pg_trgm.similarity_threshold = 0.1").Exec()
	if err != nil {
		return nil, 0, err
	}

	query := appCtx.DB().EagerPreload(
		"MTOShipments",
		"Orders.ServiceMember",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
	).
		Join("orders", "orders.id = moves.orders_id").
		Join("service_members", "service_members.id = orders.service_member_id").
		LeftJoin("duty_locations as origin_duty_locations", "origin_duty_locations.id = orders.origin_duty_location_id").
		Join("addresses as origin_addresses", "origin_addresses.id = origin_duty_locations.address_id").
		Join("duty_locations as new_duty_locations", "new_duty_locations.id = orders.new_duty_location_id").
		Join("addresses as new_addresses", "new_addresses.id = new_duty_locations.address_id").
		LeftJoin("mto_shipments", "mto_shipments.move_id = moves.id AND mto_shipments.status <> 'DRAFT'").
		GroupBy("moves.id", "service_members.id", "origin_addresses.id", "new_addresses.id").
		Where("show = TRUE")

	customerNameQuery := customerNameSearch(params.CustomerName)
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.DodID)
	branchQuery := branchFilter(params.Branch)
	originPostalCodeQuery := originPostalCodeFilter(params.OriginPostalCode)
	destinationPostalCodeQuery := destinationPostalCodeFilter(params.DestinationPostalCode)
	statusQuery := moveStatusFilter(params.Status)
	shipmentsCountQuery := shipmentsCountFilter(params.ShipmentsCount)
	orderQuery := sortOrder(params.Sort, params.Order, params.CustomerName)

	options := [10]QueryOption{customerNameQuery, locatorQuery, dodIDQuery, branchQuery, orderQuery, originPostalCodeQuery, destinationPostalCodeQuery, statusQuery, shipmentsCountQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	var moves models.Moves
	err = query.Paginate(int(params.Page), int(params.PerPage)).All(&moves)

	if err != nil {
		return models.Moves{}, 0, apperror.NewQueryError("Move", err, "")
	}
	return moves, query.Paginator.TotalEntriesSize, nil
}

var parameters = map[string]string{
	"customerName":          "service_members.last_name",
	"dodID":                 "service_members.edipi",
	"branch":                "service_members.affiliation",
	"locator":               "moves.locator",
	"status":                "moves.status",
	"originPostalCode":      "origin_addresses.postal_code",
	"destinationPostalCode": "new_addresses.postal_code",
	"shipmentsCount":        "COUNT(mto_shipments.id)",
}

func dodIDFilter(dodID *string) QueryOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func locatorFilter(locator *string) QueryOption {
	return func(query *pop.Query) {
		if locator != nil {
			query.Where("moves.locator = ?", *locator)
		}
	}
}

func branchFilter(branch *string) QueryOption {
	return func(query *pop.Query) {
		if branch != nil {
			query.Where("service_members.affiliation = ?", *branch)
		}
	}
}
func originPostalCodeFilter(postalCode *string) QueryOption {
	return func(query *pop.Query) {
		if postalCode != nil {
			query.Where("origin_addresses.postal_code = ?", *postalCode)
		}
	}
}
func destinationPostalCodeFilter(postalCode *string) QueryOption {
	return func(query *pop.Query) {
		if postalCode != nil {
			query.Where("new_addresses.postal_code = ?", *postalCode)
		}
	}
}

func moveStatusFilter(statuses []string) QueryOption {
	return func(query *pop.Query) {
		if len(statuses) > 0 {
			var translatedStatuses []string
			for _, status := range statuses {
				if strings.EqualFold(status, string(models.MoveStatusSUBMITTED)) {
					translatedStatuses = append(translatedStatuses, string(models.MoveStatusSUBMITTED), string(models.MoveStatusServiceCounselingCompleted))
				} else {
					translatedStatuses = append(translatedStatuses, status)
				}
			}
			query.Where("moves.status in (?)", translatedStatuses)
		}
	}
}

func customerNameSearch(customerName *string) QueryOption {
	return func(query *pop.Query) {
		if customerName != nil && len(*customerName) > 0 {
			query.Where("f_unaccent(lower(?)) % searchable_full_name(first_name, last_name)", *customerName)
		}
	}
}

func shipmentsCountFilter(shipmentsCount *int64) QueryOption {
	return func(query *pop.Query) {
		if shipmentsCount != nil {
			query.Having("COUNT(mto_shipments.id) = ?", *shipmentsCount)
		}
	}
}

func sortOrder(sort *string, order *string, customerNameSearch *string) QueryOption {
	return func(query *pop.Query) {
		if sort != nil && order != nil {
			sortTerm := parameters[*sort]
			if *sort == "customerName" {
				orderName(query, order)
			} else {
				query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
			}
		} else if customerNameSearch != nil {
			query.Order("similarity(searchable_full_name(first_name, last_name), f_unaccent(lower(?))) DESC", *customerNameSearch)
		} else {
			query.Order("moves.created_at DESC")
		}
	}
}

func orderName(query *pop.Query, order *string) *pop.Query {
	query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
	return query
}
