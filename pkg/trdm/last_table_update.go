package trdm

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

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
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:ret="http://ReturnTablePackage/">
   <soapenv:Header/>
   <soapenv:Body>
      <ret:getLastTableUpdateRequestElement>
         <ret:physicalName>ACFT</ret:physicalName>
      </ret:getLastTableUpdateRequestElement>
   </soapenv:Body>
</soapenv:Envelope>

SOAP Response:
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
   <soap:Body>
      <getLastTableUpdateResponseElement xmlns="http://ReturnTablePackage/">
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

type GetLastTableUpdateRequestElement struct {
	XMLName       string `xml:"ret:getLastTableUpdateReqElement"`
	PhysicalName  string `xml:"ret:physicalName"`
	soapClient    SoapCaller
	securityToken *x509.Certificate
	privateKey    *rsa.PrivateKey
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

func NewTRDMGetLastTableUpdate(physicalName string, securityToken *x509.Certificate, privateKey *rsa.PrivateKey, soapClient SoapCaller) GetLastTableUpdater {
	return &GetLastTableUpdateRequestElement{
		PhysicalName:  physicalName,
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
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ret":     "http://ReturnTablePackage/",
	})

	params := GetLastTableUpdateRequestElement{
		PhysicalName: physicalName,
	}
	marshaledBody, marshalEr := xml.Marshal(params)
	if marshalEr != nil {
		return marshalEr
	}
	signedHeader, headerSigningError := GenerateSignedHeader(d.securityToken, d.privateKey)
	if headerSigningError != nil {
		return headerSigningError
	}
	newParams := gosoap.Params{
		"header": signedHeader,
		"body":   marshaledBody,
	}
	err := lastTableUpdateSoapCall(d, newParams, appCtx)
	if err != nil {
		return fmt.Errorf("request error: %s", err.Error())
	}
	return nil
}

// Makes Soap Request for Last Table Update. If successful call GetTable to update TRDM table.
//   - *GetLastTableUpdateRequestElement - request elements
//   - params - Created SOAP element (Header and Body)
//   - appCtx - Application context
//   - returns error
func lastTableUpdateSoapCall(d *GetLastTableUpdateRequestElement, params gosoap.Params, appCtx appcontext.AppContext) error {
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
		getTable := NewGetTable(d.PhysicalName, d.securityToken, d.privateKey, d.soapClient)
		getTableErr := getTable.GetTable(appCtx, d.PhysicalName, r.LastUpdate)
		if getTableErr != nil {
			return fmt.Errorf("getTable error: %s", getTableErr.Error())
		}
	}

	appCtx.Logger().Debug("getLastTableUpdate result", zap.Any("processRequestResponse", r))
	return nil
}
func StartLastTableUpdateCron(appCtx appcontext.AppContext, certificate *x509.Certificate, privateKey *rsa.PrivateKey, physicalName string, soapCaller SoapCaller) error {
	cron := cron.New()

	cronTask := func() {
		err := NewTRDMGetLastTableUpdate(physicalName, certificate, privateKey, soapCaller).GetLastTableUpdate(appCtx, physicalName)
		if err != nil {
			fmt.Println("Error in lastTableUpdate cron task: ", err)
		}
	}

	res, err := cron.AddFunc("@every 24h00m00s", cronTask)
	if err != nil {
		return fmt.Errorf("Error adding cron task: %s, %v", err.Error(), res)
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

	trdmWSDL := v.GetString(cli.TRDMApiWSDLFlag)
	trdmURL := v.GetString(cli.TRDMApiURLFlag)
	soapClient, err := gosoap.SoapClient(trdmWSDL, httpClient)
	if err != nil {
		return fmt.Errorf("unable to create SOAP client: %w", err)
	}
	soapClient.URL = trdmURL

	x509CertString := v.GetString(cli.TRDMx509Cert)
	certificate, err := x509.ParseCertificate([]byte(x509CertString))
	if err != nil {
		return err
	}

	privateKeyString := v.GetString(cli.TRDMx509PrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyString))
	if err != nil {
		return err
	}

	getLastTableUpdateTACErr := NewTRDMGetLastTableUpdate(transportationAccountingCode, certificate, privateKey, soapClient).GetLastTableUpdate(appCtx, transportationAccountingCode)
	getLastTableUpdateLOAErr := NewTRDMGetLastTableUpdate(lineOfAccounting, certificate, privateKey, soapClient).GetLastTableUpdate(appCtx, lineOfAccounting)
	if getLastTableUpdateLOAErr != nil {
		return getLastTableUpdateLOAErr
	}
	if getLastTableUpdateTACErr != nil {
		return getLastTableUpdateTACErr
	}

	cronErrTAC := StartLastTableUpdateCron(appCtx, certificate, privateKey, transportationAccountingCode, soapClient)
	cronErrLOA := StartLastTableUpdateCron(appCtx, certificate, privateKey, lineOfAccounting, soapClient)

	if cronErrLOA != nil {
		return cronErrLOA
	}
	if cronErrTAC != nil {
		return cronErrTAC
	}
	return nil
}
