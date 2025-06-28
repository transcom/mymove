package factory

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildClientCert creates a single ClientCert
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildClientCert(db *pop.Connection, customs []Customization, traits []Trait) models.ClientCert {
	customs = setupCustomizations(customs, traits)

	var cClientCert models.ClientCert
	// default to allowing prime
	defaultAllowPrime := true
	defaultAllowPPTAS := true
	var defaultPPTASAffiliation *models.ServiceMemberAffiliation
	if result := findValidCustomization(customs, ClientCert); result != nil {
		cClientCert = result.Model.(models.ClientCert)
		if result.LinkOnly {
			return cClientCert
		}
		// if customization is provided, explicitly override it to
		// allow false to override true
		defaultAllowPrime = cClientCert.AllowPrime
		defaultAllowPPTAS = cClientCert.AllowPPTAS
		defaultPPTASAffiliation = cClientCert.PPTASAffiliation
	}

	user := BuildUserAndUsersRoles(db, customs, traits)

	id := uuid.Must(uuid.NewV4())
	s := sha256.Sum256(id.Bytes())
	clientCert := models.ClientCert{
		ID:               id,
		Sha256Digest:     hex.EncodeToString(s[:]),
		Subject:          "/C=US/ST=DC/L=Washington/O=Truss/OU=AppClientTLS/CN=factory-" + id.String(),
		UserID:           user.ID,
		AllowPrime:       defaultAllowPrime,
		AllowPPTAS:       defaultAllowPPTAS,
		PPTASAffiliation: defaultPPTASAffiliation,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&clientCert, cClientCert)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &clientCert)
	}
	return clientCert
}

func BuildPrimeClientCert(db *pop.Connection) models.ClientCert {
	return BuildClientCert(db, nil, []Trait{GetTraitPrimeUser})
}

// Create dev client cert from 20191212230438_add_devlocal-mtls_client_cert.up.sql
const devlocalSha256Digest = "2c0c1fc67a294443292a9e71de0c71cc374fe310e8073f8cdc15510f6b0ef4db"
const devlocalSubject = "/C=US/ST=DC/L=Washington/O=Truss/OU=AppClientTLS/CN=devlocal"

var devlocalID = uuid.Must(uuid.FromString("190b1e07-eef8-445a-9696-5a2b49ee488d"))

// FetchOrBuildDevlocalClientCert tries fetching an existing clientCert, then falls back to creating one
func FetchOrBuildDevlocalClientCert(db *pop.Connection) models.ClientCert {
	traits := []Trait{GetTraitClientCertDevlocal, GetTraitPrimeUser}
	if db == nil {
		return BuildClientCert(db, nil, traits)
	}

	var clientCert models.ClientCert
	err := db.Q().Where(`sha256_digest=$1`, devlocalSha256Digest).First(&clientCert)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return clientCert
	}

	return BuildClientCert(db, nil, traits)
}

// Traits
// GetTraitOfficeUserActive sets the User and OfficeUser as Active
func GetTraitClientCertDevlocal() []Customization {
	return []Customization{
		{
			Model: models.ClientCert{
				ID:               devlocalID,
				Sha256Digest:     devlocalSha256Digest,
				Subject:          devlocalSubject,
				AllowPrime:       true,
				AllowPPTAS:       true,
				PPTASAffiliation: nil,
			},
		},
	}
}
