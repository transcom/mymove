package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) loadServiceMap(dbTx *pop.Connection) error {
	var services models.ReServices
	err := dbTx.Select("id", "code").All(&services)
	if err != nil {
		return fmt.Errorf("could not read services: %w", err)
	}

	gre.serviceToIDMap = make(map[string]uuid.UUID)
	for _, service := range services {
		gre.serviceToIDMap[string(service.Code)] = service.ID
	}

	return nil
}
