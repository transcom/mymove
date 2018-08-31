package server

import (
	tls "crypto/tls"
)

// TLSCert is
type TLSCert struct {
	CertPEMBlock []byte
	KeyPEMBlock  []byte
}

// ParseTLSCert is
func ParseTLSCert(certs []TLSCert) ([]tls.Certificate, error) {
	var parsedCerts []tls.Certificate

	for _, cert := range certs {
		parsedCert, err := tls.X509KeyPair(cert.CertPEMBlock, cert.KeyPEMBlock)
		if err != nil {
			//TODO Add error message
			return nil, err
		}
		parsedCerts = append(parsedCerts, parsedCert)

	}
	return parsedCerts, nil
}
