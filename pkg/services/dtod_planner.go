package services

// DTODPlannerMileage is the exported interface for connecting to DTOD SOAP service and requesting distance mileage
//go:generate mockery -name DTODPlannerMileage
type DTODPlannerMileage interface {
	DTODZip5Distance(pickup string, destination string) (int, error)
}
