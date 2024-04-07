package pt

type Iter[T Ordered[T]] struct {
	root  *Treap[T]
	next  *Treap[T]
	stack []*Treap[T]
}

func (i *Iter[T]) Next() (ret T, ok bool) {
	// push all left nodes to stack
	for i.next != nil {
		i.stack = append(i.stack, i.next)
		i.next = i.next.left
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
	i.next = node.right
	return
}

func (i *Iter[T]) Seek(pivot T) (ret T, ok bool) {
	// push all left nodes to stack
	for i.next != nil {
		switch pivot.Compare(i.next.value) {
		case 0:
			// found
			ret = i.next.value
			ok = true
			i.next = i.next.right
			return
		case 1:
			// go right node
			i.next = i.next.right
		case -1:
			// go left node
			i.stack = append(i.stack, i.next)
			i.next = i.next.left
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
	i.next = node.right
	return
}

func (i *Iter[T]) Rewind() {
	i.next = i.root
	i.stack = i.stack[:0]
}

func (i *Iter[T]) Close() {
	putIter(i)
}

func (n *Treap[T]) NewIter() *Iter[T] {
	iter := getIter[T]()
	iter.root = n
	iter.next = n
	return iter
}
