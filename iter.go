package pt

type Iter[T any] func() (T, Iter[T])

func (n *node[T]) iter(cont Iter[T]) Iter[T] {
	if n == nil {
		return cont
	}
	return n.left.iter(func() (T, Iter[T]) {
		return n.value, n.right.iter(cont)
	})
}
