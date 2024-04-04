package pt

type Iter[T Ordered[T]] struct {
	current *Treap[T]
	stack   []*Treap[T]
}

func (i *Iter[T]) Next() (ret T, ok bool) {
	for i.current != nil {
		i.stack = append(i.stack, i.current)
		i.current = i.current.left
	}
	if len(i.stack) == 0 {
		return
	}
	node := i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	ret = node.value
	ok = true
	i.current = node.right
	return
}

func (i *Iter[T]) Close() {
	putIter(i)
}

func (n *Treap[T]) NewIter() *Iter[T] {
	iter := getIter[T]()
	iter.current = n
	return iter
}
