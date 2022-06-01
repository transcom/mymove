package v6

type Connection struct {
	ID          string
	Elapsed     int64
	eager       bool
	eagerFields []string
}
