package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchClientCert() {
	digest := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	subject := "C=US, ST=DC, L=Washington, O=Test, OU=Test Cert, CN=localhost"
	certNew := models.ClientCert{
		Sha256Digest:    digest,
		Subject:         subject,
		AllowDpsAuthAPI: true,
		AllowOrdersAPI:  true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	suite.mustSave(&certNew)

	cert, err := models.FetchClientCert(suite.db, digest)
	suite.Nil(err)
	suite.Equal(cert.Sha256Digest, digest)
	suite.Equal(cert.Subject, subject)
}

func (suite *ModelSuite) Test_FetchClientCertNotFound() {
	cert, err := models.FetchClientCert(suite.db, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	suite.Nil(cert)
	suite.Equal(models.ErrFetchNotFound, err)
}
