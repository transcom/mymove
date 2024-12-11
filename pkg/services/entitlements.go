package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// The weight allotment fetcher interface helps identify weight allotments (allowances) for a given grade
//
//go:generate mockery --name WeightAllotmentFetcher
type WeightAllotmentFetcher interface {
	GetWeightAllotment(appCtx appcontext.AppContext, grade string) (*models.HHGAllowance, error)
	GetAllWeightAllotments(appCtx appcontext.AppContext) (models.HHGAllowances, error)
}

// The weight restrictor interface helps apply weight restrictions to entitlements
//
//go:generate mockery --name WeightRestrictor
type WeightRestrictor interface {
	ApplyWeightRestrictionToEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, weightRestriction int, eTag string) (*models.Entitlement, error)
	RemoveWeightRestrictionFromEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, eTag string) (*models.Entitlement, error)
}
