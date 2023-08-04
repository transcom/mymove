package trdm

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

type GetTableRequestElement struct {
	Input struct {
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
				Text       string `xml:",chardata"`
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
