package trdm

import (
	"crypto/rsa"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
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
	securityToken string
	privateKey    *rsa.PrivateKey
}
type GetLastTableUpdateResponseElement struct {
	XMLName    xml.Name `xml:"getLastTableUpdateResponseElement"`
	LastUpdate string   `xml:"lastUpdate"`
	Status     struct {
		StatusCode string `xml:"statusCode"`
		DateTime   string `xml:"dateTime"`
	} `xml:"status"`
}

func NewTRDMGetLastTableUpdate(physicalName string, securityToken string, privateKey *rsa.PrivateKey, soapClient SoapCaller) GetLastTableUpdater {
	return &GetLastTableUpdateRequestElement{
		PhysicalName:  physicalName,
		soapClient:    soapClient,
		securityToken: securityToken,
		privateKey:    privateKey,
	}

}

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
	err := lastTableUpdateSoapCall(d, newParams, appCtx, physicalName)
	if err != nil {
		return fmt.Errorf("Request error: %s", err.Error())
	}
	return nil
}

func lastTableUpdateSoapCall(d *GetLastTableUpdateRequestElement, params gosoap.Params, appCtx appcontext.AppContext, physicalName string) error {
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
		getTable := NewGetTable(physicalName, d.securityToken, d.privateKey, d.soapClient)
		getTableErr := getTable.GetTable(appCtx, physicalName, r.LastUpdate)
		if getTableErr != nil {
			return fmt.Errorf("getTable error: %s", getTableErr.Error())
		}
	}

	appCtx.Logger().Debug("getLastTableUpdate result", zap.Any("processRequestResponse", r))
	return nil
}
func StartLastTableUpdateCron(appCtx appcontext.AppContext, physicalName string) {
	cron := cron.New()

	cronTask := func() {
		err := NewTRDMGetLastTableUpdate("", "", nil, nil).GetLastTableUpdate(appCtx, physicalName)
		if err != nil {
			fmt.Println("Error in lastTableUpdate cron task: ", err)
		}
	}

	res, err := cron.AddFunc("@every 24h00m00s", cronTask)
	if err != nil {
		fmt.Println("Error adding cron task: ", err, res)
	}
	cron.Start()
}

func LastTableUpdate() {

	tableName := ""

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// DB connection
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		fmt.Println("Error getting DB connection: ", err)
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	err2 := NewTRDMGetLastTableUpdate("", "", nil, nil).GetLastTableUpdate(appCtx, tableName)
	if err2 != nil {
		fmt.Println("Error executing GetLastTableUpdate: ", err)
	}

	StartLastTableUpdateCron(appCtx, tableName)
}
