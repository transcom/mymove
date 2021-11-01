package clientcert

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type clientCertUpdater struct {
	builder clientCertQueryBuilder
}

func (o *clientCertUpdater) UpdateClientCert(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.ClientCertUpdatePayload) (*models.ClientCert, *validate.Errors, error) {
	var foundClientCert models.ClientCert
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundClientCert, filters)

	if err != nil {
		return nil, nil, err
	}

	if payload.AllowOrdersAPI != nil {
		foundClientCert.AllowOrdersAPI = *payload.AllowOrdersAPI
	}
	if payload.AllowAirForceOrdersRead != nil {
		foundClientCert.AllowAirForceOrdersRead = *payload.AllowAirForceOrdersRead
	}
	if payload.AllowAirForceOrdersWrite != nil {
		foundClientCert.AllowAirForceOrdersWrite = *payload.AllowAirForceOrdersWrite
	}
	if payload.AllowArmyOrdersRead != nil {
		foundClientCert.AllowArmyOrdersRead = *payload.AllowArmyOrdersRead
	}
	if payload.AllowArmyOrdersWrite != nil {
		foundClientCert.AllowArmyOrdersWrite = *payload.AllowArmyOrdersWrite
	}
	if payload.AllowCoastGuardOrdersRead != nil {
		foundClientCert.AllowCoastGuardOrdersRead = *payload.AllowCoastGuardOrdersRead
	}
	if payload.AllowCoastGuardOrdersWrite != nil {
		foundClientCert.AllowCoastGuardOrdersWrite = *payload.AllowCoastGuardOrdersWrite
	}
	if payload.AllowMarineCorpsOrdersRead != nil {
		foundClientCert.AllowMarineCorpsOrdersRead = *payload.AllowMarineCorpsOrdersRead
	}
	if payload.AllowMarineCorpsOrdersWrite != nil {
		foundClientCert.AllowMarineCorpsOrdersWrite = *payload.AllowMarineCorpsOrdersWrite
	}
	if payload.AllowNavyOrdersRead != nil {
		foundClientCert.AllowNavyOrdersRead = *payload.AllowNavyOrdersRead
	}
	if payload.AllowNavyOrdersWrite != nil {
		foundClientCert.AllowNavyOrdersWrite = *payload.AllowNavyOrdersWrite
	}
	if payload.AllowPrime != nil {
		foundClientCert.AllowPrime = *payload.AllowPrime
	}

	verrs, err := o.builder.UpdateOne(appCtx, &foundClientCert, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundClientCert, nil, nil
}

// NewClientCertUpdater returns a new admin user updater builder
func NewClientCertUpdater(builder clientCertQueryBuilder) services.ClientCertUpdater {
	return &clientCertUpdater{builder}
}
