package trdm

import (
	"fmt"

	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
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

type GetTableRequestElement struct {
	soapClient SoapCaller
	Input      struct {
		TRDM struct {
			PhysicalName  string `xml:"physicalName"`
			ReturnContent string `xml:"returnContent"`
		}
	}
}

type GetTableResponseElement struct {
	Output struct {
		TRDM struct {
			Status struct {
				RowCount   string `xml:"rowCount"`
				StatusCode string `xml:"statusCode"`
				DateTime   string `xml:"dateTime"`
			}
		}
	}
	Attachment struct {
		Include struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Xop  string `xml:"xop,attr"`
		}
	}
}

type GetTableUpdater interface {
	GetTable(appCtx appcontext.AppContext, physicalName string, returnContent bool) error
}

func NewGetTable(physicalName string, returnContent bool, soapClient SoapCaller) GetTableUpdater {
	return &GetTableRequestElement{
		soapClient: soapClient,
		Input: struct {
			TRDM struct {
				PhysicalName  string "xml:\"physicalName\""
				ReturnContent string "xml:\"returnContent\""
			}
		}{
			TRDM: struct {
				PhysicalName  string "xml:\"physicalName\""
				ReturnContent string "xml:\"returnContent\""
			}{
				PhysicalName:  physicalName,
				ReturnContent: fmt.Sprintf("%t", returnContent),
			},
		},
	}
}

func (d *GetTableRequestElement) GetTable(appCtx appcontext.AppContext, physicalName string, returnContent bool) error {

	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/", //! Replace
	})

	params := gosoap.Params{
		"getTableRequestElement": map[string]interface{}{
			"Input": map[string]interface{}{
				"TRDM": map[string]interface{}{
					"physicalName":  physicalName,
					"returnContent": returnContent,
				},
			},
		},
	}
	response, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return fmt.Errorf("call error: %s", err.Error())
	}

	var r GetTableResponseElement
	unmarshalErr := response.Unmarshal(&r)
	if unmarshalErr != nil {
		return fmt.Errorf("unmarshall error: %s", unmarshalErr.Error())
	}

	if r.Output.TRDM.Status.StatusCode == successResponseString {
		println("Hi")
	}

	appCtx.Logger().Debug("getTable result", zap.Any("processRequestResponse", response))

	return nil
}
