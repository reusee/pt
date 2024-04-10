package pt

import "testing"

func BenchmarkTreapParallelUpsert(b *testing.B) {
	treap := NewTreap[Int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := Int(0); pb.Next(); i++ {
			treap.Upsert(i)
		}
	})
}

func BenchmarkTreapParallelBulkUpsert(b *testing.B) {
	treap := NewTreap[Int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := Int(0); pb.Next(); i++ {
			treap.BulkUpsert([]Int{i, i + 1, i + 2, i + 3})
		}
	})
}
