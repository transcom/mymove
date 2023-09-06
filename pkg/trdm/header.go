package trdm

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"time"

	"github.com/opencontainers/go-digest"
)

type header struct {
	XMLName  xml.Name `xml:"soap:header"`
	Security security `xml:"wsse:Security"`
}

type security struct {
	Text                string              `xml:"chardata"`
	Wsse                string              `xml:"xmlns:wsse,attr"`
	Wsu                 string              `xml:"xmlns:wsu,attr"`
	BinarySecurityToken binarySecurityToken `xml:"BinarySecurityToken"`
	Signature           signature           `xml:"Signature"`
	Timestamp           timestamp           `xml:"TimeStamp"`
}
type signature struct {
	Text           string         `xml:",chardata"`
	ID             string         `xml:"Id,attr"`
	Ds             string         `xml:"ds,attr"`
	SignedInfo     signedInfo     `xml:"ds:SignedInfo"`
	KeyInfo        keyInfo        `xml:"ds:KeyInfo"`
	SignatureValue signatureValue `xml:"ds:SignatureValue"`
}
type keyInfo struct {
	ID                     string                 `xml:"Id"`
	SecurityTokenReference securityTokenReference `xml:"wsse:SecurityTokenReference"`
}
type securityTokenReference struct {
	STReference sTReference `xml:"wsse:Reference"`
}

type sTReference struct {
	URI       string `xml:"URI,attr"`
	ValueType string `xml:"ValueType,attr"`
}

type binarySecurityToken struct {
	Text         string `xml:",chardata"`
	EncodingType string `xml:"EncodingType,attr"`
	ValueType    string `xml:"ValueType,attr"`
	ID           string `xml:"Id,attr"`
}

type signedInfo struct {
	Text                   string                 `xml:",chardata"`
	CanonicalizationMethod canonicalizationMethod `xml:"ds:CanonicalizationMethod"`
	SignatureMethod        signatureMethod        `xml:"ds:SignatureMethod"`
	Reference              reference              `xml:"ds:Reference"`
}
type reference struct {
	URI          string       `xml:"URI,attr"`
	Transforms   transforms   `xml:"Transforms"`
	DigestMethod digestMethod `xml:"DigestMethod"`
	DigetValue   digestValue  `xml:"DigestValue"`
}
type digestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type digestValue struct {
	Text string `xml:",chardata"`
}
type transforms struct {
	Transform transform `xml:"Transform"`
}
type transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type canonicalizationMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}

type signatureMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}
type timestamp struct {
	ID      string `xml:"Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
}
type signatureValue struct {
	Text string `xml:",chardata"`
}

func GenerateSignedHeader(certificate *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	const certificateID = "X509-CertificateId"
	encodedDigest := digest.FromBytes([]byte(certificate.Raw)).Encoded()

	canonicalized := digest.Canonical.Encode([]byte(certificate.Raw))

	msgHash := sha512.New()
	_, err := msgHash.Write([]byte(canonicalized))
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	signedHash, err := privateKey.Sign(rand.Reader, msgHashSum, crypto.SHA512)
	if err != nil {
		return nil, err
	}
	// canonicalize & sign private key of x509 cert -> use this value for signaturevalue

	securityHeader := header{
		Security: security{
			Wsse: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
			Wsu:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			BinarySecurityToken: binarySecurityToken{
				EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
				ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
				Text:         base64.StdEncoding.EncodeToString(certificate.Raw),
				ID:           certificateID,
			},
			Signature: signature{
				Ds: "ttp://www.w3.org/2000/09/xmldsig#",
				SignedInfo: signedInfo{
					CanonicalizationMethod: canonicalizationMethod{
						Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
					},
					SignatureMethod: signatureMethod{
						Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512",
					},
					Reference: reference{
						URI: certificateID,
						Transforms: transforms{
							Transform: transform{
								Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
							},
						},
						DigestMethod: digestMethod{
							Algorithm: "http://www.w3.org/2001/04/xmlenc#sha512",
						},
						DigetValue: digestValue{
							Text: encodedDigest,
						},
					},
				},
				SignatureValue: signatureValue{
					Text: base64.StdEncoding.EncodeToString(signedHash),
				},
				KeyInfo: keyInfo{
					ID: "KI-KeyInfoIdentification",
					SecurityTokenReference: securityTokenReference{
						STReference: sTReference{
							URI:       certificateID,
							ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
						},
					},
				},
			},
			Timestamp: timestamp{
				Created: time.Now().UTC().Format(time.RFC3339),
				Expires: time.Now().Add(time.Millisecond * 5000).UTC().Format(time.RFC3339),
			},
		},
	}
	marshaledHeader, err := xml.Marshal(securityHeader)
	if err != nil {
		return nil, err
	}
	return marshaledHeader, nil
}
