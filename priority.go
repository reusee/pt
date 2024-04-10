package pt

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
	"runtime"
	_ "unsafe"
)

type _Priority = int64

const (
	minPriority = math.MinInt64
	maxPriority = math.MaxInt64
)

type _PrioritySource = func() _Priority

func newPrioritySource() _PrioritySource {
	var shards []*rand.Rand
	numShards := runtime.GOMAXPROCS(-1)
	for i := 0; i < numShards; i++ {
		var s1 int64
		if err := binary.Read(crand.Reader, binary.LittleEndian, &s1); err != nil {
			panic(err)
		}
		r := rand.New(
			rand.NewSource(s1),
		)
		shards = append(shards, r)
	}
	return func() _Priority {
		proc := runtime_procPin()
		shard := shards[proc]
		runtime_procUnpin()
		return shard.Int63() - 1 // minus one to avoid returning the max int64
	}
}

//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin() int
