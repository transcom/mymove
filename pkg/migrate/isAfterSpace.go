package migrate

// isAfterSpace returns true if the index directly follows a space character
//
// The function also returns true if the index is at the beginning of the buffer
func isAfterSpace(in *Buffer, i int) bool {
	if i == 0 {
		return true
	}
	prev, err := in.Index(i - 1)
	return err == nil && byteIsSpace(prev)
}
