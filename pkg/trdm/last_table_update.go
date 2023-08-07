package trdm

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

/*******************************************

The struct/method getLastTableUpdate:GetLastTableUpdate implements the service TRDM.

This method GetLastTableUpdate sends a SOAP request to TRDM to get the last table update.
This code is using the gosoap lib https://github.com/tiaguinho/gosoap

The Request to GetTable
SOAP Request:
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ret="http://ReturnTablePackage/">
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
type getLastTableUpdateReq struct {
	physicalName string
	soapClient   SoapCaller
}

type GetLastTableUpdateResponseElement struct {
	XMLName    xml.Name `xml:"getLastTableUpdateResponseElement"`
	LastUpdate string   `xml:"lastUpdate"`
	Status     struct {
		StatusCode string `xml:"statusCode"`
		DateTime   string `xml:"dateTime"`
	} `xml:"status"`
}

type TACCodes struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Tac       string    `json:"tac" db:"tac"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewTRDMGetLastTableUpdate(physicalName string, soapClient SoapCaller) GetLastTableUpdater {
	return &getLastTableUpdateReq{
		physicalName: physicalName,
		soapClient:   soapClient,
	}

}

// FetchAllTACRecords queries and fetches all transportation_accounting_codes
func FetchAllTACRecords(appcontext appcontext.AppContext) ([]TACCodes, error) {
	var tacCodes []TACCodes
	query := `SELECT * FROM transportation_accounting_codes`

	err := appcontext.DB().RawQuery(query).All(&tacCodes)
	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil

}

func (d *getLastTableUpdateReq) GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string) error {

	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/", //! Replace
	})

	params := gosoap.Params{
		"getLastTableUpdateRequestElement": map[string]interface{}{
			"physicalName": physicalName,
		},
	}
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
		tacCodes, dbError := FetchAllTACRecords(appCtx)
		if dbError != nil {
			return fmt.Errorf(err.Error())
		}
		for _, tacCode := range tacCodes {
			if tacCode.UpdatedAt.String() != r.LastUpdate {
				getTable := NewGetTable(physicalName, d.soapClient)
				getTableErr := getTable.GetTable(appCtx, physicalName)
				if getTableErr != nil {
					return fmt.Errorf("getTable error: %s", getTableErr.Error())
				}
			}
		}
	}

	appCtx.Logger().Debug("getLastTableUpdate result", zap.Any("processRequestResponse", r))

	return nil
}
