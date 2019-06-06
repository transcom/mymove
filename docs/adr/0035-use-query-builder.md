# Use Query Builder for for Admin Interface

**User Story:**

System Admins (MilMove Engineers)
require flexible querying of MilMove data
to debug and fix production data issues.
For example, if an Engineer is debugging a recent bug with shipments,
they may want to view all shipments created or updated recently.
Maybe the bug seems to be related to a certain set of moves
and an Engineer wants to filter by foreign key.
This type of exploratory querying isn't built into the current APIs
because it's not relevant or authorized to users of the production applications.

Since the current apps don't require this capability,
they also weren't built in a way that can be easily extended.
The most common API pattern is to list possible filters as query params in swagger
then write custom model methods to retrieve the data.
For example, the shipment model calls `FetchShipmentsByTSP` to filter by a TSP.
Following this pattern would require writing individual methods for every column
(`Shipment` has around 50!).
Add on the requirement of combining filters
and the possibilities become exponential.

Along with being inflexible,
the current model fetching models are inconsistent across the codebase.
Filters, sorting, or joins are ad hoc based on handler or feature needs.
This causes issues with maintenance, performance consistency, and security scoping.

This ADR will mostly focus on filtering data,
because it is the primary required feature.
However, similar patterns can be applied to other querying features,
such as sorting and association embedding.

## Considered Alternatives

* **Standardize model querying and use code generation for query methods**

Let's continue with our shipment example.
Note that we'll use Pop here,
but we could also generate sqlx similarly.

Generated model methods may look like:

```go
// The main interface for defining a filter clause
type ShipmentQueryFilter interface {
    ApplyQuery(q pop.Query)  pop.Query // would apply a Where clause for the column
}

type CreatedAtFilter struct {
     createdAt time.Time
     comparator Comparator // string of possible comparators (=, >, <...)
}

func (f CreatedAtFilter) ApplyQuery(q pop.Query) pop.Query {
  // comparator is from a set of constants/enum,
  // not user defined
  // could also generate case statement here instead
  column := fmt.Sprintf("created_at %s", f.comparator)
  return query.Where(column, f.createAt)
}

type MoveFilter struct {
    moveID uuid.UUID
    comparator Comparator // string of possible comparators (=, >, <...)
}

func (f MoveFilter) ApplyQuery(q pop.Query) pop.Query {
  column := fmt.Sprintf("move_id %s", f.comparator)
  return query.Where(column, f.moveID)
}

// ... so on for every column

// FetchShipments is the only method exposed for fetching a list of shipments
func FetchShipments(db *pop.Connection, filters []ShipmentQueryFilter) []Shipments {
  var shipments Shipments
  query := db.Query()
  for _, filter := range filters {
    query := filter.ApplyQuery(query)
  }
  query.All(&shipments)
  return shipments
}
```

The handler or service (similar code in either)
would then parse all query params into filters.
This could also be generated:

```go
func (h ListShipmentsHandler) Handle(params shipmentop.ListShipmentParams) middleware.Responder {
  filters := make([]QueryFilter)
  // could also be a list
  if params.CreatedAtEq {
    filters = append(
      filters,
      CreatedAtFilter{
        time.Time(params.CreateAt).String(),
        Equals,
      })
  }
  if params.CreatedAtGreaterThan {
    filters = append(
      filters,
      CreatedAtFilter{
        time.Time(params.CreateAt).String(),
        GreaterThan,
      })
  }
  if params.MoveIDEqual {
    filters = append(
      filters,
      MoveIDFilter{
        uuid.FromString(params.MoveID.String()),
        Equals,
      })
  }
  // once again every possible column...

  shipment, err := FetchShipments(h.db, filters)
  // ... return a response
}
```

The swagger definition would then list all the possible column filters.
These would be arrays (ex. `move_ids_equal`),
rather than the singular parameters in our current API specs.
Like in the query builder below,
this filters could also use a more complex string,
such as `move_ids=eq,1` to avoid denoting every filter type.
However, we lose typing/go-swagger generation
because Open API does not support object nesting like this.

* **Use a third party query builder library**

Some Go sql libraries offer more flexibility than Pop or sqlx.
For example, [goqu](https://github.com/doug-martin/goqu)
or [squirrel](https://github.com/Masterminds/squirrel).

For example, goqu offers expression based querying like so:

```golang
sql, _, _ := db.From("items").Where(goqu.Ex{
  "col1": "a",
  "col2": 1,
  "col3": true,
  "col4": false,
  "col5": nil,
  "col6": []string{"a", "b", "c"},
}).ToSql()
```

* **Write a generic query interface that accepts any model**
Another option is to write a query builder as a dependency to handlers/services.
The proof of concept is as follows:

```go
type QueryFilter interface {
  Column() string
  Comparator() string
  Value() interface{}
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

func filteredQuery(query *pop.Query, filters []services.QueryFilter, t reflect.Type) (*pop.Query, error) {
  invalidFields := make([]string, 0)
  for _, f := range filters {
    column, ok := getDBColumn(t, f.Column())
    if !ok {
      invalidFields = append(
        invalidFields,
        fmt.Sprintf("%s %s", f.Column(), f.Comparator()),
      )
      continue
    }
    columnQuery := fmt.Sprintf("%s %s ?", column, comparator)
    query = query.Where(columnQuery, f.Value())
  }
  if len(invalidFields) != 0 {
    return query, fmt.Errorf("%v is not valid input", invalidFields)
  }
  return query, nil
}

// FetchMany fetches multiple model records using pop's All method
// Will return error if model is not pointer to slice of structs
func (p *Builder) FetchMany(model interface{}, filters []services.QueryFilter) error {
  t := reflect.TypeOf(model).Elem().Elem()
  query := p.db.Q()
  query, err := filteredQuery(query, filters, t)
  if err != nil {
    return err
  }
  return query.All(model)
}
```

And a service would utilize it as follows:

```go
type shipmentListQueryBuilder interface {
  FetchMany(model interface{}, filters []services.QueryFilter) error
}

// FetchShipmentList is uses the passed query builder to fetch a list of shipments
func (o *shipmentListFetcher) FetchShipmentList(filters []services.QueryFilter) (models.Shipments, error) {
  var shipments models.Shipments
  error := o.builder.FetchMany(&shipments, filters)
  return shipments, error
}
```

In this case, the API uses a generic filter parameter, such as:
`shipments?filter=created_at.gt:2019,moveID.eq:100`.
Note that we're going to run into object nesting issues discussed in this
[Open API Spec Issue](https://github.com/OAI/OpenAPI-Specification/issues/1706)

The handler (or middleware) marshals those filters into a `QueryFilter`
and passes them to the service

## Decision Outcome

* Chosen Alternative: **write a generic query builder**
* The main motivators of this decision are:
  * It accomplishes our desired feature set
  * It is relatively simple to implement
  * It makes our codebase more maintainable

## Pros and Cons of the Alternatives

### *Standardize model querying and use code generation for query methods*

* `+` Provides typed filters and code patterns are explicit
* `+` Generation would require in depth go-swagger and Pop/sqlx knowledge
* `-` Large API surface (individual params for each possible column)
* `-` It would be difficult to build
* `-` Due to development time, it could lock us into a specific db connector/ORM

### *Use a third party query builder library*

* `+` Wouldn't have to build/maintain it ourselves
* `-` Could lock us into a third party library
* `-` Most options don't provide much more than an ORM
* `-` Most options aren't well maintained or abandoned
* `-` Doesn't integrate into current toolset

### *Write a generic query interface that accepts any model*

* `+` Relatively lightweight/easy to build
* `+` Allows us to abstract the db connector/ORM through our own interface
* `+` Utilizes metadata already on models with little/no code generation
* `+` Reduces the API surface and boilerplate parameter definitions
* `-` Loses a lot of type safety with columns being strings.
  OpenAPI 2 does not support objects in query params.
  OpenAPI does, however it does not support them nested within an array.
* `-` Uses reflection (note this is already being used by Pop to create queries)
