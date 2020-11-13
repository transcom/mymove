package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/big"
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

// GetRandomInt takes an int and returns a cryptographically sourced random integer of type int64 (exclusive of the max int)
func GetRandomInt(max int) (int, error) {
	randMax := big.NewInt(int64(max))
	randInt, err := crand.Int(crand.Reader, randMax)
	if err != nil {
		return 0, err
	}
	return int(randInt.Int64()), nil
}

// GetRandomIntAddend takes a min and max integer and returns a cryptographically secure integer reflecting the difference between those ranges
func GetRandomIntAddend(min int, max int) (int, error) {
	diff := max - min
	addend, err := GetRandomInt(diff)
	if err != nil {
		return 0, err
	}
	return addend, nil
}
