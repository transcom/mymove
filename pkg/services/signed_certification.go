package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// SignedCertificationCreator creates a signed certification
//
//go:generate mockery --name SignedCertificationCreator --disable-version-string
type SignedCertificationCreator interface {
	CreateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification) (*models.SignedCertification, error)
}

// SignedCertificationUpdater updates a signed certification
//
//go:generate mockery --name SignedCertificationUpdater --disable-version-string
type SignedCertificationUpdater interface {
	UpdateSignedCertification(appCtx appcontext.AppContext, signedCertification models.SignedCertification, eTag string) (*models.SignedCertification, error)
}
