package pt

import "testing"

func BenchmarkTreapParallelUpsertSeries(b *testing.B) {
	treap := NewTreap[Int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := Int(0); pb.Next(); i++ {
			treap.Upsert(i)
		}
	})
}

func BenchmarkTreapParallelUpsertRandom(b *testing.B) {
	treap := NewTreap[Int]()
	ps := newPrioritySource()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			treap.Upsert(Int(ps()))
		}
	})
}

func BenchmarkTreapParallelBulkUpsertSeries(b *testing.B) {
	treap := NewTreap[Int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := Int(0); pb.Next(); i++ {
			treap.BulkUpsert([]Int{i, i + 1, i + 2, i + 3})
		}
	})
}

func BenchmarkTreapParallelBulkUpsertRandom(b *testing.B) {
	treap := NewTreap[Int]()
	ps := newPrioritySource()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			treap.BulkUpsert([]Int{
				Int(ps()),
				Int(ps()),
				Int(ps()),
				Int(ps()),
				Int(ps()),
				Int(ps()),
				Int(ps()),
				Int(ps()),
			})
		}
	})
}
