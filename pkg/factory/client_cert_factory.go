package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildClientCert creates a Client Certificate
func BuildClientCert(db *pop.Connection, customs []Customization, traits []Trait) models.ClientCert {
	customs = setupCustomizations(customs, traits)

	// Find ClientCert assertion and covert models to ClientCert
	var cClientCert models.ClientCert
	if result := findValidCustomization(customs, ClientCert); result != nil {
		cClientCert = result.Model.(models.ClientCert)
		if result.LinkOnly {
			return cClientCert
		}
	}

	// find/create user and roles
	user := BuildUserAndUsersRoles(db, customs, traits)

	// create the client certificate
	certificate := models.ClientCert{
		Subject:      "CN=example-user,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US",
		Sha256Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
		UserID:       user.ID,
	}

	// Overwrite the values wtih those from assetions
	testdatagen.MergeModels(&certificate, cClientCert)

	if db != nil {
		mustCreate(db, &certificate)
	}

	return certificate
}
