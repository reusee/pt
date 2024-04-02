package pt

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestNode(t *testing.T) {
	const num = 4096
	var n *node[Int]

	// upsert
	for _, i := range rand.Perm(num) {
		priority := NewPriority()
		existed := false
		n, existed = n.upsert(Int(i), priority)
		if existed {
			t.Fatal()
		}
	}

	// iter
	iter := n.newIter()
	defer iter.Close()
	for i := 0; i < num; i++ {
		j, ok := iter.Next()
		if !ok {
			break
		}
		if j != Int(i) {
			t.Fatal()
		}
	}

	// height
	if h := n.height(); h > int(math.Log2(float64(65536))*3) {
		t.Fatalf("got %v", n.height())
	} else {
		pt("num %v, height %v\n", num, h)
	}

	// split
	for i := 0; i < num; i++ {
		split, existed := n.split(Int(i))
		if !existed {
			t.Fatal()
		}
		if l := split.length(); l != num {
			t.Fatalf("got %v, expected %v", l, num)
		}
		if l := split.left.length(); l != i {
			t.Fatal()
		}
		if l := split.right.length(); l != num-i-1 {
			t.Fatalf("got %v, expected %v", l, num-i-1)
		}
	}

	// remove
	for _, i := range rand.Perm(num) {
		removed := false
		n, removed = n.remove(Int(i))
		if !removed {
			t.Fatal()
		}
	}
	if n.height() != 0 {
		t.Fatal()
	}
}

func TestUpsertPersistence(t *testing.T) {
	const num = 1024
	var nodes []*node[Int]
	var n *node[Int]
	for i := Int(0); i < num; i++ {
		n, _ = n.upsert(i, NewPriority())
		nodes = append(nodes, n)
	}
	for i, n := range nodes {
		iter := n.newIter()
		for expected := Int(0); expected < Int(i+1); expected++ {
			got, ok := iter.Next()
			if !ok {
				t.Fatal()
			}
			if got != expected {
				t.Fatal()
			}
		}
		_, ok := iter.Next()
		if ok {
			t.Fatal()
		}
		iter.Close()
	}
}

func BenchmarkUpsert(b *testing.B) {
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), NewPriority())
	}
}

func BenchmarkDelete(b *testing.B) {
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), NewPriority())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.remove(Int(i))
	}
}

func BenchmarkSplit(b *testing.B) {
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), NewPriority())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.split(Int(i))
	}
}

func BenchmarkUpsertPriority(b *testing.B) {
	var n *node[Int]
	n, _ = n.upsert(1, -1) // will be the left node
	n, _ = n.upsert(3, -1) // will be the right node
	for i := 0; i < b.N; i++ {
		// upsert node priority with non-empty left and right nodes
		n, _ = n.upsert(2, int64(i))
	}
}
