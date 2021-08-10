package fetch

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type fetcherQueryBuilder interface {
	FetchOne(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error
}

type fetcher struct {
	builder fetcherQueryBuilder
}

// FetchRecord uses the passed query builder to fetch a record
func (o *fetcher) FetchRecord(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(query.FetchOneReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(query.FetchOneReflectionMessage)
	}

	err := o.builder.FetchOne(appCfg, model, filters)

	elem := reflect.ValueOf(model).Elem()
	id := elem.FieldByName("ID").Interface().(uuid.UUID)
	if id == uuid.Nil {
		err = fmt.Errorf("Resource not found: %w", err)
	}

	return err
}

// NewFetcher returns an implementation of ListFetcher
func NewFetcher(builder fetcherQueryBuilder) services.Fetcher {
	return &fetcher{builder}
}
