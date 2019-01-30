package sequence

// Sequencer provides an interface for generating sequence numbers.
type Sequencer interface {
	NextVal() (int64, error)
	SetVal(val int64) error
}
