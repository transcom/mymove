package query

import (
	"fmt"
	"github.com/transcom/mymove/pkg/services"
)

type queryAssociation struct {
	model string
	column string
}

func(a queryAssociation) Association() string {
	return fmt.Sprintf("%s.%s", a.model, a.column)
}

func NewQueryAssociation(model string, column string) services.QueryAssociation {
	return queryAssociation{
		model,
		column,
	}
}

type queryAssociations struct {
	associations []queryAssociation
}

func(as queryAssociations) StringGetAssociations () []string {
	associations := make([]string, 0, len(as.associations))

	for _, a := range as.associations {
		associations = append(associations, a.Association())
	}

	return associations
}

func NewQueryAssociations(associations []queryAssociation) services.QueryAssociations {
	return queryAssociations{
		associations,
	}
}

