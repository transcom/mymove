package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchClientCert() {
	oktaID := uuid.Must(uuid.NewV4()).String()
	userForClientCert := models.User{
		ID:        uuid.Must(uuid.NewV4()),
		OktaID:    oktaID,
		OktaEmail: "prime_user_with_client_cert@okta.mil",
		Active:    true,
	}
	suite.MustCreate(&userForClientCert)

	digest := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	subject := "/C=US/ST=DC/L=Washington/O=Test/OU=Test Cert/CN=localhost"
	allowpptas := true
	pptasaffiliation := (*models.ServiceMemberAffiliation)(models.StringPointer("MARINES"))
	certNew := models.ClientCert{
		Sha256Digest:     digest,
		Subject:          subject,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		UserID:           userForClientCert.ID,
		AllowPPTAS:       allowpptas,
		PPTASAffiliation: pptasaffiliation,
	}
	suite.MustSave(&certNew)

	cert, err := models.FetchClientCert(suite.DB(), digest)
	suite.NoError(err)
	suite.Equal(cert.Sha256Digest, digest)
	suite.Equal(cert.Subject, subject)
	suite.Equal(cert.AllowPPTAS, allowpptas)
	suite.Equal(cert.PPTASAffiliation, pptasaffiliation)
}

func (suite *ModelSuite) Test_FetchClientCertNotFound() {
	cert, err := models.FetchClientCert(suite.DB(), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	suite.Nil(cert)
	suite.Equal(models.ErrFetchNotFound, err)
}

func (suite *ModelSuite) Test_ClientCertValidations() {
	cert := &models.ClientCert{}

	var expErrors = map[string][]string{
		"sha256_digest": {"Sha256Digest not in range(64, 64)"},
		"subject":       {"Subject can not be blank."},
		"user_id":       {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(cert, expErrors, nil)
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
