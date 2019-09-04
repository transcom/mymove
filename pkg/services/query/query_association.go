package query

import (
	"github.com/transcom/mymove/pkg/services"
)

type queryAssociation struct {
	field string
}

func (q queryAssociation) Field() string {
	return q.field
}

func NewQueryAssociation(field string) services.QueryAssociation {
	return queryAssociation{
		field,
	}
}

type queryAssociations struct {
	associations []services.QueryAssociation
}

func (qa queryAssociations) StringGetAssociations() []string {
	associations := make([]string, 0, len(qa.associations))

	for _, a := range qa.associations {
		associations = append(associations, a.Field())
	}

	return associations
}

func NewQueryAssociations(associations []services.QueryAssociation) services.QueryAssociations {
	return queryAssociations{
		associations,
	}
}
