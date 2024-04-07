package pt

type Iter[T Ordered[T]] struct {
	root    *Treap[T]
	current *Treap[T]
	stack   []*Treap[T]
}

func (i *Iter[T]) Next() (ret T, ok bool) {
	// push all left nodes to stack
	for i.current != nil {
		i.stack = append(i.stack, i.current)
		i.current = i.current.left
	}
	// no more
	if len(i.stack) == 0 {
		return
	}
	// pop
	node := i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	// read
	ret = node.value
	ok = true
	// read right node next
	i.current = node.right
	return
}

func (i *Iter[T]) Seek(pivot T) (ret T, ok bool) {
	// push all left nodes to stack
	for i.current != nil {
		switch pivot.Compare(i.current.value) {
		case 0:
			// found
			ret = i.current.value
			ok = true
			i.current = i.current.right
			return
		case 1:
			// go right node
			i.current = i.current.right
		case -1:
			// go left node
			i.stack = append(i.stack, i.current)
			i.current = i.current.left
		}
	}
	// no more
	if len(i.stack) == 0 {
		return
	}
	// pop
	node := i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	// read
	ret = node.value
	ok = true
	// read right node next
	i.current = node.right
	return
}

func (i *Iter[T]) Rewind() {
	i.current = i.root
	i.stack = i.stack[:0]
}

func (i *Iter[T]) Close() {
	putIter(i)
}

func (n *Treap[T]) NewIter() *Iter[T] {
	iter := getIter[T]()
	iter.root = n
	iter.current = n
	return iter
}
