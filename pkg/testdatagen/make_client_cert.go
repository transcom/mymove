package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func MakeDevClientCert(db *pop.Connection, assertions Assertions) {
	clientCert := models.ClientCert{
		ID:                          uuid.Must(uuid.FromString("190b1e07-eef8-445a-9696-5a2b49ee488d")),
		Sha256Digest:                "2c0c1fc67a294443292a9e71de0c71cc374fe310e8073f8cdc15510f6b0ef4db",
		Subject:                     "/C=US/ST=DC/L=Washington/O=Truss/OU=AppClientTLS/CN=devlocal",
		AllowDpsAuthAPI:             false,
		AllowOrdersAPI:              true,
		AllowAirForceOrdersRead:     true,
		AllowAirForceOrdersWrite:    true,
		AllowArmyOrdersRead:         true,
		AllowArmyOrdersWrite:        true,
		AllowCoastGuardOrdersRead:   true,
		AllowCoastGuardOrdersWrite:  true,
		AllowMarineCorpsOrdersRead:  true,
		AllowMarineCorpsOrdersWrite: true,
		AllowNavyOrdersRead:         true,
		AllowNavyOrdersWrite:        true,
		AllowPrime:                  true,
		UserID:                      uuid.UUID{},
	}

	mergeModels(&clientCert, assertions.ClientCert)

	existingCert, err := models.FetchClientCert(db, clientCert.Sha256Digest)
	if err != nil && errors.Cause(err).Error() == models.ErrFetchNotFound.Error() {
		// dev client cert not found; create new client cert
		mustCreate(db, &clientCert, false)
	} else if err == nil && existingCert != nil {
		// client cert already exists; update existing client cert
		mergeModels(&existingCert, clientCert)
		MustSave(db, &existingCert)
	}
}
