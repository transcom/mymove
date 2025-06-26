package factory

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildClientCert() {
	suite.Run("Successful creation of default clientCert", func() {
		// Under test:      BuildClientCert
		// Mocked:          None
		// Set up:          Create a clientCert with no customizations or traits
		// Expected outcome:ClientCert should be created with default values

		// FUNCTION UNDER TEST
		clientCert := BuildClientCert(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotEmpty(clientCert.Sha256Digest)
		suite.Contains(clientCert.Subject, "OU=AppClientTLS/CN=factory-")
		suite.True(clientCert.AllowPrime)
		suite.True(clientCert.AllowPPTAS)
		suite.Nil(clientCert.PPTASAffiliation)
	})

	suite.Run("Successful creation of prime clientCert", func() {
		// Under test:      BuildPrimeClientCert
		// Mocked:          None
		// Set up:          Create a clientCert with no customizations or traits
		// Expected outcome:ClientCert should be created with default
		// values, associated with an active user with the prime role

		// FUNCTION UNDER TEST
		clientCert := BuildPrimeClientCert(suite.DB())

		// VALIDATE RESULTS
		suite.NotEmpty(clientCert.Sha256Digest)
		suite.Contains(clientCert.Subject, "OU=AppClientTLS/CN=factory-")
		suite.True(clientCert.AllowPrime)
		var user models.User
		suite.NoError(suite.DB().Eager("Roles").Find(&user, clientCert.UserID))
		suite.NotEmpty(user.Roles)
		suite.Equal(1, len(user.Roles), user.Roles)
		suite.Equal(roles.RoleTypePrime, user.Roles[0].RoleType)
	})

	suite.Run("Successful creation of customized clientCert", func() {
		// Under test:      BuildClientCert
		// Mocked:          None
		// Set up:          Create a clientCert with customization
		// Expected outcome:ClientCert should be created with customized values

		// SETUP
		// Create a custom clientCert to compare values
		s := sha256.Sum256([]byte("custom"))
		custClientCert := models.ClientCert{
			Sha256Digest: hex.EncodeToString(s[:]),
			Subject:      "CustomSubject",
			AllowPrime:   false,
		}

		customUser := BuildUser(suite.DB(), nil, nil)

		// FUNCTION UNDER TEST
		clientCert := BuildClientCert(suite.DB(), []Customization{
			{
				Model: custClientCert,
			},
			{
				Model:    customUser,
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(custClientCert.Sha256Digest, clientCert.Sha256Digest)
		suite.Equal(custClientCert.Subject, clientCert.Subject)
		suite.Equal(custClientCert.AllowPrime, clientCert.AllowPrime)
		suite.Equal(customUser.ID, clientCert.UserID)
	})

	suite.Run("Successful return of linkOnly clientCert", func() {
		// Under test:      BuildClientCert
		// Set up:          Create a clientCert and pass in a linkOnly clientCert
		// Expected outcome:No new clientCert should be created

		// Check num clientCerts
		precount, err := suite.DB().Count(&models.ClientCert{})
		suite.NoError(err)

		customCert := models.ClientCert{
			ID:           uuid.Must(uuid.NewV4()),
			Sha256Digest: hex.EncodeToString([]byte("LinkOnly")),
			Subject:      "LinkOnly",
		}

		clientCert := BuildClientCert(suite.DB(), []Customization{
			{
				Model:    customCert,
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ClientCert{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customCert.ID, clientCert.ID)
		suite.Equal(customCert.Sha256Digest, clientCert.Sha256Digest)
		suite.Equal(customCert.Subject, clientCert.Subject)
	})

	suite.Run("Two clientCerts should not be created", func() {
		// Under test:      FetchOrBuildDevlocalClientCert
		// Set up:          Create a clientCert with no customized state or traits
		// Expected outcome:Only 1 clientCert should be created
		count, potentialErr := suite.DB().Where(`sha256_digest=$1`, devlocalSha256Digest).Count(&models.ClientCert{})
		suite.NoError(potentialErr)
		suite.Zero(count)

		firstClientCert := FetchOrBuildDevlocalClientCert(suite.DB())

		secondClientCert := FetchOrBuildDevlocalClientCert(suite.DB())

		suite.Equal(firstClientCert.ID, secondClientCert.ID)

		existingClientCertCount, err := suite.DB().Where(`sha256_digest=$1`, devlocalSha256Digest).Count(&models.ClientCert{})
		suite.NoError(err)
		suite.Equal(1, existingClientCertCount)
	})

	suite.Run("Successful creation of pptas clientCert", func() {
		// Under test:      BuildClientCert
		// Mocked:          None
		// Set up:          Create a clientCert with pptas
		// Expected outcome:ClientCert should be created with pptas values

		// SETUP
		// Create a custom pptas clientCert to compare values
		s := sha256.Sum256([]byte("custom"))
		custClientCert := models.ClientCert{
			Sha256Digest:     hex.EncodeToString(s[:]),
			Subject:          "CustomSubject",
			AllowPrime:       false,
			AllowPPTAS:       true,
			PPTASAffiliation: (*models.ServiceMemberAffiliation)(models.StringPointer("MARINES")),
		}

		customUser := BuildUser(suite.DB(), nil, nil)

		// FUNCTION UNDER TEST
		clientCert := BuildClientCert(suite.DB(), []Customization{
			{
				Model: custClientCert,
			},
			{
				Model:    customUser,
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(custClientCert.PPTASAffiliation, clientCert.PPTASAffiliation)
		suite.Equal(custClientCert.Sha256Digest, clientCert.Sha256Digest)
		suite.Equal(custClientCert.Subject, clientCert.Subject)
		suite.Equal(custClientCert.AllowPrime, clientCert.AllowPrime)
		suite.Equal(customUser.ID, clientCert.UserID)
	})
}
