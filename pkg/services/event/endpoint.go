package event

// EndpointType stores the details of the endpoint API name and operationID
type EndpointType struct {
	APIName     string
	OperationID string
}

// EndpointKeyType is used to key into the map of EndpointTypes
type EndpointKeyType string

// EndpointMapType is used to map EndpointKeyType to info about the endpoint
type EndpointMapType map[EndpointKeyType]EndpointType

// PrimeAPIName is a const string to use the EndpointTypes
const PrimeAPIName string = "primeapi"

var apiEndpoints = []EndpointMapType{
	supportEndpoints,
	ghcEndpoints,
	internalEndpoints,
}

// String returns the string representation of the endpoint name
func (e EndpointType) String() string {
	return e.APIName + "." + e.OperationID
}

// GetEndpointAPI returns the api name of the endpoint
func GetEndpointAPI(key EndpointKeyType) *string {
	for _, endpointMap := range apiEndpoints {

		if endpointInfo, ok := endpointMap[key]; ok {
			var result = endpointInfo.APIName
			return &result
		}
	}
	return nil
}

// GetEndpointOperationID retuns the operation ID of the endpoint
func GetEndpointOperationID(key EndpointKeyType) *string {
	for _, endpointMap := range apiEndpoints {

		if endpointInfo, ok := endpointMap[key]; ok {
			var result = endpointInfo.OperationID
			return &result
		}
	}
	return nil
}
