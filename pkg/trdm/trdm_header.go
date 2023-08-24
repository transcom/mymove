package trdm

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/xml"
	"time"

	"github.com/opencontainers/go-digest"
)

type Header struct {
	XMLName  xml.Name `xml:"soap:header"`
	Security Security `xml:"wsse:Security"`
}

type Security struct {
	Text                string              `xml:"chardata"`
	Wsse                string              `xml:"xmlns:wsse,attr"`
	Wsu                 string              `xml:"xmlns:wsu,attr"`
	BinarySecurityToken BinarySecurityToken `xml:"BinarySecurityToken"`
	Signature           Signature           `xml:"Signature"`
	Timestamp           Timestamp           `xml:"TimeStamp"`
}
type Signature struct {
	Text           string         `xml:",chardata"`
	ID             string         `xml:"Id,attr"`
	Ds             string         `xml:"ds,attr"`
	SignedInfo     SignedInfo     `xml:"ds:SignedInfo"`
	KeyInfo        KeyInfo        `xml:"ds:KeyInfo"`
	SignatureValue SignatureValue `xml:"ds:SignatureValue"`
}
type KeyInfo struct {
	ID                     string                 `xml:"Id"`
	SecurityTokenReference SecurityTokenReference `xml:"wsse:SecurityTokenReference"`
}
type SecurityTokenReference struct {
	STReference STReference `xml:"wsse:Reference"`
}

type STReference struct {
	URI       string `xml:"URI,attr"`
	ValueType string `xml:"ValueType,attr"`
}

type BinarySecurityToken struct {
	Text         string `xml:",chardata"`
	EncodingType string `xml:"EncodingType,attr"`
	ValueType    string `xml:"ValueType,attr"`
	ID           string `xml:"Id,attr"`
}

type SignedInfo struct {
	Text                   string                 `xml:",chardata"`
	CanonicalizationMethod CanonicalizationMethod `xml:"ds:CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"ds:SignatureMethod"`
	Reference              Reference              `xml:"ds:Reference"`
}
type Reference struct {
	URI          string       `xml:"URI,attr"`
	Transforms   Transforms   `xml:"Transforms"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigetValue   DigestValue  `xml:"DigestValue"`
}
type DigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type DigestValue struct {
	Text string `xml:",chardata"`
}
type Transforms struct {
	Transform Transform `xml:"Transform"`
}
type Transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type CanonicalizationMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}
type Timestamp struct {
	ID      string `xml:"Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
}
type SignatureValue struct {
	Text []byte `xml:",chardata"`
}

func GenerateSignedHeader(certificate string, body []byte, privateKey *rsa.PrivateKey) ([]byte, error) {

	print(privateKey)
	const certificateID = "X509-CertificateId"
	encodedDigest := digest.FromBytes(body).Encoded()

	msgHash := sha256.New()
	_, err := msgHash.Write(body)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	signedHash, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, msgHashSum)
	if err != nil {
		return nil, err
	}

	securityHeader := Header{
		Security: Security{
			Wsse: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
			Wsu:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			BinarySecurityToken: BinarySecurityToken{
				EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
				ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
				Text:         certificate,
				ID:           certificateID,
			},
			Signature: Signature{
				Ds: "ttp://www.w3.org/2000/09/xmldsig#",
				SignedInfo: SignedInfo{
					CanonicalizationMethod: CanonicalizationMethod{
						Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
					},
					SignatureMethod: SignatureMethod{
						Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
					},
					Reference: Reference{
						URI: certificateID,
						Transforms: Transforms{
							Transform: Transform{
								Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
							},
						},
						DigestMethod: DigestMethod{
							Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
						},
						DigetValue: DigestValue{
							Text: encodedDigest,
						},
					},
				},
				SignatureValue: SignatureValue{
					Text: []byte(signedHash),
				},
				KeyInfo: KeyInfo{
					ID: "KI-KeyInfoIdentification",
					SecurityTokenReference: SecurityTokenReference{
						STReference: STReference{
							URI:       certificateID,
							ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
						},
					},
				},
			},
			Timestamp: Timestamp{
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
