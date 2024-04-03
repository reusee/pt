package pt

import (
	"math"
	"math/rand/v2"
)

type Priority = int64

const (
	MinPriority = math.MinInt64
	MaxPriority = math.MaxInt64
)

func NewPriority() Priority {
	return rand.Int64() - 1 // minus one to avoid returning the max int64
}
