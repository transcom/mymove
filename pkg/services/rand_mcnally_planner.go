package services

// RandMcNallyPlannerMileage is the exported interface for looking up distances in the Rand McNally database table
//go:generate mockery -name RandMcNallyPlannerMileage
type RandMcNallyPlannerMileage interface {
	RandMcNallyZip3Distance(pickup string, destination string) (int, error)
}
