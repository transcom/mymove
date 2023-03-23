package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

const defaultSHA256Digest string = "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"
const defaultSubject string = "CN=example-user,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US"
const customSubject string = "CN=custom-user,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US"

func (suite *FactorySuite) TestBuildClientCert() {
	suite.Run("Successful creation of default user", func() {
		// Under test:       BuildClientCert
		// Mocked:           None
		// Set up:           Create a ClientCert with no customizations or traits
		// Expected outcome: ClientCert should be created with default values

		certificate := BuildClientCert(suite.DB(), nil, nil)
		suite.Equal(defaultSHA256Digest, certificate.Sha256Digest)
		suite.Equal(defaultSubject, certificate.Subject)
		suite.NotNil(certificate.UserID)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:       BuildClientCert
		// Set up:           Create a ClientCert with a customized subject and no trait
		// Expected outcome: ClientCert should be created with default SHA256 digest and custom subject with no user association
		certificate := BuildClientCert(suite.DB(), []Customization{
			{
				Model: models.ClientCert{
					Subject: customSubject,
				},
			},
		}, nil)
		suite.Equal(defaultSHA256Digest, certificate.Sha256Digest)
		suite.Equal(customSubject, certificate.Subject)
		suite.NotNil(certificate.UserID)
	})

	suite.Run("Successful creation of user with trait", func() {
		// Under test:       BuildClientCert
		// Set up:           Create a Client Certificate with a trait
		// Expected outcome: User should be created with with default SHA256,
		// Subject, and an associated UserID

		user := BuildUser(suite.DB(), nil, nil)

		certificate := BuildClientCert(suite.DB(), []Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, nil)
		suite.Equal(defaultSHA256Digest, certificate.Sha256Digest)
		suite.Equal(defaultSubject, certificate.Subject)
		suite.Equal(user.ID, certificate.UserID)
	})

	suite.Run("Successful creation of user with both", func() {
		// Under test:       BuildClientCert
		// Set up:           Create a Client Certificate with a customized subject and active trait
		// Expected outcome: Cert should be created with default Sha256Digest, custom subject, and an associated UserID

		certificate := BuildClientCert(suite.DB(), []Customization{
			{
				Model: models.ClientCert{
					Subject: customSubject,
				},
			}}, []Trait{
			GetTraitActiveUser,
		})
		suite.Equal(defaultSHA256Digest, certificate.Sha256Digest)
		suite.Equal(customSubject, certificate.Subject)
		suite.NotNil(certificate.UserID)
	})

	suite.Run("Successful creation of stubbed user", func() {
		// Under test:       BuildClientCert
		// Set up:           Create a customized user, but don't pass in a db
		// Expected outcome: Client Certifiate should be created with email and active status
		//                   No user should be created in database
		precount, err := suite.DB().Count(&models.ClientCert{})
		suite.NoError(err)

		user := models.User{
			ID: uuid.Must(uuid.NewV4()),
		}
		certificate := BuildClientCert(nil, []Customization{
			{
				Model: models.ClientCert{
					Subject: customSubject,
				},
			},
			{
				Model: user,
			},
		}, nil)

		suite.Equal(customSubject, certificate.Subject)
		suite.Equal(user.ID, certificate.UserID)
		// Count how many certificates are in the DB, no certificates should have been created.
		count, err := suite.DB().Count(&models.ClientCert{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
}
