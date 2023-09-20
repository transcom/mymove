package trdm

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/beevik/etree"
)

func genEtreeEnvelope(cert x509.Certificate, key *rsa.PrivateKey, tableName string, timeToExpire time.Duration) (*etree.Element, error) {
	// Create envelope and set namespaces
	envelope := etree.NewElement("soap:Envelope")
	envelope.CreateAttr("xmlns:ret", "http://trdm/ReturnTableService")
	envelope.CreateAttr("xmlns:soap", "http://www.w3.org/2003/05/soap-envelope")

	// Create URIs
	bodyURI, err := GenerateSOAPURIWithPrefix("id")
	if err != nil {
		return nil, err
	}
	bstURI, err := GenerateSOAPURIWithPrefix("X509")
	if err != nil {
		return nil, err
	}
	tsURI, err := GenerateSOAPURIWithPrefix("TS")
	if err != nil {
		return nil, err
	}
	sigURI, err := GenerateSOAPURIWithPrefix("SIG")
	if err != nil {
		return nil, err
	}
	keyInfoURI, err := GenerateSOAPURIWithPrefix("SIG")
	if err != nil {
		return nil, err
	}
	securityTokenReferenceURI, err := GenerateSOAPURIWithPrefix("SIG")
	if err != nil {
		return nil, err
	}

	// Create empty Header element and children
	header, err := createHeader(envelope, cert, key, bstURI, sigURI, securityTokenReferenceURI, tsURI, bodyURI, keyInfoURI, timeToExpire)
	if err != nil {
		return nil, err
	}

	// Create empty body element and children
	body, err := createBody(envelope, bodyURI, tableName)
	if err != nil {
		return nil, err
	}

	// Create digests
	bst := header.FindElement("./wsse:Security/wsse:BinarySecurityToken[@wsu:Id='" + bstURI + "']")
	if bst == nil {
		return nil, fmt.Errorf("could not find element")
	}
	bstDigest, err := retrieveElemCanonicalDigest(bst)
	if err != nil {
		return nil, err
	}
	bstReference := header.FindElement("./wsse:Security/ds:Signature/ds:SignedInfo/ds:Reference[@URI='#" + bstURI + "']")
	if bstReference == nil {
		return nil, fmt.Errorf("could not find element")
	}
	bstDigestValueElement := bstReference.FindElement("./ds:DigestValue")
	bstDigestValueElement.SetText(base64.StdEncoding.EncodeToString(bstDigest))

	ts := header.FindElement("./wsse:Security/wsu:Timestamp[@wsu:Id='" + tsURI + "']")
	if ts == nil {
		return nil, fmt.Errorf("could not find element")
	}
	tsDigest, err := retrieveElemCanonicalDigest(ts)
	if err != nil {
		return nil, err
	}
	tsReference := header.FindElement("./wsse:Security/ds:Signature/ds:SignedInfo/ds:Reference[@URI='#" + tsURI + "']")
	if tsReference == nil {
		return nil, fmt.Errorf("could not find element")
	}
	tsDigestValueElement := tsReference.FindElement("./ds:DigestValue")
	tsDigestValueElement.SetText(base64.StdEncoding.EncodeToString(tsDigest))

	bodyDigest, err := retrieveElemCanonicalDigest(body)
	if err != nil {
		return nil, err
	}
	bodyReference := header.FindElement("./wsse:Security/ds:Signature/ds:SignedInfo/ds:Reference[@URI='#" + bodyURI + "']")
	if bodyReference == nil {
		return nil, fmt.Errorf("could not find element")
	}
	bodyDigestValueElement := bodyReference.FindElement("./ds:DigestValue")
	bodyDigestValueElement.SetText(base64.StdEncoding.EncodeToString(bodyDigest))

	// Now, all three digests have been formed. Now we need to sign the signedInfo element
	signedInfoElement := header.FindElement("./wsse:Security/ds:Signature/ds:SignedInfo")
	var signedInfoBytes bytes.Buffer
	signedInfoElement.WriteTo(&signedInfoBytes, &etree.NewDocument().WriteSettings)
	signedInfoDigest, _, err := GenerateDigest(signedInfoBytes.Bytes())
	if err != nil {
		return nil, err
	}
	sigValue, err := key.Sign(rand.Reader, signedInfoDigest, crypto.SHA512)
	if err != nil {
		return nil, err
	}

	signatureValueElement := header.FindElement("./wsse:Security/ds:Signature/ds:SignatureValue")
	signatureValueElement.SetText(base64.StdEncoding.EncodeToString(sigValue))

	// Now that the signature is done, the envelope should be complete.

	return envelope, nil
}

func createHeader(parent *etree.Element, cert x509.Certificate, key *rsa.PrivateKey, bstURI string, sigURI string, securityTokenReferenceURI string, tsURI string, bodyURI string, keyInfoURI string, timeToExpire time.Duration) (*etree.Element, error) {
	header := parent.CreateElement("soap:Header")
	// Add base security element
	security := header.CreateElement("wsse:Security")
	security.CreateAttr("xmlns:wsse", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd")
	security.CreateAttr("xmlns:wsu", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd")
	// Create base BST element
	// This element will populate with text as it does not have a digest value child element
	createBinarySecurityTokenElement(security, bstURI, cert)
	// Create signature element with empty hash values
	err := createSignatureElement(security, sigURI, keyInfoURI, securityTokenReferenceURI, cert, tsURI, bstURI, bodyURI, key)
	if err != nil {
		return nil, err
	}

	// Create timestamp element with time to expire value
	// Timestamp comes after signature element
	_, err = createTimestampElement(security, timeToExpire, tsURI)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func createBody(parent *etree.Element, URI string, tableName string) (*etree.Element, error) {
	if tableName == "" {
		return nil, fmt.Errorf("tablename provided to soap body left blank")
	}
	body := parent.CreateElement("soap:Body")
	body.CreateAttr("wsu:Id", URI)
	body.CreateAttr("xmlns:wsu", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd")
	getLastTableUpdateRequestElement := body.CreateElement("ret:getLastTableUpdateRequestElement")
	physicalName := getLastTableUpdateRequestElement.CreateElement("ret:physicalName")
	physicalName.SetText(tableName)
	return body, nil
}

// Creates KeyInfo element and its children
func createKeyInfoElement(parent *etree.Element, keyInfoURI string, strURI string, x509URI string) *etree.Element {
	keyInfo := parent.CreateElement("ds:KeyInfo")
	keyInfo.CreateAttr("Id", keyInfoURI)

	securityTokenReference := keyInfo.CreateElement("wsse:SecurityTokenReference")
	securityTokenReference.CreateAttr("wsu:Id", strURI)

	reference := securityTokenReference.CreateElement("wsse:Reference")
	reference.CreateAttr("URI", "#"+x509URI)
	reference.CreateAttr("ValueType", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3")

	return keyInfo
}

func createBinarySecurityTokenElement(parent *etree.Element, bstURI string, cert x509.Certificate) *etree.Element {
	binarySecurityToken := parent.CreateElement("wsse:BinarySecurityToken")
	binarySecurityToken.CreateAttr("EncodingType", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary")
	binarySecurityToken.CreateAttr("ValueType", "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3")
	binarySecurityToken.CreateAttr("wsu:Id", bstURI)
	binarySecurityToken.SetText(base64.StdEncoding.EncodeToString(cert.Raw))
	return binarySecurityToken
}

// Creeates the Signature element and its children
func createSignatureElement(parent *etree.Element, sigURI string, keyInfoURI string, strURI string, cert x509.Certificate, tsURI string, bstURI string, bodyURI string, key *rsa.PrivateKey) error {
	signatureElement := parent.CreateElement("ds:Signature")
	signatureElement.CreateAttr("Id", sigURI)
	signatureElement.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")
	// Create empty signed info element (No digest values)
	_, err := createSignedInfoElement(signatureElement, tsURI, bstURI, bodyURI)
	if err != nil {
		return err
	}
	// Now, our signed info exists and the references have empty digest values
	// At this point, the only value entered is the binary security token (Public cert)

	// Now create the empty signature value element
	signatureElement.CreateElement("ds:SignatureValue")

	// After signature comes KeyInfo, no digest values inside
	createKeyInfoElement(signatureElement, keyInfoURI, strURI, bstURI)

	// Signature complete with empty digests

	return nil
}

// Returns the signature of the canonicalized element
func signElement(element *etree.Element, key *rsa.PrivateKey) ([]byte, error) {
	// Get the digest of the canon element
	digest, err := retrieveElemCanonicalDigest(element)
	if err != nil {
		return nil, err
	}
	// Sign the canon digest
	signature, err := key.Sign(rand.Reader, digest, crypto.SHA512)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Creates the timestamp element and applies an expiration date to the envelope.
// The expiration is proven by its digest value.
func createTimestampElement(parent *etree.Element, timeToExpire time.Duration, tsURI string) (*etree.Element, error) {
	timestamp := parent.CreateElement("wsu:Timestamp")
	timestamp.CreateAttr("wsu:Id", tsURI)

	now := time.Now()
	layout := "2006-01-02T15:04:05.000Z07:00" // Allows for millisecond precision
	formattedNow := now.Format(layout)
	future := now.Add(timeToExpire)
	formattedFuture := future.Format(layout)

	// Add 'creates' and 'expires' elements
	created := timestamp.CreateElement("wsu:Created")
	created.SetText(formattedNow)

	expires := timestamp.CreateElement("wsu:Expires")
	expires.SetText(formattedFuture)

	return timestamp, nil
}

// Returns the canonical SHA512 digest of the entire provided element, not just the content
func retrieveElemCanonicalDigest(elem *etree.Element) ([]byte, error) {
	var buffer bytes.Buffer
	writeSettings := &etree.WriteSettings{
		CanonicalText:    true,
		CanonicalEndTags: true,
		CanonicalAttrVal: true,
	}

	elem.WriteTo(&buffer, writeSettings)

	hash, _, err := GenerateDigest(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Creates the SignedInfo element and its children
func createSignedInfoElement(parent *etree.Element, tsURI string, bstURI string, bodyURI string) (*etree.Element, error) {
	signedInfo := parent.CreateElement("ds:SignedInfo")

	// Create canonicalization method
	// No digest values inside
	createCanonicalizationMethodElement(signedInfo)

	// Create signature method
	// No digest values inside
	createSignatureMethod(signedInfo)

	// Create timestamp reference
	// ! Empty digest inside
	err := createReferenceElement(signedInfo, "#"+tsURI, "wsse ret soap", "http://www.w3.org/2001/04/xmlenc#sha512")
	if err != nil {
		return nil, err
	}

	// Create BinarySecurityToken reference
	// ! Empty digest inside
	err = createReferenceElement(signedInfo, "#"+bstURI, "", "http://www.w3.org/2001/04/xmlenc#sha512")
	if err != nil {
		return nil, err
	}

	// Create body reference
	// ! Empty digest inside
	err = createReferenceElement(signedInfo, "#"+bodyURI, "ret", "http://www.w3.org/2001/04/xmlenc#sha512")
	if err != nil {
		return nil, err
	}

	return signedInfo, nil
}

func createCanonicalizationMethodElement(parent *etree.Element) {
	canonicalizationElement := parent.CreateElement("ds:CanonicalizationMethod")
	canonicalizationElement.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")
	inclusiveNameSpaces := canonicalizationElement.CreateElement("ec:InclusiveNamespaces")
	inclusiveNameSpaces.CreateAttr("PrefixList", "ret soap")
	inclusiveNameSpaces.CreateAttr("xmlns:ec", "http://www.w3.org/2001/10/xml-exc-c14n#")
}

func createSignatureMethod(parent *etree.Element) {
	sigMethod := parent.CreateElement("ds:SignatureMethod")
	sigMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512")
}

// Creates Reference element and its children, but with an empty digest value
func createReferenceElement(parent *etree.Element, uri string, prefixList string, digestAlgorithm string) error {
	reference := parent.CreateElement("ds:Reference")
	reference.CreateAttr("URI", uri)

	// Add Transforms
	transforms := reference.CreateElement("ds:Transforms")
	transform := transforms.CreateElement("ds:Transform")
	transform.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#") // Transform algorithm does not change

	// Add InclusiveNamespaces to Transform
	inclusiveNs := transform.CreateElement("ec:InclusiveNamespaces")
	inclusiveNs.CreateAttr("PrefixList", prefixList)
	inclusiveNs.CreateAttr("xmlns:ec", "http://www.w3.org/2001/10/xml-exc-c14n#")

	// Add DigestMethod
	digestMethod := reference.CreateElement("ds:DigestMethod")
	digestMethod.CreateAttr("Algorithm", digestAlgorithm)

	reference.CreateElement("ds:DigestValue")
	// Do not add digest value at this time

	return nil
}
