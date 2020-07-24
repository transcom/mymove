package events

// EndpointType stores the details of the endpoint API name and operationID
type EndpointType struct {
	APIName     string
	OperationID string
}

// EndpointKeyType is used to key into the map of EndpointTypes
type EndpointKeyType string

// -------------------- API NAMES --------------------

// primeAPIName is a const string to use the EndpointTypes
const primeAPIName string = "primeapi"

// -------------------- ENDPOINT KEYS --------------------

// PrimeFetchMTOUpdatesEndpointKey is the key for the fetchMTOUpdates endpoint in prime
const PrimeFetchMTOUpdatesEndpointKey = "Prime.FetchMTOUpdates"

// PrimeUpdateMTOPostCounselingInformationEndpointKey is the key for the updateMTOPostCounselingInformation endpoint in prime
const PrimeUpdateMTOPostCounselingInformationEndpointKey = "Prime.UpdateMTOPostCounselingInformation"

// PrimeCreateMTOShipmentEndpointKey is the key for the createMTOShipment endpoint in prime
const PrimeCreateMTOShipmentEndpointKey = "Prime.CreateMTOShipment"

// PrimeUpdateMTOShipmentEndpointKey is the key for the updateMTOShipment endpoint in prime
const PrimeUpdateMTOShipmentEndpointKey = "Prime.UpdateMTOShipment"

// PrimeCreateMTOServiceItemEndpointKey is the key for the createMTOServiceItem endpoint in prime
const PrimeCreateMTOServiceItemEndpointKey = "Prime.CreateMTOServiceItem"

// PrimeCreatePaymentRequestEndpointKey is the key for the createPaymentRequest endpoint in prime
const PrimeCreatePaymentRequestEndpointKey = "Prime.CreatePaymentRequest"

// PrimeCreateUploadEndpointKey is the key for the createUpload endpoint in prime
const PrimeCreateUploadEndpointKey = "Prime.CreateUpload"

// -------------------- ENDPOINT MAP --------------------
var endpoints map[EndpointKeyType]EndpointType = map[EndpointKeyType]EndpointType{
	PrimeFetchMTOUpdatesEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "fetchMTOUpdates",
	},
	PrimeUpdateMTOPostCounselingInformationEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "updateMTOPostCounselingInformation",
	},
	PrimeCreateMTOShipmentEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "createMTOShipment",
	},
	PrimeUpdateMTOShipmentEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "updateMTOShipment",
	},
	PrimeCreateMTOServiceItemEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "createMTOServiceItem",
	},
	PrimeCreatePaymentRequestEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "createPaymentRequest",
	},
	PrimeCreateUploadEndpointKey: {
		APIName:     primeAPIName,
		OperationID: "createUpload",
	},
}

// String returns the string representation of the endpoint name
func (e EndpointType) String() string {
	return e.APIName + "." + e.OperationID
}

// GetEndpointAPI returns the api name of the endpoint
func GetEndpointAPI(key EndpointKeyType) string {
	return endpoints[key].APIName
}

// GetEndpointOperationID retuns the operation ID of the endpoint
func GetEndpointOperationID(key EndpointKeyType) string {
	return endpoints[key].OperationID
}
