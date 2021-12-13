package baselinetest

import (
	"encoding/binary"

	"github.com/gofrs/uuid"
)

type FakeGenerator struct {
	counter uint64
}

func (g *FakeGenerator) NewFromCounter() uuid.UUID {
	g.counter++

	u := uuid.UUID{}
	binary.BigEndian.PutUint16(u[0:], 0)
	binary.BigEndian.PutUint64(u[2:], g.counter)

	u.SetVariant(uuid.VariantRFC4122)
	return u
}

func (g *FakeGenerator) NewV1() (uuid.UUID, error) {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V1)

	return u, nil
}

func (g *FakeGenerator) NewV3(ns uuid.UUID, name string) uuid.UUID {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V3)
	return u
}

func (g *FakeGenerator) NewV4() (uuid.UUID, error) {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V4)

	return u, nil
}

func (g *FakeGenerator) NewV5(ns uuid.UUID, name string) uuid.UUID {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V5)

	return u
}

func (g *FakeGenerator) NewV6() (uuid.UUID, error) {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V6)

	return u, nil
}

func (g *FakeGenerator) NewV7(uuid.Precision) (uuid.UUID, error) {
	u := g.NewFromCounter()
	u.SetVersion(uuid.V7)

	return u, nil
}

func NewFakeGenerator() uuid.Generator {
	return &FakeGenerator{}
}
