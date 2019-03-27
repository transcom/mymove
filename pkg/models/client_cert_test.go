package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchClientCert() {
	digest := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	subject := "/C=US/ST=DC/L=Washington/O=Test/OU=Test Cert/CN=localhost"
	certNew := models.ClientCert{
		Sha256Digest:    digest,
		Subject:         subject,
		AllowDpsAuthAPI: true,
		AllowOrdersAPI:  true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	suite.MustSave(&certNew)

	cert, err := models.FetchClientCert(suite.DB(), digest)
	suite.Nil(err)
	suite.Equal(cert.Sha256Digest, digest)
	suite.Equal(cert.Subject, subject)
}

func (suite *ModelSuite) Test_FetchClientCertNotFound() {
	cert, err := models.FetchClientCert(suite.DB(), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	suite.Nil(cert)
	suite.Equal(models.ErrFetchNotFound, err)
}

func (suite *ModelSuite) Test_ClientCertValidations() {
	cert := &models.ClientCert{}

	var expErrors = map[string][]string{
		"sha256_digest": {"Sha256Digest can not be blank."},
		"subject":       {"Subject can not be blank."},
	}

	suite.verifyValidationErrors(cert, expErrors)
}

func (suite *ModelSuite) Test_ClientCertGetAllowedOrdersIssuersReadNone() {
	cert := models.ClientCert{}
	suite.Empty(cert.GetAllowedOrdersIssuersRead())
}

func (suite *ModelSuite) Test_ClientCertGetAllowedOrdersIssuersReadAll() {
	cert := models.ClientCert{
		AllowAirForceOrdersRead:    true,
		AllowArmyOrdersRead:        true,
		AllowCoastGuardOrdersRead:  true,
		AllowMarineCorpsOrdersRead: true,
		AllowNavyOrdersRead:        true,
	}
	suite.ElementsMatch(
		cert.GetAllowedOrdersIssuersRead(),
		[]string{
			string(models.IssuerAirForce),
			string(models.IssuerArmy),
			string(models.IssuerCoastGuard),
			string(models.IssuerMarineCorps),
			string(models.IssuerNavy),
		})
}
