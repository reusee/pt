package pt

import (
	"fmt"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	iter := testTreap.NewIter()
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

var testTreap = &Treap[Int]{
	left: &Treap[Int]{
		left: &Treap[Int]{
			value: 1,
		},
		right: &Treap[Int]{
			value: 2,
		},
		value: 3,
	},
	right: &Treap[Int]{
		right: &Treap[Int]{
			value: 4,
		},
		value: 5,
	},
	value: 6,
}

func BenchmarkIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := testTreap.NewIter()
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
			iter := testTreap.NewIter()
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
