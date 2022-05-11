package route

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
	"fmt"

	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

// DTODPlannerMileage is the interface for connecting to DTOD SOAP service and requesting distance mileage
type DTODPlannerMileage interface {
	DTODZip5Distance(appCtx appcontext.AppContext, pickup string, destination string) (int, error)
}

type dtodZip5DistanceInfo struct {
	username   string
	password   string
	soapClient SoapCaller
}

// Response XML structs
// There's more in the returned XML, but we're only trying to get to the distance.
type processRequestResponse struct {
	ProcessRequestResult processRequestResult `xml:"ProcessRequestResult"`
}

type processRequestResult struct {
	Distance float64 `xml:"Distance"`
}

// NewDTODZip5Distance returns a new DTOD Planner Mileage interface
func NewDTODZip5Distance(username string, password string, soapClient SoapCaller) DTODPlannerMileage {
	return &dtodZip5DistanceInfo{
		username:   username,
		password:   password,
		soapClient: soapClient,
	}
}

// DTODZip5Distance returns the distance in miles between the pickup and destination zips
func (d *dtodZip5DistanceInfo) DTODZip5Distance(appCtx appcontext.AppContext, pickupZip string, destinationZip string) (int, error) {
	distance := 0

	// set custom envelope
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/",
	})

	params := gosoap.Params{
		"DtodRequest": map[string]interface{}{
			"AuthToken": map[string]interface{}{
				"Username": d.username,
				"Password": d.password,
			},
			"UserRequest": map[string]interface{}{
				"Function": "Distance",
				// TODO: Default RouteType is PcsTdyTravel, but CommercialPersonalProperty seems better. Verify.
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

	res, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return distance, fmt.Errorf("call error: %s", err.Error())
	}

	var r processRequestResponse
	err = res.Unmarshal(&r)
	if err != nil {
		return distance, fmt.Errorf("unmarshal error: %s", err.Error())
	}

	// It looks like sending a bad zip just returns a distance of -1, so test for that
	distanceFloat := r.ProcessRequestResult.Distance
	if distanceFloat <= 0 {
		return distance, fmt.Errorf("invalid distance using pickup %s and destination %s", pickupZip, destinationZip)
	}

	// TODO: DTOD gives us a float back. Should we round, floor, or ceiling? Just going to round for now.
	distance = int(distanceFloat + 0.5)

	appCtx.Logger().Debug("dtod result", zap.Any("processRequestResponse", r), zap.Int("distance", distance))

	return distance, nil
}
