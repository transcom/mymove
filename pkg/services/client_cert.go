package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// ClientCertListFetcher is the exported interface for fetching multiple client certs
//go:generate mockery --name ClientCertListFetcher --disable-version-string
type ClientCertListFetcher interface {
	FetchClientCertList(appCtx appcontext.AppContext, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) (models.ClientCerts, error)
	FetchClientCertCount(appCtx appcontext.AppContext, filters []QueryFilter) (int, error)
}

// ClientCertFetcher is the exported interface for fetching a single client cert
//go:generate mockery --name ClientCertFetcher --disable-version-string
type ClientCertFetcher interface {
	FetchClientCert(appCtx appcontext.AppContext, filters []QueryFilter) (models.ClientCert, error)
}

// ClientCertCreator is the exported interface for creating an client cert
//go:generate mockery --name ClientCertCreator --disable-version-string
type ClientCertCreator interface {
	CreateClientCert(appCtx appcontext.AppContext, user *models.ClientCert) (*models.ClientCert, *validate.Errors, error)
}

// ClientCertUpdater is the exported interface for updating an client cert
//go:generate mockery --name ClientCertUpdater --disable-version-string
type ClientCertUpdater interface {
	UpdateClientCert(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.ClientCertUpdatePayload) (*models.ClientCert, *validate.Errors, error)
}
