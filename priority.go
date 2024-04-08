package pt

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
)

type Priority = int64

const (
	MinPriority = math.MinInt64
	MaxPriority = math.MaxInt64
)

type PrioritySource = func() Priority

func NewPrioritySource() PrioritySource {
	var s1 int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &s1); err != nil {
		panic(err)
	}
	r := rand.New(
		rand.NewSource(s1),
	)
	return func() Priority {
		return r.Int63() - 1 // minus one to avoid returning the max int64
	}
}
