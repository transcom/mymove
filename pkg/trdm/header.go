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
	"fmt"
	"time"

	"github.com/beevik/etree"
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
	XMLName                xml.Name               `xml:"ds:KeyInfo"`
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
	XMLName      xml.Name `xml:"wsse:BinarySecurityToken"`
	Text         string   `xml:",chardata"`
	EncodingType string   `xml:"EncodingType,attr"`
	ValueType    string   `xml:"ValueType,attr"`
	ID           string   `xml:"wsu:Id,attr"`
}

type signedInfo struct {
	XMLName                xml.Name               `xml:"ds:SignedInfo"`
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
	XMLName xml.Name `xml:"wsu:Timestamp"`
	ID      string   `xml:"wsu:Id,attr"`
	Created string   `xml:"wsu:Created"`
	Expires string   `xml:"wsu:Expires"`
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

// Returns canon security element
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

// ! Future readers: Keep in mind the vast majority of this code is currently inoperable. The bottom of this function has an additional function call that overwrites all
// ! of the witnessed functions here. There have been multiple implementations attempted and strings printed to utilize in `curl` requests for testing - however none have
// ! worked. See etree_imp.go for the latest attempt, you will find many artifacts in this current file referencing marshaled XML generation as well as
// ! manual string concatenation and canonicalization.
func GenerateSignedHeader(certificate *x509.Certificate, privateKey *rsa.PrivateKey, bodyReferenceURI string, bodyXML []byte) ([]byte, error) {

	// ! WARNING !
	// ! Read the comment above this function before proceeding.

	// ! Read the comment above this function before proceeding.

	// ! Read the comment above this function before proceeding.

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

	bodyDigest, err := canonicalizeAndDigestBodyXML(bodyXML)
	if err != nil {
		return nil, err
	}

	// Old generation of headerXML
	_, err = generateHeaderXML(*certificate, privateKey, x509URI, bodyReferenceURI, bodyDigest, keyInfoReferenceID, securityTokenReferenceID)
	if err != nil {
		return nil, err
	}

	var envBytes bytes.Buffer
	envelope, err := genEtreeEnvelope(*certificate, privateKey, "TRNSPRTN_ACNT", time.Millisecond*5000)
	if err != nil {
		return nil, err
	}
	envelope.WriteTo(&envBytes, &etree.WriteSettings{CanonicalText: true, CanonicalEndTags: true, CanonicalAttrVal: true})
	fmt.Printf("\n My canon envelope from etree: \n %s", string(envBytes.Bytes()))
	return envBytes.Bytes(), nil
}

func generateHeaderXML(cert x509.Certificate, key *rsa.PrivateKey, certURI string, bodyURI string, bodyDigest string, keyInfoURI string, strURI string) (string, error) {
	bstXML, bstDigest, err := generateBinarySecurityToken(cert, certURI)
	if err != nil {
		return "", err
	}

	ts, timestampXML, timestampDigest, err := GenerateTimestampAndDigest()
	if err != nil {
		return "", err
	}

	// ! Cert digest is the digest of the x509 XML element, not just a digest of the certificate. Aka the binarySecurityToken
	_, signedInfoXML, _, err := GenerateSignedInfoAndDigest(ts.ID, timestampDigest, bodyURI, bodyDigest, certURI, bstDigest)
	if err != nil {
		return "", err
	}

	keyInfoXML, _, err := generateKeyInfo(keyInfoURI, strURI, certURI)
	if err != nil {
		return "", err
	}

	sigURI, err := GenerateSOAPURIWithPrefix("SIG")
	if err != nil {
		return "", err
	}
	sigXML, sigHash, err := generateSignature(signedInfoXML, keyInfoXML, keyInfoURI, sigURI, key)
	if err != nil {
		return "", err
	}

	headerXML := fmt.Sprintf(`<soap:Header>
<wsse:Security
xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
%s
%s
%s
</wsse:Security>
</soap:Header>`, bstXML, sigXML, timestampXML)
	fmt.Println(headerXML)

	signedInfoHash := sha512.New()
	_, err = signedInfoHash.Write([]byte(signedInfoXML))
	if err != nil {
		return "", err
	}
	finalHash := signedInfoHash.Sum(nil)

	rsaCert := cert.PublicKey.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(rsaCert, crypto.SHA512, finalHash, sigHash)
	if err != nil {
		return "", err
	}

	return headerXML, nil
}

// Returns
// - XML
// - Digest
// - Error
func generateBinarySecurityToken(cert x509.Certificate, x509URI string) ([]byte, string, error) {
	// Do not return struct, it should no longer be used after the XML has been generated
	t := binarySecurityToken{
		EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
		ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
		Text:         base64.StdEncoding.EncodeToString(cert.Raw),
		ID:           x509URI,
	}
	xmlByte, digest, err := GenerateSecurityElement(t)
	if err != nil {
		return nil, "", err
	}

	return xmlByte, digest, nil
}

// Returns
// - Canon XML
// - Digest
// - Error
func generateSignature(signedInfoXML []byte, keyInfoXML []byte, keyInfoURI string, sigURI string, key *rsa.PrivateKey) (string, []byte, error) {
	// Do not return struct, it should no longer be used after the XML has been generated
	signedHash, err := signXML(signedInfoXML, key)
	if err != nil {
		return "", nil, err
	}

	fmt.Printf("\nThis is the XML I have signed\n%s\n", string(signedInfoXML))

	sigXML := fmt.Sprintf(`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="%s">
%s
<ds:SignatureValue>%s</ds:SignatureValue>
%s
</ds:Signature>`, sigURI, string(signedInfoXML), base64.StdEncoding.EncodeToString(signedHash), string(keyInfoXML))

	// Do not canonicalize, the child elements already have canonicalized XML. Recanonicalizing will break the canonicalization

	return sigXML, signedHash, nil
}

// Returns
// - Signed XML hash
// - Error
func signXML(xml []byte, key *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	_, err := hash.Write(xml)
	if err != nil {
		return nil, err
	}
	finalHash := hash.Sum(nil)
	signedHash, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, finalHash)
	if err != nil {
		return nil, err
	}

	return signedHash, nil
}

// Returns
// - XML
// - Digest
// - Error
func generateKeyInfo(keyinfoURI string, strURI string, certURI string) ([]byte, string, error) {
	// str = securityTokenReference
	// Do not return struct, it should no longer be used after the XML has been generated
	t := keyInfo{
		ID: keyinfoURI,
		SecurityTokenReference: securityTokenReference{
			ID: strURI,
			STReference: sTReference{
				URI:       "#" + certURI,
				ValueType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
			},
		},
	}
	xmlByte, digest, err := GenerateSecurityElement(t)
	if err != nil {
		return nil, "", err
	}

	return xmlByte, digest, nil
}
