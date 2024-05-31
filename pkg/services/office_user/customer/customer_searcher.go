package customer

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type customerSearcher struct {
}

func NewCustomerSearcher() services.CustomerSearcher {
	return &customerSearcher{}
}

type QueryOption func(*pop.Query)

func (s customerSearcher) SearchCustomers(appCtx appcontext.AppContext, params *services.SearchCustomersParams) (models.ServiceMembers, int, error) {
	if params.DodID == nil && params.CustomerName == nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "DOD ID or customer name must be provided")
		return models.ServiceMembers{}, 0, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}

	if params.CustomerName != nil && params.DodID != nil {
		verrs := validate.NewErrors()
		verrs.Add("search key", "search by multiple keys is not supported")
		return models.ServiceMembers{}, 0, apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	}

	err := appCtx.DB().RawQuery("SET pg_trgm.similarity_threshold = 0.1").Exec()
	if err != nil {
		return nil, 0, err
	}

	privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if err != nil {
		appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
	}

	var query *pop.Query

	if appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
		rawquery := `SELECT DISTINCT ON (id)
			service_members.affiliation, service_members.backup_mailing_address_id,
			service_members.cac_validated, service_members.created_at, service_members.edipi,
			service_members.email_is_preferred, service_members.emplid,
			service_members.first_name, service_members.id, service_members.last_name,
			service_members.middle_name, service_members.personal_email,
			service_members.phone_is_preferred, service_members.residential_address_id,
			service_members.secondary_telephone, service_members.suffix,
			service_members.telephone, service_members.updated_at, service_members.user_id
		FROM service_members AS service_members
			JOIN users ON users.id = service_members.user_id
			LEFT JOIN orders ON orders.service_member_id = service_members.id`

		if !privileges.HasPrivilege(models.PrivilegeTypeSafety) {
			rawquery += ` WHERE ((orders.orders_type != 'SAFETY' or orders.orders_type IS NULL) AND`
		} else {
			rawquery += ` WHERE (`
		}

		if params.DodID != nil {
			rawquery += ` service_members.edipi = $1)`
			query = appCtx.DB().RawQuery(rawquery, params.DodID)
		} else {
			rawquery += ` f_unaccent(lower($1)) % searchable_full_name(first_name, last_name))`
			query = appCtx.DB().RawQuery(rawquery, params.CustomerName)
		}
	}

	customerNameQuery := customerNameSearch(params.CustomerName)
	dodIDQuery := dodIDSearch(params.DodID)
	orderQuery := sortOrder(params.Sort, params.Order)

	options := [3]QueryOption{customerNameQuery, dodIDQuery, orderQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	var customers models.ServiceMembers
	err = query.Paginate(int(params.Page), int(params.PerPage)).All(&customers)

	if err != nil {
		return models.ServiceMembers{}, 0, apperror.NewQueryError("Customer", err, "")
	}
	return customers, query.Paginator.TotalEntriesSize, nil
}

func dodIDSearch(dodID *string) QueryOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query.Where("service_members.edipi = ?", dodID)
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

var parameters = map[string]string{
	"customerName":  "service_members.last_name",
	"dodID":         "service_members.edipi",
	"branch":        "service_members.affiliation",
	"personalEmail": "service_members.personal_email",
	"telephone":     "service_members.telephone",
}

func sortOrder(sort *string, order *string) QueryOption {
	return func(query *pop.Query) {
		if sort != nil && order != nil {
			sortTerm := parameters[*sort]
			query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
		} else {
			query.Order("service_members.last_name ASC")
		}
	}
}
