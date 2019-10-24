package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/services"
)

// instanceOfBuilder.FetchOne(model, filters, associations, pagination, ...parameters)
// instanceOfBuilder.FetchOne(model, id).WithFilters(filtes)
// query.NewFetchMany(model interface{}).WithFilters(filters).WithPagination(pagination).WithAssociations(associations).Execute()

type FetchMany struct {
	DB      *pop.Connection
	Model   interface{}
	Filters []services.QueryFilter
}

func NewFetchMany(model interface{}) *FetchMany {
	return &FetchMany{
		Model: &model,
	}
}

func (f *FetchMany) WithFilters(filters []services.QueryFilter) *FetchMany {
	f.Filters = filters
	return f
}

func (f *FetchMany) Execute() error {
	query := f.DB.Q()
	t := reflect.TypeOf(f.Model)

	if len(f.Filters) > 0 {
		filteredQuery(query, f.Filters, t)
	}

	return query.All(f.Model)
}

// allowed comparators for this query builder implementation
const equals = "="
const greaterThan = ">"

// Error message constants
const fetchManyReflectionMessage = "Model should be pointer to slice of structs"
const fetchOneReflectionMessage = "Model should be pointer to struct"

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
	default:
		return "", false
	}
}

func buildQuery(query *pop.Query, filters []services.QueryFilter, pagination services.Pagination, t reflect.Type) (*pop.Query, error) {
	query, err := filteredQuery(query, filters, t)

	if err != nil {
		return nil, err
	}

	query, err = paginatedQuery(query, pagination, t)

	if err != nil {
		return nil, err
	}

	return query, nil
}

func paginatedQuery(query *pop.Query, pagination services.Pagination, t reflect.Type) (*pop.Query, error) {
	return query.Paginate(pagination.Page(), pagination.PerPage()), nil
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
	invalidFields := make([]string, 0)
	for _, f := range filters {
		// Validate the filter we're using is valid/safe
		invalidField := validateFilter(f, t)
		if invalidField != "" {
			invalidFields = append(invalidFields, fmt.Sprintf("%s %s", f.Column(), f.Comparator()))
		}
		// Column lookup should always adhere to SQL injection input validations
		// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.md#defense-option-3-whitelist-input-validation
		columnQuery := fmt.Sprintf("%s %s ?", f.Column(), f.Comparator())
		query = query.Where(columnQuery, f.Value())
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
		return errors.New(fetchOneReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(fetchOneReflectionMessage)
	}
	query := p.db.Q()
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return err
	}
	return query.First(model)
}

// FetchMany fetches multiple model records using pop's All method
// Will return error if model is not pointer to slice of structs
func (p *Builder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination) error {
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
	query, err := buildQuery(query, filters, pagination, t)
	if err != nil {
		return err
	}

	err = associatedQuery(query, associations, model)
	if err != nil {
		return err
	}
	return nil
}

func (p *Builder) CreateOne(model interface{}) (*validate.Errors, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New(fetchOneReflectionMessage)
	}

	verrs, err := p.db.ValidateAndCreate(model)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}
	return nil, nil
}

func (p *Builder) UpdateOne(model interface{}) (*validate.Errors, error) {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New(fetchOneReflectionMessage)
	}

	verrs, err := p.db.ValidateAndUpdate(model)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return nil, nil
}

func (p *Builder) FetchCategoricalCountsFromOneModel(model interface{}, filters []services.QueryFilter, andFilters *[]services.QueryFilter) (map[interface{}]int, error) {
	conn := p.db
	t := reflect.TypeOf(model)
	categoricalCounts, err := categoricalCountsQueryOneModel(conn, filters, andFilters, t)
	if err != nil {
		return nil, err
	}
	return categoricalCounts, nil
}

func (p *Builder) QueryForAssociations(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(fetchOneReflectionMessage)
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
	query, err := buildQuery(query, filters, pagination, t)
	if err != nil {
		return err
	}

	err = associatedQuery(query, associations, model)
	if err != nil {
		return err
	}
	return nil
}

func associatedQuery(query *pop.Query, associations services.QueryAssociations, model interface{}) error {
	query = query.Eager(associations.StringGetAssociations()...)
	return query.All(model)
}
