package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"
)

//go:generate mockery --name PrivilegeFetcher
type PrivilegeFetcher interface {
	FetchPrivilegeTypes(appCtx appcontext.AppContext) ([]roles.PrivilegeType, error)
}
