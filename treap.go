package pt

type Treap[T ordered[T]] struct {
	root *node[T]
}

type ordered[T any] interface {
	Compare(T) int
}
