package server

import (
	"crypto/x509"

	"go.mozilla.org/pkcs7"
)

// LoadCertPoolFromPkcs7Package reads the certificates in a DER-encoded PKCS7
// package and returns a newly-created x509.CertPool with those Certificates
// added.
func LoadCertPoolFromPkcs7Package(pkcs7Package []byte) (*x509.CertPool, error) {
	p7, err := pkcs7.Parse(pkcs7Package)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	for _, cert := range p7.Certificates {
		certPool.AddCert(cert)
	}
	return certPool, nil
}
