package trdm

import (
	"bytes"
	"crypto/rsa"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/parser/loa"
	"github.com/transcom/mymove/pkg/parser/tac"
)

// <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ret="http://ReturnTablePackage/">
//    <soapenv:Header/>
//    <soapenv:Body>
//       <ret:getTableRequestElement>
//          <ret:input>
//             <ret:TRDM>
//                <ret:physicalName>ACFT</ret:physicalName>
//                <ret:returnContent>true</ret:returnContent>
//             </ret:TRDM>
//          </ret:input>
//       </ret:getTableRequestElement>
//    </soapenv:Body>
// </soapenv:Envelope>

// SOAP Response:
// <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
//    <soap:Body>
//       <getTableResponseElement xmlns="http://ReturnTablePackage/">
//          <output>
//             <TRDM>
//                <status>
//                   <rowCount>28740</rowCount>
//                   <statusCode>Successful</statusCode>
//                   <dateTime>2020-01-27T19:12:25.326Z</dateTime>
//                </status>
//             </TRDM>
//          </output>
//          <attachment>
//             <xop:Include href="cid:fefe5d81-468c-4639-a543-e758a3cbceea-2@ReturnTablePackage" xmlns:xop="http://www.w3.org/2004/08/xop/include"/>
//          </attachment>
//       </getTableResponseElement>
//    </soap:Body>
// </soap:Envelope>

const successResponseString = "Successful"
const lineOfAccounting = "LN_OF_ACCT"
const transportationAccountingCode = "TRNSPRTN_ACNT"

type GetTableRequestElement struct {
	soapClient    SoapCaller
	securityToken string
	privateKey    *rsa.PrivateKey
	Input         struct {
		TRDM struct {
			PhysicalName  string `xml:"physicalName"`
			ReturnContent string `xml:"returnContent"`
		}
	}
}

type GetTableResponseElement struct {
	XMLName xml.Name `xml:"getTableResponseElement"`
	Output  struct {
		TRDM struct {
			Status struct {
				RowCount   string `xml:"rowCount"`
				StatusCode string `xml:"statusCode"`
				DateTime   string `xml:"dateTime"`
			} `xml:"status"`
		} `xml:"TRDM"`
	} `xml:"output"`
	Attachment struct {
		Include struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Xop  string `xml:"xop,attr"`
		}
	}
}

type GetTableUpdater interface {
	GetTable(appCtx appcontext.AppContext, physicalName string, lastUpdate string) error
}

func NewGetTable(physicalName string, securityToken string, privateKey *rsa.PrivateKey, soapClient SoapCaller) GetTableUpdater {
	return &GetTableRequestElement{
		securityToken: securityToken,
		privateKey:    privateKey,
		soapClient:    soapClient,
		Input: struct {
			TRDM struct {
				PhysicalName  string `xml:"physicalName"`
				ReturnContent string `xml:"returnContent"`
			}
		}{
			TRDM: struct {
				PhysicalName  string `xml:"physicalName"`
				ReturnContent string `xml:"returnContent"`
			}{
				PhysicalName:  physicalName,
				ReturnContent: fmt.Sprintf("%t", true),
			},
		},
	}
}
func FetchTACRecordsByTime(appcontext appcontext.AppContext, time string) ([]models.TransportationAccountingCode, error) {
	var tacCodes []models.TransportationAccountingCode
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&tacCodes)

	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil
}

func FetchLOARecordsByTime(appcontext appcontext.AppContext, time string) ([]models.LineOfAccounting, error) {
	var loa []models.LineOfAccounting
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&loa)

	if err != nil {
		return loa, errors.Wrap(err, "Fetch line items query failed")
	}

	return loa, nil
}
func (d *GetTableRequestElement) GetTable(appCtx appcontext.AppContext, physicalName string, lastUpdate string) error {
	switch physicalName {
	case lineOfAccounting:
		loaRecords, loaFetchErr := FetchLOARecordsByTime(appCtx, lastUpdate)

		if loaFetchErr != nil {
			return loaFetchErr
		}

		if len(loaRecords) > 0 {
			if err := setupSoapCall(d, appCtx, physicalName); err != nil {
				return err
			}
		}
	case transportationAccountingCode:
		tacRecords, fetchErr := FetchTACRecordsByTime(appCtx, lastUpdate)
		if fetchErr != nil {
			return fetchErr
		}
		if len(tacRecords) > 0 {
			if err := setupSoapCall(d, appCtx, physicalName); err != nil {
				return err
			}
		}
	}

	return nil
}

func setupSoapCall(d *GetTableRequestElement, appCtx appcontext.AppContext, physicalName string) error {
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ret":     "http://ReturnTablePackage/",
	})
	params := GetTableRequestElement{
		Input: d.Input,
	}

	marshaledBody, marshalErr := xml.Marshal(params)
	if marshalErr != nil {
		return marshalErr
	}

	operation := func() error {
		signedHeader, headerSigningError := GenerateSignedHeader(d.securityToken, marshaledBody, d.privateKey)
		if headerSigningError != nil {
			return headerSigningError
		}
		newParams := gosoap.Params{
			"header": signedHeader,
			"body":   marshaledBody,
		}
		return getTableSoapCall(d, newParams, appCtx, physicalName)
	}
	b := backoff.NewExponentialBackOff()

	// Set the max retries to 5
	b.MaxElapsedTime = 5 * time.Hour

	// Only re-call after 1 hour
	b.InitialInterval = 1 * time.Hour
	err := backoff.Retry(operation, b)
	if err != nil {
		return fmt.Errorf("Failed after retries: %s", err)
	}
	return nil
}

func getTableSoapCall(d *GetTableRequestElement, params gosoap.Params, appCtx appcontext.AppContext, physicalName string) error {
	response, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return err
	}
	var r GetTableResponseElement
	unmarshalErr := response.Unmarshal(&r)
	if unmarshalErr != nil {
		return fmt.Errorf("unmarshall error: %s", unmarshalErr.Error())
	}
	if r.Output.TRDM.Status.StatusCode == successResponseString {
		parseError := parseGetTableResponse(appCtx, response, physicalName)
		if parseError != nil {
			return parseError
		}
	} else {
		err := errors.New("GetTable response was not `Successful`")
		return fmt.Errorf(r.Output.TRDM.Status.StatusCode, err)
	}
	appCtx.Logger().Debug("getTable result", zap.Any("processRequestResponse", response))
	return nil
}

func parseGetTableResponse(appcontext appcontext.AppContext, response *gosoap.Response, physicalName string) error {
	reader := bytes.NewReader(response.Payload)
	switch physicalName {
	case lineOfAccounting:
		loaCodes, err := loa.Parse(reader)
		if err != nil {
			return err
		}
		saveErr := saveLoaCodes(appcontext, loaCodes)
		if saveErr != nil {
			return saveErr
		}
	case transportationAccountingCode:
		tacCodes, err := tac.Parse(reader)
		consolidatedTacs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(tacCodes)
		if err != nil {
			return err
		}
		if saveErr := saveTacCodes(appcontext, consolidatedTacs); saveErr != nil {
			return saveErr
		}
	default:
		return nil
	}
	return nil
}

func saveTacCodes(appcontext appcontext.AppContext, tacCodes []models.TransportationAccountingCode) error {
	saveErr := appcontext.DB().Update(tacCodes)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func saveLoaCodes(appcontext appcontext.AppContext, loa []models.LineOfAccounting) error {
	saveErr := appcontext.DB().Update(loa)
	if saveErr != nil {
		return saveErr
	}
	return nil
}
