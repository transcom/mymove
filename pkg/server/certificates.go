package server

import (
	"crypto/x509"

	"go.mozilla.org/pkcs7"
)

// AddToCertPoolFromPkcs7Package reads the certificates in a DER-encoded PKCS7
// package and adds those Certificates to the x509.CertPool
func AddToCertPoolFromPkcs7Package(certPool *x509.CertPool, pkcs7Package []byte) error {
	p7, err := pkcs7.Parse(pkcs7Package)
	if err != nil {
		return err
	}
	for _, cert := range p7.Certificates {
		certPool.AddCert(cert)
	}
	return nil
}
