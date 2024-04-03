package pt

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkUpsert(b *testing.B) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), ps(), false)
	}
}

func BenchmarkDelete(b *testing.B) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.remove(Int(i), false)
	}
}

func BenchmarkSplit(b *testing.B) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.split(Int(rand.N(b.N)), false)
	}
}

func BenchmarkUpsertPriority(b *testing.B) {
	var n *node[Int]
	n, _ = n.upsert(1, -1, false) // will be the left node
	n, _ = n.upsert(3, -1, false) // will be the right node
	for i := 0; i < b.N; i++ {
		// upsert node priority with non-empty left and right nodes
		n, _ = n.upsert(2, int64(i), false)
	}
}

func BenchmarkUnion(b *testing.B) {
	ps := NewPrioritySource()
	const l = 1024
	var n1, n2 *node[Int]
	for i := 0; i < l; i++ {
		n1, _ = n1.upsert(Int(i), ps(), false)
		n2, _ = n2.upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.union(n2, false)
	}
}

func BenchmarkGet(b *testing.B) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := n.get(Int(i))
		if !ok {
			b.Fatal()
		}
	}
}

func BenchmarkUpsert65536(b *testing.B) {
	ps := NewPrioritySource()
	for i := 0; i < b.N; i++ {
		var n *node[Int]
		for k := range 65536 {
			n, _ = n.upsert(Int(k), ps(), false)
		}
	}
}

func BenchmarkBuild65536(b *testing.B) {
	ps := NewPrioritySource()
	var slice []Int
	for i := 0; i < 65536; i++ {
		slice = append(slice, Int(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		build(ps, slice)
	}
}
