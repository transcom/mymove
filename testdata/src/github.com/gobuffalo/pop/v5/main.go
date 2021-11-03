package v5

type Connection struct {
	ID          string
	Elapsed     int64
	eager       bool
	eagerFields []string
}
