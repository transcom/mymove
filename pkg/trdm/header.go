package trdm

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"time"

	"github.com/opencontainers/go-digest"
)

type header struct {
	XMLName  xml.Name `xml:"soap:header"`
	Security security `xml:"wsse:Security"`
}

type security struct {
	// ! Text may not be needed here
	Text                string              `xml:"chardata"`
	Wsse                string              `xml:"xmlns:wsse,attr"`
	Wsu                 string              `xml:"xmlns:wsu,attr"`
	BinarySecurityToken binarySecurityToken `xml:"wsse:BinarySecurityToken"`
	Signature           signature           `xml:"ds:Signature"`
	Timestamp           timestamp           `xml:"wsu:TimeStamp"`
}
type signature struct {
	Text           string         `xml:",chardata"`
	ID             string         `xml:"Id,attr"`
	Ds             string         `xml:"xmlns:ds,attr"`
	SignedInfo     signedInfo     `xml:"ds:SignedInfo"`
	SignatureValue signatureValue `xml:"ds:SignatureValue"`
	KeyInfo        keyInfo        `xml:"ds:KeyInfo"`
}
type keyInfo struct {
	ID                     string                 `xml:"Id,attr"`
	SecurityTokenReference securityTokenReference `xml:"wsse:SecurityTokenReference"`
}
type securityTokenReference struct {
	ID          string      `xml:"wsu:Id,attr"`
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
	ID           string `xml:"wsu:Id,attr"`
}

type signedInfo struct {
	CanonicalizationMethod canonicalizationMethod `xml:"ds:CanonicalizationMethod"`
	SignatureMethod        signatureMethod        `xml:"ds:SignatureMethod"`
	// For hitting the TRDM V7 endpoints there are typically three references.
	//    1. The timestamp id at the bottom of the envelope. wsu:Timestamp
	//    2. Inside of KeyInfo -> SecurityTokenReference -> wsse:Reference
	//    3. soap:Body
	Reference []reference `xml:"ds:Reference"`
}
type reference struct {
	URI          string       `xml:"URI,attr"`
	Transforms   transforms   `xml:"ds:Transforms"`
	DigestMethod digestMethod `xml:"ds:DigestMethod"`
	DigestValue  digestValue  `xml:"ds:DigestValue"`
}
type digestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type digestValue struct {
	Text string `xml:",chardata"`
}
type transforms struct {
	Transform transform `xml:"ds:Transform"`
}
type transform struct {
	Algorithm           string              `xml:"Algorithm,attr"`
	InclusiveNamespaces inclusiveNameSpaces `xml:"ec:InclusiveNamespaces"`
}
type inclusiveNameSpaces struct {
	PrefixList string `xml:"PrefixList,attr"`
	Ec         string `xml:"xmlns:ec,attr"`
}
type canonicalizationMethod struct {
	Text                string              `xml:",chardata"`
	Algorithm           string              `xml:"Algorithm,attr"`
	InclusiveNamespaces inclusiveNameSpaces `xml:"ec:InclusiveNamespaces"`
}

type signatureMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}
type timestamp struct {
	ID      string `xml:"wsu:Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
}
type signatureValue struct {
	Text string `xml:",chardata"`
}

func GenerateSOAPURIWithPrefix(prefix string) (string, error) {
	randBytes := make([]byte, 8)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return prefix + "-" + hex.EncodeToString(randBytes), nil
}

// ! TODO: Ensure the bodyReferenceURI correctly references the body in the envelope. Includes refactoring envelope.
func GenerateSignedHeader(certificate *x509.Certificate, privateKey *rsa.PrivateKey, bodyReferenceURI string) ([]byte, error) {
	const certificateID = "X509-CertificateId"
	const signatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"
	const digestAlgorithm = "http://www.w3.org/2001/04/xmlenc#sha512"
	const signatureCanonicalizationMethod = "http://www.w3.org/2001/10/xml-exc-c14n#"
	wsseReferenceURI, err := GenerateSOAPURIWithPrefix("#X509")
	if err != nil {
		return nil, err
	}
	timestampReferenceID, err := GenerateSOAPURIWithPrefix("#TS")
	if err != nil {
		return nil, err
	}
	securityTokenReferenceID, err := GenerateSOAPURIWithPrefix("STR")
	if err != nil {
		return nil, err
	}
	keyInfoReferenceID, err := GenerateSOAPURIWithPrefix("KI")
	if err != nil {
		return nil, err
	}
	encodedDigest := digest.FromBytes([]byte(certificate.Raw)).Encoded()

	canonicalized := digest.Canonical.Encode([]byte(certificate.Raw))

	msgHash := sha512.New()
	_, err = msgHash.Write([]byte(canonicalized))
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
				ID:           wsseReferenceURI,
			},
			Signature: signature{
				// ! Update ID value
				ID: "SignatureID",
				Ds: "http://www.w3.org/2000/09/xmldsig#",
				SignedInfo: signedInfo{
					CanonicalizationMethod: canonicalizationMethod{
						Algorithm: signatureCanonicalizationMethod,
						InclusiveNamespaces: inclusiveNameSpaces{
							PrefixList: "ret soap",
							Ec:         signatureCanonicalizationMethod,
						},
					},
					SignatureMethod: signatureMethod{
						Algorithm: signatureAlgorithm,
					},
					Reference: []reference{
						{
							// References the Timestamp's wsu:Id, the one at the bottom of the envelope
							URI: timestampReferenceID,
							Transforms: transforms{
								Transform: transform{
									Algorithm: signatureCanonicalizationMethod,
									InclusiveNamespaces: inclusiveNameSpaces{
										PrefixList: "wsse ret soap",
										Ec:         signatureCanonicalizationMethod,
									},
								},
							},
							DigestMethod: digestMethod{
								Algorithm: digestAlgorithm,
							},
							DigestValue: digestValue{
								Text: encodedDigest,
							},
						},
						{
							// References KeyInfo's wsse:Reference URI
							URI: wsseReferenceURI,
							Transforms: transforms{
								Transform: transform{
									Algorithm: signatureCanonicalizationMethod,
									InclusiveNamespaces: inclusiveNameSpaces{
										// Prefix intentionaly left blank
										PrefixList: "",
										Ec:         signatureCanonicalizationMethod,
									},
								},
							},
							DigestMethod: digestMethod{
								Algorithm: digestAlgorithm,
							},
							DigestValue: digestValue{
								Text: encodedDigest,
							},
						},
						{
							// References the body
							URI: bodyReferenceURI,
							Transforms: transforms{
								Transform: transform{
									Algorithm: signatureCanonicalizationMethod,
									InclusiveNamespaces: inclusiveNameSpaces{
										PrefixList: "ret",
										Ec:         signatureCanonicalizationMethod,
									},
								},
							},
							DigestMethod: digestMethod{
								Algorithm: digestAlgorithm,
							},
							DigestValue: digestValue{
								Text: encodedDigest,
							},
						},
					},
				},
				SignatureValue: signatureValue{
					Text: base64.StdEncoding.EncodeToString(signedHash),
				},
				KeyInfo: keyInfo{
					ID: keyInfoReferenceID,
					SecurityTokenReference: securityTokenReference{
						ID: securityTokenReferenceID,
						STReference: sTReference{
							URI:       certificateID,
							ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
						},
					},
				},
			},
			Timestamp: timestamp{
				ID:      timestampReferenceID,
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
