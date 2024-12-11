package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// The weight allotment fetcher interface helps identify weight allotments (allowances) for a given grade
//
//go:generate mockery --name WeightAllotmentFetcher
type WeightAllotmentFetcher interface {
	GetWeightAllotment(appCtx appcontext.AppContext, grade string) (*models.HHGAllowance, error)
}

// The weight restrictor interface helps apply weight restrictions to entitlements
//
//go:generate mockery --name WeightRestrictor
type WeightRestrictor interface {
	ApplyWeightRestrictionToEntitlement(appCtx appcontext.AppContext, entitlementID uuid.UUID, weightRestriction int) error
	RemoveWeightRestrictionFromEntitlement(appCtx appcontext.AppContext, entitlementID uuid.UUID) error
}
