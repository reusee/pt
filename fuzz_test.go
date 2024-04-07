package pt

import (
	"math/rand/v2"
	"testing"
)

func FuzzTreap(f *testing.F) {
	f.Fuzz(func(t *testing.T, i uint64) {
		r := rand.New(rand.NewPCG(i, i))
		var node *Treap[Int]
		for i := 0; i < 2048; i++ {
			v := Int(r.Int64())
			node, _ = node.Upsert(
				v,
				r.Int64()-1,
				false,
			)
			_, ok := node.Get(v)
			if !ok {
				t.Fatal()
			}
			iter := node.NewIter(nil)
			value, ok := iter.Seek(v)
			if !ok {
				t.Fatal()
			}
			if value != v {
				t.Fatal()
			}
			iter.Close()
		}
		if node.Length() != 2048 {
			t.Fatal()
		}
	})
}
