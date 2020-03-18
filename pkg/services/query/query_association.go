package query

import (
	"github.com/transcom/mymove/pkg/services"
)

type queryAssociation struct {
	field string
}

// Field returns a field on a query association
func (q queryAssociation) Field() string {
	return q.field
}

// NewQueryAssociation creates a new query association
func NewQueryAssociation(field string) services.QueryAssociation {
	return queryAssociation{
		field,
	}
}

type queryAssociations struct {
	associations []services.QueryAssociation
}

// StringGetAssociations returns a slice of string associations
func (qa queryAssociations) StringGetAssociations() []string {
	associations := make([]string, 0, len(qa.associations))

	for _, a := range qa.associations {
		associations = append(associations, a.Field())
	}

	return associations
}

// NewQueryAssociations returns new query associations
func NewQueryAssociations(associations []services.QueryAssociation) services.QueryAssociations {
	return queryAssociations{
		associations,
	}
}
