package services

// DTODPlannerMileage is the exported interface for connecting to DTOD Planner and get
//go:generate mockery -name DTODPlannerMileage
type DTODPlannerMileage interface {
	DTODZip5Distance(pickup string, destination string) (int, error)
}
