package pt

import (
	"io"
	"math"
	"math/rand/v2"
	"testing"
)

func (n *Treap[T]) checkHeap() {
	if n == nil {
		return
	}
	if n.left != nil && n.left.priority > n.priority {
		panic("bad heap")
	}
	if n.right != nil && n.right.priority > n.priority {
		panic("bad heap")
	}
	n.left.checkHeap()
	n.right.checkHeap()
}

func TestTreap(t *testing.T) {
	ps := NewPrioritySource()
	const num = 8192
	var n *Treap[Int]

	// upsert
	for _, i := range rand.Perm(num) {
		priority := ps()
		existed := false
		n, existed = n.Upsert(Int(i), priority, false)
		if existed {
			t.Fatal()
		}
		n, existed = n.Upsert(Int(i), priority, false)
		if !existed {
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
		// get
		k, ok := n.Get(j)
		if !ok {
			t.Fatal()
		}
		if k != j {
			t.Fatal()
		}
	}

	// height
	if h := n.Height(); h > int(math.Log2(float64(num))*4) {
		// bad luck or bad implementation
		t.Fatalf("got %v", n.Height())
	} else {
		pt("num %v, height %v\n", num, h)
	}

	// split
	for i := 0; i < num; i++ {
		split, existed := n.Split(Int(i), false)
		if !existed {
			t.Fatal()
		}
		if l := split.Length(); l != num {
			t.Fatalf("got %v, expected %v", l, num)
		}
		if l := split.left.Length(); l != i {
			t.Fatal()
		}
		if l := split.right.Length(); l != num-i-1 {
			t.Fatalf("got %v, expected %v", l, num-i-1)
		}
	}

	// union
	n = n.Union(n, false)
	if n.Length() != num {
		t.Fatal()
	}
	for i := 0; i < num; i++ {
		u := n.Union(&Treap[Int]{
			value:    Int(i),
			priority: ps(),
		}, false)
		if u.Length() != num {
			t.Fatal()
		}
	}

	// remove
	for _, i := range rand.Perm(num) {
		removed := false
		n, removed = n.Remove(Int(i), false)
		if !removed {
			t.Fatal()
		}
	}
	if n.Height() != 0 {
		t.Fatal()
	}
}

func TestUpsertPersistence(t *testing.T) {
	ps := NewPrioritySource()
	const num = 1024
	var treaps []*Treap[Int]
	var n *Treap[Int]
	for i := Int(0); i < num; i++ {
		n, _ = n.Upsert(i, ps(), false)
		treaps = append(treaps, n)
	}
	for i, n := range treaps {
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

func TestBuild(t *testing.T) {
	ps := NewPrioritySource()
	var slice []Int
	for i := 0; i < 65536; i++ {
		slice = append(slice, Int(i))
	}
	n := Build(ps, slice)
	if n.Length() != 65536 {
		t.Fatal()
	}
	n.checkHeap()
}

func TestDump(t *testing.T) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < 8; i++ {
		n, _ = n.Upsert(Int(i), ps(), false)
	}
	n.Dump(io.Discard, 0)
}

func TestGetFromNil(t *testing.T) {
	var n *Treap[Int]
	_, ok := n.Get(42)
	if ok {
		t.Fatal()
	}
}

func TestMutateUpsert(t *testing.T) {
	ps := NewPrioritySource()
	var n *Treap[Int]
	for i := 0; i < 4096; i++ {
		n, _ = n.Upsert(Int(i), ps(), true)
		if n.Length() != i+1 {
			t.Fatal()
		}
		j, ok := n.Get(Int(i))
		if !ok {
			t.Fatal()
		}
		if j != Int(i) {
			t.Fatal()
		}
	}
}
