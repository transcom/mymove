package paymentrequest

import (
	"database/sql"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/models"
)

type paymentRequestListFetcher struct {
	db *pop.Connection
}

// NewPaymentRequestListFetcher returns a new payment request list fetcher
func NewPaymentRequestListFetcher(db *pop.Connection) services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{db}
}

func (f *paymentRequestListFetcher) FetchPaymentRequestList(officeUserID uuid.UUID, options ...func(query *pop.Query)) (*models.PaymentRequests, error) {
	gblocFetcher := officeuser.NewOfficeUserGblocFetcher(f.db)
	gbloc, gblocErr := gblocFetcher.FetchGblocForOfficeUser(officeUserID)
	if gblocErr != nil {
		return &models.PaymentRequests{}, gblocErr
	}

	paymentRequests := models.PaymentRequests{}
	query := f.db.Q().Eager(
		"MoveTaskOrder.Orders.OriginDutyStation",
		"MoveTaskOrder.Orders.ServiceMember",
	).
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("duty_stations", "duty_stations.id = orders.origin_duty_station_id").
		InnerJoin("transportation_offices", "transportation_offices.id = duty_stations.transportation_office_id").
		Where("transportation_offices.gbloc = ?", gbloc)

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	err := query.GroupBy("payment_requests.id").All(&paymentRequests)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, services.NotFoundError{}
		default:
			return nil, err
		}
	}

	if err != nil {
		return &models.PaymentRequests{}, err
	}

	for i := range paymentRequests {
		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		// cannot eager load the address as "OriginDutyStation.Address" because
		// OriginDutyStation is a pointer.
		if originDutyStation := paymentRequests[i].MoveTaskOrder.Orders.OriginDutyStation; originDutyStation != nil {
			f.db.Load(originDutyStation, "TransportationOffice")
		}
	}

	return &paymentRequests, nil
}
