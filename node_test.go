package pt

import (
	"io"
	"math"
	"math/rand/v2"
	"testing"
)

func TestNode(t *testing.T) {
	ps := NewPrioritySource()
	const num = 4096
	var n *node[Int]

	// upsert
	for _, i := range rand.Perm(num) {
		priority := ps()
		existed := false
		n, existed = n.upsert(Int(i), priority)
		if existed {
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
		split, existed := n.split(Int(i))
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
	n = n.union(n)
	if n.length() != num {
		t.Fatal()
	}
	for i := 0; i < num; i++ {
		split, _ := n.split(Int(i))
		u := n.union(split)
		if u.length() != num {
			t.Fatal()
		}
		u = n.union(&node[Int]{
			value:    Int(i),
			priority: ps(),
		})
		if u.length() != num {
			t.Fatal()
		}
	}

	// remove
	for _, i := range rand.Perm(num) {
		removed := false
		n, removed = n.remove(Int(i))
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
		n, _ = n.upsert(i, ps())
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
		n, _ = n.upsert(Int(i), ps())
	}
	n.dump(io.Discard, 0)
}
