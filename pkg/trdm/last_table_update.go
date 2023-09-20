package trdm

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/beevik/etree"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
)

/*******************************************

The struct/method getLastTableUpdate:GetLastTableUpdate implements the service TRDM.

This method GetLastTableUpdate sends a SOAP request to TRDM to get the last table update.
This code is using the gosoap lib https://github.com/tiaguinho/gosoap

SOAP Request:
<soapenv:Envelope xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope"
xmlns:ret="http://trdm/ReturnTableService">
   <soapenv:Header/>
   <soapenv:Body>
      <ret:getLastTableUpdateRequestElement>
         <ret:physicalName>ACFT</ret:physicalName>
      </ret:getLastTableUpdateRequestElement>
   </soapenv:Body>
</soapenv:Envelope>

SOAP Response:
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
   <soap:Body>
      <getLastTableUpdateResponseElement xmlns="http://trdm/ReturnTableService">
         <lastUpdate>2020-01-27T16:14:20.000Z</lastUpdate>
         <status>
            <statusCode>Successful</statusCode>
            <dateTime>2020-01-27T20:18:34.226Z</dateTime>
         </status>
      </getLastTableUpdateResponseElement>
   </soap:Body>
</soap:Envelope>
 *******************************************/

const successfulStatusCode = "Successful"

// Date/time value is used in conjunction with the contentUpdatedSinceDateTime column in the getTable method.
type GetLastTableUpdater interface {
	GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string) error
}

// ReturnContent is to sometimes be left blank intentionally
// For example. it is required for getTable but should not be provided
// in getLastTableUpdate
type getTableTRDMInput struct {
	PhysicalName  string `xml:"ret:physicalName"`
	ReturnContent string `xml:"ret:returnContent"`
}

type lastTableUpdateTRDMInput struct {
	PhysicalName string `xml:"ret:physicalName"`
}

type input struct {
	TRDMInput getTableTRDMInput `xml:"ret:TRDM"`
}

// ret:input is only for getTable, not getLastTable
// Directly embed trdmInput here
type lastTableUpdateRequestElement struct {
	lastTableUpdateTRDMInput
}

type GetLastTableUpdateRequestElement struct {
	XMLName                       xml.Name                      `xml:"soap:Body"`
	ID                            string                        `xml:"wsu:Id,attr"`
	Wsu                           string                        `xml:"xmlns:wsu,attr"`
	LastTableUpdateRequestElement lastTableUpdateRequestElement `xml:"ret:getLastTableUpdateRequestElement"`
	soapClient                    SoapCaller
	securityToken                 *x509.Certificate
	privateKey                    *rsa.PrivateKey
}
type GetLastTableUpdateResponseElement struct {
	XMLName    xml.Name  `xml:"getLastTableUpdateResponseElement"`
	LastUpdate time.Time `xml:"lastUpdate"`
	Status     status    `xml:"status"`
}

type Status struct {
	StatusCode string `xml:"statusCode"`
	DateTime   string `xml:"dateTime"`
}

func NewTRDMGetLastTableUpdate(physicalName string, bodyID string, securityToken *x509.Certificate, privateKey *rsa.PrivateKey, soapClient SoapCaller) GetLastTableUpdater {
	return &GetLastTableUpdateRequestElement{
		LastTableUpdateRequestElement: lastTableUpdateRequestElement{
			// Remember, leaving returnContent in trdmInput blank is intentional
			lastTableUpdateTRDMInput{
				PhysicalName: physicalName,
			},
		},
		ID:            bodyID,
		Wsu:           "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
		soapClient:    soapClient,
		securityToken: securityToken,
		privateKey:    privateKey,
	}
}

// FetchAllTACRecords queries and fetches all transportation_accounting_codes
func FetchAllTACRecords(appcontext appcontext.AppContext) ([]models.TransportationAccountingCode, error) {
	var tacCodes []models.TransportationAccountingCode
	query := `SELECT * FROM transportation_accounting_codes`
	err := appcontext.DB().RawQuery(query).All(&tacCodes)
	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil

}

// Sets up GetLastTableUpdate Soap Request
//   - appCtx - Application Context.
//   - physicalName - TableName
//   - Generates custom soap envelope, soap body, soap header.
//   - Returns Error
func (d *GetLastTableUpdateRequestElement) GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string) error {
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://www.w3.org/2003/05/soap-envelope",
		"xmlns:ret":     "http://trdm/ReturnTableService",
	})
	bodyID, err := GenerateSOAPURIWithPrefix("id")
	if err != nil {
		return err
	}
	// Needs to be nested in test:input ret:TRDM
	params := GetLastTableUpdateRequestElement{
		ID:  bodyID,
		Wsu: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
		// Remember, leaving returnContent in trdmInput blank is intentional
		LastTableUpdateRequestElement: lastTableUpdateRequestElement{
			lastTableUpdateTRDMInput{
				PhysicalName: physicalName,
			},
		},
	}
	marshaledBody, marshalEr := xml.Marshal(params)
	if marshalEr != nil {
		return marshalEr
	}
	signedHeader, err := GenerateSignedHeader(d.securityToken, d.privateKey, bodyID, marshaledBody)
	if err != nil {
		return err
	}
	newParams := gosoap.Params{
		"header": signedHeader,
		"body":   marshaledBody,
	}

	// ! This is being utilized because the vscode debugger does not support
	// ! strings above 64 bytes
	// Start printing
	headerStr := string(newParams["header"].([]byte))
	canonBodyXML, err := CanonicalizeXML(newParams["body"].([]byte))
	if err != nil {
		return err
	}
	//bodyStr := string(newParams["body"].([]byte))
	//bodyStr := string(marshaledBody)
	soapEnvelope := fmt.Sprintf(
		`<soap:Envelope xmlns:ret="http://trdm/ReturnTableService" xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
%s
%s
</soap:Envelope>`,
		headerStr, canonBodyXML,
	)

	fmt.Println(soapEnvelope)

	if err = verifySignedInfoXML(soapEnvelope); err != nil {
		return err
	}
	/*
		err = verifyXML(soapEnvelope)
		if err != nil {
			return err
		}
	*/
	// End printing
	err = lastTableUpdateSoapCall(d, newParams, appCtx)
	if err != nil {
		return fmt.Errorf("request error: %s", err.Error())
	}
	return nil
}
func verifySignedInfoXML(xmlContent string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlContent); err != nil {
		return err
	}
	signedInfoElem := doc.FindElement("//SignedInfo")
	if signedInfoElem == nil {
		return fmt.Errorf("could not find signed info elem")
	}
	doc.SetRoot(signedInfoElem)

	return nil
}
func verifyXML(xmlContent string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlContent); err != nil {
		return err
	}

	// Locate the SignatureValue and extract the signature
	signatureElement := doc.FindElement("//ds:SignatureValue")
	if signatureElement == nil {
		return fmt.Errorf("could not find signature element")
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signatureElement.Text())
	if err != nil {
		return err
	}

	// Locate the certificate and extract it
	certElement := doc.FindElement("//wsse:BinarySecurityToken")
	if certElement == nil {
		return fmt.Errorf("could not find x509 cert")
	}

	decodedCert, err := base64.StdEncoding.DecodeString(certElement.Text())
	if err != nil {
		return err
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate([]byte(decodedCert))
	if err != nil {
		return fmt.Errorf("Failed to parse certificate PEM")
	}

	// Compute the hash of the SignedInfo element
	signedInfoElement := doc.FindElement("//SignedInfo")
	if signedInfoElement == nil {
		fmt.Errorf("SignedInfo not found")
	}

	headerElem := doc.FindElement("//Header")
	headerElem.CreateAttr("ret", "http://trdm/ReturnTableService")
	headerElem.CreateAttr("soap", "http://www.w3.org/2003/05/soap-envelope")
	// signedInfoElement.CreateAttr("xmlns:ds", "http://www.w3.org/2000/09/xmldsig#")
	// signedInfoElement.CreateAttr("Id", "SIG-8796c3e3fa1d1183")

	// Create a new document to write just the SignedInfo element
	docFragment := etree.NewDocument()
	docFragment.SetRoot(signedInfoElement.Copy())

	strBuffer := &bytes.Buffer{}
	if _, err := docFragment.WriteTo(strBuffer); err != nil {
		return err
	}

	// canonicalSignedInfo, err := signedInfoElement.WriteToString()
	// if err != nil {
	// 	return err
	// }

	hashed := sha512.Sum512([]byte(strBuffer.Bytes()))
	// Verify the signature
	if err := rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA512, hashed[:], decodedSignature); err != nil {
		return err
	}

	return nil
}

// Makes Soap Request for Last Table Update. If successful call GetTable to update TRDM table.
//   - *GetLastTableUpdateRequestElement - request elements
//   - params - Created SOAP element (Header and Body)
//   - appCtx - Application context
//   - returns error
func lastTableUpdateSoapCall(d *GetLastTableUpdateRequestElement, params gosoap.Params, appCtx appcontext.AppContext) error {
	// This will hit the ?WSDL endpoint with the marshaled body and header
	res, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return fmt.Errorf("call error: %s", err.Error())
	}

	var r GetLastTableUpdateResponseElement
	err = res.Unmarshal(&r)
	if err != nil {
		return fmt.Errorf("unmarshal error: %s", err.Error())
	}

	if r.Status.StatusCode == successfulStatusCode {
		getTable := NewGetTable(d.LastTableUpdateRequestElement.PhysicalName, d.securityToken, d.privateKey, d.soapClient)
		getTableErr := getTable.GetTable(appCtx, d.LastTableUpdateRequestElement.PhysicalName, r.LastUpdate)
		if getTableErr != nil {
			return fmt.Errorf("getTable error: %s", getTableErr.Error())
		}
	}

	appCtx.Logger().Debug("getLastTableUpdate result", zap.Any("processRequestResponse", r))
	return nil
}
func StartLastTableUpdateCron(appCtx appcontext.AppContext, certificate *x509.Certificate, publicPem *pem.Block, privateKey *rsa.PrivateKey, physicalName string, soapCaller SoapCaller) error {
	cron := cron.New()
	bodyID, err := GenerateSOAPURIWithPrefix("#id")
	if err != nil {
		return err
	}
	cronTask := func() {
		err = NewTRDMGetLastTableUpdate(physicalName, bodyID, certificate, privateKey, soapCaller).GetLastTableUpdate(appCtx, physicalName)
		if err != nil {
			fmt.Println("Error in lastTableUpdate cron task: ", err)
		}
	}

	res, err := cron.AddFunc("@every 24h00m00s", cronTask)
	if err != nil {
		return fmt.Errorf("error adding cron task: %s, %v", err.Error(), res)
	}
	cron.Start()
	return nil
}

func LastTableUpdate(v *viper.Viper, tlsConfig *tls.Config) error {
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	// DB connection
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		return err
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: tr, Timeout: time.Duration(30) * time.Second}

	trdmWSDL := v.GetString(cli.TRDMApiReturnTableV7WSDLFlag)
	trdmURL := v.GetString(cli.TRDMApiURLFlag)
	soapClient, err := gosoap.SoapClient(trdmWSDL, httpClient)
	if err != nil {
		return fmt.Errorf("unable to create SOAP client: %w", err)
	}
	soapClient.URL = trdmURL

	// Base64 DER encoded x509 certificate
	x509CertString := v.GetString(cli.MoveMilDoDTLSCertFlag)

	publicPem, rest := pem.Decode([]byte(x509CertString))
	if len(rest) != 0 {
		return fmt.Errorf("unable to properly decode public key, something is leftover: %w", err)
	}

	certificate, err := x509.ParseCertificate(publicPem.Bytes)
	if err != nil {
		return err
	}

	// Currently the private key is in PKSC1 format
	privateKeyString := v.GetString(cli.MoveMilDoDTLSKeyFlag)

	privatePem, rest := pem.Decode([]byte(privateKeyString))
	if len(rest) != 0 {
		return fmt.Errorf("unable to properly decode private key, something is leftover: %w", err)
	}

	// Declare key here for PKCS1 and PKCS8 handling
	var key interface{}

	// ! If this line is failing, it's because app-devlocal uses PKCS8 and app-stg and app-prd use PKCS1
	// Try to parse as PKCS1
	key, err = x509.ParsePKCS1PrivateKey([]byte(privatePem.Bytes))
	if err != nil {
		// Try to parse as PKCS8
		var pkcs8err error
		key, pkcs8err = x509.ParsePKCS8PrivateKey([]byte(privatePem.Bytes))
		if pkcs8err != nil {
			return fmt.Errorf("failed parsing private keys, \n PKCS1 err: %s \n PKCS8 err: %s", err, pkcs8err)
		}
	}

	// Type assert
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("key is not of type *rsa.PrivateKey, got: %T", key)
	}

	tacBodyID, err := GenerateSOAPURIWithPrefix("#id")
	if err != nil {
		return err
	}
	loaBodyID, err := GenerateSOAPURIWithPrefix("#id")
	if err != nil {
		return err
	}
	getLastTableUpdateTACErr := NewTRDMGetLastTableUpdate(transportationAccountingCode, tacBodyID, certificate, rsaKey, soapClient).GetLastTableUpdate(appCtx, transportationAccountingCode)
	getLastTableUpdateLOAErr := NewTRDMGetLastTableUpdate(lineOfAccounting, loaBodyID, certificate, rsaKey, soapClient).GetLastTableUpdate(appCtx, lineOfAccounting)
	if getLastTableUpdateLOAErr != nil {
		return getLastTableUpdateLOAErr
	}
	if getLastTableUpdateTACErr != nil {
		return getLastTableUpdateTACErr
	}

	cronErrTAC := StartLastTableUpdateCron(appCtx, certificate, publicPem, rsaKey, transportationAccountingCode, soapClient)
	cronErrLOA := StartLastTableUpdateCron(appCtx, certificate, publicPem, rsaKey, lineOfAccounting, soapClient)

	if cronErrLOA != nil {
		return cronErrLOA
	}
	if cronErrTAC != nil {
		return cronErrTAC
	}
	return nil
}
