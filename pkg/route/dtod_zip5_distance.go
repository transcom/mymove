package route

import (
	"fmt"

	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/notifications"
)

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

// DTODPlannerMileage is the interface for connecting to DTOD SOAP service and requesting distance mileage
// NOTE: Placing this in a separate package/directory to avoid a circular dependency from an existing mock.
//
//go:generate mockery --name DTODPlannerMileage --outpkg ghcmocks --output ./ghcmocks
type DTODPlannerMileage interface {
	DTODZip5Distance(appCtx appcontext.AppContext, pickup string, destination string) (int, error)
}

type dtodZip5DistanceInfo struct {
	username       string
	password       string
	soapClient     SoapCaller
	simulateOutage bool
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
func NewDTODZip5Distance(username string, password string, soapClient SoapCaller, simulateOutage bool) DTODPlannerMileage {
	return &dtodZip5DistanceInfo{
		username:       username,
		password:       password,
		soapClient:     soapClient,
		simulateOutage: simulateOutage,
	}
}

// DTODZip5Distance returns the distance in miles between the pickup and destination zips
func (d *dtodZip5DistanceInfo) DTODZip5Distance(appCtx appcontext.AppContext, pickupZip string, destinationZip string) (int, error) {
	distance := 0

	params := createDTODParams(d.username, d.password, pickupZip, destinationZip)

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

	if d.simulateOutage {
		distanceFloat = -1
	}

	if distanceFloat <= 0 {
		dtodAvailable, _ := validateDTODServiceAvailable(*d)
		if !dtodAvailable && appCtx.Session().IsServiceMember() {
			return distance, nil
		}

		return distance, apperror.NewEventError(notifications.DtodErrorMessage, nil)
	}

	// TODO: DTOD gives us a float back. Should we round, floor, or ceiling? Just going to round for now.
	distance = int(distanceFloat + 0.5)

	appCtx.Logger().Debug("dtod result", zap.Any("processRequestResponse", r), zap.Int("distance", distance))

	return distance, nil
}

// validateDTODServiceAvailable pings the DTOD service with zips that are known to be accepted.
// This is used to verify that the DTOD service is live.
func validateDTODServiceAvailable(d dtodZip5DistanceInfo) (bool, error) {

	if d.simulateOutage {
		return false, nil
	}

	params := createDTODParams(d.username, d.password, "20001", "20301")

	res, err := d.soapClient.Call("ProcessRequest", params)
	if err != nil {
		return false, fmt.Errorf("call error: %s", err.Error())
	}

	var r processRequestResponse
	err = res.Unmarshal(&r)
	if err != nil {
		return false, fmt.Errorf("unmarshal error: %s", err.Error())
	}

	distanceFloat := r.ProcessRequestResult.Distance

	if distanceFloat > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func createDTODParams(username string, password string, pickupZip string, destinationZip string) gosoap.Params {
	// set custom envelope
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:ser":     "https://dtod.sddc.army.mil/service/",
	})

	params := gosoap.Params{
		"DtodRequest": map[string]interface{}{
			"AuthToken": map[string]interface{}{
				"Username": username,
				"Password": password,
			},
			"UserRequest": map[string]interface{}{
				"Function":  "Distance",
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

	return params
}
