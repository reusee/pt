package pt

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestNode(t *testing.T) {
	const num = 65536
	var n *node[Int]

	// upsert
	for _, i := range rand.Perm(num) {
		priority := NewPriority()
		n = n.upsert(Int(i), priority)
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

	// remove
	for _, i := range rand.Perm(num) {
		n = n.remove(Int(i))
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
		n = n.upsert(i, NewPriority())
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
		n = n.upsert(Int(i), NewPriority())
	}
}

func BenchmarkDelete(b *testing.B) {
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n = n.upsert(Int(i), NewPriority())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n = n.remove(Int(i))
	}
}
