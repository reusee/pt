package pt

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
)

type _Priority = int64

const (
	minPriority = math.MinInt64
	maxPriority = math.MaxInt64
)

type _PrioritySource = func() _Priority

func newPrioritySource() _PrioritySource {
	var s1 int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &s1); err != nil {
		panic(err)
	}
	r := rand.New(
		rand.NewSource(s1),
	)
	return func() _Priority {
		return r.Int63() - 1 // minus one to avoid returning the max int64
	}
}
