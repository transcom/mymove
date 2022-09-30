package services

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ServiceMemberAssociator is the service object interface for CreateServiceMember
//
//go:generate mockery --name ServiceMemberAssociator --disable-version-string
type ServiceMemberAssociator interface {
	CreateServiceMember(appCtx appcontext.AppContext, newServiceMember models.ServiceMember) (*validate.Errors, error)
}
