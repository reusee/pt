package pt

import (
	"fmt"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	iter := testTreap.NewIter(nil)
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
		iter := testTreap.NewIter(nil)
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
			iter := testTreap.NewIter(nil)
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
	ps := NewPrioritySource()
	var node *Treap[Int]
	const num = 128
	// insert
	var expected []Int
	for i := 0; i < num; i += 2 {
		node, _ = node.Upsert(Int(i), ps(), false)
		expected = append(expected, Int(i))
	}
	// seek
	for i := 0; i < num; i++ {
		iter := node.NewIter(nil)

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
