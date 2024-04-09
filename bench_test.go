package pt

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkUpsert(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.Upsert(Int(i), ps(), false)
	}
}

func BenchmarkDelete(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.Upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.Remove(Int(i), false)
	}
}

func BenchmarkSplit(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.Upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.Split(Int(rand.Intn(b.N)), false)
	}
}

func BenchmarkUpsertPriority(b *testing.B) {
	var n *Treap[Int]
	n, _ = n.Upsert(1, -1, false) // will be the left node
	n, _ = n.Upsert(3, -1, false) // will be the right node
	for i := 0; i < b.N; i++ {
		// upsert node priority with non-empty left and right nodes
		n, _ = n.Upsert(2, int64(i), false)
	}
}

func BenchmarkUnion(b *testing.B) {
	ps := NewPrioritySource()
	const l = 1024
	var n1, n2 *Treap[Int]
	for i := 0; i < l; i++ {
		n1, _ = n1.Upsert(Int(i), ps(), false)
		n2, _ = n2.Upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.Union(n2, false)
	}
}

func BenchmarkGet(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < b.N; i++ {
		n, _ = n.Upsert(Int(i), ps(), false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := n.Get(Int(i))
		if !ok {
			b.Fatal()
		}
	}
}

func BenchmarkUpsert65536(b *testing.B) {
	ps := NewPrioritySource()
	for i := 0; i < b.N; i++ {
		var n *Treap[Int]
		for k := 0; k < 65536; k++ {
			n, _ = n.Upsert(Int(k), ps(), false)
		}
	}
}

func BenchmarkParallelUpsert65536(b *testing.B) {
	// workers
	jobs := make(chan func(PrioritySource))
	quit := make(chan bool)
	defer func() {
		close(quit)
	}()
	for i := 0; i < b.N; i++ {
		go func() {
			ps := NewPrioritySource()
			for {
				select {
				case job := <-jobs:
					job(ps)
				case <-quit:
					return
				}
			}
		}()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		var n atomic.Pointer[Treap[Int]]
		var wg sync.WaitGroup
		wg.Add(65536)
		for x := Int(0); x < 65536; x++ {
			x := x
			jobs <- func(ps PrioritySource) {
				defer wg.Done()
				for {
					node := n.Load()
					newNode, _ := node.Upsert(Int(x), ps(), false)
					if n.CompareAndSwap(node, newNode) {
						break
					}
				}
			}
		}
		wg.Wait()
		if n.Load().Length() != 65536 {
			b.Fatal()
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
		Build(ps, slice)
	}
}

func BenchmarkBuildUnion65536(b *testing.B) {
	ps := NewPrioritySource()
	var slice []Int
	for i := 0; i < 65536; i++ {
		slice = append(slice, Int(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Build(ps, slice[:len(slice)/2]).Union(
			Build(ps, slice[len(slice)/2:]),
			false,
		)
	}
}

func BenchmarkMutateUpsert65536(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 65536; i++ {
			n, _ = n.Upsert(Int(i), ps(), true)
		}
	}
}

func BenchmarkMutateUpsert(b *testing.B) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.Upsert(Int(i), ps(), true)
	}
}

func BenchmarkRemove(b *testing.B) {
	var slice []Int
	for i := 0; i < b.N; i++ {
		slice = append(slice, Int(i))
	}
	ps := NewPrioritySource()
	n := Build(ps, slice)
	s := rand.Perm(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = n.Remove(Int(s[i]), false)
	}
	if n.Length() != 0 {
		b.Fatal()
	}
}
