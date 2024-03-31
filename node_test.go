package pt

import (
	"math/rand"
	"testing"
)

func TestInsert(t *testing.T) {
	const num = 65536
	var n *node[Int]
	for _, i := range rand.Perm(num) {
		n = n.insert(Int(i), rand.Int63())
	}

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
	}

}

func BenchmarkInsert(b *testing.B) {
	var n *node[Int]
	for i := 0; i < b.N; i++ {
		n = n.insert(Int(i), rand.Int63())
	}
}
