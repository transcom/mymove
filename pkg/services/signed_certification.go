package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// SignedCertificationCreator creates a signed certification
//
//go:generate mockery --name SignedCertificationCreator
type SignedCertificationCreator interface {
	CreateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification) (*models.SignedCertification, error)
}

// SignedCertificationUpdater updates a signed certification
//
//go:generate mockery --name SignedCertificationUpdater
type SignedCertificationUpdater interface {
	UpdateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification, eTag string) (*models.SignedCertification, error)
}
