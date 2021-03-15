package query

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"

	"github.com/transcom/mymove/pkg/services"
)

// allowed comparators for this query builder implementation
const equals = "="
const greaterThan = ">"
const ilike = "ILIKE" // Case insensitive

// allowed sorting order for this query builder implmentation
const asc = "asc"
const desc = "desc"

// Error message constants
const fetchManyReflectionMessage = "Model should be pointer to slice of structs"

// FetchOneReflectionMessage error message for fetching one struct by reflection
const FetchOneReflectionMessage = "Model should be pointer to struct"

// Builder is a wrapper around pop
// with more flexible query patterns to MilMove
type Builder struct {
	db *pop.Connection
}

// NewQueryBuilder returns a new query builder implemented with pop
// constructor is for Dependency Injection frameworks requiring a function instead of struct
func NewQueryBuilder(db *pop.Connection) *Builder {
	return &Builder{db}
}

// Lookup to check if a specific string is inside the db field tags of the type
func getDBColumn(t reflect.Type, field string) (string, bool) {
	for i := 0; i < t.NumField(); i++ {
		dbTag, ok := t.Field(i).Tag.Lookup("db")
		if ok && dbTag == field {
			return dbTag, true
		}
	}
	return "", false
}

// check that we have a valid comparator
func getComparator(comparator string) (string, bool) {
	switch comparator {
	case equals:
		return equals, true
	case greaterThan:
		return greaterThan, true
	case ilike:
		return ilike, true
	default:
		return "", false
	}
}

func buildQuery(query *pop.Query, filters []services.QueryFilter, pagination services.Pagination, order services.QueryOrder, t reflect.Type) (*pop.Query, error) {
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return nil, err
	}

	query, err = paginatedQuery(query, pagination, t)
	if err != nil {
		return nil, err
	}

	query, err = orderedQuery(query, order, t)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func paginatedQuery(query *pop.Query, pagination services.Pagination, t reflect.Type) (*pop.Query, error) {
	if pagination == nil {
		return query, nil
	}

	return query.Paginate(pagination.Page(), pagination.PerPage()), nil
}

func orderedQuery(query *pop.Query, order services.QueryOrder, t reflect.Type) (*pop.Query, error) {
	//omit sorting if no column specified
	if order == nil || order.Column() == nil || order.SortOrder() == nil {
		return query, nil
	}

	// Validate the filter we're using is valid/safe
	invalidField := validateOrder(order, t)

	// Column lookup should always adhere to SQL injection input validations
	// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.md#defense-option-3-whitelist-input-validation
	sortOrder := desc
	if *order.SortOrder() {
		sortOrder = asc
	}

	orderQuery := fmt.Sprintf("%s %s", *order.Column(), sortOrder)
	query = query.Order(orderQuery)

	if len(invalidField) != 0 {
		return query, fmt.Errorf("%v is not valid input", invalidField)
	}

	return query, nil
}

// Validate that the QueryFilter is valid using getDBColumn and getComparator
func validateFilter(f services.QueryFilter, t reflect.Type) string {
	invalidField := ""
	_, ok := getDBColumn(t, f.Column())
	if !ok {
		invalidField = fmt.Sprintf("%s %s", f.Column(), f.Comparator())
	}
	_, ok = getComparator(f.Comparator())
	if !ok {
		invalidField = fmt.Sprintf("%s %s", f.Column(), f.Comparator())
	}
	return invalidField
}

// Validate that the QueryOrder is valid using getDBColumn
func validateOrder(s services.QueryOrder, t reflect.Type) string {
	invalidField := ""
	if s.Column() == nil || s.SortOrder() == nil {
		return "QueryOrder Column or SortOrder is nil"
	}

	_, ok := getDBColumn(t, *s.Column())
	if !ok {
		invalidField = fmt.Sprintf("%s", *s.Column())
	}

	return invalidField
}

// Currently this can select counts for 'categories' based on a field comparison using an array of QueryFilters. Additionally it supports adding in AND logic
// by including a list of AND clauses, also via an array of QueryFilters. TODO: Add in functionality for OR when a use case for it comes up.
func categoricalCountsQueryOneModel(conn *pop.Connection, filters []services.QueryFilter, andFilters *[]services.QueryFilter, t reflect.Type) (map[interface{}]int, error) {
	invalidFields := make([]string, 0)
	counts := make(map[interface{}]int, 0)

	for _, f := range filters {
		// Set up an empty query for us to use to get the count
		query := conn.Q()

		// Validate the filter we're using is valid/safe
		invalidField := validateFilter(f, t)
		if invalidField != "" {
			invalidFields = append(invalidFields, fmt.Sprintf("%s %s", f.Column(), f.Comparator()))
		}

		queryColumn := fmt.Sprintf("%s %s ?", f.Column(), f.Comparator())
		query = query.Where(queryColumn, f.Value())

		if andFilters != nil {
			for _, af := range *andFilters {
				invalidField := validateFilter(af, t)
				if invalidField != "" {
					return nil, fmt.Errorf("%v is not valid input", invalidField)
				}

				queryColumn := fmt.Sprintf("%s %s ?", af.Column(), af.Comparator())
				query = query.Where(queryColumn, af.Value())
			}
		}
		if len(invalidFields) != 0 {
			return nil, fmt.Errorf("%v is not valid input", invalidFields)
		}

		count, err := query.Count(reflect.Zero(t).Interface())
		if err != nil {
			return nil, err
		}
		counts[f.Value()] = count
	}

	return counts, nil
}

func filteredQuery(query *pop.Query, filters []services.QueryFilter, t reflect.Type) (*pop.Query, error) {
	if filters == nil {
		return query, nil
	}

	invalidFields := make([]string, 0)
	likeFilters := []services.QueryFilter{}

	for _, f := range filters {
		// Validate the filter we're using is valid/safe
		invalidField := validateFilter(f, t)
		if invalidField != "" {
			invalidFields = append(invalidFields, fmt.Sprintf("%s %s", f.Column(), f.Comparator()))
		}

		if f.Comparator() == ilike {
			likeFilters = append(likeFilters, f)
			continue
		}

		// Column lookup should always adhere to SQL injection input validations
		// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.md#defense-option-3-whitelist-input-validation
		columnQuery := fmt.Sprintf("%s %s ?", f.Column(), f.Comparator())
		query = query.Where(columnQuery, f.Value())
	}

	// Hacky way to get ILIKE filters to work with OR instead of AND

	if len(likeFilters) > 0 {
		var likeQuery string
		var vals []interface{}
		var queries []string

		for _, f := range likeFilters {
			vals = append(vals, f.Value())
			columnQuery := fmt.Sprintf("%s %s ?", f.Column(), f.Comparator())
			queries = append(queries, columnQuery)

			likeQuery = strings.Join(queries, " OR ")
		}

		query = query.Where(likeQuery, vals...)
	}

	if len(invalidFields) != 0 {
		return query, fmt.Errorf("%v is not valid input", invalidFields)
	}
	return query, nil
}

// FetchOne fetches a single model record using pop's First method
// Will return error if model is not pointer to struct
func (p *Builder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(FetchOneReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(FetchOneReflectionMessage)
	}
	query := p.db.Q()
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return err
	}

	// ToDo: Refactor once we expand the query builder for this function
	// For now, opt-in for eager loading
	queryAssociations := []services.QueryAssociation{}
	associations := NewQueryAssociations(queryAssociations)
	query = associatedQuery(query, associations, model)

	err = query.First(model)
	if err != nil {
		return err
	}

	return nil
}

// FetchMany fetches multiple model records using pop's All method
// Will return error if model is not pointer to slice of structs
func (p *Builder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(fetchManyReflectionMessage)
	}
	query := p.db.Q()
	query, err := buildQuery(query, filters, pagination, ordering, t)
	if err != nil {
		return err
	}

	query = associatedQuery(query, associations, model)
	err = query.All(model)
	if err != nil {
		return err
	}
	return nil
}

// Count returns a count from a filter
func (p *Builder) Count(model interface{}, filters []services.QueryFilter) (int, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return 0, errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return 0, errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return 0, errors.New(fetchManyReflectionMessage)
	}
	query := p.db.Q()
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return 0, err
	}

	count, err := query.Count(model)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateOne creates exactly one model
func (p *Builder) CreateOne(model interface{}) (*validate.Errors, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New(FetchOneReflectionMessage)
	}

	verrs, err := p.db.ValidateAndCreate(model)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}
	return nil, nil
}

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
}

// UpdateOne updates exactly one model
func (p *Builder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New(FetchOneReflectionMessage)
	}

	var verrs *validate.Errors
	var err error

	if eTag != nil {
		err = p.db.Transaction(func(tx *pop.Connection) error {
			t = t.Elem()
			v := reflect.ValueOf(model).Elem()
			var id uuid.UUID
			for i := 0; i < t.NumField(); i++ {
				if t.Field(i).Name == "ID" {
					id = v.Field(i).Interface().(uuid.UUID)
					break
				}
			}

			var tableName string
			tableNameable, ok := model.(pop.TableNameAble)

			if ok {
				tableName = tableNameable.TableName()
			} else {
				tableName = flect.Underscore(flect.Pluralize(t.Name()))
			}

			sqlString := fmt.Sprintf("SELECT updated_at from %s WHERE id = $1 FOR UPDATE", pq.QuoteIdentifier(tableName))
			var updatedAt time.Time
			errExec := tx.RawQuery(sqlString, id.String()).First(&updatedAt)
			if errExec != nil {
				return errExec
			}

			encodedUpdatedAt := etag.GenerateEtag(updatedAt)

			if encodedUpdatedAt != *eTag {
				return StaleIdentifierError{StaleIdentifier: *eTag}
			}

			verrs, err = tx.ValidateAndUpdate(model)

			return nil
		})
	} else {
		verrs, err = p.db.ValidateAndUpdate(model)
	}

	if err != nil {
		return nil, err
	}

	if verrs != nil && verrs.HasAny() {
		return verrs, nil
	}

	return nil, nil
}

// FetchCategoricalCountsFromOneModel returns categorical counts from exactly one model
func (p *Builder) FetchCategoricalCountsFromOneModel(model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	conn := p.db
	t := reflect.TypeOf(model)
	categoricalCounts, err := categoricalCountsQueryOneModel(conn, filters, andFilters, t)
	if err != nil {
		return nil, err
	}
	return categoricalCounts, nil
}

// QueryForAssociations builds a query for associations
func (p *Builder) QueryForAssociations(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(FetchOneReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(fetchManyReflectionMessage)
	}
	query := p.db.Q()
	query, err := buildQuery(query, filters, pagination, ordering, t)
	if err != nil {
		return err
	}

	query = associatedQuery(query, associations, model)
	err = query.All(model)
	if err != nil {
		return err
	}
	return nil
}

func associatedQuery(query *pop.Query, associations services.QueryAssociations, model interface{}) *pop.Query {
	return query.Eager(associations.StringGetAssociations()...)
}

// Transaction will create a new transaction on the connection. Will rollback if
// fn returns error, otherwise, will commit.
// TODO: Will need to revisit the design for this. Use internal txDB state instead?
func (p *Builder) Transaction(fn func(tx *pop.Connection) error) error {
	return p.db.Transaction(fn)
}
