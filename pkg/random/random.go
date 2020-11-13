package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	mrand "math/rand"
)

// NewCryptoSeededSource utilizes the cryptographically secure generator to create a seed source
func NewCryptoSeededSource() mrand.Source {
	var seed int64
	err := binary.Read(crand.Reader, binary.BigEndian, &seed)
	if err != nil {
		log.Panicf("failed to create crypto seeded source: %v", err)
	}
	return mrand.NewSource(seed)
}
