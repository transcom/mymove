package closeout

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type serviceMember struct {
	FirstName     string
	LastName      string
	MiddleName    string
	Rank          string
	DodID         string
	Affiliation   string
	PersonalEmail string
	Telephone     string
}

// NewCloseoutData creates a new struct
func NewCloseoutData(db *pop.Connection, logger *zap.Logger) (*Data, error) {
	return &Data{
		db:     db,
		logger: logger,
	}, nil
}

// Data is a service object to add missing fuel prices to db
type Data struct {
	db     *pop.Connection
	logger *zap.Logger
}

func buildJSON(sm models.ServiceMember) error {
	data := serviceMember{
		FirstName:     handlers.DerefStringTypes(sm.FirstName),
		LastName:      handlers.DerefStringTypes(sm.LastName),
		MiddleName:    handlers.DerefStringTypes(sm.MiddleName),
		Rank:          handlers.DerefStringTypes(sm.Rank),
		DodID:         handlers.DerefStringTypes(sm.Edipi),
		Affiliation:   handlers.DerefStringTypes(sm.Affiliation),
		PersonalEmail: handlers.DerefStringTypes(sm.PersonalEmail),
		Telephone:     handlers.DerefStringTypes(sm.Telephone),
	}

	file, _ := json.MarshalIndent(data, "", " ")
	return ioutil.WriteFile("pptas.json", file, 0644)
}

// FetchCloseoutDetails testing
func (d Data) FetchCloseoutDetails(moveIDs []string) error {
	fmt.Println("abc 123")
	for i, m := range moveIDs {
		fmt.Println(i, m)
		move, err := d.fetchMove(m)
		if err != nil {
			return err
		}
		fmt.Println(move.ID)

		order, err := d.fetchOrder(move.OrdersID)
		if err != nil {
			return err
		}
		fmt.Println(order.ServiceMemberID)

		sm, err := d.fetchServiceMember(order.ServiceMemberID)
		if err != nil {
			return err
		}
		err = buildJSON(sm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Data) fetchMove(moveLocator string) (models.Move, error) {
	var move models.Move

	builder := query.NewQueryBuilder(d.db)
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("locator", "=", moveLocator),
	}

	err := builder.FetchOne(&move, queryFilters)
	if err != nil {
		d.logger.Error("Move not found")
		return models.Move{}, err
	}

	return move, nil
}

func (d Data) fetchOrder(ordersID uuid.UUID) (models.Order, error) {

	var order models.Order

	builder := query.NewQueryBuilder(d.db)
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", ordersID),
	}

	err := builder.FetchOne(&order, queryFilters)
	if err != nil {
		d.logger.Error("Order not found")
		return models.Order{}, err
	}

	return order, nil
}

func (d Data) fetchServiceMember(serviceMemberID uuid.UUID) (models.ServiceMember, error) {

	var sm models.ServiceMember

	builder := query.NewQueryBuilder(d.db)
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", serviceMemberID),
	}

	err := builder.FetchOne(&sm, queryFilters)
	if err != nil {
		d.logger.Error("ServiceMember not found")
		return models.ServiceMember{}, err
	}

	return sm, nil
}
