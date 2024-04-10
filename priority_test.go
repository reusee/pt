package pt

import "testing"

func BenchmarkPrioritySource(b *testing.B) {
	ps := newPrioritySource()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps()
	}
}

func BenchmarkParallelPrioritySource(b *testing.B) {
	ps := newPrioritySource()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ps()
		}
	})
}
