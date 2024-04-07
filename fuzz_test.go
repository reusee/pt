package pt

import (
	"math/rand/v2"
	"testing"
)

func FuzzTreap(f *testing.F) {
	f.Fuzz(func(t *testing.T, i uint64) {
		r := rand.New(rand.NewPCG(i, i))
		data := r.Perm(4096)

		var node *Treap[Int]
		for _, v := range data {
			// upsert
			node, ok := node.Upsert(
				Int(v),
				r.Int64()-1,
				false,
			)
			if ok {
				t.Fatal()
			}

			// upsert again
			node, ok = node.Upsert(
				Int(v),
				r.Int64()-1,
				false,
			)
			if !ok {
				t.Fatal()
			}

			// get
			_, ok = node.Get(Int(v))
			if !ok {
				t.Fatal()
			}

			// iter
			iter := node.NewIter(nil)
			// seek
			value, ok := iter.Seek(Int(v))
			if !ok {
				t.Fatal()
			}
			if value != Int(v) {
				t.Fatal()
			}
			iter.Close()
		}

		// length
		if node.Length() != len(data) {
			t.Fatalf("got %v, expected %v", node.Length(), len(data))
		}

		// iter all
		var iterResult []int
		iter := node.NewIter(nil)
		for n, ok := iter.Next(); ok; n, ok = iter.Next() {
			iterResult = append(iterResult, int(n))
		}
		if len(iterResult) != len(data) {
			t.Fatal()
		}
		for i, v := range iterResult {
			if i > 0 {
				if v != iterResult[i] {
					t.Fatal()
				}
			}
		}

	})
}
