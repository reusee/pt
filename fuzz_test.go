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
		for i, v := range data {
			// upsert
			var existed bool
			node, existed = node.Upsert(
				Int(v),
				r.Int64()-1,
				false,
			)
			if existed {
				t.Fatal()
			}
			if node.Length() != i+1 {
				t.Fatalf("got %v, expected %v", node.Length(), i+1)
			}

			// upsert again
			node, existed = node.Upsert(
				Int(v),
				r.Int64()-1,
				false,
			)
			if !existed {
				t.Fatal()
			}
			if node.Length() != i+1 {
				t.Fatal()
			}

			// get
			_, existed = node.Get(Int(v))
			if !existed {
				t.Fatal()
			}

			// iter
			iter := node.NewIter()
			// seek
			value, existed := iter.Seek(Int(v))
			if !existed {
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
		iter := node.NewIter()
		for n, ok := iter.Next(); ok; n, ok = iter.Next() {
			iterResult = append(iterResult, int(n))
		}
		iter.Close()
		if len(iterResult) != len(data) {
			t.Fatal()
		}
		// asc
		for i, v := range iterResult {
			if i > 0 {
				if v != iterResult[i-1]+1 {
					t.Fatal()
				}
			}
		}

	})
}
