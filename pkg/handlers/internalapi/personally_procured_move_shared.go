package internalapi

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

// GetDestinationDutyStationPostalCode returns the postal code associated with orders->new_duty_station->address
func GetDestinationDutyStationPostalCode(db *pop.Connection, ordersID uuid.UUID) (string, error) {
	queryBuilder := query.NewQueryBuilder(db)

	var orders models.Orders
	filters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", ordersID),
	}
	associations := query.NewQueryAssociations([]services.QueryAssociation{
		query.NewQueryAssociation("NewDutyStation.Address"),
	})
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	pagination := pagination.NewPagination(&page, &perPage)
	ordering := query.NewQueryOrder(nil, nil)

	err := queryBuilder.FetchMany(&orders, filters, associations, pagination, ordering)
	if err != nil {
		return "", err
	}

	if len(orders) == 0 {
		return "", models.ErrFetchNotFound
	}

	return orders[0].NewDutyStation.Address.PostalCode, nil
}
