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

type ordered[T any] interface {
	Compare(T) int
}

func (n *node[T]) upsert(value T, priority Priority, mutate bool) (ret *node[T], existed bool) {
	if n != nil {
		return n.upsertSlow(value, priority, mutate)
	}
	// new node
	return &node[T]{
		value:    value,
		priority: priority,
	}, false
}

func (n *node[T]) upsertSlow(value T, priority Priority, mutate bool) (ret *node[T], existed bool) {
	switch value.Compare(n.value) {

	case -1:
		left, existed := n.left.upsert(value, priority, mutate)
		return join(
			n,
			left,
			n.right,
			mutate,
		), existed

	case 1:
		right, existed := n.right.upsert(value, priority, mutate)
		return join(
			n,
			n.left,
			right,
			mutate,
		), existed

	case 0:
		// exists
		if n.priority == priority {
			// same
			return n, true
		}
		if mutate {
			n.priority = priority
			return join(
				n,
				n.left,
				n.right,
				true,
			), true
		} else {
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
				false,
			), true
		}

	}
	panic("bad Compare result") // NOCOVER
}

func join[T ordered[T]](middle, left, right *node[T], mutate bool) *node[T] {
	if middle.priority == MinPriority && left == nil && right == nil {
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
		if mutate {
			// use middle
			middle.left = left
			middle.right = right
			return middle
		} else {
			// new node
			return &node[T]{
				value:    middle.value,
				priority: middle.priority,
				left:     left,
				right:    right,
			}
		}
	}

	if left != nil && left.priority > middle.priority && (right == nil || left.priority >= right.priority) {
		// rotate right
		if mutate {
			// use left
			left.right = join(
				middle,
				left.right,
				right,
				true,
			)
			return left
		} else {
			// new node
			return &node[T]{
				value:    left.value,
				priority: left.priority,
				left:     left.left,
				right: join(
					middle,
					left.right,
					right,
					false,
				),
			}
		}
	}

	if right != nil && right.priority > middle.priority && (left == nil || right.priority >= left.priority) {
		// rotate left
		if mutate {
			// use right
			right.left = join(
				middle,
				left,
				right.left,
				true,
			)
			return right
		} else {
			// new node
			return &node[T]{
				value:    right.value,
				priority: right.priority,
				left: join(
					middle,
					left,
					right.left,
					false,
				),
				right: right.right,
			}
		}
	}

	panic("impossible") // NOCOVER
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

func (n *node[T]) remove(value T, mutate bool) (ret *node[T], removed bool) {
	return n.upsert(value, MinPriority, mutate)
}

func (n *node[T]) split(value T, mutate bool) (ret *node[T], existed bool) {
	return n.upsert(value, MaxPriority, mutate)
}

func (n *node[T]) union(n2 *node[T], mutate bool) *node[T] {
	if n2 == nil {
		return n
	}
	if n == nil {
		return n2
	}
	if n2.priority > n.priority {
		return n2.union(n, mutate)
	}
	n2Split, _ := n2.split(n.value, mutate)
	return join(
		n,
		n.left.union(n2Split.left, mutate),
		n.right.union(n2Split.right, mutate),
		mutate,
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
	panic("bad Compare result") // NOCOVER
}

func build[T ordered[T]](source PrioritySource, slice []T) *node[T] {
	if len(slice) == 0 {
		return nil
	}
	return buildSlow(source, slice)
}

func buildSlow[T ordered[T]](source PrioritySource, slice []T) *node[T] {
	i := len(slice) / 2
	left := slice[:i]
	right := slice[i+1:]
	ret := &node[T]{
		value:    slice[i],
		priority: source(),
		left:     build(source, left),
		right:    build(source, right),
	}
	heapify(ret)
	return ret
}

func heapify[T ordered[T]](n *node[T]) {
	max := n
	if n.left != nil && n.left.priority > max.priority {
		max = n.left
	}
	if n.right != nil && n.right.priority > max.priority {
		max = n.right
	}
	if max != n {
		max.priority, n.priority = n.priority, max.priority
		heapify(max)
	}
}
