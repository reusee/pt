package pt

import (
	"io"
	"math"
	"math/rand/v2"
	"testing"
)

func (n *node[T]) checkHeap() {
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

func TestNode(t *testing.T) {
	ps := NewPrioritySource()
	const num = 8192
	var n *node[Int]

	// upsert
	for _, i := range rand.Perm(num) {
		priority := ps()
		existed := false
		n, existed = n.upsert(Int(i), priority, false)
		if existed {
			t.Fatal()
		}
		n, existed = n.upsert(Int(i), priority, false)
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
		k, ok := n.get(j)
		if !ok {
			t.Fatal()
		}
		if k != j {
			t.Fatal()
		}
	}

	// height
	if h := n.height(); h > int(math.Log2(float64(num))*4) {
		// bad luck or bad implementation
		t.Fatalf("got %v", n.height())
	} else {
		pt("num %v, height %v\n", num, h)
	}

	// split
	for i := 0; i < num; i++ {
		split, existed := n.split(Int(i), false)
		if !existed {
			t.Fatal()
		}
		if l := split.length(); l != num {
			t.Fatalf("got %v, expected %v", l, num)
		}
		if l := split.left.length(); l != i {
			t.Fatal()
		}
		if l := split.right.length(); l != num-i-1 {
			t.Fatalf("got %v, expected %v", l, num-i-1)
		}
	}

	// union
	n = n.union(n, false)
	if n.length() != num {
		t.Fatal()
	}
	for i := 0; i < num; i++ {
		u := n.union(&node[Int]{
			value:    Int(i),
			priority: ps(),
		}, false)
		if u.length() != num {
			t.Fatal()
		}
	}

	// remove
	for _, i := range rand.Perm(num) {
		removed := false
		n, removed = n.remove(Int(i), false)
		if !removed {
			t.Fatal()
		}
	}
	if n.height() != 0 {
		t.Fatal()
	}
}

func TestUpsertPersistence(t *testing.T) {
	ps := NewPrioritySource()
	const num = 1024
	var nodes []*node[Int]
	var n *node[Int]
	for i := Int(0); i < num; i++ {
		n, _ = n.upsert(i, ps(), false)
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

func TestBuild(t *testing.T) {
	ps := NewPrioritySource()
	var slice []Int
	for i := 0; i < 65536; i++ {
		slice = append(slice, Int(i))
	}
	n := build(ps, slice)
	if n.length() != 65536 {
		t.Fatal()
	}
	n.checkHeap()
}

func TestDump(t *testing.T) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < 8; i++ {
		n, _ = n.upsert(Int(i), ps(), false)
	}
	n.dump(io.Discard, 0)
}

func TestGetFromNil(t *testing.T) {
	var n *node[Int]
	_, ok := n.get(42)
	if ok {
		t.Fatal()
	}
}

func TestMutateUpsert(t *testing.T) {
	ps := NewPrioritySource()
	var n *node[Int]
	for i := 0; i < 4096; i++ {
		n, _ = n.upsert(Int(i), ps(), true)
		if n.length() != i+1 {
			t.Fatal()
		}
		j, ok := n.get(Int(i))
		if !ok {
			t.Fatal()
		}
		if j != Int(i) {
			t.Fatal()
		}
	}
}
