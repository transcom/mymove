package dtod

/*******************************************

The struct/method dtodZip5DistanceInfo:DTODZip5Distance implements the service DTODPlannerMileage.

This method DTODZip5Distance sends a SOAP request to DTOD to get the mileage between two ZIP 5 locations.
The DTOD web service is described by this document https://docs.google.com/document/d/1yUsk8JWj1u-EBfdLiCBHOtrUWaYgBXFv/edit
This code is using the gosoap lib https://github.com/tiaguinho/gosoap

The Request to DTOD using the service ProcessRequest which looks like

<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="https://dtod.sddc.army.mil/service/">
    <soapenv:Body>
        <ProcessRequest xmlns="https://dtod.sddc.army.mil/service/">
            <DtodRequest>
                <UserRequest>
                    <Function>Distance</Function>
                    <Origin>
                        <ZipCode>05030</ZipCode>
                    </Origin>
                    <Destination>
                        <ZipCode>05091</ZipCode>
                    </Destination>
                </UserRequest>
                <AuthToken>
                    <Username>theusername</Username>
                    <Password>thepassword</Password>
                </AuthToken>
            </DtodRequest>
        </ProcessRequest>
    </soapenv:Body>
</soapenv:Envelope>

 *******************************************/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/services"
)

const (
	dtodRequestTimeout = time.Duration(30) * time.Second
)

var (
	// UsernameFlag DTOD username env flag
	UsernameFlag = map[string]string{
		"cli": "dtod-api-username",
		"env": "DTOD_API_USERNAME",
	}
	// PasswordFlag DTOD password env flag
	PasswordFlag = map[string]string{
		"cli": "dtod-api-password",
		"env": "DTOD_API_PASSWORD",
	}
	// URLFlag DTOD URL env flag
	URLFlag = map[string]string{
		"cli": "dtod-api-url",
		"env": "DTOD_API_URL",
	}
	// WSDLFlag DTOD WSDL env flag
	WSDLFlag = map[string]string{
		"cli": "dtod-api-wsdl",
		"env": "DTOD_API_WSDL",
	}
)

type dtodZip5DistanceInfo struct {
	logger    Logger
	tlsConfig *tls.Config
	username  string
	password  string
	url       string
	wsdl      string
}

// Response XML structs
// There's more in the returned XML, but we're only trying to get to the distance.
type processRequestResponse struct {
	ProcessRequestResult processRequestResult `xml:"ProcessRequestResult"`
}

type processRequestResult struct {
	Distance float64 `xml:"Distance"`
}

// InitDTODFlags initializes DTOD command line flags
func InitDTODFlags(flag *pflag.FlagSet) {
	flag.String(UsernameFlag["cli"], "", "DTOD api auth username")
	flag.String(PasswordFlag["cli"], "", "DTOD api auth password")
	flag.String(URLFlag["cli"], "", "URL for sending an SOAP request to DTOD")
	flag.String(WSDLFlag["cli"], "", "WSDL for sending an SOAP request to DTOD")
}

// GetDTODFlags return the DTOD flag values
func GetDTODFlags(v *viper.Viper) (string, string, string, string, error) {
	username := v.GetString(UsernameFlag["cli"])
	if len(username) == 0 {
		username = os.Getenv(UsernameFlag["env"])
		if len(username) == 0 {
			return "", "", "", "", fmt.Errorf("%s not set", UsernameFlag["env"])
		}
	}
	password := v.GetString(PasswordFlag["cli"])
	if len(password) == 0 {
		password = os.Getenv(PasswordFlag["env"])
		if len(password) == 0 {
			return "", "", "", "", fmt.Errorf("%s not set", PasswordFlag["env"])
		}
	}
	url := v.GetString(URLFlag["cli"])
	if len(url) == 0 {
		url = os.Getenv(URLFlag["env"])
		if len(url) == 0 {
			return "", "", "", "", fmt.Errorf("%s not set", URLFlag["env"])
		}
	}
	wsdl := v.GetString(WSDLFlag["cli"])
	if len(url) == 0 {
		url = os.Getenv(WSDLFlag["env"])
		if len(url) == 0 {
			return "", "", "", "", fmt.Errorf("%s not set", WSDLFlag["env"])
		}
	}
	return username, password, url, wsdl, nil
}

// NewDTODZip5Distance returns a new DTOD Planner Mileage interface
func NewDTODZip5Distance(logger Logger, tlsConfig *tls.Config, username string, password string, url string, wsdl string) services.DTODPlannerMileage {
	return &dtodZip5DistanceInfo{
		logger:    logger,
		tlsConfig: tlsConfig,
		username:  username,
		password:  password,
		url:       url,
		wsdl:      wsdl,
	}
}

func (d *dtodZip5DistanceInfo) DTODZip5Distance(pickupZip string, destinationZip string) (int, error) {
	distance := 0

	tr := &http.Transport{TLSClientConfig: d.tlsConfig}
	httpClient := &http.Client{Transport: tr, Timeout: dtodRequestTimeout}

	// set custom envelope
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/",
	})

	// TODO: for dev uncomment and use SoapClientWithConfig to see both request and response
	//config := gosoap.Config{
	//	Dump: true,
	//}
	//soap, err := gosoap.SoapClientWithConfig(d.wsdl, httpClient, &config)

	soap, err := gosoap.SoapClient(d.wsdl, httpClient)
	if err != nil {
		return distance, fmt.Errorf("SoapClient error: %s", err.Error())
	}

	params := gosoap.Params{
		"DtodRequest": map[string]interface{}{
			"AuthToken": map[string]interface{}{
				"Username": d.username,
				"Password": d.password,
			},
			"UserRequest": map[string]interface{}{
				"Function": "Distance",
				// TODO: We are currently defaulting RouteType to PcsTdyTravel. Should we be using
				//   DityPersonalProperty or CommercialPersonalProperty instead?
				"RouteType": "CommercialPersonalProperty",
				"Origin": map[string]interface{}{
					"ZipCode": pickupZip,
				},
				"Destination": map[string]interface{}{
					"ZipCode": destinationZip,
				},
			},
		},
	}

	soap.URL = d.url
	res, err := soap.Call("ProcessRequest", params)
	if err != nil {
		return distance, fmt.Errorf("call error: %s", err.Error())
	}

	var r processRequestResponse
	err = res.Unmarshal(&r)
	if err != nil {
		fmt.Printf("Unmarshal error: %s", err.Error())
		d.logger.Fatal("Unmarshal error: ", zap.Error(err))
	}

	// It looks like sending a bad zip just returns a distance of -1.
	distanceFloat := r.ProcessRequestResult.Distance
	if distanceFloat <= 0 {
		return distance, fmt.Errorf("invalid distance using pickup %s and destination %s", pickupZip, destinationZip)
	}

	// TODO: DTOD gives us a float back. Should we round, floor, or ceiling? Just going to round for now.
	distance = int(distanceFloat + 0.5)

	d.logger.Debug("dtod result", zap.Any("processRequestResponse", r), zap.Int("distance", distance))

	return distance, nil
}
