package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// The weight allotment fetcher interface helps identify weight allotments (allowances) for a given grade
//
//go:generate mockery --name WeightAllotmentFetcher
type WeightAllotmentFetcher interface {
	GetWeightAllotment(appCtx appcontext.AppContext, grade string, ordersType internalmessages.OrdersType) (models.WeightAllotment, error)
	GetAllWeightAllotments(appCtx appcontext.AppContext) (map[internalmessages.OrderPayGrade]models.WeightAllotment, error)
	GetWeightAllotmentByOrdersType(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType) (models.WeightAllotment, error)
	GetTotalWeightAllotment(appCtx appcontext.AppContext, order models.Order, entitlement models.Entitlement) (int, error)
}

// The weight restrictor interface helps apply weight restrictions to entitlements
//
//go:generate mockery --name WeightRestrictor
type WeightRestrictor interface {
	ApplyWeightRestrictionToEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, weightRestriction int, eTag string) (*models.Entitlement, error)
	RemoveWeightRestrictionFromEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, eTag string) (*models.Entitlement, error)
}
