package migrate

func isAfterSpace(in *Buffer, i int) bool {
	if i == 0 {
		return true
	}
	prev, err := in.Index(i - 1)
	return err == nil && byteIsSpace(prev)
}
