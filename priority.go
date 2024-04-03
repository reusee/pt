package pt

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand/v2"
)

type Priority = int64

const (
	MinPriority = math.MinInt64
	MaxPriority = math.MaxInt64
)

type PrioritySource = func() Priority

func NewPrioritySource() PrioritySource {
	var s1, s2 uint64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &s1); err != nil {
		panic(err)
	}
	if err := binary.Read(crand.Reader, binary.LittleEndian, &s2); err != nil {
		panic(err)
	}
	r := rand.New(
		rand.NewPCG(s1, s2),
	)
	return func() Priority {
		return r.Int64() - 1 // minus one to avoid returning the max int64
	}
}
