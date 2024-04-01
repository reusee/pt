package pt

import (
	"math"
	"math/rand/v2"
)

type Priority = int64

const minPriority = math.MinInt64

func NewPriority() Priority {
	return rand.Int64()
}
