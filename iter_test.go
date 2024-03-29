package pt

import (
	"fmt"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	iter := testNode.iter(nil)

	out := new(strings.Builder)
	var n Int
	for {
		n, iter = iter()
		fmt.Fprintf(out, "%v", n)
		if iter == nil {
			break
		}
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
	iter := testNode.iter(nil)
	for i := 0; i < b.N; i++ {
		_, iter = iter()
		if iter == nil {
			iter = testNode.iter(nil)
		}
	}
}
