package pt

import (
	"fmt"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	iter := testNode.newIter()
	defer iter.Close()

	out := new(strings.Builder)
	for {
		n, ok := iter.Next()
		if !ok {
			break
		}
		fmt.Fprintf(out, "%v", n)
	}

	if out.String() != "132654" {
		t.Fatalf("got %v", out.String())
	}
}

var testNode = &node[Int]{
	left: &node[Int]{
		left: &node[Int]{
			value: 1,
		},
		right: &node[Int]{
			value: 2,
		},
		value: 3,
	},
	right: &node[Int]{
		right: &node[Int]{
			value: 4,
		},
		value: 5,
	},
	value: 6,
}

func BenchmarkIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := testNode.newIter()
		for {
			_, ok := iter.Next()
			if !ok {
				break
			}
		}
		iter.Close()
	}
}

func BenchmarkParallelIter(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			iter := testNode.newIter()
			for {
				_, ok := iter.Next()
				if !ok {
					break
				}
			}
			iter.Close()
		}
	})
}
