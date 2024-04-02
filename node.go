package pt

import (
	"fmt"
	"io"
)

type node[T ordered[T]] struct {
	value    T
	priority Priority
	left     *node[T]
	right    *node[T]
}

func (n *node[T]) upsert(value T, priority Priority) (ret *node[T], existed bool) {
	if n != nil {
		return n.upsertSlow(value, priority)
	}
	// new node
	return &node[T]{
		value:    value,
		priority: priority,
	}, false
}

func (n *node[T]) upsertSlow(value T, priority Priority) (ret *node[T], existed bool) {
	switch value.Compare(n.value) {

	case -1:
		left, existed := n.left.upsert(value, priority)
		return join(
			n,
			left,
			n.right,
		), existed

	case 1:
		right, existed := n.right.upsert(value, priority)
		return join(
			n,
			n.left,
			right,
		), existed

	case 0:
		// exists
		if n.priority == priority {
			// same
			return n, true
		}
		return join(
			// new node
			&node[T]{
				value:    n.value,
				priority: priority,
				// setting these fields is not required for correctness, but will save a node allocation in join
				left:  n.left,
				right: n.right,
			},
			n.left,
			n.right,
		), true

	}
	panic("bad Compare result")
}

func join[T ordered[T]](middle, left, right *node[T]) *node[T] {
	if middle.priority == minPriority && left == nil && right == nil {
		// leaf node to be deleted
		return nil
	}

	if (left == nil || middle.priority >= left.priority) &&
		(right == nil || middle.priority >= right.priority) {
		// no rotation
		if middle.left == left && middle.right == right {
			// no change
			return middle
		}
		// new node
		return &node[T]{
			value:    middle.value,
			priority: middle.priority,
			left:     left,
			right:    right,
		}
	}

	if left != nil && left.priority > middle.priority && (right == nil || left.priority > right.priority) {
		// rotate right
		return &node[T]{
			value:    left.value,
			priority: left.priority,
			left:     left.left,
			right: join(
				middle,
				left.right,
				right,
			),
		}
	}

	if right != nil && right.priority > middle.priority && (left == nil || right.priority > left.priority) {
		// rotate left
		return &node[T]{
			value:    right.value,
			priority: right.priority,
			left: join(
				middle,
				left,
				right.left,
			),
			right: right.right,
		}
	}

	panic("impossible")
}

func (n *node[T]) height() int {
	if n == nil {
		return 0
	}
	return 1 + max(n.left.height(), n.right.height())
}

func (n *node[T]) dump(out io.Writer, level int) {
	if n == nil {
		return
	}
	for i := 0; i < level; i++ {
		out.Write([]byte("\t"))
	}
	fmt.Fprintf(out, "value %v, priority %v\n", n.value, n.priority)
	n.left.dump(out, level+1)
	n.right.dump(out, level+1)
}

func (n *node[T]) remove(value T) (ret *node[T], removed bool) {
	return n.upsert(value, minPriority)
}

func (n *node[T]) split(value T) (ret *node[T], existed bool) {
	return n.upsert(value, maxPriority)
}

func (n *node[T]) union(n2 *node[T]) *node[T] {
	if n2 == nil {
		return n
	}
	if n == nil {
		return n2
	}
	if n2.priority > n.priority {
		return n2.union(n)
	}
	n2Split, _ := n2.split(n.value)
	return join(
		n,
		n.left.union(n2Split.left),
		n.right.union(n2Split.right),
	)
}

func (n *node[T]) length() int {
	if n == nil {
		return 0
	}
	return 1 + n.left.length() + n.right.length()
}

func (n *node[T]) get(pivot T) (ret T, ok bool) {
	if n == nil {
		return
	}
	switch pivot.Compare(n.value) {
	case 0:
		return n.value, true
	case -1:
		return n.left.get(pivot)
	case 1:
		return n.right.get(pivot)
	}
	panic("bad Compare result")
}
