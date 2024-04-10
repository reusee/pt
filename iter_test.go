package pt

import (
	"fmt"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	iter := testNode.NewIter()
	defer iter.Close()

	for i := 0; i < 8; i++ {
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

		iter.Rewind()
	}
}

var testNode = &_Node[Int]{
	left: &_Node[Int]{
		left: &_Node[Int]{
			value: 1,
		},
		right: &_Node[Int]{
			value: 2,
		},
		value: 3,
	},
	right: &_Node[Int]{
		right: &_Node[Int]{
			value: 4,
		},
		value: 5,
	},
	value: 6,
}

func BenchmarkIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := testNode.NewIter()
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
			iter := testNode.NewIter()
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

func TestSeek(t *testing.T) {
	ps := newPrioritySource()
	var node *_Node[Int]
	const num = 128
	// insert
	var expected []Int
	for i := 0; i < num; i += 2 {
		node, _ = node.Upsert(Int(i), ps(), false)
		expected = append(expected, Int(i))
	}
	// seek
	for i := 0; i < num; i++ {
		iter := node.NewIter()

		// seek
		switch i % 2 {

		case 0:
			n, ok := iter.Seek(Int(i))
			if !ok {
				t.Fatal()
			}
			if n != Int(i) {
				t.Fatal()
			}

		case 1:
			n, ok := iter.Seek(Int(i))
			if i == num-1 {
				// last
				if ok {
					t.Fatal()
				}
			} else {
				if !ok {
					t.Fatal()
				}
				if n != Int(i+1) {
					t.Fatal()
				}
			}

		}

		// next
		if i != num-1 {
			idx := (i+1)/2 + 1
			for {
				n, ok := iter.Next()
				if !ok {
					break
				}
				if n != expected[idx] {
					t.Fatal()
				}
				idx++
			}
			if idx != num/2 {
				t.Fatal()
			}
		}
	}
}
