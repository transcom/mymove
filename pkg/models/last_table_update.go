package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
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
// Date/time value is used in conjunction with the contentUpdatedSinceDateTime column in the getTable method.
type TRDMGetLastTableUpdater interface {
	GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string, returnContent bool) error
}
type getLastTableUpdateReq struct {
	physicalName  string
	returnContent bool
	soapClient    SoapCaller
}

type SoapCaller interface {
	Call(m string, p gosoap.Params) (res *gosoap.Response, err error)
}

type getTableResponse struct {
	GetTableResponseElement getTableResponseElement `xml:"getTableResponseElement"`
}

// Response XML Struct
type getTableResponseElement struct {
	LastUpdate time.Time `xml:"lastUpdate"`
	RowCount   int       `xml:"rowCount"`
	StatusCode string    `xml:"statusCode"`
	Message    string    `xml:"message"`
	DateTime   time.Time `xml:"dateTime"`
}
type TACCodes struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Tac       string    `json:"tac" db:"tac"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewTRDMGetLastTableUpdate(physicalName string, returnContent bool, soapClient SoapCaller) TRDMGetLastTableUpdater {
	return &getLastTableUpdateReq{
		physicalName:  physicalName,
		returnContent: returnContent,
		soapClient:    soapClient,
	}

}

// FetchAllTACRecords queries and fetches all transportation_accounting_codes
func fetchAllTACRecords(dbConnection *pop.Connection) ([]TACCodes, error) {
	var tacCodes []TACCodes
	query := `Select * from transportation_accounting_codes`

	err := dbConnection.RawQuery(query).All(&tacCodes)
	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil

}

func (d *getLastTableUpdateReq) GetLastTableUpdate(appCtx appcontext.AppContext, physicalName string, returnContent bool) error {

	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/", //! Replace
	})

	params := gosoap.Params{
		"TRDM": map[string]interface{}{
			"physicalName":  physicalName,
			"returnContent": returnContent,
		},
	}
	res, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return fmt.Errorf("call error: %s", err.Error())
	}

	var r getTableResponse
	err = res.Unmarshal(&r)

	if err != nil {
		return fmt.Errorf("unmarshal error: %s", err.Error())
	}

	if r.GetTableResponseElement.RowCount != 0 {
		tacCodes, dbError := fetchAllTACRecords(appCtx.DB())
		if dbError != nil {
			return fmt.Errorf(err.Error())
		}
		for _, tacCode := range tacCodes {
			if tacCode.UpdatedAt != r.GetTableResponseElement.LastUpdate {
				return nil
			}
		}
	}

	appCtx.Logger().Debug("getLastTableUpdate result", zap.Any("processRequestResponse", r))

	return nil
}
