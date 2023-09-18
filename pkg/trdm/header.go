package trdm

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"time"

	"github.com/ucarion/c14n"
)

const SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"
const DigestAlgorithm = "http://www.w3.org/2001/04/xmlenc#sha512"
const SignatureCanonicalizationMethod = "http://www.w3.org/2001/10/xml-exc-c14n#"

type header struct {
	XMLName  xml.Name `xml:"soap:Header"`
	Security security `xml:"wsse:Security"`
}

type security struct {
	Wsse                string              `xml:"xmlns:wsse,attr"`
	Wsu                 string              `xml:"xmlns:wsu,attr"`
	BinarySecurityToken binarySecurityToken `xml:"wsse:BinarySecurityToken"`
	Signature           signature           `xml:"ds:Signature"`
	Timestamp           timestamp           `xml:"wsu:Timestamp"`
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
	Created string `xml:"wsu:Created"`
	Expires string `xml:"wsu:Expires"`
}
type signatureValue struct {
	Text string `xml:",chardata"`
}

// Generate SHA-512 digest
// Returns
// - hashValue
// - encodedHash
// - Error
func GenerateDigest(data []byte) ([]byte, string, error) {
	hasher := sha512.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, "", err
	}
	hashValue := hasher.Sum(nil)
	encodedHash := base64.StdEncoding.EncodeToString(hashValue)
	return hashValue, encodedHash, nil
}

func CanonicalizeXML(xmlByte []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlByte))
	out, err := c14n.Canonicalize(decoder)
	if err != nil {
		return nil, nil
	}
	return out, nil
}

// Generate timestamp XML and digest
// Returns
// - Struct
// - Canon XML Byte
// - Digest Value
// - Error
func GenerateTimestampAndDigest() (timestamp, []byte, string, error) {
	tsID, err := GenerateSOAPURIWithPrefix("TS")
	if err != nil {
		return timestamp{}, nil, "", err
	}
	ts := timestamp{
		ID:      tsID,
		Created: time.Now().UTC().Format(time.RFC3339),
		// Currently 10 minutes for testing
		Expires: time.Now().Add(time.Minute * 10).UTC().Format(time.RFC3339),
	}

	xmlByte, digest, err := GenerateSecurityElement(ts)
	if err != nil {
		return timestamp{}, nil, "", err
	}

	return ts, xmlByte, digest, nil
}
func GenerateSignedInfoAndDigest(timestampID string, timestampDigest string, bodyID string, bodyDigest string, x509ID string, x509Digest string) (signedInfo, []byte, string, error) {

	signedInfoStruct := signedInfo{
		CanonicalizationMethod: canonicalizationMethod{
			Algorithm: SignatureCanonicalizationMethod,
			InclusiveNamespaces: inclusiveNameSpaces{
				PrefixList: "ret soap",
				Ec:         SignatureCanonicalizationMethod,
			},
		},
		SignatureMethod: signatureMethod{
			Algorithm: SignatureAlgorithm,
		},
		Reference: []reference{
			{
				// References the Timestamp's wsu:Id, the one at the bottom of the envelope
				// Prepend with # to reference
				URI: "#" + timestampID,
				Transforms: transforms{
					Transform: transform{
						Algorithm: SignatureCanonicalizationMethod,
						InclusiveNamespaces: inclusiveNameSpaces{
							PrefixList: "wsse ret soap",
							Ec:         SignatureCanonicalizationMethod,
						},
					},
				},
				DigestMethod: digestMethod{
					Algorithm: DigestAlgorithm,
				},
				DigestValue: digestValue{
					Text: timestampDigest,
				},
			},
			{
				// References the bodyID with #
				URI: "#" + bodyID,
				Transforms: transforms{
					Transform: transform{
						Algorithm: SignatureCanonicalizationMethod,
						InclusiveNamespaces: inclusiveNameSpaces{
							PrefixList: "ret",
							Ec:         SignatureCanonicalizationMethod,
						},
					},
				},
				DigestMethod: digestMethod{
					Algorithm: DigestAlgorithm,
				},
				DigestValue: digestValue{
					Text: bodyDigest,
				},
			},
			{
				// References KeyInfo's wsse:Reference URI
				// Prepend with # to reference
				URI: "#" + x509ID,
				Transforms: transforms{
					Transform: transform{
						Algorithm: SignatureCanonicalizationMethod,
						InclusiveNamespaces: inclusiveNameSpaces{
							// Prefix intentionaly left blank
							PrefixList: "",
							Ec:         SignatureCanonicalizationMethod,
						},
					},
				},
				DigestMethod: digestMethod{
					Algorithm: DigestAlgorithm,
				},
				DigestValue: digestValue{
					Text: x509Digest,
				},
			},
		},
	}

	xmlByte, digest, err := GenerateSecurityElement(signedInfoStruct)
	if err != nil {
		return signedInfo{}, nil, "", err
	}

	return signedInfoStruct, xmlByte, digest, nil
}
func GenerateSOAPURIWithPrefix(prefix string) (string, error) {
	randBytes := make([]byte, 8)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return prefix + "-" + hex.EncodeToString(randBytes), nil
}

// Returns
// - XML Byte
// - XML Digest
// - Error
func GenerateSecurityElement(t interface{}) ([]byte, string, error) {
	// ID should already be located inside of the interface
	noncanonXML, err := xml.Marshal(t)
	if err != nil {
		return nil, "", err
	}
	canonXML, err := CanonicalizeXML(noncanonXML)
	if err != nil {
		return nil, "", err
	}
	_, digest, err := GenerateDigest(canonXML)
	if err != nil {
		return nil, "", err
	}
	return canonXML, digest, nil
}

// Does not return canonXML
// Returns canon digest
func canonicalizeAndDigestBodyXML(body []byte) (string, error) {
	canon, err := CanonicalizeXML(body)
	if err != nil {
		return "", err
	}
	_, digest, err := GenerateDigest(canon)
	if err != nil {
		return "", err
	}
	return digest, nil
}

func GenerateSignedHeader(certificate *x509.Certificate, privateKey *rsa.PrivateKey, bodyReferenceURI string, bodyXML []byte) ([]byte, error) {
	// Generate URIs
	// These URIs should only have '#' in front of them when they are a
	// reference. It must exist once in the XML without
	// being referenced first, and should be unique.
	// Prepend the '#' reference where necessary.
	x509URI, err := GenerateSOAPURIWithPrefix("X509")
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
	signatureID, err := GenerateSOAPURIWithPrefix("SIG")
	if err != nil {
		return nil, err
	}

	ts, _, timestampDigest, err := GenerateTimestampAndDigest()
	if err != nil {
		return nil, err
	}

	bodyDigest, err := canonicalizeAndDigestBodyXML(bodyXML)
	if err != nil {
		return nil, err
	}

	_, x509Digest, err := GenerateDigest([]byte(certificate.Raw))
	if err != nil {
		return nil, err
	}

	signedInfoStruct, signedInfoXML, _, err := GenerateSignedInfoAndDigest(ts.ID, timestampDigest, bodyReferenceURI, bodyDigest, x509URI, x509Digest)
	if err != nil {
		return nil, err
	}

	signedInfoHash := sha512.New()
	_, err = signedInfoHash.Write(signedInfoXML)
	if err != nil {
		return nil, err
	}
	finalHash := signedInfoHash.Sum(nil)

	signedHash, err := privateKey.Sign(rand.Reader, finalHash, crypto.SHA512)
	if err != nil {
		return nil, err
	}

	securityHeader := header{
		Security: security{
			Wsse: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
			Wsu:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			BinarySecurityToken: binarySecurityToken{
				EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
				ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
				Text:         base64.StdEncoding.EncodeToString(certificate.Raw),
				ID:           x509URI,
			},
			Signature: signature{
				ID:         signatureID,
				Ds:         "http://www.w3.org/2000/09/xmldsig#",
				SignedInfo: signedInfoStruct,
				SignatureValue: signatureValue{
					Text: base64.StdEncoding.EncodeToString(signedHash),
				},
				KeyInfo: keyInfo{
					ID: keyInfoReferenceID,
					SecurityTokenReference: securityTokenReference{
						ID: securityTokenReferenceID,
						STReference: sTReference{
							URI:       "#" + x509URI,
							ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
						},
					},
				},
			},
			Timestamp: ts,
		},
	}

	// Canonicalizing the entire header here will be rejected by the server after successful TLS handshake. It must be put together earlier.

	marshaledHeader, err := xml.Marshal(securityHeader)
	if err != nil {
		return nil, err
	}

	return marshaledHeader, nil
}
